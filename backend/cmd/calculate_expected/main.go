package main

import (
	"fmt"
)

func main() {
	// AAPL 範例數據
	quantity := 1.69698
	transactionAmountUSD := 440.50 // 交易時的 USD 金額
	transactionExchangeRate := 31.5 // 交易當天的匯率（假設）
	currentPriceUSD := 268.81
	currentExchangeRate := 31.5

	fmt.Println("=== AAPL 持倉計算驗證（修正後）===")

	// 計算 total_cost（交易時的 USD 金額轉換為 TWD）
	totalCostTWD := transactionAmountUSD * transactionExchangeRate
	fmt.Printf("1. Total Cost (TWD):\n")
	fmt.Printf("   %.2f USD × %.2f = %.2f TWD\n\n", transactionAmountUSD, transactionExchangeRate, totalCostTWD)

	// 計算 current_price_twd
	currentPriceTWD := currentPriceUSD * currentExchangeRate
	fmt.Printf("2. Current Price (TWD):\n")
	fmt.Printf("   %.2f USD × %.2f = %.2f TWD\n\n", currentPriceUSD, currentExchangeRate, currentPriceTWD)

	// 計算 market_value
	marketValue := quantity * currentPriceTWD
	fmt.Printf("3. Market Value (TWD):\n")
	fmt.Printf("   %.5f × %.2f = %.2f TWD\n\n", quantity, currentPriceTWD, marketValue)

	// 計算 unrealized_pl
	unrealizedPL := marketValue - totalCostTWD
	fmt.Printf("4. Unrealized P/L (TWD):\n")
	fmt.Printf("   %.2f - %.2f = %.2f TWD\n\n", marketValue, totalCostTWD, unrealizedPL)

	// 計算 unrealized_pl_pct
	unrealizedPLPct := (unrealizedPL / totalCostTWD) * 100
	fmt.Printf("5. Unrealized P/L %%:\n")
	fmt.Printf("   (%.2f / %.2f) × 100 = %.2f%%\n\n", unrealizedPL, totalCostTWD, unrealizedPLPct)

	fmt.Println("=== 修正後預期 API Response ===")
	fmt.Printf(`{
  "symbol": "AAPL",
  "quantity": %.5f,
  "total_cost": %.2f,          // 修正：440.50 USD × 31.5 = 13,875.75 TWD
  "current_price": %.2f,
  "currency": "USD",
  "current_price_twd": %.2f,
  "market_value": %.2f,
  "unrealized_pl": %.2f,       // 修正：14,369.20 - 13,875.75 = 493.45 TWD
  "unrealized_pl_pct": %.2f    // 修正：3.56%% (合理！)
}
`, quantity, totalCostTWD, currentPriceUSD, currentPriceTWD, marketValue, unrealizedPL, unrealizedPLPct)
}

