package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAllocationService 用於測試的 Mock AllocationService
type MockAllocationService struct {
	mock.Mock
}

func (m *MockAllocationService) GetCurrentAllocation() (*models.AllocationSummary, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AllocationSummary), args.Error(1)
}

func (m *MockAllocationService) GetAllocationByType() ([]models.AllocationByType, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AllocationByType), args.Error(1)
}

func (m *MockAllocationService) GetAllocationByAsset(limit int) ([]models.AllocationByAsset, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AllocationByAsset), args.Error(1)
}

// TestAllocationHandler_GetCurrentAllocation 測試取得當前資產配置
func TestAllocationHandler_GetCurrentAllocation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAllocationService)
	handler := NewAllocationHandler(mockService)

	// 準備測試資料
	summary := &models.AllocationSummary{
		TotalMarketValue: 100000,
		ByType: []models.AllocationByType{
			{AssetType: models.AssetTypeTWStock, Name: "台股", MarketValue: 60000, Percentage: 60, Count: 2},
		},
		ByAsset: []models.AllocationByAsset{
			{Symbol: "2330.TW", Name: "台積電", MarketValue: 60000, Percentage: 60},
		},
		Currency: "TWD",
		AsOfDate: time.Now(),
	}

	mockService.On("GetCurrentAllocation").Return(summary, nil)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/allocation/current", nil)

	// 執行測試
	handler.GetCurrentAllocation(c)

	// 驗證結果
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestAllocationHandler_GetCurrentAllocation_ServiceError 測試服務錯誤
func TestAllocationHandler_GetCurrentAllocation_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAllocationService)
	handler := NewAllocationHandler(mockService)

	mockService.On("GetCurrentAllocation").Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/allocation/current", nil)

	handler.GetCurrentAllocation(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "GET_CURRENT_ALLOCATION_FAILED", response.Error.Code)

	mockService.AssertExpectations(t)
}

// TestAllocationHandler_GetAllocationByType 測試取得按資產類型的配置
func TestAllocationHandler_GetAllocationByType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAllocationService)
	handler := NewAllocationHandler(mockService)

	allocations := []models.AllocationByType{
		{AssetType: models.AssetTypeTWStock, Name: "台股", MarketValue: 60000, Percentage: 60, Count: 2},
		{AssetType: models.AssetTypeUSStock, Name: "美股", MarketValue: 40000, Percentage: 40, Count: 1},
	}

	mockService.On("GetAllocationByType").Return(allocations, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/allocation/by-type", nil)

	handler.GetAllocationByType(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestAllocationHandler_GetAllocationByAsset 測試取得按個別資產的配置
func TestAllocationHandler_GetAllocationByAsset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAllocationService)
	handler := NewAllocationHandler(mockService)

	allocations := []models.AllocationByAsset{
		{Symbol: "2330.TW", Name: "台積電", MarketValue: 60000, Percentage: 60},
		{Symbol: "AAPL", Name: "Apple Inc.", MarketValue: 40000, Percentage: 40},
	}

	mockService.On("GetAllocationByAsset", 10).Return(allocations, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/allocation/by-asset?limit=10", nil)

	handler.GetAllocationByAsset(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestAllocationHandler_GetAllocationByAsset_DefaultLimit 測試預設限制
func TestAllocationHandler_GetAllocationByAsset_DefaultLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAllocationService)
	handler := NewAllocationHandler(mockService)

	allocations := []models.AllocationByAsset{}

	mockService.On("GetAllocationByAsset", 20).Return(allocations, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/allocation/by-asset", nil)

	handler.GetAllocationByAsset(c)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestAllocationHandler_GetAllocationByAsset_InvalidLimit 測試無效限制
func TestAllocationHandler_GetAllocationByAsset_InvalidLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAllocationService)
	handler := NewAllocationHandler(mockService)

	allocations := []models.AllocationByAsset{}

	mockService.On("GetAllocationByAsset", 20).Return(allocations, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/allocation/by-asset?limit=invalid", nil)

	handler.GetAllocationByAsset(c)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

