package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/chienchuanw/asset-manager/internal/api"
	"github.com/chienchuanw/asset-manager/internal/cache"
	"github.com/chienchuanw/asset-manager/internal/client"
	"github.com/chienchuanw/asset-manager/internal/db"
	"github.com/chienchuanw/asset-manager/internal/middleware"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/scheduler"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
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
	settingsRepo := repository.NewSettingsRepository(database)
	cashFlowRepo := repository.NewCashFlowRepository(database)
	categoryRepo := repository.NewCategoryRepository(database)
	subscriptionRepo := repository.NewSubscriptionRepository(database)
	installmentRepo := repository.NewInstallmentRepository(database)
	bankAccountRepo := repository.NewBankAccountRepository(database)
	creditCardRepo := repository.NewCreditCardRepository(database)
	creditCardGroupRepo := repository.NewCreditCardGroupRepository(database)

	// PerformanceSnapshotRepository 需要 sqlx.DB
	dbx := sqlx.NewDb(database, "postgres")
	performanceSnapshotRepo := repository.NewPerformanceSnapshotRepository(dbx)
	schedulerLogRepo := repository.NewSchedulerLogRepository(dbx)
	cashFlowReportLogRepo := repository.NewCashFlowReportLogRepository(database)

	authService := service.NewAuthService()

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
		exchangeRateClient := client.NewExchangeRateAPIClient()
		exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, exchangeRateClient, nil)

		// 初始化 FIFO Calculator（需要 exchangeRateService）
		fifoCalculator := service.NewFIFOCalculator(exchangeRateService)

		// 初始化 TransactionService
		transactionService := service.NewTransactionService(transactionRepo, realizedProfitRepo, fifoCalculator, exchangeRateService)

		holdingService := service.NewHoldingService(transactionRepo, fifoCalculator, priceService, exchangeRateService)

		// 初始化 Analytics Service
		analyticsService := service.NewAnalyticsService(realizedProfitRepo)
		unrealizedAnalyticsService := service.NewUnrealizedAnalyticsService(holdingService)
		allocationService := service.NewAllocationService(holdingService)
		performanceTrendService := service.NewPerformanceTrendService(performanceSnapshotRepo, unrealizedAnalyticsService, analyticsService)
		settingsService := service.NewSettingsService(settingsRepo)
		discordService := service.NewDiscordService()
		rebalanceService := service.NewRebalanceService(settingsService, holdingService)
		cashFlowService := service.NewCashFlowService(cashFlowRepo, categoryRepo, bankAccountRepo, creditCardRepo)
		categoryService := service.NewCategoryService(categoryRepo)
		subscriptionService := service.NewSubscriptionService(subscriptionRepo, categoryRepo)
		installmentService := service.NewInstallmentService(installmentRepo, categoryRepo)
		billingService := service.NewBillingService(subscriptionRepo, installmentRepo, cashFlowRepo)
		bankAccountService := service.NewBankAccountService(bankAccountRepo)
		creditCardService := service.NewCreditCardService(creditCardRepo)
		creditCardGroupService := service.NewCreditCardGroupService(creditCardGroupRepo, creditCardRepo)

		// 初始化 Asset Snapshot Service（不帶排程器）
		assetSnapshotService := service.NewAssetSnapshotServiceWithDeps(assetSnapshotRepo, holdingService)

		// 初始化 CSV Import Service
		csvImportService := service.NewCSVImportService()

		// 初始化 Handler
		authHandler := api.NewAuthHandler(authService)
		transactionHandler := api.NewTransactionHandler(transactionService, csvImportService)
		holdingHandler := api.NewHoldingHandler(holdingService)
		analyticsHandler := api.NewAnalyticsHandler(analyticsService)
		unrealizedAnalyticsHandler := api.NewUnrealizedAnalyticsHandler(unrealizedAnalyticsService)
		allocationHandler := api.NewAllocationHandler(allocationService)
		performanceTrendHandler := api.NewPerformanceTrendHandler(performanceTrendService)
		settingsHandler := api.NewSettingsHandler(settingsService)
		assetSnapshotHandler := api.NewAssetSnapshotHandler(assetSnapshotService)
		discordHandler := api.NewDiscordHandler(discordService, settingsService, holdingService, rebalanceService)
		rebalanceHandler := api.NewRebalanceHandler(rebalanceService)
		cashFlowHandler := api.NewCashFlowHandler(cashFlowService)
		cashFlowHandler.SetDiscordService(discordService) // 設定 Discord service 用於發送報告
		categoryHandler := api.NewCategoryHandler(categoryService)
		subscriptionHandler := api.NewSubscriptionHandler(subscriptionService)
		installmentHandler := api.NewInstallmentHandler(installmentService)
		billingHandler := api.NewBillingHandler(billingService)
		bankAccountHandler := api.NewBankAccountHandler(bankAccountService)
		creditCardHandler := api.NewCreditCardHandler(creditCardService)
		creditCardGroupHandler := api.NewCreditCardGroupHandler(creditCardGroupService)
		exchangeRateHandler := api.NewExchangeRateHandler(exchangeRateService)

		// 初始化排程器管理器（不啟動）
		schedulerManagerConfig := scheduler.SchedulerManagerConfig{
			Enabled:           false, // Redis 不可用時停用排程器
			DailySnapshotTime: "23:59",
		}
		schedulerManager := scheduler.NewSchedulerManager(
			assetSnapshotService,
			discordService,
			settingsService,
			holdingService,
			rebalanceService,
			billingService,
			exchangeRateService,
			creditCardService,
			cashFlowService,
			nil, // schedulerLogRepo 設為 nil（因為 Redis 不可用時也不記錄）
			nil, // cashFlowReportLogRepo 設為 nil
			schedulerManagerConfig,
		)
		schedulerHandler := api.NewSchedulerHandler(schedulerManager)

		// 建立 router 並啟動（簡化版，不啟動排程器）
		log.Println("Warning: Scheduler is disabled (Redis not available)")
		startServer(authHandler, transactionHandler, holdingHandler, analyticsHandler, unrealizedAnalyticsHandler, allocationHandler, performanceTrendHandler, settingsHandler, assetSnapshotHandler, discordHandler, schedulerHandler, rebalanceHandler, cashFlowHandler, categoryHandler, subscriptionHandler, installmentHandler, billingHandler, bankAccountHandler, creditCardHandler, creditCardGroupHandler, exchangeRateHandler, nil)
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
	exchangeRateClient := client.NewExchangeRateAPIClient()
	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, exchangeRateClient, redisCache.GetClient())

	// 初始化 FIFO Calculator（需要 exchangeRateService）
	fifoCalculator := service.NewFIFOCalculator(exchangeRateService)

	// 初始化 TransactionService
	transactionService := service.NewTransactionService(transactionRepo, realizedProfitRepo, fifoCalculator, exchangeRateService)

	// 初始化 Holding Service
	holdingService := service.NewHoldingService(transactionRepo, fifoCalculator, priceService, exchangeRateService)

	// 初始化 Analytics Service
	analyticsService := service.NewAnalyticsService(realizedProfitRepo)
	unrealizedAnalyticsService := service.NewUnrealizedAnalyticsService(holdingService)
	allocationService := service.NewAllocationService(holdingService)
	performanceTrendService := service.NewPerformanceTrendService(performanceSnapshotRepo, unrealizedAnalyticsService, analyticsService)
	settingsService := service.NewSettingsService(settingsRepo)
	discordService := service.NewDiscordService()
	rebalanceService := service.NewRebalanceService(settingsService, holdingService)
	cashFlowService := service.NewCashFlowService(cashFlowRepo, categoryRepo, bankAccountRepo, creditCardRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, categoryRepo)
	installmentService := service.NewInstallmentService(installmentRepo, categoryRepo)
	billingService := service.NewBillingService(subscriptionRepo, installmentRepo, cashFlowRepo)
	bankAccountService := service.NewBankAccountService(bankAccountRepo)
	creditCardService := service.NewCreditCardService(creditCardRepo)
	creditCardGroupService := service.NewCreditCardGroupService(creditCardGroupRepo, creditCardRepo)

	// 初始化 Asset Snapshot Service（包含依賴）
	assetSnapshotService := service.NewAssetSnapshotServiceWithDeps(assetSnapshotRepo, holdingService)

	// 初始化 CSV Import Service
	csvImportService := service.NewCSVImportService()

	// 初始化 Handler
	authHandler := api.NewAuthHandler(authService)
	transactionHandler := api.NewTransactionHandler(transactionService, csvImportService)
	holdingHandler := api.NewHoldingHandler(holdingService)
	analyticsHandler := api.NewAnalyticsHandler(analyticsService)
	unrealizedAnalyticsHandler := api.NewUnrealizedAnalyticsHandler(unrealizedAnalyticsService)
	allocationHandler := api.NewAllocationHandler(allocationService)
	performanceTrendHandler := api.NewPerformanceTrendHandler(performanceTrendService)
	settingsHandler := api.NewSettingsHandler(settingsService)
	assetSnapshotHandler := api.NewAssetSnapshotHandler(assetSnapshotService)
	discordHandler := api.NewDiscordHandler(discordService, settingsService, holdingService, rebalanceService)
	rebalanceHandler := api.NewRebalanceHandler(rebalanceService)
	cashFlowHandler := api.NewCashFlowHandler(cashFlowService)
	cashFlowHandler.SetDiscordService(discordService) // 設定 Discord service 用於發送報告
	categoryHandler := api.NewCategoryHandler(categoryService)
	subscriptionHandler := api.NewSubscriptionHandler(subscriptionService)
	installmentHandler := api.NewInstallmentHandler(installmentService)
	billingHandler := api.NewBillingHandler(billingService)
	bankAccountHandler := api.NewBankAccountHandler(bankAccountService)
	creditCardHandler := api.NewCreditCardHandler(creditCardService)
	creditCardGroupHandler := api.NewCreditCardGroupHandler(creditCardGroupService)
	exchangeRateHandler := api.NewExchangeRateHandler(exchangeRateService)

	// 初始化並啟動排程器管理器
	schedulerManagerConfig := scheduler.SchedulerManagerConfig{
		Enabled:           os.Getenv("SNAPSHOT_SCHEDULER_ENABLED") == "true",
		DailySnapshotTime: getEnvOrDefault("SCHEDULER_SNAPSHOT_TIME", "23:59"),
	}
	schedulerManager := scheduler.NewSchedulerManager(
		assetSnapshotService,
		discordService,
		settingsService,
		holdingService,
		rebalanceService,
		billingService,
		exchangeRateService,
		creditCardService,
		cashFlowService,
		schedulerLogRepo,
		cashFlowReportLogRepo,
		schedulerManagerConfig,
	)
	if err := schedulerManager.Start(); err != nil {
		log.Printf("Warning: Failed to start scheduler manager: %v", err)
	}

	// 初始化排程器 Handler
	schedulerHandler := api.NewSchedulerHandler(schedulerManager)

	// 啟動伺服器（會在內部處理 graceful shutdown）
	startServer(authHandler, transactionHandler, holdingHandler, analyticsHandler, unrealizedAnalyticsHandler, allocationHandler, performanceTrendHandler, settingsHandler, assetSnapshotHandler, discordHandler, schedulerHandler, rebalanceHandler, cashFlowHandler, categoryHandler, subscriptionHandler, installmentHandler, billingHandler, bankAccountHandler, creditCardHandler, creditCardGroupHandler, exchangeRateHandler, schedulerManager)
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
func startServer(authHandler *api.AuthHandler, transactionHandler *api.TransactionHandler, holdingHandler *api.HoldingHandler, analyticsHandler *api.AnalyticsHandler, unrealizedAnalyticsHandler *api.UnrealizedAnalyticsHandler, allocationHandler *api.AllocationHandler, performanceTrendHandler *api.PerformanceTrendHandler, settingsHandler *api.SettingsHandler, assetSnapshotHandler *api.AssetSnapshotHandler, discordHandler *api.DiscordHandler, schedulerHandler *api.SchedulerHandler, rebalanceHandler *api.RebalanceHandler, cashFlowHandler *api.CashFlowHandler, categoryHandler *api.CategoryHandler, subscriptionHandler *api.SubscriptionHandler, installmentHandler *api.InstallmentHandler, billingHandler *api.BillingHandler, bankAccountHandler *api.BankAccountHandler, creditCardHandler *api.CreditCardHandler, creditCardGroupHandler *api.CreditCardGroupHandler, exchangeRateHandler *api.ExchangeRateHandler, schedulerManager *scheduler.SchedulerManager) {
	// 建立 Gin router
	router := gin.Default()

	// 設定 CORS
	// 從環境變數讀取允許的來源，預設為 localhost:3000
	allowedOrigins := []string{"http://localhost:3000"}
	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 重要：允許發送 cookies
		MaxAge:           12 * 3600,
	}))

	// Health check endpoint (不需要驗證)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Asset Manager API Server is running.",
		})
	})

	// Auth routes (不需要驗證)
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
	}

	// API routes (需要驗證)
	apiGroup := router.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware())
	{
		// Transactions 路由
		transactions := apiGroup.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.POST("/batch", transactionHandler.CreateTransactionsBatch)
			transactions.GET("", transactionHandler.ListTransactions)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
			transactions.GET("/template", transactionHandler.DownloadCSVTemplate)
			transactions.POST("/parse-csv", transactionHandler.ParseCSV)
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

		// Performance Trends 路由
		performanceTrends := apiGroup.Group("/performance-trends")
		{
			performanceTrends.POST("/snapshot", performanceTrendHandler.CreateDailySnapshot)
			performanceTrends.GET("/range", performanceTrendHandler.GetTrendByDateRange)
			performanceTrends.GET("/latest", performanceTrendHandler.GetLatestTrend)
		}

		// Settings 路由
		settings := apiGroup.Group("/settings")
		{
			settings.GET("", settingsHandler.GetSettings)
			settings.PUT("", settingsHandler.UpdateSettings)
		}

		// Discord 路由
		discord := apiGroup.Group("/discord")
		{
			discord.POST("/test", discordHandler.TestDiscord)
			discord.POST("/daily-report", discordHandler.SendDailyReport)
		}

		// Scheduler 路由
		schedulerGroup := apiGroup.Group("/scheduler")
		{
			schedulerGroup.GET("/status", schedulerHandler.GetStatus)
			schedulerGroup.GET("/summaries", schedulerHandler.GetTaskSummaries)
			schedulerGroup.POST("/trigger/snapshot", schedulerHandler.TriggerSnapshot)
			schedulerGroup.POST("/trigger/discord-report", schedulerHandler.TriggerDiscordReport)
			schedulerGroup.POST("/reload/discord", schedulerHandler.ReloadDiscordSchedule)
		}

		// Rebalance 路由
		rebalance := apiGroup.Group("/rebalance")
		{
			rebalance.GET("/check", rebalanceHandler.CheckRebalance)
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

		// Cash Flows 路由
		cashFlows := apiGroup.Group("/cash-flows")
		{
			cashFlows.POST("", cashFlowHandler.CreateCashFlow)
			cashFlows.GET("", cashFlowHandler.ListCashFlows)
			cashFlows.GET("/summary", cashFlowHandler.GetSummary)
			cashFlows.GET("/monthly-summary", cashFlowHandler.GetMonthlySummary)
			cashFlows.GET("/yearly-summary", cashFlowHandler.GetYearlySummary)
			cashFlows.POST("/send-monthly-report", cashFlowHandler.SendMonthlyReport)
			cashFlows.POST("/send-yearly-report", cashFlowHandler.SendYearlyReport)
			cashFlows.GET("/:id", cashFlowHandler.GetCashFlow)
			cashFlows.PUT("/:id", cashFlowHandler.UpdateCashFlow)
			cashFlows.DELETE("/:id", cashFlowHandler.DeleteCashFlow)
		}

		// Categories 路由
		categories := apiGroup.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("", categoryHandler.ListCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		// Subscriptions 路由
		subscriptions := apiGroup.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.CreateSubscription)
			subscriptions.GET("", subscriptionHandler.ListSubscriptions)
			subscriptions.GET("/:id", subscriptionHandler.GetSubscription)
			subscriptions.PUT("/:id", subscriptionHandler.UpdateSubscription)
			subscriptions.DELETE("/:id", subscriptionHandler.DeleteSubscription)
			subscriptions.POST("/:id/cancel", subscriptionHandler.CancelSubscription)
		}

		// Installments 路由
		installments := apiGroup.Group("/installments")
		{
			installments.POST("", installmentHandler.CreateInstallment)
			installments.GET("", installmentHandler.ListInstallments)
			installments.GET("/completing-soon", installmentHandler.GetCompletingSoon)
			installments.GET("/:id", installmentHandler.GetInstallment)
			installments.PUT("/:id", installmentHandler.UpdateInstallment)
			installments.DELETE("/:id", installmentHandler.DeleteInstallment)
		}

		// Billing 路由
		billing := apiGroup.Group("/billing")
		{
			billing.POST("/process-daily", billingHandler.ProcessDailyBilling)
			billing.POST("/process-subscriptions", billingHandler.ProcessSubscriptionBilling)
			billing.POST("/process-installments", billingHandler.ProcessInstallmentBilling)
		}

		// Bank Accounts 路由
		bankAccounts := apiGroup.Group("/bank-accounts")
		{
			bankAccounts.POST("", bankAccountHandler.CreateBankAccount)
			bankAccounts.GET("", bankAccountHandler.ListBankAccounts)
			bankAccounts.GET("/:id", bankAccountHandler.GetBankAccount)
			bankAccounts.PUT("/:id", bankAccountHandler.UpdateBankAccount)
			bankAccounts.DELETE("/:id", bankAccountHandler.DeleteBankAccount)
		}

		// Credit Cards 路由
		creditCards := apiGroup.Group("/credit-cards")
		{
			creditCards.POST("", creditCardHandler.CreateCreditCard)
			creditCards.GET("", creditCardHandler.ListCreditCards)
			creditCards.GET("/upcoming-billing", creditCardHandler.GetUpcomingBilling)
			creditCards.GET("/upcoming-payment", creditCardHandler.GetUpcomingPayment)
			creditCards.GET("/:id", creditCardHandler.GetCreditCard)
			creditCards.PUT("/:id", creditCardHandler.UpdateCreditCard)
			creditCards.DELETE("/:id", creditCardHandler.DeleteCreditCard)
		}

		// Credit Card Groups 路由
		creditCardGroups := apiGroup.Group("/credit-card-groups")
		{
			creditCardGroups.POST("", creditCardGroupHandler.CreateCreditCardGroup)
			creditCardGroups.GET("", creditCardGroupHandler.ListCreditCardGroups)
			creditCardGroups.GET("/:id", creditCardGroupHandler.GetCreditCardGroup)
			creditCardGroups.PUT("/:id", creditCardGroupHandler.UpdateCreditCardGroup)
			creditCardGroups.DELETE("/:id", creditCardGroupHandler.DeleteCreditCardGroup)
			creditCardGroups.POST("/:id/cards", creditCardGroupHandler.AddCardsToGroup)
			creditCardGroups.DELETE("/:id/cards", creditCardGroupHandler.RemoveCardsFromGroup)
		}

		// Exchange Rates 路由
		exchangeRates := apiGroup.Group("/exchange-rates")
		{
			exchangeRates.POST("/refresh", exchangeRateHandler.RefreshExchangeRate)
		}
	}

	// 建立 HTTP 伺服器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 在 goroutine 中啟動伺服器
	go func() {
		log.Println("Starting server on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	// 停止排程器
	if schedulerManager != nil {
		log.Println("Stopping scheduler...")
		schedulerManager.Stop()
		// 等待正在執行的任務完成
		time.Sleep(5 * time.Second)
	}

	// 設定 5 秒的超時時間來關閉伺服器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}