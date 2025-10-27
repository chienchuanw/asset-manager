# 🚀 快速開始指南

這份指南會帶你在 10 分鐘內完成 Asset Manager 的部署。

---

## ✅ 前置檢查清單

在開始之前,確認你已經完成:

- [ ] AWS EC2 Instance 已啟動 (IP: 43.213.77.244)
- [ ] SSH Key 已下載並設定權限
- [ ] 可以 SSH 連線到 EC2
- [ ] Docker 和 Docker Compose 已安裝

如果還沒完成,請先參考 [DEPLOYMENT.md](DEPLOYMENT.md)。

---

## 📝 Step 1: 準備環境變數

### 1.1 產生 JWT Secret

在你的**本機**執行:

```bash
openssl rand -base64 32
```

複製產生的字串,例如: `XkVdiQpHuvmD8EL/b7izSs/ZD9AadgGEVvi95jsL6ko=`

### 1.2 準備 API Keys

確認你已經申請以下 API Keys:

- [ ] FinMind API Key - https://finmind.github.io/
- [ ] CoinGecko API Key - https://www.coingecko.com/en/api
- [ ] Alpha Vantage API Key - https://www.alphavantage.co/support/#api-key

---

## 🖥️ Step 2: 在 EC2 上部署

### 2.1 SSH 連線到 EC2

```bash
ssh -i ~/.ssh/asset-manager-key.pem ubuntu@43.213.77.244
```

### 2.2 Clone 專案程式碼

```bash
cd /home/ubuntu
git clone https://github.com/chienchuanw/asset-manager.git
cd asset-manager
```

### 2.3 建立環境變數檔案

```bash
cp .env.production.example .env.production
vim .env.production
```

**填入以下必要的值:**

```bash
# 資料庫密碼 (自己設定一個強密碼)
DB_PASSWORD=YOUR_STRONG_PASSWORD_HERE

# 身份驗證
AUTH_USERNAME=admin
AUTH_PASSWORD=YOUR_STRONG_PASSWORD_HERE
JWT_SECRET=YOUR_JWT_SECRET_FROM_STEP_1

# API Keys
FINMIND_API_KEY=YOUR_FINMIND_KEY
COINGECKO_API_KEY=YOUR_COINGECKO_KEY
ALPHA_VANTAGE_API_KEY=YOUR_ALPHA_VANTAGE_KEY

# CORS 和 API URL (使用你的 EC2 IP)
CORS_ALLOWED_ORIGINS=http://43.213.77.244
NEXT_PUBLIC_API_URL=http://43.213.77.244/api
```

**儲存並離開:** 按 `Esc`,輸入 `:wq`,按 `Enter`

### 2.4 啟動服務

```bash
# 建置並啟動所有容器
docker-compose --env-file .env.production up -d

# 查看容器狀態
docker-compose ps

# 查看日誌 (確認沒有錯誤)
docker-compose logs -f
```

**按 `Ctrl+C` 停止查看日誌**

### 2.5 驗證部署

```bash
# 檢查 Backend API
curl http://localhost:8080/health

# 應該看到: {"status":"OK","message":"Asset Manager API Server is running."}

# 檢查 Frontend
curl http://localhost:3000

# 應該看到 HTML 內容

# 檢查 Nginx
curl http://localhost/health

# 應該看到: healthy
```

---

## 🌐 Step 3: 從瀏覽器訪問

### 3.1 開啟瀏覽器

前往: `http://43.213.77.244`

### 3.2 登入

- 帳號: `admin` (或你在 .env.production 設定的)
- 密碼: 你在 .env.production 設定的密碼

### 3.3 測試功能

- 查看 Dashboard
- 新增一筆交易記錄
- 查看持倉資訊

---

## 🔄 Step 4: 設定自動備份

### 4.1 設定 Cron Job

```bash
# 編輯 crontab
crontab -e

# 如果是第一次使用,選擇編輯器 (建議選 vim)

# 加入以下這行 (每天凌晨 2 點自動備份)
0 2 * * * /home/ubuntu/asset-manager/scripts/backup-db.sh >> /home/ubuntu/backup.log 2>&1

# 儲存並離開
```

### 4.2 測試備份

```bash
# 手動執行備份
bash scripts/backup-db.sh

# 查看備份檔案
ls -lh /home/ubuntu/backups/
```

---

## 🤖 Step 5: 設定 GitHub Actions 自動部署

### 5.1 建立 Discord Webhook (可選)

1. 開啟 Discord,選擇一個頻道
2. 頻道設定 → 整合 → Webhooks → 新增 Webhook
3. 複製 Webhook URL

### 5.2 設定 GitHub Secrets

前往: `https://github.com/chienchuanw/asset-manager/settings/secrets/actions`

點擊 **"New repository secret"**,新增以下 Secrets:

**EC2 連線:**

```
Name: EC2_HOST
Value: 43.213.77.244

Name: EC2_USERNAME
Value: ubuntu

Name: EC2_SSH_KEY
Value: (貼上 ~/.ssh/asset-manager-key.pem 的完整內容)
```

**資料庫:**

```
Name: PROD_DB_USER
Value: postgres

Name: PROD_DB_PASSWORD
Value: (你在 .env.production 設定的密碼)

Name: PROD_DB_NAME
Value: asset_manager
```

**Redis:**

```
Name: PROD_REDIS_PASSWORD
Value: (留空或設定密碼)
```

**身份驗證:**

```
Name: PROD_AUTH_USERNAME
Value: admin

Name: PROD_AUTH_PASSWORD
Value: (你在 .env.production 設定的密碼)

Name: PROD_JWT_SECRET
Value: (你在 Step 1 產生的 JWT Secret)
```

**API Keys:**

```
Name: PROD_FINMIND_API_KEY
Value: (你的 FinMind API Key)

Name: PROD_COINGECKO_API_KEY
Value: (你的 CoinGecko API Key)

Name: PROD_ALPHA_VANTAGE_API_KEY
Value: (你的 Alpha Vantage API Key)
```

**Discord (可選):**

```
Name: DISCORD_WEBHOOK_URL
Value: (你的 Discord Webhook URL)
```

### 5.3 測試自動部署

```bash
# 在本機修改程式碼
git add .
git commit -m "test: 測試自動部署"
git push origin main

# 前往 GitHub Actions 查看部署進度
# https://github.com/chienchuanw/asset-manager/actions
```

---

## ✅ 完成!

恭喜!你已經成功部署 Asset Manager 了! 🎉

### 接下來可以做什麼?

1. **設定 HTTPS**

   - 註冊網域名稱
   - 使用 Let's Encrypt 取得免費 SSL 憑證
   - 更新 nginx.conf

2. **設定 S3 備份**

   - 建立 S3 Bucket
   - 設定 AWS CLI
   - 使用 `backup-to-s3.sh` 腳本

3. **監控和優化**
   - 設定 CloudWatch 監控
   - 調整容器資源限制
   - 優化資料庫查詢

---

## 📚 更多資源

- [完整部署指南](DEPLOYMENT.md)
- [腳本使用說明](scripts/README.md)
- [GitHub Issues](https://github.com/chienchuanw/asset-manager/issues)

---

## ⚠️ 常見問題

### Q: 容器無法啟動

**A:** 檢查日誌:

```bash
docker-compose logs backend
docker-compose logs frontend
```

### Q: 無法從瀏覽器訪問

**A:** 檢查 Security Group 是否開放 Port 80:

- 前往 AWS Console → EC2 → Security Groups
- 確認有 Port 80 的 Inbound Rule

### Q: 登入失敗

**A:** 檢查環境變數:

```bash
cat .env.production | grep AUTH_
```

---

## 🆘 需要幫助?

如果遇到問題:

1. 查看日誌: `docker-compose logs -f`
2. 檢查容器狀態: `docker-compose ps`
3. 查看 [DEPLOYMENT.md](DEPLOYMENT.md) 的故障排除章節
4. 在 GitHub 開 Issue

祝你使用愉快! 🚀
