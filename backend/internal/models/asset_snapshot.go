package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SnapshotAssetType 快照資產類型
type SnapshotAssetType string

const (
	SnapshotAssetTypeTWStock SnapshotAssetType = "tw-stock"
	SnapshotAssetTypeUSStock SnapshotAssetType = "us-stock"
	SnapshotAssetTypeCrypto  SnapshotAssetType = "crypto"
	SnapshotAssetTypeTotal   SnapshotAssetType = "total"
)

// AssetSnapshot 資產快照模型
type AssetSnapshot struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	SnapshotDate time.Time         `json:"snapshot_date" db:"snapshot_date"`
	AssetType    SnapshotAssetType `json:"asset_type" db:"asset_type"`
	ValueTWD     float64           `json:"value_twd" db:"value_twd"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
}

// CreateAssetSnapshotInput 建立資產快照的輸入
type CreateAssetSnapshotInput struct {
	SnapshotDate time.Time         `json:"snapshot_date"`
	AssetType    SnapshotAssetType `json:"asset_type"`
	ValueTWD     float64           `json:"value_twd"`
}

// AssetSnapshotFilters 資產快照篩選條件
type AssetSnapshotFilters struct {
	StartDate *time.Time         `json:"start_date,omitempty"`
	EndDate   *time.Time         `json:"end_date,omitempty"`
	AssetType *SnapshotAssetType `json:"asset_type,omitempty"`
}

// Validate 驗證建立資產快照的輸入
func (input *CreateAssetSnapshotInput) Validate() error {
	if input.SnapshotDate.IsZero() {
		return fmt.Errorf("snapshot_date is required")
	}

	if input.AssetType == "" {
		return fmt.Errorf("asset_type is required")
	}

	// 驗證資產類型
	validTypes := map[SnapshotAssetType]bool{
		SnapshotAssetTypeTWStock: true,
		SnapshotAssetTypeUSStock: true,
		SnapshotAssetTypeCrypto:  true,
		SnapshotAssetTypeTotal:   true,
	}
	if !validTypes[input.AssetType] {
		return fmt.Errorf("invalid asset_type")
	}

	if input.ValueTWD < 0 {
		return fmt.Errorf("value_twd must be non-negative")
	}

	return nil
}

