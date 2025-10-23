package service

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==================== Mock Objects ====================

// MockTransactionRepositoryForHolding Transaction Repository 的 Mock
type MockTransactionRepositoryForHolding struct {
	mock.Mock
}

func (m *MockTransactionRepositoryForHolding) Create(input *models.CreateTransactionInput) (*models.Transaction, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForHolding) GetByID(id uuid.UUID) (*models.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForHolding) GetAll(filters repository.TransactionFilters) ([]*models.Transaction, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForHolding) Update(id uuid.UUID, input *models.UpdateTransactionInput) (*models.Transaction, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForHolding) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockPriceService Price Service 的 Mock
type MockPriceService struct {
	mock.Mock
}

func (m *MockPriceService) GetPrice(symbol string, assetType models.AssetType) (*models.Price, error) {
	args := m.Called(symbol, assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Price), args.Error(1)
}

func (m *MockPriceService) GetPrices(symbols []string, assetTypes map[string]models.AssetType) (map[string]*models.Price, error) {
	args := m.Called(symbols, assetTypes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]*models.Price), args.Error(1)
}

func (m *MockPriceService) RefreshPrice(symbol string, assetType models.AssetType) (*models.Price, error) {
	args := m.Called(symbol, assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Price), args.Error(1)
}

// ==================== 測試案例 ====================

// TestGetAllHoldings_Success 測試成功取得所有持倉
func TestGetAllHoldings_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepositoryForHolding)
	mockPriceService := new(MockPriceService)
	fifoCalculator := NewFIFOCalculator()
	service := NewHoldingService(mockRepo, fifoCalculator, mockPriceService)

	// 準備交易記錄
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeUSStock,
			Symbol:          "AAPL",
			Name:            "Apple Inc.",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        50,
			Price:           150,
			Amount:          7500,
			Fee:             ptrFloat64(10),
		},
	}

	// 準備價格資料
	prices := map[string]*models.Price{
		"2330": {
			Symbol:    "2330",
			AssetType: models.AssetTypeTWStock,
			Price:     620,
			Currency:  "TWD",
			Source:    "mock",
			UpdatedAt: time.Now(),
		},
		"AAPL": {
			Symbol:    "AAPL",
			AssetType: models.AssetTypeUSStock,
			Price:     175,
			Currency:  "USD",
			Source:    "mock",
			UpdatedAt: time.Now(),
		},
	}

	// Mock 設定
	mockRepo.On("GetAll", mock.Anything).Return(transactions, nil)
	mockPriceService.On("GetPrices", mock.Anything, mock.Anything).Return(prices, nil)

	// Act
	holdings, err := service.GetAllHoldings(models.HoldingFilters{})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(holdings))

	// 驗證台積電持倉
	var tsmc *models.Holding
	for _, h := range holdings {
		if h.Symbol == "2330" {
			tsmc = h
			break
		}
	}
	assert.NotNil(t, tsmc)
	assert.Equal(t, 100.0, tsmc.Quantity)
	assert.InDelta(t, 500.28, tsmc.AvgCost, 0.01)
	assert.Equal(t, 620.0, tsmc.CurrentPrice)
	assert.InDelta(t, 62000.0, tsmc.MarketValue, 0.01)
	assert.InDelta(t, 11972.0, tsmc.UnrealizedPL, 0.01) // 62000 - 50028
	assert.InDelta(t, 23.93, tsmc.UnrealizedPLPct, 0.01) // (11972 / 50028) * 100

	mockRepo.AssertExpectations(t)
	mockPriceService.AssertExpectations(t)
}

// TestGetAllHoldings_EmptyTransactions 測試空交易記錄
func TestGetAllHoldings_EmptyTransactions(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepositoryForHolding)
	mockPriceService := new(MockPriceService)
	fifoCalculator := NewFIFOCalculator()
	service := NewHoldingService(mockRepo, fifoCalculator, mockPriceService)

	// Mock 設定
	mockRepo.On("GetAll", mock.Anything).Return([]*models.Transaction{}, nil)

	// Act
	holdings, err := service.GetAllHoldings(models.HoldingFilters{})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(holdings))

	mockRepo.AssertExpectations(t)
}

// TestGetAllHoldings_WithAssetTypeFilter 測試按資產類型篩選
func TestGetAllHoldings_WithAssetTypeFilter(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepositoryForHolding)
	mockPriceService := new(MockPriceService)
	fifoCalculator := NewFIFOCalculator()
	service := NewHoldingService(mockRepo, fifoCalculator, mockPriceService)

	assetType := models.AssetTypeTWStock
	filters := models.HoldingFilters{
		AssetType: &assetType,
	}

	// 準備交易記錄（只有台股）
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
	}

	prices := map[string]*models.Price{
		"2330": {
			Symbol:    "2330",
			AssetType: models.AssetTypeTWStock,
			Price:     620,
			Currency:  "TWD",
			Source:    "mock",
			UpdatedAt: time.Now(),
		},
	}

	// Mock 設定
	mockRepo.On("GetAll", mock.MatchedBy(func(f repository.TransactionFilters) bool {
		return f.AssetType != nil && *f.AssetType == models.AssetTypeTWStock
	})).Return(transactions, nil)
	mockPriceService.On("GetPrices", mock.Anything, mock.Anything).Return(prices, nil)

	// Act
	holdings, err := service.GetAllHoldings(filters)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(holdings))
	assert.Equal(t, "2330", holdings[0].Symbol)

	mockRepo.AssertExpectations(t)
	mockPriceService.AssertExpectations(t)
}

// TestGetHoldingBySymbol_Success 測試成功取得單一持倉
func TestGetHoldingBySymbol_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepositoryForHolding)
	mockPriceService := new(MockPriceService)
	fifoCalculator := NewFIFOCalculator()
	service := NewHoldingService(mockRepo, fifoCalculator, mockPriceService)

	symbol := "2330"
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
	}

	price := &models.Price{
		Symbol:    "2330",
		AssetType: models.AssetTypeTWStock,
		Price:     620,
		Currency:  "TWD",
		Source:    "mock",
		UpdatedAt: time.Now(),
	}

	// Mock 設定
	mockRepo.On("GetAll", mock.MatchedBy(func(f repository.TransactionFilters) bool {
		return f.Symbol != nil && *f.Symbol == symbol
	})).Return(transactions, nil)
	mockPriceService.On("GetPrice", symbol, models.AssetTypeTWStock).Return(price, nil)

	// Act
	holding, err := service.GetHoldingBySymbol(symbol)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, holding)
	assert.Equal(t, "2330", holding.Symbol)
	assert.Equal(t, 100.0, holding.Quantity)
	assert.Equal(t, 620.0, holding.CurrentPrice)

	mockRepo.AssertExpectations(t)
	mockPriceService.AssertExpectations(t)
}

// TestGetHoldingBySymbol_NotFound 測試標的不存在
func TestGetHoldingBySymbol_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepositoryForHolding)
	mockPriceService := new(MockPriceService)
	fifoCalculator := NewFIFOCalculator()
	service := NewHoldingService(mockRepo, fifoCalculator, mockPriceService)

	symbol := "9999"

	// Mock 設定：沒有交易記錄
	mockRepo.On("GetAll", mock.Anything).Return([]*models.Transaction{}, nil)

	// Act
	holding, err := service.GetHoldingBySymbol(symbol)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, holding)
	assert.Contains(t, err.Error(), "holding not found")

	mockRepo.AssertExpectations(t)
}

