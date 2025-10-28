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

	// 測試連線
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database")

	// 建立 services
	exchangeRateRepo := repository.NewExchangeRateRepository(db)
	bankClient := client.NewTaiwanBankClient()
	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, bankClient, nil)

	// 測試 1: 查看資料庫中的匯率
	log.Println("\n=== 測試 1: 查看資料庫中的匯率 ===")
	today := time.Now().Truncate(24 * time.Hour)
	dbRate, err := exchangeRateRepo.GetByDate(models.CurrencyUSD, models.CurrencyTWD, today)
	if err != nil {
		log.Printf("❌ Error getting rate from DB: %v", err)
	} else if dbRate == nil {
		log.Printf("⚠️  No rate found in DB for today (%s)", today.Format("2006-01-02"))
	} else {
		log.Printf("✅ Found rate in DB: USD/TWD = %.4f (date: %s, updated: %s)",
			dbRate.Rate, dbRate.Date.Format("2006-01-02"), dbRate.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	// 測試 2: 測試 ConvertToTWD
	log.Println("\n=== 測試 2: 測試 ConvertToTWD ===")
	testAmount := 268.81
	testDate := time.Now()

	result, err := exchangeRateService.ConvertToTWD(testAmount, models.CurrencyUSD, testDate)
	if err != nil {
		log.Printf("❌ ConvertToTWD failed: %v", err)
	} else {
		log.Printf("✅ ConvertToTWD success: %.2f USD = %.2f TWD", testAmount, result)
	}

	// 測試 3: 測試 GetRate
	log.Println("\n=== 測試 3: 測試 GetRate ===")
	rate, err := exchangeRateService.GetRate(models.CurrencyUSD, models.CurrencyTWD, testDate)
	if err != nil {
		log.Printf("❌ GetRate failed: %v", err)
	} else {
		log.Printf("✅ GetRate success: USD/TWD = %.4f", rate)
		log.Printf("   Manual calculation: %.2f USD × %.4f = %.2f TWD", testAmount, rate, testAmount*rate)
	}

	// 測試 4: 更新今日匯率
	log.Println("\n=== 測試 4: 更新今日匯率 ===")
	if err := exchangeRateService.RefreshTodayRate(); err != nil {
		log.Printf("❌ RefreshTodayRate failed: %v", err)
	} else {
		log.Println("✅ RefreshTodayRate success")

		// 再次查詢
		dbRate, _ = exchangeRateRepo.GetByDate(models.CurrencyUSD, models.CurrencyTWD, today)
		if dbRate != nil {
			log.Printf("   Updated rate: USD/TWD = %.4f (updated: %s)",
				dbRate.Rate, dbRate.UpdatedAt.Format("2006-01-02 15:04:05"))
		}
	}
}

