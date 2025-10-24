# Mock Data Seeding 實作總結

## 📋 概述

實作了一個完整的 Mock 資料匯入系統，可以從 CSV 檔案讀取真實的資產配置資料並匯入資料庫。

## ✅ 完成項目

### 1. Seed 程式 (`cmd/seed/main.go`)

**功能：**

- 讀取 CSV 檔案
- 解析資產配置資料
- 轉換為交易記錄
- 匯入資料庫
- 支援清空資料庫選項

**特色：**

- ✅ 自動載入環境變數（`.env.local` 或 `.env`）
- ✅ 支援自訂 CSV 檔案路徑
- ✅ 完整的錯誤處理和日誌輸出
- ✅ 保留數量的完整精度（`DECIMAL(18, 8)`）
- ✅ 支援台股、美股、加密貨幣三種資產類型

### 2. Makefile 指令

新增兩個 Make 指令：

```makefile
# 匯入 Mock 資料（保留現有資料）
make seed

# 清空資料庫並匯入 Mock 資料
make seed-clean
```

### 3. 文檔

- `backend/mock/README.md` - 使用說明
- `backend/doc/SEED_DATA_IMPLEMENTATION.md` - 實作總結（本文件）

## 📊 資料統計

### CSV 資料

- **總計**：28 筆資產
- **台股**：16 筆
- **美股**：8 筆
- **加密貨幣**：4 筆

### 匯入結果

```text
✅ Import completed: 28/28 transactions imported successfully
```

**各資產類型統計：**

| 資產類型          | 數量 | 總金額 (原幣別)  |
| ----------------- | ---- | ---------------- |
| 台股 (tw-stock)   | 16   | 2,496,154.51 TWD |
| 美股 (us-stock)   | 8    | 4,182.26 USD     |
| 加密貨幣 (crypto) | 4    | 793.75 USD       |

## 🔧 技術實作細節

### 資料轉換規則

#### 1. 資產類型對應

```go
"Taiwan Stocks" → models.AssetTypeTWStock  // "tw-stock"
"US Stocks"     → models.AssetTypeUSStock  // "us-stock"
"Crypto"        → models.AssetTypeCrypto   // "crypto"
```

#### 2. 幣別對應

```go
"TWD" → models.CurrencyTWD  // "TWD"
"USD" → models.CurrencyUSD  // "USD"
```

#### 3. 交易記錄生成

```go
CreateTransactionInput{
    Date:            2025-10-20,           // CSV 中的 Update Date
    AssetType:       convertAssetType(),   // 轉換後的資產類型
    Symbol:          record.Ticker,        // 股票代碼
    Name:            record.Name,          // 資產名稱
    TransactionType: TransactionTypeBuy,   // 固定為「買入」
    Quantity:        parseFloat(Units),    // 持有數量
    Price:           parseFloat(Price),    // 價格
    Amount:          Quantity × Price,     // 自動計算
    Fee:             0.0,                  // 手續費為 0
    Currency:        convertCurrency(),    // 轉換後的幣別
    Note:            nil,                  // 無備註
}
```

### 數量精度處理

**資料庫 Schema：**

```sql
quantity DECIMAL(18, 8)  -- 保留 8 位小數
```

**實際儲存範例：**

- BTC: `0.00215432`（完整精度）
- ETH: `0.05863319`（完整精度）
- AAPL: `1.69698000`（完整精度）

**Log 輸出格式：**

```go
log.Printf("✓ Imported %s - %.2f %s @ %.2f\n", ...)
// BTC 顯示為 0.00（僅用於顯示）
```

### CSV 解析

**處理千分位逗號：**

```go
func parseFloat(s string) (float64, error) {
    s = strings.ReplaceAll(s, ",", "")  // 移除逗號
    return strconv.ParseFloat(s, 64)
}
```

**範例：**

- `"885,806.45"` → `885806.45`
- `"14207"` → `14207.0`

### 資料庫清空

**執行順序（避免外鍵約束錯誤）：**

```go
1. DELETE FROM realized_profits  // 先刪除子表
2. DELETE FROM transactions       // 再刪除父表
```

## 📝 使用範例

### 基本使用

```bash
# 1. 清空資料庫並匯入 Mock 資料
make seed-clean

# 2. 啟動後端伺服器
make run

# 3. 啟動前端（另一個終端）
cd ../frontend
pnpm dev

# 4. 瀏覽器訪問 http://localhost:3000
```

### 進階使用

```bash
# 使用自訂 CSV 檔案
go run cmd/seed/main.go --csv=path/to/custom.csv

# 只匯入不清空
go run cmd/seed/main.go

# 清空並匯入
go run cmd/seed/main.go --clean
```

### 驗證資料

```bash
# 查看所有交易記錄
psql -U chienchuanw -d db_asset_manager -c "SELECT * FROM transactions;"

# 查看各資產類型統計
psql -U chienchuanw -d db_asset_manager -c "
  SELECT asset_type, COUNT(*), SUM(amount)
  FROM transactions
  GROUP BY asset_type;
"

# 查看加密貨幣的完整精度
psql -U chienchuanw -d db_asset_manager -c "
  SELECT symbol, quantity, price
  FROM transactions
  WHERE asset_type = 'crypto';
"
```

## 🎯 設計決策

### 1. 為什麼選擇 Go 程式而非 SQL 腳本？

**優點：**

- ✅ 可以重用現有的 Repository 和 Service
- ✅ 可以進行資料驗證和錯誤處理
- ✅ 可以讀取 CSV 並自動轉換
- ✅ 可以輕鬆擴展功能（例如：支援多種格式）

**缺點：**

- ❌ 需要編譯（但 `go run` 可以直接執行）

### 2. 為什麼保留 DECIMAL(18, 8) 而非改為 DECIMAL(18, 2)？

**原因：**

- 加密貨幣需要高精度（例如：BTC 0.00215432）
- 資料庫保留完整精度，前端顯示時再決定精度
- 避免精度損失（0.00215432 → 0.00 會變成 0）

**解決方案：**

- 資料庫：`DECIMAL(18, 8)`（保留完整精度）
- 前端顯示：根據資產類型決定顯示精度
  - 台股、美股：2 位小數
  - 加密貨幣：8 位小數

### 3. 為什麼所有交易都是「買入」？

**原因：**

- CSV 只包含目前持倉狀態，沒有歷史交易記錄
- 使用者要求「選項 A：只建立買入交易」
- 簡化實作，專注於資料匯入功能

**未來擴展：**

- 可以支援匯入完整的交易歷史（包含買入和賣出）
- 可以從券商匯出的交易明細 CSV 匯入

## 🚀 後續建議

### 1. 前端顯示優化

建議在前端根據資產類型顯示不同精度：

```typescript
function formatQuantity(quantity: number, assetType: string): string {
  if (assetType === "crypto") {
    return quantity.toFixed(8); // 加密貨幣顯示 8 位
  }
  return quantity.toFixed(2); // 股票顯示 2 位
}
```

### 2. 支援更多資料來源

可以擴展 Seed 程式支援：

- 券商交易明細 CSV
- Excel 檔案（`.xlsx`）
- JSON 格式
- API 匯入（例如：從券商 API 自動同步）

### 3. 資料驗證

可以加入更多驗證規則：

- 數量必須 > 0
- 價格必須 > 0
- 日期不能是未來
- 股票代碼格式驗證

### 4. 批次處理

對於大量資料，可以使用批次插入提升效能：

```go
// 使用 transaction 批次插入
tx, _ := db.Begin()
for _, record := range records {
    repo.CreateWithTx(tx, record)
}
tx.Commit()
```

## 📚 相關文件

- [Mock Data README](../mock/README.md) - 使用說明
- [Analytics TDD Roadmap](./ANALYTICS_TDD_ROADMAP.md) - Analytics 功能實作
- [Project Overview](../../.augment/rules/project-overview.md) - 專案概述

## 🎉 總結

成功實作了一個完整的 Mock 資料匯入系統：

- ✅ 28 筆真實資產配置資料成功匯入
- ✅ 支援台股、美股、加密貨幣三種資產類型
- ✅ 保留數量的完整精度（8 位小數）
- ✅ 提供簡單易用的 Make 指令
- ✅ 完整的文檔和使用說明

**現在可以使用真實資料測試所有功能！** 🚀
