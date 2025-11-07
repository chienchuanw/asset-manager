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
	GetMonthlySummary(year, month int) (*models.MonthlyCashFlowSummary, error)
	GetYearlySummary(year int) (*models.YearlyCashFlowSummary, error)
	GetCategorySummary(startDate, endDate time.Time, cashFlowType models.CashFlowType) ([]*models.CategorySummary, error)
	GetTopExpenses(startDate, endDate time.Time, limit int) ([]*models.CashFlow, error)
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
		INSERT INTO cash_flows (date, type, category_id, amount, currency, description, note, source_type, source_id, target_type, target_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, date, type, category_id, amount, currency, description, note, source_type, source_id, target_type, target_id, created_at, updated_at
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
		input.TargetType,
		input.TargetID,
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
		&cashFlow.TargetType,
		&cashFlow.TargetID,
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
			cf.description, cf.note, cf.source_type, cf.source_id, cf.target_type, cf.target_id, cf.created_at, cf.updated_at,
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
		&cashFlow.TargetType,
		&cashFlow.TargetID,
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
			cf.description, cf.note, cf.source_type, cf.source_id, cf.target_type, cf.target_id, cf.created_at, cf.updated_at,
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
			&cashFlow.TargetType,
			&cashFlow.TargetID,
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

// GetMonthlySummary 取得指定月份的現金流摘要
func (r *cashFlowRepository) GetMonthlySummary(year, month int) (*models.MonthlyCashFlowSummary, error) {
	// 計算月份的開始和結束日期
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	// 取得基本統計
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense,
			COALESCE(SUM(CASE WHEN type = 'income' THEN 1 ELSE 0 END), 0) as income_count,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN 1 ELSE 0 END), 0) as expense_count
		FROM cash_flows
		WHERE date >= $1 AND date <= $2
	`

	summary := &models.MonthlyCashFlowSummary{
		Year:  year,
		Month: month,
	}

	err := r.db.QueryRow(query, startDate, endDate).Scan(
		&summary.TotalIncome,
		&summary.TotalExpense,
		&summary.IncomeCount,
		&summary.ExpenseCount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get monthly summary: %w", err)
	}

	// 計算淨現金流
	summary.NetCashFlow = summary.TotalIncome - summary.TotalExpense

	// 取得收入分類摘要
	incomeSummary, err := r.GetCategorySummary(startDate, endDate, models.CashFlowTypeIncome)
	if err != nil {
		return nil, fmt.Errorf("failed to get income category summary: %w", err)
	}
	summary.IncomeCategoryBreakdown = incomeSummary

	// 取得支出分類摘要
	expenseSummary, err := r.GetCategorySummary(startDate, endDate, models.CashFlowTypeExpense)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense category summary: %w", err)
	}
	summary.ExpenseCategoryBreakdown = expenseSummary

	// 取得前 10 大支出
	topExpenses, err := r.GetTopExpenses(startDate, endDate, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top expenses: %w", err)
	}
	summary.TopExpenses = topExpenses

	return summary, nil
}

// GetYearlySummary 取得指定年度的現金流摘要
func (r *cashFlowRepository) GetYearlySummary(year int) (*models.YearlyCashFlowSummary, error) {
	// 計算年度的開始和結束日期
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.Local)

	// 取得基本統計
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense,
			COALESCE(SUM(CASE WHEN type = 'income' THEN 1 ELSE 0 END), 0) as income_count,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN 1 ELSE 0 END), 0) as expense_count
		FROM cash_flows
		WHERE date >= $1 AND date <= $2
	`

	summary := &models.YearlyCashFlowSummary{
		Year: year,
	}

	err := r.db.QueryRow(query, startDate, endDate).Scan(
		&summary.TotalIncome,
		&summary.TotalExpense,
		&summary.IncomeCount,
		&summary.ExpenseCount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly summary: %w", err)
	}

	// 計算淨現金流
	summary.NetCashFlow = summary.TotalIncome - summary.TotalExpense

	// 取得收入分類摘要
	incomeSummary, err := r.GetCategorySummary(startDate, endDate, models.CashFlowTypeIncome)
	if err != nil {
		return nil, fmt.Errorf("failed to get income category summary: %w", err)
	}
	summary.IncomeCategoryBreakdown = incomeSummary

	// 取得支出分類摘要
	expenseSummary, err := r.GetCategorySummary(startDate, endDate, models.CashFlowTypeExpense)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense category summary: %w", err)
	}
	summary.ExpenseCategoryBreakdown = expenseSummary

	// 取得前 10 大支出
	topExpenses, err := r.GetTopExpenses(startDate, endDate, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top expenses: %w", err)
	}
	summary.TopExpenses = topExpenses

	// 取得月度細分
	monthlyBreakdown, err := r.getMonthlyBreakdown(year)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly breakdown: %w", err)
	}
	summary.MonthlyBreakdown = monthlyBreakdown

	return summary, nil
}

// GetCategorySummary 取得指定日期區間和類型的分類摘要
func (r *cashFlowRepository) GetCategorySummary(startDate, endDate time.Time, cashFlowType models.CashFlowType) ([]*models.CategorySummary, error) {
	query := `
		SELECT
			c.id as category_id,
			c.name as category_name,
			COALESCE(SUM(cf.amount), 0) as amount,
			COUNT(cf.id) as count
		FROM cash_flow_categories c
		LEFT JOIN cash_flows cf ON c.id = cf.category_id
			AND cf.date >= $1
			AND cf.date <= $2
			AND cf.type = $3
		WHERE c.type = $3
		GROUP BY c.id, c.name
		HAVING COUNT(cf.id) > 0
		ORDER BY amount DESC
	`

	rows, err := r.db.Query(query, startDate, endDate, cashFlowType)
	if err != nil {
		return nil, fmt.Errorf("failed to query category summary: %w", err)
	}
	defer rows.Close()

	var summaries []*models.CategorySummary
	for rows.Next() {
		summary := &models.CategorySummary{}
		err := rows.Scan(
			&summary.CategoryID,
			&summary.CategoryName,
			&summary.Amount,
			&summary.Count,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category summary: %w", err)
		}
		summaries = append(summaries, summary)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating category summary rows: %w", err)
	}

	return summaries, nil
}

// GetTopExpenses 取得指定日期區間的前 N 大支出
func (r *cashFlowRepository) GetTopExpenses(startDate, endDate time.Time, limit int) ([]*models.CashFlow, error) {
	query := `
		SELECT
			cf.id, cf.date, cf.type, cf.category_id, cf.amount, cf.currency,
			cf.description, cf.note, cf.source_type, cf.source_id, cf.target_type, cf.target_id,
			cf.created_at, cf.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM cash_flows cf
		LEFT JOIN cash_flow_categories c ON cf.category_id = c.id
		WHERE cf.date >= $1 AND cf.date <= $2 AND cf.type = 'expense'
		ORDER BY cf.amount DESC
		LIMIT $3
	`

	rows, err := r.db.Query(query, startDate, endDate, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top expenses: %w", err)
	}
	defer rows.Close()

	var expenses []*models.CashFlow
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
			&cashFlow.TargetType,
			&cashFlow.TargetID,
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

		expenses = append(expenses, cashFlow)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cash flow rows: %w", err)
	}

	return expenses, nil
}

// getMonthlyBreakdown 取得指定年度的月度細分
func (r *cashFlowRepository) getMonthlyBreakdown(year int) ([]*models.MonthlyBreakdown, error) {
	query := `
		SELECT
			EXTRACT(MONTH FROM date)::int as month,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense
		FROM cash_flows
		WHERE EXTRACT(YEAR FROM date) = $1
		GROUP BY EXTRACT(MONTH FROM date)
		ORDER BY month
	`

	rows, err := r.db.Query(query, year)
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly breakdown: %w", err)
	}
	defer rows.Close()

	var breakdowns []*models.MonthlyBreakdown
	for rows.Next() {
		breakdown := &models.MonthlyBreakdown{}
		err := rows.Scan(
			&breakdown.Month,
			&breakdown.Income,
			&breakdown.Expense,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan monthly breakdown: %w", err)
		}

		// 計算淨現金流
		breakdown.NetCashFlow = breakdown.Income - breakdown.Expense

		breakdowns = append(breakdowns, breakdown)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating monthly breakdown rows: %w", err)
	}

	return breakdowns, nil
}


