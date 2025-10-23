package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FinMindClient FinMind API 客戶端（台股價格）
type FinMindClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// FinMindResponse FinMind API 回應結構
type FinMindResponse struct {
	Msg    string          `json:"msg"`
	Status int             `json:"status"`
	Data   []FinMindPrice  `json:"data"`
}

// FinMindPrice FinMind 價格資料
type FinMindPrice struct {
	Date        string  `json:"date"`
	StockID     string  `json:"stock_id"`
	Open        float64 `json:"open"`
	Max         float64 `json:"max"`
	Min         float64 `json:"min"`
	Close       float64 `json:"close"`
	SpreadRate  float64 `json:"spread_rate"`
	Volume      int64   `json:"Trading_Volume"`
	Transaction int64   `json:"Trading_turnover"`
	Change      float64 `json:"change"`
}

// NewFinMindClient 建立 FinMind 客戶端
func NewFinMindClient(apiKey string) *FinMindClient {
	return &FinMindClient{
		apiKey:  apiKey,
		baseURL: "https://api.finmindtrade.com/api/v4",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetStockPrice 取得台股即時價格
// symbol: 股票代碼（例如：2330）
func (c *FinMindClient) GetStockPrice(symbol string) (float64, error) {
	// FinMind API 端點：取得最新收盤價
	url := fmt.Sprintf("%s/data?dataset=TaiwanStockPrice&data_id=%s&start_date=%s&token=%s",
		c.baseURL,
		symbol,
		c.getRecentDate(), // 取得最近 7 天的日期
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
		return 0, fmt.Errorf("FinMind API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// 解析 JSON 回應
	var result FinMindResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// 檢查 API 回應狀態
	if result.Status != 200 {
		return 0, fmt.Errorf("FinMind API error: status=%d, msg=%s", result.Status, result.Msg)
	}

	// 檢查是否有資料
	if len(result.Data) == 0 {
		return 0, fmt.Errorf("no price data found for symbol: %s", symbol)
	}

	// 返回最新的收盤價（最後一筆資料）
	latestPrice := result.Data[len(result.Data)-1]
	return latestPrice.Close, nil
}

// getRecentDate 取得最近 7 天的日期（YYYY-MM-DD 格式）
// 用於查詢最新的股價資料
func (c *FinMindClient) getRecentDate() string {
	// 往前推 7 天，確保能取得最新資料（考慮週末和假日）
	date := time.Now().AddDate(0, 0, -7)
	return date.Format("2006-01-02")
}

// GetMultipleStockPrices 批次取得多個台股價格
func (c *FinMindClient) GetMultipleStockPrices(symbols []string) (map[string]float64, error) {
	prices := make(map[string]float64)
	
	for _, symbol := range symbols {
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

