package service

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRealizedProfitRepositoryForAnalytics 模擬的 RealizedProfitRepository（用於 Analytics）
type MockRealizedProfitRepositoryForAnalytics struct {
	mock.Mock
}

func (m *MockRealizedProfitRepositoryForAnalytics) Create(input *models.CreateRealizedProfitInput) (*models.RealizedProfit, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RealizedProfit), args.Error(1)
}

func (m *MockRealizedProfitRepositoryForAnalytics) GetByTransactionID(transactionID string) (*models.RealizedProfit, error) {
	args := m.Called(transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RealizedProfit), args.Error(1)
}

func (m *MockRealizedProfitRepositoryForAnalytics) GetAll(filters models.RealizedProfitFilters) ([]*models.RealizedProfit, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.RealizedProfit), args.Error(1)
}

func (m *MockRealizedProfitRepositoryForAnalytics) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestAnalyticsService_GetSummary 測試取得分析摘要
func TestAnalyticsService_GetSummary(t *testing.T) {
	// Arrange
	mockRepo := new(MockRealizedProfitRepositoryForAnalytics)
	service := NewAnalyticsService(mockRepo)

	// 模擬已實現損益記錄
	mockRecords := []*models.RealizedProfit{
		{
			ID:            "1",
			Symbol:        "2330",
			AssetType:     models.AssetTypeTWStock,
			SellDate:      time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC),
			Quantity:      100,
			SellPrice:     620,
			SellAmount:    62000,
			SellFee:       28,
			CostBasis:     50028,
			RealizedPL:    11944,
			RealizedPLPct: 23.87,
			Currency:      "TWD",
		},
		{
			ID:            "2",
			Symbol:        "AAPL",
			AssetType:     models.AssetTypeUSStock,
			SellDate:      time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
			Quantity:      10,
			SellPrice:     180,
			SellAmount:    1800,
			SellFee:       5,
			CostBasis:     1500,
			RealizedPL:    295,
			RealizedPLPct: 19.67,
			Currency:      "USD",
		},
	}

	// 使用 mock.MatchedBy 來匹配任何包含 StartDate 和 EndDate 的 filters
	mockRepo.On("GetAll", mock.MatchedBy(func(filters models.RealizedProfitFilters) bool {
		return filters.StartDate != nil && filters.EndDate != nil
	})).Return(mockRecords, nil)

	// Act
	summary, err := service.GetSummary(models.TimeRangeMonth)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 2, summary.TransactionCount)
	assert.Equal(t, 12239.0, summary.TotalRealizedPL) // 11944 + 295
	assert.Equal(t, 51528.0, summary.TotalCostBasis)  // 50028 + 1500
	assert.Equal(t, 63800.0, summary.TotalSellAmount) // 62000 + 1800
	assert.Equal(t, 33.0, summary.TotalSellFee)       // 28 + 5
	assert.InDelta(t, 23.75, summary.TotalRealizedPLPct, 0.01)
	assert.Equal(t, "month", summary.TimeRange)
	mockRepo.AssertExpectations(t)
}

// TestAnalyticsService_GetPerformance 測試取得績效資料
func TestAnalyticsService_GetPerformance(t *testing.T) {
	// Arrange
	mockRepo := new(MockRealizedProfitRepositoryForAnalytics)
	service := NewAnalyticsService(mockRepo)

	mockRecords := []*models.RealizedProfit{
		{
			Symbol:        "2330",
			AssetType:     models.AssetTypeTWStock,
			SellDate:      time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC),
			SellAmount:    62000,
			SellFee:       28,
			CostBasis:     50028,
			RealizedPL:    11944,
			RealizedPLPct: 23.87,
		},
		{
			Symbol:        "2317",
			AssetType:     models.AssetTypeTWStock,
			SellDate:      time.Date(2025, 10, 12, 0, 0, 0, 0, time.UTC),
			SellAmount:    30000,
			SellFee:       14,
			CostBasis:     32000,
			RealizedPL:    -2014,
			RealizedPLPct: -6.29,
		},
		{
			Symbol:        "AAPL",
			AssetType:     models.AssetTypeUSStock,
			SellDate:      time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
			SellAmount:    1800,
			SellFee:       5,
			CostBasis:     1500,
			RealizedPL:    295,
			RealizedPLPct: 19.67,
		},
	}

	mockRepo.On("GetAll", mock.MatchedBy(func(filters models.RealizedProfitFilters) bool {
		return filters.StartDate != nil && filters.EndDate != nil
	})).Return(mockRecords, nil)

	// Act
	performance, err := service.GetPerformance(models.TimeRangeMonth)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, performance)
	assert.Len(t, performance, 2) // 台股和美股

	// 檢查台股績效
	twStock := findPerformanceByAssetType(performance, models.AssetTypeTWStock)
	assert.NotNil(t, twStock)
	assert.Equal(t, "台股", twStock.Name)
	assert.Equal(t, 9930.0, twStock.RealizedPL) // 11944 + (-2014)
	assert.Equal(t, 82028.0, twStock.CostBasis) // 50028 + 32000
	assert.Equal(t, 2, twStock.TransactionCount)

	// 檢查美股績效
	usStock := findPerformanceByAssetType(performance, models.AssetTypeUSStock)
	assert.NotNil(t, usStock)
	assert.Equal(t, "美股", usStock.Name)
	assert.Equal(t, 295.0, usStock.RealizedPL)
	assert.Equal(t, 1500.0, usStock.CostBasis)
	assert.Equal(t, 1, usStock.TransactionCount)

	mockRepo.AssertExpectations(t)
}

// TestAnalyticsService_GetTopAssets 測試取得最佳/最差表現資產
func TestAnalyticsService_GetTopAssets(t *testing.T) {
	// Arrange
	mockRepo := new(MockRealizedProfitRepositoryForAnalytics)
	service := NewAnalyticsService(mockRepo)

	mockRecords := []*models.RealizedProfit{
		{
			Symbol:        "2330",
			AssetType:     models.AssetTypeTWStock,
			SellAmount:    62000,
			CostBasis:     50028,
			RealizedPL:    11972,
			RealizedPLPct: 23.93,
		},
		{
			Symbol:        "2317",
			AssetType:     models.AssetTypeTWStock,
			SellAmount:    30000,
			CostBasis:     32000,
			RealizedPL:    -2000,
			RealizedPLPct: -6.25,
		},
		{
			Symbol:        "AAPL",
			AssetType:     models.AssetTypeUSStock,
			SellAmount:    1800,
			CostBasis:     1500,
			RealizedPL:    300,
			RealizedPLPct: 20.0,
		},
		{
			Symbol:        "BTC",
			AssetType:     models.AssetTypeCrypto,
			SellAmount:    500000,
			CostBasis:     300000,
			RealizedPL:    200000,
			RealizedPLPct: 66.67,
		},
	}

	mockRepo.On("GetAll", mock.MatchedBy(func(filters models.RealizedProfitFilters) bool {
		return filters.StartDate != nil && filters.EndDate != nil
	})).Return(mockRecords, nil)

	// Act
	topAssets, err := service.GetTopAssets(models.TimeRangeMonth, 3)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, topAssets)
	assert.Len(t, topAssets, 3)

	// 檢查排序（應該按已實現損益由高到低）
	assert.Equal(t, "BTC", topAssets[0].Symbol)
	assert.Equal(t, 200000.0, topAssets[0].RealizedPL)

	assert.Equal(t, "2330", topAssets[1].Symbol)
	assert.Equal(t, 11972.0, topAssets[1].RealizedPL)

	assert.Equal(t, "AAPL", topAssets[2].Symbol)
	assert.Equal(t, 300.0, topAssets[2].RealizedPL)

	mockRepo.AssertExpectations(t)
}

// 輔助函式：根據資產類型尋找績效資料
func findPerformanceByAssetType(performance []*models.PerformanceData, assetType models.AssetType) *models.PerformanceData {
	for _, p := range performance {
		if p.AssetType == assetType {
			return p
		}
	}
	return nil
}

