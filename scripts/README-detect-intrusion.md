# Intrusion & Availability Alerting

Detector + runbook added after the 2026-05-08 Next.js MCP RCE incident, where the site was down for ~4 days before being noticed. The goal is the cheapest possible signal: "something is wrong, look at the host now."

## Components

### 1. Host-side intrusion detector (`scripts/detect-intrusion.sh`)

Runs on the host (not inside any container), watches the containers, and alerts to the existing Discord webhook.

Checks per run:

- **Bad process names** — greps `docker top` output for `xmrig`, `kdevtmpfsi`, `kinsing`, `javae`, `c3pool`, `supportxmr`, `minerd`, `cpuminer`. These are the families seen in the 2026-05-08 incident plus the most common XMRig droppers.
- **/tmp executables** — `find /tmp -type f -perm -u+x` inside each container. Hardened services (frontend, backend after #81) have `noexec` tmpfs so this should always be empty; the check still applies to `nginx`/`redis`/`postgres` until those are tightened too.
- **Memory pressure** — `docker stats --no-stream`; if any watched container exceeds `MEM_THRESHOLD_PCT` (default 90%), alert. The miner's runaway CPU/memory was the first observable signal during the incident.

#### Install

```bash
# Copy script onto the EC2 host (deploy.sh / your provisioning of choice)
sudo install -m 0755 scripts/detect-intrusion.sh /opt/asset-manager/scripts/detect-intrusion.sh

# Install jq if missing
sudo apt-get install -y jq

# Add to root crontab — every 5 minutes
sudo crontab -e
# Append:
*/5 * * * * DISCORD_WEBHOOK_URL='<copy from .env.production>' /opt/asset-manager/scripts/detect-intrusion.sh >> /var/log/asset-manager-detect.log 2>&1
```

#### Tune

- `CONTAINERS` (space-separated): default `asset-manager-frontend asset-manager-backend asset-manager-nginx`. Add postgres/redis if you also want to watch them.
- `MEM_THRESHOLD_PCT`: default `90`. Lower for early warning, higher to reduce noise.

#### Verify (acceptance criteria from issue #85)

```bash
# Fake a memory-pressure alert: stress the frontend until it's near limit
docker exec asset-manager-frontend sh -c 'node -e "let a=[]; setInterval(()=>{a.push(Buffer.alloc(1024*1024))},10)"' &

# Within 5 minutes, expect a Discord message:
#   🚨 asset-manager 入侵偵測告警 (...)
#   **[asset-manager-frontend] 記憶體用量 95.x% 超過閾值 90%**
```

### 2. External uptime monitor

The detector is host-resident; if the host or network is the failure, it can't alert. Pair with an external check:

**Cloudflare Health Checks** (free, recommended since DNS is already on Cloudflare):

1. Cloudflare Dashboard → your zone → Traffic → Health Checks → Create.
2. Type: HTTPS, Path: `/health`, Method: `GET`, Expected codes: `200`, Expected response body: `healthy`.
3. Interval: 60s, Retries: 2, Timeout: 5s. Region: closest to EC2 (e.g. Western US).
4. Notify: Discord webhook (Cloudflare → Notifications → Add → Health Check). Use the same `DISCORD_WEBHOOK_URL`.

Alternative: UptimeRobot free tier (5-min interval, 50 monitors). Configure the same way against `https://asset.chienchuanw.com/health`.

#### Verify

```bash
# Force a 503 by stopping nginx briefly
ssh ec2 'docker stop asset-manager-nginx'
# Wait 2-3 minutes — Cloudflare / UptimeRobot should fire a "DOWN" alert
ssh ec2 'docker start asset-manager-nginx'
# "UP" / recovery notification follows
```

## Runbook — receiving an alert

1. **Confirm it's real.** `curl -I https://asset.chienchuanw.com/health` from your workstation. If 200, the host detector likely fired on memory/CPU; if 5xx, both detectors should be screaming.
2. **SSH to the host.** `docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'` — anything restarting, anything CPU-pegged in `docker stats`?
3. **For the "bad process" alert specifically (highest priority):**
   - `docker top <container>` — identify the process and parent.
   - Capture evidence before killing: `docker exec <c> ps auxf > /tmp/incident-$(date +%s).txt`.
   - Stop the container: `docker stop <c>`. Do **not** `docker rm` — keep the filesystem layer for forensics.
   - Rotate any secrets the container had access to (`.env.production` contents).
   - File a follow-up incident issue referencing the original 2026-05-08 incident (PRs #79, #80).
4. **For uptime / 5xx:** `docker compose logs --tail=200 <service>` — check for OOM kills, panic traces. If the host is fine but containers crashed, check disk: `df -h`.
5. **Document.** Add a one-paragraph note to the project journal so future you / future me can correlate.

## Why this design

- **No agents, no extra daemons.** The 2026-05-08 root cause was an unmaintained dependency; piling on a Datadog agent or similar adds the same kind of attack surface this script is supposed to detect.
- **Two independent layers.** Host detector catches what the external monitor cannot (process names, /tmp drops); external monitor catches what the host detector cannot (host-level outage).
- **Reuses existing Discord webhook.** No new secrets, no new notification channel to forget about.

## Related

- Issue #24 — broader observability (structured logging, metrics). This is the minimum viable subset.
- Issue #81 — container hardening (read-only fs + noexec /tmp). Once merged, the `/tmp executable` check becomes redundant on hardened services but still useful on nginx/postgres/redis.
