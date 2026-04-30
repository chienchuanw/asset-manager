package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
)

// HoldingReconciliationService 提供持倉對帳業務邏輯
type HoldingReconciliationService interface {
	// ReconcileHoldings 計算每個 item 的 diff；dryRun=true 僅回傳預覽，
	// dryRun=false 在單一 DB 交易內寫入 adjustment（不更動 cash_flows / realized_profits）。
	ReconcileHoldings(ctx context.Context, items []models.HoldingReconcileItem, dryRun bool) (*models.HoldingReconcilePreview, error)
}

type holdingReconciliationService struct {
	db     *sql.DB
	txRepo repository.TransactionRepository
	fifo   FIFOCalculator
}

// NewHoldingReconciliationService 建立 HoldingReconciliationService
func NewHoldingReconciliationService(db *sql.DB, txRepo repository.TransactionRepository, fifo FIFOCalculator) HoldingReconciliationService {
	return &holdingReconciliationService{db: db, txRepo: txRepo, fifo: fifo}
}

// ReconcileHoldings 詳見介面註解
func (s *holdingReconciliationService) ReconcileHoldings(
	ctx context.Context,
	items []models.HoldingReconcileItem,
	dryRun bool,
) (*models.HoldingReconcilePreview, error) {
	batch := models.HoldingReconcileBatchInput{Items: items}
	if err := batch.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	preview := &models.HoldingReconcilePreview{Items: make([]models.HoldingReconcilePreviewItem, 0, len(items))}
	now := time.Now()

	for _, item := range items {
		txs, existingCurrency, err := s.lockSymbolTxs(ctx, tx, item.AssetType, item.Symbol)
		if err != nil {
			return nil, err
		}

		if existingCurrency != "" && existingCurrency != item.Currency {
			return nil, models.ErrReconcileCurrencyMismatch
		}

		var prevQty, prevAvg float64
		var existingName string
		if len(txs) > 0 {
			h, err := s.fifo.CalculateHoldingForSymbol(item.Symbol, txs)
			if err != nil {
				return nil, fmt.Errorf("compute current holding: %w", err)
			}
			if h != nil {
				prevQty = h.Quantity
				prevAvg = h.AvgCostOriginal
				existingName = h.Name
			} else {
				existingName = txs[0].Name
			}
		}

		action := classifyReconcileAction(prevQty, item.TargetQuantity, prevAvg, item.TargetAvgCost)

		if action == models.ReconcileActionCreate && item.Name == "" {
			return nil, models.ErrReconcileNameRequired
		}

		preview.Items = append(preview.Items, models.HoldingReconcilePreviewItem{
			Item:          item,
			PrevQuantity:  prevQty,
			PrevAvgCost:   prevAvg,
			QuantityDelta: item.TargetQuantity - prevQty,
			AvgCostDelta:  item.TargetAvgCost - prevAvg,
			Action:        action,
		})

		if dryRun || action == models.ReconcileActionNoop {
			continue
		}

		name := item.Name
		if name == "" {
			name = existingName
		}

		prevQtyCopy := prevQty
		prevAvgCopy := prevAvg
		var reasonPtr *string
		if item.Reason != "" {
			r := item.Reason
			reasonPtr = &r
		}

		row := &models.Transaction{
			Date:                   now,
			AssetType:              item.AssetType,
			Symbol:                 item.Symbol,
			Name:                   name,
			TransactionType:        models.TransactionTypeAdjustment,
			Quantity:               item.TargetQuantity,
			Price:                  item.TargetAvgCost,
			Amount:                 0,
			Currency:               item.Currency,
			AdjustmentPrevQuantity: &prevQtyCopy,
			AdjustmentPrevAvgCost:  &prevAvgCopy,
			AdjustmentReason:       reasonPtr,
		}
		if _, err := s.txRepo.CreateAdjustmentTx(tx, row); err != nil {
			return nil, fmt.Errorf("write adjustment: %w", err)
		}
	}

	if !dryRun {
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit: %w", err)
		}
	}
	return preview, nil
}

// lockSymbolTxs 取得並鎖定指定 (asset_type, symbol) 的所有交易；同時回傳既有 currency。
func (s *holdingReconciliationService) lockSymbolTxs(
	ctx context.Context,
	tx *sql.Tx,
	assetType models.AssetType,
	symbol string,
) ([]*models.Transaction, models.Currency, error) {
	const q = `
		SELECT id, date, asset_type, symbol, name, transaction_type, quantity, price, amount, fee, tax, currency, exchange_rate_id, note,
		       adjustment_prev_quantity, adjustment_prev_avg_cost, adjustment_reason, broker_account_id,
		       created_at, updated_at
		FROM transactions
		WHERE asset_type = $1 AND symbol = $2
		ORDER BY date ASC, created_at ASC
		FOR UPDATE`
	rows, err := tx.QueryContext(ctx, q, assetType, symbol)
	if err != nil {
		return nil, "", fmt.Errorf("lock symbol txs: %w", err)
	}
	defer rows.Close()

	var out []*models.Transaction
	var currency models.Currency
	for rows.Next() {
		t := &models.Transaction{}
		if err := rows.Scan(
			&t.ID, &t.Date, &t.AssetType, &t.Symbol, &t.Name, &t.TransactionType,
			&t.Quantity, &t.Price, &t.Amount, &t.Fee, &t.Tax, &t.Currency, &t.ExchangeRateID, &t.Note,
			&t.AdjustmentPrevQuantity, &t.AdjustmentPrevAvgCost, &t.AdjustmentReason, &t.BrokerAccountID,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, "", err
		}
		if currency == "" {
			currency = t.Currency
		}
		out = append(out, t)
	}
	return out, currency, rows.Err()
}

func classifyReconcileAction(prevQty, targetQty, prevAvg, targetAvg float64) models.ReconcileAction {
	if prevQty == 0 && targetQty > 0 {
		return models.ReconcileActionCreate
	}
	if prevQty > 0 && targetQty == 0 {
		return models.ReconcileActionLiquidate
	}
	if prevQty == targetQty && prevAvg == targetAvg {
		return models.ReconcileActionNoop
	}
	return models.ReconcileActionUpdate
}

// IsReconcileValidationError 回報 err 是否屬於可預期的對帳輸入驗證錯誤（供 handler 映射 400）
func IsReconcileValidationError(err error) bool {
	return errors.Is(err, models.ErrInvalidReconcileQuantity) ||
		errors.Is(err, models.ErrInvalidReconcileAvgCost) ||
		errors.Is(err, models.ErrInvalidReconcileCurrency) ||
		errors.Is(err, models.ErrInvalidReconcileAssetType) ||
		errors.Is(err, models.ErrReconcileNameRequired) ||
		errors.Is(err, models.ErrReconcileCurrencyMismatch) ||
		errors.Is(err, models.ErrReconcileEmptyBatch)
}
