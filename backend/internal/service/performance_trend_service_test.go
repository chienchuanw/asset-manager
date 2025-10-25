package service

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockPerformanceSnapshotRepository 用於測試的 Mock Repository
type MockPerformanceSnapshotRepository struct {
	mock.Mock
}

// MockUnrealizedAnalyticsService 用於測試的 Mock Service
type MockUnrealizedAnalyticsService struct {
	mock.Mock
}

func (m *MockUnrealizedAnalyticsService) GetSummary() (*models.UnrealizedSummary, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UnrealizedSummary), args.Error(1)
}

func (m *MockUnrealizedAnalyticsService) GetPerformance() ([]models.UnrealizedPerformance, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UnrealizedPerformance), args.Error(1)
}

func (m *MockUnrealizedAnalyticsService) GetTopAssets(limit int) ([]models.UnrealizedTopAsset, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UnrealizedTopAsset), args.Error(1)
}

// MockAnalyticsService 用於測試的 Mock Service
type MockAnalyticsService struct {
	mock.Mock
}

func (m *MockAnalyticsService) GetSummary(timeRange models.TimeRange) (*models.AnalyticsSummary, error) {
	args := m.Called(timeRange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AnalyticsSummary), args.Error(1)
}

func (m *MockAnalyticsService) GetPerformance(timeRange models.TimeRange) ([]*models.PerformanceData, error) {
	args := m.Called(timeRange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PerformanceData), args.Error(1)
}

func (m *MockAnalyticsService) GetTopAssets(timeRange models.TimeRange, limit int) ([]*models.TopAsset, error) {
	args := m.Called(timeRange, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TopAsset), args.Error(1)
}

func (m *MockPerformanceSnapshotRepository) Create(input *models.CreateDailyPerformanceSnapshotInput) (*models.DailyPerformanceSnapshot, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DailyPerformanceSnapshot), args.Error(1)
}

func (m *MockPerformanceSnapshotRepository) GetByDate(date time.Time) (*models.DailyPerformanceSnapshot, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DailyPerformanceSnapshot), args.Error(1)
}

func (m *MockPerformanceSnapshotRepository) GetByDateRange(startDate, endDate time.Time) ([]*models.DailyPerformanceSnapshot, error) {
	args := m.Called(startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DailyPerformanceSnapshot), args.Error(1)
}

func (m *MockPerformanceSnapshotRepository) GetLatest(limit int) ([]*models.DailyPerformanceSnapshot, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DailyPerformanceSnapshot), args.Error(1)
}

func (m *MockPerformanceSnapshotRepository) GetDetailsBySnapshotID(snapshotID uuid.UUID) ([]*models.DailyPerformanceSnapshotDetail, error) {
	args := m.Called(snapshotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DailyPerformanceSnapshotDetail), args.Error(1)
}

func (m *MockPerformanceSnapshotRepository) GetDetailsByDateRange(startDate, endDate time.Time) (map[uuid.UUID][]*models.DailyPerformanceSnapshotDetail, error) {
	args := m.Called(startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uuid.UUID][]*models.DailyPerformanceSnapshotDetail), args.Error(1)
}

// TestPerformanceTrendService_CreateDailySnapshot 測試建立每日快照
func TestPerformanceTrendService_CreateDailySnapshot(t *testing.T) {
	mockRepo := new(MockPerformanceSnapshotRepository)
	mockUnrealizedService := new(MockUnrealizedAnalyticsService)
	mockAnalyticsService := new(MockAnalyticsService)

	service := NewPerformanceTrendService(mockRepo, mockUnrealizedService, mockAnalyticsService)

	// Mock 未實現損益摘要
	unrealizedSummary := &models.UnrealizedSummary{
		TotalCost:        900000,
		TotalMarketValue: 1000000,
		TotalUnrealizedPL: 100000,
		TotalUnrealizedPct: 11.11,
		HoldingCount:     10,
	}
	mockUnrealizedService.On("GetSummary").Return(unrealizedSummary, nil)

	// Mock 已實現損益摘要
	analyticsSummary := &models.AnalyticsSummary{
		TotalRealizedPL:    50000,
		TotalRealizedPLPct: 5.56,
		TotalCostBasis:     900000,
	}
	mockAnalyticsService.On("GetSummary", mock.Anything).Return(analyticsSummary, nil)

	// Mock 未實現損益績效（按資產類型）
	unrealizedPerformance := []models.UnrealizedPerformance{
		{
			AssetType:     models.AssetTypeTWStock,
			Cost:          450000,
			MarketValue:   500000,
			UnrealizedPL:  50000,
			UnrealizedPct: 11.11,
			HoldingCount:  5,
		},
		{
			AssetType:     models.AssetTypeUSStock,
			Cost:          450000,
			MarketValue:   500000,
			UnrealizedPL:  50000,
			UnrealizedPct: 11.11,
			HoldingCount:  5,
		},
	}
	mockUnrealizedService.On("GetPerformance").Return(unrealizedPerformance, nil)

	// Mock 已實現損益績效（按資產類型）
	analyticsPerformance := []*models.PerformanceData{
		{
			AssetType:     models.AssetTypeTWStock,
			RealizedPL:    25000,
			RealizedPLPct: 5.56,
			CostBasis:     450000,
		},
		{
			AssetType:     models.AssetTypeUSStock,
			RealizedPL:    25000,
			RealizedPLPct: 5.56,
			CostBasis:     450000,
		},
	}
	mockAnalyticsService.On("GetPerformance", mock.Anything).Return(analyticsPerformance, nil)

	// Mock Repository Create
	expectedSnapshot := &models.DailyPerformanceSnapshot{
		ID:                 uuid.New(),
		SnapshotDate:       time.Now().Truncate(24 * time.Hour),
		TotalMarketValue:   1000000,
		TotalCost:          900000,
		TotalUnrealizedPL:  100000,
		TotalUnrealizedPct: 11.11,
		TotalRealizedPL:    50000,
		TotalRealizedPct:   5.56,
		HoldingCount:       10,
		Currency:           "TWD",
	}
	mockRepo.On("Create", mock.AnythingOfType("*models.CreateDailyPerformanceSnapshotInput")).Return(expectedSnapshot, nil)

	// 執行測試
	snapshot, err := service.CreateDailySnapshot()

	// 驗證
	require.NoError(t, err)
	assert.NotNil(t, snapshot)
	assert.Equal(t, expectedSnapshot.TotalMarketValue, snapshot.TotalMarketValue)
	assert.Equal(t, expectedSnapshot.TotalCost, snapshot.TotalCost)
	assert.Equal(t, expectedSnapshot.HoldingCount, snapshot.HoldingCount)

	mockRepo.AssertExpectations(t)
	mockUnrealizedService.AssertExpectations(t)
	mockAnalyticsService.AssertExpectations(t)
}

// TestPerformanceTrendService_GetTrendByDateRange 測試取得日期範圍內的趨勢
func TestPerformanceTrendService_GetTrendByDateRange(t *testing.T) {
	mockRepo := new(MockPerformanceSnapshotRepository)
	mockUnrealizedService := new(MockUnrealizedAnalyticsService)
	mockAnalyticsService := new(MockAnalyticsService)

	service := NewPerformanceTrendService(mockRepo, mockUnrealizedService, mockAnalyticsService)

	startDate := time.Date(2025, 10, 23, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC)

	// Mock 快照資料
	snapshots := []*models.DailyPerformanceSnapshot{
		{
			ID:                 uuid.New(),
			SnapshotDate:       time.Date(2025, 10, 23, 0, 0, 0, 0, time.UTC),
			TotalMarketValue:   900000,
			TotalCost:          850000,
			TotalUnrealizedPL:  50000,
			TotalUnrealizedPct: 5.88,
			TotalRealizedPL:    30000,
			TotalRealizedPct:   3.53,
			HoldingCount:       8,
			Currency:           "TWD",
		},
		{
			ID:                 uuid.New(),
			SnapshotDate:       time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC),
			TotalMarketValue:   950000,
			TotalCost:          875000,
			TotalUnrealizedPL:  75000,
			TotalUnrealizedPct: 8.57,
			TotalRealizedPL:    40000,
			TotalRealizedPct:   4.57,
			HoldingCount:       9,
			Currency:           "TWD",
		},
		{
			ID:                 uuid.New(),
			SnapshotDate:       time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
			TotalMarketValue:   1000000,
			TotalCost:          900000,
			TotalUnrealizedPL:  100000,
			TotalUnrealizedPct: 11.11,
			TotalRealizedPL:    50000,
			TotalRealizedPct:   5.56,
			HoldingCount:       10,
			Currency:           "TWD",
		},
	}
	mockRepo.On("GetByDateRange", startDate, endDate).Return(snapshots, nil)

	// Mock 明細資料（空的也可以）
	mockRepo.On("GetDetailsByDateRange", startDate, endDate).Return(make(map[uuid.UUID][]*models.DailyPerformanceSnapshotDetail), nil)

	// 執行測試
	summary, err := service.GetTrendByDateRange(startDate, endDate)

	// 驗證
	require.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, startDate, summary.StartDate)
	assert.Equal(t, endDate, summary.EndDate)
	assert.Len(t, summary.TotalData, 3)
	assert.Equal(t, 3, summary.DataPointCount)
	assert.Equal(t, "TWD", summary.Currency)

	// 驗證資料點
	assert.Equal(t, float64(900000), summary.TotalData[0].MarketValue)
	assert.Equal(t, float64(950000), summary.TotalData[1].MarketValue)
	assert.Equal(t, float64(1000000), summary.TotalData[2].MarketValue)

	mockRepo.AssertExpectations(t)
}

// TestPerformanceTrendService_GetLatestTrend 測試取得最新趨勢
func TestPerformanceTrendService_GetLatestTrend(t *testing.T) {
	mockRepo := new(MockPerformanceSnapshotRepository)
	mockUnrealizedService := new(MockUnrealizedAnalyticsService)
	mockAnalyticsService := new(MockAnalyticsService)

	service := NewPerformanceTrendService(mockRepo, mockUnrealizedService, mockAnalyticsService)

	// Mock 快照資料
	snapshots := []*models.DailyPerformanceSnapshot{
		{
			ID:                 uuid.New(),
			SnapshotDate:       time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
			TotalMarketValue:   1000000,
			TotalCost:          900000,
			TotalUnrealizedPL:  100000,
			TotalUnrealizedPct: 11.11,
			TotalRealizedPL:    50000,
			TotalRealizedPct:   5.56,
			HoldingCount:       10,
			Currency:           "TWD",
		},
	}
	mockRepo.On("GetLatest", 30).Return(snapshots, nil)

	// 執行測試
	data, err := service.GetLatestTrend(30)

	// 驗證
	require.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data, 1)
	assert.Equal(t, float64(1000000), data[0].MarketValue)

	mockRepo.AssertExpectations(t)
}

