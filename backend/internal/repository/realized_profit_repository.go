package repository

import (
	"database/sql"
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// RealizedProfitRepository 已實現損益資料存取介面
type RealizedProfitRepository interface {
	// Create 建立已實現損益記錄
	Create(input *models.CreateRealizedProfitInput) (*models.RealizedProfit, error)

	// GetByTransactionID 根據交易 ID 取得已實現損益
	GetByTransactionID(transactionID string) (*models.RealizedProfit, error)

	// GetAll 取得所有已實現損益記錄（支援篩選）
	GetAll(filters models.RealizedProfitFilters) ([]*models.RealizedProfit, error)

	// Delete 刪除已實現損益記錄（當交易被刪除時）
	Delete(id string) error
}

// realizedProfitRepository 已實現損益資料存取實作
type realizedProfitRepository struct {
	db *sql.DB
}

// NewRealizedProfitRepository 建立新的已實現損益 repository
func NewRealizedProfitRepository(db *sql.DB) RealizedProfitRepository {
	return &realizedProfitRepository{db: db}
}

// Create 建立已實現損益記錄
func (r *realizedProfitRepository) Create(input *models.CreateRealizedProfitInput) (*models.RealizedProfit, error) {
	// 計算已實現損益
	// 已實現損益 = (賣出金額 - 賣出手續費) - 成本基礎
	realizedPL := (input.SellAmount - input.SellFee) - input.CostBasis

	// 計算已實現損益百分比
	var realizedPLPct float64
	if input.CostBasis > 0 {
		realizedPLPct = (realizedPL / input.CostBasis) * 100
	}

	query := `
		INSERT INTO realized_profits (
			transaction_id, symbol, asset_type, sell_date, quantity,
			sell_price, sell_amount, sell_fee, cost_basis,
			realized_pl, realized_pl_pct, currency
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, transaction_id, symbol, asset_type, sell_date, quantity,
		          sell_price, sell_amount, sell_fee, cost_basis,
		          realized_pl, realized_pl_pct, currency, created_at, updated_at
	`

	var result models.RealizedProfit
	err := r.db.QueryRow(
		query,
		input.TransactionID,
		input.Symbol,
		input.AssetType,
		input.SellDate,
		input.Quantity,
		input.SellPrice,
		input.SellAmount,
		input.SellFee,
		input.CostBasis,
		realizedPL,
		realizedPLPct,
		input.Currency,
	).Scan(
		&result.ID,
		&result.TransactionID,
		&result.Symbol,
		&result.AssetType,
		&result.SellDate,
		&result.Quantity,
		&result.SellPrice,
		&result.SellAmount,
		&result.SellFee,
		&result.CostBasis,
		&result.RealizedPL,
		&result.RealizedPLPct,
		&result.Currency,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create realized profit: %w", err)
	}

	return &result, nil
}

// GetByTransactionID 根據交易 ID 取得已實現損益
func (r *realizedProfitRepository) GetByTransactionID(transactionID string) (*models.RealizedProfit, error) {
	query := `
		SELECT id, transaction_id, symbol, asset_type, sell_date, quantity,
		       sell_price, sell_amount, sell_fee, cost_basis,
		       realized_pl, realized_pl_pct, currency, created_at, updated_at
		FROM realized_profits
		WHERE transaction_id = $1
	`

	var result models.RealizedProfit
	err := r.db.QueryRow(query, transactionID).Scan(
		&result.ID,
		&result.TransactionID,
		&result.Symbol,
		&result.AssetType,
		&result.SellDate,
		&result.Quantity,
		&result.SellPrice,
		&result.SellAmount,
		&result.SellFee,
		&result.CostBasis,
		&result.RealizedPL,
		&result.RealizedPLPct,
		&result.Currency,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("realized profit not found for transaction: %s", transactionID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get realized profit: %w", err)
	}

	return &result, nil
}

// GetAll 取得所有已實現損益記錄（支援篩選）
func (r *realizedProfitRepository) GetAll(filters models.RealizedProfitFilters) ([]*models.RealizedProfit, error) {
	query := `
		SELECT id, transaction_id, symbol, asset_type, sell_date, quantity,
		       sell_price, sell_amount, sell_fee, cost_basis,
		       realized_pl, realized_pl_pct, currency, created_at, updated_at
		FROM realized_profits
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// 動態建立 WHERE 條件
	if filters.AssetType != nil {
		query += fmt.Sprintf(" AND asset_type = $%d", argCount)
		args = append(args, *filters.AssetType)
		argCount++
	}

	if filters.Symbol != nil {
		query += fmt.Sprintf(" AND symbol = $%d", argCount)
		args = append(args, *filters.Symbol)
		argCount++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND sell_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
		argCount++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND sell_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
		argCount++
	}

	// 排序
	query += " ORDER BY sell_date DESC, created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query realized profits: %w", err)
	}
	defer rows.Close()

	results := []*models.RealizedProfit{}
	for rows.Next() {
		var rp models.RealizedProfit
		err := rows.Scan(
			&rp.ID,
			&rp.TransactionID,
			&rp.Symbol,
			&rp.AssetType,
			&rp.SellDate,
			&rp.Quantity,
			&rp.SellPrice,
			&rp.SellAmount,
			&rp.SellFee,
			&rp.CostBasis,
			&rp.RealizedPL,
			&rp.RealizedPLPct,
			&rp.Currency,
			&rp.CreatedAt,
			&rp.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan realized profit: %w", err)
		}
		results = append(results, &rp)
	}

	return results, nil
}

// Delete 刪除已實現損益記錄
func (r *realizedProfitRepository) Delete(id string) error {
	query := `DELETE FROM realized_profits WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete realized profit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("realized profit not found: %s", id)
	}

	return nil
}

