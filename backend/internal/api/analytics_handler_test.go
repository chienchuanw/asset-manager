package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAnalyticsService 模擬的 AnalyticsService
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

// setupAnalyticsTestRouter 設定測試用的 router
func setupAnalyticsTestRouter(handler *AnalyticsHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api")
	{
		analytics := api.Group("/analytics")
		{
			analytics.GET("/summary", handler.GetSummary)
			analytics.GET("/performance", handler.GetPerformance)
			analytics.GET("/top-assets", handler.GetTopAssets)
		}
	}

	return router
}

// TestAnalyticsHandler_GetSummary 測試取得分析摘要
func TestAnalyticsHandler_GetSummary(t *testing.T) {
	// Arrange
	mockService := new(MockAnalyticsService)
	handler := NewAnalyticsHandler(mockService)
	router := setupAnalyticsTestRouter(handler)

	mockSummary := &models.AnalyticsSummary{
		TotalRealizedPL:    12239.0,
		TotalRealizedPLPct: 23.75,
		TotalCostBasis:     51528.0,
		TotalSellAmount:    63800.0,
		TotalSellFee:       33.0,
		TransactionCount:   2,
		Currency:           "TWD",
		TimeRange:          "month",
		StartDate:          "2025-10-01",
		EndDate:            "2025-10-31",
	}

	mockService.On("GetSummary", models.TimeRangeMonth).Return(mockSummary, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/summary?time_range=month", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.AnalyticsSummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockSummary.TotalRealizedPL, response.TotalRealizedPL)
	assert.Equal(t, mockSummary.TransactionCount, response.TransactionCount)
	assert.Equal(t, "month", response.TimeRange)

	mockService.AssertExpectations(t)
}

// TestAnalyticsHandler_GetSummary_InvalidTimeRange 測試無效的時間範圍
func TestAnalyticsHandler_GetSummary_InvalidTimeRange(t *testing.T) {
	// Arrange
	mockService := new(MockAnalyticsService)
	handler := NewAnalyticsHandler(mockService)
	router := setupAnalyticsTestRouter(handler)

	mockService.On("GetSummary", models.TimeRange("invalid")).Return(nil, fmt.Errorf("invalid time range: invalid"))

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/summary?time_range=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

// TestAnalyticsHandler_GetPerformance 測試取得績效資料
func TestAnalyticsHandler_GetPerformance(t *testing.T) {
	// Arrange
	mockService := new(MockAnalyticsService)
	handler := NewAnalyticsHandler(mockService)
	router := setupAnalyticsTestRouter(handler)

	mockPerformance := []*models.PerformanceData{
		{
			AssetType:        models.AssetTypeTWStock,
			Name:             "台股",
			RealizedPL:       9930.0,
			RealizedPLPct:    12.11,
			CostBasis:        82028.0,
			SellAmount:       92000.0,
			TransactionCount: 2,
		},
		{
			AssetType:        models.AssetTypeUSStock,
			Name:             "美股",
			RealizedPL:       295.0,
			RealizedPLPct:    19.67,
			CostBasis:        1500.0,
			SellAmount:       1800.0,
			TransactionCount: 1,
		},
	}

	mockService.On("GetPerformance", models.TimeRangeMonth).Return(mockPerformance, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/performance?time_range=month", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.PerformanceData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "台股", response[0].Name)
	assert.Equal(t, 9930.0, response[0].RealizedPL)

	mockService.AssertExpectations(t)
}

// TestAnalyticsHandler_GetTopAssets 測試取得最佳表現資產
func TestAnalyticsHandler_GetTopAssets(t *testing.T) {
	// Arrange
	mockService := new(MockAnalyticsService)
	handler := NewAnalyticsHandler(mockService)
	router := setupAnalyticsTestRouter(handler)

	mockTopAssets := []*models.TopAsset{
		{
			Symbol:        "BTC",
			Name:          "BTC",
			AssetType:     models.AssetTypeCrypto,
			RealizedPL:    200000.0,
			RealizedPLPct: 66.67,
			CostBasis:     300000.0,
			SellAmount:    500000.0,
		},
		{
			Symbol:        "2330",
			Name:          "2330",
			AssetType:     models.AssetTypeTWStock,
			RealizedPL:    11972.0,
			RealizedPLPct: 23.93,
			CostBasis:     50028.0,
			SellAmount:    62000.0,
		},
	}

	mockService.On("GetTopAssets", models.TimeRangeMonth, 5).Return(mockTopAssets, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/top-assets?time_range=month&limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.TopAsset
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "BTC", response[0].Symbol)
	assert.Equal(t, 200000.0, response[0].RealizedPL)

	mockService.AssertExpectations(t)
}

// TestAnalyticsHandler_GetTopAssets_DefaultLimit 測試預設 limit
func TestAnalyticsHandler_GetTopAssets_DefaultLimit(t *testing.T) {
	// Arrange
	mockService := new(MockAnalyticsService)
	handler := NewAnalyticsHandler(mockService)
	router := setupAnalyticsTestRouter(handler)

	mockTopAssets := []*models.TopAsset{}

	mockService.On("GetTopAssets", models.TimeRangeMonth, 5).Return(mockTopAssets, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/top-assets?time_range=month", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

