package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chienchuanw/asset-manager/internal/db"
	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/joho/godotenv"
)

// CSVRecord 代表 CSV 中的一筆記錄
type CSVRecord struct {
	UpdateDate string
	Category   string
	Ticker     string
	Name       string
	Market     string
	Currency   string
	Units      string
	Price      string
}

func main() {
	// 解析命令列參數
	csvPath := flag.String("csv", "mock/Asset_Allocation - Data.csv", "Path to CSV file")
	clean := flag.Bool("clean", false, "Clean database before seeding")
	flag.Parse()

	// 載入環境變數（優先使用 .env.local）
	if err := godotenv.Load(".env.local"); err != nil {
		// 如果 .env.local 不存在，嘗試載入 .env
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

	// 連接資料庫
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 建立 Repository
	transactionRepo := repository.NewTransactionRepository(database)

	// 如果需要清空資料庫
	if *clean {
		log.Println("Cleaning database...")
		if err := cleanDatabase(database); err != nil {
			log.Fatalf("Failed to clean database: %v", err)
		}
		log.Println("Database cleaned successfully")
	}

	// 讀取 CSV
	log.Printf("Reading CSV file: %s\n", *csvPath)
	records, err := readCSV(*csvPath)
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}
	log.Printf("Found %d records in CSV\n", len(records))

	// 匯入資料
	log.Println("Importing transactions...")
	successCount := 0
	for i, record := range records {
		transaction, err := convertToTransaction(record)
		if err != nil {
			log.Printf("Warning: Failed to convert record %d (%s): %v\n", i+1, record.Ticker, err)
			continue
		}

		createdTransaction, err := transactionRepo.Create(transaction)
		if err != nil {
			log.Printf("Warning: Failed to create transaction for %s: %v\n", record.Ticker, err)
			continue
		}

		successCount++
		log.Printf("✓ Imported %s (%s) - %.2f %s @ %.2f\n",
			record.Ticker, record.Name, createdTransaction.Quantity, record.Currency, createdTransaction.Price)
	}

	log.Printf("\n✅ Import completed: %d/%d transactions imported successfully\n", successCount, len(records))
}

// readCSV 讀取 CSV 檔案
func readCSV(path string) ([]CSVRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("CSV file is empty or has no data rows")
	}

	// 跳過標題列
	records := make([]CSVRecord, 0, len(rows)-1)
	for i, row := range rows[1:] {
		if len(row) < 8 {
			log.Printf("Warning: Row %d has insufficient columns, skipping\n", i+2)
			continue
		}

		records = append(records, CSVRecord{
			UpdateDate: row[0],
			Category:   row[1],
			Ticker:     row[2],
			Name:       row[3],
			Market:     row[4],
			Currency:   row[5],
			Units:      row[6],
			Price:      row[7],
		})
	}

	return records, nil
}

// convertToTransaction 將 CSV 記錄轉換為 Transaction
func convertToTransaction(record CSVRecord) (*models.CreateTransactionInput, error) {
	// 解析日期
	date, err := time.Parse("2006-01-02", record.UpdateDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	// 解析數量
	quantity, err := parseFloat(record.Units)
	if err != nil {
		return nil, fmt.Errorf("invalid quantity: %w", err)
	}

	// 解析價格
	price, err := parseFloat(record.Price)
	if err != nil {
		return nil, fmt.Errorf("invalid price: %w", err)
	}

	// 計算總金額
	amount := quantity * price

	// 轉換資產類型
	assetType, err := convertAssetType(record.Category)
	if err != nil {
		return nil, err
	}

	// 轉換幣別
	currency, err := convertCurrency(record.Currency)
	if err != nil {
		return nil, err
	}

	// 手續費為 0
	fee := 0.0

	return &models.CreateTransactionInput{
		Date:            date,
		AssetType:       assetType,
		Symbol:          record.Ticker,
		Name:            record.Name,
		TransactionType: models.TransactionTypeBuy,
		Quantity:        quantity,
		Price:           price,
		Amount:          amount,
		Fee:             &fee,
		Currency:        currency,
		Note:            nil,
	}, nil
}

// convertAssetType 轉換資產類型
func convertAssetType(category string) (models.AssetType, error) {
	switch category {
	case "Taiwan Stocks":
		return models.AssetTypeTWStock, nil
	case "US Stocks":
		return models.AssetTypeUSStock, nil
	case "Crypto":
		return models.AssetTypeCrypto, nil
	default:
		return "", fmt.Errorf("unknown asset category: %s", category)
	}
}

// convertCurrency 轉換幣別
func convertCurrency(currency string) (models.Currency, error) {
	switch currency {
	case "TWD":
		return models.CurrencyTWD, nil
	case "USD":
		return models.CurrencyUSD, nil
	default:
		return "", fmt.Errorf("unknown currency: %s", currency)
	}
}

// parseFloat 解析浮點數（處理千分位逗號）
func parseFloat(s string) (float64, error) {
	// 移除千分位逗號
	s = strings.ReplaceAll(s, ",", "")
	return strconv.ParseFloat(s, 64)
}

// cleanDatabase 清空資料庫
func cleanDatabase(database *sql.DB) error {
	// 先刪除 realized_profits（因為有外鍵約束）
	if _, err := database.Exec("DELETE FROM realized_profits"); err != nil {
		return fmt.Errorf("failed to delete realized_profits: %w", err)
	}

	// 再刪除 transactions
	if _, err := database.Exec("DELETE FROM transactions"); err != nil {
		return fmt.Errorf("failed to delete transactions: %w", err)
	}

	return nil
}

