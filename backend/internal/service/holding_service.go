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
	transactionRepo   repository.TransactionRepository
	fifoCalculator    FIFOCalculator
	priceService      PriceService
	exchangeRateService ExchangeRateService
}

// NewHoldingService 建立新的持倉服務
func NewHoldingService(
	transactionRepo repository.TransactionRepository,
	fifoCalculator FIFOCalculator,
	priceService PriceService,
	exchangeRateService ExchangeRateService,
) HoldingService {
	return &holdingService{
		transactionRepo:   transactionRepo,
		fifoCalculator:    fifoCalculator,
		priceService:      priceService,
		exchangeRateService: exchangeRateService,
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

	// 4. 批次取得價格（採用優雅降級策略）
	prices, err := s.priceService.GetPrices(symbols, assetTypes)
	if err != nil {
		// 如果完全無法取得價格，記錄警告但繼續處理（使用成本價作為市值）
		fmt.Printf("Warning: failed to get prices: %v\n", err)
		prices = make(map[string]*models.Price) // 空的價格 map
	}

	// 5. 整合價格資訊並計算損益（統一轉換為 TWD）
	holdings := make([]*models.Holding, 0, len(holdingsMap))
	for symbol, holding := range holdingsMap {
		price, exists := prices[symbol]
		if exists && price.Price > 0 {
			// 有價格資訊且價格有效
			holding.CurrentPrice = price.Price

			// 根據資產類型決定幣別
			currency := s.getCurrencyForAssetType(holding.AssetType)
			holding.Currency = currency

			// 將價格轉換為 TWD
			priceTWD, err := s.exchangeRateService.ConvertToTWD(price.Price, currency, price.UpdatedAt)
			if err != nil {
				// 如果匯率轉換失敗，使用原始價格（假設為 TWD）
				priceTWD = price.Price
			}
			holding.CurrentPriceTWD = priceTWD

			// 計算市值（TWD）
			holding.MarketValue = holding.Quantity * priceTWD
			holding.UnrealizedPL = holding.MarketValue - holding.TotalCost

			// 計算未實現損益百分比
			if holding.TotalCost > 0 {
				holding.UnrealizedPLPct = (holding.UnrealizedPL / holding.TotalCost) * 100
			}

			// 傳遞價格來源資訊
			holding.PriceSource = price.Source
			holding.IsPriceStale = price.IsStale
			holding.PriceStaleReason = price.StaleReason
		} else {
			// 無價格資訊或價格為 0，使用成本價作為市值（保守估計）
			holding.CurrentPrice = 0
			holding.CurrentPriceTWD = 0
			holding.MarketValue = holding.TotalCost // 使用成本價作為市值
			holding.UnrealizedPL = 0
			holding.UnrealizedPLPct = 0
			holding.PriceSource = "unavailable"
			holding.IsPriceStale = true
			holding.PriceStaleReason = "Price not available"
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

	// 4. 整合價格資訊並計算損益（統一轉換為 TWD）
	holding.CurrentPrice = price.Price

	// 根據資產類型決定幣別
	currency := s.getCurrencyForAssetType(holding.AssetType)
	holding.Currency = currency

	// 將價格轉換為 TWD
	priceTWD, err := s.exchangeRateService.ConvertToTWD(price.Price, currency, price.UpdatedAt)
	if err != nil {
		// 如果匯率轉換失敗，使用原始價格（假設為 TWD）
		priceTWD = price.Price
	}
	holding.CurrentPriceTWD = priceTWD

	// 計算市值（TWD）
	holding.MarketValue = holding.Quantity * priceTWD
	holding.UnrealizedPL = holding.MarketValue - holding.TotalCost

	// 計算未實現損益百分比
	if holding.TotalCost > 0 {
		holding.UnrealizedPLPct = (holding.UnrealizedPL / holding.TotalCost) * 100
	}

	// 傳遞價格來源資訊
	holding.PriceSource = price.Source
	holding.IsPriceStale = price.IsStale
	holding.PriceStaleReason = price.StaleReason

	return holding, nil
}

// getCurrencyForAssetType 根據資產類型取得幣別
func (s *holdingService) getCurrencyForAssetType(assetType models.AssetType) models.Currency {
	switch assetType {
	case models.AssetTypeTWStock:
		return models.CurrencyTWD
	case models.AssetTypeUSStock:
		return models.CurrencyUSD
	case models.AssetTypeCrypto:
		return models.CurrencyUSD // 加密貨幣使用 USD 計價
	default:
		return models.CurrencyTWD
	}
}

