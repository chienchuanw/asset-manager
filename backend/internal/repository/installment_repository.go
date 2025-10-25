package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
)

// InstallmentRepository 分期資料存取介面
type InstallmentRepository interface {
	Create(input *models.CreateInstallmentInput) (*models.Installment, error)
	GetByID(id uuid.UUID) (*models.Installment, error)
	List(filters InstallmentFilters) ([]*models.Installment, error)
	Update(id uuid.UUID, input *models.UpdateInstallmentInput) (*models.Installment, error)
	Delete(id uuid.UUID) error
	GetDueBillings(date time.Time) ([]*models.Installment, error)
	GetCompletingSoon(remainingCount int) ([]*models.Installment, error)
}

// InstallmentFilters 分期查詢篩選條件
type InstallmentFilters struct {
	Status     *models.InstallmentStatus `json:"status,omitempty"`
	CategoryID *uuid.UUID                `json:"category_id,omitempty"`
	Limit      int                       `json:"limit,omitempty"`
	Offset     int                       `json:"offset,omitempty"`
}

// installmentRepository 分期資料存取實作
type installmentRepository struct {
	db *sql.DB
}

// NewInstallmentRepository 建立新的分期 repository
func NewInstallmentRepository(db *sql.DB) InstallmentRepository {
	return &installmentRepository{db: db}
}

// Create 建立新的分期
func (r *installmentRepository) Create(input *models.CreateInstallmentInput) (*models.Installment, error) {
	// 建立 Installment 模型並計算利息
	installment := &models.Installment{
		TotalAmount:      input.TotalAmount,
		InstallmentCount: input.InstallmentCount,
		InterestRate:     input.InterestRate,
	}
	installment.CalculateInterest()

	query := `
		INSERT INTO installments (
			name, total_amount, currency, installment_count, installment_amount,
			interest_rate, total_interest, paid_count, billing_day,
			category_id, start_date, note
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, name, total_amount, currency, installment_count, installment_amount,
			interest_rate, total_interest, paid_count, billing_day,
			category_id, start_date, status, note, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		input.Name,
		input.TotalAmount,
		models.CurrencyTWD, // 固定為 TWD
		input.InstallmentCount,
		installment.InstallmentAmount,
		input.InterestRate,
		installment.TotalInterest,
		0, // paid_count 初始為 0
		input.BillingDay,
		input.CategoryID,
		input.StartDate,
		input.Note,
	).Scan(
		&installment.ID,
		&installment.Name,
		&installment.TotalAmount,
		&installment.Currency,
		&installment.InstallmentCount,
		&installment.InstallmentAmount,
		&installment.InterestRate,
		&installment.TotalInterest,
		&installment.PaidCount,
		&installment.BillingDay,
		&installment.CategoryID,
		&installment.StartDate,
		&installment.Status,
		&installment.Note,
		&installment.CreatedAt,
		&installment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create installment: %w", err)
	}

	return installment, nil
}

// GetByID 根據 ID 取得分期（包含分類資訊）
func (r *installmentRepository) GetByID(id uuid.UUID) (*models.Installment, error) {
	query := `
		SELECT 
			i.id, i.name, i.total_amount, i.currency, i.installment_count, i.installment_amount,
			i.interest_rate, i.total_interest, i.paid_count, i.billing_day,
			i.category_id, i.start_date, i.status, i.note, i.created_at, i.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM installments i
		LEFT JOIN cash_flow_categories c ON i.category_id = c.id
		WHERE i.id = $1
	`

	installment := &models.Installment{
		Category: &models.CashFlowCategory{},
	}

	err := r.db.QueryRow(query, id).Scan(
		&installment.ID,
		&installment.Name,
		&installment.TotalAmount,
		&installment.Currency,
		&installment.InstallmentCount,
		&installment.InstallmentAmount,
		&installment.InterestRate,
		&installment.TotalInterest,
		&installment.PaidCount,
		&installment.BillingDay,
		&installment.CategoryID,
		&installment.StartDate,
		&installment.Status,
		&installment.Note,
		&installment.CreatedAt,
		&installment.UpdatedAt,
		&installment.Category.ID,
		&installment.Category.Name,
		&installment.Category.Type,
		&installment.Category.IsSystem,
		&installment.Category.CreatedAt,
		&installment.Category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return installment, nil
}

// List 取得分期列表（包含分類資訊）
func (r *installmentRepository) List(filters InstallmentFilters) ([]*models.Installment, error) {
	query := `
		SELECT 
			i.id, i.name, i.total_amount, i.currency, i.installment_count, i.installment_amount,
			i.interest_rate, i.total_interest, i.paid_count, i.billing_day,
			i.category_id, i.start_date, i.status, i.note, i.created_at, i.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM installments i
		LEFT JOIN cash_flow_categories c ON i.category_id = c.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// 套用篩選條件
	if filters.Status != nil {
		query += fmt.Sprintf(" AND i.status = $%d", argCount)
		args = append(args, *filters.Status)
		argCount++
	}

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND i.category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
		argCount++
	}

	// 排序
	query += " ORDER BY i.created_at DESC"

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
		return nil, fmt.Errorf("failed to list installments: %w", err)
	}
	defer rows.Close()

	installments := []*models.Installment{}
	for rows.Next() {
		installment := &models.Installment{
			Category: &models.CashFlowCategory{},
		}

		err := rows.Scan(
			&installment.ID,
			&installment.Name,
			&installment.TotalAmount,
			&installment.Currency,
			&installment.InstallmentCount,
			&installment.InstallmentAmount,
			&installment.InterestRate,
			&installment.TotalInterest,
			&installment.PaidCount,
			&installment.BillingDay,
			&installment.CategoryID,
			&installment.StartDate,
			&installment.Status,
			&installment.Note,
			&installment.CreatedAt,
			&installment.UpdatedAt,
			&installment.Category.ID,
			&installment.Category.Name,
			&installment.Category.Type,
			&installment.Category.IsSystem,
			&installment.Category.CreatedAt,
			&installment.Category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan installment: %w", err)
		}

		installments = append(installments, installment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating installments: %w", err)
	}

	return installments, nil
}

// Update 更新分期
func (r *installmentRepository) Update(id uuid.UUID, input *models.UpdateInstallmentInput) (*models.Installment, error) {
	// 建立動態更新語句
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if input.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *input.Name)
		argCount++
	}

	if input.CategoryID != nil {
		updates = append(updates, fmt.Sprintf("category_id = $%d", argCount))
		args = append(args, *input.CategoryID)
		argCount++
	}

	if input.Note != nil {
		updates = append(updates, fmt.Sprintf("note = $%d", argCount))
		args = append(args, *input.Note)
		argCount++
	}

	if len(updates) == 0 {
		return r.GetByID(id)
	}

	// 加入 ID 參數
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE installments
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE id = $%d
		RETURNING id, name, total_amount, currency, installment_count, installment_amount,
			interest_rate, total_interest, paid_count, billing_day,
			category_id, start_date, status, note, created_at, updated_at
	`, strings.Join(updates, ", "), argCount)

	installment := &models.Installment{}
	err := r.db.QueryRow(query, args...).Scan(
		&installment.ID,
		&installment.Name,
		&installment.TotalAmount,
		&installment.Currency,
		&installment.InstallmentCount,
		&installment.InstallmentAmount,
		&installment.InterestRate,
		&installment.TotalInterest,
		&installment.PaidCount,
		&installment.BillingDay,
		&installment.CategoryID,
		&installment.StartDate,
		&installment.Status,
		&installment.Note,
		&installment.CreatedAt,
		&installment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update installment: %w", err)
	}

	return installment, nil
}

// Delete 刪除分期
func (r *installmentRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM installments WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete installment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetDueBillings 取得指定日期需要扣款的分期
func (r *installmentRepository) GetDueBillings(date time.Time) ([]*models.Installment, error) {
	day := date.Day()

	query := `
		SELECT
			i.id, i.name, i.total_amount, i.currency, i.installment_count, i.installment_amount,
			i.interest_rate, i.total_interest, i.paid_count, i.billing_day,
			i.category_id, i.start_date, i.status, i.note, i.created_at, i.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM installments i
		LEFT JOIN cash_flow_categories c ON i.category_id = c.id
		WHERE i.status = $1
			AND i.billing_day = $2
			AND i.start_date <= $3
			AND i.paid_count < i.installment_count
		ORDER BY i.name
	`

	rows, err := r.db.Query(query, models.InstallmentStatusActive, day, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get due billings: %w", err)
	}
	defer rows.Close()

	installments := []*models.Installment{}
	for rows.Next() {
		installment := &models.Installment{
			Category: &models.CashFlowCategory{},
		}

		err := rows.Scan(
			&installment.ID,
			&installment.Name,
			&installment.TotalAmount,
			&installment.Currency,
			&installment.InstallmentCount,
			&installment.InstallmentAmount,
			&installment.InterestRate,
			&installment.TotalInterest,
			&installment.PaidCount,
			&installment.BillingDay,
			&installment.CategoryID,
			&installment.StartDate,
			&installment.Status,
			&installment.Note,
			&installment.CreatedAt,
			&installment.UpdatedAt,
			&installment.Category.ID,
			&installment.Category.Name,
			&installment.Category.Type,
			&installment.Category.IsSystem,
			&installment.Category.CreatedAt,
			&installment.Category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan installment: %w", err)
		}

		installments = append(installments, installment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating installments: %w", err)
	}

	return installments, nil
}

// GetCompletingSoon 取得即將完成的分期
func (r *installmentRepository) GetCompletingSoon(remainingCount int) ([]*models.Installment, error) {
	query := `
		SELECT
			i.id, i.name, i.total_amount, i.currency, i.installment_count, i.installment_amount,
			i.interest_rate, i.total_interest, i.paid_count, i.billing_day,
			i.category_id, i.start_date, i.status, i.note, i.created_at, i.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM installments i
		LEFT JOIN cash_flow_categories c ON i.category_id = c.id
		WHERE i.status = $1
			AND (i.installment_count - i.paid_count) <= $2
			AND (i.installment_count - i.paid_count) > 0
		ORDER BY (i.installment_count - i.paid_count)
	`

	rows, err := r.db.Query(query, models.InstallmentStatusActive, remainingCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get completing installments: %w", err)
	}
	defer rows.Close()

	installments := []*models.Installment{}
	for rows.Next() {
		installment := &models.Installment{
			Category: &models.CashFlowCategory{},
		}

		err := rows.Scan(
			&installment.ID,
			&installment.Name,
			&installment.TotalAmount,
			&installment.Currency,
			&installment.InstallmentCount,
			&installment.InstallmentAmount,
			&installment.InterestRate,
			&installment.TotalInterest,
			&installment.PaidCount,
			&installment.BillingDay,
			&installment.CategoryID,
			&installment.StartDate,
			&installment.Status,
			&installment.Note,
			&installment.CreatedAt,
			&installment.UpdatedAt,
			&installment.Category.ID,
			&installment.Category.Name,
			&installment.Category.Type,
			&installment.Category.IsSystem,
			&installment.Category.CreatedAt,
			&installment.Category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan installment: %w", err)
		}

		installments = append(installments, installment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating installments: %w", err)
	}

	return installments, nil
}

