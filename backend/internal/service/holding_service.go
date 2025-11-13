package service

import (
	"fmt"
	"log"

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
	log.Println("[DEBUG] HoldingService.GetAllHoldings started")

	// 1. 從 Repository 取得交易記錄
	txFilters := repository.TransactionFilters{
		AssetType: filters.AssetType,
		Symbol:    filters.Symbol,
	}

	log.Println("[DEBUG] Step 1: Fetching transactions from repository...")
	transactions, err := s.transactionRepo.GetAll(txFilters)
	if err != nil {
		log.Printf("[ERROR] Failed to get transactions: %v", err)
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	log.Printf("[DEBUG] Step 1: Got %d transactions", len(transactions))

	// 如果沒有交易記錄，返回空列表
	if len(transactions) == 0 {
		log.Println("[DEBUG] No transactions found, returning empty list")
		return []*models.Holding{}, nil
	}

	// 2. 使用 FIFO Calculator 計算持倉
	log.Println("[DEBUG] Step 2: Calculating holdings using FIFO...")
	holdingsMap, err := s.fifoCalculator.CalculateAllHoldings(transactions)
	if err != nil {
		log.Printf("[ERROR] Failed to calculate holdings: %v", err)
		return nil, fmt.Errorf("failed to calculate holdings: %w", err)
	}
	log.Printf("[DEBUG] Step 2: Calculated %d holdings", len(holdingsMap))

	// 如果沒有持倉，返回空列表
	if len(holdingsMap) == 0 {
		log.Println("[DEBUG] No holdings after FIFO calculation, returning empty list")
		return []*models.Holding{}, nil
	}

	// 3. 準備批次取得價格
	log.Println("[DEBUG] Step 3: Preparing to fetch prices...")
	symbols := make([]string, 0, len(holdingsMap))
	assetTypes := make(map[string]models.AssetType)

	for symbol, holding := range holdingsMap {
		symbols = append(symbols, symbol)
		assetTypes[symbol] = holding.AssetType
	}
	log.Printf("[DEBUG] Step 3: Need to fetch prices for %d symbols: %v", len(symbols), symbols)

	// 4. 批次取得價格（採用優雅降級策略）
	log.Println("[DEBUG] Step 4: Fetching prices from price service...")
	prices, err := s.priceService.GetPrices(symbols, assetTypes)
	if err != nil {
		// 如果完全無法取得價格，記錄警告但繼續處理（使用成本價作為市值）
		log.Printf("[WARNING] Failed to get prices: %v", err)
		prices = make(map[string]*models.Price) // 空的價格 map
	} else {
		log.Printf("[DEBUG] Step 4: Successfully fetched prices for %d symbols", len(prices))
	}

	// 5. 整合價格資訊並計算損益（統一轉換為 TWD）
	log.Println("[DEBUG] Step 5: Integrating prices and calculating P&L...")
	holdings := make([]*models.Holding, 0, len(holdingsMap))
	skippedCount := 0

	for symbol, holding := range holdingsMap {
		log.Printf("[DEBUG] Processing holding: %s (AssetType: %s, Quantity: %.4f)",
			symbol, holding.AssetType, holding.Quantity)

		price, exists := prices[symbol]
		if exists && price.Price > 0 {
			log.Printf("[DEBUG] Price found for %s: %.4f (Source: %s, IsStale: %v)",
				symbol, price.Price, price.Source, price.IsStale)

			// 有價格資訊且價格有效
			holding.CurrentPrice = price.Price

			// 根據資產類型決定幣別
			currency := s.getCurrencyForAssetType(holding.AssetType)
			holding.Currency = currency
			log.Printf("[DEBUG] Currency for %s: %s", symbol, currency)

			// 將價格轉換為 TWD
			log.Printf("[DEBUG] Converting price to TWD for %s...", symbol)
			priceTWD, err := s.exchangeRateService.ConvertToTWD(price.Price, currency, price.UpdatedAt)
			if err != nil {
				// ConvertToTWD 內部已有完整的 fallback 機制（最新匯率 → 預設匯率）
				// 理論上不應該會失敗，但為了安全起見還是處理錯誤
				log.Printf("[ERROR] Failed to convert %s to TWD for %s: %v", currency, symbol, err)
				log.Printf("[WARNING] Skipping holding %s due to conversion error", symbol)
				skippedCount++
				// 如果真的失敗，跳過這個持倉
				continue
			}
			log.Printf("[DEBUG] Converted price for %s: %.4f %s -> %.4f TWD",
				symbol, price.Price, currency, priceTWD)

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

			log.Printf("[DEBUG] Calculated P&L for %s: MarketValue=%.2f, UnrealizedPL=%.2f (%.2f%%)",
				symbol, holding.MarketValue, holding.UnrealizedPL, holding.UnrealizedPLPct)
		} else {
			log.Printf("[WARNING] No valid price for %s (exists: %v, price: %.4f)",
				symbol, exists, func() float64 { if exists { return price.Price } else { return 0 } }())

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

	log.Printf("[DEBUG] Step 5 completed: %d holdings processed, %d skipped", len(holdings), skippedCount)

	if len(holdings) == 0 && len(holdingsMap) > 0 {
		log.Printf("[ERROR] All holdings were skipped! Original count: %d, Final count: %d",
			len(holdingsMap), len(holdings))
		return nil, fmt.Errorf("all holdings were skipped due to conversion errors")
	}

	log.Printf("[DEBUG] GetAllHoldings completed successfully, returning %d holdings", len(holdings))
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
		// ConvertToTWD 內部已有完整的 fallback 機制（最新匯率 → 預設匯率）
		// 理論上不應該會失敗
		return nil, fmt.Errorf("failed to convert price to TWD for %s: %w", symbol, err)
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

