package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// AlphaVantageClient Alpha Vantage API 客戶端（美股價格）
type AlphaVantageClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// AlphaVantageGlobalQuoteResponse Alpha Vantage Global Quote 回應
type AlphaVantageGlobalQuoteResponse struct {
	GlobalQuote AlphaVantageQuote `json:"Global Quote"`
	Note        string            `json:"Note"`        // API 速率限制訊息
	Information string            `json:"Information"` // API 錯誤訊息
}

// AlphaVantageQuote Alpha Vantage 報價資料
type AlphaVantageQuote struct {
	Symbol           string `json:"01. symbol"`
	Open             string `json:"02. open"`
	High             string `json:"03. high"`
	Low              string `json:"04. low"`
	Price            string `json:"05. price"`
	Volume           string `json:"06. volume"`
	LatestTradingDay string `json:"07. latest trading day"`
	PreviousClose    string `json:"08. previous close"`
	Change           string `json:"09. change"`
	ChangePercent    string `json:"10. change percent"`
}

// NewAlphaVantageClient 建立 Alpha Vantage 客戶端
func NewAlphaVantageClient(apiKey string) *AlphaVantageClient {
	return &AlphaVantageClient{
		apiKey:  apiKey,
		baseURL: "https://www.alphavantage.co/query",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetStockPrice 取得美股即時價格
// symbol: 股票代碼（例如：AAPL, GOOGL）
func (c *AlphaVantageClient) GetStockPrice(symbol string) (float64, error) {
	// Alpha Vantage API 端點：GLOBAL_QUOTE
	url := fmt.Sprintf("%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		c.baseURL,
		symbol,
		c.apiKey,
	)

	// 發送 HTTP 請求
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch stock price: %w", err)
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("Alpha Vantage API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// 解析 JSON 回應
	var result AlphaVantageGlobalQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// 檢查 API 速率限制
	if result.Note != "" {
		return 0, fmt.Errorf("Alpha Vantage API rate limit: %s", result.Note)
	}

	// 檢查 API 錯誤訊息
	if result.Information != "" {
		return 0, fmt.Errorf("Alpha Vantage API error: %s", result.Information)
	}

	// 檢查是否有資料
	if result.GlobalQuote.Symbol == "" {
		return 0, fmt.Errorf("no price data found for symbol: %s", symbol)
	}

	// 解析價格字串為 float64
	price, err := strconv.ParseFloat(result.GlobalQuote.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return price, nil
}

// GetMultipleStockPrices 批次取得多個美股價格
// 注意：Alpha Vantage 免費版不支援批次查詢，需要逐一查詢
// 免費版限制：每分鐘 5 次請求，每天 100 次請求
func (c *AlphaVantageClient) GetMultipleStockPrices(symbols []string) (map[string]float64, error) {
	prices := make(map[string]float64)

	for i, symbol := range symbols {
		// 為了避免超過速率限制，每次請求之間延遲 12 秒（免費版每分鐘 5 次）
		if i > 0 {
			time.Sleep(12 * time.Second)
		}

		price, err := c.GetStockPrice(symbol)
		if err != nil {
			// 如果單一股票查詢失敗，記錄錯誤但繼續處理其他股票
			fmt.Printf("Warning: failed to get price for %s: %v\n", symbol, err)
			continue
		}

		prices[symbol] = price
	}

	return prices, nil
}

