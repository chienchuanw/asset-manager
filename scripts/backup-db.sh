#!/bin/bash

# ========================================
# 資料庫備份腳本
# ========================================
# 
# 功能:
# 1. 備份 PostgreSQL 資料庫
# 2. 壓縮備份檔案
# 3. 刪除舊備份 (保留最近 7 天)
# 4. (可選) 上傳到 S3
#
# 使用方式:
# ./scripts/backup-db.sh
#
# ========================================

set -e  # 遇到錯誤立即停止

# 顏色定義
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 設定變數
BACKUP_DIR="${BACKUP_DIR:-/home/ubuntu/backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/asset_manager_$TIMESTAMP.sql"
CONTAINER_NAME="${CONTAINER_NAME:-asset-manager-postgres}"
DB_NAME="${DB_NAME:-asset_manager}"
DB_USER="${DB_USER:-postgres}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"

# 建立備份目錄
mkdir -p "$BACKUP_DIR"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}資料庫備份開始${NC}"
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

# 執行備份
echo -e "${YELLOW}正在備份資料庫...${NC}"
if docker exec "$CONTAINER_NAME" pg_dump -U "$DB_USER" "$DB_NAME" > "$BACKUP_FILE"; then
    echo -e "${GREEN}✓ 資料庫備份成功${NC}"
else
    echo -e "${RED}✗ 資料庫備份失敗${NC}"
    exit 1
fi

# 檢查備份檔案大小
BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
echo "備份檔案大小: $BACKUP_SIZE"

# 壓縮備份檔案
echo -e "${YELLOW}正在壓縮備份檔案...${NC}"
if gzip "$BACKUP_FILE"; then
    echo -e "${GREEN}✓ 壓縮成功${NC}"
    COMPRESSED_FILE="${BACKUP_FILE}.gz"
    COMPRESSED_SIZE=$(du -h "$COMPRESSED_FILE" | cut -f1)
    echo "壓縮後大小: $COMPRESSED_SIZE"
else
    echo -e "${RED}✗ 壓縮失敗${NC}"
    exit 1
fi

# 刪除舊備份
echo -e "${YELLOW}正在清理舊備份 (保留最近 $RETENTION_DAYS 天)...${NC}"
DELETED_COUNT=$(find "$BACKUP_DIR" -name "asset_manager_*.sql.gz" -mtime +$RETENTION_DAYS -delete -print | wc -l)
if [ "$DELETED_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓ 已刪除 $DELETED_COUNT 個舊備份${NC}"
else
    echo "沒有需要刪除的舊備份"
fi

# 列出現有備份
echo ""
echo "現有備份檔案:"
ls -lh "$BACKUP_DIR"/asset_manager_*.sql.gz 2>/dev/null | tail -5 || echo "無備份檔案"

# (可選) 上傳到 S3
if [ -n "$AWS_S3_BUCKET" ]; then
    echo ""
    echo -e "${YELLOW}正在上傳備份到 S3...${NC}"
    
    # 檢查 AWS CLI 是否安裝
    if ! command -v aws &> /dev/null; then
        echo -e "${YELLOW}警告: AWS CLI 未安裝,跳過 S3 上傳${NC}"
    else
        S3_PATH="s3://$AWS_S3_BUCKET/backups/$(date +%Y/%m)/asset_manager_$TIMESTAMP.sql.gz"
        
        if aws s3 cp "$COMPRESSED_FILE" "$S3_PATH"; then
            echo -e "${GREEN}✓ 已上傳到 S3: $S3_PATH${NC}"
        else
            echo -e "${RED}✗ S3 上傳失敗${NC}"
        fi
    fi
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}備份完成!${NC}"
echo -e "${GREEN}========================================${NC}"
echo "備份檔案: $COMPRESSED_FILE"
echo "完成時間: $(date '+%Y-%m-%d %H:%M:%S')"

