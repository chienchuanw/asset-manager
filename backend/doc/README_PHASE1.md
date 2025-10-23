# Phase 1: 後端 Transactions API 實作指南

## 📋 已完成的檔案

### 1. 資料庫 Migration
- `migrations/000001_create_transactions_table.up.sql` - 建立 transactions 資料表
- `migrations/000001_create_transactions_table.down.sql` - 刪除 transactions 資料表

### 2. Models
- `internal/models/transaction.go` - Transaction 模型定義

### 3. Repository 層（資料存取層）
- `internal/repository/transaction_repository.go` - Repository 實作
- `internal/repository/transaction_repository_test.go` - Repository 測試
- `internal/repository/test_helper.go` - 測試輔助函式

### 4. Service 層（業務邏輯層）
- `internal/service/transaction_service.go` - Service 實作
- `internal/service/transaction_service_test.go` - Service 測試

### 5. API Handler 層
- `internal/api/transaction_handler.go` - API Handler 實作
- `internal/api/transaction_handler_test.go` - API Handler 測試

### 6. Main Application
- `cmd/api/main.go` - 主程式（已更新整合所有元件）

---

## 🚀 執行步驟

### Step 1: 安裝 Go 依賴套件

```bash
cd backend
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/suite
go get github.com/stretchr/testify/mock
go get github.com/google/uuid
go mod tidy
```

### Step 2: 設定資料庫

#### 2.1 建立開發資料庫

```bash
# 使用 psql 連接到 PostgreSQL
psql -U postgres

# 建立資料庫
CREATE DATABASE asset_manager;

# 建立測試資料庫
CREATE DATABASE asset_manager_test;

# 退出 psql
\q
```

#### 2.2 設定環境變數

複製 `.env.example` 並修改為你的資料庫設定：

```bash
cp .env.example .env
```

編輯 `.env`：
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=asset_manager

APP_PORT=8080
GIN_MODE=debug
```

### Step 3: 執行 Migration

```bash
# 安裝 migrate CLI（如果還沒安裝）
# macOS
brew install golang-migrate

# 或使用 go install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 執行 migration（開發資料庫）
migrate -path migrations -database "postgresql://postgres:your_password@localhost:5432/asset_manager?sslmode=disable" up

# 執行 migration（測試資料庫）
migrate -path migrations -database "postgresql://postgres:your_password@localhost:5432/asset_manager_test?sslmode=disable" up
```

### Step 4: 執行測試

#### 4.1 執行 Service 層測試（使用 Mock，不需要資料庫）

```bash
cd backend
go test ./internal/service/... -v
```

#### 4.2 執行 API Handler 測試（使用 Mock，不需要資料庫）

```bash
go test ./internal/api/... -v
```

#### 4.3 執行 Repository 測試（需要測試資料庫）

先設定測試資料庫環境變數：

```bash
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=your_password
export TEST_DB_NAME=asset_manager_test
```

然後執行測試：

```bash
go test ./internal/repository/... -v
```

#### 4.4 執行所有測試

```bash
# 設定環境變數後執行
go test ./... -v
```

### Step 5: 啟動 API 伺服器

```bash
# 確保已設定 .env 檔案
cd backend
go run cmd/api/main.go
```

伺服器會在 `http://localhost:8080` 啟動。

### Step 6: 測試 API

#### 6.1 Health Check

```bash
curl http://localhost:8080/health
```

#### 6.2 建立交易記錄

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

#### 6.3 取得所有交易記錄

```bash
curl http://localhost:8080/api/transactions
```

#### 6.4 取得單筆交易記錄

```bash
# 將 {id} 替換為實際的 UUID
curl http://localhost:8080/api/transactions/{id}
```

#### 6.5 更新交易記錄

```bash
# 將 {id} 替換為實際的 UUID
curl -X PUT http://localhost:8080/api/transactions/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 20,
    "price": 630,
    "amount": 12600
  }'
```

#### 6.6 刪除交易記錄

```bash
# 將 {id} 替換為實際的 UUID
curl -X DELETE http://localhost:8080/api/transactions/{id}
```

#### 6.7 使用篩選條件查詢

```bash
# 只查詢台股
curl "http://localhost:8080/api/transactions?asset_type=tw-stock"

# 查詢特定日期範圍
curl "http://localhost:8080/api/transactions?start_date=2025-10-01&end_date=2025-10-31"

# 分頁查詢
curl "http://localhost:8080/api/transactions?limit=10&offset=0"
```

---

## 🧪 TDD 開發流程說明

我們遵循了 TDD 的開發流程：

### 1. Repository 層
- ✅ 先寫測試 (`transaction_repository_test.go`)
- ✅ 再寫實作 (`transaction_repository.go`)
- ✅ 執行測試確認通過

### 2. Service 層
- ✅ 先寫測試 (`transaction_service_test.go`)，使用 Mock Repository
- ✅ 再寫實作 (`transaction_service.go`)
- ✅ 執行測試確認通過

### 3. API Handler 層
- ✅ 先寫測試 (`transaction_handler_test.go`)，使用 Mock Service
- ✅ 再寫實作 (`transaction_handler.go`)
- ✅ 執行測試確認通過

---

## 📊 API 端點總覽

| 方法 | 路徑 | 說明 |
|------|------|------|
| POST | `/api/transactions` | 建立交易記錄 |
| GET | `/api/transactions` | 取得交易記錄列表（支援篩選） |
| GET | `/api/transactions/:id` | 取得單筆交易記錄 |
| PUT | `/api/transactions/:id` | 更新交易記錄 |
| DELETE | `/api/transactions/:id` | 刪除交易記錄 |

---

## 🔍 常見問題

### Q1: Migration 執行失敗
**A:** 確認資料庫連線設定是否正確，以及資料庫是否已建立。

### Q2: Repository 測試失敗
**A:** 確認測試資料庫是否已建立，並且已執行 migration。

### Q3: 找不到 go 指令
**A:** 請先安裝 Go（建議版本 1.21 或以上）。

---

## ✅ 下一步

Phase 1 完成後，可以進行 Phase 2：前端整合

1. 安裝前端必要套件（React Query、react-hook-form、zod）
2. 建立 API Client
3. 實作交易列表顯示
4. 實作新增交易功能

