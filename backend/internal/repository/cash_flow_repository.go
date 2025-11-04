package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
)

// CashFlowRepository 現金流記錄資料存取介面
type CashFlowRepository interface {
	Create(input *models.CreateCashFlowInput) (*models.CashFlow, error)
	GetByID(id uuid.UUID) (*models.CashFlow, error)
	GetAll(filters CashFlowFilters) ([]*models.CashFlow, error)
	Update(id uuid.UUID, input *models.UpdateCashFlowInput) (*models.CashFlow, error)
	Delete(id uuid.UUID) error
	GetSummary(startDate, endDate time.Time) (*CashFlowSummary, error)
}

// CashFlowFilters 現金流查詢篩選條件
type CashFlowFilters struct {
	Type       *models.CashFlowType `json:"type,omitempty"`
	CategoryID *uuid.UUID           `json:"category_id,omitempty"`
	StartDate  *time.Time           `json:"start_date,omitempty"`
	EndDate    *time.Time           `json:"end_date,omitempty"`
	Limit      int                  `json:"limit,omitempty"`
	Offset     int                  `json:"offset,omitempty"`
}

// CashFlowSummary 現金流摘要
type CashFlowSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetCashFlow  float64 `json:"net_cash_flow"`
}

// cashFlowRepository 現金流記錄資料存取實作
type cashFlowRepository struct {
	db *sql.DB
}

// NewCashFlowRepository 建立新的現金流記錄 repository
func NewCashFlowRepository(db *sql.DB) CashFlowRepository {
	return &cashFlowRepository{db: db}
}

// Create 建立新的現金流記錄
func (r *cashFlowRepository) Create(input *models.CreateCashFlowInput) (*models.CashFlow, error) {
	query := `
		INSERT INTO cash_flows (date, type, category_id, amount, currency, description, note, source_type, source_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, date, type, category_id, amount, currency, description, note, source_type, source_id, created_at, updated_at
	`

	cashFlow := &models.CashFlow{}
	err := r.db.QueryRow(
		query,
		input.Date,
		input.Type,
		input.CategoryID,
		input.Amount,
		models.CurrencyTWD, // 固定為 TWD
		input.Description,
		input.Note,
		input.SourceType,
		input.SourceID,
	).Scan(
		&cashFlow.ID,
		&cashFlow.Date,
		&cashFlow.Type,
		&cashFlow.CategoryID,
		&cashFlow.Amount,
		&cashFlow.Currency,
		&cashFlow.Description,
		&cashFlow.Note,
		&cashFlow.SourceType,
		&cashFlow.SourceID,
		&cashFlow.CreatedAt,
		&cashFlow.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create cash flow: %w", err)
	}

	return cashFlow, nil
}

// GetByID 根據 ID 取得現金流記錄（包含分類資訊）
func (r *cashFlowRepository) GetByID(id uuid.UUID) (*models.CashFlow, error) {
	query := `
		SELECT
			cf.id, cf.date, cf.type, cf.category_id, cf.amount, cf.currency,
			cf.description, cf.note, cf.source_type, cf.source_id, cf.created_at, cf.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM cash_flows cf
		LEFT JOIN cash_flow_categories c ON cf.category_id = c.id
		WHERE cf.id = $1
	`

	cashFlow := &models.CashFlow{
		Category: &models.CashFlowCategory{},
	}

	err := r.db.QueryRow(query, id).Scan(
		&cashFlow.ID,
		&cashFlow.Date,
		&cashFlow.Type,
		&cashFlow.CategoryID,
		&cashFlow.Amount,
		&cashFlow.Currency,
		&cashFlow.Description,
		&cashFlow.Note,
		&cashFlow.SourceType,
		&cashFlow.SourceID,
		&cashFlow.CreatedAt,
		&cashFlow.UpdatedAt,
		&cashFlow.Category.ID,
		&cashFlow.Category.Name,
		&cashFlow.Category.Type,
		&cashFlow.Category.IsSystem,
		&cashFlow.Category.CreatedAt,
		&cashFlow.Category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cash flow not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cash flow: %w", err)
	}

	return cashFlow, nil
}

// GetAll 取得所有現金流記錄（支援篩選，包含分類資訊）
func (r *cashFlowRepository) GetAll(filters CashFlowFilters) ([]*models.CashFlow, error) {
	query := `
		SELECT
			cf.id, cf.date, cf.type, cf.category_id, cf.amount, cf.currency,
			cf.description, cf.note, cf.source_type, cf.source_id, cf.created_at, cf.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM cash_flows cf
		LEFT JOIN cash_flow_categories c ON cf.category_id = c.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// 動態建立 WHERE 條件
	if filters.Type != nil {
		query += fmt.Sprintf(" AND cf.type = $%d", argCount)
		args = append(args, *filters.Type)
		argCount++
	}

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND cf.category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
		argCount++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND cf.date >= $%d", argCount)
		args = append(args, *filters.StartDate)
		argCount++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND cf.date <= $%d", argCount)
		args = append(args, *filters.EndDate)
		argCount++
	}

	// 排序
	query += " ORDER BY cf.date DESC, cf.created_at DESC"

	// 分頁
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.Limit)
		argCount++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filters.Offset)
		argCount++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query cash flows: %w", err)
	}
	defer rows.Close()

	cashFlows := []*models.CashFlow{}
	for rows.Next() {
		cashFlow := &models.CashFlow{
			Category: &models.CashFlowCategory{},
		}
		err := rows.Scan(
			&cashFlow.ID,
			&cashFlow.Date,
			&cashFlow.Type,
			&cashFlow.CategoryID,
			&cashFlow.Amount,
			&cashFlow.Currency,
			&cashFlow.Description,
			&cashFlow.Note,
			&cashFlow.SourceType,
			&cashFlow.SourceID,
			&cashFlow.CreatedAt,
			&cashFlow.UpdatedAt,
			&cashFlow.Category.ID,
			&cashFlow.Category.Name,
			&cashFlow.Category.Type,
			&cashFlow.Category.IsSystem,
			&cashFlow.Category.CreatedAt,
			&cashFlow.Category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cash flow: %w", err)
		}
		cashFlows = append(cashFlows, cashFlow)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cash flows: %w", err)
	}

	return cashFlows, nil
}

// Update 更新現金流記錄
func (r *cashFlowRepository) Update(id uuid.UUID, input *models.UpdateCashFlowInput) (*models.CashFlow, error) {
	// 動態建立 UPDATE 語句
	setClauses := []string{}
	args := []interface{}{}
	argCount := 1

	if input.Date != nil {
		setClauses = append(setClauses, fmt.Sprintf("date = $%d", argCount))
		args = append(args, *input.Date)
		argCount++
	}

	if input.CategoryID != nil {
		setClauses = append(setClauses, fmt.Sprintf("category_id = $%d", argCount))
		args = append(args, *input.CategoryID)
		argCount++
	}

	if input.Amount != nil {
		setClauses = append(setClauses, fmt.Sprintf("amount = $%d", argCount))
		args = append(args, *input.Amount)
		argCount++
	}

	if input.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argCount))
		args = append(args, *input.Description)
		argCount++
	}

	if input.Note != nil {
		setClauses = append(setClauses, fmt.Sprintf("note = $%d", argCount))
		args = append(args, *input.Note)
		argCount++
	}

	if input.SourceType != nil {
		setClauses = append(setClauses, fmt.Sprintf("source_type = $%d", argCount))
		args = append(args, *input.SourceType)
		argCount++
	}

	if input.SourceID != nil {
		setClauses = append(setClauses, fmt.Sprintf("source_id = $%d", argCount))
		args = append(args, *input.SourceID)
		argCount++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// 加入 ID 參數
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE cash_flows
		SET %s
		WHERE id = $%d
		RETURNING id, date, type, category_id, amount, currency, description, note, source_type, source_id, created_at, updated_at
	`, strings.Join(setClauses, ", "), argCount)

	cashFlow := &models.CashFlow{}
	err := r.db.QueryRow(query, args...).Scan(
		&cashFlow.ID,
		&cashFlow.Date,
		&cashFlow.Type,
		&cashFlow.CategoryID,
		&cashFlow.Amount,
		&cashFlow.Currency,
		&cashFlow.Description,
		&cashFlow.Note,
		&cashFlow.SourceType,
		&cashFlow.SourceID,
		&cashFlow.CreatedAt,
		&cashFlow.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cash flow not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update cash flow: %w", err)
	}

	return cashFlow, nil
}

// Delete 刪除現金流記錄
func (r *cashFlowRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM cash_flows WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cash flow: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cash flow not found")
	}

	return nil
}

// GetSummary 取得指定日期區間的現金流摘要
func (r *cashFlowRepository) GetSummary(startDate, endDate time.Time) (*CashFlowSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense
		FROM cash_flows
		WHERE date >= $1 AND date <= $2
	`

	summary := &CashFlowSummary{}
	err := r.db.QueryRow(query, startDate, endDate).Scan(
		&summary.TotalIncome,
		&summary.TotalExpense,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get cash flow summary: %w", err)
	}

	// 計算淨現金流
	summary.NetCashFlow = summary.TotalIncome - summary.TotalExpense

	return summary, nil
}

