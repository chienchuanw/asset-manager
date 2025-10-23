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
	repo repository.TransactionRepository
}

// NewTransactionService 建立新的交易記錄 service
func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
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
	return s.repo.Create(input)
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

