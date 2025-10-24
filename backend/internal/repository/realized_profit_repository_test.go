package repository

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRealizedProfitRepository_Create 測試建立已實現損益記錄
func TestRealizedProfitRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := NewRealizedProfitRepository(db)
	transactionRepo := NewTransactionRepository(db)

	// 先建立一筆賣出交易
	sellDate := time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC)
	transaction, err := transactionRepo.Create(&models.CreateTransactionInput{
		Date:            sellDate,
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeSell,
		Quantity:        100,
		Price:           620,
		Amount:          62000,
		Fee:             28,
		Currency:        "TWD",
	})
	require.NoError(t, err)

	// 建立已實現損益記錄
	input := &models.CreateRealizedProfitInput{
		TransactionID: transaction.ID,
		Symbol:        "2330",
		AssetType:     models.AssetTypeTWStock,
		SellDate:      sellDate,
		Quantity:      100,
		SellPrice:     620,
		SellAmount:    62000,
		SellFee:       28,
		CostBasis:     50000, // 假設成本基礎為 50000
		Currency:      "TWD",
	}

	result, err := repo.Create(input)

	// 驗證
	require.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, transaction.ID, result.TransactionID)
	assert.Equal(t, "2330", result.Symbol)
	assert.Equal(t, models.AssetTypeTWStock, result.AssetType)
	assert.Equal(t, 100.0, result.Quantity)
	assert.Equal(t, 620.0, result.SellPrice)
	assert.Equal(t, 62000.0, result.SellAmount)
	assert.Equal(t, 28.0, result.SellFee)
	assert.Equal(t, 50000.0, result.CostBasis)
	
	// 驗證已實現損益計算
	// 已實現損益 = (賣出金額 - 賣出手續費) - 成本基礎
	// = (62000 - 28) - 50000 = 11972
	expectedPL := 11972.0
	assert.Equal(t, expectedPL, result.RealizedPL)
	
	// 驗證已實現損益百分比
	// = (11972 / 50000) × 100 = 23.944%
	expectedPLPct := 23.944
	assert.InDelta(t, expectedPLPct, result.RealizedPLPct, 0.001)
	
	assert.Equal(t, "TWD", result.Currency)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
}

// TestRealizedProfitRepository_GetByTransactionID 測試根據交易 ID 取得已實現損益
func TestRealizedProfitRepository_GetByTransactionID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := NewRealizedProfitRepository(db)
	transactionRepo := NewTransactionRepository(db)

	// 建立交易和已實現損益
	sellDate := time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC)
	transaction, _ := transactionRepo.Create(&models.CreateTransactionInput{
		Date:            sellDate,
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeSell,
		Quantity:        100,
		Price:           620,
		Amount:          62000,
		Fee:             28,
		Currency:        "TWD",
	})

	created, _ := repo.Create(&models.CreateRealizedProfitInput{
		TransactionID: transaction.ID,
		Symbol:        "2330",
		AssetType:     models.AssetTypeTWStock,
		SellDate:      sellDate,
		Quantity:      100,
		SellPrice:     620,
		SellAmount:    62000,
		SellFee:       28,
		CostBasis:     50000,
		Currency:      "TWD",
	})

	// 測試：根據交易 ID 取得
	result, err := repo.GetByTransactionID(transaction.ID)

	// 驗證
	require.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, transaction.ID, result.TransactionID)
	assert.Equal(t, "2330", result.Symbol)
}

// TestRealizedProfitRepository_GetByTransactionID_NotFound 測試取得不存在的記錄
func TestRealizedProfitRepository_GetByTransactionID_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := NewRealizedProfitRepository(db)

	// 測試：取得不存在的交易 ID
	result, err := repo.GetByTransactionID("non-existent-id")

	// 驗證
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestRealizedProfitRepository_GetAll 測試取得所有已實現損益記錄
func TestRealizedProfitRepository_GetAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := NewRealizedProfitRepository(db)
	transactionRepo := NewTransactionRepository(db)

	// 建立多筆已實現損益記錄
	sellDate1 := time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC)
	sellDate2 := time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC)

	// 台股賣出
	tx1, _ := transactionRepo.Create(&models.CreateTransactionInput{
		Date:            sellDate1,
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeSell,
		Quantity:        100,
		Price:           620,
		Amount:          62000,
		Fee:             28,
		Currency:        "TWD",
	})
	repo.Create(&models.CreateRealizedProfitInput{
		TransactionID: tx1.ID,
		Symbol:        "2330",
		AssetType:     models.AssetTypeTWStock,
		SellDate:      sellDate1,
		Quantity:      100,
		SellPrice:     620,
		SellAmount:    62000,
		SellFee:       28,
		CostBasis:     50000,
		Currency:      "TWD",
	})

	// 美股賣出
	tx2, _ := transactionRepo.Create(&models.CreateTransactionInput{
		Date:            sellDate2,
		AssetType:       models.AssetTypeUSStock,
		Symbol:          "AAPL",
		Name:            "Apple Inc.",
		TransactionType: models.TransactionTypeSell,
		Quantity:        50,
		Price:           180,
		Amount:          9000,
		Fee:             5,
		Currency:        "USD",
	})
	repo.Create(&models.CreateRealizedProfitInput{
		TransactionID: tx2.ID,
		Symbol:        "AAPL",
		AssetType:     models.AssetTypeUSStock,
		SellDate:      sellDate2,
		Quantity:      50,
		SellPrice:     180,
		SellAmount:    9000,
		SellFee:       5,
		CostBasis:     8000,
		Currency:      "USD",
	})

	// 測試：取得所有記錄
	t.Run("GetAll_NoFilter", func(t *testing.T) {
		results, err := repo.GetAll(models.RealizedProfitFilters{})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	// 測試：按資產類型篩選
	t.Run("GetAll_FilterByAssetType", func(t *testing.T) {
		assetType := models.AssetTypeTWStock
		results, err := repo.GetAll(models.RealizedProfitFilters{
			AssetType: &assetType,
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "2330", results[0].Symbol)
	})

	// 測試：按標的代碼篩選
	t.Run("GetAll_FilterBySymbol", func(t *testing.T) {
		symbol := "AAPL"
		results, err := repo.GetAll(models.RealizedProfitFilters{
			Symbol: &symbol,
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "AAPL", results[0].Symbol)
	})

	// 測試：按時間範圍篩選
	t.Run("GetAll_FilterByDateRange", func(t *testing.T) {
		startDate := time.Date(2025, 10, 18, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC)
		results, err := repo.GetAll(models.RealizedProfitFilters{
			StartDate: &startDate,
			EndDate:   &endDate,
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "AAPL", results[0].Symbol)
	})
}

// TestRealizedProfitRepository_Delete 測試刪除已實現損益記錄
func TestRealizedProfitRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := NewRealizedProfitRepository(db)
	transactionRepo := NewTransactionRepository(db)

	// 建立記錄
	sellDate := time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC)
	transaction, _ := transactionRepo.Create(&models.CreateTransactionInput{
		Date:            sellDate,
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeSell,
		Quantity:        100,
		Price:           620,
		Amount:          62000,
		Fee:             28,
		Currency:        "TWD",
	})

	created, _ := repo.Create(&models.CreateRealizedProfitInput{
		TransactionID: transaction.ID,
		Symbol:        "2330",
		AssetType:     models.AssetTypeTWStock,
		SellDate:      sellDate,
		Quantity:      100,
		SellPrice:     620,
		SellAmount:    62000,
		SellFee:       28,
		CostBasis:     50000,
		Currency:      "TWD",
	})

	// 測試：刪除記錄
	err := repo.Delete(created.ID)
	require.NoError(t, err)

	// 驗證：記錄已被刪除
	result, err := repo.GetByTransactionID(transaction.ID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

