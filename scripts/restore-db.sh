#!/bin/bash

# ========================================
# 資料庫還原腳本
# ========================================
# 
# 功能:
# 1. 從備份檔案還原 PostgreSQL 資料庫
# 2. 支援壓縮和未壓縮的備份檔案
#
# 使用方式:
# ./scripts/restore-db.sh <backup_file>
#
# 範例:
# ./scripts/restore-db.sh /home/ubuntu/backups/asset_manager_20241027_120000.sql.gz
#
# ========================================

set -e  # 遇到錯誤立即停止

# 顏色定義
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 檢查參數
if [ $# -eq 0 ]; then
    echo -e "${RED}錯誤: 請提供備份檔案路徑${NC}"
    echo ""
    echo "使用方式:"
    echo "  $0 <backup_file>"
    echo ""
    echo "範例:"
    echo "  $0 /home/ubuntu/backups/asset_manager_20241027_120000.sql.gz"
    echo ""
    echo "可用的備份檔案:"
    ls -lh /home/ubuntu/backups/asset_manager_*.sql.gz 2>/dev/null || echo "  無備份檔案"
    exit 1
fi

BACKUP_FILE="$1"
CONTAINER_NAME="${CONTAINER_NAME:-asset-manager-postgres}"
DB_NAME="${DB_NAME:-asset_manager}"
DB_USER="${DB_USER:-postgres}"

# 檢查備份檔案是否存在
if [ ! -f "$BACKUP_FILE" ]; then
    echo -e "${RED}錯誤: 備份檔案不存在: $BACKUP_FILE${NC}"
    exit 1
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}資料庫還原${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "時間: $(date '+%Y-%m-%d %H:%M:%S')"
echo "容器: $CONTAINER_NAME"
echo "資料庫: $DB_NAME"
echo "備份檔案: $BACKUP_FILE"
echo ""

# 檢查容器是否運行
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo -e "${RED}錯誤: 容器 $CONTAINER_NAME 未運行${NC}"
    exit 1
fi

# 確認還原操作
echo -e "${YELLOW}警告: 此操作會覆蓋現有資料庫!${NC}"
read -p "確定要繼續嗎? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "已取消還原操作"
    exit 0
fi

# 建立當前資料庫的備份 (安全措施)
echo ""
echo -e "${YELLOW}正在建立當前資料庫的安全備份...${NC}"
SAFETY_BACKUP="/tmp/asset_manager_before_restore_$(date +%Y%m%d_%H%M%S).sql"
docker exec "$CONTAINER_NAME" pg_dump -U "$DB_USER" "$DB_NAME" > "$SAFETY_BACKUP"
echo -e "${GREEN}✓ 安全備份已建立: $SAFETY_BACKUP${NC}"

# 解壓縮備份檔案 (如果是 .gz 格式)
RESTORE_FILE="$BACKUP_FILE"
if [[ "$BACKUP_FILE" == *.gz ]]; then
    echo ""
    echo -e "${YELLOW}正在解壓縮備份檔案...${NC}"
    RESTORE_FILE="${BACKUP_FILE%.gz}"
    gunzip -c "$BACKUP_FILE" > "$RESTORE_FILE"
    echo -e "${GREEN}✓ 解壓縮完成${NC}"
fi

# 刪除現有資料庫並重新建立
echo ""
echo -e "${YELLOW}正在重新建立資料庫...${NC}"
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -c "DROP DATABASE IF EXISTS $DB_NAME;"
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;"
echo -e "${GREEN}✓ 資料庫已重新建立${NC}"

# 還原資料庫
echo ""
echo -e "${YELLOW}正在還原資料庫...${NC}"
if docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" "$DB_NAME" < "$RESTORE_FILE"; then
    echo -e "${GREEN}✓ 資料庫還原成功${NC}"
else
    echo -e "${RED}✗ 資料庫還原失敗${NC}"
    echo ""
    echo -e "${YELLOW}正在從安全備份還原...${NC}"
    docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" "$DB_NAME" < "$SAFETY_BACKUP"
    echo -e "${GREEN}✓ 已還原到還原前的狀態${NC}"
    exit 1
fi

# 清理臨時檔案
if [[ "$BACKUP_FILE" == *.gz ]] && [ -f "$RESTORE_FILE" ]; then
    rm "$RESTORE_FILE"
fi

# 驗證還原結果
echo ""
echo -e "${YELLOW}正在驗證還原結果...${NC}"
TABLE_COUNT=$(docker exec "$CONTAINER_NAME" psql -U "$DB_USER" "$DB_NAME" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';")
echo "資料表數量: $TABLE_COUNT"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}還原完成!${NC}"
echo -e "${GREEN}========================================${NC}"
echo "完成時間: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""
echo -e "${YELLOW}提示: 安全備份保存在 $SAFETY_BACKUP${NC}"
echo -e "${YELLOW}如果還原結果正確,可以刪除此檔案${NC}"

