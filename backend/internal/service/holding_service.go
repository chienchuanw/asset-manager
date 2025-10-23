package service

import (
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
)

// HoldingService 持倉服務介面
type HoldingService interface {
	// GetAllHoldings 取得所有持倉
	GetAllHoldings(filters models.HoldingFilters) ([]*models.Holding, error)

	// GetHoldingBySymbol 取得單一標的持倉
	GetHoldingBySymbol(symbol string) (*models.Holding, error)
}

// holdingService 持倉服務實作
type holdingService struct {
	transactionRepo repository.TransactionRepository
	fifoCalculator  FIFOCalculator
	priceService    PriceService
}

// NewHoldingService 建立新的持倉服務
func NewHoldingService(
	transactionRepo repository.TransactionRepository,
	fifoCalculator FIFOCalculator,
	priceService PriceService,
) HoldingService {
	return &holdingService{
		transactionRepo: transactionRepo,
		fifoCalculator:  fifoCalculator,
		priceService:    priceService,
	}
}

// GetAllHoldings 取得所有持倉
func (s *holdingService) GetAllHoldings(filters models.HoldingFilters) ([]*models.Holding, error) {
	// 1. 從 Repository 取得交易記錄
	txFilters := repository.TransactionFilters{
		AssetType: filters.AssetType,
		Symbol:    filters.Symbol,
	}

	transactions, err := s.transactionRepo.GetAll(txFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// 如果沒有交易記錄，返回空列表
	if len(transactions) == 0 {
		return []*models.Holding{}, nil
	}

	// 2. 使用 FIFO Calculator 計算持倉
	holdingsMap, err := s.fifoCalculator.CalculateAllHoldings(transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate holdings: %w", err)
	}

	// 如果沒有持倉，返回空列表
	if len(holdingsMap) == 0 {
		return []*models.Holding{}, nil
	}

	// 3. 準備批次取得價格
	symbols := make([]string, 0, len(holdingsMap))
	assetTypes := make(map[string]models.AssetType)

	for symbol, holding := range holdingsMap {
		symbols = append(symbols, symbol)
		assetTypes[symbol] = holding.AssetType
	}

	// 4. 批次取得價格
	prices, err := s.priceService.GetPrices(symbols, assetTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to get prices: %w", err)
	}

	// 5. 整合價格資訊並計算損益
	holdings := make([]*models.Holding, 0, len(holdingsMap))
	for symbol, holding := range holdingsMap {
		price, exists := prices[symbol]
		if exists {
			holding.CurrentPrice = price.Price
			holding.MarketValue = holding.Quantity * price.Price
			holding.UnrealizedPL = holding.MarketValue - holding.TotalCost

			// 計算未實現損益百分比
			if holding.TotalCost > 0 {
				holding.UnrealizedPLPct = (holding.UnrealizedPL / holding.TotalCost) * 100
			}
		}

		holdings = append(holdings, holding)
	}

	return holdings, nil
}

// GetHoldingBySymbol 取得單一標的持倉
func (s *holdingService) GetHoldingBySymbol(symbol string) (*models.Holding, error) {
	// 1. 從 Repository 取得該標的的交易記錄
	symbolFilter := symbol
	txFilters := repository.TransactionFilters{
		Symbol: &symbolFilter,
	}

	transactions, err := s.transactionRepo.GetAll(txFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// 如果沒有交易記錄，返回錯誤
	if len(transactions) == 0 {
		return nil, fmt.Errorf("holding not found for symbol: %s", symbol)
	}

	// 2. 使用 FIFO Calculator 計算持倉
	holding, err := s.fifoCalculator.CalculateHoldingForSymbol(symbol, transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate holding: %w", err)
	}

	// 如果沒有持倉（已全部賣出），返回錯誤
	if holding == nil {
		return nil, fmt.Errorf("holding not found for symbol: %s (all sold)", symbol)
	}

	// 3. 取得價格
	price, err := s.priceService.GetPrice(symbol, holding.AssetType)
	if err != nil {
		return nil, fmt.Errorf("failed to get price: %w", err)
	}

	// 4. 整合價格資訊並計算損益
	holding.CurrentPrice = price.Price
	holding.MarketValue = holding.Quantity * price.Price
	holding.UnrealizedPL = holding.MarketValue - holding.TotalCost

	// 計算未實現損益百分比
	if holding.TotalCost > 0 {
		holding.UnrealizedPLPct = (holding.UnrealizedPL / holding.TotalCost) * 100
	}

	return holding, nil
}

