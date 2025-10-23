package service

import (
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// PriceService 價格服務介面
type PriceService interface {
	// GetPrice 取得單一標的價格
	GetPrice(symbol string, assetType models.AssetType) (*models.Price, error)

	// GetPrices 批次取得多個標的價格
	GetPrices(symbols []string, assetTypes map[string]models.AssetType) (map[string]*models.Price, error)

	// RefreshPrice 手動更新價格（清除快取並重新取得）
	RefreshPrice(symbol string, assetType models.AssetType) (*models.Price, error)
}

// mockPriceService Mock 價格服務（暫時使用固定價格）
// 後續會替換成真實的價格 API 整合
type mockPriceService struct {
	// 固定價格表（用於測試和開發）
	mockPrices map[string]float64
}

// NewMockPriceService 建立 Mock 價格服務
func NewMockPriceService() PriceService {
	return &mockPriceService{
		mockPrices: map[string]float64{
			// 台股
			"2330": 620.0,  // 台積電
			"2317": 110.0,  // 鴻海
			"2454": 1050.0, // 聯發科
			"2412": 95.0,   // 中華電

			// 美股
			"AAPL": 175.0,  // Apple
			"GOOGL": 140.0, // Google
			"MSFT": 380.0,  // Microsoft
			"TSLA": 250.0,  // Tesla

			// 加密貨幣
			"BTC": 1200000.0, // Bitcoin (TWD)
			"ETH": 60000.0,   // Ethereum (TWD)
			"USDT": 31.5,     // USDT (TWD)
		},
	}
}

// GetPrice 取得單一標的價格
func (s *mockPriceService) GetPrice(symbol string, assetType models.AssetType) (*models.Price, error) {
	price, exists := s.mockPrices[symbol]
	if !exists {
		// 如果沒有預設價格，返回一個合理的預設值
		price = 100.0
	}

	currency := "TWD"
	if assetType == models.AssetTypeUSStock {
		currency = "USD"
	}

	return &models.Price{
		Symbol:    symbol,
		AssetType: assetType,
		Price:     price,
		Currency:  currency,
		Source:    "mock",
		UpdatedAt: time.Now(),
	}, nil
}

// GetPrices 批次取得多個標的價格
func (s *mockPriceService) GetPrices(symbols []string, assetTypes map[string]models.AssetType) (map[string]*models.Price, error) {
	prices := make(map[string]*models.Price)

	for _, symbol := range symbols {
		assetType, exists := assetTypes[symbol]
		if !exists {
			return nil, fmt.Errorf("asset type not found for symbol: %s", symbol)
		}

		price, err := s.GetPrice(symbol, assetType)
		if err != nil {
			return nil, fmt.Errorf("failed to get price for %s: %w", symbol, err)
		}

		prices[symbol] = price
	}

	return prices, nil
}

// RefreshPrice 手動更新價格
func (s *mockPriceService) RefreshPrice(symbol string, assetType models.AssetType) (*models.Price, error) {
	// Mock 實作：直接返回當前價格
	return s.GetPrice(symbol, assetType)
}

