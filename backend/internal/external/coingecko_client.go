package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// CoinGeckoClient CoinGecko API 客戶端（加密貨幣價格）
type CoinGeckoClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// CoinGeckoPriceResponse CoinGecko 價格回應
type CoinGeckoPriceResponse struct {
	// 格式：{ "bitcoin": { "twd": 1234567.89, "usd": 43210.12 } }
	Prices map[string]map[string]float64
}

// NewCoinGeckoClient 建立 CoinGecko 客戶端
func NewCoinGeckoClient(apiKey string) *CoinGeckoClient {
	return &CoinGeckoClient{
		apiKey:  apiKey,
		baseURL: "https://api.coingecko.com/api/v3",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// symbolToCoinID 將常見的加密貨幣代碼轉換為 CoinGecko ID
func (c *CoinGeckoClient) symbolToCoinID(symbol string) string {
	// CoinGecko 使用小寫的完整名稱作為 ID
	symbolMap := map[string]string{
		"BTC":  "bitcoin",
		"ETH":  "ethereum",
		"USDT": "tether",
		"USDC": "usd-coin",
		"BNB":  "binancecoin",
		"XRP":  "ripple",
		"ADA":  "cardano",
		"DOGE": "dogecoin",
		"SOL":  "solana",
		"MATIC": "matic-network",
		"DOT":  "polkadot",
		"AVAX": "avalanche-2",
	}

	if coinID, exists := symbolMap[strings.ToUpper(symbol)]; exists {
		return coinID
	}

	// 如果找不到對應，返回小寫的 symbol（可能是完整名稱）
	return strings.ToLower(symbol)
}

// GetCryptoPrice 取得加密貨幣價格
// symbol: 加密貨幣代碼（例如：BTC, ETH）
// currency: 目標貨幣（twd 或 usd）
func (c *CoinGeckoClient) GetCryptoPrice(symbol string, currency string) (float64, error) {
	coinID := c.symbolToCoinID(symbol)
	currency = strings.ToLower(currency)

	// CoinGecko API 端點：simple/price
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s&x_cg_demo_api_key=%s",
		c.baseURL,
		coinID,
		currency,
		c.apiKey,
	)

	// 發送 HTTP 請求
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch crypto price: %w", err)
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("CoinGecko API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// 解析 JSON 回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析為 map[string]map[string]float64
	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// 檢查是否有資料
	priceData, exists := result[coinID]
	if !exists {
		return 0, fmt.Errorf("no price data found for coin: %s", symbol)
	}

	price, exists := priceData[currency]
	if !exists {
		return 0, fmt.Errorf("no price data found for currency: %s", currency)
	}

	return price, nil
}

// GetMultipleCryptoPrices 批次取得多個加密貨幣價格
// symbols: 加密貨幣代碼列表
// currency: 目標貨幣（twd 或 usd）
func (c *CoinGeckoClient) GetMultipleCryptoPrices(symbols []string, currency string) (map[string]float64, error) {
	if len(symbols) == 0 {
		return make(map[string]float64), nil
	}

	// 將 symbols 轉換為 CoinGecko IDs
	coinIDs := make([]string, 0, len(symbols))
	symbolToCoinIDMap := make(map[string]string)
	
	for _, symbol := range symbols {
		coinID := c.symbolToCoinID(symbol)
		coinIDs = append(coinIDs, coinID)
		symbolToCoinIDMap[coinID] = symbol
	}

	currency = strings.ToLower(currency)

	// CoinGecko 支援批次查詢（用逗號分隔）
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s&x_cg_demo_api_key=%s",
		c.baseURL,
		strings.Join(coinIDs, ","),
		currency,
		c.apiKey,
	)

	// 發送 HTTP 請求
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch crypto prices: %w", err)
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CoinGecko API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// 解析 JSON 回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 將結果轉換回原始 symbol
	prices := make(map[string]float64)
	for coinID, priceData := range result {
		symbol := symbolToCoinIDMap[coinID]
		if price, exists := priceData[currency]; exists {
			prices[symbol] = price
		}
	}

	return prices, nil
}

