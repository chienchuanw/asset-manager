package service

import (
	"database/sql"
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// ReconciliationService 校準操作業務邏輯
type ReconciliationService interface {
	ReconcileBankAccount(id uuid.UUID, input *models.ReconcileBankAccountInput) (*models.ReconcileResult, error)
	ReconcileCreditCard(id uuid.UUID, input *models.ReconcileCreditCardInput) (*models.ReconcileResult, error)
	ReconcileBatch(input *models.ReconcileBatchInput) ([]*models.ReconcileResult, error)
}

type reconciliationService struct {
	db           *sql.DB
	categoryRepo repository.CategoryRepository
}

// NewReconciliationService 建立 ReconciliationService
func NewReconciliationService(db *sql.DB, categoryRepo repository.CategoryRepository) ReconciliationService {
	return &reconciliationService{db: db, categoryRepo: categoryRepo}
}

// adjustmentCategoryIDs 暫存查到的兩個調整分類 ID
type adjustmentCategoryIDs struct {
	income  uuid.UUID
	expense uuid.UUID
}

func (s *reconciliationService) loadAdjustmentCategories() (*adjustmentCategoryIDs, error) {
	in, err := s.categoryRepo.GetByNameAndType("餘額調整-收入", models.CashFlowTypeIncome)
	if err != nil {
		return nil, fmt.Errorf("load income adjustment category: %w", err)
	}
	out, err := s.categoryRepo.GetByNameAndType("餘額調整-支出", models.CashFlowTypeExpense)
	if err != nil {
		return nil, fmt.Errorf("load expense adjustment category: %w", err)
	}
	return &adjustmentCategoryIDs{income: in.ID, expense: out.ID}, nil
}

// ReconcileBankAccount 校準銀行帳戶餘額
func (s *reconciliationService) ReconcileBankAccount(id uuid.UUID, input *models.ReconcileBankAccountInput) (*models.ReconcileResult, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	cats, err := s.loadAdjustmentCategories()
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := s.reconcileBankAccountTx(tx, id, input, cats)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return res, nil
}

// ReconcileCreditCard 校準信用卡 used_credit / credit_limit
func (s *reconciliationService) ReconcileCreditCard(id uuid.UUID, input *models.ReconcileCreditCardInput) (*models.ReconcileResult, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	cats, err := s.loadAdjustmentCategories()
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := s.reconcileCreditCardTx(tx, id, input, cats)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return res, nil
}

// reconcileCreditCardTx 在指定 tx 內校準單一信用卡
func (s *reconciliationService) reconcileCreditCardTx(
	tx *sql.Tx,
	id uuid.UUID,
	input *models.ReconcileCreditCardInput,
	cats *adjustmentCategoryIDs,
) (*models.ReconcileResult, error) {
	var (
		currUsed    float64
		currLimit   float64
		issuingBank string
		cardName    string
		last4       string
	)
	err := tx.QueryRow(
		`SELECT used_credit, credit_limit, issuing_bank, card_name, card_number_last4
		 FROM credit_cards WHERE id = $1`,
		id,
	).Scan(&currUsed, &currLimit, &issuingBank, &cardName, &last4)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("credit card not found")
	}
	if err != nil {
		return nil, fmt.Errorf("fetch credit card: %w", err)
	}

	nextUsed := currUsed
	if input.NewUsedCredit != nil {
		nextUsed = *input.NewUsedCredit
	}
	nextLimit := currLimit
	if input.NewCreditLimit != nil {
		nextLimit = *input.NewCreditLimit
	}
	if nextUsed > nextLimit {
		return nil, fmt.Errorf("used_credit (%v) cannot exceed credit_limit (%v)", nextUsed, nextLimit)
	}

	if _, err := tx.Exec(
		`UPDATE credit_cards
		 SET used_credit = $1, credit_limit = $2, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $3`,
		nextUsed, nextLimit, id,
	); err != nil {
		return nil, fmt.Errorf("update credit card: %w", err)
	}

	delta := nextUsed - currUsed
	res := &models.ReconcileResult{
		TargetType:     models.SourceTypeCreditCard,
		TargetID:       id,
		Delta:          delta,
		NewUsedCredit:  ptrFloat(nextUsed),
		NewCreditLimit: ptrFloat(nextLimit),
	}
	if delta == 0 {
		return res, nil
	}

	flowType := models.CashFlowTypeExpense
	categoryID := cats.expense
	amount := delta
	if delta < 0 {
		flowType = models.CashFlowTypeIncome
		categoryID = cats.income
		amount = -delta
	}
	description := fmt.Sprintf("[餘額調整] %s %s (****%s)", issuingBank, cardName, last4)
	srcType := models.SourceTypeCreditCard
	srcID := id

	var newID uuid.UUID
	err = tx.QueryRow(
		`INSERT INTO cash_flows
		 (date, type, category_id, amount, currency, description, note, source_type, source_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id`,
		input.Date, flowType, categoryID, amount, models.CurrencyTWD,
		description, input.Note, srcType, srcID,
	).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("insert adjustment cash flow: %w", err)
	}
	res.CashFlowID = &newID
	return res, nil
}

// ReconcileBatch 批次校準（共用單一 transaction）
func (s *reconciliationService) ReconcileBatch(input *models.ReconcileBatchInput) ([]*models.ReconcileResult, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	cats, err := s.loadAdjustmentCategories()
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	results := make([]*models.ReconcileResult, 0, len(input.Items))
	for i, item := range input.Items {
		switch item.TargetType {
		case models.SourceTypeBankAccount:
			r, err := s.reconcileBankAccountTx(tx, item.TargetID, &models.ReconcileBankAccountInput{
				NewBalance: *item.NewBalance,
				Date:       input.Date,
				Note:       input.Note,
			}, cats)
			if err != nil {
				return nil, fmt.Errorf("items[%d]: %w", i, err)
			}
			results = append(results, r)
		case models.SourceTypeCreditCard:
			r, err := s.reconcileCreditCardTx(tx, item.TargetID, &models.ReconcileCreditCardInput{
				NewUsedCredit:  item.NewUsedCredit,
				NewCreditLimit: item.NewCreditLimit,
				Date:           input.Date,
				Note:           input.Note,
			}, cats)
			if err != nil {
				return nil, fmt.Errorf("items[%d]: %w", i, err)
			}
			results = append(results, r)
		default:
			return nil, fmt.Errorf("items[%d]: unsupported target_type %s", i, item.TargetType)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return results, nil
}

// reconcileBankAccountTx 在指定 tx 內校準單一銀行帳戶
func (s *reconciliationService) reconcileBankAccountTx(
	tx *sql.Tx,
	id uuid.UUID,
	input *models.ReconcileBankAccountInput,
	cats *adjustmentCategoryIDs,
) (*models.ReconcileResult, error) {
	var (
		currentBalance float64
		bankName       string
		last4          string
	)
	err := tx.QueryRow(
		`SELECT balance, bank_name, account_number_last4 FROM bank_accounts WHERE id = $1`,
		id,
	).Scan(&currentBalance, &bankName, &last4)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("bank account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("fetch bank account: %w", err)
	}

	delta := input.NewBalance - currentBalance

	if _, err := tx.Exec(
		`UPDATE bank_accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		input.NewBalance, id,
	); err != nil {
		return nil, fmt.Errorf("update bank account balance: %w", err)
	}

	res := &models.ReconcileResult{
		TargetType: models.SourceTypeBankAccount,
		TargetID:   id,
		Delta:      delta,
		NewBalance: ptrFloat(input.NewBalance),
	}

	if delta == 0 {
		return res, nil
	}

	flowType := models.CashFlowTypeIncome
	categoryID := cats.income
	amount := delta
	if delta < 0 {
		flowType = models.CashFlowTypeExpense
		categoryID = cats.expense
		amount = -delta
	}
	description := fmt.Sprintf("[餘額調整] %s (****%s)", bankName, last4)
	srcType := models.SourceTypeBankAccount
	srcID := id

	var newID uuid.UUID
	err = tx.QueryRow(
		`INSERT INTO cash_flows
		 (date, type, category_id, amount, currency, description, note, source_type, source_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id`,
		input.Date, flowType, categoryID, amount, models.CurrencyTWD,
		description, input.Note, srcType, srcID,
	).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("insert adjustment cash flow: %w", err)
	}
	res.CashFlowID = &newID
	return res, nil
}

func ptrFloat(v float64) *float64 { return &v }
