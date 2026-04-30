package service

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// reconciliationTestEnv holds the shared DB + repos for service-layer integration tests.
type reconciliationTestEnv struct {
	db       *sql.DB
	svc      ReconciliationService
	bankRepo repository.BankAccountRepository
	cardRepo repository.CreditCardRepository
	cfRepo   repository.CashFlowRepository
}

func newReconciliationTestEnv(t *testing.T) *reconciliationTestEnv {
	t.Helper()
	host := envOrDefault("TEST_DB_HOST", "localhost")
	port := envOrDefault("TEST_DB_PORT", "5432")
	user := envOrDefault("TEST_DB_USER", "postgres")
	pwd := envOrDefault("TEST_DB_PASSWORD", "postgres")
	name := envOrDefault("TEST_DB_NAME", "asset_manager_test")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pwd, name)
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	require.NoError(t, db.Ping())

	// Ensure adjustment categories exist (idempotent).
	_, err = db.Exec(`
		INSERT INTO cash_flow_categories (name, type, is_system, sort_order)
		VALUES ('餘額調整-收入', 'income', true, 9000),
		       ('餘額調整-支出', 'expense', true, 9000)
		ON CONFLICT (name, type) DO NOTHING
	`)
	require.NoError(t, err)

	// Wipe per-test state.
	for _, tbl := range []string{"cash_flows", "bank_accounts", "credit_cards"} {
		_, err := db.Exec("DELETE FROM " + tbl)
		require.NoError(t, err)
	}

	categoryRepo := repository.NewCategoryRepository(db)
	return &reconciliationTestEnv{
		db:       db,
		svc:      NewReconciliationService(db, categoryRepo),
		bankRepo: repository.NewBankAccountRepository(db),
		cardRepo: repository.NewCreditCardRepository(db),
		cfRepo:   repository.NewCashFlowRepository(db),
	}
}

func envOrDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func seedBankAccount(t *testing.T, env *reconciliationTestEnv, balance float64) *models.BankAccount {
	t.Helper()
	acc, err := env.bankRepo.Create(&models.CreateBankAccountInput{
		BankName:           "玉山",
		AccountType:        "活存",
		AccountNumberLast4: "1234",
		Currency:           models.CurrencyTWD,
		Balance:            balance,
	})
	require.NoError(t, err)
	return acc
}

func TestReconcileBankAccount_DeltaPositive_CreatesIncomeCashFlow(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	acc := seedBankAccount(t, env, 30000)
	res, err := env.svc.ReconcileBankAccount(acc.ID, &models.ReconcileBankAccountInput{
		NewBalance: 45000,
		Date:       time.Date(2026, 4, 29, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	assert.Equal(t, 15000.0, res.Delta)
	require.NotNil(t, res.CashFlowID)

	updated, err := env.bankRepo.GetByID(acc.ID)
	require.NoError(t, err)
	assert.Equal(t, 45000.0, updated.Balance)

	cf, err := env.cfRepo.GetByID(*res.CashFlowID)
	require.NoError(t, err)
	assert.Equal(t, models.CashFlowTypeIncome, cf.Type)
	assert.Equal(t, 15000.0, cf.Amount)
	require.NotNil(t, cf.SourceType)
	assert.Equal(t, models.SourceTypeBankAccount, *cf.SourceType)
	require.NotNil(t, cf.SourceID)
	assert.Equal(t, acc.ID, *cf.SourceID)
}

func TestReconcileBankAccount_DeltaNegative_CreatesExpenseCashFlow(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	acc := seedBankAccount(t, env, 30000)
	res, err := env.svc.ReconcileBankAccount(acc.ID, &models.ReconcileBankAccountInput{
		NewBalance: 20000,
		Date:       time.Now(),
	})
	require.NoError(t, err)
	assert.Equal(t, -10000.0, res.Delta)
	require.NotNil(t, res.CashFlowID)

	cf, err := env.cfRepo.GetByID(*res.CashFlowID)
	require.NoError(t, err)
	assert.Equal(t, models.CashFlowTypeExpense, cf.Type)
	assert.Equal(t, 10000.0, cf.Amount)
}

func TestReconcileBankAccount_DeltaZero_NoCashFlow(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	acc := seedBankAccount(t, env, 30000)
	res, err := env.svc.ReconcileBankAccount(acc.ID, &models.ReconcileBankAccountInput{
		NewBalance: 30000,
		Date:       time.Now(),
	})
	require.NoError(t, err)
	assert.Equal(t, 0.0, res.Delta)
	assert.Nil(t, res.CashFlowID)
}

func TestReconcileBankAccount_NotFound(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	_, err := env.svc.ReconcileBankAccount(uuid.New(), &models.ReconcileBankAccountInput{
		NewBalance: 100,
		Date:       time.Now(),
	})
	assert.Error(t, err)
}

func seedCreditCard(t *testing.T, env *reconciliationTestEnv, used, limit float64) *models.CreditCard {
	t.Helper()
	card, err := env.cardRepo.Create(&models.CreateCreditCardInput{
		IssuingBank:     "玉山",
		CardName:        "Pi 卡",
		CardNumberLast4: "5678",
		BillingDay:      5,
		PaymentDueDay:   15,
		CreditLimit:     limit,
		UsedCredit:      used,
	})
	require.NoError(t, err)
	return card
}

func TestReconcileCreditCard_UsedCreditIncrease_CreatesExpenseCashFlow(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	card := seedCreditCard(t, env, 5000, 100000)
	newUsed := 8000.0
	res, err := env.svc.ReconcileCreditCard(card.ID, &models.ReconcileCreditCardInput{
		NewUsedCredit: &newUsed,
		Date:          time.Now(),
	})
	require.NoError(t, err)
	assert.Equal(t, 3000.0, res.Delta)
	require.NotNil(t, res.CashFlowID)

	cf, err := env.cfRepo.GetByID(*res.CashFlowID)
	require.NoError(t, err)
	assert.Equal(t, models.CashFlowTypeExpense, cf.Type)
	assert.Equal(t, 3000.0, cf.Amount)

	updated, err := env.cardRepo.GetByID(card.ID)
	require.NoError(t, err)
	assert.Equal(t, 8000.0, updated.UsedCredit)
}

func TestReconcileCreditCard_UsedCreditDecrease_CreatesIncomeCashFlow(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	card := seedCreditCard(t, env, 5000, 100000)
	newUsed := 2000.0
	res, err := env.svc.ReconcileCreditCard(card.ID, &models.ReconcileCreditCardInput{
		NewUsedCredit: &newUsed,
		Date:          time.Now(),
	})
	require.NoError(t, err)
	assert.Equal(t, -3000.0, res.Delta)
	require.NotNil(t, res.CashFlowID)

	cf, err := env.cfRepo.GetByID(*res.CashFlowID)
	require.NoError(t, err)
	assert.Equal(t, models.CashFlowTypeIncome, cf.Type)
	assert.Equal(t, 3000.0, cf.Amount)
}

func TestReconcileCreditCard_OnlyCreditLimitChanged_NoCashFlow(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	card := seedCreditCard(t, env, 5000, 100000)
	newLimit := 120000.0
	res, err := env.svc.ReconcileCreditCard(card.ID, &models.ReconcileCreditCardInput{
		NewCreditLimit: &newLimit,
		Date:           time.Now(),
	})
	require.NoError(t, err)
	assert.Nil(t, res.CashFlowID)

	updated, err := env.cardRepo.GetByID(card.ID)
	require.NoError(t, err)
	assert.Equal(t, 120000.0, updated.CreditLimit)
	assert.Equal(t, 5000.0, updated.UsedCredit)
}

func TestReconcileCreditCard_BothFieldsChanged_OnlyOneCashFlow(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	card := seedCreditCard(t, env, 5000, 100000)
	newUsed := 8000.0
	newLimit := 120000.0
	res, err := env.svc.ReconcileCreditCard(card.ID, &models.ReconcileCreditCardInput{
		NewUsedCredit:  &newUsed,
		NewCreditLimit: &newLimit,
		Date:           time.Now(),
	})
	require.NoError(t, err)
	require.NotNil(t, res.CashFlowID)
	assert.Equal(t, 3000.0, res.Delta)

	updated, err := env.cardRepo.GetByID(card.ID)
	require.NoError(t, err)
	assert.Equal(t, 120000.0, updated.CreditLimit)
	assert.Equal(t, 8000.0, updated.UsedCredit)
}

func TestReconcileCreditCard_UsedCreditExceedsLimit_Rejected(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	card := seedCreditCard(t, env, 5000, 100000)
	newUsed := 120000.0
	_, err := env.svc.ReconcileCreditCard(card.ID, &models.ReconcileCreditCardInput{
		NewUsedCredit: &newUsed,
		Date:          time.Now(),
	})
	assert.Error(t, err)

	updated, err := env.cardRepo.GetByID(card.ID)
	require.NoError(t, err)
	assert.Equal(t, 5000.0, updated.UsedCredit)
}

func TestReconcileBatch_PartialChanges(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	a, err := env.bankRepo.Create(&models.CreateBankAccountInput{
		BankName: "玉山", AccountType: "活存", AccountNumberLast4: "1111",
		Currency: models.CurrencyTWD, Balance: 10000,
	})
	require.NoError(t, err)
	b, err := env.bankRepo.Create(&models.CreateBankAccountInput{
		BankName: "國泰", AccountType: "活存", AccountNumberLast4: "2222",
		Currency: models.CurrencyTWD, Balance: 20000,
	})
	require.NoError(t, err)
	c := seedCreditCard(t, env, 5000, 50000)

	bal := 10000.0  // unchanged
	bal2 := 21000.0 // +1000
	used := 4500.0  // -500
	res, err := env.svc.ReconcileBatch(&models.ReconcileBatchInput{
		Date: time.Now(),
		Items: []models.ReconcileBatchItem{
			{TargetType: models.SourceTypeBankAccount, TargetID: a.ID, NewBalance: &bal},
			{TargetType: models.SourceTypeBankAccount, TargetID: b.ID, NewBalance: &bal2},
			{TargetType: models.SourceTypeCreditCard, TargetID: c.ID, NewUsedCredit: &used},
		},
	})
	require.NoError(t, err)
	require.Len(t, res, 3)

	assert.Nil(t, res[0].CashFlowID)
	require.NotNil(t, res[1].CashFlowID)
	cf1, err := env.cfRepo.GetByID(*res[1].CashFlowID)
	require.NoError(t, err)
	assert.Equal(t, models.CashFlowTypeIncome, cf1.Type)
	require.NotNil(t, res[2].CashFlowID)
	cf2, err := env.cfRepo.GetByID(*res[2].CashFlowID)
	require.NoError(t, err)
	assert.Equal(t, models.CashFlowTypeIncome, cf2.Type)
}

func TestReconcileBatch_OneItemFails_RollsBackAll(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	a := seedBankAccount(t, env, 10000)
	c := seedCreditCard(t, env, 5000, 50000)

	balOK := 12000.0
	tooMuch := 99999.0
	_, err := env.svc.ReconcileBatch(&models.ReconcileBatchInput{
		Date: time.Now(),
		Items: []models.ReconcileBatchItem{
			{TargetType: models.SourceTypeBankAccount, TargetID: a.ID, NewBalance: &balOK},
			{TargetType: models.SourceTypeCreditCard, TargetID: c.ID, NewUsedCredit: &tooMuch},
		},
	})
	require.Error(t, err)

	gotA, err := env.bankRepo.GetByID(a.ID)
	require.NoError(t, err)
	assert.Equal(t, 10000.0, gotA.Balance)
	gotC, err := env.cardRepo.GetByID(c.ID)
	require.NoError(t, err)
	assert.Equal(t, 5000.0, gotC.UsedCredit)
}

func TestReconcileBankAccount_NonTWD_CashFlowUsesAccountCurrency(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	acc, err := env.bankRepo.Create(&models.CreateBankAccountInput{
		BankName:           "玉山",
		AccountType:        "活存",
		AccountNumberLast4: "9999",
		Currency:           models.CurrencyUSD,
		Balance:            1000,
	})
	require.NoError(t, err)

	res, err := env.svc.ReconcileBankAccount(acc.ID, &models.ReconcileBankAccountInput{
		NewBalance: 1500,
		Date:       time.Now(),
	})
	require.NoError(t, err)
	require.NotNil(t, res.CashFlowID)

	cf, err := env.cfRepo.GetByID(*res.CashFlowID)
	require.NoError(t, err)
	assert.Equal(t, models.CurrencyUSD, cf.Currency)
}

func TestReconcileBankAccount_ConcurrentReconcile_CashFlowMatchesFinalDelta(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	startBalance := 1000.0
	acc := seedBankAccount(t, env, startBalance)

	var wg sync.WaitGroup
	wg.Add(2)
	errs := make([]error, 2)
	targets := []float64{1500, 1200}
	for i, target := range targets {
		i, target := i, target
		go func() {
			defer wg.Done()
			_, errs[i] = env.svc.ReconcileBankAccount(acc.ID, &models.ReconcileBankAccountInput{
				NewBalance: target,
				Date:       time.Now(),
			})
		}()
	}
	wg.Wait()
	for _, err := range errs {
		require.NoError(t, err)
	}

	updated, err := env.bankRepo.GetByID(acc.ID)
	require.NoError(t, err)

	var totalDelta float64
	rows, err := env.db.Query(
		`SELECT type, amount FROM cash_flows
		 WHERE source_type = $1 AND source_id = $2`,
		models.SourceTypeBankAccount, acc.ID,
	)
	require.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		var typ string
		var amt float64
		require.NoError(t, rows.Scan(&typ, &amt))
		if typ == string(models.CashFlowTypeIncome) {
			totalDelta += amt
		} else {
			totalDelta -= amt
		}
	}
	require.NoError(t, rows.Err())

	assert.InDelta(t, updated.Balance-startBalance, totalDelta, 0.001,
		"sum of cash_flow deltas must equal final balance - starting balance")
}

func TestReconcileBankAccount_NegativeBalance(t *testing.T) {
	env := newReconciliationTestEnv(t)
	defer env.db.Close()

	acc := seedBankAccount(t, env, 30000)
	_, err := env.svc.ReconcileBankAccount(acc.ID, &models.ReconcileBankAccountInput{
		NewBalance: -100,
		Date:       time.Now(),
	})
	assert.Error(t, err)

	updated, err := env.bankRepo.GetByID(acc.ID)
	require.NoError(t, err)
	assert.Equal(t, 30000.0, updated.Balance)
}
