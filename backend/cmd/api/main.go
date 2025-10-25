package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/chienchuanw/asset-manager/internal/api"
	"github.com/chienchuanw/asset-manager/internal/cache"
	"github.com/chienchuanw/asset-manager/internal/client"
	"github.com/chienchuanw/asset-manager/internal/db"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/scheduler"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 載入環境變數
	if err := godotenv.Load(".env.local"); err != nil {
		log.Printf("Warning: .env.local file not found: %v", err)
	}
	
	// 初始化資料庫連線
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 初始化 Repository
	transactionRepo := repository.NewTransactionRepository(database)
	exchangeRateRepo := repository.NewExchangeRateRepository(database)
	realizedProfitRepo := repository.NewRealizedProfitRepository(database)
	assetSnapshotRepo := repository.NewAssetSnapshotRepository(database)

	// 初始化 FIFO Calculator（需要在 TransactionService 之前初始化）
	fifoCalculator := service.NewFIFOCalculator()

	// 初始化 Service
	transactionService := service.NewTransactionService(transactionRepo, realizedProfitRepo, fifoCalculator)

	// 初始化 Redis Cache
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	redisCache, err := cache.NewRedisCache(redisAddr, redisPassword, redisDB)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Using price service without cache.", err)
		// 如果 Redis 連線失敗，使用不帶快取的 Price Service

		// 初始化 Price Service（真實 API 或 Mock）
		var priceService service.PriceService
		finmindAPIKey := os.Getenv("FINMIND_API_KEY")
		coingeckoAPIKey := os.Getenv("COINGECKO_API_KEY")
		alphaVantageAPIKey := os.Getenv("ALPHA_VANTAGE_API_KEY")

		if finmindAPIKey != "" && coingeckoAPIKey != "" && alphaVantageAPIKey != "" {
			priceService = service.NewRealPriceService(finmindAPIKey, coingeckoAPIKey, alphaVantageAPIKey)
			log.Println("Using real price API without cache (FinMind + CoinGecko + Alpha Vantage)")
		} else {
			priceService = service.NewMockPriceService()
			log.Println("Using mock price service without cache")
		}

		// 初始化匯率服務（不帶 Redis 快取）
		bankClient := client.NewTaiwanBankClient()
		exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, bankClient, nil)

		holdingService := service.NewHoldingService(transactionRepo, fifoCalculator, priceService, exchangeRateService)

		// 初始化 Analytics Service
		analyticsService := service.NewAnalyticsService(realizedProfitRepo)
		unrealizedAnalyticsService := service.NewUnrealizedAnalyticsService(holdingService)
		allocationService := service.NewAllocationService(holdingService)

		// 初始化 Asset Snapshot Service（不帶排程器）
		assetSnapshotService := service.NewAssetSnapshotServiceWithDeps(assetSnapshotRepo, holdingService)

		// 初始化 Handler
		transactionHandler := api.NewTransactionHandler(transactionService)
		holdingHandler := api.NewHoldingHandler(holdingService)
		analyticsHandler := api.NewAnalyticsHandler(analyticsService)
		unrealizedAnalyticsHandler := api.NewUnrealizedAnalyticsHandler(unrealizedAnalyticsService)
		allocationHandler := api.NewAllocationHandler(allocationService)
		assetSnapshotHandler := api.NewAssetSnapshotHandler(assetSnapshotService)

		// 建立 router 並啟動（簡化版，不啟動排程器）
		log.Println("Warning: Snapshot scheduler is disabled (Redis not available)")
		startServer(transactionHandler, holdingHandler, analyticsHandler, unrealizedAnalyticsHandler, allocationHandler, assetSnapshotHandler)
		return
	}
	defer redisCache.Close()

	// 解析快取過期時間
	cacheExpiration := 5 * time.Minute
	if expStr := os.Getenv("PRICE_CACHE_EXPIRATION"); expStr != "" {
		if duration, err := time.ParseDuration(expStr); err == nil {
			cacheExpiration = duration
		}
	}

	// 初始化 Price Service（真實 API 或 Mock）
	var basePriceService service.PriceService

	finmindAPIKey := os.Getenv("FINMIND_API_KEY")
	coingeckoAPIKey := os.Getenv("COINGECKO_API_KEY")
	alphaVantageAPIKey := os.Getenv("ALPHA_VANTAGE_API_KEY")

	if finmindAPIKey != "" && coingeckoAPIKey != "" && alphaVantageAPIKey != "" {
		// 使用真實 API
		basePriceService = service.NewRealPriceService(finmindAPIKey, coingeckoAPIKey, alphaVantageAPIKey)
		log.Println("Using real price API (FinMind + CoinGecko + Alpha Vantage)")
	} else {
		// 使用 Mock Service
		basePriceService = service.NewMockPriceService()
		log.Println("Warning: API keys not found. Using mock price service.")
	}

	// 加上 Redis 快取層
	priceService := service.NewCachedPriceService(redisCache, basePriceService, cacheExpiration)

	log.Printf("Redis cache enabled: default=%v, US stocks=1h (to avoid Alpha Vantage API limits)", cacheExpiration)

	// 初始化匯率服務（帶 Redis 快取）
	bankClient := client.NewTaiwanBankClient()
	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, bankClient, redisCache.GetClient())

	// 初始化 Holding Service
	holdingService := service.NewHoldingService(transactionRepo, fifoCalculator, priceService, exchangeRateService)

	// 初始化 Analytics Service
	analyticsService := service.NewAnalyticsService(realizedProfitRepo)
	unrealizedAnalyticsService := service.NewUnrealizedAnalyticsService(holdingService)
	allocationService := service.NewAllocationService(holdingService)

	// 初始化 Asset Snapshot Service（包含依賴）
	assetSnapshotService := service.NewAssetSnapshotServiceWithDeps(assetSnapshotRepo, holdingService)

	// 初始化 Handler
	transactionHandler := api.NewTransactionHandler(transactionService)
	holdingHandler := api.NewHoldingHandler(holdingService)
	analyticsHandler := api.NewAnalyticsHandler(analyticsService)
	unrealizedAnalyticsHandler := api.NewUnrealizedAnalyticsHandler(unrealizedAnalyticsService)
	allocationHandler := api.NewAllocationHandler(allocationService)
	assetSnapshotHandler := api.NewAssetSnapshotHandler(assetSnapshotService)

	// 初始化並啟動排程器
	snapshotSchedulerConfig := scheduler.SnapshotSchedulerConfig{
		Enabled:           os.Getenv("SNAPSHOT_SCHEDULER_ENABLED") == "true",
		DailySnapshotTime: getEnvOrDefault("SNAPSHOT_SCHEDULER_TIME", "23:59"),
	}
	snapshotScheduler := scheduler.NewSnapshotScheduler(assetSnapshotService, snapshotSchedulerConfig)
	if err := snapshotScheduler.Start(); err != nil {
		log.Printf("Warning: Failed to start snapshot scheduler: %v", err)
	}
	defer snapshotScheduler.Stop()

	// 啟動伺服器
	startServer(transactionHandler, holdingHandler, analyticsHandler, unrealizedAnalyticsHandler, allocationHandler, assetSnapshotHandler)
}

// getEnvOrDefault 取得環境變數，如果不存在則使用預設值
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// startServer 啟動 HTTP 伺服器
func startServer(transactionHandler *api.TransactionHandler, holdingHandler *api.HoldingHandler, analyticsHandler *api.AnalyticsHandler, unrealizedAnalyticsHandler *api.UnrealizedAnalyticsHandler, allocationHandler *api.AllocationHandler, assetSnapshotHandler *api.AssetSnapshotHandler) {
	// 建立 Gin router
	router := gin.Default()

	// 設定 CORS
	router.Use(cors.Default())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Asset Manager API Server is running.",
		})
	})

	// API routes
	apiGroup := router.Group("/api")
	{
		// Transactions 路由
		transactions := apiGroup.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.ListTransactions)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
		}

		// Holdings 路由
		holdings := apiGroup.Group("/holdings")
		{
			holdings.GET("", holdingHandler.GetAllHoldings)
			holdings.GET("/:symbol", holdingHandler.GetHoldingBySymbol)
		}

		// Analytics 路由
		analytics := apiGroup.Group("/analytics")
		{
			analytics.GET("/summary", analyticsHandler.GetSummary)
			analytics.GET("/performance", analyticsHandler.GetPerformance)
			analytics.GET("/top-assets", analyticsHandler.GetTopAssets)

			// Unrealized Analytics 路由
			unrealized := analytics.Group("/unrealized")
			{
				unrealized.GET("/summary", unrealizedAnalyticsHandler.GetSummary)
				unrealized.GET("/performance", unrealizedAnalyticsHandler.GetPerformance)
				unrealized.GET("/top-assets", unrealizedAnalyticsHandler.GetTopAssets)
			}
		}

		// Allocation 路由
		allocation := apiGroup.Group("/allocation")
		{
			allocation.GET("/current", allocationHandler.GetCurrentAllocation)
			allocation.GET("/by-type", allocationHandler.GetAllocationByType)
			allocation.GET("/by-asset", allocationHandler.GetAllocationByAsset)
		}

		// Asset Snapshots 路由
		snapshots := apiGroup.Group("/snapshots")
		{
			snapshots.POST("", assetSnapshotHandler.CreateSnapshot)
			snapshots.POST("/trigger", assetSnapshotHandler.TriggerDailySnapshots) // 手動觸發每日快照
			snapshots.GET("/trend", assetSnapshotHandler.GetAssetTrend)
			snapshots.GET("/latest", assetSnapshotHandler.GetLatestSnapshot)
			snapshots.PUT("", assetSnapshotHandler.UpdateSnapshot)
			snapshots.DELETE("", assetSnapshotHandler.DeleteSnapshot)
		}
	}

	// 啟動伺服器
	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}