# 🚀 Asset Manager 部署指南

## 📋 目錄

1. [前置需求](#前置需求)
2. [EC2 設定](#ec2-設定)
3. [首次部署](#首次部署)
4. [GitHub Actions 設定](#github-actions-設定)
5. [Discord Webhook 設定](#discord-webhook-設定)
6. [自動化備份](#自動化備份)
7. [常見問題](#常見問題)

---

## 前置需求

### AWS 資源

- ✅ EC2 Instance (t3.small, Ubuntu 22.04)
- ✅ SSH Key Pair
- ✅ Security Group (開放 Port 22, 80, 443)
- ✅ Elastic IP (可選,建議使用)

### 本地工具

- Git
- SSH Client
- (可選) AWS CLI

### API Keys

- FinMind API Key (台股資料)
- CoinGecko API Key (加密貨幣資料)
- Alpha Vantage API Key (美股資料)

---

## EC2 設定

### 1. SSH 連線到 EC2

```bash
ssh -i ~/.ssh/asset-manager-key.pem ubuntu@43.213.77.244
```

### 2. 執行初始化腳本

```bash
# 下載並執行初始化腳本
curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/asset-manager/main/scripts/setup-ec2.sh | bash

# 或手動執行
git clone https://github.com/YOUR_USERNAME/asset-manager.git /home/ubuntu/asset-manager
cd /home/ubuntu/asset-manager
bash scripts/setup-ec2.sh
```

### 3. 登出並重新登入

```bash
exit
ssh -i ~/.ssh/asset-manager-key.pem ubuntu@43.213.77.244
```

### 4. 驗證安裝

```bash
# 檢查 Docker
docker --version
docker run hello-world

# 檢查 Docker Compose
docker-compose --version
```

---

## 首次部署

### 1. Clone 專案程式碼

```bash
cd /home/ubuntu/asset-manager
git clone https://github.com/YOUR_USERNAME/asset-manager.git .
```

### 2. 建立環境變數檔案

```bash
# 複製範本
cp .env.production.example .env.production

# 編輯環境變數
vim .env.production
```

**必須填寫的變數:**

```bash
# 資料庫密碼
DB_PASSWORD=YOUR_STRONG_PASSWORD

# 身份驗證
AUTH_USERNAME=admin
AUTH_PASSWORD=YOUR_STRONG_PASSWORD
JWT_SECRET=YOUR_JWT_SECRET  # 使用 openssl rand -base64 32 產生

# API Keys
FINMIND_API_KEY=YOUR_KEY
COINGECKO_API_KEY=YOUR_KEY
ALPHA_VANTAGE_API_KEY=YOUR_KEY

# CORS 和 API URL
CORS_ALLOWED_ORIGINS=http://43.213.77.244
NEXT_PUBLIC_API_URL=http://43.213.77.244/api
```

### 3. 啟動服務

```bash
# 建置並啟動所有容器
docker-compose --env-file .env.production up -d

# 查看容器狀態
docker-compose ps

# 查看日誌
docker-compose logs -f
```

### 4. 驗證部署

```bash
# 檢查 Backend API
curl http://localhost:8080/health

# 檢查 Frontend
curl http://localhost:3000

# 檢查 Nginx
curl http://localhost/health
```

### 5. 從外部訪問

在瀏覽器開啟:

- Frontend: `http://43.213.77.244`
- Backend API: `http://43.213.77.244/api/health`

---

## GitHub Actions 設定

### 1. 設定 GitHub Secrets

前往 GitHub Repository → Settings → Secrets and variables → Actions

新增以下 Secrets:

**EC2 連線資訊:**

```
EC2_HOST=43.213.77.244
EC2_USERNAME=ubuntu
EC2_SSH_KEY=<貼上 ~/.ssh/asset-manager-key.pem 的內容>
```

**資料庫設定:**

```
PROD_DB_USER=postgres
PROD_DB_PASSWORD=YOUR_STRONG_PASSWORD
PROD_DB_NAME=asset_manager
```

**Redis 設定:**

```
PROD_REDIS_PASSWORD=  (留空或設定密碼)
```

**身份驗證:**

```
PROD_AUTH_USERNAME=admin
PROD_AUTH_PASSWORD=YOUR_STRONG_PASSWORD
PROD_JWT_SECRET=YOUR_JWT_SECRET
```

**API Keys:**

```
PROD_FINMIND_API_KEY=YOUR_KEY
PROD_COINGECKO_API_KEY=YOUR_KEY
PROD_ALPHA_VANTAGE_API_KEY=YOUR_KEY
```

**Discord Webhook:**

```
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
```

**Discord Bot:**

```
DISCORD_BOT_ENABLED=true
DISCORD_BOT_TOKEN=MTIz...
DISCORD_CHANNEL_IDS=123456789,987654321
DISCORD_BOT_LANG=zh-TW
GEMINI_API_KEY=AIza...
```

### 2. 測試自動部署

```bash
# 在本地修改程式碼
git add .
git commit -m "test: 測試自動部署"
git push origin main

# 前往 GitHub Actions 查看部署進度
# https://github.com/YOUR_USERNAME/asset-manager/actions
```

---

## Discord Webhook 設定

### 1. 建立 Discord Webhook

1. 開啟 Discord,選擇一個頻道
2. 點擊頻道設定 (齒輪圖示)
3. 選擇「整合」→「Webhooks」
4. 點擊「新增 Webhook」
5. 設定名稱 (例如: Asset Manager Deploy)
6. 複製 Webhook URL
7. 貼到 GitHub Secrets 的 `DISCORD_WEBHOOK_URL`

### 2. 測試 Webhook

```bash
curl -H "Content-Type: application/json" \
  -d '{
    "embeds": [{
      "title": "測試通知",
      "description": "Discord Webhook 設定成功!",
      "color": 3066993
    }]
  }' \
  YOUR_DISCORD_WEBHOOK_URL
```

---

## Discord Bot 環境變數

| 變數 | 說明 | 範例 |
|------|------|------|
| `DISCORD_BOT_ENABLED` | 啟用 Discord Bot | `true` |
| `DISCORD_BOT_TOKEN` | Bot Token | `MTIz...` |
| `DISCORD_CHANNEL_IDS` | 監聽頻道 ID | `123456789,987654321` |
| `DISCORD_BOT_LANG` | Bot 語言 | `zh-TW` |
| `GEMINI_API_KEY` | Gemini API Key | `AIza...` |

---

## 自動化備份

### 1. 自動備份設定

自動備份已在 `setup-ec2.sh` 中設定,每天凌晨 2 點執行。

查看 Cron Job:

```bash
crontab -l
```

### 2. 手動備份

```bash
cd /home/ubuntu/asset-manager
bash scripts/backup-db.sh
```

### 3. 查看備份檔案

```bash
ls -lh /home/ubuntu/backups/
```

### 4. 還原備份

```bash
cd /home/ubuntu/asset-manager
bash scripts/restore-db.sh /home/ubuntu/backups/asset_manager_YYYYMMDD_HHMMSS.sql.gz
```

### 5. 上傳備份到 S3 (可選)

**安裝 AWS CLI:**

```bash
sudo apt install -y awscli
aws configure
```

**設定環境變數:**

```bash
export AWS_S3_BUCKET=asset-manager-backups
export AWS_REGION=ap-northeast-1
```

**執行備份 (會自動上傳到 S3):**

```bash
bash scripts/backup-db.sh
```

---

## 常見問題

### Q1: 容器無法啟動

**檢查日誌:**

```bash
docker-compose logs backend
docker-compose logs frontend
docker-compose logs postgres
```

**常見原因:**

- 環境變數設定錯誤
- Port 被佔用
- 記憶體不足

### Q2: 資料庫連線失敗

**檢查 PostgreSQL 容器:**

```bash
docker-compose ps postgres
docker-compose logs postgres
```

**進入容器檢查:**

```bash
docker exec -it asset-manager-postgres psql -U postgres -d asset_manager
```

### Q3: 前端無法連線到後端

**檢查環境變數:**

```bash
# 確認 NEXT_PUBLIC_API_URL 設定正確
cat .env.production | grep NEXT_PUBLIC_API_URL

# 確認 CORS 設定
cat .env.production | grep CORS_ALLOWED_ORIGINS
```

### Q4: GitHub Actions 部署失敗

**檢查 Secrets:**

- 確認所有必要的 Secrets 都已設定
- 確認 SSH Key 格式正確 (包含 BEGIN 和 END)

**查看 Actions 日誌:**

- 前往 GitHub Actions 查看詳細錯誤訊息

### Q5: 記憶體不足

**檢查記憶體使用:**

```bash
free -h
docker stats
```

**解決方法:**

- 升級 EC2 規格 (t3.small → t3.medium)
- 調整容器資源限制 (docker-compose.yml)
- 重啟容器釋放記憶體

---

## 維護指令

### 查看容器狀態

```bash
docker-compose ps
```

### 查看日誌

```bash
# 所有服務
docker-compose logs -f

# 特定服務
docker-compose logs -f backend
docker-compose logs -f frontend
```

### 重啟服務

```bash
# 重啟所有服務
docker-compose restart

# 重啟特定服務
docker-compose restart backend
```

### 更新程式碼

```bash
cd /home/ubuntu/asset-manager
bash scripts/deploy.sh
```

### 清理資源

```bash
# 清理未使用的映像檔
docker image prune -f

# 清理未使用的容器
docker container prune -f

# 清理未使用的 volumes (小心!)
docker volume prune -f
```

---

## 監控和日誌

### 系統資源監控

```bash
# 即時監控
htop

# Docker 資源使用
docker stats
```

### 應用程式日誌

```bash
# Backend 日誌
docker-compose logs -f backend

# Frontend 日誌
docker-compose logs -f frontend

# Nginx 日誌
docker-compose logs -f nginx
```

### 備份日誌

```bash
tail -f /home/ubuntu/backup.log
```

---

## 安全性建議

1. **定期更新系統套件**

   ```bash
   sudo apt update && sudo apt upgrade -y
   ```

2. **定期更新 Docker 映像檔**

   ```bash
   docker-compose pull
   docker-compose up -d
   ```

3. **定期檢查備份**

   ```bash
   ls -lh /home/ubuntu/backups/
   ```

4. **使用強密碼**

   - 資料庫密碼至少 16 個字元
   - 登入密碼至少 12 個字元
   - JWT Secret 使用 `openssl rand -base64 32` 產生

5. **限制 SSH 存取**
   - 只允許特定 IP 連線
   - 使用 SSH Key 而非密碼
   - 定期更換 SSH Key

---

## 聯絡資訊

如有問題,請聯絡:

- GitHub Issues: https://github.com/YOUR_USERNAME/asset-manager/issues
- Email: your-email@example.com
