package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestFixInsufficientQuantity_Success 測試成功修復數量不足
func TestFixInsufficientQuantity_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	mockPriceService := new(MockPriceService)
	mockExchangeRateService := new(MockExchangeRateService)

	service := NewHoldingService(mockRepo, mockFIFOCalc, mockPriceService, mockExchangeRateService)

	symbol := "2330"
	currentHolding := 100.0

	// Mock 交易記錄
	existingTransactions := []*models.Transaction{
		{
			ID:              uuid.New(),
			Date:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          symbol,
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        50,
			Price:           500,
			Amount:          25000,
			Currency:        models.CurrencyTWD,
		},
	}

	// Mock FIFO 計算結果（當前只有 50 股）
	fifoResult := &FIFOCalculatorResult{
		Holdings: map[string]*models.Holding{
			symbol: {
				Symbol:   symbol,
				Quantity: 50,
			},
		},
		Warnings: []*models.Warning{},
	}

	// Mock 價格服務（返回當前價格）
	currentPrice := &models.Price{
		Symbol:    symbol,
		AssetType: models.AssetTypeTWStock,
		Price:     600,
		Currency:  "TWD",
		Source:    "api",
		UpdatedAt: time.Now(),
	}

	// Mock 新增的交易記錄
	newTransaction := &models.Transaction{
		ID:              uuid.New(),
		Date:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          symbol,
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        50, // 補足 50 股
		Price:           600,
		Amount:          30000,
		Currency:        models.CurrencyTWD,
	}

	// 設定 Mock 期望
	mockRepo.On("GetAll", mock.MatchedBy(func(filters repository.TransactionFilters) bool {
		return filters.Symbol != nil && *filters.Symbol == symbol
	})).Return(existingTransactions, nil)

	mockFIFOCalc.On("CalculateAllHoldings", existingTransactions).Return(fifoResult, nil)
	mockPriceService.On("GetPrice", symbol, models.AssetTypeTWStock).Return(currentPrice, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.CreateTransactionInput")).Return(newTransaction, nil)

	// Act
	input := &models.FixInsufficientQuantityInput{
		Symbol:         symbol,
		CurrentHolding: currentHolding,
	}

	result, err := service.FixInsufficientQuantity(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, symbol, result.Symbol)
	assert.Equal(t, 50.0, result.Quantity) // 補足的數量
	assert.Equal(t, 600.0, result.Price)   // 使用當前價格

	mockRepo.AssertExpectations(t)
	mockFIFOCalc.AssertExpectations(t)
	mockPriceService.AssertExpectations(t)
}

// TestFixInsufficientQuantity_WithEstimatedCost 測試使用估計成本（價格 API 失敗）
func TestFixInsufficientQuantity_WithEstimatedCost(t *testing.T) {
	// Arrange
	mockRepo := new(MockTransactionRepository)
	mockFIFOCalc := new(MockFIFOCalculator)
	mockPriceService := new(MockPriceService)
	mockExchangeRateService := new(MockExchangeRateService)

	service := NewHoldingService(mockRepo, mockFIFOCalc, mockPriceService, mockExchangeRateService)

	symbol := "2330"
	currentHolding := 100.0
	estimatedCost := 550.0

	// Mock 交易記錄
	existingTransactions := []*models.Transaction{
		{
			ID:              uuid.New(),
			Date:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          symbol,
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        50,
			Price:           500,
			Amount:          25000,
			Currency:        models.CurrencyTWD,
		},
	}

	// Mock FIFO 計算結果
	fifoResult := &FIFOCalculatorResult{
		Holdings: map[string]*models.Holding{
			symbol: {
				Symbol:   symbol,
				Quantity: 50,
			},
		},
		Warnings: []*models.Warning{},
	}

	// Mock 價格服務失敗
	mockRepo.On("GetAll", mock.Anything).Return(existingTransactions, nil)
	mockFIFOCalc.On("CalculateAllHoldings", existingTransactions).Return(fifoResult, nil)
	mockPriceService.On("GetPrice", symbol, models.AssetTypeTWStock).Return(nil, fmt.Errorf("price API failed"))
	mockRepo.On("Create", mock.AnythingOfType("*models.CreateTransactionInput")).Return(&models.Transaction{
		ID:       uuid.New(),
		Symbol:   symbol,
		Quantity: 50,
		Price:    estimatedCost,
	}, nil)

	// Act
	input := &models.FixInsufficientQuantityInput{
		Symbol:         symbol,
		CurrentHolding: currentHolding,
		EstimatedCost:  &estimatedCost,
	}

	result, err := service.FixInsufficientQuantity(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, estimatedCost, result.Price) // 使用估計成本

	mockRepo.AssertExpectations(t)
	mockFIFOCalc.AssertExpectations(t)
	mockPriceService.AssertExpectations(t)
}

