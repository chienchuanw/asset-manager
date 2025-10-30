package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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

	// 建立資料庫連線
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database")

	// 建立 services
	exchangeRateRepo := repository.NewExchangeRateRepository(db)
	exchangeRateAPIClient := client.NewExchangeRateAPIClient()
	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, exchangeRateAPIClient, nil)

	// 測試場景 1: 有今天的匯率
	log.Println("\n=== 場景 1: 資料庫有今天的匯率 ===")
	today := time.Now().Truncate(24 * time.Hour)
	
	// 插入今天的匯率
	input := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         31.5,
		Date:         today,
	}
	exchangeRateRepo.Upsert(input)
	
	result, err := exchangeRateService.ConvertToTWD(100.0, models.CurrencyUSD, today)
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		log.Printf("✅ 100 USD = %.2f TWD (expected: 3150.00)", result)
	}

	// 測試場景 2: 沒有今天的匯率，但有歷史匯率
	log.Println("\n=== 場景 2: 沒有今天的匯率，使用最新匯率 ===")
	
	// 刪除今天的匯率
	db.Exec("DELETE FROM exchange_rates WHERE date = $1", today)
	
	// 插入昨天的匯率
	yesterday := today.AddDate(0, 0, -1)
	input2 := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         32.0,
		Date:         yesterday,
	}
	exchangeRateRepo.Upsert(input2)
	
	result, err = exchangeRateService.ConvertToTWD(100.0, models.CurrencyUSD, today)
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		log.Printf("✅ 100 USD = %.2f TWD (expected: 3200.00, using yesterday's rate)", result)
	}

	// 測試場景 3: 資料庫完全沒有匯率，使用預設值
	log.Println("\n=== 場景 3: 資料庫完全沒有匯率，使用預設值 ===")
	
	// 清空所有匯率
	db.Exec("DELETE FROM exchange_rates")
	
	result, err = exchangeRateService.ConvertToTWD(100.0, models.CurrencyUSD, today)
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		log.Printf("✅ 100 USD = %.2f TWD (expected: 3000.00, using default rate 30.0)", result)
	}

	// 恢復今天的匯率
	log.Println("\n=== 恢復資料 ===")
	exchangeRateRepo.Upsert(input)
	log.Println("✅ Restored today's exchange rate")
}

