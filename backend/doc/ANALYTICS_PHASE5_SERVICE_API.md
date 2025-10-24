# Analytics Feature - Phase 5: Analytics Service & API

## 📋 概述

Phase 5 建立了 Analytics Service 和 API Handler，提供已實現損益的分析報表功能，包含摘要、績效分析和最佳表現資產排行。

## ✅ 完成項目

### 1. Analytics Models

**檔案：** `backend/internal/models/analytics.go`

#### TimeRange 時間範圍類型

```go
type TimeRange string

const (
    TimeRangeWeek    TimeRange = "week"    // 本週
    TimeRangeMonth   TimeRange = "month"   // 本月
    TimeRangeQuarter TimeRange = "quarter" // 本季
    TimeRangeYear    TimeRange = "year"    // 本年
    TimeRangeAll     TimeRange = "all"     // 全部
)
```

#### AnalyticsSummary 分析摘要

```go
type AnalyticsSummary struct {
    TotalRealizedPL    float64 `json:"total_realized_pl"`     // 總已實現損益
    TotalRealizedPLPct float64 `json:"total_realized_pl_pct"` // 總已實現損益百分比
    TotalCostBasis     float64 `json:"total_cost_basis"`      // 總成本基礎
    TotalSellAmount    float64 `json:"total_sell_amount"`     // 總賣出金額
    TotalSellFee       float64 `json:"total_sell_fee"`        // 總賣出手續費
    TransactionCount   int     `json:"transaction_count"`     // 交易筆數
    Currency           string  `json:"currency"`              // 幣別
    TimeRange          string  `json:"time_range"`            // 時間範圍
    StartDate          string  `json:"start_date"`            // 起始日期
    EndDate            string  `json:"end_date"`              // 結束日期
}
```

#### PerformanceData 績效資料

```go
type PerformanceData struct {
    AssetType        AssetType `json:"asset_type"`        // 資產類型
    Name             string    `json:"name"`              // 資產類型名稱
    RealizedPL       float64   `json:"realized_pl"`       // 已實現損益
    RealizedPLPct    float64   `json:"realized_pl_pct"`   // 已實現損益百分比
    CostBasis        float64   `json:"cost_basis"`        // 成本基礎
    SellAmount       float64   `json:"sell_amount"`       // 賣出金額
    TransactionCount int       `json:"transaction_count"` // 交易筆數
}
```

#### TopAsset 最佳/最差表現資產

```go
type TopAsset struct {
    Symbol        string    `json:"symbol"`          // 標的代碼
    Name          string    `json:"name"`            // 標的名稱
    AssetType     AssetType `json:"asset_type"`      // 資產類型
    RealizedPL    float64   `json:"realized_pl"`     // 已實現損益
    RealizedPLPct float64   `json:"realized_pl_pct"` // 已實現損益百分比
    CostBasis     float64   `json:"cost_basis"`      // 成本基礎
    SellAmount    float64   `json:"sell_amount"`     // 賣出金額
}
```

### 2. Analytics Service

**檔案：** `backend/internal/service/analytics_service.go`

#### Interface 定義

```go
type AnalyticsService interface {
    GetSummary(timeRange models.TimeRange) (*models.AnalyticsSummary, error)
    GetPerformance(timeRange models.TimeRange) ([]*models.PerformanceData, error)
    GetTopAssets(timeRange models.TimeRange, limit int) ([]*models.TopAsset, error)
}
```

#### GetSummary 實作

- 驗證時間範圍
- 根據時間範圍查詢已實現損益記錄
- 計算總已實現損益、成本基礎、賣出金額等
- 計算總已實現損益百分比

#### GetPerformance 實作

- 按資產類型分組統計
- 計算各資產類型的已實現損益和百分比
- 按已實現損益由高到低排序

#### GetTopAssets 實作

- 按標的分組統計
- 計算各標的的已實現損益和百分比
- 按已實現損益由高到低排序
- 限制回傳數量

### 3. Analytics API Handler

**檔案：** `backend/internal/api/analytics_handler.go`

#### API 端點

1. **GET /api/analytics/summary**

   - 查詢參數：`time_range` (week, month, quarter, year, all)
   - 回傳：`AnalyticsSummary`

2. **GET /api/analytics/performance**

   - 查詢參數：`time_range`
   - 回傳：`[]PerformanceData`

3. **GET /api/analytics/top-assets**
   - 查詢參數：`time_range`, `limit` (預設 5)
   - 回傳：`[]TopAsset`

### 4. Main.go 更新

**檔案：** `backend/cmd/api/main.go`

```go
// 初始化 Analytics Service
analyticsService := service.NewAnalyticsService(realizedProfitRepo)

// 初始化 Handler
analyticsHandler := api.NewAnalyticsHandler(analyticsService)

// 註冊路由
analytics := apiGroup.Group("/analytics")
{
    analytics.GET("/summary", analyticsHandler.GetSummary)
    analytics.GET("/performance", analyticsHandler.GetPerformance)
    analytics.GET("/top-assets", analyticsHandler.GetTopAssets)
}
```

## 📊 測試結果

### Analytics Service 測試

```bash
=== RUN   TestAnalyticsService_GetSummary
--- PASS: TestAnalyticsService_GetSummary (0.00s)
=== RUN   TestAnalyticsService_GetPerformance
--- PASS: TestAnalyticsService_GetPerformance (0.00s)
=== RUN   TestAnalyticsService_GetTopAssets
--- PASS: TestAnalyticsService_GetTopAssets (0.00s)
PASS
ok command-line-arguments 0.313s
```

### Analytics Handler 測試

```bash
=== RUN   TestAnalyticsHandler_GetSummary
--- PASS: TestAnalyticsHandler_GetSummary (0.00s)
=== RUN   TestAnalyticsHandler_GetSummary_InvalidTimeRange
--- PASS: TestAnalyticsHandler_GetSummary_InvalidTimeRange (0.00s)
=== RUN   TestAnalyticsHandler_GetPerformance
--- PASS: TestAnalyticsHandler_GetPerformance (0.00s)
=== RUN   TestAnalyticsHandler_GetTopAssets
--- PASS: TestAnalyticsHandler_GetTopAssets (0.00s)
=== RUN   TestAnalyticsHandler_GetTopAssets_DefaultLimit
--- PASS: TestAnalyticsHandler_GetTopAssets_DefaultLimit (0.00s)
PASS
ok  command-line-arguments  0.427s
```

**✅ 所有測試通過！**

## 🔍 實作細節

### 時間範圍計算

`TimeRange.GetDateRange()` 方法根據不同的時間範圍計算起始和結束日期：

- **week**: 從本週一開始
- **month**: 從本月 1 號開始
- **quarter**: 從本季第一個月的 1 號開始
- **year**: 從本年 1/1 開始
- **all**: 從 2000-01-01 開始

### 資料聚合

**按資產類型聚合（GetPerformance）：**

```go
performanceMap := make(map[models.AssetType]*models.PerformanceData)

for _, record := range records {
    if _, exists := performanceMap[record.AssetType]; !exists {
        performanceMap[record.AssetType] = &models.PerformanceData{
            AssetType: record.AssetType,
            Name:      models.GetAssetTypeName(record.AssetType),
        }
    }

    perf := performanceMap[record.AssetType]
    perf.RealizedPL += record.RealizedPL
    perf.CostBasis += record.CostBasis
    perf.SellAmount += record.SellAmount
    perf.TransactionCount++
}
```

**按標的聚合（GetTopAssets）：**

```go
assetMap := make(map[string]*models.TopAsset)

for _, record := range records {
    if _, exists := assetMap[record.Symbol]; !exists {
        assetMap[record.Symbol] = &models.TopAsset{
            Symbol:    record.Symbol,
            Name:      record.Symbol,
            AssetType: record.AssetType,
        }
    }

    asset := assetMap[record.Symbol]
    asset.RealizedPL += record.RealizedPL
    asset.CostBasis += record.CostBasis
    asset.SellAmount += record.SellAmount
}
```

## 🎯 API 使用範例

### 1. 取得本月分析摘要

```bash
GET /api/analytics/summary?time_range=month

Response:
{
  "total_realized_pl": 12239.0,
  "total_realized_pl_pct": 23.75,
  "total_cost_basis": 51528.0,
  "total_sell_amount": 63800.0,
  "total_sell_fee": 33.0,
  "transaction_count": 2,
  "currency": "TWD",
  "time_range": "month",
  "start_date": "2025-10-01",
  "end_date": "2025-10-31"
}
```

### 2. 取得各資產類型績效

```bash
GET /api/analytics/performance?time_range=month

Response:
[
  {
    "asset_type": "tw-stock",
    "name": "台股",
    "realized_pl": 9930.0,
    "realized_pl_pct": 12.11,
    "cost_basis": 82028.0,
    "sell_amount": 92000.0,
    "transaction_count": 2
  },
  {
    "asset_type": "us-stock",
    "name": "美股",
    "realized_pl": 295.0,
    "realized_pl_pct": 19.67,
    "cost_basis": 1500.0,
    "sell_amount": 1800.0,
    "transaction_count": 1
  }
]
```

### 3. 取得最佳表現資產（Top 5）

```bash
GET /api/analytics/top-assets?time_range=month&limit=5

Response:
[
  {
    "symbol": "BTC",
    "name": "BTC",
    "asset_type": "crypto",
    "realized_pl": 200000.0,
    "realized_pl_pct": 66.67,
    "cost_basis": 300000.0,
    "sell_amount": 500000.0
  },
  {
    "symbol": "2330",
    "name": "2330",
    "asset_type": "tw-stock",
    "realized_pl": 11972.0,
    "realized_pl_pct": 23.93,
    "cost_basis": 50028.0,
    "sell_amount": 62000.0
  }
]
```

## 📝 Phase 5 學習重點

1. **時間範圍處理**：動態計算不同時間範圍的起始和結束日期
2. **資料聚合**：使用 Map 進行分組統計
3. **排序**：使用 `sort.Slice` 進行自訂排序
4. **API 設計**：RESTful API 設計原則
5. **測試驅動開發**：先寫測試，確保功能正確

## 🎯 下一步：前端整合

Phase 5 後端已完成，接下來可以進行前端整合：

1. **建立 API Client**

   - `frontend/src/lib/api/analytics.ts`

2. **建立 Hooks**

   - `frontend/src/hooks/useAnalytics.ts`

3. **更新頁面**
   - 更新 `frontend/src/app/analytics/page.tsx`
   - 移除 Mock 資料依賴
   - 加入 Loading 和錯誤處理

---

**Phase 5 完成時間：** 2025-10-24  
**測試通過率：** 100% (8/8)  
**編譯狀態：** ✅ 成功
