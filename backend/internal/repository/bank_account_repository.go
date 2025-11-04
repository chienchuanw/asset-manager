package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
)

// BankAccountRepository 銀行帳戶資料存取介面
type BankAccountRepository interface {
	Create(input *models.CreateBankAccountInput) (*models.BankAccount, error)
	GetByID(id uuid.UUID) (*models.BankAccount, error)
	GetAll(currency *models.Currency) ([]*models.BankAccount, error)
	Update(id uuid.UUID, input *models.UpdateBankAccountInput) (*models.BankAccount, error)
	UpdateBalance(id uuid.UUID, amount float64) (*models.BankAccount, error)
	Delete(id uuid.UUID) error
}

// bankAccountRepository 銀行帳戶資料存取實作
type bankAccountRepository struct {
	db *sql.DB
}

// NewBankAccountRepository 建立新的銀行帳戶 repository
func NewBankAccountRepository(db *sql.DB) BankAccountRepository {
	return &bankAccountRepository{db: db}
}

// Create 建立新的銀行帳戶
func (r *bankAccountRepository) Create(input *models.CreateBankAccountInput) (*models.BankAccount, error) {
	query := `
		INSERT INTO bank_accounts (bank_name, account_type, account_number_last4, currency, balance, note)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, bank_name, account_type, account_number_last4, currency, balance, note, created_at, updated_at
	`

	account := &models.BankAccount{}
	err := r.db.QueryRow(
		query,
		input.BankName,
		input.AccountType,
		input.AccountNumberLast4,
		input.Currency,
		input.Balance,
		input.Note,
	).Scan(
		&account.ID,
		&account.BankName,
		&account.AccountType,
		&account.AccountNumberLast4,
		&account.Currency,
		&account.Balance,
		&account.Note,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create bank account: %w", err)
	}

	return account, nil
}

// GetByID 根據 ID 取得銀行帳戶
func (r *bankAccountRepository) GetByID(id uuid.UUID) (*models.BankAccount, error) {
	query := `
		SELECT id, bank_name, account_type, account_number_last4, currency, balance, note, created_at, updated_at
		FROM bank_accounts
		WHERE id = $1
	`

	account := &models.BankAccount{}
	err := r.db.QueryRow(query, id).Scan(
		&account.ID,
		&account.BankName,
		&account.AccountType,
		&account.AccountNumberLast4,
		&account.Currency,
		&account.Balance,
		&account.Note,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("bank account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get bank account: %w", err)
	}

	return account, nil
}

// GetAll 取得所有銀行帳戶（可選擇性篩選幣別）
func (r *bankAccountRepository) GetAll(currency *models.Currency) ([]*models.BankAccount, error) {
	query := `
		SELECT id, bank_name, account_type, account_number_last4, currency, balance, note, created_at, updated_at
		FROM bank_accounts
	`

	var args []interface{}
	var conditions []string

	// 如果有指定幣別，加入篩選條件
	if currency != nil {
		conditions = append(conditions, "currency = $1")
		args = append(args, *currency)
	}

	// 組合 WHERE 條件
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// 排序：依建立時間降序
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*models.BankAccount
	for rows.Next() {
		account := &models.BankAccount{}
		err := rows.Scan(
			&account.ID,
			&account.BankName,
			&account.AccountType,
			&account.AccountNumberLast4,
			&account.Currency,
			&account.Balance,
			&account.Note,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bank account: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bank accounts: %w", err)
	}

	return accounts, nil
}

// Update 更新銀行帳戶
func (r *bankAccountRepository) Update(id uuid.UUID, input *models.UpdateBankAccountInput) (*models.BankAccount, error) {
	// 動態建立 UPDATE 語句
	var setClauses []string
	var args []interface{}
	argPosition := 1

	if input.BankName != nil {
		setClauses = append(setClauses, fmt.Sprintf("bank_name = $%d", argPosition))
		args = append(args, *input.BankName)
		argPosition++
	}
	if input.AccountType != nil {
		setClauses = append(setClauses, fmt.Sprintf("account_type = $%d", argPosition))
		args = append(args, *input.AccountType)
		argPosition++
	}
	if input.AccountNumberLast4 != nil {
		setClauses = append(setClauses, fmt.Sprintf("account_number_last4 = $%d", argPosition))
		args = append(args, *input.AccountNumberLast4)
		argPosition++
	}
	if input.Currency != nil {
		setClauses = append(setClauses, fmt.Sprintf("currency = $%d", argPosition))
		args = append(args, *input.Currency)
		argPosition++
	}
	if input.Balance != nil {
		setClauses = append(setClauses, fmt.Sprintf("balance = $%d", argPosition))
		args = append(args, *input.Balance)
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
		UPDATE bank_accounts
		SET %s
		WHERE id = $%d
		RETURNING id, bank_name, account_type, account_number_last4, currency, balance, note, created_at, updated_at
	`, strings.Join(setClauses, ", "), argPosition)

	account := &models.BankAccount{}
	err := r.db.QueryRow(query, args...).Scan(
		&account.ID,
		&account.BankName,
		&account.AccountType,
		&account.AccountNumberLast4,
		&account.Currency,
		&account.Balance,
		&account.Note,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("bank account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update bank account: %w", err)
	}

	return account, nil
}

// UpdateBalance 更新銀行帳戶餘額（增加或減少指定金額）
func (r *bankAccountRepository) UpdateBalance(id uuid.UUID, amount float64) (*models.BankAccount, error) {
	query := `
		UPDATE bank_accounts
		SET balance = balance + $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, bank_name, account_type, account_number_last4, currency, balance, note, created_at, updated_at
	`

	account := &models.BankAccount{}
	err := r.db.QueryRow(query, amount, id).Scan(
		&account.ID,
		&account.BankName,
		&account.AccountType,
		&account.AccountNumberLast4,
		&account.Currency,
		&account.Balance,
		&account.Note,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("bank account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update bank account balance: %w", err)
	}

	return account, nil
}

// Delete 刪除銀行帳戶
func (r *bankAccountRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM bank_accounts WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete bank account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("bank account not found")
	}

	return nil
}

