package models

import "time"

// TimeRange 時間範圍類型
type TimeRange string

const (
	TimeRangeWeek    TimeRange = "week"    // 本週
	TimeRangeMonth   TimeRange = "month"   // 本月
	TimeRangeQuarter TimeRange = "quarter" // 本季
	TimeRangeYear    TimeRange = "year"    // 本年
	TimeRangeAll     TimeRange = "all"     // 全部
)

// Validate 驗證時間範圍是否有效
func (tr TimeRange) Validate() bool {
	switch tr {
	case TimeRangeWeek, TimeRangeMonth, TimeRangeQuarter, TimeRangeYear, TimeRangeAll:
		return true
	default:
		return false
	}
}

// GetDateRange 根據時間範圍取得起始和結束日期
func (tr TimeRange) GetDateRange() (startDate, endDate time.Time) {
	now := time.Now()
	endDate = now

	switch tr {
	case TimeRangeWeek:
		// 本週：從週一開始
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // 週日視為第 7 天
		}
		startDate = now.AddDate(0, 0, -(weekday - 1))
	case TimeRangeMonth:
		// 本月：從 1 號開始
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case TimeRangeQuarter:
		// 本季：從季度第一個月的 1 號開始
		month := now.Month()
		quarterStartMonth := ((month-1)/3)*3 + 1
		startDate = time.Date(now.Year(), quarterStartMonth, 1, 0, 0, 0, 0, now.Location())
	case TimeRangeYear:
		// 本年：從 1/1 開始
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	case TimeRangeAll:
		// 全部：從很久以前開始（2000-01-01）
		startDate = time.Date(2000, 1, 1, 0, 0, 0, 0, now.Location())
	}

	return startDate, endDate
}

// AnalyticsSummary 分析摘要資料
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

// PerformanceData 績效資料（按資產類型）
type PerformanceData struct {
	AssetType       AssetType `json:"asset_type"`        // 資產類型
	Name            string    `json:"name"`              // 資產類型名稱
	RealizedPL      float64   `json:"realized_pl"`       // 已實現損益
	RealizedPLPct   float64   `json:"realized_pl_pct"`   // 已實現損益百分比
	CostBasis       float64   `json:"cost_basis"`        // 成本基礎
	SellAmount      float64   `json:"sell_amount"`       // 賣出金額
	TransactionCount int      `json:"transaction_count"` // 交易筆數
}

// TopAsset 最佳/最差表現資產
type TopAsset struct {
	Symbol        string    `json:"symbol"`          // 標的代碼
	Name          string    `json:"name"`            // 標的名稱
	AssetType     AssetType `json:"asset_type"`      // 資產類型
	RealizedPL    float64   `json:"realized_pl"`     // 已實現損益
	RealizedPLPct float64   `json:"realized_pl_pct"` // 已實現損益百分比
	CostBasis     float64   `json:"cost_basis"`      // 成本基礎
	SellAmount    float64   `json:"sell_amount"`     // 賣出金額
}

// AssetTypeNameMap 資產類型名稱對應
var AssetTypeNameMap = map[AssetType]string{
	AssetTypeTWStock: "台股",
	AssetTypeUSStock: "美股",
	AssetTypeCrypto:  "加密貨幣",
}

// GetAssetTypeName 取得資產類型的顯示名稱
func GetAssetTypeName(assetType AssetType) string {
	if name, ok := AssetTypeNameMap[assetType]; ok {
		return name
	}
	return string(assetType)
}

