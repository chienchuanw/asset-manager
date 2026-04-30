package features_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/cucumber/godog"
	_ "github.com/lib/pq"
)

type stubExchangeRate struct{}

func (stubExchangeRate) GetRate(_, _ models.Currency, _ time.Time) (float64, error) { return 1.0, nil }
func (stubExchangeRate) GetRateRecord(_, _ models.Currency, _ time.Time) (*models.ExchangeRate, error) {
	return nil, nil
}
func (stubExchangeRate) GetTodayRate(_, _ models.Currency) (float64, error) { return 1.0, nil }
func (stubExchangeRate) RefreshTodayRate() error                            { return nil }
func (stubExchangeRate) ConvertToTWD(amount float64, _ models.Currency, _ time.Time) (float64, error) {
	return amount, nil
}

type scenarioState struct {
	db      *sql.DB
	txRepo  repository.TransactionRepository
	fifo    service.FIFOCalculator
	svc     service.HoldingReconciliationService
	lastErr error
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func openTestDB() (*sql.DB, error) {
	host := envOrDefault("TEST_DB_HOST", "localhost")
	port := envOrDefault("TEST_DB_PORT", "5432")
	user := envOrDefault("TEST_DB_USER", "postgres")
	pwd := envOrDefault("TEST_DB_PASSWORD", "postgres")
	name := envOrDefault("TEST_DB_NAME", "asset_manager_test")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pwd, name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (s *scenarioState) reset() error {
	if s.db == nil {
		db, err := openTestDB()
		if err != nil {
			return err
		}
		s.db = db
		s.txRepo = repository.NewTransactionRepository(db)
		s.fifo = service.NewFIFOCalculator(stubExchangeRate{})
		s.svc = service.NewHoldingReconciliationService(db, s.txRepo, s.fifo)
	}
	for _, tbl := range []string{"transactions", "cash_flows", "realized_profits"} {
		if _, err := s.db.Exec("DELETE FROM " + tbl); err != nil {
			return fmt.Errorf("clean %s: %w", tbl, err)
		}
	}
	s.lastErr = nil
	return nil
}

func (s *scenarioState) givenHolding(symbol string, qty, price float64) error {
	_, err := s.txRepo.Create(&models.CreateTransactionInput{
		Date:            time.Now().Add(-30 * 24 * time.Hour),
		AssetType:       models.AssetTypeUSStock,
		Symbol:          symbol,
		Name:            symbol,
		TransactionType: models.TransactionTypeBuy,
		Quantity:        qty,
		Price:           price,
		Amount:          qty * price,
		Currency:        models.CurrencyUSD,
	})
	return err
}

func (s *scenarioState) givenNoHolding(symbol string) error {
	_, err := s.db.Exec("DELETE FROM transactions WHERE symbol = $1", symbol)
	return err
}

func (s *scenarioState) reconcile(symbol string, qty, avgCost float64) error {
	_, err := s.svc.ReconcileHoldings(context.Background(), []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: symbol, Currency: models.CurrencyUSD, TargetQuantity: qty, TargetAvgCost: avgCost},
	}, false)
	s.lastErr = err
	return nil
}

func (s *scenarioState) reconcileLiquidate(symbol string) error {
	return s.reconcile(symbol, 0, 0)
}

func (s *scenarioState) reconcileNew(symbol, name string, qty, avgCost float64) error {
	_, err := s.svc.ReconcileHoldings(context.Background(), []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: symbol, Name: name, Currency: models.CurrencyUSD, TargetQuantity: qty, TargetAvgCost: avgCost},
	}, false)
	s.lastErr = err
	return nil
}

func (s *scenarioState) submitInvalidBatch() error {
	_, err := s.svc.ReconcileHoldings(context.Background(), []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD, TargetQuantity: 200, TargetAvgCost: 180},
		{AssetType: models.AssetTypeUSStock, Symbol: "BAD", Currency: models.CurrencyUSD, TargetQuantity: -1, TargetAvgCost: 100},
	}, false)
	s.lastErr = err
	return nil
}

func (s *scenarioState) sellAtPrice(symbol string, qty, price float64) error {
	// 故意拉到一小時後，確保 sell 嚴格晚於 adjustment（避免 sort.Slice 在同日期時不穩定）
	_, err := s.txRepo.Create(&models.CreateTransactionInput{
		Date:            time.Now().Add(time.Hour),
		AssetType:       models.AssetTypeUSStock,
		Symbol:          symbol,
		Name:            symbol,
		TransactionType: models.TransactionTypeSell,
		Quantity:        qty,
		Price:           price,
		Amount:          qty * price,
		Currency:        models.CurrencyUSD,
	})
	return err
}

func (s *scenarioState) assertSingleAdjustment() error {
	var n int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE transaction_type = 'adjustment'").Scan(&n); err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("expected exactly 1 adjustment, got %d", n)
	}
	return nil
}

func (s *scenarioState) assertAuditFields(prevQty, prevAvg float64) error {
	var pq, pa sql.NullFloat64
	if err := s.db.QueryRow(`
		SELECT adjustment_prev_quantity, adjustment_prev_avg_cost
		FROM transactions
		WHERE transaction_type = 'adjustment'
		LIMIT 1
	`).Scan(&pq, &pa); err != nil {
		return err
	}
	if !pq.Valid || pq.Float64 != prevQty {
		return fmt.Errorf("prev_quantity expected %v, got %v", prevQty, pq)
	}
	if !pa.Valid || pa.Float64 != prevAvg {
		return fmt.Errorf("prev_avg_cost expected %v, got %v", prevAvg, pa)
	}
	return nil
}

func (s *scenarioState) assertHoldingZero(symbol string) error {
	txs, err := allSymbolTxs(s.db, symbol)
	if err != nil {
		return err
	}
	h, err := s.fifo.CalculateHoldingForSymbol(symbol, txs)
	if err != nil {
		return err
	}
	if h != nil && h.Quantity > 0 {
		return fmt.Errorf("expected 0 holdings, got %v", h.Quantity)
	}
	return nil
}

func (s *scenarioState) assertHoldingExists(symbol string, qty, avgCost float64) error {
	txs, err := allSymbolTxs(s.db, symbol)
	if err != nil {
		return err
	}
	h, err := s.fifo.CalculateHoldingForSymbol(symbol, txs)
	if err != nil {
		return err
	}
	if h == nil {
		return fmt.Errorf("holding %s not found", symbol)
	}
	if h.Quantity != qty {
		return fmt.Errorf("quantity expected %v, got %v", qty, h.Quantity)
	}
	if h.AvgCostOriginal != avgCost {
		return fmt.Errorf("avg cost expected %v, got %v", avgCost, h.AvgCostOriginal)
	}
	return nil
}

func (s *scenarioState) assertNoRealizedProfit() error {
	var n int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM realized_profits").Scan(&n); err != nil {
		return err
	}
	if n != 0 {
		return fmt.Errorf("expected 0 realized profits, got %d", n)
	}
	return nil
}

func (s *scenarioState) assertNoCashFlow() error {
	var n int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM cash_flows").Scan(&n); err != nil {
		return err
	}
	if n != 0 {
		return fmt.Errorf("expected 0 cash flows, got %d", n)
	}
	return nil
}

func (s *scenarioState) assertBatchRejected() error {
	if s.lastErr == nil {
		return errors.New("expected batch error, got nil")
	}
	return nil
}

func (s *scenarioState) assertNoAdjustments() error {
	var n int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE transaction_type = 'adjustment'").Scan(&n); err != nil {
		return err
	}
	if n != 0 {
		return fmt.Errorf("expected 0 adjustments, got %d", n)
	}
	return nil
}

func (s *scenarioState) assertCostBasis(symbol string, expected float64) error {
	txs, err := allSymbolTxs(s.db, symbol)
	if err != nil {
		return err
	}
	var sellTx *models.Transaction
	for _, t := range txs {
		if t.TransactionType == models.TransactionTypeSell {
			sellTx = t
			break
		}
	}
	if sellTx == nil {
		return errors.New("no sell transaction found")
	}
	cb, err := s.fifo.CalculateCostBasis(symbol, sellTx, txs)
	if err != nil {
		return err
	}
	if cb != expected {
		return fmt.Errorf("cost basis expected %v, got %v", expected, cb)
	}
	return nil
}

func allSymbolTxs(db *sql.DB, symbol string) ([]*models.Transaction, error) {
	rows, err := db.Query(`
		SELECT id, date, asset_type, symbol, name, transaction_type, quantity, price, amount, fee, tax, currency, exchange_rate_id, note,
		       adjustment_prev_quantity, adjustment_prev_avg_cost, adjustment_reason, broker_account_id,
		       created_at, updated_at
		FROM transactions
		WHERE symbol = $1
		ORDER BY date ASC, created_at ASC
	`, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Transaction
	for rows.Next() {
		t := &models.Transaction{}
		if err := rows.Scan(
			&t.ID, &t.Date, &t.AssetType, &t.Symbol, &t.Name, &t.TransactionType,
			&t.Quantity, &t.Price, &t.Amount, &t.Fee, &t.Tax, &t.Currency, &t.ExchangeRateID, &t.Note,
			&t.AdjustmentPrevQuantity, &t.AdjustmentPrevAvgCost, &t.AdjustmentReason, &t.BrokerAccountID,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func InitializeScenario(sc *godog.ScenarioContext) {
	state := &scenarioState{}

	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return ctx, state.reset()
	})

	sc.Step(`^一個乾淨的測試資料庫$`, func() error { return nil })

	sc.Step(`^目前持有 "([^"]+)" (\d+(?:\.\d+)?) 股，平均成本 (\d+(?:\.\d+)?) USD$`, state.givenHolding)
	sc.Step(`^我目前未持有 "([^"]+)"$`, state.givenNoHolding)

	sc.Step(`^將 "([^"]+)" 對帳為 (\d+(?:\.\d+)?) 股，平均成本 (\d+(?:\.\d+)?) USD$`, state.reconcile)
	sc.Step(`^將 "([^"]+)" 對帳為 0 股$`, state.reconcileLiquidate)
	sc.Step(`^將新標的 "([^"]+)"（名稱 "([^"]+)"）對帳為 (\d+(?:\.\d+)?) 股，平均成本 (\d+(?:\.\d+)?) USD$`, state.reconcileNew)
	sc.Step(`^提交一個批次，包含一筆合法更新與一筆 quantity 為負的非法項目$`, state.submitInvalidBatch)
	sc.Step(`^賣出 "([^"]+)" (\d+(?:\.\d+)?) 股，每股 (\d+(?:\.\d+)?) USD$`, state.sellAtPrice)

	sc.Step(`^系統會新增一筆 adjustment 交易$`, state.assertSingleAdjustment)
	sc.Step(`^該 adjustment 的稽核欄位記錄 prev_quantity 為 (\d+(?:\.\d+)?)，prev_avg_cost 為 (\d+(?:\.\d+)?)$`, state.assertAuditFields)
	sc.Step(`^"([^"]+)" 的持倉數量為 0$`, state.assertHoldingZero)
	sc.Step(`^"([^"]+)" 的持倉存在，數量為 (\d+(?:\.\d+)?)，平均成本為 (\d+(?:\.\d+)?)$`, state.assertHoldingExists)
	sc.Step(`^沒有任何 realized profit 紀錄$`, state.assertNoRealizedProfit)
	sc.Step(`^沒有任何 cash flow 紀錄$`, state.assertNoCashFlow)
	sc.Step(`^整批請求被拒絕$`, state.assertBatchRejected)
	sc.Step(`^沒有任何 adjustment 交易被寫入$`, state.assertNoAdjustments)
	sc.Step(`^該賣出的成本基礎為 (\d+(?:\.\d+)?)$`, func(expected float64) error {
		return state.assertCostBasis("AAPL", expected)
	})

	sc.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return ctx, nil
	})
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "holding-reconciliation",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("godog test suite failed")
	}
}
