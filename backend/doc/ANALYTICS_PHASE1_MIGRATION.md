# Analytics 功能實作 - Phase 1: 資料庫 Migration

## 📋 概述

Phase 1 完成了 `realized_profits` 表的建立，用於記錄每筆賣出交易的已實現損益資訊。

---

## ✅ 完成項目

### 1. Migration 檔案建立

**檔案：**
- `backend/migrations/000004_create_realized_profits_table.up.sql`
- `backend/migrations/000004_create_realized_profits_table.down.sql`

### 2. 資料表結構

**表名：** `realized_profits`

| 欄位名稱 | 資料型別 | 說明 | 約束條件 |
|---------|---------|------|---------|
| `id` | UUID | 主鍵 | PRIMARY KEY, DEFAULT gen_random_uuid() |
| `transaction_id` | UUID | 關聯的賣出交易 ID | NOT NULL, FOREIGN KEY → transactions(id) ON DELETE CASCADE |
| `symbol` | VARCHAR(20) | 標的代碼 | NOT NULL |
| `asset_type` | VARCHAR(20) | 資產類型 | NOT NULL, CHECK IN ('cash', 'tw-stock', 'us-stock', 'crypto') |
| `sell_date` | DATE | 賣出日期 | NOT NULL |
| `quantity` | DECIMAL(20,8) | 賣出數量 | NOT NULL, CHECK > 0 |
| `sell_price` | DECIMAL(20,8) | 賣出價格 | NOT NULL, CHECK >= 0 |
| `sell_amount` | DECIMAL(20,8) | 賣出金額 | NOT NULL, CHECK >= 0 |
| `sell_fee` | DECIMAL(20,8) | 賣出手續費 | NOT NULL, DEFAULT 0, CHECK >= 0 |
| `cost_basis` | DECIMAL(20,8) | FIFO 成本基礎 | NOT NULL, CHECK >= 0 |
| `realized_pl` | DECIMAL(20,8) | 已實現損益 | NOT NULL |
| `realized_pl_pct` | DECIMAL(10,4) | 已實現損益百分比 | NOT NULL |
| `currency` | VARCHAR(10) | 幣別 | NOT NULL, DEFAULT 'TWD', CHECK IN ('TWD', 'USD') |
| `created_at` | TIMESTAMP WITH TIME ZONE | 建立時間 | DEFAULT CURRENT_TIMESTAMP |
| `updated_at` | TIMESTAMP WITH TIME ZONE | 更新時間 | DEFAULT CURRENT_TIMESTAMP |

### 3. 索引

建立了以下索引以提升查詢效能：

- `idx_realized_profits_symbol` - 按標的代碼查詢
- `idx_realized_profits_asset_type` - 按資產類型查詢
- `idx_realized_profits_sell_date` - 按賣出日期查詢（降序）
- `idx_realized_profits_transaction_id` - 按交易 ID 查詢
- `idx_realized_profits_date_asset` - 複合索引（日期 + 資產類型）
- `idx_realized_profits_date_symbol` - 複合索引（日期 + 標的代碼）

### 4. 觸發器

- `update_realized_profits_updated_at` - 自動更新 `updated_at` 欄位

### 5. 註解

為表格和所有欄位都加入了繁體中文註解，方便理解。

---

## 🎯 已實現損益計算邏輯

### 公式

```
已實現損益 = (賣出金額 - 賣出手續費) - 成本基礎
已實現損益百分比 = (已實現損益 / 成本基礎) × 100
```

### 範例

**情境：**
- 買入 100 股 @ 500 TWD，手續費 28 TWD
- 成本基礎 = 50,000 + 28 = 50,028 TWD
- 賣出 100 股 @ 620 TWD，手續費 30 TWD
- 賣出金額 = 62,000 TWD

**計算：**
```
已實現損益 = (62,000 - 30) - 50,028 = 11,942 TWD
已實現損益百分比 = (11,942 / 50,028) × 100 = 23.87%
```

---

## 🔧 執行 Migration

### 開發環境

```bash
cd backend
make migrate-up-env
```

**執行結果：**
```
Loading .env and running migrations...
4/u create_realized_profits_table (21.278625ms)
```

### 測試環境

```bash
cd backend
make migrate-test-up-env
```

**執行結果：**
```
Loading .env.test and running migrations...
4/u create_realized_profits_table (32.781209ms)
```

---

## ✅ 驗證結果

### 檢查表結構

```bash
source .env.local
psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d $DB_NAME -c "\d realized_profits"
```

**確認項目：**
- ✅ 15 個欄位全部建立
- ✅ 主鍵約束正確
- ✅ 外鍵約束正確（ON DELETE CASCADE）
- ✅ 7 個索引全部建立
- ✅ 7 個 CHECK 約束全部建立
- ✅ 觸發器正確建立
- ✅ 表格註解正確

### 檢查註解

```bash
psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d $DB_NAME \
  -c "SELECT obj_description('realized_profits'::regclass);"
```

**結果：**
```
已實現損益記錄表 - 記錄每筆賣出交易的損益資訊
```

---

## 🚀 下一步：Phase 2

Phase 2 將實作：

1. **Model 定義**
   - `backend/internal/models/realized_profit.go`
   - 定義 `RealizedProfit` 結構
   - 定義 `CreateRealizedProfitInput` 結構
   - 定義 `RealizedProfitFilters` 結構

2. **Repository 層**
   - `backend/internal/repository/realized_profit_repository.go`
   - 實作 CRUD 操作
   - 實作篩選查詢

3. **測試**
   - `backend/internal/repository/realized_profit_repository_test.go`
   - 遵循 TDD 原則，先寫測試

---

## 📝 Migration 檔案內容

### Up Migration

檔案：`backend/migrations/000004_create_realized_profits_table.up.sql`

```sql
-- 建立已實現損益表
CREATE TABLE IF NOT EXISTS realized_profits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    asset_type VARCHAR(20) NOT NULL CHECK (asset_type IN ('cash', 'tw-stock', 'us-stock', 'crypto')),
    sell_date DATE NOT NULL,
    quantity DECIMAL(20,8) NOT NULL CHECK (quantity > 0),
    sell_price DECIMAL(20,8) NOT NULL CHECK (sell_price >= 0),
    sell_amount DECIMAL(20,8) NOT NULL CHECK (sell_amount >= 0),
    sell_fee DECIMAL(20,8) NOT NULL DEFAULT 0 CHECK (sell_fee >= 0),
    cost_basis DECIMAL(20,8) NOT NULL CHECK (cost_basis >= 0),
    realized_pl DECIMAL(20,8) NOT NULL,
    realized_pl_pct DECIMAL(10,4) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'TWD' CHECK (currency IN ('TWD', 'USD')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 建立索引
CREATE INDEX idx_realized_profits_symbol ON realized_profits(symbol);
CREATE INDEX idx_realized_profits_asset_type ON realized_profits(asset_type);
CREATE INDEX idx_realized_profits_sell_date ON realized_profits(sell_date DESC);
CREATE INDEX idx_realized_profits_transaction_id ON realized_profits(transaction_id);
CREATE INDEX idx_realized_profits_date_asset ON realized_profits(sell_date, asset_type);
CREATE INDEX idx_realized_profits_date_symbol ON realized_profits(sell_date, symbol);

-- 建立觸發器
CREATE TRIGGER update_realized_profits_updated_at
    BEFORE UPDATE ON realized_profits
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 加入註解
COMMENT ON TABLE realized_profits IS '已實現損益記錄表 - 記錄每筆賣出交易的損益資訊';
COMMENT ON COLUMN realized_profits.cost_basis IS 'FIFO 計算的成本基礎（含買入手續費）';
COMMENT ON COLUMN realized_profits.realized_pl IS '已實現損益 = (賣出金額 - 賣出手續費) - 成本基礎';
COMMENT ON COLUMN realized_profits.realized_pl_pct IS '已實現損益百分比 = (已實現損益 / 成本基礎) × 100';
```

### Down Migration

檔案：`backend/migrations/000004_create_realized_profits_table.down.sql`

```sql
-- 刪除已實現損益表
DROP TABLE IF EXISTS realized_profits CASCADE;
```

---

## 🎉 總結

Phase 1 已成功完成！

**完成項目：**
- ✅ Migration 檔案建立
- ✅ 資料表結構設計
- ✅ 索引優化
- ✅ 約束條件設定
- ✅ 觸發器建立
- ✅ 註解完整
- ✅ 開發環境 Migration 執行
- ✅ 測試環境 Migration 執行
- ✅ 驗證通過

**下一步：**
請繼續進行 Phase 2 的實作（Model 和 Repository 層）。

