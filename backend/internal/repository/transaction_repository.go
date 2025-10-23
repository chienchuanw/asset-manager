package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
)

// TransactionRepository 交易記錄資料存取介面
type TransactionRepository interface {
	Create(input *models.CreateTransactionInput) (*models.Transaction, error)
	GetByID(id uuid.UUID) (*models.Transaction, error)
	GetAll(filters TransactionFilters) ([]*models.Transaction, error)
	Update(id uuid.UUID, input *models.UpdateTransactionInput) (*models.Transaction, error)
	Delete(id uuid.UUID) error
}

// TransactionFilters 查詢篩選條件
type TransactionFilters struct {
	AssetType       *models.AssetType       `json:"asset_type,omitempty"`
	TransactionType *models.TransactionType `json:"transaction_type,omitempty"`
	Symbol          *string                 `json:"symbol,omitempty"`
	StartDate       *time.Time              `json:"start_date,omitempty"`
	EndDate         *time.Time              `json:"end_date,omitempty"`
	Limit           int                     `json:"limit,omitempty"`
	Offset          int                     `json:"offset,omitempty"`
}

// transactionRepository 交易記錄資料存取實作
type transactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository 建立新的交易記錄 repository
func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Create 建立新的交易記錄
func (r *transactionRepository) Create(input *models.CreateTransactionInput) (*models.Transaction, error) {
	query := `
		INSERT INTO transactions (date, asset_type, symbol, name, transaction_type, quantity, price, amount, fee, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, date, asset_type, symbol, name, transaction_type, quantity, price, amount, fee, note, created_at, updated_at
	`

	transaction := &models.Transaction{}
	err := r.db.QueryRow(
		query,
		input.Date,
		input.AssetType,
		input.Symbol,
		input.Name,
		input.TransactionType,
		input.Quantity,
		input.Price,
		input.Amount,
		input.Fee,
		input.Note,
	).Scan(
		&transaction.ID,
		&transaction.Date,
		&transaction.AssetType,
		&transaction.Symbol,
		&transaction.Name,
		&transaction.TransactionType,
		&transaction.Quantity,
		&transaction.Price,
		&transaction.Amount,
		&transaction.Fee,
		&transaction.Note,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

// GetByID 根據 ID 取得交易記錄
func (r *transactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	query := `
		SELECT id, date, asset_type, symbol, name, transaction_type, quantity, price, amount, fee, note, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	transaction := &models.Transaction{}
	err := r.db.QueryRow(query, id).Scan(
		&transaction.ID,
		&transaction.Date,
		&transaction.AssetType,
		&transaction.Symbol,
		&transaction.Name,
		&transaction.TransactionType,
		&transaction.Quantity,
		&transaction.Price,
		&transaction.Amount,
		&transaction.Fee,
		&transaction.Note,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transaction not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

// GetAll 取得所有交易記錄（支援篩選）
func (r *transactionRepository) GetAll(filters TransactionFilters) ([]*models.Transaction, error) {
	query := `
		SELECT id, date, asset_type, symbol, name, transaction_type, quantity, price, amount, fee, note, created_at, updated_at
		FROM transactions
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

	if filters.TransactionType != nil {
		query += fmt.Sprintf(" AND transaction_type = $%d", argCount)
		args = append(args, *filters.TransactionType)
		argCount++
	}

	if filters.Symbol != nil {
		query += fmt.Sprintf(" AND symbol = $%d", argCount)
		args = append(args, *filters.Symbol)
		argCount++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND date >= $%d", argCount)
		args = append(args, *filters.StartDate)
		argCount++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND date <= $%d", argCount)
		args = append(args, *filters.EndDate)
		argCount++
	}

	// 排序
	query += " ORDER BY date DESC, created_at DESC"

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
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*models.Transaction{}
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.Date,
			&transaction.AssetType,
			&transaction.Symbol,
			&transaction.Name,
			&transaction.TransactionType,
			&transaction.Quantity,
			&transaction.Price,
			&transaction.Amount,
			&transaction.Fee,
			&transaction.Note,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// Update 更新交易記錄
func (r *transactionRepository) Update(id uuid.UUID, input *models.UpdateTransactionInput) (*models.Transaction, error) {
	// 動態建立 UPDATE 語句
	setClauses := []string{}
	args := []interface{}{}
	argCount := 1

	if input.Date != nil {
		setClauses = append(setClauses, fmt.Sprintf("date = $%d", argCount))
		args = append(args, *input.Date)
		argCount++
	}

	if input.AssetType != nil {
		setClauses = append(setClauses, fmt.Sprintf("asset_type = $%d", argCount))
		args = append(args, *input.AssetType)
		argCount++
	}

	if input.Symbol != nil {
		setClauses = append(setClauses, fmt.Sprintf("symbol = $%d", argCount))
		args = append(args, *input.Symbol)
		argCount++
	}

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *input.Name)
		argCount++
	}

	if input.TransactionType != nil {
		setClauses = append(setClauses, fmt.Sprintf("transaction_type = $%d", argCount))
		args = append(args, *input.TransactionType)
		argCount++
	}

	if input.Quantity != nil {
		setClauses = append(setClauses, fmt.Sprintf("quantity = $%d", argCount))
		args = append(args, *input.Quantity)
		argCount++
	}

	if input.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argCount))
		args = append(args, *input.Price)
		argCount++
	}

	if input.Amount != nil {
		setClauses = append(setClauses, fmt.Sprintf("amount = $%d", argCount))
		args = append(args, *input.Amount)
		argCount++
	}

	if input.Fee != nil {
		setClauses = append(setClauses, fmt.Sprintf("fee = $%d", argCount))
		args = append(args, *input.Fee)
		argCount++
	}

	if input.Note != nil {
		setClauses = append(setClauses, fmt.Sprintf("note = $%d", argCount))
		args = append(args, *input.Note)
		argCount++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// 加入 ID 參數
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE transactions
		SET %s
		WHERE id = $%d
		RETURNING id, date, asset_type, symbol, name, transaction_type, quantity, price, amount, fee, note, created_at, updated_at
	`, strings.Join(setClauses, ", "), argCount)

	transaction := &models.Transaction{}
	err := r.db.QueryRow(query, args...).Scan(
		&transaction.ID,
		&transaction.Date,
		&transaction.AssetType,
		&transaction.Symbol,
		&transaction.Name,
		&transaction.TransactionType,
		&transaction.Quantity,
		&transaction.Price,
		&transaction.Amount,
		&transaction.Fee,
		&transaction.Note,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transaction not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return transaction, nil
}

// Delete 刪除交易記錄
func (r *transactionRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM transactions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

