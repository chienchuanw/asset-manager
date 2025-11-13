package service

import (
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSettingsServiceForRebalance Settings Service 的 Mock
type MockSettingsServiceForRebalance struct {
	mock.Mock
}

func (m *MockSettingsServiceForRebalance) GetSettings() (*models.SettingsGroup, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SettingsGroup), args.Error(1)
}

func (m *MockSettingsServiceForRebalance) UpdateSettings(input *models.UpdateSettingsGroupInput) (*models.SettingsGroup, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SettingsGroup), args.Error(1)
}

// MockHoldingServiceForRebalance Holding Service 的 Mock
type MockHoldingServiceForRebalance struct {
	mock.Mock
}

func (m *MockHoldingServiceForRebalance) GetAllHoldings(filters models.HoldingFilters) (*HoldingServiceResult, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*HoldingServiceResult), args.Error(1)
}

func (m *MockHoldingServiceForRebalance) GetHoldingBySymbol(symbol string) (*models.Holding, error) {
	args := m.Called(symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Holding), args.Error(1)
}

// ==================== 測試案例 ====================

// TestCheckRebalance_NoRebalanceNeeded 測試不需要再平衡的情況
func TestCheckRebalance_NoRebalanceNeeded(t *testing.T) {
	// Arrange
	mockSettings := new(MockSettingsServiceForRebalance)
	mockHolding := new(MockHoldingServiceForRebalance)
	service := NewRebalanceService(mockSettings, mockHolding)

	// Mock 設定：目標配置 50% 台股、30% 美股、20% 加密貨幣，閾值 5%
	settings := &models.SettingsGroup{
		Allocation: models.AllocationSettings{
			TWStock:            50.0,
			USStock:            30.0,
			Crypto:             20.0,
			RebalanceThreshold: 5.0,
		},
	}

	// Mock 持倉：總資產 1,000,000，實際配置 51% 台股、29% 美股、20% 加密貨幣（偏離都在 5% 以內）
	holdings := []*models.Holding{
		{Symbol: "2330", AssetType: models.AssetTypeTWStock, MarketValue: 510000}, // 51%
		{Symbol: "AAPL", AssetType: models.AssetTypeUSStock, MarketValue: 290000}, // 29%
		{Symbol: "BTC", AssetType: models.AssetTypeCrypto, MarketValue: 200000},   // 20%
	}

	mockSettings.On("GetSettings").Return(settings, nil)
	mockHolding.On("GetAllHoldings", mock.Anything).Return(&HoldingServiceResult{
		Holdings: holdings,
		Warnings: []*models.Warning{},
	}, nil)

	// Act
	result, err := service.CheckRebalance()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.NeedsRebalance)
	assert.Equal(t, 5.0, result.Threshold)
	assert.Equal(t, 1000000.0, result.CurrentTotal)
	assert.Equal(t, 3, len(result.Deviations))

	// 驗證台股偏離（deviations 按資產類型字母排序：crypto, tw-stock, us-stock）
	var twDeviation *models.AssetTypeDeviation
	for i := range result.Deviations {
		if result.Deviations[i].AssetType == "tw-stock" {
			twDeviation = &result.Deviations[i]
			break
		}
	}
	assert.NotNil(t, twDeviation)
	assert.Equal(t, "tw-stock", twDeviation.AssetType)
	assert.Equal(t, 50.0, twDeviation.TargetPercent)
	assert.Equal(t, 51.0, twDeviation.CurrentPercent)
	assert.Equal(t, 1.0, twDeviation.Deviation)
	assert.False(t, twDeviation.ExceedsThreshold)

	mockSettings.AssertExpectations(t)
	mockHolding.AssertExpectations(t)
}

// TestCheckRebalance_RebalanceNeeded 測試需要再平衡的情況
func TestCheckRebalance_RebalanceNeeded(t *testing.T) {
	// Arrange
	mockSettings := new(MockSettingsServiceForRebalance)
	mockHolding := new(MockHoldingServiceForRebalance)
	service := NewRebalanceService(mockSettings, mockHolding)

	// Mock 設定：目標配置 50% 台股、30% 美股、20% 加密貨幣，閾值 5%
	settings := &models.SettingsGroup{
		Allocation: models.AllocationSettings{
			TWStock:            50.0,
			USStock:            30.0,
			Crypto:             20.0,
			RebalanceThreshold: 5.0,
		},
	}

	// Mock 持倉：總資產 1,000,000，實際配置 60% 台股、25% 美股、15% 加密貨幣（台股偏離 10%，超過閾值）
	holdings := []*models.Holding{
		{Symbol: "2330", AssetType: models.AssetTypeTWStock, MarketValue: 600000}, // 60%
		{Symbol: "AAPL", AssetType: models.AssetTypeUSStock, MarketValue: 250000}, // 25%
		{Symbol: "BTC", AssetType: models.AssetTypeCrypto, MarketValue: 150000},   // 15%
	}

	mockSettings.On("GetSettings").Return(settings, nil)
	mockHolding.On("GetAllHoldings", mock.Anything).Return(holdings, nil)

	// Act
	result, err := service.CheckRebalance()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.NeedsRebalance)
	assert.Equal(t, 5.0, result.Threshold)
	assert.Equal(t, 1000000.0, result.CurrentTotal)

	// 驗證台股偏離（超過閾值）
	var twDeviation *models.AssetTypeDeviation
	for i := range result.Deviations {
		if result.Deviations[i].AssetType == "tw-stock" {
			twDeviation = &result.Deviations[i]
			break
		}
	}
	assert.NotNil(t, twDeviation)
	assert.Equal(t, 50.0, twDeviation.TargetPercent)
	assert.Equal(t, 60.0, twDeviation.CurrentPercent)
	assert.Equal(t, 10.0, twDeviation.Deviation)
	assert.True(t, twDeviation.ExceedsThreshold)

	// 驗證建議（應該有賣出台股、買入美股和加密貨幣的建議）
	assert.Greater(t, len(result.Suggestions), 0)

	mockSettings.AssertExpectations(t)
	mockHolding.AssertExpectations(t)
}

// TestCheckRebalance_EmptyHoldings 測試空持倉的情況
func TestCheckRebalance_EmptyHoldings(t *testing.T) {
	// Arrange
	mockSettings := new(MockSettingsServiceForRebalance)
	mockHolding := new(MockHoldingServiceForRebalance)
	service := NewRebalanceService(mockSettings, mockHolding)

	settings := &models.SettingsGroup{
		Allocation: models.AllocationSettings{
			TWStock:            50.0,
			USStock:            30.0,
			Crypto:             20.0,
			RebalanceThreshold: 5.0,
		},
	}

	holdings := []*models.Holding{}

	mockSettings.On("GetSettings").Return(settings, nil)
	mockHolding.On("GetAllHoldings", mock.Anything).Return(&HoldingServiceResult{
		Holdings: holdings,
		Warnings: []*models.Warning{},
	}, nil)

	// Act
	result, err := service.CheckRebalance()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.NeedsRebalance)
	assert.Equal(t, 0.0, result.CurrentTotal)
	assert.Equal(t, 0, len(result.Deviations))
	assert.Equal(t, 0, len(result.Suggestions))

	mockSettings.AssertExpectations(t)
	mockHolding.AssertExpectations(t)
}

// TestCheckRebalance_ZeroThreshold 測試閾值為 0 的情況（任何偏離都需要再平衡）
func TestCheckRebalance_ZeroThreshold(t *testing.T) {
	// Arrange
	mockSettings := new(MockSettingsServiceForRebalance)
	mockHolding := new(MockHoldingServiceForRebalance)
	service := NewRebalanceService(mockSettings, mockHolding)

	// Mock 設定：閾值為 0
	settings := &models.SettingsGroup{
		Allocation: models.AllocationSettings{
			TWStock:            50.0,
			USStock:            30.0,
			Crypto:             20.0,
			RebalanceThreshold: 0.0,
		},
	}

	// Mock 持倉：即使只偏離 1%，也應該觸發再平衡
	holdings := []*models.Holding{
		{Symbol: "2330", AssetType: models.AssetTypeTWStock, MarketValue: 510000}, // 51%
		{Symbol: "AAPL", AssetType: models.AssetTypeUSStock, MarketValue: 290000}, // 29%
		{Symbol: "BTC", AssetType: models.AssetTypeCrypto, MarketValue: 200000},   // 20%
	}

	mockSettings.On("GetSettings").Return(settings, nil)
	mockHolding.On("GetAllHoldings", mock.Anything).Return(&HoldingServiceResult{
		Holdings: holdings,
		Warnings: []*models.Warning{},
	}, nil)

	// Act
	result, err := service.CheckRebalance()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.NeedsRebalance) // 因為閾值為 0，任何偏離都需要再平衡
	assert.Equal(t, 0.0, result.Threshold)

	mockSettings.AssertExpectations(t)
	mockHolding.AssertExpectations(t)
}

