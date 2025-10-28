package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// 載入環境變數
	if err := godotenv.Load(".env.local"); err != nil {
		log.Printf("Warning: .env.local file not found: %v", err)
	}

	// 取得 API Keys
	finmindAPIKey := os.Getenv("FINMIND_API_KEY")
	coingeckoAPIKey := os.Getenv("COINGECKO_API_KEY")
	alphaVantageAPIKey := os.Getenv("ALPHA_VANTAGE_API_KEY")

	if finmindAPIKey == "" || coingeckoAPIKey == "" || alphaVantageAPIKey == "" {
		log.Fatal("❌ Missing API keys")
	}

	// 初始化真實 Price Service
	priceService := service.NewRealPriceService(finmindAPIKey, coingeckoAPIKey, alphaVantageAPIKey)

	fmt.Println("=== 測試真實價格 API ===")

	// 測試 BTC
	fmt.Println("【BTC】")
	btcPrice, err := priceService.GetPrice("BTC", models.AssetTypeCrypto)
	if err != nil {
		log.Printf("❌ Failed to get BTC price: %v\n", err)
	} else {
		fmt.Printf("  價格: %.2f USD\n", btcPrice.Price)
		fmt.Printf("  更新時間: %s\n", btcPrice.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	// 測試 AAPL
	fmt.Println("\n【AAPL】")
	aaplPrice, err := priceService.GetPrice("AAPL", models.AssetTypeUSStock)
	if err != nil {
		log.Printf("❌ Failed to get AAPL price: %v\n", err)
	} else {
		fmt.Printf("  價格: %.2f USD\n", aaplPrice.Price)
		fmt.Printf("  更新時間: %s\n", aaplPrice.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	// 測試 ETH
	fmt.Println("\n【ETH】")
	ethPrice, err := priceService.GetPrice("ETH", models.AssetTypeCrypto)
	if err != nil {
		log.Printf("❌ Failed to get ETH price: %v\n", err)
	} else {
		fmt.Printf("  價格: %.2f USD\n", ethPrice.Price)
		fmt.Printf("  更新時間: %s\n", ethPrice.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
}

