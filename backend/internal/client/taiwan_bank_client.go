package client

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TaiwanBankClient 台灣銀行匯率 API 客戶端
type TaiwanBankClient struct {
	baseURL    string
	httpClient *http.Client
}

// TaiwanBankRate 台灣銀行匯率資料
type TaiwanBankRate struct {
	Currency   string  // 幣別代碼（例如：USD）
	BuyRate    float64 // 現金買入匯率
	SellRate   float64 // 現金賣出匯率
	SpotBuy    float64 // 即期買入匯率
	SpotSell   float64 // 即期賣出匯率
	UpdateTime string  // 更新時間
}

// NewTaiwanBankClient 建立新的台灣銀行 API 客戶端
func NewTaiwanBankClient() *TaiwanBankClient {
	return &TaiwanBankClient{
		baseURL: "https://rate.bot.com.tw/xrt/flcsv/0/day",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetExchangeRates 取得所有匯率資料
func (c *TaiwanBankClient) GetExchangeRates() (map[string]*TaiwanBankRate, error) {
	// 發送 HTTP 請求
	resp, err := c.httpClient.Get(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 讀取 CSV 資料
	reader := csv.NewReader(resp.Body)
	rates := make(map[string]*TaiwanBankRate)

	// 跳過標題行
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// 解析每一行資料
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		// CSV 格式：日期,幣別,現金買入,現金賣出,即期買入,即期賣出
		// 範例：2025/10/24,USD,30.5,31.5,30.8,31.2
		if len(record) < 6 {
			continue
		}

		updateTime := strings.TrimSpace(record[0])
		currency := strings.TrimSpace(record[1])

		// 解析匯率（移除可能的逗號）
		buyRate, _ := parseFloat(record[2])
		sellRate, _ := parseFloat(record[3])
		spotBuy, _ := parseFloat(record[4])
		spotSell, _ := parseFloat(record[5])

		rates[currency] = &TaiwanBankRate{
			Currency:   currency,
			BuyRate:    buyRate,
			SellRate:   sellRate,
			SpotBuy:    spotBuy,
			SpotSell:   spotSell,
			UpdateTime: updateTime,
		}
	}

	return rates, nil
}

// GetUSDToTWDRate 取得 USD 到 TWD 的匯率（使用即期賣出匯率）
func (c *TaiwanBankClient) GetUSDToTWDRate() (float64, error) {
	rates, err := c.GetExchangeRates()
	if err != nil {
		return 0, err
	}

	usdRate, exists := rates["USD"]
	if !exists {
		return 0, fmt.Errorf("USD rate not found")
	}

	// 使用即期賣出匯率（銀行賣出 USD 給客戶的匯率）
	// 這是最常用的匯率，代表客戶用 TWD 買 USD 的價格
	if usdRate.SpotSell > 0 {
		return usdRate.SpotSell, nil
	}

	// 如果即期賣出匯率不存在，使用現金賣出匯率
	if usdRate.SellRate > 0 {
		return usdRate.SellRate, nil
	}

	return 0, fmt.Errorf("no valid USD rate found")
}

// parseFloat 解析浮點數（處理逗號分隔符）
func parseFloat(s string) (float64, error) {
	// 移除逗號和空白
	s = strings.ReplaceAll(strings.TrimSpace(s), ",", "")
	if s == "" || s == "-" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

