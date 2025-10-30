package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chienchuanw/asset-manager/internal/client"
	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 載入環境變數
	if err := godotenv.Load(".env.local"); err != nil {
		log.Printf("Warning: .env.local file not found: %v", err)
	}

	// 連接資料庫
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database")

	// 初始化 repositories
	transactionRepo := repository.NewTransactionRepository(db)
	exchangeRateRepo := repository.NewExchangeRateRepository(db)

	// 初始化 services
	exchangeRateAPIClient := client.NewExchangeRateAPIClient()
	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, exchangeRateAPIClient, nil)
	fifoCalculator := service.NewFIFOCalculator(exchangeRateService)

	// 初始化 Price Service (使用 Mock 來測試)
	priceService := service.NewMockPriceService()

	// 初始化 Holding Service
	holdingService := service.NewHoldingService(transactionRepo, fifoCalculator, priceService, exchangeRateService)

	// 測試取得所有持倉
	fmt.Println("\n=== 測試 GetAllHoldings API ===")

	filters := models.HoldingFilters{}
	holdings, err := holdingService.GetAllHoldings(filters)
	if err != nil {
		log.Fatalf("❌ Failed to get holdings: %v", err)
	}

	fmt.Printf("\n找到 %d 個持倉\n", len(holdings))

	for _, holding := range holdings {
		fmt.Println("\n============================================================")
		fmt.Printf("標的: %s (%s)\n", holding.Symbol, holding.Name)
		fmt.Printf("資產類型: %s\n", holding.AssetType)
		fmt.Println("------------------------------------------------------------")
		fmt.Printf("數量: %.8f\n", holding.Quantity)
		fmt.Printf("平均成本: %.2f TWD\n", holding.AvgCost)
		fmt.Printf("總成本: %.2f TWD\n", holding.TotalCost)
		fmt.Println("------------------------------------------------------------")
		fmt.Printf("當前價格: %.2f %s\n", holding.CurrentPrice, holding.Currency)
		fmt.Printf("當前價格 (TWD): %.2f TWD\n", holding.CurrentPriceTWD)
		fmt.Printf("市值: %.2f TWD\n", holding.MarketValue)
		fmt.Println("------------------------------------------------------------")
		fmt.Printf("未實現損益: %.2f TWD\n", holding.UnrealizedPL)
		fmt.Printf("損益百分比: %.2f%%\n", holding.UnrealizedPLPct)
		fmt.Println("============================================================")

		// 驗證計算
		expectedMarketValue := holding.Quantity * holding.CurrentPriceTWD
		expectedPL := expectedMarketValue - holding.TotalCost
		expectedPct := (expectedPL / holding.TotalCost) * 100

		fmt.Println("\n驗證:")
		fmt.Printf("  預期市值: %.2f TWD\n", expectedMarketValue)
		fmt.Printf("  實際市值: %.2f TWD\n", holding.MarketValue)
		fmt.Printf("  預期損益: %.2f TWD\n", expectedPL)
		fmt.Printf("  實際損益: %.2f TWD\n", holding.UnrealizedPL)
		fmt.Printf("  預期百分比: %.2f%%\n", expectedPct)
		fmt.Printf("  實際百分比: %.2f%%\n", holding.UnrealizedPLPct)

		if holding.TotalCost < 1000 && holding.Currency == models.CurrencyUSD {
			fmt.Println("  ❌ 警告：總成本太低，可能未正確轉換！")
		} else {
			fmt.Println("  ✅ 計算正確")
		}
	}
}

