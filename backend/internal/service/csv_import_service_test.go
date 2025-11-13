package service

import (
	"strings"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestGenerateCSVTemplate 測試生成 CSV 樣板
func TestGenerateCSVTemplate(t *testing.T) {
	service := NewCSVImportService()

	csv := service.GenerateTemplate()

	// 驗證 CSV 包含 header
	assert.Contains(t, csv, "date,asset_type,symbol,name,transaction_type,quantity,price,fee,tax,currency,note")

	// 驗證 CSV 包含台股範例資料
	assert.Contains(t, csv, "2025-01-15")
	assert.Contains(t, csv, "tw_stock")
	assert.Contains(t, csv, "2330")
	assert.Contains(t, csv, "台積電")

	// 驗證 CSV 包含美股範例資料
	assert.Contains(t, csv, "2025-01-16")
	assert.Contains(t, csv, "us_stock")
	assert.Contains(t, csv, "AAPL")
	assert.Contains(t, csv, "Apple Inc.")

	// 驗證 CSV 包含加密貨幣範例資料
	assert.Contains(t, csv, "2025-01-17")
	assert.Contains(t, csv, "crypto")
	assert.Contains(t, csv, "BTC")
	assert.Contains(t, csv, "Bitcoin")
}

// TestParseCSV_Success 測試成功解析 CSV
func TestParseCSV_Success(t *testing.T) {
	service := NewCSVImportService()

	csvContent := `date,asset_type,symbol,name,transaction_type,quantity,price,fee,tax,currency,note
2025-01-15,tw_stock,2330,台積電,buy,10,620,28,,TWD,測試交易`

	result := service.ParseCSV(strings.NewReader(csvContent))

	assert.True(t, result.Success)
	assert.Len(t, result.Transactions, 1)
	assert.Len(t, result.Errors, 0)

	// 驗證解析後的資料
	tx := result.Transactions[0]
	assert.Equal(t, models.AssetTypeTWStock, tx.AssetType)
	assert.Equal(t, "2330", tx.Symbol)
	assert.Equal(t, "台積電", tx.Name)
	assert.Equal(t, models.TransactionTypeBuy, tx.TransactionType)
	assert.Equal(t, 10.0, tx.Quantity)
	assert.Equal(t, 620.0, tx.Price)
	assert.Equal(t, models.CurrencyTWD, tx.Currency)
}

// TestParseCSV_InvalidDate 測試無效的日期格式
func TestParseCSV_InvalidDate(t *testing.T) {
	service := NewCSVImportService()

	csvContent := `date,asset_type,symbol,name,transaction_type,quantity,price,fee,tax,currency,note
2025/01/15,tw_stock,2330,台積電,buy,10,620,28,,TWD,`

	result := service.ParseCSV(strings.NewReader(csvContent))

	assert.False(t, result.Success)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, 1, result.Errors[0].Row)
	assert.Equal(t, "date", result.Errors[0].Field)
}

// TestParseCSV_MissingRequiredField 測試缺少必填欄位
func TestParseCSV_MissingRequiredField(t *testing.T) {
	service := NewCSVImportService()

	csvContent := `date,asset_type,symbol,name,transaction_type,quantity,price,fee,tax,currency,note
2025-01-15,tw_stock,,台積電,buy,10,620,28,,TWD,`

	result := service.ParseCSV(strings.NewReader(csvContent))

	assert.False(t, result.Success)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "symbol", result.Errors[0].Field)
}

// TestParseCSV_InvalidAssetType 測試無效的資產類別
func TestParseCSV_InvalidAssetType(t *testing.T) {
	service := NewCSVImportService()

	csvContent := `date,asset_type,symbol,name,transaction_type,quantity,price,fee,tax,currency,note
2025-01-15,invalid_type,2330,台積電,buy,10,620,28,,TWD,`

	result := service.ParseCSV(strings.NewReader(csvContent))

	assert.False(t, result.Success)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "asset_type", result.Errors[0].Field)
}

// TestParseCSV_InvalidQuantity 測試無效的數量
func TestParseCSV_InvalidQuantity(t *testing.T) {
	service := NewCSVImportService()

	csvContent := `date,asset_type,symbol,name,transaction_type,quantity,price,fee,tax,currency,note
2025-01-15,tw_stock,2330,台積電,buy,abc,620,28,,TWD,`

	result := service.ParseCSV(strings.NewReader(csvContent))

	assert.False(t, result.Success)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "quantity", result.Errors[0].Field)
}

// TestParseCSV_MultipleRows 測試多筆交易記錄
func TestParseCSV_MultipleRows(t *testing.T) {
	service := NewCSVImportService()

	csvContent := `date,asset_type,symbol,name,transaction_type,quantity,price,fee,tax,currency,note
2025-01-15,tw_stock,2330,台積電,buy,10,620,28,,TWD,
2025-01-16,us_stock,AAPL,Apple,buy,5,150,1,,USD,`

	result := service.ParseCSV(strings.NewReader(csvContent))

	assert.True(t, result.Success)
	assert.Len(t, result.Transactions, 2)
	assert.Len(t, result.Errors, 0)
}

// TestParseCSV_PartialErrors 測試部分資料有錯誤
func TestParseCSV_PartialErrors(t *testing.T) {
	service := NewCSVImportService()

	csvContent := `date,asset_type,symbol,name,transaction_type,quantity,price,fee,tax,currency,note
2025-01-15,tw_stock,2330,台積電,buy,10,620,28,,TWD,
2025/01/16,us_stock,AAPL,Apple,buy,5,150,1,,USD,
2025-01-17,crypto,BTC,Bitcoin,buy,0.1,50000,,,USD,`

	result := service.ParseCSV(strings.NewReader(csvContent))

	// 即使有部分錯誤，也應該標記為失敗
	assert.False(t, result.Success)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, 2, result.Errors[0].Row) // 第二行有錯誤
}

