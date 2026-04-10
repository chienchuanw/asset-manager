package service

import (
	"database/sql"
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// TransactionService 交易記錄業務邏輯介面
type TransactionService interface {
	CreateTransaction(input *models.CreateTransactionInput) (*models.Transaction, error)
	CreateTransactionsBatch(inputs []*models.CreateTransactionInput) ([]*models.Transaction, error)
	GetTransaction(id uuid.UUID) (*models.Transaction, error)
	ListTransactions(filters repository.TransactionFilters) ([]*models.Transaction, error)
	UpdateTransaction(id uuid.UUID, input *models.UpdateTransactionInput) (*models.Transaction, error)
	DeleteTransaction(id uuid.UUID) error
}

// transactionService 交易記錄業務邏輯實作
type transactionService struct {
	repo               repository.TransactionRepository
	realizedProfitRepo repository.RealizedProfitRepository
	fifoCalculator     FIFOCalculator
	exchangeRateService ExchangeRateService
}

// NewTransactionService 建立新的交易記錄 service
func NewTransactionService(
	repo repository.TransactionRepository,
	realizedProfitRepo repository.RealizedProfitRepository,
	fifoCalculator FIFOCalculator,
	exchangeRateService ExchangeRateService,
) TransactionService {
	return &transactionService{
		repo:               repo,
		realizedProfitRepo: realizedProfitRepo,
		fifoCalculator:     fifoCalculator,
		exchangeRateService: exchangeRateService,
	}
}

// CreateTransaction 建立新的交易記錄
func (s *transactionService) CreateTransaction(input *models.CreateTransactionInput) (*models.Transaction, error) {
	// 驗證資產類型
	if !input.AssetType.Validate() {
		return nil, fmt.Errorf("invalid asset type: %s", input.AssetType)
	}

	// 驗證交易類型
	if !input.TransactionType.Validate() {
		return nil, fmt.Errorf("invalid transaction type: %s", input.TransactionType)
	}

	// 驗證數量和價格
	if input.Quantity < 0 {
		return nil, fmt.Errorf("quantity must be non-negative")
	}

	if input.Price < 0 {
		return nil, fmt.Errorf("price must be non-negative")
	}

	// 驗證手續費
	if input.Fee != nil && *input.Fee < 0 {
		return nil, fmt.Errorf("fee must be non-negative")
	}

	// 驗證交易稅
	if input.Tax != nil && *input.Tax < 0 {
		return nil, fmt.Errorf("tax must be non-negative")
	}

	// 驗證幣別
	if !input.Currency.Validate() {
		return nil, fmt.Errorf("invalid currency: %s", input.Currency)
	}

	// 賣出交易需要在資料庫事務中同時建立交易和已實現損益
	if input.TransactionType == models.TransactionTypeSell {
		return s.createSellTransactionAtomically(input)
	}

	// 非賣出交易（買入/股息/手續費），不需要事務
	return s.createNonSellTransaction(input)
}

// CreateTransactionsBatch 批次建立交易記錄（全有或全無）
func (s *transactionService) CreateTransactionsBatch(inputs []*models.CreateTransactionInput) ([]*models.Transaction, error) {
	// 驗證輸入
	if len(inputs) == 0 {
		return nil, fmt.Errorf("no transactions to create")
	}

	// 驗證每筆交易的資料
	for i, input := range inputs {
		if !input.AssetType.Validate() {
			return nil, fmt.Errorf("transaction %d: invalid asset type: %s", i, input.AssetType)
		}
		if !input.TransactionType.Validate() {
			return nil, fmt.Errorf("transaction %d: invalid transaction type: %s", i, input.TransactionType)
		}
		if input.Quantity < 0 {
			return nil, fmt.Errorf("transaction %d: quantity must be non-negative", i)
		}
		if input.Price < 0 {
			return nil, fmt.Errorf("transaction %d: price must be non-negative", i)
		}
		if input.Fee != nil && *input.Fee < 0 {
			return nil, fmt.Errorf("transaction %d: fee must be non-negative", i)
		}
		if input.Tax != nil && *input.Tax < 0 {
			return nil, fmt.Errorf("transaction %d: tax must be non-negative", i)
		}
	}

	// 建立交易記錄陣列
	transactions := make([]*models.Transaction, 0, len(inputs))

	// 逐筆建立交易（使用現有的 CreateTransaction 方法）
	// 如果任一筆失敗，返回錯誤（呼叫方需要處理回滾）
	for i, input := range inputs {
		transaction, err := s.CreateTransaction(input)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction %d: %w", i, err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// GetTransaction 取得單筆交易記錄
func (s *transactionService) GetTransaction(id uuid.UUID) (*models.Transaction, error) {
	return s.repo.GetByID(id)
}

// ListTransactions 取得交易記錄列表
func (s *transactionService) ListTransactions(filters repository.TransactionFilters) ([]*models.Transaction, error) {
	// 驗證篩選條件
	if filters.AssetType != nil && !filters.AssetType.Validate() {
		return nil, fmt.Errorf("invalid asset type filter: %s", *filters.AssetType)
	}

	if filters.TransactionType != nil && !filters.TransactionType.Validate() {
		return nil, fmt.Errorf("invalid transaction type filter: %s", *filters.TransactionType)
	}

	return s.repo.GetAll(filters)
}

// UpdateTransaction 更新交易記錄
func (s *transactionService) UpdateTransaction(id uuid.UUID, input *models.UpdateTransactionInput) (*models.Transaction, error) {
	// 驗證資產類型
	if input.AssetType != nil && !input.AssetType.Validate() {
		return nil, fmt.Errorf("invalid asset type: %s", *input.AssetType)
	}

	// 驗證交易類型
	if input.TransactionType != nil && !input.TransactionType.Validate() {
		return nil, fmt.Errorf("invalid transaction type: %s", *input.TransactionType)
	}

	// 驗證數量和價格
	if input.Quantity != nil && *input.Quantity < 0 {
		return nil, fmt.Errorf("quantity must be non-negative")
	}

	if input.Price != nil && *input.Price < 0 {
		return nil, fmt.Errorf("price must be non-negative")
	}

	// 驗證手續費
	if input.Fee != nil && *input.Fee < 0 {
		return nil, fmt.Errorf("fee must be non-negative")
	}

	return s.repo.Update(id, input)
}

// DeleteTransaction 刪除交易記錄
func (s *transactionService) DeleteTransaction(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// createNonSellTransaction 建立非賣出交易（買入/股息/手續費）
func (s *transactionService) createNonSellTransaction(input *models.CreateTransactionInput) (*models.Transaction, error) {
	if input.Currency == models.CurrencyUSD {
		rate, err := s.exchangeRateService.GetRate(models.CurrencyUSD, models.CurrencyTWD, input.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to get exchange rate for USD transaction: %w", err)
		}

		exchangeRate, err := s.exchangeRateService.GetRateRecord(models.CurrencyUSD, models.CurrencyTWD, input.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to get exchange rate record: %w", err)
		}

		transaction, err := s.repo.CreateWithExchangeRate(input, exchangeRate.ID)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Created USD transaction with exchange rate %.4f (ID: %d)\n", rate, exchangeRate.ID)
		return transaction, nil
	}

	return s.repo.Create(input)
}

// createSellTransactionAtomically 在資料庫事務中建立賣出交易和已實現損益
func (s *transactionService) createSellTransactionAtomically(input *models.CreateTransactionInput) (*models.Transaction, error) {
	dbTx, err := s.repo.DB().Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer dbTx.Rollback()

	// 在事務中建立交易記錄
	var transaction *models.Transaction
	if input.Currency == models.CurrencyUSD {
		rate, err := s.exchangeRateService.GetRate(models.CurrencyUSD, models.CurrencyTWD, input.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to get exchange rate for USD transaction: %w", err)
		}

		exchangeRate, err := s.exchangeRateService.GetRateRecord(models.CurrencyUSD, models.CurrencyTWD, input.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to get exchange rate record: %w", err)
		}

		transaction, err = s.repo.CreateWithExchangeRateTx(dbTx, input, exchangeRate.ID)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Created USD sell transaction with exchange rate %.4f (ID: %d)\n", rate, exchangeRate.ID)
	} else {
		transaction, err = s.repo.CreateTx(dbTx, input)
		if err != nil {
			return nil, err
		}
	}

	// 在同一事務中建立已實現損益
	if err := s.createRealizedProfitTx(dbTx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create realized profit: %w", err)
	}

	if err := dbTx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// createRealizedProfitTx 在指定的資料庫事務中建立已實現損益記錄
func (s *transactionService) createRealizedProfitTx(dbTx *sql.Tx, sellTransaction *models.Transaction) error {
	// 取得該標的的所有交易記錄
	filters := repository.TransactionFilters{
		Symbol: &sellTransaction.Symbol,
	}
	allTransactions, err := s.repo.GetAll(filters)
	if err != nil {
		return fmt.Errorf("failed to get transactions for symbol %s: %w", sellTransaction.Symbol, err)
	}

	// 使用 FIFO Calculator 計算成本基礎
	costBasis, err := s.fifoCalculator.CalculateCostBasis(
		sellTransaction.Symbol,
		sellTransaction,
		allTransactions,
	)
	if err != nil {
		return fmt.Errorf("failed to calculate cost basis: %w", err)
	}

	// 準備賣出手續費
	sellFee := 0.0
	if sellTransaction.Fee != nil {
		sellFee = *sellTransaction.Fee
	}

	// 建立已實現損益記錄
	input := &models.CreateRealizedProfitInput{
		TransactionID: sellTransaction.ID.String(),
		Symbol:        sellTransaction.Symbol,
		AssetType:     sellTransaction.AssetType,
		SellDate:      sellTransaction.Date,
		Quantity:      sellTransaction.Quantity,
		SellPrice:     sellTransaction.Price,
		SellAmount:    sellTransaction.Amount,
		SellFee:       sellFee,
		CostBasis:     costBasis,
		Currency:      string(sellTransaction.Currency),
	}

	_, err = s.realizedProfitRepo.CreateTx(dbTx, input)
	if err != nil {
		return fmt.Errorf("failed to create realized profit record: %w", err)
	}

	return nil
}

