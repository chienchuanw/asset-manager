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

// MockExchangeRateService 是 ExchangeRateService 的 mock 實作
type MockExchangeRateService struct {
	mock.Mock
}

func (m *MockExchangeRateService) GetRate(fromCurrency, toCurrency models.Currency, date time.Time) (float64, error) {
	args := m.Called(fromCurrency, toCurrency, date)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExchangeRateService) GetRateRecord(fromCurrency, toCurrency models.Currency, date time.Time) (*models.ExchangeRate, error) {
	args := m.Called(fromCurrency, toCurrency, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExchangeRate), args.Error(1)
}

func (m *MockExchangeRateService) GetTodayRate(fromCurrency, toCurrency models.Currency) (float64, error) {
	args := m.Called(fromCurrency, toCurrency)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExchangeRateService) RefreshTodayRate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockExchangeRateService) ConvertToTWD(amount float64, currency models.Currency, date time.Time) (float64, error) {
	args := m.Called(amount, currency, date)
	return args.Get(0).(float64), args.Error(1)
}

// TestRefreshExchangeRate_Success 測試成功更新匯率
func TestRefreshExchangeRate_Success(t *testing.T) {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)

	// 建立 mock service
	mockService := new(MockExchangeRateService)
	
	// 設定 mock 行為
	mockService.On("RefreshTodayRate").Return(nil)
	
	// 模擬更新後的匯率記錄
	now := time.Now()
	mockService.On("GetRateRecord", models.CurrencyUSD, models.CurrencyTWD, mock.Anything).Return(&models.ExchangeRate{
		ID:           1,
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         30.6,
		Date:         now.Truncate(24 * time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil)

	// 建立 handler
	handler := NewExchangeRateHandler(mockService)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/exchange-rates/refresh", nil)

	// 執行 handler
	handler.RefreshExchangeRate(c)

	// 驗證回應
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	// 驗證回應資料
	dataMap, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "USD", dataMap["from_currency"])
	assert.Equal(t, "TWD", dataMap["to_currency"])
	assert.Equal(t, 30.6, dataMap["rate"])
	assert.NotEmpty(t, dataMap["updated_at"])

	// 驗證 mock 被呼叫
	mockService.AssertExpectations(t)
}

// TestRefreshExchangeRate_RefreshFailed 測試更新匯率失敗
func TestRefreshExchangeRate_RefreshFailed(t *testing.T) {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)

	// 建立 mock service
	mockService := new(MockExchangeRateService)
	
	// 設定 mock 行為 - 更新失敗
	mockService.On("RefreshTodayRate").Return(errors.New("API connection failed"))

	// 建立 handler
	handler := NewExchangeRateHandler(mockService)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/exchange-rates/refresh", nil)

	// 執行 handler
	handler.RefreshExchangeRate(c)

	// 驗證回應
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "REFRESH_RATE_FAILED", response.Error.Code)
	assert.Contains(t, response.Error.Message, "API connection failed")

	// 驗證 mock 被呼叫
	mockService.AssertExpectations(t)
}

// TestRefreshExchangeRate_GetRecordFailed 測試取得更新後的記錄失敗
func TestRefreshExchangeRate_GetRecordFailed(t *testing.T) {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)

	// 建立 mock service
	mockService := new(MockExchangeRateService)
	
	// 設定 mock 行為 - 更新成功但取得記錄失敗
	mockService.On("RefreshTodayRate").Return(nil)
	mockService.On("GetRateRecord", models.CurrencyUSD, models.CurrencyTWD, mock.Anything).Return(nil, errors.New("database error"))

	// 建立 handler
	handler := NewExchangeRateHandler(mockService)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/exchange-rates/refresh", nil)

	// 執行 handler
	handler.RefreshExchangeRate(c)

	// 驗證回應
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "GET_RATE_FAILED", response.Error.Code)
	assert.Contains(t, response.Error.Message, "database error")

	// 驗證 mock 被呼叫
	mockService.AssertExpectations(t)
}

