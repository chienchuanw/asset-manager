# Container Security Hardening

This document records the runtime security hardening applied to the Docker Compose stack after the 2026-05-08 Next.js MCP-server RCE incident (PRs #79, #80). The goal is to reduce the blast radius if any single container is compromised again.

## Applied controls (docker-compose.yml)

The `backend` and `frontend` services run with:

| Control | Setting | Why |
|---|---|---|
| Read-only root filesystem | `read_only: true` | Attacker cannot drop a payload anywhere outside explicitly mounted tmpfs. |
| Hardened `/tmp` | `tmpfs: /tmp` with `rw,noexec,nosuid,size=64m` | The XMRig dropper (`javae`, `kinsing`) needed `/tmp` writable AND executable. `noexec` blocks the launch step entirely. |
| Frontend Next.js cache | `tmpfs: /app/.next/cache` with `rw,nosuid,size=128m` | Next.js standalone needs `.next/cache` writable for the image optimization cache. Kept separate so we don't need to make rootfs writable. |
| Drop all capabilities | `cap_drop: [ALL]` | Neither service needs Linux capabilities at runtime. |
| No privilege escalation | `security_opt: [no-new-privileges:true]` | Blocks setuid binaries from elevating. |
| Process count limit | `deploy.resources.limits.pids: 200` | Caps fork-bombs and runaway miners. |

The `postgres`, `redis`, `nginx`, and `certbot` services are unchanged for now — they need writable state directories (`/var/lib/postgresql/data`, `/data`, `/var/cache/nginx`, `/etc/letsencrypt`) and tightening them is a separate exercise.

## Verification

After `docker compose up -d`:

```bash
# Health check passes
make health

# /app is read-only inside the frontend container
docker compose exec frontend sh -c 'touch /app/x' # expect: Read-only file system

# /tmp is writable but non-executable inside the frontend container
docker compose exec frontend sh -c 'echo \#!/bin/sh > /tmp/t && chmod +x /tmp/t && /tmp/t' # expect: Permission denied

# Same checks for backend
docker compose exec backend sh -c 'touch /x' # expect: Read-only file system
docker compose exec backend sh -c 'echo data > /tmp/t && chmod +x /tmp/t && /tmp/t' # expect: Permission denied
```

## Why these specific values

- `tmpfs size=64m` for `/tmp`: enough for legitimate temp files (sockets, lock files) but too small for the ~5–10 MB miner binaries seen in the incident plus their config and logs.
- `pids: 200`: comfortably above normal worker counts (Next.js + Node worker threads, Gin goroutines stay as one OS process), but well below the levels a coin-miner with worker fork-out would need.
- `/app/.next/cache size=128m`: Next.js image cache for a small SPA stays well under this. If you start serving large image catalogs and see eviction churn in logs, raise it.

## Related work

- Issue #82 — nginx attack surface reduction (proxy paths, rate limiting, SSRF blocks)
- Issue #83 — dependency CVE monitoring + base image digest pinning
- Issue #85 — intrusion and availability alerting
