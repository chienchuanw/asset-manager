package service

import (
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/external"
	"github.com/chienchuanw/asset-manager/internal/models"
)

// realPriceService 真實價格服務（整合多個外部 API）
type realPriceService struct {
	finmindClient       *external.FinMindClient
	coingeckoClient     *external.CoinGeckoClient
	alphaVantageClient  *external.AlphaVantageClient
}

// NewRealPriceService 建立真實價格服務
func NewRealPriceService(finmindAPIKey, coingeckoAPIKey, alphaVantageAPIKey string) PriceService {
	return &realPriceService{
		finmindClient:       external.NewFinMindClient(finmindAPIKey),
		coingeckoClient:     external.NewCoinGeckoClient(coingeckoAPIKey),
		alphaVantageClient:  external.NewAlphaVantageClient(alphaVantageAPIKey),
	}
}

// GetPrice 取得單一標的價格
func (s *realPriceService) GetPrice(symbol string, assetType models.AssetType) (*models.Price, error) {
	var price float64
	var currency string
	var err error

	switch assetType {
	case models.AssetTypeTWStock:
		// 台股：使用 FinMind API
		price, err = s.finmindClient.GetStockPrice(symbol)
		if err != nil {
			return nil, fmt.Errorf("failed to get TW stock price: %w", err)
		}
		currency = "TWD"

	case models.AssetTypeUSStock:
		// 美股：使用 Alpha Vantage API
		price, err = s.alphaVantageClient.GetStockPrice(symbol)
		if err != nil {
			return nil, fmt.Errorf("failed to get US stock price: %w", err)
		}
		currency = "USD"

	case models.AssetTypeCrypto:
		// 加密貨幣：使用 CoinGecko API（使用 USD 作為基準貨幣）
		price, err = s.coingeckoClient.GetCryptoPrice(symbol, "usd")
		if err != nil {
			return nil, fmt.Errorf("failed to get crypto price: %w", err)
		}
		currency = "USD"

	default:
		return nil, fmt.Errorf("unsupported asset type: %s", assetType)
	}

	return &models.Price{
		Symbol:    symbol,
		AssetType: assetType,
		Price:     price,
		Currency:  currency,
		Source:    "api",
		UpdatedAt: time.Now(),
	}, nil
}

// GetPrices 批次取得多個標的價格
func (s *realPriceService) GetPrices(symbols []string, assetTypes map[string]models.AssetType) (map[string]*models.Price, error) {
	// 按資產類型分組
	twStocks := []string{}
	usStocks := []string{}
	cryptos := []string{}

	for _, symbol := range symbols {
		assetType, exists := assetTypes[symbol]
		if !exists {
			return nil, fmt.Errorf("asset type not found for symbol: %s", symbol)
		}

		switch assetType {
		case models.AssetTypeTWStock:
			twStocks = append(twStocks, symbol)
		case models.AssetTypeUSStock:
			usStocks = append(usStocks, symbol)
		case models.AssetTypeCrypto:
			cryptos = append(cryptos, symbol)
		}
	}

	prices := make(map[string]*models.Price)

	// 批次取得台股價格
	if len(twStocks) > 0 {
		twPrices, err := s.finmindClient.GetMultipleStockPrices(twStocks)
		if err != nil {
			fmt.Printf("Warning: failed to get TW stock prices: %v\n", err)
		} else {
			for symbol, price := range twPrices {
				prices[symbol] = &models.Price{
					Symbol:    symbol,
					AssetType: models.AssetTypeTWStock,
					Price:     price,
					Currency:  "TWD",
					Source:    "api",
					UpdatedAt: time.Now(),
				}
			}
		}
	}

	// 批次取得美股價格
	// 注意：Alpha Vantage 免費版有速率限制（每分鐘 5 次），批次查詢會較慢
	if len(usStocks) > 0 {
		usPrices, err := s.alphaVantageClient.GetMultipleStockPrices(usStocks)
		if err != nil {
			fmt.Printf("Warning: failed to get US stock prices: %v\n", err)
		} else {
			for symbol, price := range usPrices {
				prices[symbol] = &models.Price{
					Symbol:    symbol,
					AssetType: models.AssetTypeUSStock,
					Price:     price,
					Currency:  "USD",
					Source:    "api",
					UpdatedAt: time.Now(),
				}
			}
		}
	}

	// 批次取得加密貨幣價格
	if len(cryptos) > 0 {
		cryptoPrices, err := s.coingeckoClient.GetMultipleCryptoPrices(cryptos, "usd")
		if err != nil {
			fmt.Printf("Warning: failed to get crypto prices: %v\n", err)
		} else {
			for symbol, price := range cryptoPrices {
				prices[symbol] = &models.Price{
					Symbol:    symbol,
					AssetType: models.AssetTypeCrypto,
					Price:     price,
					Currency:  "USD",
					Source:    "api",
					UpdatedAt: time.Now(),
				}
			}
		}
	}

	// 檢查是否所有標的都有價格
	if len(prices) == 0 {
		return nil, fmt.Errorf("failed to get any prices")
	}

	return prices, nil
}

// RefreshPrice 手動更新價格（直接從 API 取得）
func (s *realPriceService) RefreshPrice(symbol string, assetType models.AssetType) (*models.Price, error) {
	// 與 GetPrice 相同，因為每次都是從 API 取得最新價格
	return s.GetPrice(symbol, assetType)
}

