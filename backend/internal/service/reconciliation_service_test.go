package service

import (
	"database/sql"
	"fmt"
	"os"
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
