package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
)

// SubscriptionRepository 訂閱資料存取介面
type SubscriptionRepository interface {
	Create(input *models.CreateSubscriptionInput) (*models.Subscription, error)
	GetByID(id uuid.UUID) (*models.Subscription, error)
	List(filters SubscriptionFilters) ([]*models.Subscription, error)
	Update(id uuid.UUID, input *models.UpdateSubscriptionInput) (*models.Subscription, error)
	Delete(id uuid.UUID) error
	GetDueBillings(date time.Time) ([]*models.Subscription, error)
	GetExpiringSoon(days int) ([]*models.Subscription, error)
}

// SubscriptionFilters 訂閱查詢篩選條件
type SubscriptionFilters struct {
	Status     *models.SubscriptionStatus `json:"status,omitempty"`
	CategoryID *uuid.UUID                 `json:"category_id,omitempty"`
	Limit      int                        `json:"limit,omitempty"`
	Offset     int                        `json:"offset,omitempty"`
}

// subscriptionRepository 訂閱資料存取實作
type subscriptionRepository struct {
	db *sql.DB
}

// NewSubscriptionRepository 建立新的訂閱 repository
func NewSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

// Create 建立新的訂閱
func (r *subscriptionRepository) Create(input *models.CreateSubscriptionInput) (*models.Subscription, error) {
	query := `
		INSERT INTO subscriptions (
			name, amount, currency, billing_cycle, billing_day, 
			category_id, start_date, end_date, auto_renew, note
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, name, amount, currency, billing_cycle, billing_day, 
			category_id, start_date, end_date, auto_renew, status, note, 
			created_at, updated_at
	`

	subscription := &models.Subscription{}
	err := r.db.QueryRow(
		query,
		input.Name,
		input.Amount,
		models.CurrencyTWD, // 固定為 TWD
		input.BillingCycle,
		input.BillingDay,
		input.CategoryID,
		input.StartDate,
		input.EndDate,
		input.AutoRenew,
		input.Note,
	).Scan(
		&subscription.ID,
		&subscription.Name,
		&subscription.Amount,
		&subscription.Currency,
		&subscription.BillingCycle,
		&subscription.BillingDay,
		&subscription.CategoryID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.AutoRenew,
		&subscription.Status,
		&subscription.Note,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return subscription, nil
}

// GetByID 根據 ID 取得訂閱（包含分類資訊）
func (r *subscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT 
			s.id, s.name, s.amount, s.currency, s.billing_cycle, s.billing_day,
			s.category_id, s.start_date, s.end_date, s.auto_renew, s.status, s.note,
			s.created_at, s.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM subscriptions s
		LEFT JOIN cash_flow_categories c ON s.category_id = c.id
		WHERE s.id = $1
	`

	subscription := &models.Subscription{
		Category: &models.CashFlowCategory{},
	}

	err := r.db.QueryRow(query, id).Scan(
		&subscription.ID,
		&subscription.Name,
		&subscription.Amount,
		&subscription.Currency,
		&subscription.BillingCycle,
		&subscription.BillingDay,
		&subscription.CategoryID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.AutoRenew,
		&subscription.Status,
		&subscription.Note,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&subscription.Category.ID,
		&subscription.Category.Name,
		&subscription.Category.Type,
		&subscription.Category.IsSystem,
		&subscription.Category.CreatedAt,
		&subscription.Category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// List 取得訂閱列表（包含分類資訊）
func (r *subscriptionRepository) List(filters SubscriptionFilters) ([]*models.Subscription, error) {
	query := `
		SELECT 
			s.id, s.name, s.amount, s.currency, s.billing_cycle, s.billing_day,
			s.category_id, s.start_date, s.end_date, s.auto_renew, s.status, s.note,
			s.created_at, s.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM subscriptions s
		LEFT JOIN cash_flow_categories c ON s.category_id = c.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// 套用篩選條件
	if filters.Status != nil {
		query += fmt.Sprintf(" AND s.status = $%d", argCount)
		args = append(args, *filters.Status)
		argCount++
	}

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND s.category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
		argCount++
	}

	// 排序
	query += " ORDER BY s.created_at DESC"

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
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions := []*models.Subscription{}
	for rows.Next() {
		subscription := &models.Subscription{
			Category: &models.CashFlowCategory{},
		}

		err := rows.Scan(
			&subscription.ID,
			&subscription.Name,
			&subscription.Amount,
			&subscription.Currency,
			&subscription.BillingCycle,
			&subscription.BillingDay,
			&subscription.CategoryID,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.AutoRenew,
			&subscription.Status,
			&subscription.Note,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
			&subscription.Category.ID,
			&subscription.Category.Name,
			&subscription.Category.Type,
			&subscription.Category.IsSystem,
			&subscription.Category.CreatedAt,
			&subscription.Category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subscriptions, nil
}

// Update 更新訂閱
func (r *subscriptionRepository) Update(id uuid.UUID, input *models.UpdateSubscriptionInput) (*models.Subscription, error) {
	// 建立動態更新語句
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if input.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *input.Name)
		argCount++
	}

	if input.Amount != nil {
		updates = append(updates, fmt.Sprintf("amount = $%d", argCount))
		args = append(args, *input.Amount)
		argCount++
	}

	if input.BillingCycle != nil {
		updates = append(updates, fmt.Sprintf("billing_cycle = $%d", argCount))
		args = append(args, *input.BillingCycle)
		argCount++
	}

	if input.BillingDay != nil {
		updates = append(updates, fmt.Sprintf("billing_day = $%d", argCount))
		args = append(args, *input.BillingDay)
		argCount++
	}

	if input.CategoryID != nil {
		updates = append(updates, fmt.Sprintf("category_id = $%d", argCount))
		args = append(args, *input.CategoryID)
		argCount++
	}

	if input.EndDate != nil {
		updates = append(updates, fmt.Sprintf("end_date = $%d", argCount))
		args = append(args, *input.EndDate)
		argCount++
	}

	if input.AutoRenew != nil {
		updates = append(updates, fmt.Sprintf("auto_renew = $%d", argCount))
		args = append(args, *input.AutoRenew)
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
		UPDATE subscriptions
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE id = $%d
		RETURNING id, name, amount, currency, billing_cycle, billing_day,
			category_id, start_date, end_date, auto_renew, status, note,
			created_at, updated_at
	`, strings.Join(updates, ", "), argCount)

	subscription := &models.Subscription{}
	err := r.db.QueryRow(query, args...).Scan(
		&subscription.ID,
		&subscription.Name,
		&subscription.Amount,
		&subscription.Currency,
		&subscription.BillingCycle,
		&subscription.BillingDay,
		&subscription.CategoryID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.AutoRenew,
		&subscription.Status,
		&subscription.Note,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return subscription, nil
}

// Delete 刪除訂閱
func (r *subscriptionRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
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

// GetDueBillings 取得指定日期需要扣款的訂閱
func (r *subscriptionRepository) GetDueBillings(date time.Time) ([]*models.Subscription, error) {
	day := date.Day()

	query := `
		SELECT
			s.id, s.name, s.amount, s.currency, s.billing_cycle, s.billing_day,
			s.category_id, s.start_date, s.end_date, s.auto_renew, s.status, s.note,
			s.created_at, s.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM subscriptions s
		LEFT JOIN cash_flow_categories c ON s.category_id = c.id
		WHERE s.status = $1
			AND s.billing_day = $2
			AND s.start_date <= $3
			AND (s.end_date IS NULL OR s.end_date >= $3)
		ORDER BY s.name
	`

	rows, err := r.db.Query(query, models.SubscriptionStatusActive, day, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get due billings: %w", err)
	}
	defer rows.Close()

	subscriptions := []*models.Subscription{}
	for rows.Next() {
		subscription := &models.Subscription{
			Category: &models.CashFlowCategory{},
		}

		err := rows.Scan(
			&subscription.ID,
			&subscription.Name,
			&subscription.Amount,
			&subscription.Currency,
			&subscription.BillingCycle,
			&subscription.BillingDay,
			&subscription.CategoryID,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.AutoRenew,
			&subscription.Status,
			&subscription.Note,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
			&subscription.Category.ID,
			&subscription.Category.Name,
			&subscription.Category.Type,
			&subscription.Category.IsSystem,
			&subscription.Category.CreatedAt,
			&subscription.Category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subscriptions, nil
}

// GetExpiringSoon 取得即將到期的訂閱
func (r *subscriptionRepository) GetExpiringSoon(days int) ([]*models.Subscription, error) {
	now := time.Now()
	futureDate := now.AddDate(0, 0, days)

	query := `
		SELECT
			s.id, s.name, s.amount, s.currency, s.billing_cycle, s.billing_day,
			s.category_id, s.start_date, s.end_date, s.auto_renew, s.status, s.note,
			s.created_at, s.updated_at,
			c.id, c.name, c.type, c.is_system, c.created_at, c.updated_at
		FROM subscriptions s
		LEFT JOIN cash_flow_categories c ON s.category_id = c.id
		WHERE s.status = $1
			AND s.end_date IS NOT NULL
			AND s.end_date BETWEEN $2 AND $3
		ORDER BY s.end_date
	`

	rows, err := r.db.Query(query, models.SubscriptionStatusActive, now, futureDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions := []*models.Subscription{}
	for rows.Next() {
		subscription := &models.Subscription{
			Category: &models.CashFlowCategory{},
		}

		err := rows.Scan(
			&subscription.ID,
			&subscription.Name,
			&subscription.Amount,
			&subscription.Currency,
			&subscription.BillingCycle,
			&subscription.BillingDay,
			&subscription.CategoryID,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.AutoRenew,
			&subscription.Status,
			&subscription.Note,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
			&subscription.Category.ID,
			&subscription.Category.Name,
			&subscription.Category.Type,
			&subscription.Category.IsSystem,
			&subscription.Category.CreatedAt,
			&subscription.Category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subscriptions, nil
}

