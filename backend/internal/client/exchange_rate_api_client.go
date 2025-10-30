package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ExchangeRateAPIClient ExchangeRate-API 客戶端
type ExchangeRateAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

// ExchangeRateAPIResponse ExchangeRate-API 的回應格式
type ExchangeRateAPIResponse struct {
	Provider        string             `json:"provider"`
	Base            string             `json:"base"`
	Date            string             `json:"date"`
	TimeLastUpdated int64              `json:"time_last_updated"`
	Rates           map[string]float64 `json:"rates"`
}

// NewExchangeRateAPIClient 建立新的 ExchangeRate-API 客戶端
func NewExchangeRateAPIClient() *ExchangeRateAPIClient {
	return &ExchangeRateAPIClient{
		baseURL: "https://api.exchangerate-api.com/v4/latest",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetUSDToTWDRate 取得 USD 到 TWD 的匯率
func (c *ExchangeRateAPIClient) GetUSDToTWDRate() (float64, error) {
	// 從 USD 為基準取得所有匯率
	resp, err := c.httpClient.Get(c.baseURL + "/USD")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 解析 JSON 回應
	var data ExchangeRateAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// 取得 TWD 匯率
	twdRate, exists := data.Rates["TWD"]
	if !exists {
		return 0, fmt.Errorf("TWD rate not found")
	}

	return twdRate, nil
}

