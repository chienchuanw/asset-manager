# Testing Guide

## 🚀 Quick Start

### 一鍵設定（推薦）

```bash
cd backend
chmod +x scripts/quick-setup.sh
./scripts/quick-setup.sh
```

這個腳本會自動：
1. ✅ 載入環境變數
2. ✅ 建立開發和測試資料庫
3. ✅ 執行所有 migrations
4. ✅ 執行所有測試

---

## 📋 手動設定步驟

### Step 1: 建立資料庫

```bash
# 使用 Makefile（推薦）
make db-create

# 或手動建立
psql -U postgres -c "CREATE DATABASE asset_manager;"
psql -U postgres -c "CREATE DATABASE asset_manager_test;"
```

### Step 2: 執行 Migrations

```bash
# 開發資料庫
source .env.local
make migrate-up

# 測試資料庫
source .env.test
make migrate-test-up

# 或使用自動載入環境變數的版本
make migrate-up-env
make migrate-test-up-env
```

### Step 3: 執行測試

```bash
# 執行所有測試（彩色輸出）
make test

# 只執行單元測試（不需要資料庫）
make test-unit

# 只執行整合測試（需要測試資料庫）
source .env.test
make test-integration
```

---

## 🎨 測試輸出說明

測試結果會以彩色顯示：

- **綠色 (PASS)**: 測試通過 ✅
- **紅色 (FAIL)**: 測試失敗 ❌
- **藍色 (coverage)**: 測試覆蓋率 📊
- **黃色 (warning)**: 警告訊息 ⚠️

---

## 🧪 測試類型

### 1. 單元測試（Unit Tests）

**位置**: `internal/service/`, `internal/api/`

**特點**:
- 使用 Mock 模擬依賴
- 不需要資料庫
- 執行速度快

**執行**:
```bash
make test-unit
```

**測試內容**:
- Service 層業務邏輯驗證
- API Handler 層 HTTP 請求處理
- 輸入驗證和錯誤處理

### 2. 整合測試（Integration Tests）

**位置**: `internal/repository/`

**特點**:
- 需要真實的測試資料庫
- 測試資料庫互動
- 執行速度較慢

**執行**:
```bash
# 確保測試資料庫已建立並執行 migration
source .env.test
make test-integration
```

**測試內容**:
- Repository 層 CRUD 操作
- SQL 查詢正確性
- 資料庫約束驗證

---

## 📊 測試覆蓋率

查看測試覆蓋率：

```bash
# 所有測試的覆蓋率
go test ./... -cover

# 產生詳細的覆蓋率報告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 🔍 常見問題排除

### 問題 1: 測試失敗 - "database does not exist"

**原因**: 測試資料庫不存在

**解決方法**:
```bash
make db-create
make migrate-test-up-env
```

### 問題 2: 測試失敗 - "connection refused"

**原因**: PostgreSQL 沒有執行

**解決方法**:
```bash
# macOS
brew services start postgresql

# 檢查狀態
brew services list | grep postgresql
```

### 問題 3: 環境變數沒有載入

**原因**: 沒有 source .env.local 或 .env.test

**解決方法**:
```bash
# 載入開發環境變數
source .env.local

# 載入測試環境變數
source .env.test

# 或使用 Makefile 的 *-env 版本
make migrate-up-env
make migrate-test-up-env
```

### 問題 4: 找不到 migrate 指令

**原因**: golang-migrate 沒有安裝

**解決方法**:
```bash
# macOS
brew install golang-migrate

# 驗證安裝
migrate -version
```

### 問題 5: 語法錯誤 - "missing import path"

**原因**: Go 檔案有語法錯誤

**解決方法**:
```bash
# 檢查語法
go vet ./...

# 格式化程式碼
go fmt ./...
```

---

## 🎯 測試最佳實踐

### 1. 執行測試前

```bash
# 確保程式碼格式正確
go fmt ./...

# 檢查語法錯誤
go vet ./...

# 確保依賴是最新的
go mod tidy
```

### 2. 測試隔離

- 每個測試應該獨立運行
- 使用 `SetupTest()` 和 `TearDownTest()` 清理資料
- 不要依賴測試執行順序

### 3. 測試命名

- 測試函式名稱應該清楚描述測試內容
- 格式: `Test<FunctionName>_<Scenario>`
- 例如: `TestCreateTransaction_Success`, `TestCreateTransaction_InvalidInput`

### 4. 使用 Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid input", "valid", false},
        {"invalid input", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

---

## 📝 Makefile 指令總覽

| 指令 | 說明 | 需要資料庫 |
|------|------|-----------|
| `make test` | 執行所有測試（彩色輸出） | ✅ |
| `make test-unit` | 執行單元測試 | ❌ |
| `make test-integration` | 執行整合測試 | ✅ |
| `make db-create` | 建立開發和測試資料庫 | - |
| `make db-drop` | 刪除開發和測試資料庫 | - |
| `make migrate-up` | 執行開發資料庫 migration | - |
| `make migrate-up-env` | 載入 .env.local 並執行 migration | - |
| `make migrate-test-up` | 執行測試資料庫 migration | - |
| `make migrate-test-up-env` | 載入 .env.test 並執行 migration | - |

---

## 🚀 完整測試流程

```bash
# 1. 進入 backend 目錄
cd backend

# 2. 建立資料庫
make db-create

# 3. 執行 migrations
make migrate-up-env
make migrate-test-up-env

# 4. 執行所有測試
make test

# 5. 如果測試通過，啟動伺服器
make run
```

---

## 📞 需要幫助？

如果遇到問題：

1. 檢查 `.env.local` 和 `.env.test` 是否正確設定
2. 確認 PostgreSQL 正在執行
3. 確認 golang-migrate 已安裝
4. 查看錯誤訊息並參考「常見問題排除」章節
5. 執行 `make help` 查看所有可用指令

---

**祝測試順利！** 🎉

