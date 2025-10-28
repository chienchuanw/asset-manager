package main

import (
	"fmt"
)

type CryptoHolding struct {
	Symbol              string
	Quantity            float64
	TransactionAmountUSD float64
	CurrentPriceUSD     float64
}

func main() {
	// 交易匯率（2025-10-20）
	transactionRate := 31.5
	// 當前匯率（2025-10-28）
	currentRate := 31.5

	cryptos := []CryptoHolding{
		{
			Symbol:              "BTC",
			Quantity:            0.00215432,
			TransactionAmountUSD: 232.43,
			CurrentPriceUSD:     107888.35, // 假設價格沒變
		},
		{
			Symbol:              "ETH",
			Quantity:            0.05863319,
			TransactionAmountUSD: 230.61,
			CurrentPriceUSD:     3933.02,
		},
		{
			Symbol:              "SOL",
			Quantity:            0.61845231,
			TransactionAmountUSD: 116.78,
			CurrentPriceUSD:     188.82,
		},
		{
			Symbol:              "USDT",
			Quantity:            213.93292890,
			TransactionAmountUSD: 216.07,
			CurrentPriceUSD:     1.00,
		},
	}

	fmt.Println("=== 加密貨幣持倉計算驗證（修正後）===")
	fmt.Println()

	for _, crypto := range cryptos {
		fmt.Printf("【%s】\n", crypto.Symbol)
		fmt.Println("----------------------------------------")

		// 計算 total_cost（交易時的 USD 金額轉換為 TWD）
		totalCostTWD := crypto.TransactionAmountUSD * transactionRate
		fmt.Printf("總成本 (TWD): %.2f USD × %.2f = %.2f TWD\n",
			crypto.TransactionAmountUSD, transactionRate, totalCostTWD)

		// 計算 current_price_twd
		currentPriceTWD := crypto.CurrentPriceUSD * currentRate
		fmt.Printf("當前價格 (TWD): %.2f USD × %.2f = %.2f TWD\n",
			crypto.CurrentPriceUSD, currentRate, currentPriceTWD)

		// 計算 market_value
		marketValue := crypto.Quantity * currentPriceTWD
		fmt.Printf("市值 (TWD): %.8f × %.2f = %.2f TWD\n",
			crypto.Quantity, currentPriceTWD, marketValue)

		// 計算 unrealized_pl
		unrealizedPL := marketValue - totalCostTWD
		fmt.Printf("未實現損益 (TWD): %.2f - %.2f = %.2f TWD\n",
			marketValue, totalCostTWD, unrealizedPL)

		// 計算 unrealized_pl_pct
		unrealizedPLPct := (unrealizedPL / totalCostTWD) * 100
		fmt.Printf("損益百分比: (%.2f / %.2f) × 100 = %.2f%%\n",
			unrealizedPL, totalCostTWD, unrealizedPLPct)

		fmt.Println()
	}

	fmt.Println("=== 總結 ===")
	fmt.Println("修正前：total_cost 使用 USD 金額（錯誤）")
	fmt.Println("修正後：total_cost = USD 金額 × 交易當天匯率（正確）")
	fmt.Println()
	fmt.Println("這樣損益百分比才會合理！")
}

