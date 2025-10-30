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

	// 測試 BTC
	fmt.Println("\n=== 測試 BTC 持倉計算 ===")
	testSymbol(transactionRepo, fifoCalculator, "BTC")

	// 測試 AAPL
	fmt.Println("\n=== 測試 AAPL 持倉計算 ===")
	testSymbol(transactionRepo, fifoCalculator, "AAPL")
}

func testSymbol(transactionRepo repository.TransactionRepository, fifoCalculator service.FIFOCalculator, symbol string) {
	// 取得交易記錄
	symbolFilter := symbol
	filters := repository.TransactionFilters{
		Symbol: &symbolFilter,
	}

	transactions, err := transactionRepo.GetAll(filters)
	if err != nil {
		log.Printf("❌ Failed to get transactions: %v", err)
		return
	}

	if len(transactions) == 0 {
		log.Printf("⚠️  No transactions found for %s", symbol)
		return
	}

	fmt.Printf("找到 %d 筆交易記錄\n", len(transactions))

	// 顯示交易詳情
	for i, tx := range transactions {
		fmt.Printf("\n交易 %d:\n", i+1)
		fmt.Printf("  日期: %s\n", tx.Date.Format("2006-01-02"))
		fmt.Printf("  類型: %s\n", tx.TransactionType)
		fmt.Printf("  數量: %.8f\n", tx.Quantity)
		fmt.Printf("  價格: %.2f %s\n", tx.Price, tx.Currency)
		fmt.Printf("  金額: %.2f %s\n", tx.Amount, tx.Currency)
		if tx.Fee != nil {
			fmt.Printf("  手續費: %.2f %s\n", *tx.Fee, tx.Currency)
		}
	}

	// 計算持倉
	holding, err := fifoCalculator.CalculateHoldingForSymbol(symbol, transactions)
	if err != nil {
		log.Printf("❌ Failed to calculate holding: %v", err)
		return
	}

	if holding == nil {
		log.Printf("⚠️  No holding found (all sold)")
		return
	}

	fmt.Printf("\n計算結果:\n")
	fmt.Printf("  標的: %s\n", holding.Symbol)
	fmt.Printf("  數量: %.8f\n", holding.Quantity)
	fmt.Printf("  平均成本: %.2f TWD\n", holding.AvgCost)
	fmt.Printf("  總成本: %.2f TWD\n", holding.TotalCost)

	// 檢查第一筆交易的幣別
	firstTx := transactions[0]
	if firstTx.Currency == models.CurrencyUSD {
		expectedCost := firstTx.Amount * 31.5 // 假設匯率 31.5
		if firstTx.Fee != nil {
			expectedCost += *firstTx.Fee * 31.5
		}
		fmt.Printf("\n驗證:\n")
		fmt.Printf("  交易金額: %.2f USD\n", firstTx.Amount)
		fmt.Printf("  預期成本 (31.5匯率): %.2f TWD\n", expectedCost)
		fmt.Printf("  實際成本: %.2f TWD\n", holding.TotalCost)

		if holding.TotalCost > 1000 {
			fmt.Printf("  ✅ 成本已正確轉換為 TWD\n")
		} else {
			fmt.Printf("  ❌ 成本未轉換！仍然是 USD 金額\n")
		}
	}
}

