# 📜 Scripts 使用說明

這個目錄包含所有部署和維護相關的腳本。

## 📋 腳本列表

### 1. `setup-ec2.sh` - EC2 初始化腳本

**功能:**
- 安裝 Docker 和 Docker Compose
- 建立專案目錄
- 設定自動備份 Cron Job
- 設定 Git

**使用方式:**
```bash
# 在 EC2 上執行
bash scripts/setup-ec2.sh
```

**執行後需要:**
- 登出並重新登入 (讓 Docker 群組設定生效)
- Clone 專案程式碼
- 建立 .env.production 檔案

---

### 2. `backup-db.sh` - 資料庫備份腳本

**功能:**
- 備份 PostgreSQL 資料庫
- 壓縮備份檔案
- 刪除舊備份 (預設保留 7 天)
- (可選) 上傳到 S3

**使用方式:**
```bash
# 基本使用
bash scripts/backup-db.sh

# 自訂備份目錄
BACKUP_DIR=/path/to/backups bash scripts/backup-db.sh

# 自訂保留天數
RETENTION_DAYS=14 bash scripts/backup-db.sh

# 上傳到 S3
AWS_S3_BUCKET=my-bucket bash scripts/backup-db.sh
```

**環境變數:**
- `BACKUP_DIR` - 備份目錄 (預設: /home/ubuntu/backups)
- `CONTAINER_NAME` - PostgreSQL 容器名稱 (預設: asset-manager-postgres)
- `DB_NAME` - 資料庫名稱 (預設: asset_manager)
- `DB_USER` - 資料庫使用者 (預設: postgres)
- `RETENTION_DAYS` - 保留天數 (預設: 7)
- `AWS_S3_BUCKET` - S3 Bucket 名稱 (可選)
- `AWS_REGION` - AWS Region (預設: ap-northeast-1)

**輸出:**
- 備份檔案: `/home/ubuntu/backups/asset_manager_YYYYMMDD_HHMMSS.sql.gz`

---

### 3. `restore-db.sh` - 資料庫還原腳本

**功能:**
- 從備份檔案還原資料庫
- 支援壓縮和未壓縮的備份檔案
- 還原前自動建立安全備份

**使用方式:**
```bash
# 列出可用的備份檔案
ls -lh /home/ubuntu/backups/

# 還原指定的備份
bash scripts/restore-db.sh /home/ubuntu/backups/asset_manager_20241027_120000.sql.gz
```

**注意事項:**
- ⚠️ 此操作會覆蓋現有資料庫
- 還原前會建立安全備份
- 需要手動確認才會執行

---

### 4. `deploy.sh` - 部署腳本

**功能:**
- 部署前自動備份資料庫
- 拉取最新程式碼
- 重新建置並啟動容器
- 執行健康檢查

**使用方式:**
```bash
# 基本使用 (含備份)
bash scripts/deploy.sh

# 跳過備份
BACKUP_BEFORE_DEPLOY=false bash scripts/deploy.sh

# 自訂專案目錄
PROJECT_DIR=/path/to/project bash scripts/deploy.sh
```

**環境變數:**
- `PROJECT_DIR` - 專案目錄 (預設: /home/ubuntu/asset-manager)
- `BACKUP_BEFORE_DEPLOY` - 是否在部署前備份 (預設: true)

**執行流程:**
1. 備份資料庫 (如果啟用)
2. 拉取最新程式碼
3. 停止現有容器
4. 建置新映像檔
5. 啟動容器
6. 等待服務啟動
7. 健康檢查
8. 清理舊映像檔

---

### 5. `backup-to-s3.sh` - 上傳備份到 S3

**功能:**
- 備份資料庫
- 上傳到 S3
- 按年月組織備份檔案

**使用方式:**
```bash
# 設定環境變數
export AWS_S3_BUCKET=asset-manager-backups
export AWS_REGION=ap-northeast-1

# 執行備份並上傳
bash scripts/backup-to-s3.sh
```

**前置需求:**
- 安裝 AWS CLI: `sudo apt install -y awscli`
- 設定 AWS 憑證: `aws configure`
- 建立 S3 Bucket

**S3 路徑結構:**
```
s3://asset-manager-backups/
  └── backups/
      └── 2024/
          └── 10/
              ├── asset_manager_20241027_020000.sql.gz
              ├── asset_manager_20241028_020000.sql.gz
              └── ...
```

---

## 🔄 自動化備份設定

### Cron Job 設定

`setup-ec2.sh` 會自動設定每天凌晨 2 點執行備份:

```bash
# 查看 Cron Job
crontab -l

# 編輯 Cron Job
crontab -e
```

**預設設定:**
```cron
0 2 * * * /home/ubuntu/asset-manager/scripts/backup-db.sh >> /home/ubuntu/backup.log 2>&1
```

### 查看備份日誌

```bash
tail -f /home/ubuntu/backup.log
```

---

## 🚀 常見使用場景

### 場景 1: 首次部署

```bash
# 1. SSH 連線到 EC2
ssh -i ~/.ssh/asset-manager-key.pem ubuntu@43.213.77.244

# 2. 執行初始化腳本
bash scripts/setup-ec2.sh

# 3. 登出並重新登入
exit
ssh -i ~/.ssh/asset-manager-key.pem ubuntu@43.213.77.244

# 4. Clone 專案
cd /home/ubuntu/asset-manager
git clone https://github.com/YOUR_USERNAME/asset-manager.git .

# 5. 建立環境變數
cp .env.production.example .env.production
vim .env.production

# 6. 啟動服務
docker-compose --env-file .env.production up -d
```

### 場景 2: 更新程式碼

```bash
# 使用部署腳本 (推薦)
bash scripts/deploy.sh

# 或手動執行
git pull origin main
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### 場景 3: 定期備份

```bash
# 手動備份
bash scripts/backup-db.sh

# 查看備份檔案
ls -lh /home/ubuntu/backups/

# 上傳到 S3
export AWS_S3_BUCKET=asset-manager-backups
bash scripts/backup-to-s3.sh
```

### 場景 4: 災難復原

```bash
# 1. 列出可用的備份
ls -lh /home/ubuntu/backups/

# 2. 還原最新的備份
bash scripts/restore-db.sh /home/ubuntu/backups/asset_manager_LATEST.sql.gz

# 3. 重啟服務
docker-compose restart backend
```

---

## ⚠️ 注意事項

### 權限設定

所有腳本需要執行權限:

```bash
chmod +x scripts/*.sh
```

### 環境變數

確保 `.env.production` 檔案存在且包含所有必要的變數。

### 備份空間

定期檢查備份目錄的空間使用:

```bash
du -sh /home/ubuntu/backups/
df -h
```

### S3 費用

使用 S3 備份會產生費用:
- 儲存費用: ~$0.023/GB/月
- 上傳流量: 免費
- 下載流量: ~$0.09/GB

---

## 🔧 故障排除

### 問題 1: 備份腳本執行失敗

**檢查容器狀態:**
```bash
docker ps | grep postgres
```

**檢查日誌:**
```bash
docker-compose logs postgres
```

### 問題 2: 還原失敗

**檢查備份檔案:**
```bash
# 檢查檔案是否存在
ls -lh /home/ubuntu/backups/asset_manager_*.sql.gz

# 檢查檔案內容
gunzip -c backup.sql.gz | head -20
```

### 問題 3: S3 上傳失敗

**檢查 AWS CLI:**
```bash
aws --version
aws s3 ls
```

**檢查權限:**
```bash
aws iam get-user
```

---

## 📞 需要幫助?

如有問題,請查看:
- [DEPLOYMENT.md](../DEPLOYMENT.md) - 完整部署指南
- [GitHub Issues](https://github.com/YOUR_USERNAME/asset-manager/issues)

