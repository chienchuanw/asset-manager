# 手動建立資產快照指令

## 概述

`make snapshot` 指令可以手動觸發建立當日的資產價值快照，並將資料儲存到以下兩個資料表：

1. **`asset_snapshots`** - 各類資產的市值快照（台股、美股、加密貨幣、總資產）
2. **`daily_performance_snapshots`** - 每日績效快照（包含已實現/未實現損益）

## 使用方式

```bash
cd backend
make snapshot
```

## 執行流程

此指令會依序執行以下步驟：

### 1. 更新今日匯率
- 從 ExchangeRate-API 取得最新的 USD/TWD 匯率
- 如果 API 失敗，會使用快取或預設匯率（31.5）繼續執行

### 2. 建立資產快照
建立四種類型的快照到 `asset_snapshots` 資料表：
- **total** - 總資產價值
- **tw-stock** - 台股總價值
- **us-stock** - 美股總價值（已轉換為 TWD）
- **crypto** - 加密貨幣總價值（已轉換為 TWD）

### 3. 建立績效快照
建立績效快照到 `daily_performance_snapshots` 資料表，包含：
- 總市值、總成本
- 未實現損益（金額與百分比）
- 已實現損益（金額與百分比）
- 持倉數量
- 各資產類型的明細資料

## 自動更新機制

如果當日（`snapshot_date`）已經有資料，系統會**自動更新**該筆資料，而不是建立新的記錄。

這是透過以下機制實現的：

### asset_snapshots 表
```go
// 檢查是否已存在今日快照
existing, err := s.repo.GetByDateAndType(today, snapshot.assetType)
if err == nil && existing != nil {
    // 已存在，更新
    _, err = s.repo.Update(today, snapshot.assetType, snapshot.value)
} else {
    // 不存在，建立新的
    _, err = s.repo.Create(input)
}
```

### daily_performance_snapshots 表
```sql
INSERT INTO daily_performance_snapshots (...)
VALUES (...)
ON CONFLICT (snapshot_date) DO UPDATE SET
    total_market_value = EXCLUDED.total_market_value,
    total_cost = EXCLUDED.total_cost,
    ...
```

## 輸出範例

```
✓ Database connected
✓ Using real price API (FinMind + CoinGecko + Alpha Vantage)
✓ Services initialized

📊 Step 1: Refreshing today's exchange rate...
✓ Exchange rate refreshed successfully

📊 Step 2: Creating asset snapshots...
✓ Asset snapshots created successfully

📊 Step 3: Creating performance snapshot...
✓ Performance snapshot created successfully

============================================================
📈 Snapshot Summary
============================================================
Date:              2025-01-13
Total Market Value: 1234567.89 TWD
Total Cost:         1000000.00 TWD
Unrealized P/L:     234567.89 TWD (23.46%)
Realized P/L:       50000.00 TWD (5.00%)
Holdings Count:     15
============================================================

✅ All snapshots created successfully!
```

## 注意事項

1. **環境變數**：需要正確設定 `.env.local` 檔案，包含資料庫連線資訊
2. **API Keys**：
   - 如果有設定 `FINMIND_API_KEY`、`COINGECKO_API_KEY`、`ALPHA_VANTAGE_API_KEY`，會使用真實價格 API
   - 如果沒有設定，會使用 Mock 價格服務
3. **執行時機**：可以在任何時間執行，系統會自動處理重複執行的情況
4. **資料一致性**：建議在市場收盤後執行，以確保價格資料的準確性

## 相關檔案

- **CLI 工具**：`backend/cmd/snapshot/main.go`
- **Makefile**：`backend/Makefile`
- **服務層**：
  - `backend/internal/service/asset_snapshot_service.go`
  - `backend/internal/service/performance_trend_service.go`
- **Repository 層**：
  - `backend/internal/repository/asset_snapshot_repository.go`
  - `backend/internal/repository/performance_snapshot_repository.go`

## 排程執行

如果需要自動化執行，可以使用系統內建的排程器：

```bash
# 在 .env.local 中設定
SNAPSHOT_SCHEDULER_ENABLED=true
SCHEDULER_SNAPSHOT_TIME=23:59  # 每天 23:59 執行
```

或使用系統的 cron job：

```bash
# 每天 23:59 執行
59 23 * * * cd /path/to/backend && make snapshot >> /var/log/snapshot.log 2>&1
```

