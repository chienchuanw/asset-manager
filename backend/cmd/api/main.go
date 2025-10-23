package main

import (
	"log"

	"github.com/chienchuanw/asset-manager/internal/api"
	"github.com/chienchuanw/asset-manager/internal/db"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化資料庫連線
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 初始化 Repository
	transactionRepo := repository.NewTransactionRepository(database)

	// 初始化 Service
	transactionService := service.NewTransactionService(transactionRepo)

	// 初始化 Handler
	transactionHandler := api.NewTransactionHandler(transactionService)

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
		transactions := apiGroup.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.ListTransactions)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
		}
	}

	// 啟動伺服器
	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}