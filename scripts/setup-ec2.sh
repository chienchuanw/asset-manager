#!/bin/bash

# ========================================
# EC2 初始化腳本
# ========================================
# 
# 功能:
# 1. 安裝 Docker 和 Docker Compose
# 2. 設定自動備份 Cron Job
# 3. 設定環境變數
# 4. Clone 專案程式碼
#
# 使用方式:
# 在 EC2 上執行:
# curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/asset-manager/main/scripts/setup-ec2.sh | bash
#
# 或手動執行:
# ./scripts/setup-ec2.sh
#
# ========================================

set -e  # 遇到錯誤立即停止

# 顏色定義
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}EC2 初始化腳本${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Step 1: 更新系統套件
echo -e "${YELLOW}Step 1: 更新系統套件${NC}"
echo "----------------------------------------"
sudo apt update
sudo apt upgrade -y
echo -e "${GREEN}✓ 系統套件更新完成${NC}"
echo ""

# Step 2: 安裝必要工具
echo -e "${YELLOW}Step 2: 安裝必要工具${NC}"
echo "----------------------------------------"
sudo apt install -y curl git vim htop
echo -e "${GREEN}✓ 必要工具安裝完成${NC}"
echo ""

# Step 3: 安裝 Docker
echo -e "${YELLOW}Step 3: 安裝 Docker${NC}"
echo "----------------------------------------"
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker ubuntu
    rm get-docker.sh
    echo -e "${GREEN}✓ Docker 安裝完成${NC}"
else
    echo -e "${GREEN}✓ Docker 已安裝${NC}"
fi
echo ""

# Step 4: 安裝 Docker Compose
echo -e "${YELLOW}Step 4: 安裝 Docker Compose${NC}"
echo "----------------------------------------"
if ! command -v docker-compose &> /dev/null; then
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo -e "${GREEN}✓ Docker Compose 安裝完成${NC}"
else
    echo -e "${GREEN}✓ Docker Compose 已安裝${NC}"
fi
echo ""

# Step 5: 建立專案目錄
echo -e "${YELLOW}Step 5: 建立專案目錄${NC}"
echo "----------------------------------------"
mkdir -p /home/ubuntu/asset-manager
mkdir -p /home/ubuntu/backups
echo -e "${GREEN}✓ 目錄建立完成${NC}"
echo ""

# Step 6: 設定自動備份 Cron Job
echo -e "${YELLOW}Step 6: 設定自動備份 Cron Job${NC}"
echo "----------------------------------------"
CRON_JOB="0 2 * * * /home/ubuntu/asset-manager/scripts/backup-db.sh >> /home/ubuntu/backup.log 2>&1"

# 檢查 cron job 是否已存在
if ! crontab -l 2>/dev/null | grep -q "backup-db.sh"; then
    (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
    echo -e "${GREEN}✓ 自動備份已設定 (每天凌晨 2 點)${NC}"
else
    echo -e "${GREEN}✓ 自動備份已存在${NC}"
fi
echo ""

# Step 7: 設定 Git (可選)
echo -e "${YELLOW}Step 7: 設定 Git${NC}"
echo "----------------------------------------"
read -p "請輸入 Git 使用者名稱 (或按 Enter 跳過): " GIT_NAME
if [ -n "$GIT_NAME" ]; then
    git config --global user.name "$GIT_NAME"
    read -p "請輸入 Git Email: " GIT_EMAIL
    git config --global user.email "$GIT_EMAIL"
    echo -e "${GREEN}✓ Git 設定完成${NC}"
else
    echo "跳過 Git 設定"
fi
echo ""

# 完成
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}EC2 初始化完成!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "接下來的步驟:"
echo ""
echo "1. 登出並重新登入 (讓 Docker 群組設定生效):"
echo "   exit"
echo "   ssh -i ~/.ssh/asset-manager-key.pem ubuntu@YOUR_EC2_IP"
echo ""
echo "2. Clone 專案程式碼:"
echo "   cd /home/ubuntu/asset-manager"
echo "   git clone https://github.com/YOUR_USERNAME/asset-manager.git ."
echo ""
echo "3. 建立環境變數檔案:"
echo "   cp .env.production.example .env.production"
echo "   vim .env.production  # 填入實際的值"
echo ""
echo "4. 啟動服務:"
echo "   docker-compose up -d"
echo ""
echo "5. 查看日誌:"
echo "   docker-compose logs -f"
echo ""

