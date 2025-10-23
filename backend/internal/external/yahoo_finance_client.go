package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// YahooFinanceClient Yahoo Finance API 客戶端（美股價格）
// 使用免費的 Yahoo Finance API v8
type YahooFinanceClient struct {
	baseURL    string
	httpClient *http.Client
}

// YahooQuoteResponse Yahoo Finance Quote 回應
type YahooQuoteResponse struct {
	QuoteResponse struct {
		Result []YahooQuote `json:"result"`
		Error  interface{}  `json:"error"`
	} `json:"quoteResponse"`
}

// YahooQuote Yahoo Finance 報價資料
type YahooQuote struct {
	Symbol                string  `json:"symbol"`
	RegularMarketPrice    float64 `json:"regularMarketPrice"`
	RegularMarketTime     int64   `json:"regularMarketTime"`
	RegularMarketChange   float64 `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketDayHigh  float64 `json:"regularMarketDayHigh"`
	RegularMarketDayLow   float64 `json:"regularMarketDayLow"`
	RegularMarketVolume   int64   `json:"regularMarketVolume"`
	RegularMarketOpen     float64 `json:"regularMarketOpen"`
	RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
}

// NewYahooFinanceClient 建立 Yahoo Finance 客戶端
func NewYahooFinanceClient() *YahooFinanceClient {
	return &YahooFinanceClient{
		baseURL: "https://query1.finance.yahoo.com/v8/finance",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetStockPrice 取得美股即時價格
// symbol: 股票代碼（例如：AAPL, GOOGL）
func (c *YahooFinanceClient) GetStockPrice(symbol string) (float64, error) {
	// Yahoo Finance API 端點：quote
	url := fmt.Sprintf("%s/quote?symbols=%s", c.baseURL, symbol)

	// 建立 HTTP 請求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// 設定 User-Agent（Yahoo Finance 需要）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	// 發送 HTTP 請求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch stock price: %w", err)
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("Yahoo Finance API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// 解析 JSON 回應
	var result YahooQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// 檢查是否有錯誤
	if result.QuoteResponse.Error != nil {
		return 0, fmt.Errorf("Yahoo Finance API error: %v", result.QuoteResponse.Error)
	}

	// 檢查是否有資料
	if len(result.QuoteResponse.Result) == 0 {
		return 0, fmt.Errorf("no price data found for symbol: %s", symbol)
	}

	// 返回即時市場價格
	quote := result.QuoteResponse.Result[0]
	return quote.RegularMarketPrice, nil
}

// GetMultipleStockPrices 批次取得多個美股價格
// symbols: 股票代碼列表（例如：["AAPL", "GOOGL", "MSFT"]）
func (c *YahooFinanceClient) GetMultipleStockPrices(symbols []string) (map[string]float64, error) {
	if len(symbols) == 0 {
		return make(map[string]float64), nil
	}

	// Yahoo Finance 支援批次查詢（用逗號分隔）
	symbolsStr := ""
	for i, symbol := range symbols {
		if i > 0 {
			symbolsStr += ","
		}
		symbolsStr += symbol
	}

	url := fmt.Sprintf("%s/quote?symbols=%s", c.baseURL, symbolsStr)

	// 建立 HTTP 請求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 設定 User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	// 發送 HTTP 請求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock prices: %w", err)
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Yahoo Finance API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// 解析 JSON 回應
	var result YahooQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 檢查是否有錯誤
	if result.QuoteResponse.Error != nil {
		return nil, fmt.Errorf("Yahoo Finance API error: %v", result.QuoteResponse.Error)
	}

	// 將結果轉換為 map
	prices := make(map[string]float64)
	for _, quote := range result.QuoteResponse.Result {
		prices[quote.Symbol] = quote.RegularMarketPrice
	}

	return prices, nil
}

