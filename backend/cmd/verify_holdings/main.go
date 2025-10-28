package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type HoldingResponse struct {
	Data  []models.Holding `json:"data"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

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

	// 建立 repository
	exchangeRateRepo := repository.NewExchangeRateRepository(db)

	// 確認今天有匯率資料
	log.Println("\n=== 檢查匯率資料 ===")
	today := time.Now().Truncate(24 * time.Hour)
	rate, err := exchangeRateRepo.GetByDate(models.CurrencyUSD, models.CurrencyTWD, today)
	if err != nil {
		log.Printf("❌ Error getting exchange rate: %v", err)
	} else if rate == nil {
		log.Printf("⚠️  No exchange rate found for today, inserting default rate...")
		input := &models.ExchangeRateInput{
			FromCurrency: models.CurrencyUSD,
			ToCurrency:   models.CurrencyTWD,
			Rate:         31.5,
			Date:         today,
		}
		rate, err = exchangeRateRepo.Upsert(input)
		if err != nil {
			log.Fatalf("Failed to insert exchange rate: %v", err)
		}
		log.Printf("✅ Inserted exchange rate: USD/TWD = %.4f", rate.Rate)
	} else {
		log.Printf("✅ Found exchange rate: USD/TWD = %.4f (date: %s)", rate.Rate, rate.Date.Format("2006-01-02"))
	}

	// 從 API 取得所有持倉
	log.Println("\n=== 從 API 取得所有持倉 ===")
	resp, err := http.Get("http://localhost:8080/api/holdings")
	if err != nil {
		log.Fatalf("Failed to call API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	var holdingResp HoldingResponse
	if err := json.Unmarshal(body, &holdingResp); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	if holdingResp.Error != nil {
		log.Fatalf("API error: %s - %s", holdingResp.Error.Code, holdingResp.Error.Message)
	}

	holdings := holdingResp.Data
	log.Printf("Found %d holdings\n", len(holdings))

	// 顯示每個持倉的詳細資訊
	for _, holding := range holdings {
		log.Println("\n" + strings.Repeat("=", 80))
		log.Printf("標的: %s (%s)", holding.Symbol, holding.Name)
		log.Printf("資產類型: %s", holding.AssetType)
		log.Printf("幣別: %s", holding.Currency)
		log.Println(strings.Repeat("-", 80))

		// 基本資訊
		log.Printf("持有數量: %.4f", holding.Quantity)
		log.Printf("平均成本: %.2f TWD", holding.AvgCost)
		log.Printf("總成本: %.2f TWD", holding.TotalCost)

		// 價格資訊
		log.Println(strings.Repeat("-", 80))
		log.Printf("當前價格 (%s): %.2f", holding.Currency, holding.CurrentPrice)
		log.Printf("當前價格 (TWD): %.2f", holding.CurrentPriceTWD)

		// 驗證匯率轉換
		if holding.Currency != models.CurrencyTWD && holding.CurrentPrice > 0 {
			calculatedRate := holding.CurrentPriceTWD / holding.CurrentPrice
			log.Printf("計算出的匯率: %.4f", calculatedRate)
			log.Printf("預期匯率: %.4f", rate.Rate)
			if abs(calculatedRate-rate.Rate) > 0.01 {
				log.Printf("⚠️  WARNING: 匯率不符！")
			} else {
				log.Printf("✅ 匯率轉換正確")
			}
		}

		// 市值與損益
		log.Println(strings.Repeat("-", 80))
		log.Printf("市值 (TWD): %.2f", holding.MarketValue)
		log.Printf("未實現損益 (TWD): %.2f", holding.UnrealizedPL)
		log.Printf("損益百分比: %.2f%%", holding.UnrealizedPLPct)

		// 驗證計算
		log.Println(strings.Repeat("-", 80))
		expectedMarketValue := holding.Quantity * holding.CurrentPriceTWD
		expectedUnrealizedPL := expectedMarketValue - holding.TotalCost
		var expectedPLPct float64
		if holding.TotalCost > 0 {
			expectedPLPct = (expectedUnrealizedPL / holding.TotalCost) * 100
		}

		log.Printf("驗證市值: %.2f (預期) vs %.2f (實際)", expectedMarketValue, holding.MarketValue)
		log.Printf("驗證損益: %.2f (預期) vs %.2f (實際)", expectedUnrealizedPL, holding.UnrealizedPL)
		log.Printf("驗證損益%%: %.2f%% (預期) vs %.2f%% (實際)", expectedPLPct, holding.UnrealizedPLPct)

		if abs(expectedMarketValue-holding.MarketValue) > 0.01 {
			log.Printf("❌ 市值計算錯誤！")
		} else if abs(expectedUnrealizedPL-holding.UnrealizedPL) > 0.01 {
			log.Printf("❌ 損益計算錯誤！")
		} else if abs(expectedPLPct-holding.UnrealizedPLPct) > 0.01 {
			log.Printf("❌ 損益百分比計算錯誤！")
		} else {
			log.Printf("✅ 所有計算正確！")
		}

		// 價格來源資訊
		log.Println(strings.Repeat("-", 80))
		log.Printf("價格來源: %s", holding.PriceSource)
		if holding.IsPriceStale {
			log.Printf("⚠️  價格過期: %s", holding.PriceStaleReason)
		}
	}

	log.Println("\n" + strings.Repeat("=", 80))
	log.Println("驗證完成！")
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

