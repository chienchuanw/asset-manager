package models_test

import (
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHoldingReconcileItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    models.HoldingReconcileItem
		wantErr error
	}{
		{
			name: "valid update",
			item: models.HoldingReconcileItem{
				AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD,
				TargetQuantity: 150, TargetAvgCost: 180.5,
			},
		},
		{
			name: "valid liquidate (zero qty + zero avg cost)",
			item: models.HoldingReconcileItem{
				AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD,
				TargetQuantity: 0, TargetAvgCost: 0,
			},
		},
		{
			name: "negative quantity rejected",
			item: models.HoldingReconcileItem{
				AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD,
				TargetQuantity: -1, TargetAvgCost: 100,
			},
			wantErr: models.ErrInvalidReconcileQuantity,
		},
		{
			name: "negative avg cost rejected",
			item: models.HoldingReconcileItem{
				AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD,
				TargetQuantity: 10, TargetAvgCost: -1,
			},
			wantErr: models.ErrInvalidReconcileAvgCost,
		},
		{
			name: "zero avg cost on non-liquidate rejected",
			item: models.HoldingReconcileItem{
				AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD,
				TargetQuantity: 10, TargetAvgCost: 0,
			},
			wantErr: models.ErrInvalidReconcileAvgCost,
		},
		{
			name: "missing currency rejected",
			item: models.HoldingReconcileItem{
				AssetType: models.AssetTypeUSStock, Symbol: "AAPL",
				TargetQuantity: 10, TargetAvgCost: 50,
			},
			wantErr: models.ErrInvalidReconcileCurrency,
		},
		{
			name: "cash asset type rejected",
			item: models.HoldingReconcileItem{
				AssetType: models.AssetTypeCash, Symbol: "AAPL", Currency: models.CurrencyUSD,
				TargetQuantity: 10, TargetAvgCost: 50,
			},
			wantErr: models.ErrInvalidReconcileAssetType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestHoldingReconcileBatchInput_Validate(t *testing.T) {
	t.Run("empty batch rejected", func(t *testing.T) {
		batch := models.HoldingReconcileBatchInput{}
		assert.ErrorIs(t, batch.Validate(), models.ErrReconcileEmptyBatch)
	})

	t.Run("propagates first item error", func(t *testing.T) {
		batch := models.HoldingReconcileBatchInput{Items: []models.HoldingReconcileItem{
			{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD, TargetQuantity: 100, TargetAvgCost: 150},
			{AssetType: models.AssetTypeUSStock, Symbol: "TSLA", Currency: models.CurrencyUSD, TargetQuantity: -1, TargetAvgCost: 100},
		}}
		assert.ErrorIs(t, batch.Validate(), models.ErrInvalidReconcileQuantity)
	})

	t.Run("all valid", func(t *testing.T) {
		batch := models.HoldingReconcileBatchInput{Items: []models.HoldingReconcileItem{
			{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD, TargetQuantity: 100, TargetAvgCost: 150},
		}}
		assert.NoError(t, batch.Validate())
	})
}

func TestTransactionType_AdjustmentValidates(t *testing.T) {
	assert.True(t, models.TransactionTypeAdjustment.Validate())
}
