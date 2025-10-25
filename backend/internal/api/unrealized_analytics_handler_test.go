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

// MockUnrealizedAnalyticsService 模擬的 UnrealizedAnalyticsService
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

// setupUnrealizedAnalyticsTestRouter 設定測試用的 router
func setupUnrealizedAnalyticsTestRouter(handler *UnrealizedAnalyticsHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api")
	{
		analytics := api.Group("/analytics")
		{
			unrealized := analytics.Group("/unrealized")
			{
				unrealized.GET("/summary", handler.GetSummary)
				unrealized.GET("/performance", handler.GetPerformance)
				unrealized.GET("/top-assets", handler.GetTopAssets)
			}
		}
	}

	return router
}

// TestUnrealizedAnalyticsHandler_GetSummary 測試取得未實現損益摘要
func TestUnrealizedAnalyticsHandler_GetSummary(t *testing.T) {
	// Arrange
	mockService := new(MockUnrealizedAnalyticsService)
	handler := NewUnrealizedAnalyticsHandler(mockService)
	router := setupUnrealizedAnalyticsTestRouter(handler)

	mockSummary := &models.UnrealizedSummary{
		TotalCost:          100000.0,
		TotalMarketValue:   120000.0,
		TotalUnrealizedPL:  20000.0,
		TotalUnrealizedPct: 20.0,
		HoldingCount:       5,
		Currency:           "TWD",
	}

	mockService.On("GetSummary").Return(mockSummary, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/unrealized/summary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	// 驗證回傳的資料
	dataBytes, _ := json.Marshal(response.Data)
	var summary models.UnrealizedSummary
	json.Unmarshal(dataBytes, &summary)

	assert.Equal(t, 100000.0, summary.TotalCost)
	assert.Equal(t, 120000.0, summary.TotalMarketValue)
	assert.Equal(t, 20000.0, summary.TotalUnrealizedPL)
	assert.Equal(t, 20.0, summary.TotalUnrealizedPct)
	assert.Equal(t, 5, summary.HoldingCount)
	assert.Equal(t, "TWD", summary.Currency)

	mockService.AssertExpectations(t)
}

// TestUnrealizedAnalyticsHandler_GetSummary_ServiceError 測試服務錯誤
func TestUnrealizedAnalyticsHandler_GetSummary_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockUnrealizedAnalyticsService)
	handler := NewUnrealizedAnalyticsHandler(mockService)
	router := setupUnrealizedAnalyticsTestRouter(handler)

	mockService.On("GetSummary").Return((*models.UnrealizedSummary)(nil), fmt.Errorf("service error"))

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/unrealized/summary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "GET_SUMMARY_FAILED", response.Error.Code)

	mockService.AssertExpectations(t)
}

// TestUnrealizedAnalyticsHandler_GetPerformance 測試取得各資產類型績效
func TestUnrealizedAnalyticsHandler_GetPerformance(t *testing.T) {
	// Arrange
	mockService := new(MockUnrealizedAnalyticsService)
	handler := NewUnrealizedAnalyticsHandler(mockService)
	router := setupUnrealizedAnalyticsTestRouter(handler)

	mockPerformance := []models.UnrealizedPerformance{
		{
			AssetType:     models.AssetTypeTWStock,
			Name:          "台股",
			Cost:          50000.0,
			MarketValue:   60000.0,
			UnrealizedPL:  10000.0,
			UnrealizedPct: 20.0,
			HoldingCount:  3,
		},
		{
			AssetType:     models.AssetTypeUSStock,
			Name:          "美股",
			Cost:          30000.0,
			MarketValue:   33000.0,
			UnrealizedPL:  3000.0,
			UnrealizedPct: 10.0,
			HoldingCount:  2,
		},
	}

	mockService.On("GetPerformance").Return(mockPerformance, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/unrealized/performance", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestUnrealizedAnalyticsHandler_GetTopAssets 測試取得 Top 資產
func TestUnrealizedAnalyticsHandler_GetTopAssets(t *testing.T) {
	// Arrange
	mockService := new(MockUnrealizedAnalyticsService)
	handler := NewUnrealizedAnalyticsHandler(mockService)
	router := setupUnrealizedAnalyticsTestRouter(handler)

	mockTopAssets := []models.UnrealizedTopAsset{
		{
			Symbol:        "2330",
			Name:          "台積電",
			AssetType:     models.AssetTypeTWStock,
			Quantity:      100,
			AvgCost:       500,
			CurrentPrice:  600,
			Cost:          50000,
			MarketValue:   60000,
			UnrealizedPL:  10000,
			UnrealizedPct: 20.0,
		},
	}

	mockService.On("GetTopAssets", 10).Return(mockTopAssets, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/unrealized/top-assets?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestUnrealizedAnalyticsHandler_GetTopAssets_DefaultLimit 測試預設 limit
func TestUnrealizedAnalyticsHandler_GetTopAssets_DefaultLimit(t *testing.T) {
	// Arrange
	mockService := new(MockUnrealizedAnalyticsService)
	handler := NewUnrealizedAnalyticsHandler(mockService)
	router := setupUnrealizedAnalyticsTestRouter(handler)

	mockService.On("GetTopAssets", 10).Return([]models.UnrealizedTopAsset{}, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/unrealized/top-assets", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// TestUnrealizedAnalyticsHandler_GetTopAssets_InvalidLimit 測試無效的 limit
func TestUnrealizedAnalyticsHandler_GetTopAssets_InvalidLimit(t *testing.T) {
	// Arrange
	mockService := new(MockUnrealizedAnalyticsService)
	handler := NewUnrealizedAnalyticsHandler(mockService)
	router := setupUnrealizedAnalyticsTestRouter(handler)

	// 無效的 limit 應該使用預設值 10
	mockService.On("GetTopAssets", 10).Return([]models.UnrealizedTopAsset{}, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/analytics/unrealized/top-assets?limit=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

