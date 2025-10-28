package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestTransactionStructure 測試 Transaction 結構是否包含所有必要欄位
func TestTransactionStructure(t *testing.T) {
	fee := 10.5
	tax := 15.0
	note := "test note"

	transaction := Transaction{
		ID:              uuid.New(),
		Date:            time.Now(),
		AssetType:       AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: TransactionTypeBuy,
		Quantity:        100,
		Price:           500.0,
		Amount:          50000.0,
		Fee:             &fee,
		Tax:             &tax,
		Currency:        CurrencyTWD,
		Note:            &note,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 驗證所有欄位都能正確設定
	assert.NotEqual(t, uuid.Nil, transaction.ID)
	assert.Equal(t, AssetTypeTWStock, transaction.AssetType)
	assert.Equal(t, "2330", transaction.Symbol)
	assert.Equal(t, "台積電", transaction.Name)
	assert.Equal(t, TransactionTypeBuy, transaction.TransactionType)
	assert.Equal(t, 100.0, transaction.Quantity)
	assert.Equal(t, 500.0, transaction.Price)
	assert.Equal(t, 50000.0, transaction.Amount)
	assert.NotNil(t, transaction.Fee)
	assert.Equal(t, 10.5, *transaction.Fee)
	assert.NotNil(t, transaction.Tax)
	assert.Equal(t, 15.0, *transaction.Tax)
	assert.Equal(t, CurrencyTWD, transaction.Currency)
	assert.NotNil(t, transaction.Note)
	assert.Equal(t, "test note", *transaction.Note)
}

// TestCreateTransactionInputWithTax 測試建立交易輸入包含 tax 欄位
func TestCreateTransactionInputWithTax(t *testing.T) {
	fee := 10.5
	tax := 15.0
	note := "test note"

	input := CreateTransactionInput{
		Date:            time.Now(),
		AssetType:       AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: TransactionTypeSell,
		Quantity:        100,
		Price:           500.0,
		Amount:          50000.0,
		Fee:             &fee,
		Tax:             &tax,
		Currency:        CurrencyTWD,
		Note:            &note,
	}

	// 驗證 tax 欄位可以正確設定
	assert.NotNil(t, input.Tax)
	assert.Equal(t, 15.0, *input.Tax)
}

// TestCreateTransactionInputWithoutTax 測試建立交易輸入不包含 tax 欄位
func TestCreateTransactionInputWithoutTax(t *testing.T) {
	input := CreateTransactionInput{
		Date:            time.Now(),
		AssetType:       AssetTypeUSStock,
		Symbol:          "AAPL",
		Name:            "Apple Inc.",
		TransactionType: TransactionTypeBuy,
		Quantity:        10,
		Price:           150.0,
		Amount:          1500.0,
		Currency:        CurrencyUSD,
	}

	// 驗證 tax 欄位可以為 nil
	assert.Nil(t, input.Tax)
}

// TestUpdateTransactionInputWithTax 測試更新交易輸入包含 tax 欄位
func TestUpdateTransactionInputWithTax(t *testing.T) {
	tax := 20.0

	input := UpdateTransactionInput{
		Tax: &tax,
	}

	// 驗證 tax 欄位可以正確設定
	assert.NotNil(t, input.Tax)
	assert.Equal(t, 20.0, *input.Tax)
}

// TestUpdateTransactionInputClearTax 測試更新交易時清除 tax 欄位
func TestUpdateTransactionInputClearTax(t *testing.T) {
	zero := 0.0

	input := UpdateTransactionInput{
		Tax: &zero,
	}

	// 驗證 tax 欄位可以設為 0
	assert.NotNil(t, input.Tax)
	assert.Equal(t, 0.0, *input.Tax)
}

