#!/usr/bin/env bash
# 入侵與資源異常偵測（2026-05-08 RCE 事件後新增）
#
# 設計目標：以最少的依賴（docker + curl + awk）在主機 cron 跑，
# 偵測下列幾類訊號並透過 Discord webhook 通知：
#   1. 已知挖礦/惡意程式名出現在任何容器內
#   2. /tmp 出現可執行檔（這是上次事件中 XMRig 落地的具體手法）
#   3. 容器記憶體用量長時間貼近 limit
#
# Usage:
#   DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/... ./scripts/detect-intrusion.sh
#
# Cron 範例（每 5 分鐘跑一次）：
#   */5 * * * * DISCORD_WEBHOOK_URL=... /opt/asset-manager/scripts/detect-intrusion.sh >> /var/log/asset-manager-detect.log 2>&1

set -uo pipefail

if [[ -z "${DISCORD_WEBHOOK_URL:-}" ]]; then
  echo "ERROR: DISCORD_WEBHOOK_URL not set" >&2
  exit 1
fi

HOSTNAME_S=$(hostname)
ALERTS=()

# 觀察哪些容器；可由 CONTAINERS 環境變數覆寫
CONTAINERS=${CONTAINERS:-"asset-manager-frontend asset-manager-backend asset-manager-nginx"}

# 已知挖礦/惡意程式名（小寫，會 grep -i）
BAD_PROCS_REGEX='xmrig|kdevtmpfsi|kinsing|javae|c3pool|supportxmr|minerd|cpuminer'

# 記憶體閾值：超過 limit 的此百分比就告警
MEM_THRESHOLD_PCT=${MEM_THRESHOLD_PCT:-90}

check_bad_processes() {
  local c="$1"
  # docker top 失敗（容器不存在/未跑）就跳過
  local procs
  if ! procs=$(docker top "$c" -eo pid,comm,args 2>/dev/null); then
    return
  fi
  if echo "$procs" | grep -iE "$BAD_PROCS_REGEX" >/dev/null; then
    local matches
    matches=$(echo "$procs" | grep -iE "$BAD_PROCS_REGEX" | head -3)
    ALERTS+=("**[$c] 已知挖礦/惡意程式名命中**\n\`\`\`\n$matches\n\`\`\`")
  fi
}

check_tmp_executables() {
  local c="$1"
  # 容器掛 read-only + tmpfs noexec 後此項應永遠為空；但 nginx/postgres/redis 等未強化容器仍適用。
  local exe_files
  exe_files=$(docker exec "$c" sh -c 'find /tmp -type f -perm -u+x 2>/dev/null | head -5' 2>/dev/null || true)
  if [[ -n "$exe_files" ]]; then
    ALERTS+=("**[$c] /tmp 出現可執行檔**\n\`\`\`\n$exe_files\n\`\`\`")
  fi
}

check_memory_pressure() {
  # docker stats: 一次取一個 snapshot
  local stats
  stats=$(docker stats --no-stream --format '{{.Name}} {{.MemPerc}}' 2>/dev/null || true)
  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    local name pct
    name=$(echo "$line" | awk '{print $1}')
    pct=$(echo "$line" | awk '{print $2}' | tr -d '%')
    # 只看我們關心的容器
    if ! [[ " $CONTAINERS " =~ " $name " ]]; then continue; fi
    # bash 不會做小數比較，用 awk
    if awk -v p="$pct" -v t="$MEM_THRESHOLD_PCT" 'BEGIN{exit !(p+0 > t+0)}'; then
      ALERTS+=("**[$name] 記憶體用量 ${pct}% 超過閾值 ${MEM_THRESHOLD_PCT}%**")
    fi
  done <<<"$stats"
}

for c in $CONTAINERS; do
  check_bad_processes "$c"
  check_tmp_executables "$c"
done
check_memory_pressure

if [[ ${#ALERTS[@]} -eq 0 ]]; then
  echo "$(date -u +%FT%TZ) ok"
  exit 0
fi

# 拼 payload
joined=$(printf '%s\n\n' "${ALERTS[@]}")
content=$(jq -Rn --arg c "🚨 asset-manager 入侵偵測告警 (${HOSTNAME_S})" --arg b "$joined" '{content: ($c + "\n\n" + $b)}')

curl -sS -X POST "$DISCORD_WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -d "$content" >/dev/null

echo "$(date -u +%FT%TZ) alerted: ${#ALERTS[@]} item(s)"
