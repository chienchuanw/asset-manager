package main

import (
	"log"
	"os"
	"strings"

	"github.com/chienchuanw/asset-manager/internal/client"
	"github.com/chienchuanw/asset-manager/internal/db"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	// 載入環境變數
	if err := godotenv.Load(".env.local"); err != nil {
		log.Printf("Warning: .env.local file not found, using environment variables")
	}

	// 連接資料庫
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("✓ Database connected")

	// 初始化 repositories
	transactionRepo := repository.NewTransactionRepository(database)
	exchangeRateRepo := repository.NewExchangeRateRepository(database)
	assetSnapshotRepo := repository.NewAssetSnapshotRepository(database)
	realizedProfitRepo := repository.NewRealizedProfitRepository(database)

	// PerformanceSnapshotRepository 需要 sqlx.DB
	dbx := sqlx.NewDb(database, "postgres")
	performanceSnapshotRepo := repository.NewPerformanceSnapshotRepository(dbx)

	// 初始化 Price Service
	var priceService service.PriceService
	finmindAPIKey := os.Getenv("FINMIND_API_KEY")
	coingeckoAPIKey := os.Getenv("COINGECKO_API_KEY")
	alphaVantageAPIKey := os.Getenv("ALPHA_VANTAGE_API_KEY")

	if finmindAPIKey != "" && coingeckoAPIKey != "" && alphaVantageAPIKey != "" {
		priceService = service.NewRealPriceService(finmindAPIKey, coingeckoAPIKey, alphaVantageAPIKey)
		log.Println("✓ Using real price API (FinMind + CoinGecko + Alpha Vantage)")
	} else {
		priceService = service.NewMockPriceService()
		log.Println("⚠️  Using mock price service (API keys not configured)")
	}

	// 初始化匯率服務
	exchangeRateClient := client.NewExchangeRateAPIClient()
	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, exchangeRateClient, nil)

	// 初始化 FIFO Calculator
	fifoCalculator := service.NewFIFOCalculator(exchangeRateService)

	// 初始化 HoldingService
	holdingService := service.NewHoldingService(transactionRepo, fifoCalculator, priceService, exchangeRateService)

	// 初始化 AssetSnapshotService
	assetSnapshotService := service.NewAssetSnapshotServiceWithDeps(assetSnapshotRepo, holdingService)

	// 初始化 Analytics Services
	unrealizedAnalyticsService := service.NewUnrealizedAnalyticsService(holdingService)
	analyticsService := service.NewAnalyticsService(realizedProfitRepo)
	performanceTrendService := service.NewPerformanceTrendService(performanceSnapshotRepo, unrealizedAnalyticsService, analyticsService)

	log.Println("✓ Services initialized")

	// 1. 更新今日匯率
	log.Println("\n📊 Step 1: Refreshing today's exchange rate...")
	if err := exchangeRateService.RefreshTodayRate(); err != nil {
		log.Printf("⚠️  Warning: Failed to refresh exchange rate: %v", err)
		log.Println("   Continuing with cached/default rate...")
	} else {
		log.Println("✓ Exchange rate refreshed successfully")
	}

	// 2. 建立資產快照（asset_snapshots）
	log.Println("\n📊 Step 2: Creating asset snapshots...")
	if err := assetSnapshotService.CreateDailySnapshots(); err != nil {
		log.Fatalf("❌ Failed to create asset snapshots: %v", err)
	}
	log.Println("✓ Asset snapshots created successfully")

	// 3. 建立績效快照（daily_performance_snapshots）
	log.Println("\n📊 Step 3: Creating performance snapshot...")
	snapshot, err := performanceTrendService.CreateDailySnapshot()
	if err != nil {
		log.Fatalf("❌ Failed to create performance snapshot: %v", err)
	}
	log.Println("✓ Performance snapshot created successfully")

	// 顯示摘要
	log.Println("\n" + strings.Repeat("=", 60))
	log.Println("📈 Snapshot Summary")
	log.Println(strings.Repeat("=", 60))
	log.Printf("Date:              %s\n", snapshot.SnapshotDate.Format("2006-01-02"))
	log.Printf("Total Market Value: %.2f TWD\n", snapshot.TotalMarketValue)
	log.Printf("Total Cost:         %.2f TWD\n", snapshot.TotalCost)
	log.Printf("Unrealized P/L:     %.2f TWD (%.2f%%)\n", snapshot.TotalUnrealizedPL, snapshot.TotalUnrealizedPct)
	log.Printf("Realized P/L:       %.2f TWD (%.2f%%)\n", snapshot.TotalRealizedPL, snapshot.TotalRealizedPct)
	log.Printf("Holdings Count:     %d\n", snapshot.HoldingCount)
	log.Println(strings.Repeat("=", 60))

	log.Println("\n✅ All snapshots created successfully!")
	os.Exit(0)
}

