package models

import (
	"time"
)

// Holding 持倉資料
// 代表某個標的（symbol）的當前持倉狀況
// 所有金額欄位（成本、市值、損益）統一以 TWD 計價
type Holding struct {
	Symbol           string    `json:"symbol"`                       // 標的代碼（例如：2330, AAPL, BTC）
	Name             string    `json:"name"`                         // 標的名稱
	AssetType        AssetType `json:"asset_type"`                   // 資產類型
	Quantity         float64   `json:"quantity"`                     // 當前持有數量
	AvgCost          float64   `json:"avg_cost"`                     // FIFO 計算的平均成本（含手續費，TWD）
	AvgCostOriginal  float64   `json:"avg_cost_original"`            // FIFO 計算的平均成本（含手續費，原幣別）
	TotalCost        float64   `json:"total_cost"`                   // 總成本 = AvgCost * Quantity（TWD）
	CurrentPrice     float64   `json:"current_price"`                // 當前市場價格（原始幣別）
	Currency         Currency  `json:"currency"`                     // 價格幣別
	CurrentPriceTWD  float64   `json:"current_price_twd"`            // 當前市場價格（TWD）
	MarketValue      float64   `json:"market_value"`                 // 市值 = CurrentPriceTWD * Quantity（TWD）
	UnrealizedPL     float64   `json:"unrealized_pl"`                // 未實現損益 = MarketValue - TotalCost（TWD）
	UnrealizedPLPct  float64   `json:"unrealized_pl_pct"`            // 未實現損益百分比
	LastUpdated      time.Time `json:"last_updated"`                 // 最後更新時間
	PriceSource      string    `json:"price_source,omitempty"`       // 價格來源（cache, api, stale-cache）
	IsPriceStale     bool      `json:"is_price_stale,omitempty"`     // 價格是否過期
	PriceStaleReason string    `json:"price_stale_reason,omitempty"` // 價格過期原因
}

// CostBatch FIFO 成本批次
// 用於追蹤每一批買入的成本，實作 FIFO 計算
// 成本統一以 TWD 計價
type CostBatch struct {
	Date             time.Time `json:"date"`               // 買入日期
	Quantity         float64   `json:"quantity"`           // 該批次剩餘數量
	UnitCost         float64   `json:"unit_cost"`          // 單位成本（含手續費，TWD）
	UnitCostOriginal float64   `json:"unit_cost_original"` // 單位成本（含手續費，原幣別）
	OriginalQty      float64   `json:"original_qty"`       // 原始買入數量
	Currency         Currency  `json:"currency"`           // 原始交易幣別
	ExchangeRate     float64   `json:"exchange_rate"`      // 交易時的匯率（TWD/原幣別）
}

// HoldingFilters 持倉篩選條件
type HoldingFilters struct {
	AssetType *AssetType `json:"asset_type,omitempty"` // 按資產類型篩選
	Symbol    *string    `json:"symbol,omitempty"`     // 按標的代碼篩選
}

// CorporateActionType 公司行動類型（股票分割/合併）
type CorporateActionType string

const (
	ActionTypeSplit CorporateActionType = "split" // 股票分割（例如 1:2）
	ActionTypeMerge CorporateActionType = "merge" // 股票合併（例如 2:1）
)

// CorporateAction 公司行動記錄
// 用於處理股票分割、合併等事件
type CorporateAction struct {
	Symbol string              `json:"symbol"` // 標的代碼
	Type   CorporateActionType `json:"type"`   // 行動類型
	Date   time.Time           `json:"date"`   // 生效日期
	Ratio  float64             `json:"ratio"`  // 比例（例如 1:2 分割 = 2.0, 2:1 合併 = 0.5）
}

// Validate 驗證 CorporateActionType 是否有效
func (c CorporateActionType) Validate() bool {
	switch c {
	case ActionTypeSplit, ActionTypeMerge:
		return true
	}
	return false
}

