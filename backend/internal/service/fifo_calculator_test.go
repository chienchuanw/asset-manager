package service

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
)

// ==================== 基本 FIFO 測試 ====================

// TestFIFO_SingleBuy 測試單次買入，無賣出
func TestFIFO_SingleBuy(t *testing.T) {
	// Arrange: 準備測試資料
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28), // 手續費 28 元
		},
	}

	calculator := NewFIFOCalculator()

	// Act: 執行計算
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert: 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, holding)
	assert.Equal(t, "2330", holding.Symbol)
	assert.Equal(t, "台積電", holding.Name)
	assert.Equal(t, models.AssetTypeTWStock, holding.AssetType)
	assert.Equal(t, 100.0, holding.Quantity)

	// 平均成本 = (50000 + 28) / 100 = 500.28
	assert.InDelta(t, 500.28, holding.AvgCost, 0.01)
	assert.InDelta(t, 50028.0, holding.TotalCost, 0.01)
}

// TestFIFO_MultipleBuys 測試多次買入，無賣出
func TestFIFO_MultipleBuys(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        50,
			Price:           520,
			Amount:          26000,
			Fee:             ptrFloat64(15),
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 150.0, holding.Quantity)

	// 總成本 = (50000 + 28) + (26000 + 15) = 76043
	// 平均成本 = 76043 / 150 = 506.95
	assert.InDelta(t, 506.95, holding.AvgCost, 0.01)
	assert.InDelta(t, 76043.0, holding.TotalCost, 0.01)
}

// TestFIFO_BuyThenPartialSell 測試買入後部分賣出
func TestFIFO_BuyThenPartialSell(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeSell,
			Quantity:        30,
			Price:           550,
			Amount:          16500,
			Fee:             ptrFloat64(10),
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 70.0, holding.Quantity)

	// 賣出 30 股不影響剩餘 70 股的成本
	// 剩餘成本 = 500.28 * 70 = 35019.6
	assert.InDelta(t, 500.28, holding.AvgCost, 0.01)
	assert.InDelta(t, 35019.6, holding.TotalCost, 0.01)
}

// TestFIFO_BuyThenFullSell 測試買入後全部賣出
func TestFIFO_BuyThenFullSell(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeSell,
			Quantity:        100,
			Price:           550,
			Amount:          55000,
			Fee:             ptrFloat64(30),
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, holding) // 全部賣出後應該沒有持倉
}

// TestFIFO_MultipleBuysAndSells 測試多次買入、多次賣出（跨批次）
func TestFIFO_MultipleBuysAndSells(t *testing.T) {
	// Arrange
	// 情境：
	// 1/1: 買入 100 股 @ 500 (費用 28)
	// 1/5: 買入 50 股 @ 520 (費用 15)
	// 1/10: 賣出 120 股 @ 550 (費用 30)
	//
	// FIFO 邏輯：
	// - 賣出 100 股來自第一批 (500.28)
	// - 賣出 20 股來自第二批 (520.30)
	// - 剩餘 30 股來自第二批 (520.30)

	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        50,
			Price:           520,
			Amount:          26000,
			Fee:             ptrFloat64(15),
		},
		{
			Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeSell,
			Quantity:        120,
			Price:           550,
			Amount:          66000,
			Fee:             ptrFloat64(30),
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 30.0, holding.Quantity)

	// 剩餘 30 股來自第二批
	// 第二批單位成本 = (26000 + 15) / 50 = 520.30
	assert.InDelta(t, 520.30, holding.AvgCost, 0.01)
	assert.InDelta(t, 15609.0, holding.TotalCost, 0.01)
}

// ==================== 手續費測試 ====================

// TestFIFO_BuyFeeIncludedInCost 測試買入手續費計入成本
func TestFIFO_BuyFeeIncludedInCost(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(100), // 較高的手續費
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert
	assert.NoError(t, err)
	// 平均成本應包含手續費: (50000 + 100) / 100 = 501
	assert.InDelta(t, 501.0, holding.AvgCost, 0.01)
}

// TestFIFO_SellFeeNotAffectRemainingCost 測試賣出手續費不影響剩餘持倉成本
func TestFIFO_SellFeeNotAffectRemainingCost(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeSell,
			Quantity:        50,
			Price:           550,
			Amount:          27500,
			Fee:             ptrFloat64(200), // 很高的賣出手續費
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert
	assert.NoError(t, err)
	// 賣出手續費不應影響剩餘 50 股的成本
	// 剩餘成本仍然是 500.28
	assert.InDelta(t, 500.28, holding.AvgCost, 0.01)
}

// ==================== 股利測試 ====================

// TestFIFO_DividendNotAffectCost 測試股利不影響持倉成本
func TestFIFO_DividendNotAffectCost(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeDividend,
			Quantity:        0,    // 股利不影響數量
			Price:           0,
			Amount:          5000, // 收到 5000 元股利
			Fee:             nil,
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 100.0, holding.Quantity)
	// 股利不應影響成本
	assert.InDelta(t, 500.28, holding.AvgCost, 0.01)
	assert.InDelta(t, 50028.0, holding.TotalCost, 0.01)
}

// ==================== 邊界情況測試 ====================

// TestFIFO_EmptyTransactions 測試空交易記錄
func TestFIFO_EmptyTransactions(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{}
	calculator := NewFIFOCalculator()

	// Act
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, holding) // 沒有交易記錄應該返回 nil
}

// TestFIFO_SellMoreThanHolding 測試賣出數量大於持有（錯誤處理）
func TestFIFO_SellMoreThanHolding(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeSell,
			Quantity:        150, // 賣出超過持有數量
			Price:           550,
			Amount:          82500,
			Fee:             ptrFloat64(30),
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, holding)
	assert.Contains(t, err.Error(), "insufficient quantity")
}

// TestFIFO_SameDayMultipleTransactions 測試同一天多筆交易
func TestFIFO_SameDayMultipleTransactions(t *testing.T) {
	// Arrange
	sameDay := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	transactions := []*models.Transaction{
		{
			Date:            sameDay,
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            sameDay,
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        50,
			Price:           510,
			Amount:          25500,
			Fee:             ptrFloat64(15),
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holding, err := calculator.CalculateHoldingForSymbol("2330", transactions)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 150.0, holding.Quantity)
	// 總成本 = (50000 + 28) + (25500 + 15) = 75543
	// 平均成本 = 75543 / 150 = 503.62
	assert.InDelta(t, 503.62, holding.AvgCost, 0.01)
}

// ==================== 多標的測試 ====================

// TestFIFO_CalculateAllHoldings 測試計算所有標的持倉
func TestFIFO_CalculateAllHoldings(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		// 台積電交易
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		// Apple 交易
		{
			Date:            time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeUSStock,
			Symbol:          "AAPL",
			Name:            "Apple Inc.",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        50,
			Price:           150,
			Amount:          7500,
			Fee:             ptrFloat64(10),
		},
		// Bitcoin 交易
		{
			Date:            time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeCrypto,
			Symbol:          "BTC",
			Name:            "Bitcoin",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        0.5,
			Price:           900000,
			Amount:          450000,
			Fee:             ptrFloat64(100),
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holdings, err := calculator.CalculateAllHoldings(transactions)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 3, len(holdings))

	// 驗證台積電
	tsmc := holdings["2330"]
	assert.NotNil(t, tsmc)
	assert.Equal(t, 100.0, tsmc.Quantity)
	assert.InDelta(t, 500.28, tsmc.AvgCost, 0.01)

	// 驗證 Apple
	apple := holdings["AAPL"]
	assert.NotNil(t, apple)
	assert.Equal(t, 50.0, apple.Quantity)
	assert.InDelta(t, 150.20, apple.AvgCost, 0.01)

	// 驗證 Bitcoin
	btc := holdings["BTC"]
	assert.NotNil(t, btc)
	assert.Equal(t, 0.5, btc.Quantity)
	assert.InDelta(t, 900200.0, btc.AvgCost, 0.01)
}

// TestFIFO_CalculateAllHoldings_WithSoldOut 測試包含已賣出標的
func TestFIFO_CalculateAllHoldings_WithSoldOut(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		// 台積電：持有
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		// 鴻海：買入後全部賣出
		{
			Date:            time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2317",
			Name:            "鴻海",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        50,
			Price:           100,
			Amount:          5000,
			Fee:             ptrFloat64(10),
		},
		{
			Date:            time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2317",
			Name:            "鴻海",
			TransactionType: models.TransactionTypeSell,
			Quantity:        50,
			Price:           110,
			Amount:          5500,
			Fee:             ptrFloat64(10),
		},
	}

	calculator := NewFIFOCalculator()

	// Act
	holdings, err := calculator.CalculateAllHoldings(transactions)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(holdings)) // 只有台積電有持倉

	// 驗證台積電
	tsmc := holdings["2330"]
	assert.NotNil(t, tsmc)
	assert.Equal(t, 100.0, tsmc.Quantity)

	// 鴻海應該不在持倉中
	_, exists := holdings["2317"]
	assert.False(t, exists)
}

// ==================== CalculateCostBasis 測試 ====================

// TestCalculateCostBasis_SingleBatch 測試單一批次賣出
func TestCalculateCostBasis_SingleBatch(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
	}

	sellTransaction := &models.Transaction{
		Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeSell,
		Quantity:        30,
		Price:           550,
		Amount:          16500,
		Fee:             ptrFloat64(10),
	}

	calculator := NewFIFOCalculator()

	// Act
	costBasis, err := calculator.CalculateCostBasis("2330", sellTransaction, transactions)

	// Assert
	assert.NoError(t, err)
	// 成本基礎 = 30 股 × 500.28 = 15008.4
	assert.InDelta(t, 15008.4, costBasis, 0.01)
}

// TestCalculateCostBasis_MultipleBatches 測試跨多個批次賣出
func TestCalculateCostBasis_MultipleBatches(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        50,
			Price:           520,
			Amount:          26000,
			Fee:             ptrFloat64(15),
		},
	}

	sellTransaction := &models.Transaction{
		Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeSell,
		Quantity:        120,
		Price:           550,
		Amount:          66000,
		Fee:             ptrFloat64(30),
	}

	calculator := NewFIFOCalculator()

	// Act
	costBasis, err := calculator.CalculateCostBasis("2330", sellTransaction, transactions)

	// Assert
	assert.NoError(t, err)
	// FIFO: 賣出 100 股來自第一批 (500.28) + 20 股來自第二批 (520.30)
	// 成本基礎 = (100 × 500.28) + (20 × 520.30) = 50028 + 10406 = 60434
	assert.InDelta(t, 60434.0, costBasis, 0.01)
}

// TestCalculateCostBasis_WithPreviousSell 測試考慮之前的賣出交易
func TestCalculateCostBasis_WithPreviousSell(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
		{
			Date:            time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeSell,
			Quantity:        30,
			Price:           520,
			Amount:          15600,
			Fee:             ptrFloat64(10),
		},
		{
			Date:            time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        50,
			Price:           530,
			Amount:          26500,
			Fee:             ptrFloat64(15),
		},
	}

	sellTransaction := &models.Transaction{
		Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeSell,
		Quantity:        50,
		Price:           550,
		Amount:          27500,
		Fee:             ptrFloat64(20),
	}

	calculator := NewFIFOCalculator()

	// Act
	costBasis, err := calculator.CalculateCostBasis("2330", sellTransaction, transactions)

	// Assert
	assert.NoError(t, err)
	// 第一批買入 100 股 @ 500.28
	// 第一次賣出 30 股，剩餘 70 股 @ 500.28
	// 第二批買入 50 股 @ 530.30
	// 第二次賣出 50 股：全部來自第一批剩餘的 70 股
	// 成本基礎 = 50 × 500.28 = 25014
	assert.InDelta(t, 25014.0, costBasis, 0.01)
}

// TestCalculateCostBasis_InsufficientQuantity 測試賣出數量超過持有
func TestCalculateCostBasis_InsufficientQuantity(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
	}

	sellTransaction := &models.Transaction{
		Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeSell,
		Quantity:        150,
		Price:           550,
		Amount:          82500,
		Fee:             ptrFloat64(30),
	}

	calculator := NewFIFOCalculator()

	// Act
	costBasis, err := calculator.CalculateCostBasis("2330", sellTransaction, transactions)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0.0, costBasis)
	assert.Contains(t, err.Error(), "insufficient quantity")
}

// TestCalculateCostBasis_NotSellTransaction 測試非賣出交易
func TestCalculateCostBasis_NotSellTransaction(t *testing.T) {
	// Arrange
	transactions := []*models.Transaction{
		{
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        100,
			Price:           500,
			Amount:          50000,
			Fee:             ptrFloat64(28),
		},
	}

	buyTransaction := &models.Transaction{
		Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        50,
		Price:           520,
		Amount:          26000,
		Fee:             ptrFloat64(15),
	}

	calculator := NewFIFOCalculator()

	// Act
	costBasis, err := calculator.CalculateCostBasis("2330", buyTransaction, transactions)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0.0, costBasis)
	assert.Contains(t, err.Error(), "not a sell transaction")
}

// ==================== 輔助函式 ====================

// ptrFloat64 建立 float64 指標（方便測試）
func ptrFloat64(v float64) *float64 {
	return &v
}

