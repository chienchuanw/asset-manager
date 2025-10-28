package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
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

	// 測試連線
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database")

	// 建立 repository
	exchangeRateRepo := repository.NewExchangeRateRepository(db)

	// 插入今天的匯率（使用合理的匯率值）
	today := time.Now().Truncate(24 * time.Hour)
	rate := 31.5 // 假設的 USD/TWD 匯率

	input := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         rate,
		Date:         today,
	}

	result, err := exchangeRateRepo.Upsert(input)
	if err != nil {
		log.Fatalf("❌ Failed to insert rate: %v", err)
	}

	log.Printf("✅ Successfully inserted/updated exchange rate:")
	log.Printf("   Date: %s", result.Date.Format("2006-01-02"))
	log.Printf("   Rate: USD/TWD = %.4f", result.Rate)
	log.Printf("   Created: %s", result.CreatedAt.Format("2006-01-02 15:04:05"))
	log.Printf("   Updated: %s", result.UpdatedAt.Format("2006-01-02 15:04:05"))
}

