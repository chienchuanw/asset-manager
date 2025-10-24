package service

import (
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// TransactionService 交易記錄業務邏輯介面
type TransactionService interface {
	CreateTransaction(input *models.CreateTransactionInput) (*models.Transaction, error)
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
}

// NewTransactionService 建立新的交易記錄 service
func NewTransactionService(
	repo repository.TransactionRepository,
	realizedProfitRepo repository.RealizedProfitRepository,
	fifoCalculator FIFOCalculator,
) TransactionService {
	return &transactionService{
		repo:               repo,
		realizedProfitRepo: realizedProfitRepo,
		fifoCalculator:     fifoCalculator,
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

	// 呼叫 repository 建立交易記錄
	transaction, err := s.repo.Create(input)
	if err != nil {
		return nil, err
	}

	// 如果是賣出交易，自動計算並記錄已實現損益
	if input.TransactionType == models.TransactionTypeSell {
		if err := s.createRealizedProfit(transaction); err != nil {
			// 記錄錯誤但不影響交易建立
			// TODO: 考慮是否要回滾交易或使用事務
			fmt.Printf("Warning: failed to create realized profit for transaction %s: %v\n", transaction.ID, err)
		}
	}

	return transaction, nil
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

// createRealizedProfit 建立已實現損益記錄（賣出交易時自動呼叫）
func (s *transactionService) createRealizedProfit(sellTransaction *models.Transaction) error {
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

	_, err = s.realizedProfitRepo.Create(input)
	if err != nil {
		return fmt.Errorf("failed to create realized profit record: %w", err)
	}

	return nil
}

