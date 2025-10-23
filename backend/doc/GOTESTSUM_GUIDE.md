# Gotestsum 使用指南

## 📖 什麼是 Gotestsum？

`gotestsum` 是一個增強版的 Go 測試執行器，提供：
- ✅ 彩色輸出
- ✅ 更清晰的測試結果顯示
- ✅ 多種輸出格式
- ✅ Watch 模式（檔案變更時自動重新執行）
- ✅ 失敗測試自動重試
- ✅ 測試覆蓋率報告

---

## 🚀 快速開始

### 基本測試指令

```bash
# 執行所有測試（預設格式：testname）
make test

# 執行所有測試（詳細模式）
make test-verbose

# 只執行單元測試
make test-unit

# 只執行整合測試
make test-integration

# 執行測試並產生覆蓋率報告
make test-coverage

# Watch 模式（檔案變更時自動重新執行）
make test-watch
```

---

## 🎨 輸出格式說明

### 1. `testname` 格式（預設）

簡潔的測試名稱列表，適合快速查看測試結果。

```
✓ TestCreateTransaction_Success
✓ TestCreateTransaction_InvalidInput
✓ TestGetTransaction_Success
✗ TestGetTransaction_InvalidID
```

### 2. `standard-verbose` 格式

顯示完整的測試輸出，包括所有 log 訊息。

```
=== RUN   TestCreateTransaction_Success
--- PASS: TestCreateTransaction_Success (0.00s)
=== RUN   TestCreateTransaction_InvalidInput
--- PASS: TestCreateTransaction_InvalidInput (0.00s)
```

### 3. `dots` 格式

每個測試用一個點表示，非常簡潔。

```
..........
```

### 4. `pkgname` 格式

按套件分組顯示測試結果。

```
✓ github.com/chienchuanw/asset-manager/internal/api
✓ github.com/chienchuanw/asset-manager/internal/service
✗ github.com/chienchuanw/asset-manager/internal/repository
```

---

## 🔧 進階用法

### 直接使用 gotestsum

```bash
# 基本用法
gotestsum --format testname -- -cover ./...

# 指定輸出格式
gotestsum --format standard-verbose -- -cover ./...

# 只執行特定套件
gotestsum --format testname -- -cover ./internal/service/...

# 執行特定測試
gotestsum --format testname -- -run TestCreateTransaction ./...

# Watch 模式
gotestsum --watch --format testname -- -cover ./...

# 產生 JSON 輸出
gotestsum --jsonfile test-output.json --format testname -- -cover ./...

# 失敗時重新執行
gotestsum --rerun-fails --format testname -- -cover ./...
```

---

## 📊 測試覆蓋率

### 產生覆蓋率報告

```bash
# 使用 Makefile（推薦）
make test-coverage

# 直接使用 gotestsum
gotestsum --format testname -- -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

這會產生兩個檔案：
- `coverage.out` - 覆蓋率資料
- `coverage.html` - HTML 格式的覆蓋率報告

### 查看覆蓋率報告

```bash
# 在瀏覽器中開啟
open coverage.html

# 或在終端中查看
go tool cover -func=coverage.out
```

---

## 👀 Watch 模式

Watch 模式會監控檔案變更，自動重新執行測試。

```bash
# 使用 Makefile
make test-watch

# 直接使用 gotestsum
gotestsum --watch --format testname -- -cover ./...
```

**使用技巧**：
- 修改程式碼後儲存，測試會自動執行
- 按 `Ctrl+C` 停止 watch 模式
- 適合 TDD 開發流程

---

## 🎯 測試篩選

### 執行特定測試

```bash
# 執行名稱包含 "Create" 的測試
gotestsum --format testname -- -run Create ./...

# 執行特定套件的測試
gotestsum --format testname -- -cover ./internal/service/...

# 執行多個套件的測試
gotestsum --format testname -- -cover ./internal/service/... ./internal/api/...
```

### 排除特定測試

```bash
# 排除名稱包含 "Integration" 的測試
gotestsum --format testname -- -skip Integration ./...
```

---

## 🔄 失敗測試重試

gotestsum 可以自動重新執行失敗的測試。

```bash
# 失敗時重新執行（最多 2 次）
gotestsum --rerun-fails --rerun-fails-max-attempts=2 --format testname -- -cover ./...
```

這在以下情況很有用：
- 測試有時會因為時序問題而失敗
- 網路相關的測試
- 資料庫連線測試

---

## 📝 設定檔

gotestsum 可以使用設定檔 `.gotestsum.yml`：

```yaml
# .gotestsum.yml
format: testname
show-elapsed: true
hide-summary: false
rerun-fails: failed
rerun-fails-max-attempts: 2
timeout: 0
```

有了設定檔，只需要執行：

```bash
gotestsum
```

---

## 🎨 自訂輸出

### 產生 JUnit XML 報告（用於 CI/CD）

```bash
gotestsum --junitfile junit.xml --format testname -- -cover ./...
```

### 產生 JSON 輸出

```bash
gotestsum --jsonfile test-output.json --format testname -- -cover ./...
```

### 同時產生多種報告

```bash
gotestsum \
  --format testname \
  --jsonfile test-output.json \
  --junitfile junit.xml \
  -- -coverprofile=coverage.out ./...
```

---

## 🐛 除錯技巧

### 顯示詳細的測試輸出

```bash
# 使用 standard-verbose 格式
gotestsum --format standard-verbose -- -v ./...

# 顯示所有 log 訊息
gotestsum --format testname -- -v ./...
```

### 只執行失敗的測試

```bash
# 第一次執行，記錄失敗的測試
gotestsum --format testname -- -cover ./...

# 只重新執行失敗的測試
gotestsum --rerun-fails-only --format testname -- -cover ./...
```

---

## 📋 常用指令速查表

| 指令 | 說明 |
|------|------|
| `make test` | 執行所有測試 |
| `make test-verbose` | 執行所有測試（詳細模式） |
| `make test-unit` | 執行單元測試 |
| `make test-integration` | 執行整合測試 |
| `make test-coverage` | 產生覆蓋率報告 |
| `make test-watch` | Watch 模式 |
| `gotestsum --format testname` | 簡潔格式 |
| `gotestsum --format standard-verbose` | 詳細格式 |
| `gotestsum --watch` | Watch 模式 |
| `gotestsum --rerun-fails` | 重新執行失敗的測試 |

---

## 🎯 最佳實踐

### 1. 開發時使用 Watch 模式

```bash
make test-watch
```

這樣可以即時看到程式碼變更的影響。

### 2. CI/CD 使用詳細模式

```bash
gotestsum --format standard-verbose --junitfile junit.xml -- -cover ./...
```

這樣可以在 CI/CD 系統中看到完整的測試輸出。

### 3. 本地開發使用簡潔格式

```bash
make test
```

快速查看測試結果，不需要太多細節。

### 4. 定期檢查覆蓋率

```bash
make test-coverage
open coverage.html
```

確保測試覆蓋率保持在合理水平。

---

## 🔗 相關資源

- [gotestsum GitHub](https://github.com/gotestyourself/gotestsum)
- [Go Testing 官方文件](https://golang.org/pkg/testing/)
- [測試最佳實踐](https://go.dev/doc/tutorial/add-a-test)

---

## 💡 提示

1. **彩色輸出**：gotestsum 會自動偵測終端是否支援彩色輸出
2. **效能**：gotestsum 不會影響測試執行速度
3. **相容性**：gotestsum 完全相容 `go test` 的所有參數
4. **CI/CD**：可以在 CI/CD 環境中使用 gotestsum 產生報告

---

**祝測試順利！** 🎉

