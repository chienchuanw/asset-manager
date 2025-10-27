#!/bin/bash

# ========================================
# 部署腳本 (含自動備份)
# ========================================
# 
# 功能:
# 1. 在部署前自動備份資料庫
# 2. 拉取最新程式碼
# 3. 重新建置並啟動容器
# 4. 執行健康檢查
#
# 使用方式:
# ./scripts/deploy.sh
#
# ========================================

set -e  # 遇到錯誤立即停止

# 顏色定義
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 設定變數
PROJECT_DIR="${PROJECT_DIR:-/home/ubuntu/asset-manager}"
BACKUP_BEFORE_DEPLOY="${BACKUP_BEFORE_DEPLOY:-true}"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Asset Manager 部署腳本${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "開始時間: $(date '+%Y-%m-%d %H:%M:%S')"
echo "專案目錄: $PROJECT_DIR"
echo ""

# 切換到專案目錄
cd "$PROJECT_DIR"

# Step 1: 備份資料庫 (如果容器正在運行)
if [ "$BACKUP_BEFORE_DEPLOY" = "true" ]; then
    if docker ps | grep -q "asset-manager-postgres"; then
        echo -e "${YELLOW}Step 1: 備份資料庫${NC}"
        echo "----------------------------------------"
        
        # 執行備份腳本
        if [ -f "./scripts/backup-db.sh" ]; then
            bash ./scripts/backup-db.sh
        else
            echo -e "${YELLOW}警告: 備份腳本不存在,跳過備份${NC}"
        fi
        echo ""
    else
        echo -e "${YELLOW}Step 1: 跳過備份 (資料庫容器未運行)${NC}"
        echo ""
    fi
else
    echo -e "${YELLOW}Step 1: 跳過備份 (BACKUP_BEFORE_DEPLOY=false)${NC}"
    echo ""
fi

# Step 2: 拉取最新程式碼
echo -e "${YELLOW}Step 2: 拉取最新程式碼${NC}"
echo "----------------------------------------"
git pull origin main
echo -e "${GREEN}✓ 程式碼更新完成${NC}"
echo ""

# Step 3: 停止現有容器
echo -e "${YELLOW}Step 3: 停止現有容器${NC}"
echo "----------------------------------------"
docker-compose down
echo -e "${GREEN}✓ 容器已停止${NC}"
echo ""

# Step 4: 建置新映像檔
echo -e "${YELLOW}Step 4: 建置 Docker 映像檔${NC}"
echo "----------------------------------------"
docker-compose build --no-cache
echo -e "${GREEN}✓ 映像檔建置完成${NC}"
echo ""

# Step 5: 啟動容器
echo -e "${YELLOW}Step 5: 啟動容器${NC}"
echo "----------------------------------------"
docker-compose up -d
echo -e "${GREEN}✓ 容器已啟動${NC}"
echo ""

# Step 6: 等待服務啟動
echo -e "${YELLOW}Step 6: 等待服務啟動${NC}"
echo "----------------------------------------"
echo "等待 30 秒讓服務完全啟動..."
sleep 30
echo -e "${GREEN}✓ 等待完成${NC}"
echo ""

# Step 7: 健康檢查
echo -e "${YELLOW}Step 7: 健康檢查${NC}"
echo "----------------------------------------"

# 檢查容器狀態
echo "檢查容器狀態..."
docker-compose ps

echo ""
echo "檢查服務健康狀態..."

# 檢查 Backend API
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Backend API 正常運行${NC}"
else
    echo -e "${RED}✗ Backend API 健康檢查失敗${NC}"
    echo "查看 Backend 日誌:"
    docker-compose logs --tail=50 backend
    exit 1
fi

# 檢查 Frontend
if curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Frontend 正常運行${NC}"
else
    echo -e "${RED}✗ Frontend 健康檢查失敗${NC}"
    echo "查看 Frontend 日誌:"
    docker-compose logs --tail=50 frontend
    exit 1
fi

# 檢查 Nginx
if curl -f http://localhost/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Nginx 正常運行${NC}"
else
    echo -e "${RED}✗ Nginx 健康檢查失敗${NC}"
    echo "查看 Nginx 日誌:"
    docker-compose logs --tail=50 nginx
    exit 1
fi

echo ""

# Step 8: 清理舊映像檔
echo -e "${YELLOW}Step 8: 清理舊映像檔${NC}"
echo "----------------------------------------"
docker image prune -f
echo -e "${GREEN}✓ 清理完成${NC}"
echo ""

# 完成
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}部署完成!${NC}"
echo -e "${BLUE}========================================${NC}"
echo "完成時間: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""
echo "服務狀態:"
echo "  Backend API: http://localhost:8080"
echo "  Frontend: http://localhost:3000"
echo "  Nginx: http://localhost"
echo ""
echo "查看日誌:"
echo "  docker-compose logs -f [service_name]"
echo ""

