package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
)

// CreditCardGroupRepository 信用卡群組資料存取介面
type CreditCardGroupRepository interface {
	Create(input *models.CreateCreditCardGroupInput) (*models.CreditCardGroup, error)
	GetByID(id uuid.UUID) (*models.CreditCardGroupWithCards, error)
	GetAll() ([]*models.CreditCardGroupWithCards, error)
	Update(id uuid.UUID, input *models.UpdateCreditCardGroupInput) (*models.CreditCardGroup, error)
	Delete(id uuid.UUID) error
	AddCardsToGroup(groupID uuid.UUID, cardIDs []uuid.UUID) error
	RemoveCardsFromGroup(cardIDs []uuid.UUID) error
	GetCardsByGroupID(groupID uuid.UUID) ([]*models.CreditCard, error)
}

// creditCardGroupRepository 信用卡群組資料存取實作
type creditCardGroupRepository struct {
	db *sql.DB
}

// NewCreditCardGroupRepository 建立新的信用卡群組 repository
func NewCreditCardGroupRepository(db *sql.DB) CreditCardGroupRepository {
	return &creditCardGroupRepository{db: db}
}

// Create 建立新的信用卡群組
func (r *creditCardGroupRepository) Create(input *models.CreateCreditCardGroupInput) (*models.CreditCardGroup, error) {
	query := `
		INSERT INTO credit_card_groups (name, issuing_bank, shared_credit_limit, note)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, issuing_bank, shared_credit_limit, note, created_at, updated_at
	`

	group := &models.CreditCardGroup{}
	err := r.db.QueryRow(
		query,
		input.Name,
		input.IssuingBank,
		input.SharedCreditLimit,
		input.Note,
	).Scan(
		&group.ID,
		&group.Name,
		&group.IssuingBank,
		&group.SharedCreditLimit,
		&group.Note,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create credit card group: %w", err)
	}

	return group, nil
}

// GetByID 根據 ID 取得信用卡群組及其包含的卡片
func (r *creditCardGroupRepository) GetByID(id uuid.UUID) (*models.CreditCardGroupWithCards, error) {
	// 取得群組基本資料
	groupQuery := `
		SELECT id, name, issuing_bank, shared_credit_limit, note, created_at, updated_at
		FROM credit_card_groups
		WHERE id = $1
	`

	group := &models.CreditCardGroupWithCards{}
	err := r.db.QueryRow(groupQuery, id).Scan(
		&group.ID,
		&group.Name,
		&group.IssuingBank,
		&group.SharedCreditLimit,
		&group.Note,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("credit card group not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get credit card group: %w", err)
	}

	// 取得群組內的卡片
	cards, err := r.GetCardsByGroupID(id)
	if err != nil {
		return nil, err
	}

	group.Cards = cards

	// 計算總已使用額度
	var totalUsedCredit float64
	for _, card := range cards {
		totalUsedCredit += card.UsedCredit
	}
	group.TotalUsedCredit = totalUsedCredit

	return group, nil
}

// GetAll 取得所有信用卡群組
func (r *creditCardGroupRepository) GetAll() ([]*models.CreditCardGroupWithCards, error) {
	query := `
		SELECT id, name, issuing_bank, shared_credit_limit, note, created_at, updated_at
		FROM credit_card_groups
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit card groups: %w", err)
	}
	defer rows.Close()

	var groups []*models.CreditCardGroupWithCards
	for rows.Next() {
		group := &models.CreditCardGroupWithCards{}
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.IssuingBank,
			&group.SharedCreditLimit,
			&group.Note,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credit card group: %w", err)
		}

		// 取得群組內的卡片
		cards, err := r.GetCardsByGroupID(group.ID)
		if err != nil {
			return nil, err
		}
		group.Cards = cards

		// 計算總已使用額度
		var totalUsedCredit float64
		for _, card := range cards {
			totalUsedCredit += card.UsedCredit
		}
		group.TotalUsedCredit = totalUsedCredit

		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating credit card groups: %w", err)
	}

	return groups, nil
}

// Update 更新信用卡群組
func (r *creditCardGroupRepository) Update(id uuid.UUID, input *models.UpdateCreditCardGroupInput) (*models.CreditCardGroup, error) {
	// 動態建立 UPDATE 語句
	var setClauses []string
	var args []interface{}
	argPosition := 1

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPosition))
		args = append(args, *input.Name)
		argPosition++
	}
	if input.SharedCreditLimit != nil {
		setClauses = append(setClauses, fmt.Sprintf("shared_credit_limit = $%d", argPosition))
		args = append(args, *input.SharedCreditLimit)
		argPosition++
	}
	if input.Note != nil {
		setClauses = append(setClauses, fmt.Sprintf("note = $%d", argPosition))
		args = append(args, *input.Note)
		argPosition++
	}

	// 如果沒有任何欄位需要更新,直接返回現有資料
	if len(setClauses) == 0 {
		groupWithCards, err := r.GetByID(id)
		if err != nil {
			return nil, err
		}
		return &groupWithCards.CreditCardGroup, nil
	}

	// 加入 updated_at
	setClauses = append(setClauses, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))

	// 加入 ID 參數
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE credit_card_groups
		SET %s
		WHERE id = $%d
		RETURNING id, name, issuing_bank, shared_credit_limit, note, created_at, updated_at
	`, strings.Join(setClauses, ", "), argPosition)

	group := &models.CreditCardGroup{}
	err := r.db.QueryRow(query, args...).Scan(
		&group.ID,
		&group.Name,
		&group.IssuingBank,
		&group.SharedCreditLimit,
		&group.Note,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("credit card group not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update credit card group: %w", err)
	}

	return group, nil
}

// Delete 刪除信用卡群組
func (r *creditCardGroupRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM credit_card_groups WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete credit card group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("credit card group not found")
	}

	return nil
}

// AddCardsToGroup 將卡片加入群組
func (r *creditCardGroupRepository) AddCardsToGroup(groupID uuid.UUID, cardIDs []uuid.UUID) error {
	// 使用 transaction 確保資料一致性
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE credit_cards
		SET group_id = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	for _, cardID := range cardIDs {
		_, err := tx.Exec(query, groupID, cardID)
		if err != nil {
			return fmt.Errorf("failed to add card to group: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RemoveCardsFromGroup 從群組移除卡片
func (r *creditCardGroupRepository) RemoveCardsFromGroup(cardIDs []uuid.UUID) error {
	// 使用 transaction 確保資料一致性
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE credit_cards
		SET group_id = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	for _, cardID := range cardIDs {
		_, err := tx.Exec(query, cardID)
		if err != nil {
			return fmt.Errorf("failed to remove card from group: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetCardsByGroupID 取得群組內的所有卡片
func (r *creditCardGroupRepository) GetCardsByGroupID(groupID uuid.UUID) ([]*models.CreditCard, error) {
	query := `
		SELECT id, issuing_bank, card_name, card_number_last4, billing_day, payment_due_day, 
		       credit_limit, used_credit, group_id, note, created_at, updated_at
		FROM credit_cards
		WHERE group_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards by group ID: %w", err)
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

