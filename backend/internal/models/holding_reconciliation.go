package models

import "errors"

// 持倉對帳的 sentinel errors（handler 透過 errors.Is 映射至 HTTP 狀態碼）
var (
	ErrInvalidReconcileQuantity  = errors.New("invalid reconcile quantity")
	ErrInvalidReconcileAvgCost   = errors.New("invalid reconcile avg cost")
	ErrInvalidReconcileCurrency  = errors.New("invalid reconcile currency")
	ErrInvalidReconcileAssetType = errors.New("invalid reconcile asset type")
	ErrReconcileNameRequired     = errors.New("name required for first-time reconciliation")
	ErrReconcileEmptyBatch       = errors.New("reconcile batch is empty")
	ErrReconcileCurrencyMismatch = errors.New("currency mismatch with existing holding")
)

// HoldingReconcileItem 描述單一 (asset_type, symbol, currency) 的目標狀態
type HoldingReconcileItem struct {
	AssetType      AssetType `json:"asset_type" binding:"required"`
	Symbol         string    `json:"symbol" binding:"required"`
	Name           string    `json:"name"`
	Currency       Currency  `json:"currency" binding:"required"`
	TargetQuantity float64   `json:"target_quantity"`
	TargetAvgCost  float64   `json:"target_avg_cost"`
	Reason         string    `json:"reason"`
}

// Validate 檢查單一 reconcile item 的欄位範圍
func (i *HoldingReconcileItem) Validate() error {
	switch i.AssetType {
	case AssetTypeTWStock, AssetTypeUSStock, AssetTypeCrypto:
	default:
		return ErrInvalidReconcileAssetType
	}
	if !i.Currency.Validate() {
		return ErrInvalidReconcileCurrency
	}
	if i.TargetQuantity < 0 {
		return ErrInvalidReconcileQuantity
	}
	if i.TargetAvgCost < 0 {
		return ErrInvalidReconcileAvgCost
	}
	// avg_cost == 0 僅在清倉（qty == 0）時合法
	if i.TargetAvgCost == 0 && i.TargetQuantity > 0 {
		return ErrInvalidReconcileAvgCost
	}
	return nil
}

// ReconcileAction 對帳結果分類
type ReconcileAction string

const (
	ReconcileActionCreate    ReconcileAction = "create"
	ReconcileActionUpdate    ReconcileAction = "update"
	ReconcileActionLiquidate ReconcileAction = "liquidate"
	ReconcileActionNoop      ReconcileAction = "noop"
)

// HoldingReconcilePreviewItem 單筆對帳的 diff 預覽
type HoldingReconcilePreviewItem struct {
	Item          HoldingReconcileItem `json:"item"`
	PrevQuantity  float64              `json:"prev_quantity"`
	PrevAvgCost   float64              `json:"prev_avg_cost"`
	QuantityDelta float64              `json:"quantity_delta"`
	AvgCostDelta  float64              `json:"avg_cost_delta"`
	Action        ReconcileAction      `json:"action"`
}

// HoldingReconcileBatchInput 批次對帳請求 body
type HoldingReconcileBatchInput struct {
	Items []HoldingReconcileItem `json:"items" binding:"required"`
}

// Validate 檢查整個批次（空陣列拒絕，逐項驗證）
func (b *HoldingReconcileBatchInput) Validate() error {
	if len(b.Items) == 0 {
		return ErrReconcileEmptyBatch
	}
	for idx := range b.Items {
		if err := b.Items[idx].Validate(); err != nil {
			return err
		}
	}
	return nil
}

// HoldingReconcilePreview 批次對帳預覽結果
type HoldingReconcilePreview struct {
	Items []HoldingReconcilePreviewItem `json:"items"`
}
