package service

import (
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHoldingServiceForAllocation 用於測試的 Mock HoldingService
type MockHoldingServiceForAllocation struct {
	mock.Mock
}

func (m *MockHoldingServiceForAllocation) GetAllHoldings(filters models.HoldingFilters) (*HoldingServiceResult, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*HoldingServiceResult), args.Error(1)
}

func (m *MockHoldingServiceForAllocation) GetHoldingBySymbol(symbol string) (*models.Holding, error) {
	args := m.Called(symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Holding), args.Error(1)
}

func (m *MockHoldingServiceForAllocation) FixInsufficientQuantity(input *models.FixInsufficientQuantityInput) (*models.Transaction, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

// TestAllocationService_GetCurrentAllocation 測試取得當前資產配置
func TestAllocationService_GetCurrentAllocation(t *testing.T) {
	mockHoldingService := new(MockHoldingServiceForAllocation)
	service := NewAllocationService(mockHoldingService)

	// 準備測試資料
	holdings := []*models.Holding{
		{
			Symbol:       "2330.TW",
			Name:         "台積電",
			AssetType:    models.AssetTypeTWStock,
			Quantity:     100,
			AvgCost:      500,
			CurrentPrice: 600,
			MarketValue:  60000,
		},
		{
			Symbol:       "AAPL",
			Name:         "Apple Inc.",
			AssetType:    models.AssetTypeUSStock,
			Quantity:     50,
			AvgCost:      150,
			CurrentPrice: 180,
			MarketValue:  9000,
		},
		{
			Symbol:       "BTC",
			Name:         "Bitcoin",
			AssetType:    models.AssetTypeCrypto,
			Quantity:     0.5,
			AvgCost:      40000,
			CurrentPrice: 50000,
			MarketValue:  25000,
		},
	}

	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return(&HoldingServiceResult{
		Holdings: holdings,
		Warnings: []*models.Warning{},
	}, nil)

	// 執行測試
	result, err := service.GetCurrentAllocation()

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 94000.0, result.TotalMarketValue)
	assert.Equal(t, "TWD", result.Currency)
	assert.Len(t, result.ByType, 3)
	assert.Len(t, result.ByAsset, 3)

	// 驗證按資產類型分類
	assert.Equal(t, models.AssetTypeTWStock, result.ByType[0].AssetType)
	assert.Equal(t, 60000.0, result.ByType[0].MarketValue)
	assert.InDelta(t, 63.83, result.ByType[0].Percentage, 0.01)
	assert.Equal(t, 1, result.ByType[0].Count)

	// 驗證按個別資產分類
	assert.Equal(t, "2330.TW", result.ByAsset[0].Symbol)
	assert.Equal(t, 60000.0, result.ByAsset[0].MarketValue)
	assert.InDelta(t, 63.83, result.ByAsset[0].Percentage, 0.01)

	mockHoldingService.AssertExpectations(t)
}

// TestAllocationService_GetCurrentAllocation_EmptyHoldings 測試空持倉情況
func TestAllocationService_GetCurrentAllocation_EmptyHoldings(t *testing.T) {
	mockHoldingService := new(MockHoldingServiceForAllocation)
	service := NewAllocationService(mockHoldingService)

	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return(&HoldingServiceResult{
		Holdings: []*models.Holding{},
		Warnings: []*models.Warning{},
	}, nil)

	result, err := service.GetCurrentAllocation()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0.0, result.TotalMarketValue)
	assert.Len(t, result.ByType, 0)
	assert.Len(t, result.ByAsset, 0)

	mockHoldingService.AssertExpectations(t)
}

// TestAllocationService_GetAllocationByType 測試取得按資產類型的配置
func TestAllocationService_GetAllocationByType(t *testing.T) {
	mockHoldingService := new(MockHoldingServiceForAllocation)
	service := NewAllocationService(mockHoldingService)

	holdings := []*models.Holding{
		{
			Symbol:       "2330.TW",
			Name:         "台積電",
			AssetType:    models.AssetTypeTWStock,
			Quantity:     100,
			AvgCost:      500,
			CurrentPrice: 600,
			MarketValue:  60000,
		},
		{
			Symbol:       "2317.TW",
			Name:         "鴻海",
			AssetType:    models.AssetTypeTWStock,
			Quantity:     200,
			AvgCost:      100,
			CurrentPrice: 120,
			MarketValue:  24000,
		},
		{
			Symbol:       "AAPL",
			Name:         "Apple Inc.",
			AssetType:    models.AssetTypeUSStock,
			Quantity:     50,
			AvgCost:      150,
			CurrentPrice: 180,
			MarketValue:  9000,
		},
	}

	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return(&HoldingServiceResult{
		Holdings: holdings,
		Warnings: []*models.Warning{},
	}, nil)

	result, err := service.GetAllocationByType()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	// 驗證台股配置（應該合併兩個台股）
	assert.Equal(t, models.AssetTypeTWStock, result[0].AssetType)
	assert.Equal(t, 84000.0, result[0].MarketValue)
	assert.InDelta(t, 90.32, result[0].Percentage, 0.01)
	assert.Equal(t, 2, result[0].Count)

	// 驗證美股配置
	assert.Equal(t, models.AssetTypeUSStock, result[1].AssetType)
	assert.Equal(t, 9000.0, result[1].MarketValue)
	assert.InDelta(t, 9.68, result[1].Percentage, 0.01)
	assert.Equal(t, 1, result[1].Count)

	mockHoldingService.AssertExpectations(t)
}

// TestAllocationService_GetAllocationByAsset 測試取得按個別資產的配置
func TestAllocationService_GetAllocationByAsset(t *testing.T) {
	mockHoldingService := new(MockHoldingServiceForAllocation)
	service := NewAllocationService(mockHoldingService)

	holdings := []*models.Holding{
		{
			Symbol:       "2330.TW",
			Name:         "台積電",
			AssetType:    models.AssetTypeTWStock,
			Quantity:     100,
			AvgCost:      500,
			CurrentPrice: 600,
			MarketValue:  60000,
		},
		{
			Symbol:       "AAPL",
			Name:         "Apple Inc.",
			AssetType:    models.AssetTypeUSStock,
			Quantity:     50,
			AvgCost:      150,
			CurrentPrice: 180,
			MarketValue:  9000,
		},
	}

	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return(&HoldingServiceResult{
		Holdings: holdings,
		Warnings: []*models.Warning{},
	}, nil)

	result, err := service.GetAllocationByAsset(10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	// 驗證排序（按市值降序）
	assert.Equal(t, "2330.TW", result[0].Symbol)
	assert.Equal(t, 60000.0, result[0].MarketValue)
	assert.InDelta(t, 86.96, result[0].Percentage, 0.01)

	assert.Equal(t, "AAPL", result[1].Symbol)
	assert.Equal(t, 9000.0, result[1].MarketValue)
	assert.InDelta(t, 13.04, result[1].Percentage, 0.01)

	mockHoldingService.AssertExpectations(t)
}

// TestAllocationService_GetAllocationByAsset_WithLimit 測試限制回傳數量
func TestAllocationService_GetAllocationByAsset_WithLimit(t *testing.T) {
	mockHoldingService := new(MockHoldingServiceForAllocation)
	service := NewAllocationService(mockHoldingService)

	holdings := []*models.Holding{
		{Symbol: "A", MarketValue: 100},
		{Symbol: "B", MarketValue: 90},
		{Symbol: "C", MarketValue: 80},
		{Symbol: "D", MarketValue: 70},
		{Symbol: "E", MarketValue: 60},
	}

	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return(&HoldingServiceResult{
		Holdings: holdings,
		Warnings: []*models.Warning{},
	}, nil)

	result, err := service.GetAllocationByAsset(3)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "A", result[0].Symbol)
	assert.Equal(t, "B", result[1].Symbol)
	assert.Equal(t, "C", result[2].Symbol)

	mockHoldingService.AssertExpectations(t)
}

