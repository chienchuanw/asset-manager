package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("=== 最終驗證：預期損益百分比 ===")

	// 匯率
	exchangeRate := 31.5

	// BTC
	fmt.Println("【BTC】")
	btcQuantity := 0.00215432
	btcTransactionAmount := 232.43 // USD
	btcCurrentPrice := 115196.00   // USD

	btcTotalCost := btcTransactionAmount * exchangeRate
	btcCurrentPriceTWD := btcCurrentPrice * exchangeRate
	btcMarketValue := btcQuantity * btcCurrentPriceTWD
	btcUnrealizedPL := btcMarketValue - btcTotalCost
	btcUnrealizedPLPct := (btcUnrealizedPL / btcTotalCost) * 100

	fmt.Printf("  交易金額: %.2f USD\n", btcTransactionAmount)
	fmt.Printf("  總成本: %.2f TWD (%.2f × %.1f)\n", btcTotalCost, btcTransactionAmount, exchangeRate)
	fmt.Printf("  當前價格: %.2f USD\n", btcCurrentPrice)
	fmt.Printf("  當前價格 TWD: %.2f TWD (%.2f × %.1f)\n", btcCurrentPriceTWD, btcCurrentPrice, exchangeRate)
	fmt.Printf("  市值: %.2f TWD (%.8f × %.2f)\n", btcMarketValue, btcQuantity, btcCurrentPriceTWD)
	fmt.Printf("  未實現損益: %.2f TWD\n", btcUnrealizedPL)
	fmt.Printf("  損益百分比: %.2f%%\n", btcUnrealizedPLPct)

	// AAPL
	fmt.Println("\n【AAPL】")
	aaplQuantity := 1.69698
	aaplTransactionAmount := 440.50 // USD
	aaplCurrentPrice := 268.81      // USD

	aaplTotalCost := aaplTransactionAmount * exchangeRate
	aaplCurrentPriceTWD := aaplCurrentPrice * exchangeRate
	aaplMarketValue := aaplQuantity * aaplCurrentPriceTWD
	aaplUnrealizedPL := aaplMarketValue - aaplTotalCost
	aaplUnrealizedPLPct := (aaplUnrealizedPL / aaplTotalCost) * 100

	fmt.Printf("  交易金額: %.2f USD\n", aaplTransactionAmount)
	fmt.Printf("  總成本: %.2f TWD (%.2f × %.1f)\n", aaplTotalCost, aaplTransactionAmount, exchangeRate)
	fmt.Printf("  當前價格: %.2f USD\n", aaplCurrentPrice)
	fmt.Printf("  當前價格 TWD: %.2f TWD (%.2f × %.1f)\n", aaplCurrentPriceTWD, aaplCurrentPrice, exchangeRate)
	fmt.Printf("  市值: %.2f TWD (%.8f × %.2f)\n", aaplMarketValue, aaplQuantity, aaplCurrentPriceTWD)
	fmt.Printf("  未實現損益: %.2f TWD\n", aaplUnrealizedPL)
	fmt.Printf("  損益百分比: %.2f%%\n", aaplUnrealizedPLPct)

	// ETH
	fmt.Println("\n【ETH】")
	ethQuantity := 0.05863319
	ethTransactionAmount := 230.61 // USD
	ethCurrentPrice := 4128.57     // USD

	ethTotalCost := ethTransactionAmount * exchangeRate
	ethCurrentPriceTWD := ethCurrentPrice * exchangeRate
	ethMarketValue := ethQuantity * ethCurrentPriceTWD
	ethUnrealizedPL := ethMarketValue - ethTotalCost
	ethUnrealizedPLPct := (ethUnrealizedPL / ethTotalCost) * 100

	fmt.Printf("  交易金額: %.2f USD\n", ethTransactionAmount)
	fmt.Printf("  總成本: %.2f TWD (%.2f × %.1f)\n", ethTotalCost, ethTransactionAmount, exchangeRate)
	fmt.Printf("  當前價格: %.2f USD\n", ethCurrentPrice)
	fmt.Printf("  當前價格 TWD: %.2f TWD (%.2f × %.1f)\n", ethCurrentPriceTWD, ethCurrentPrice, exchangeRate)
	fmt.Printf("  市值: %.2f TWD (%.8f × %.2f)\n", ethMarketValue, ethQuantity, ethCurrentPriceTWD)
	fmt.Printf("  未實現損益: %.2f TWD\n", ethUnrealizedPL)
	fmt.Printf("  損益百分比: %.2f%%\n", ethUnrealizedPLPct)

	fmt.Println("\n=== 總結 ===")
	fmt.Println("✅ 所有損益百分比都應該在合理範圍內（不是 3000%+）")
	fmt.Println("✅ BTC 應該是正報酬（買入價 107,888 USD，現價 115,196 USD）")
	fmt.Println("✅ AAPL 應該是正報酬（買入價 259.58 USD，現價 268.81 USD）")
	fmt.Println("✅ ETH 應該是正報酬（買入價 3,933.02 USD，現價 4,128.57 USD）")

	// 檢查是否合理
	if btcUnrealizedPLPct > 100 || btcUnrealizedPLPct < -50 {
		log.Printf("⚠️  BTC 損益百分比不合理: %.2f%%", btcUnrealizedPLPct)
	}
	if aaplUnrealizedPLPct > 100 || aaplUnrealizedPLPct < -50 {
		log.Printf("⚠️  AAPL 損益百分比不合理: %.2f%%", aaplUnrealizedPLPct)
	}
	if ethUnrealizedPLPct > 100 || ethUnrealizedPLPct < -50 {
		log.Printf("⚠️  ETH 損益百分比不合理: %.2f%%", ethUnrealizedPLPct)
	}
}

