package models

import (
	"time"
)

// RealizedProfit 已實現損益記錄
type RealizedProfit struct {
	ID            string    `json:"id" db:"id"`
	TransactionID string    `json:"transaction_id" db:"transaction_id"`
	Symbol        string    `json:"symbol" db:"symbol"`
	AssetType     AssetType `json:"asset_type" db:"asset_type"`
	SellDate      time.Time `json:"sell_date" db:"sell_date"`
	Quantity      float64   `json:"quantity" db:"quantity"`
	SellPrice     float64   `json:"sell_price" db:"sell_price"`
	SellAmount    float64   `json:"sell_amount" db:"sell_amount"`
	SellFee       float64   `json:"sell_fee" db:"sell_fee"`
	CostBasis     float64   `json:"cost_basis" db:"cost_basis"`
	RealizedPL    float64   `json:"realized_pl" db:"realized_pl"`
	RealizedPLPct float64   `json:"realized_pl_pct" db:"realized_pl_pct"`
	Currency      string    `json:"currency" db:"currency"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// CreateRealizedProfitInput 建立已實現損益的輸入
type CreateRealizedProfitInput struct {
	TransactionID string    `json:"transaction_id"`
	Symbol        string    `json:"symbol"`
	AssetType     AssetType `json:"asset_type"`
	SellDate      time.Time `json:"sell_date"`
	Quantity      float64   `json:"quantity"`
	SellPrice     float64   `json:"sell_price"`
	SellAmount    float64   `json:"sell_amount"`
	SellFee       float64   `json:"sell_fee"`
	CostBasis     float64   `json:"cost_basis"`
	Currency      string    `json:"currency"`
}

// RealizedProfitFilters 已實現損益查詢篩選條件
type RealizedProfitFilters struct {
	AssetType *AssetType `json:"asset_type,omitempty"`
	Symbol    *string    `json:"symbol,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

