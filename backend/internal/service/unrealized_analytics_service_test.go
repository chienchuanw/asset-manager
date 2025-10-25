package service

import (
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHoldingService for testing
type MockHoldingService struct {
	mock.Mock
}

func (m *MockHoldingService) GetAllHoldings(filters models.HoldingFilters) ([]*models.Holding, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Holding), args.Error(1)
}

func (m *MockHoldingService) GetHoldingBySymbol(symbol string) (*models.Holding, error) {
	args := m.Called(symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Holding), args.Error(1)
}

func TestUnrealizedAnalyticsService_GetSummary(t *testing.T) {
	// Arrange
	mockHoldingService := new(MockHoldingService)
	service := NewUnrealizedAnalyticsService(mockHoldingService)

	mockHoldings := []*models.Holding{
		{
			Symbol:           "2330",
			Name:             "台積電",
			AssetType:        models.AssetTypeTWStock,
			Quantity:         100,
			AvgCost:          500,
			TotalCost:        50000,
			CurrentPriceTWD:  600,
			MarketValue:      60000,
			UnrealizedPL:     10000,
			UnrealizedPLPct:  20.0,
		},
		{
			Symbol:           "AAPL",
			Name:             "Apple Inc.",
			AssetType:        models.AssetTypeUSStock,
			Quantity:         10,
			AvgCost:          3000,
			TotalCost:        30000,
			CurrentPriceTWD:  3300,
			MarketValue:      33000,
			UnrealizedPL:     3000,
			UnrealizedPLPct:  10.0,
		},
	}

	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return(mockHoldings, nil)

	// Act
	summary, err := service.GetSummary()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 80000.0, summary.TotalCost)
	assert.Equal(t, 93000.0, summary.TotalMarketValue)
	assert.Equal(t, 13000.0, summary.TotalUnrealizedPL)
	assert.InDelta(t, 16.25, summary.TotalUnrealizedPct, 0.01) // 13000/80000*100
	assert.Equal(t, 2, summary.HoldingCount)
	assert.Equal(t, "TWD", summary.Currency)

	mockHoldingService.AssertExpectations(t)
}

func TestUnrealizedAnalyticsService_GetSummary_EmptyHoldings(t *testing.T) {
	// Arrange
	mockHoldingService := new(MockHoldingService)
	service := NewUnrealizedAnalyticsService(mockHoldingService)

	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return([]*models.Holding{}, nil)

	// Act
	summary, err := service.GetSummary()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 0.0, summary.TotalCost)
	assert.Equal(t, 0.0, summary.TotalMarketValue)
	assert.Equal(t, 0.0, summary.TotalUnrealizedPL)
	assert.Equal(t, 0.0, summary.TotalUnrealizedPct)
	assert.Equal(t, 0, summary.HoldingCount)
	assert.Equal(t, "TWD", summary.Currency)

	mockHoldingService.AssertExpectations(t)
}

func TestUnrealizedAnalyticsService_GetPerformance(t *testing.T) {
	// Arrange
	mockHoldingService := new(MockHoldingService)
	service := NewUnrealizedAnalyticsService(mockHoldingService)

	mockHoldings := []*models.Holding{
		{
			Symbol:           "2330",
			AssetType:        models.AssetTypeTWStock,
			TotalCost:        50000,
			MarketValue:      60000,
			UnrealizedPL:     10000,
		},
		{
			Symbol:           "2317",
			AssetType:        models.AssetTypeTWStock,
			TotalCost:        30000,
			MarketValue:      27000,
			UnrealizedPL:     -3000,
		},
		{
			Symbol:           "AAPL",
			AssetType:        models.AssetTypeUSStock,
			TotalCost:        40000,
			MarketValue:      44000,
			UnrealizedPL:     4000,
		},
	}

	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return(mockHoldings, nil)

	// Act
	performance, err := service.GetPerformance()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, performance)
	assert.Len(t, performance, 2) // 台股和美股

	// 驗證台股績效
	var twStockPerf *models.UnrealizedPerformance
	for i := range performance {
		if performance[i].AssetType == models.AssetTypeTWStock {
			twStockPerf = &performance[i]
			break
		}
	}
	assert.NotNil(t, twStockPerf)
	assert.Equal(t, "台股", twStockPerf.Name)
	assert.Equal(t, 80000.0, twStockPerf.Cost)
	assert.Equal(t, 87000.0, twStockPerf.MarketValue)
	assert.Equal(t, 7000.0, twStockPerf.UnrealizedPL)
	assert.InDelta(t, 8.75, twStockPerf.UnrealizedPct, 0.01) // 7000/80000*100
	assert.Equal(t, 2, twStockPerf.HoldingCount)

	// 驗證美股績效
	var usStockPerf *models.UnrealizedPerformance
	for i := range performance {
		if performance[i].AssetType == models.AssetTypeUSStock {
			usStockPerf = &performance[i]
			break
		}
	}
	assert.NotNil(t, usStockPerf)
	assert.Equal(t, "美股", usStockPerf.Name)
	assert.Equal(t, 40000.0, usStockPerf.Cost)
	assert.Equal(t, 44000.0, usStockPerf.MarketValue)
	assert.Equal(t, 4000.0, usStockPerf.UnrealizedPL)
	assert.InDelta(t, 10.0, usStockPerf.UnrealizedPct, 0.01)
	assert.Equal(t, 1, usStockPerf.HoldingCount)

	mockHoldingService.AssertExpectations(t)
}

func TestUnrealizedAnalyticsService_GetTopAssets(t *testing.T) {
	// Arrange
	mockHoldingService := new(MockHoldingService)
	service := NewUnrealizedAnalyticsService(mockHoldingService)

	mockHoldings := []*models.Holding{
		{
			Symbol:          "2330",
			Name:            "台積電",
			AssetType:       models.AssetTypeTWStock,
			Quantity:        100,
			AvgCost:         500,
			CurrentPriceTWD: 600,
			TotalCost:       50000,
			MarketValue:     60000,
			UnrealizedPL:    10000,
			UnrealizedPLPct: 20.0,
		},
		{
			Symbol:          "AAPL",
			Name:            "Apple Inc.",
			AssetType:       models.AssetTypeUSStock,
			Quantity:        10,
			AvgCost:         3000,
			CurrentPriceTWD: 3500,
			TotalCost:       30000,
			MarketValue:     35000,
			UnrealizedPL:    5000,
			UnrealizedPLPct: 16.67,
		},
		{
			Symbol:          "BTC",
			Name:            "Bitcoin",
			AssetType:       models.AssetTypeCrypto,
			Quantity:        0.5,
			AvgCost:         1000000,
			CurrentPriceTWD: 900000,
			TotalCost:       500000,
			MarketValue:     450000,
			UnrealizedPL:    -50000,
			UnrealizedPLPct: -10.0,
		},
	}

	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return(mockHoldings, nil)

	// Act
	topAssets, err := service.GetTopAssets(10)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, topAssets)
	assert.Len(t, topAssets, 3)

	// 驗證排序（依未實現損益降序）
	assert.Equal(t, "2330", topAssets[0].Symbol)
	assert.Equal(t, 10000.0, topAssets[0].UnrealizedPL)

	assert.Equal(t, "AAPL", topAssets[1].Symbol)
	assert.Equal(t, 5000.0, topAssets[1].UnrealizedPL)

	assert.Equal(t, "BTC", topAssets[2].Symbol)
	assert.Equal(t, -50000.0, topAssets[2].UnrealizedPL)

	mockHoldingService.AssertExpectations(t)
}

func TestUnrealizedAnalyticsService_GetTopAssets_WithLimit(t *testing.T) {
	// Arrange
	mockHoldingService := new(MockHoldingService)
	service := NewUnrealizedAnalyticsService(mockHoldingService)

	mockHoldings := []*models.Holding{
		{Symbol: "A", UnrealizedPL: 1000},
		{Symbol: "B", UnrealizedPL: 2000},
		{Symbol: "C", UnrealizedPL: 3000},
		{Symbol: "D", UnrealizedPL: 4000},
		{Symbol: "E", UnrealizedPL: 5000},
	}

	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return(mockHoldings, nil)

	// Act
	topAssets, err := service.GetTopAssets(3)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, topAssets, 3)
	assert.Equal(t, "E", topAssets[0].Symbol) // 最高
	assert.Equal(t, "D", topAssets[1].Symbol)
	assert.Equal(t, "C", topAssets[2].Symbol)

	mockHoldingService.AssertExpectations(t)
}

