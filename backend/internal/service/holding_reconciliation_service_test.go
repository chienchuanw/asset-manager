package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type holdingReconcileTestEnv struct {
	db     *sql.DB
	svc    HoldingReconciliationService
	txRepo repository.TransactionRepository
}

// stubExchangeRateForReconcile 對 USD 一律以 1:1 回傳，用於不關心匯率的對帳測試。
type stubExchangeRateForReconcile struct{}

func (stubExchangeRateForReconcile) GetRate(_, _ models.Currency, _ time.Time) (float64, error) {
	return 1.0, nil
}
func (stubExchangeRateForReconcile) GetRateRecord(_, _ models.Currency, _ time.Time) (*models.ExchangeRate, error) {
	return nil, nil
}
func (stubExchangeRateForReconcile) GetTodayRate(_, _ models.Currency) (float64, error) {
	return 1.0, nil
}
func (stubExchangeRateForReconcile) RefreshTodayRate() error { return nil }
func (stubExchangeRateForReconcile) ConvertToTWD(amount float64, _ models.Currency, _ time.Time) (float64, error) {
	return amount, nil
}

func newHoldingReconcileTestEnv(t *testing.T) *holdingReconcileTestEnv {
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

	for _, tbl := range []string{"transactions", "cash_flows", "realized_profits"} {
		_, _ = db.Exec("DELETE FROM " + tbl)
	}

	txRepo := repository.NewTransactionRepository(db)
	fifo := NewFIFOCalculator(stubExchangeRateForReconcile{})
	return &holdingReconcileTestEnv{
		db:     db,
		txRepo: txRepo,
		svc:    NewHoldingReconciliationService(db, txRepo, fifo),
	}
}

func (env *holdingReconcileTestEnv) seedBuy(t *testing.T, assetType models.AssetType, symbol, name string, currency models.Currency, qty, price float64, date time.Time) {
	t.Helper()
	_, err := env.txRepo.Create(&models.CreateTransactionInput{
		Date:            date,
		AssetType:       assetType,
		Symbol:          symbol,
		Name:            name,
		TransactionType: models.TransactionTypeBuy,
		Quantity:        qty,
		Price:           price,
		Amount:          qty * price,
		Currency:        currency,
	})
	require.NoError(t, err)
}

func (env *holdingReconcileTestEnv) countTxByType(t *testing.T, txType models.TransactionType) int {
	t.Helper()
	var n int
	require.NoError(t, env.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE transaction_type = $1", txType).Scan(&n))
	return n
}

func (env *holdingReconcileTestEnv) countCashFlows(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, env.db.QueryRow("SELECT COUNT(*) FROM cash_flows").Scan(&n))
	return n
}

func (env *holdingReconcileTestEnv) countRealizedProfits(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, env.db.QueryRow("SELECT COUNT(*) FROM realized_profits").Scan(&n))
	return n
}

func TestHoldingReconciliation_BatchMixedActions(t *testing.T) {
	env := newHoldingReconcileTestEnv(t)
	defer env.db.Close()

	earlier := time.Now().Add(-30 * 24 * time.Hour)
	env.seedBuy(t, models.AssetTypeUSStock, "AAPL", "Apple Inc.", models.CurrencyUSD, 100, 150, earlier)
	env.seedBuy(t, models.AssetTypeUSStock, "TSLA", "Tesla", models.CurrencyUSD, 50, 250, earlier)
	env.seedBuy(t, models.AssetTypeTWStock, "0050", "元大台灣 50", models.CurrencyTWD, 200, 110, earlier)

	items := []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD, TargetQuantity: 150, TargetAvgCost: 180},
		{AssetType: models.AssetTypeUSStock, Symbol: "TSLA", Currency: models.CurrencyUSD, TargetQuantity: 0, TargetAvgCost: 0},
		{AssetType: models.AssetTypeUSStock, Symbol: "NVDA", Name: "NVIDIA", Currency: models.CurrencyUSD, TargetQuantity: 25, TargetAvgCost: 800},
		{AssetType: models.AssetTypeTWStock, Symbol: "0050", Currency: models.CurrencyTWD, TargetQuantity: 200, TargetAvgCost: 110},
	}

	preview, err := env.svc.ReconcileHoldings(context.Background(), items, false)
	require.NoError(t, err)
	require.Len(t, preview.Items, 4)

	bySymbol := map[string]models.ReconcileAction{}
	for _, p := range preview.Items {
		bySymbol[p.Item.Symbol] = p.Action
	}
	assert.Equal(t, models.ReconcileActionUpdate, bySymbol["AAPL"])
	assert.Equal(t, models.ReconcileActionLiquidate, bySymbol["TSLA"])
	assert.Equal(t, models.ReconcileActionCreate, bySymbol["NVDA"])
	assert.Equal(t, models.ReconcileActionNoop, bySymbol["0050"])

	// 三筆 adjustment 寫入（noop 不寫）
	assert.Equal(t, 3, env.countTxByType(t, models.TransactionTypeAdjustment))
	// cash flow 與 realized P&L 完全未受影響
	assert.Equal(t, 0, env.countCashFlows(t))
	assert.Equal(t, 0, env.countRealizedProfits(t))
}

func TestHoldingReconciliation_DryRunWritesNothing(t *testing.T) {
	env := newHoldingReconcileTestEnv(t)
	defer env.db.Close()
	env.seedBuy(t, models.AssetTypeUSStock, "AAPL", "Apple Inc.", models.CurrencyUSD, 100, 150, time.Now().Add(-24*time.Hour))

	preview, err := env.svc.ReconcileHoldings(context.Background(), []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD, TargetQuantity: 150, TargetAvgCost: 180},
	}, true)
	require.NoError(t, err)
	assert.Len(t, preview.Items, 1)
	assert.Equal(t, models.ReconcileActionUpdate, preview.Items[0].Action)
	assert.Equal(t, 0, env.countTxByType(t, models.TransactionTypeAdjustment))
}

func TestHoldingReconciliation_NameRequiredForCreate(t *testing.T) {
	env := newHoldingReconcileTestEnv(t)
	defer env.db.Close()

	_, err := env.svc.ReconcileHoldings(context.Background(), []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: "NVDA", Currency: models.CurrencyUSD, TargetQuantity: 10, TargetAvgCost: 800},
	}, false)
	require.ErrorIs(t, err, models.ErrReconcileNameRequired)
	assert.Equal(t, 0, env.countTxByType(t, models.TransactionTypeAdjustment))
}

func TestHoldingReconciliation_CurrencyMismatch(t *testing.T) {
	env := newHoldingReconcileTestEnv(t)
	defer env.db.Close()
	env.seedBuy(t, models.AssetTypeUSStock, "AAPL", "Apple Inc.", models.CurrencyUSD, 100, 150, time.Now().Add(-24*time.Hour))

	_, err := env.svc.ReconcileHoldings(context.Background(), []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyTWD, TargetQuantity: 150, TargetAvgCost: 180},
	}, false)
	require.ErrorIs(t, err, models.ErrReconcileCurrencyMismatch)
	assert.Equal(t, 0, env.countTxByType(t, models.TransactionTypeAdjustment))
}

func TestHoldingReconciliation_EmptyBatchRejected(t *testing.T) {
	env := newHoldingReconcileTestEnv(t)
	defer env.db.Close()

	_, err := env.svc.ReconcileHoldings(context.Background(), nil, false)
	require.ErrorIs(t, err, models.ErrReconcileEmptyBatch)
}

func TestHoldingReconciliation_ValidationFailureRollsBackEntireBatch(t *testing.T) {
	env := newHoldingReconcileTestEnv(t)
	defer env.db.Close()
	env.seedBuy(t, models.AssetTypeUSStock, "AAPL", "Apple Inc.", models.CurrencyUSD, 100, 150, time.Now().Add(-24*time.Hour))

	_, err := env.svc.ReconcileHoldings(context.Background(), []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD, TargetQuantity: 150, TargetAvgCost: 180},
		{AssetType: models.AssetTypeUSStock, Symbol: "AMD", Currency: models.CurrencyUSD, TargetQuantity: -1, TargetAvgCost: 100},
	}, false)
	require.ErrorIs(t, err, models.ErrInvalidReconcileQuantity)
	// 即使第一筆合法，整批應 rollback
	assert.Equal(t, 0, env.countTxByType(t, models.TransactionTypeAdjustment))
}

func TestHoldingReconciliation_AuditFieldsCaptured(t *testing.T) {
	env := newHoldingReconcileTestEnv(t)
	defer env.db.Close()
	env.seedBuy(t, models.AssetTypeUSStock, "AAPL", "Apple Inc.", models.CurrencyUSD, 100, 150, time.Now().Add(-24*time.Hour))

	_, err := env.svc.ReconcileHoldings(context.Background(), []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD, TargetQuantity: 150, TargetAvgCost: 180, Reason: "broker statement"},
	}, false)
	require.NoError(t, err)

	var prevQty, prevAvg sql.NullFloat64
	var reason sql.NullString
	require.NoError(t, env.db.QueryRow(`
		SELECT adjustment_prev_quantity, adjustment_prev_avg_cost, adjustment_reason
		FROM transactions
		WHERE transaction_type = 'adjustment' AND symbol = 'AAPL'
	`).Scan(&prevQty, &prevAvg, &reason))
	assert.True(t, prevQty.Valid)
	assert.InDelta(t, 100.0, prevQty.Float64, 1e-9)
	assert.True(t, prevAvg.Valid)
	assert.InDelta(t, 150.0, prevAvg.Float64, 1e-6)
	assert.True(t, reason.Valid)
	assert.Equal(t, "broker statement", reason.String)
}

// 確保編譯依賴 os 套件（避免 lint 警告）
var _ = os.Getenv
