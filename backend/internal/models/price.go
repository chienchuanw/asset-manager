package models

import (
	"time"
)

// Price 價格資料
type Price struct {
	Symbol      string    `json:"symbol"`       // 標的代碼
	AssetType   AssetType `json:"asset_type"`   // 資產類型
	Price       float64   `json:"price"`        // 價格
	Currency    string    `json:"currency"`     // 幣別（TWD, USD）
	Source      string    `json:"source"`       // 資料來源（例如：cache, api）
	UpdatedAt   time.Time `json:"updated_at"`   // 更新時間
}

// PriceCache 價格快取資料（用於 Redis）
type PriceCache struct {
	Symbol    string    `json:"symbol"`
	AssetType AssetType `json:"asset_type"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	CachedAt  time.Time `json:"cached_at"`
}

