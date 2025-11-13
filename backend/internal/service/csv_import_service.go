package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// CSVImportService CSV 匯入服務介面
type CSVImportService interface {
	GenerateTemplate() string
	ParseCSV(reader io.Reader) *models.CSVImportResult
}

type csvImportService struct{}

// NewCSVImportService 建立新的 CSV 匯入服務
func NewCSVImportService() CSVImportService {
	return &csvImportService{}
}

// GenerateTemplate 生成 CSV 樣板
func (s *csvImportService) GenerateTemplate() string {
	var sb strings.Builder

	// Header
	sb.WriteString("date,asset_type,symbol,name,transaction_type,quantity,price,fee,tax,currency,note\n")

	// 範例資料 - 台股
	sb.WriteString("2025-01-15,tw_stock,2330,台積電,buy,10,620,28,,TWD,台股買入範例\n")

	// 範例資料 - 美股
	sb.WriteString("2025-01-16,us_stock,AAPL,Apple Inc.,buy,5,185.5,1,,USD,美股買入範例\n")

	// 範例資料 - 加密貨幣
	sb.WriteString("2025-01-17,crypto,BTC,Bitcoin,buy,0.01,45000,5,,USD,加密貨幣買入範例\n")

	return sb.String()
}

// ParseCSV 解析 CSV 檔案
func (s *csvImportService) ParseCSV(reader io.Reader) *models.CSVImportResult {
	result := &models.CSVImportResult{
		Success:      true,
		Transactions: []*models.CreateTransactionInput{},
		Errors:       []models.CSVValidationError{},
	}

	csvReader := csv.NewReader(reader)

	// 讀取 header
	headers, err := csvReader.Read()
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, models.CSVValidationError{
			Row:     0,
			Field:   "header",
			Message: "無法讀取 CSV header",
		})
		return result
	}

	// 驗證 header
	expectedHeaders := []string{"date", "asset_type", "symbol", "name", "transaction_type", "quantity", "price", "fee", "tax", "currency", "note"}
	if !s.validateHeaders(headers, expectedHeaders) {
		result.Success = false
		result.Errors = append(result.Errors, models.CSVValidationError{
			Row:     0,
			Field:   "header",
			Message: "CSV header 格式不正確",
		})
		return result
	}

	// 讀取資料行
	rowNum := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Success = false
			result.Errors = append(result.Errors, models.CSVValidationError{
				Row:     rowNum + 1,
				Field:   "row",
				Message: fmt.Sprintf("讀取第 %d 行時發生錯誤: %v", rowNum+1, err),
			})
			continue
		}

		rowNum++

		// 解析並驗證這一行
		tx, errs := s.parseRow(record, rowNum)
		if len(errs) > 0 {
			result.Success = false
			result.Errors = append(result.Errors, errs...)
		} else {
			result.Transactions = append(result.Transactions, tx)
		}
	}

	return result
}

// validateHeaders 驗證 CSV header
func (s *csvImportService) validateHeaders(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i, h := range expected {
		if actual[i] != h {
			return false
		}
	}
	return true
}

// parseRow 解析單一行資料
func (s *csvImportService) parseRow(record []string, rowNum int) (*models.CreateTransactionInput, []models.CSVValidationError) {
	errors := []models.CSVValidationError{}

	// 確保有足夠的欄位
	if len(record) < 11 {
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "row",
			Message: "欄位數量不足",
		})
		return nil, errors
	}

	tx := &models.CreateTransactionInput{}

	// 解析日期
	date, err := time.Parse("2006-01-02", strings.TrimSpace(record[0]))
	if err != nil {
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "date",
			Message: "日期格式錯誤，應為 YYYY-MM-DD",
		})
	} else {
		tx.Date = date
	}

	// 解析資產類別
	assetType := strings.TrimSpace(record[1])
	switch assetType {
	case "tw_stock":
		tx.AssetType = models.AssetTypeTWStock
	case "us_stock":
		tx.AssetType = models.AssetTypeUSStock
	case "crypto":
		tx.AssetType = models.AssetTypeCrypto
	default:
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "asset_type",
			Message: "資產類別無效，應為 tw_stock, us_stock 或 crypto",
		})
	}

	// 解析 Symbol（必填）
	symbol := strings.TrimSpace(record[2])
	if symbol == "" {
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "symbol",
			Message: "交易標的代碼為必填欄位",
		})
	} else {
		tx.Symbol = symbol
	}

	// 解析 Name（必填）
	name := strings.TrimSpace(record[3])
	if name == "" {
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "name",
			Message: "交易標的名稱為必填欄位",
		})
	} else {
		tx.Name = name
	}

	// 解析交易類型
	transactionType := strings.TrimSpace(record[4])
	switch transactionType {
	case "buy":
		tx.TransactionType = models.TransactionTypeBuy
	case "sell":
		tx.TransactionType = models.TransactionTypeSell
	case "dividend":
		tx.TransactionType = models.TransactionTypeDividend
	case "fee":
		tx.TransactionType = models.TransactionTypeFee
	default:
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "transaction_type",
			Message: "交易類型無效，應為 buy, sell, dividend 或 fee",
		})
	}

	// 解析數量
	quantity, err := strconv.ParseFloat(strings.TrimSpace(record[5]), 64)
	if err != nil {
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "quantity",
			Message: "數量格式錯誤，應為數字",
		})
	} else if quantity < 0 {
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "quantity",
			Message: "數量不可為負數",
		})
	} else {
		tx.Quantity = quantity
	}

	// 解析單價
	price, err := strconv.ParseFloat(strings.TrimSpace(record[6]), 64)
	if err != nil {
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "price",
			Message: "單價格式錯誤，應為數字",
		})
	} else if price < 0 {
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "price",
			Message: "單價不可為負數",
		})
	} else {
		tx.Price = price
	}

	// 計算金額
	tx.Amount = tx.Quantity * tx.Price

	// 解析手續費（選填）
	feeStr := strings.TrimSpace(record[7])
	if feeStr != "" {
		fee, err := strconv.ParseFloat(feeStr, 64)
		if err != nil {
			errors = append(errors, models.CSVValidationError{
				Row:     rowNum,
				Field:   "fee",
				Message: "手續費格式錯誤，應為數字",
			})
		} else if fee < 0 {
			errors = append(errors, models.CSVValidationError{
				Row:     rowNum,
				Field:   "fee",
				Message: "手續費不可為負數",
			})
		} else {
			tx.Fee = &fee
		}
	}

	// 解析交易稅（選填）
	taxStr := strings.TrimSpace(record[8])
	if taxStr != "" {
		tax, err := strconv.ParseFloat(taxStr, 64)
		if err != nil {
			errors = append(errors, models.CSVValidationError{
				Row:     rowNum,
				Field:   "tax",
				Message: "交易稅格式錯誤，應為數字",
			})
		} else if tax < 0 {
			errors = append(errors, models.CSVValidationError{
				Row:     rowNum,
				Field:   "tax",
				Message: "交易稅不可為負數",
			})
		} else {
			tx.Tax = &tax
		}
	}

	// 解析幣別
	currency := strings.TrimSpace(record[9])
	switch currency {
	case "TWD":
		tx.Currency = models.CurrencyTWD
	case "USD":
		tx.Currency = models.CurrencyUSD
	default:
		errors = append(errors, models.CSVValidationError{
			Row:     rowNum,
			Field:   "currency",
			Message: "幣別無效，應為 TWD 或 USD",
		})
	}

	// 解析備註（選填）
	note := strings.TrimSpace(record[10])
	if note != "" {
		tx.Note = &note
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return tx, nil
}

