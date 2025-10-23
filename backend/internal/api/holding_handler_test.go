package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==================== Mock Objects ====================

// MockHoldingService Holdings Service 的 Mock
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

// ==================== 測試案例 ====================

// TestGetAllHoldings_Success 測試成功取得所有持倉
func TestGetAllHoldings_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockHoldingService)
	handler := NewHoldingHandler(mockService)

	// 準備測試資料
	holdings := []*models.Holding{
		{
			Symbol:          "2330",
			Name:            "台積電",
			AssetType:       models.AssetTypeTWStock,
			Quantity:        100,
			AvgCost:         500.28,
			TotalCost:       50028,
			CurrentPrice:    620,
			MarketValue:     62000,
			UnrealizedPL:    11972,
			UnrealizedPLPct: 23.93,
			LastUpdated:     time.Now(),
		},
		{
			Symbol:          "AAPL",
			Name:            "Apple Inc.",
			AssetType:       models.AssetTypeUSStock,
			Quantity:        50,
			AvgCost:         150.2,
			TotalCost:       7510,
			CurrentPrice:    175,
			MarketValue:     8750,
			UnrealizedPL:    1240,
			UnrealizedPLPct: 16.51,
			LastUpdated:     time.Now(),
		},
	}

	// Mock 設定
	mockService.On("GetAllHoldings", mock.Anything).Return(holdings, nil)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/holdings", nil)

	// Act
	handler.GetAllHoldings(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))

	// 驗證第一筆資料
	holding1 := data[0].(map[string]interface{})
	assert.Equal(t, "2330", holding1["symbol"])
	assert.Equal(t, "台積電", holding1["name"])
	assert.Equal(t, 100.0, holding1["quantity"])

	mockService.AssertExpectations(t)
}

// TestGetAllHoldings_WithAssetTypeFilter 測試按資產類型篩選
func TestGetAllHoldings_WithAssetTypeFilter(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockHoldingService)
	handler := NewHoldingHandler(mockService)

	holdings := []*models.Holding{
		{
			Symbol:          "2330",
			Name:            "台積電",
			AssetType:       models.AssetTypeTWStock,
			Quantity:        100,
			AvgCost:         500.28,
			TotalCost:       50028,
			CurrentPrice:    620,
			MarketValue:     62000,
			UnrealizedPL:    11972,
			UnrealizedPLPct: 23.93,
			LastUpdated:     time.Now(),
		},
	}

	// Mock 設定：驗證 filter 參數
	mockService.On("GetAllHoldings", mock.MatchedBy(func(f models.HoldingFilters) bool {
		return f.AssetType != nil && *f.AssetType == models.AssetTypeTWStock
	})).Return(holdings, nil)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/holdings?asset_type=tw-stock", nil)

	// Act
	handler.GetAllHoldings(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].([]interface{})
	assert.Equal(t, 1, len(data))

	mockService.AssertExpectations(t)
}

// TestGetAllHoldings_EmptyResult 測試空結果
func TestGetAllHoldings_EmptyResult(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockHoldingService)
	handler := NewHoldingHandler(mockService)

	// Mock 設定：返回空列表
	mockService.On("GetAllHoldings", mock.Anything).Return([]*models.Holding{}, nil)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/holdings", nil)

	// Act
	handler.GetAllHoldings(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].([]interface{})
	assert.Equal(t, 0, len(data))

	mockService.AssertExpectations(t)
}

// TestGetHoldingBySymbol_Success 測試成功取得單一持倉
func TestGetHoldingBySymbol_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockHoldingService)
	handler := NewHoldingHandler(mockService)

	holding := &models.Holding{
		Symbol:          "2330",
		Name:            "台積電",
		AssetType:       models.AssetTypeTWStock,
		Quantity:        100,
		AvgCost:         500.28,
		TotalCost:       50028,
		CurrentPrice:    620,
		MarketValue:     62000,
		UnrealizedPL:    11972,
		UnrealizedPLPct: 23.93,
		LastUpdated:     time.Now(),
	}

	// Mock 設定
	mockService.On("GetHoldingBySymbol", "2330").Return(holding, nil)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "symbol", Value: "2330"}}
	c.Request = httptest.NewRequest("GET", "/api/holdings/2330", nil)

	// Act
	handler.GetHoldingBySymbol(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "2330", data["symbol"])
	assert.Equal(t, "台積電", data["name"])
	assert.Equal(t, 100.0, data["quantity"])

	mockService.AssertExpectations(t)
}

// TestGetHoldingBySymbol_NotFound 測試標的不存在
func TestGetHoldingBySymbol_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockHoldingService)
	handler := NewHoldingHandler(mockService)

	// Mock 設定：返回錯誤
	mockService.On("GetHoldingBySymbol", "9999").Return(nil, assert.AnError)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "symbol", Value: "9999"}}
	c.Request = httptest.NewRequest("GET", "/api/holdings/9999", nil)

	// Act
	handler.GetHoldingBySymbol(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Nil(t, response["data"])
	assert.NotNil(t, response["error"])

	mockService.AssertExpectations(t)
}

// TestGetHoldingBySymbol_MissingSymbol 測試缺少 symbol 參數
func TestGetHoldingBySymbol_MissingSymbol(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockHoldingService)
	handler := NewHoldingHandler(mockService)

	// 建立測試請求（沒有 symbol 參數）
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/holdings/", nil)

	// Act
	handler.GetHoldingBySymbol(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Nil(t, response["data"])
	assert.NotNil(t, response["error"])
}

