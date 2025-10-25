package models

import (
	"time"

	"github.com/google/uuid"
)

// DailyPerformanceSnapshot 每日績效快照
type DailyPerformanceSnapshot struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	SnapshotDate       time.Time `json:"snapshot_date" db:"snapshot_date"`
	TotalMarketValue   float64   `json:"total_market_value" db:"total_market_value"`
	TotalCost          float64   `json:"total_cost" db:"total_cost"`
	TotalUnrealizedPL  float64   `json:"total_unrealized_pl" db:"total_unrealized_pl"`
	TotalUnrealizedPct float64   `json:"total_unrealized_pct" db:"total_unrealized_pct"`
	TotalRealizedPL    float64   `json:"total_realized_pl" db:"total_realized_pl"`
	TotalRealizedPct   float64   `json:"total_realized_pct" db:"total_realized_pct"`
	HoldingCount       int       `json:"holding_count" db:"holding_count"`
	Currency           string    `json:"currency" db:"currency"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// DailyPerformanceSnapshotDetail 每日績效快照明細（按資產類型）
type DailyPerformanceSnapshotDetail struct {
	ID            uuid.UUID `json:"id" db:"id"`
	SnapshotID    uuid.UUID `json:"snapshot_id" db:"snapshot_id"`
	AssetType     AssetType `json:"asset_type" db:"asset_type"`
	MarketValue   float64   `json:"market_value" db:"market_value"`
	Cost          float64   `json:"cost" db:"cost"`
	UnrealizedPL  float64   `json:"unrealized_pl" db:"unrealized_pl"`
	UnrealizedPct float64   `json:"unrealized_pct" db:"unrealized_pct"`
	RealizedPL    float64   `json:"realized_pl" db:"realized_pl"`
	RealizedPct   float64   `json:"realized_pct" db:"realized_pct"`
	HoldingCount  int       `json:"holding_count" db:"holding_count"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// CreateDailyPerformanceSnapshotInput 建立每日績效快照的輸入
type CreateDailyPerformanceSnapshotInput struct {
	SnapshotDate       time.Time                                  `json:"snapshot_date"`
	TotalMarketValue   float64                                    `json:"total_market_value"`
	TotalCost          float64                                    `json:"total_cost"`
	TotalUnrealizedPL  float64                                    `json:"total_unrealized_pl"`
	TotalUnrealizedPct float64                                    `json:"total_unrealized_pct"`
	TotalRealizedPL    float64                                    `json:"total_realized_pl"`
	TotalRealizedPct   float64                                    `json:"total_realized_pct"`
	HoldingCount       int                                        `json:"holding_count"`
	Currency           string                                     `json:"currency"`
	Details            []CreateDailyPerformanceSnapshotDetailInput `json:"details"`
}

// CreateDailyPerformanceSnapshotDetailInput 建立每日績效快照明細的輸入
type CreateDailyPerformanceSnapshotDetailInput struct {
	AssetType     AssetType `json:"asset_type"`
	MarketValue   float64   `json:"market_value"`
	Cost          float64   `json:"cost"`
	UnrealizedPL  float64   `json:"unrealized_pl"`
	UnrealizedPct float64   `json:"unrealized_pct"`
	RealizedPL    float64   `json:"realized_pl"`
	RealizedPct   float64   `json:"realized_pct"`
	HoldingCount  int       `json:"holding_count"`
}

// PerformanceTrendPoint 績效趨勢資料點
type PerformanceTrendPoint struct {
	Date           time.Time `json:"date"`
	MarketValue    float64   `json:"market_value"`
	Cost           float64   `json:"cost"`
	UnrealizedPL   float64   `json:"unrealized_pl"`
	UnrealizedPct  float64   `json:"unrealized_pct"`
	RealizedPL     float64   `json:"realized_pl"`
	RealizedPct    float64   `json:"realized_pct"`
	TotalPL        float64   `json:"total_pl"`        // 總損益 = 已實現 + 未實現
	TotalPct       float64   `json:"total_pct"`       // 總報酬率
	HoldingCount   int       `json:"holding_count"`
}

// PerformanceTrendByType 按資產類型的績效趨勢
type PerformanceTrendByType struct {
	AssetType AssetType               `json:"asset_type"`
	Name      string                  `json:"name"`
	Data      []PerformanceTrendPoint `json:"data"`
}

// PerformanceTrendSummary 績效趨勢摘要
type PerformanceTrendSummary struct {
	StartDate      time.Time                 `json:"start_date"`
	EndDate        time.Time                 `json:"end_date"`
	TotalData      []PerformanceTrendPoint   `json:"total_data"`       // 總體趨勢
	ByType         []PerformanceTrendByType  `json:"by_type"`          // 按資產類型的趨勢
	Currency       string                    `json:"currency"`
	DataPointCount int                       `json:"data_point_count"` // 資料點數量
}

