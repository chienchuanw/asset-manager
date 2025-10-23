# 🚀 Quick Start Guide

## 前置需求

- Go 1.21 或以上
- PostgreSQL 12 或以上
- golang-migrate CLI

---

## 快速開始（5 分鐘）

### 1. 自動化設定（推薦）

```bash
cd backend
chmod +x scripts/setup.sh
./scripts/setup.sh
```

這個腳本會自動：
- ✅ 檢查 Go 和 PostgreSQL 安裝
- ✅ 安裝 golang-migrate（如果需要）
- ✅ 安裝所有 Go 依賴套件
- ✅ 建立 .env 檔案
- ✅ 建立資料庫（可選）
- ✅ 執行 migration（可選）

### 2. 手動設定

如果你想手動設定，請按照以下步驟：

#### Step 1: 安裝依賴
```bash
cd backend
make install
```

#### Step 2: 設定環境變數
```bash
cp .env.example .env
# 編輯 .env 檔案，設定你的資料庫連線資訊
```

#### Step 3: 建立資料庫
```bash
# 使用 psql 連接到 PostgreSQL
psql -U postgres

# 建立資料庫
CREATE DATABASE asset_manager;
CREATE DATABASE asset_manager_test;

# 退出
\q
```

#### Step 4: 執行 Migration
```bash
# 方法 1: 使用 Makefile（需要先設定環境變數）
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=asset_manager

make migrate-up

# 方法 2: 直接使用 migrate CLI
migrate -path migrations \
  -database "postgresql://postgres:your_password@localhost:5432/asset_manager?sslmode=disable" \
  up
```

---

## 執行測試

### 執行所有測試
```bash
make test
```

### 只執行單元測試（不需要資料庫）
```bash
make test-unit
```

### 只執行整合測試（需要資料庫）
```bash
# 先設定測試資料庫環境變數
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=your_password
export TEST_DB_NAME=asset_manager_test

# 執行測試
make test-integration
```

---

## 啟動 API 伺服器

```bash
make run
```

伺服器會在 `http://localhost:8080` 啟動。

---

## 測試 API

### 方法 1: 使用測試腳本（推薦）

```bash
chmod +x scripts/test-api.sh
./scripts/test-api.sh
```

### 方法 2: 手動測試

#### Health Check
```bash
curl http://localhost:8080/health
```

#### 建立交易記錄
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2025-10-22T00:00:00Z",
    "asset_type": "tw-stock",
    "symbol": "2330",
    "name": "台積電",
    "type": "buy",
    "quantity": 10,
    "price": 620,
    "amount": 6200,
    "fee": 28,
    "note": "定期定額買入"
  }'
```

#### 取得所有交易記錄
```bash
curl http://localhost:8080/api/transactions
```

#### 取得單筆交易記錄
```bash
# 將 {id} 替換為實際的 UUID
curl http://localhost:8080/api/transactions/{id}
```

#### 更新交易記錄
```bash
curl -X PUT http://localhost:8080/api/transactions/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 20,
    "price": 630,
    "amount": 12600
  }'
```

#### 刪除交易記錄
```bash
curl -X DELETE http://localhost:8080/api/transactions/{id}
```

#### 使用篩選條件
```bash
# 只查詢台股
curl "http://localhost:8080/api/transactions?asset_type=tw-stock"

# 查詢日期範圍
curl "http://localhost:8080/api/transactions?start_date=2025-10-01&end_date=2025-10-31"

# 分頁查詢
curl "http://localhost:8080/api/transactions?limit=10&offset=0"
```

---

## 常用指令

```bash
# 查看所有可用指令
make help

# 安裝依賴
make install

# 執行測試
make test

# 執行單元測試
make test-unit

# 執行整合測試
make test-integration

# 執行 migration
make migrate-up

# 回滾 migration
make migrate-down

# 啟動伺服器
make run

# 編譯應用程式
make build

# 清理編譯產物
make clean
```

---

## 故障排除

### 問題 1: 找不到 go 指令
**解決方法：** 安裝 Go 1.21 或以上版本
- macOS: `brew install go`
- Ubuntu: `sudo apt-get install golang-go`
- 或從官網下載：https://golang.org/dl/

### 問題 2: 找不到 migrate 指令
**解決方法：** 安裝 golang-migrate
- macOS: `brew install golang-migrate`
- Linux: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

### 問題 3: 資料庫連線失敗
**解決方法：**
1. 確認 PostgreSQL 是否正在執行
2. 檢查 .env 檔案中的資料庫設定是否正確
3. 確認資料庫是否已建立

### 問題 4: Migration 執行失敗
**解決方法：**
1. 確認資料庫連線設定正確
2. 確認資料庫已建立
3. 檢查 migration 檔案是否存在於 `migrations/` 目錄

### 問題 5: 測試失敗
**解決方法：**
1. 單元測試失敗：檢查程式碼邏輯
2. 整合測試失敗：
   - 確認測試資料庫已建立
   - 確認測試資料庫已執行 migration
   - 確認測試環境變數已設定

---

## 下一步

✅ Phase 1 完成後，可以進行 Phase 2：前端整合

詳細資訊請參考：
- `README_PHASE1.md` - 完整的實作指南
- `PHASE1_SUMMARY.md` - Phase 1 完成總結

