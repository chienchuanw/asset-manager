package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionRepository 模擬的 repository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(input *models.CreateTransactionInput) (*models.Transaction, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetAll(filters repository.TransactionFilters) ([]*models.Transaction, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(id uuid.UUID, input *models.UpdateTransactionInput) (*models.Transaction, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestCreateTransaction_Success 測試成功建立交易記錄
func TestCreateTransaction_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	fee := 28.0
	input := &models.CreateTransactionInput{
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        10,
		Price:           620,
		Amount:          6200,
		Fee:             &fee,
	}

	expectedTransaction := &models.Transaction{
		ID:              uuid.New(),
		Date:            input.Date,
		AssetType:       input.AssetType,
		Symbol:          input.Symbol,
		Name:            input.Name,
		TransactionType: input.TransactionType,
		Quantity:        input.Quantity,
		Price:           input.Price,
		Amount:          input.Amount,
		Fee:             input.Fee,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockRepo.On("Create", input).Return(expectedTransaction, nil)

	// Act
	result, err := service.CreateTransaction(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransaction.ID, result.ID)
	assert.Equal(t, expectedTransaction.Symbol, result.Symbol)
	mockRepo.AssertExpectations(t)
}

// TestCreateTransaction_InvalidAssetType 測試無效的資產類型
func TestCreateTransaction_InvalidAssetType(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	input := &models.CreateTransactionInput{
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetType("invalid"),
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        10,
		Price:           620,
		Amount:          6200,
	}

	// Act
	result, err := service.CreateTransaction(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid asset type")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestCreateTransaction_InvalidTransactionType 測試無效的交易類型
func TestCreateTransaction_InvalidTransactionType(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	input := &models.CreateTransactionInput{
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionType("invalid"),
		Quantity:        10,
		Price:           620,
		Amount:          6200,
	}

	// Act
	result, err := service.CreateTransaction(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid transaction type")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestCreateTransaction_NegativeQuantity 測試負數數量
func TestCreateTransaction_NegativeQuantity(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	input := &models.CreateTransactionInput{
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        -10,
		Price:           620,
		Amount:          6200,
	}

	// Act
	result, err := service.CreateTransaction(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "quantity must be non-negative")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestGetTransaction_Success 測試成功取得交易記錄
func TestGetTransaction_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	transactionID := uuid.New()
	expectedTransaction := &models.Transaction{
		ID:              transactionID,
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        10,
		Price:           620,
		Amount:          6200,
	}

	mockRepo.On("GetByID", transactionID).Return(expectedTransaction, nil)

	// Act
	result, err := service.GetTransaction(transactionID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransaction.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

// TestGetTransaction_NotFound 測試取得不存在的交易記錄
func TestGetTransaction_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	transactionID := uuid.New()
	mockRepo.On("GetByID", transactionID).Return(nil, fmt.Errorf("transaction not found"))

	// Act
	result, err := service.GetTransaction(transactionID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// TestListTransactions_Success 測試成功取得交易記錄列表
func TestListTransactions_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	filters := repository.TransactionFilters{}
	expectedTransactions := []*models.Transaction{
		{
			ID:              uuid.New(),
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
		},
		{
			ID:              uuid.New(),
			Symbol:          "ETH",
			Name:            "Ethereum",
			TransactionType: models.TransactionTypeBuy,
		},
	}

	mockRepo.On("GetAll", filters).Return(expectedTransactions, nil)

	// Act
	result, err := service.ListTransactions(filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

// TestDeleteTransaction_Success 測試成功刪除交易記錄
func TestDeleteTransaction_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	transactionID := uuid.New()
	mockRepo.On("Delete", transactionID).Return(nil)

	// Act
	err := service.DeleteTransaction(transactionID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

