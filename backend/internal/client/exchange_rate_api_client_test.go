package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGetUSDToTWDRate_Success 測試成功取得 USD 到 TWD 的匯率
func TestGetUSDToTWDRate_Success(t *testing.T) {
	// Arrange - 建立模擬的 HTTP 伺服器
	mockResponse := ExchangeRateAPIResponse{
		Provider: "https://www.exchangerate-api.com",
		Base:     "USD",
		Date:     "2025-10-30",
		Rates: map[string]float64{
			"USD": 1.0,
			"TWD": 30.6,
			"EUR": 0.86,
			"JPY": 152.42,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 驗證請求路徑
		assert.Equal(t, "/USD", r.URL.Path)
		
		// 回傳模擬的 JSON 資料
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// 建立 client，使用測試伺服器的 URL
	client := &ExchangeRateAPIClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	// Act - 呼叫方法
	rate, err := client.GetUSDToTWDRate()

	// Assert - 驗證結果
	assert.NoError(t, err)
	assert.Equal(t, 30.6, rate)
}

// TestGetUSDToTWDRate_HTTPError 測試 HTTP 請求失敗
func TestGetUSDToTWDRate_HTTPError(t *testing.T) {
	// Arrange - 建立會回傳錯誤的伺服器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &ExchangeRateAPIClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	// Act
	rate, err := client.GetUSDToTWDRate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0.0, rate)
	assert.Contains(t, err.Error(), "unexpected status code")
}

// TestGetUSDToTWDRate_InvalidJSON 測試無效的 JSON 回應
func TestGetUSDToTWDRate_InvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := &ExchangeRateAPIClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	// Act
	rate, err := client.GetUSDToTWDRate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0.0, rate)
	assert.Contains(t, err.Error(), "failed to decode JSON")
}

// TestGetUSDToTWDRate_MissingTWDRate 測試回應中缺少 TWD 匯率
func TestGetUSDToTWDRate_MissingTWDRate(t *testing.T) {
	// Arrange
	mockResponse := ExchangeRateAPIResponse{
		Provider: "https://www.exchangerate-api.com",
		Base:     "USD",
		Date:     "2025-10-30",
		Rates: map[string]float64{
			"USD": 1.0,
			"EUR": 0.86,
			// 故意不包含 TWD
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := &ExchangeRateAPIClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	// Act
	rate, err := client.GetUSDToTWDRate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0.0, rate)
	assert.Contains(t, err.Error(), "TWD rate not found")
}

// TestGetUSDToTWDRate_NetworkError 測試網路錯誤
func TestGetUSDToTWDRate_NetworkError(t *testing.T) {
	// Arrange - 使用無效的 URL
	client := &ExchangeRateAPIClient{
		baseURL:    "http://invalid-url-that-does-not-exist-12345.com",
		httpClient: &http.Client{Timeout: 1 * time.Second},
	}

	// Act
	rate, err := client.GetUSDToTWDRate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0.0, rate)
	// 網路錯誤會包含 "failed to fetch exchange rates"
}

// TestNewExchangeRateAPIClient 測試建立新的 client
func TestNewExchangeRateAPIClient(t *testing.T) {
	// Act
	client := NewExchangeRateAPIClient()

	// Assert
	assert.NotNil(t, client)
	assert.Equal(t, "https://api.exchangerate-api.com/v4/latest", client.baseURL)
	assert.NotNil(t, client.httpClient)
}

