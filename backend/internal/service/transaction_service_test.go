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

// MockTransactionRepository 模擬的 TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

// MockRealizedProfitRepository 模擬的 RealizedProfitRepository
type MockRealizedProfitRepository struct {
	mock.Mock
}

// MockFIFOCalculator 模擬的 FIFOCalculator
type MockFIFOCalculator struct {
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

// MockRealizedProfitRepository 方法實作
func (m *MockRealizedProfitRepository) Create(input *models.CreateRealizedProfitInput) (*models.RealizedProfit, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RealizedProfit), args.Error(1)
}

func (m *MockRealizedProfitRepository) GetByTransactionID(transactionID string) (*models.RealizedProfit, error) {
	args := m.Called(transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RealizedProfit), args.Error(1)
}

func (m *MockRealizedProfitRepository) GetAll(filters models.RealizedProfitFilters) ([]*models.RealizedProfit, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.RealizedProfit), args.Error(1)
}

func (m *MockRealizedProfitRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockFIFOCalculator 方法實作
func (m *MockFIFOCalculator) CalculateHoldingForSymbol(symbol string, transactions []*models.Transaction) (*models.Holding, error) {
	args := m.Called(symbol, transactions)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Holding), args.Error(1)
}

func (m *MockFIFOCalculator) CalculateAllHoldings(transactions []*models.Transaction) (map[string]*models.Holding, error) {
	args := m.Called(transactions)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]*models.Holding), args.Error(1)
}

func (m *MockFIFOCalculator) CalculateCostBasis(symbol string, sellTransaction *models.Transaction, allTransactions []*models.Transaction) (float64, error) {
	args := m.Called(symbol, sellTransaction, allTransactions)
	return args.Get(0).(float64), args.Error(1)
}

// TestCreateTransaction_Success 測試成功建立買入交易記錄
func TestCreateTransaction_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	mockRealizedProfitRepo := new(MockRealizedProfitRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	service := NewTransactionService(mockRepo, mockRealizedProfitRepo, mockFIFOCalc)

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
		Currency:        models.CurrencyTWD,
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
		Currency:        input.Currency,
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
	// 買入交易不應該呼叫 RealizedProfitRepo 或 FIFOCalculator
	mockRealizedProfitRepo.AssertNotCalled(t, "Create")
	mockFIFOCalc.AssertNotCalled(t, "CalculateCostBasis")
}

// TestCreateTransaction_InvalidAssetType 測試無效的資產類型
func TestCreateTransaction_InvalidAssetType(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	mockRealizedProfitRepo := new(MockRealizedProfitRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	service := NewTransactionService(mockRepo, mockRealizedProfitRepo, mockFIFOCalc)

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
	mockRealizedProfitRepo := new(MockRealizedProfitRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	service := NewTransactionService(mockRepo, mockRealizedProfitRepo, mockFIFOCalc)

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
	mockRealizedProfitRepo := new(MockRealizedProfitRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	service := NewTransactionService(mockRepo, mockRealizedProfitRepo, mockFIFOCalc)

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
	mockRealizedProfitRepo := new(MockRealizedProfitRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	service := NewTransactionService(mockRepo, mockRealizedProfitRepo, mockFIFOCalc)

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
	mockRealizedProfitRepo := new(MockRealizedProfitRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	service := NewTransactionService(mockRepo, mockRealizedProfitRepo, mockFIFOCalc)

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
	mockRealizedProfitRepo := new(MockRealizedProfitRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	service := NewTransactionService(mockRepo, mockRealizedProfitRepo, mockFIFOCalc)

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
	mockRealizedProfitRepo := new(MockRealizedProfitRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	service := NewTransactionService(mockRepo, mockRealizedProfitRepo, mockFIFOCalc)

	transactionID := uuid.New()
	mockRepo.On("Delete", transactionID).Return(nil)

	// Act
	err := service.DeleteTransaction(transactionID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreateTransaction_SellWithRealizedProfit 測試建立賣出交易並自動建立已實現損益
func TestCreateTransaction_SellWithRealizedProfit(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	mockRealizedProfitRepo := new(MockRealizedProfitRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	service := NewTransactionService(mockRepo, mockRealizedProfitRepo, mockFIFOCalc)

	fee := 28.0
	sellInput := &models.CreateTransactionInput{
		Date:            time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeSell,
		Quantity:        100,
		Price:           620,
		Amount:          62000,
		Fee:             &fee,
		Currency:        models.CurrencyTWD,
	}

	sellTransaction := &models.Transaction{
		ID:              uuid.New(),
		Date:            sellInput.Date,
		AssetType:       sellInput.AssetType,
		Symbol:          sellInput.Symbol,
		Name:            sellInput.Name,
		TransactionType: sellInput.TransactionType,
		Quantity:        sellInput.Quantity,
		Price:           sellInput.Price,
		Amount:          sellInput.Amount,
		Fee:             sellInput.Fee,
		Currency:        sellInput.Currency,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 模擬之前的買入交易
	previousTransactions := []*models.Transaction{
		{
			ID:              uuid.New(),
			Date:            time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
			Currency:        models.CurrencyTWD,
		},
	}

	// Mock 期望
	mockRepo.On("Create", sellInput).Return(sellTransaction, nil)

	filters := repository.TransactionFilters{Symbol: &sellInput.Symbol}
	mockRepo.On("GetAll", filters).Return(previousTransactions, nil)

	costBasis := 50028.0 // (50000 + 28)
	mockFIFOCalc.On("CalculateCostBasis", "2330", sellTransaction, previousTransactions).Return(costBasis, nil)

	mockRealizedProfitRepo.On("Create", mock.MatchedBy(func(input *models.CreateRealizedProfitInput) bool {
		return input.Symbol == "2330" &&
			input.Quantity == 100 &&
			input.SellAmount == 62000 &&
			input.SellFee == 28 &&
			input.CostBasis == 50028.0
	})).Return(&models.RealizedProfit{}, nil)

	// Act
	result, err := service.CreateTransaction(sellInput)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sellTransaction.ID, result.ID)
	mockRepo.AssertExpectations(t)
	mockFIFOCalc.AssertExpectations(t)
	mockRealizedProfitRepo.AssertExpectations(t)
}

