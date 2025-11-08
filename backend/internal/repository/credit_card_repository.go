package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
)

// CreditCardRepository 信用卡資料存取介面
type CreditCardRepository interface {
	Create(input *models.CreateCreditCardInput) (*models.CreditCard, error)
	GetByID(id uuid.UUID) (*models.CreditCard, error)
	GetAll() ([]*models.CreditCard, error)
	GetByBillingDay(day int) ([]*models.CreditCard, error)
	GetByPaymentDueDay(day int) ([]*models.CreditCard, error)
	GetUpcomingBilling(daysAhead int) ([]*models.CreditCard, error)
	GetUpcomingPayment(daysAhead int) ([]*models.CreditCard, error)
	Update(id uuid.UUID, input *models.UpdateCreditCardInput) (*models.CreditCard, error)
	UpdateUsedCredit(id uuid.UUID, amount float64) (*models.CreditCard, error)
	Delete(id uuid.UUID) error
}

// creditCardRepository 信用卡資料存取實作
type creditCardRepository struct {
	db *sql.DB
}

// NewCreditCardRepository 建立新的信用卡 repository
func NewCreditCardRepository(db *sql.DB) CreditCardRepository {
	return &creditCardRepository{db: db}
}

// Create 建立新的信用卡
func (r *creditCardRepository) Create(input *models.CreateCreditCardInput) (*models.CreditCard, error) {
	query := `
		INSERT INTO credit_cards (issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, credit_limit, used_credit, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, credit_limit, used_credit, group_id, note, created_at, updated_at
	`

	card := &models.CreditCard{}
	err := r.db.QueryRow(
		query,
		input.IssuingBank,
		input.CardName,
		input.CardNumberLast4,
		input.BillingDay,
		input.PaymentDueDay,
		input.CreditLimit,
		input.UsedCredit,
		input.Note,
	).Scan(
		&card.ID,
		&card.IssuingBank,
		&card.CardName,
		&card.CardNumberLast4,
		&card.BillingDay,
		&card.PaymentDueDay,
		&card.CreditLimit,
		&card.UsedCredit,
		&card.GroupID,
		&card.Note,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create credit card: %w", err)
	}

	return card, nil
}

// GetByID 根據 ID 取得信用卡
func (r *creditCardRepository) GetByID(id uuid.UUID) (*models.CreditCard, error) {
	query := `
		SELECT id, issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, credit_limit, used_credit, group_id, note, created_at, updated_at
		FROM credit_cards
		WHERE id = $1
	`

	card := &models.CreditCard{}
	err := r.db.QueryRow(query, id).Scan(
		&card.ID,
		&card.IssuingBank,
		&card.CardName,
		&card.CardNumberLast4,
		&card.BillingDay,
		&card.PaymentDueDay,
		&card.CreditLimit,
		&card.UsedCredit,
		&card.GroupID,
		&card.Note,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("credit card not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get credit card: %w", err)
	}

	return card, nil
}

// GetAll 取得所有信用卡
func (r *creditCardRepository) GetAll() ([]*models.CreditCard, error) {
	query := `
		SELECT id, issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, credit_limit, used_credit, group_id, note, created_at, updated_at
		FROM credit_cards
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit cards: %w", err)
	}
	defer rows.Close()

	var cards []*models.CreditCard
	for rows.Next() {
		card := &models.CreditCard{}
		err := rows.Scan(
			&card.ID,
			&card.IssuingBank,
			&card.CardName,
			&card.CardNumberLast4,
			&card.BillingDay,
			&card.PaymentDueDay,
			&card.CreditLimit,
			&card.UsedCredit,
			&card.GroupID,
			&card.Note,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credit card: %w", err)
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating credit cards: %w", err)
	}

	return cards, nil
}

// GetByBillingDay 根據帳單日取得信用卡
func (r *creditCardRepository) GetByBillingDay(day int) ([]*models.CreditCard, error) {
	query := `
		SELECT id, issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, credit_limit, used_credit, group_id, note, created_at, updated_at
		FROM credit_cards
		WHERE billing_day = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, day)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit cards by billing day: %w", err)
	}
	defer rows.Close()

	var cards []*models.CreditCard
	for rows.Next() {
		card := &models.CreditCard{}
		err := rows.Scan(
			&card.ID,
			&card.IssuingBank,
			&card.CardName,
			&card.CardNumberLast4,
			&card.BillingDay,
			&card.PaymentDueDay,
			&card.CreditLimit,
			&card.UsedCredit,
			&card.GroupID,
			&card.Note,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credit card: %w", err)
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating credit cards: %w", err)
	}

	return cards, nil
}

// GetByPaymentDueDay 根據繳款截止日取得信用卡
func (r *creditCardRepository) GetByPaymentDueDay(day int) ([]*models.CreditCard, error) {
	query := `
		SELECT id, issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, credit_limit, used_credit, group_id, note, created_at, updated_at
		FROM credit_cards
		WHERE payment_due_day = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, day)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit cards by payment due day: %w", err)
	}
	defer rows.Close()

	var cards []*models.CreditCard
	for rows.Next() {
		card := &models.CreditCard{}
		err := rows.Scan(
			&card.ID,
			&card.IssuingBank,
			&card.CardName,
			&card.CardNumberLast4,
			&card.BillingDay,
			&card.PaymentDueDay,
			&card.CreditLimit,
			&card.UsedCredit,
			&card.GroupID,
			&card.Note,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credit card: %w", err)
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating credit cards: %w", err)
	}

	return cards, nil
}

// GetUpcomingBilling 取得即將到來的帳單日信用卡（未來 N 天內）
func (r *creditCardRepository) GetUpcomingBilling(daysAhead int) ([]*models.CreditCard, error) {
	now := time.Now()
	currentDay := now.Day()
	targetDay := now.AddDate(0, 0, daysAhead).Day()

	var cards []*models.CreditCard
	var err error

	// 如果跨月，需要分兩次查詢
	if targetDay < currentDay {
		// 查詢本月剩餘天數
		cards1, err := r.getCardsByDayRange(currentDay, 31)
		if err != nil {
			return nil, err
		}
		// 查詢下月開始到目標日
		cards2, err := r.getCardsByDayRange(1, targetDay)
		if err != nil {
			return nil, err
		}
		cards = append(cards1, cards2...)
	} else {
		// 同一個月內
		cards, err = r.getCardsByDayRange(currentDay, targetDay)
		if err != nil {
			return nil, err
		}
	}

	return cards, nil
}

// GetUpcomingPayment 取得即將到來的繳款截止日信用卡（未來 N 天內）
func (r *creditCardRepository) GetUpcomingPayment(daysAhead int) ([]*models.CreditCard, error) {
	now := time.Now()
	currentDay := now.Day()
	targetDay := now.AddDate(0, 0, daysAhead).Day()

	var cards []*models.CreditCard
	var err error

	// 如果跨月，需要分兩次查詢
	if targetDay < currentDay {
		// 查詢本月剩餘天數
		cards1, err := r.getCardsByPaymentDayRange(currentDay, 31)
		if err != nil {
			return nil, err
		}
		// 查詢下月開始到目標日
		cards2, err := r.getCardsByPaymentDayRange(1, targetDay)
		if err != nil {
			return nil, err
		}
		cards = append(cards1, cards2...)
	} else {
		// 同一個月內
		cards, err = r.getCardsByPaymentDayRange(currentDay, targetDay)
		if err != nil {
			return nil, err
		}
	}

	return cards, nil
}

// getCardsByDayRange 根據帳單日範圍取得信用卡（輔助函式）
func (r *creditCardRepository) getCardsByDayRange(startDay, endDay int) ([]*models.CreditCard, error) {
	query := `
		SELECT id, issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, credit_limit, used_credit, group_id, note, created_at, updated_at
		FROM credit_cards
		WHERE billing_day >= $1 AND billing_day <= $2
		ORDER BY billing_day ASC
	`

	rows, err := r.db.Query(query, startDay, endDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit cards by day range: %w", err)
	}
	defer rows.Close()

	return r.scanCards(rows)
}

// getCardsByPaymentDayRange 根據繳款截止日範圍取得信用卡（輔助函式）
func (r *creditCardRepository) getCardsByPaymentDayRange(startDay, endDay int) ([]*models.CreditCard, error) {
	query := `
		SELECT id, issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, credit_limit, used_credit, group_id, note, created_at, updated_at
		FROM credit_cards
		WHERE payment_due_day >= $1 AND payment_due_day <= $2
		ORDER BY payment_due_day ASC
	`

	rows, err := r.db.Query(query, startDay, endDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit cards by payment day range: %w", err)
	}
	defer rows.Close()

	return r.scanCards(rows)
}

// scanCards 掃描信用卡資料（輔助函式）
func (r *creditCardRepository) scanCards(rows *sql.Rows) ([]*models.CreditCard, error) {
	var cards []*models.CreditCard
	for rows.Next() {
		card := &models.CreditCard{}
		err := rows.Scan(
			&card.ID,
			&card.IssuingBank,
			&card.CardName,
			&card.CardNumberLast4,
			&card.BillingDay,
			&card.PaymentDueDay,
			&card.CreditLimit,
			&card.UsedCredit,
			&card.GroupID,
			&card.Note,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credit card: %w", err)
		}
		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating credit cards: %w", err)
	}

	return cards, nil
}

// Update 更新信用卡
func (r *creditCardRepository) Update(id uuid.UUID, input *models.UpdateCreditCardInput) (*models.CreditCard, error) {
	// 動態建立 UPDATE 語句
	var setClauses []string
	var args []interface{}
	argPosition := 1

	if input.IssuingBank != nil {
		setClauses = append(setClauses, fmt.Sprintf("issuing_bank = $%d", argPosition))
		args = append(args, *input.IssuingBank)
		argPosition++
	}
	if input.CardName != nil {
		setClauses = append(setClauses, fmt.Sprintf("card_name = $%d", argPosition))
		args = append(args, *input.CardName)
		argPosition++
	}
	if input.CardNumberLast4 != nil {
		setClauses = append(setClauses, fmt.Sprintf("card_number_last4 = $%d", argPosition))
		args = append(args, *input.CardNumberLast4)
		argPosition++
	}
	if input.BillingDay != nil {
		setClauses = append(setClauses, fmt.Sprintf("billing_day = $%d", argPosition))
		args = append(args, *input.BillingDay)
		argPosition++
	}
	if input.PaymentDueDay != nil {
		setClauses = append(setClauses, fmt.Sprintf("payment_due_day = $%d", argPosition))
		args = append(args, *input.PaymentDueDay)
		argPosition++
	}
	if input.CreditLimit != nil {
		setClauses = append(setClauses, fmt.Sprintf("credit_limit = $%d", argPosition))
		args = append(args, *input.CreditLimit)
		argPosition++
	}
	if input.UsedCredit != nil {
		setClauses = append(setClauses, fmt.Sprintf("used_credit = $%d", argPosition))
		args = append(args, *input.UsedCredit)
		argPosition++
	}
	if input.Note != nil {
		setClauses = append(setClauses, fmt.Sprintf("note = $%d", argPosition))
		args = append(args, *input.Note)
		argPosition++
	}

	// 如果沒有任何欄位需要更新，直接返回現有資料
	if len(setClauses) == 0 {
		return r.GetByID(id)
	}

	// 加入 ID 參數
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE credit_cards
		SET %s
		WHERE id = $%d
		RETURNING id, issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, credit_limit, used_credit, group_id, note, created_at, updated_at
	`, strings.Join(setClauses, ", "), argPosition)

	card := &models.CreditCard{}
	err := r.db.QueryRow(query, args...).Scan(
		&card.ID,
		&card.IssuingBank,
		&card.CardName,
		&card.CardNumberLast4,
		&card.BillingDay,
		&card.PaymentDueDay,
		&card.CreditLimit,
		&card.UsedCredit,
		&card.GroupID,
		&card.Note,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("credit card not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update credit card: %w", err)
	}

	return card, nil
}

// UpdateUsedCredit 更新信用卡已使用額度（增加或減少指定金額）
func (r *creditCardRepository) UpdateUsedCredit(id uuid.UUID, amount float64) (*models.CreditCard, error) {
	query := `
		UPDATE credit_cards
		SET used_credit = used_credit + $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, credit_limit, used_credit, group_id, note, created_at, updated_at
	`

	card := &models.CreditCard{}
	err := r.db.QueryRow(query, amount, id).Scan(
		&card.ID,
		&card.IssuingBank,
		&card.CardName,
		&card.CardNumberLast4,
		&card.BillingDay,
		&card.PaymentDueDay,
		&card.CreditLimit,
		&card.UsedCredit,
		&card.GroupID,
		&card.Note,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("credit card not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update credit card used credit: %w", err)
	}

	return card, nil
}

// Delete 刪除信用卡
func (r *creditCardRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM credit_cards WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete credit card: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("credit card not found")
	}

	return nil
}

