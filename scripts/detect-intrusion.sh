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

# 硬依賴：jq 用來組 JSON、curl 用來送 webhook。少任何一個就讓 cron 看到 stderr 而非靜默失敗。
for dep in jq curl docker; do
  if ! command -v "$dep" >/dev/null 2>&1; then
    echo "ERROR: required command '$dep' not found in PATH" >&2
    exit 1
  fi
done

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
    # printf -v 把實體換行寫進變數；單純的 "...\n..." 在 bash 雙引號中是字面 backslash-n
    local entry
    printf -v entry '**[%s] 已知挖礦/惡意程式名命中**\n```\n%s\n```' "$c" "$matches"
    ALERTS+=("$entry")
  fi
}

check_tmp_executables() {
  local c="$1"
  # 容器掛 read-only + tmpfs noexec 後此項應永遠為空；但 nginx/postgres/redis 等未強化容器仍適用。
  local exe_files
  exe_files=$(docker exec "$c" sh -c 'find /tmp -type f -perm -u+x 2>/dev/null | head -5' 2>/dev/null || true)
  if [[ -n "$exe_files" ]]; then
    local entry
    printf -v entry '**[%s] /tmp 出現可執行檔**\n```\n%s\n```' "$c" "$exe_files"
    ALERTS+=("$entry")
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

# 拼 payload；Discord content 欄位上限 2000 字元，多重告警同時觸發會超過。
header="🚨 asset-manager 入侵偵測告警 (${HOSTNAME_S})"
joined=$(printf '%s\n\n' "${ALERTS[@]}")
# 預留 ~150 字給 header + 截斷提示
max_body=$((2000 - ${#header} - 150))
if [[ ${#joined} -gt $max_body ]]; then
  joined="${joined:0:$max_body}

…（訊息過長已截斷，共 ${#ALERTS[@]} 項，請登入主機查看完整輸出）"
fi
content=$(jq -Rn --arg c "$header" --arg b "$joined" '{content: ($c + "\n\n" + $b)}')

# 把 curl 的 HTTP code 抓出來，非 2xx 時讓 cron 看到 stderr
http_code=$(curl -sS -o /tmp/discord-resp.$$ -w '%{http_code}' \
  -X POST "$DISCORD_WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -d "$content")
if [[ "$http_code" != 2* ]]; then
  echo "ERROR: Discord webhook returned HTTP $http_code" >&2
  cat /tmp/discord-resp.$$ >&2
  rm -f /tmp/discord-resp.$$
  exit 2
fi
rm -f /tmp/discord-resp.$$

echo "$(date -u +%FT%TZ) alerted: ${#ALERTS[@]} item(s)"
