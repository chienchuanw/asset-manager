#!/bin/bash

# ========================================
# 上傳備份到 S3 腳本
# ========================================
# 
# 功能:
# 1. 備份資料庫
# 2. 上傳到 S3
# 3. 設定生命週期 (可選)
#
# 使用方式:
# export AWS_S3_BUCKET=asset-manager-backups
# export AWS_REGION=ap-northeast-1
# ./scripts/backup-to-s3.sh
#
# ========================================

set -e

# 顏色定義
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# 檢查 AWS CLI
if ! command -v aws &> /dev/null; then
    echo -e "${RED}錯誤: AWS CLI 未安裝${NC}"
    echo ""
    echo "安裝方式:"
    echo "  sudo apt install -y awscli"
    echo "  aws configure"
    exit 1
fi

# 檢查環境變數
if [ -z "$AWS_S3_BUCKET" ]; then
    echo -e "${RED}錯誤: 請設定 AWS_S3_BUCKET 環境變數${NC}"
    exit 1
fi

# 設定變數
BACKUP_DIR="${BACKUP_DIR:-/home/ubuntu/backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/asset_manager_$TIMESTAMP.sql"
CONTAINER_NAME="${CONTAINER_NAME:-asset-manager-postgres}"
DB_NAME="${DB_NAME:-asset_manager}"
DB_USER="${DB_USER:-postgres}"
AWS_REGION="${AWS_REGION:-ap-northeast-1}"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}備份並上傳到 S3${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "S3 Bucket: $AWS_S3_BUCKET"
echo "Region: $AWS_REGION"
echo ""

# 建立備份目錄
mkdir -p "$BACKUP_DIR"

# 執行備份
echo -e "${YELLOW}正在備份資料庫...${NC}"
docker exec "$CONTAINER_NAME" pg_dump -U "$DB_USER" "$DB_NAME" > "$BACKUP_FILE"
echo -e "${GREEN}✓ 資料庫備份成功${NC}"

# 壓縮備份
echo -e "${YELLOW}正在壓縮備份檔案...${NC}"
gzip "$BACKUP_FILE"
COMPRESSED_FILE="${BACKUP_FILE}.gz"
echo -e "${GREEN}✓ 壓縮成功${NC}"

# 上傳到 S3
S3_PATH="s3://$AWS_S3_BUCKET/backups/$(date +%Y/%m)/asset_manager_$TIMESTAMP.sql.gz"

echo -e "${YELLOW}正在上傳到 S3...${NC}"
if aws s3 cp "$COMPRESSED_FILE" "$S3_PATH" --region "$AWS_REGION"; then
    echo -e "${GREEN}✓ 上傳成功: $S3_PATH${NC}"
else
    echo -e "${RED}✗ 上傳失敗${NC}"
    exit 1
fi

# 列出 S3 上的備份
echo ""
echo "S3 上的備份檔案:"
aws s3 ls "s3://$AWS_S3_BUCKET/backups/" --recursive --human-readable --summarize | tail -10

echo ""
echo -e "${GREEN}完成!${NC}"

