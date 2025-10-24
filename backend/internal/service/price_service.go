package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/chienchuanw/asset-manager/internal/cache"
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

// ==================== Cached Price Service ====================

// cachedPriceService 帶 Redis 快取的價格服務
type cachedPriceService struct {
	cache               *cache.RedisCache
	fallback            PriceService // 當快取失效時使用的備用服務（例如 Mock 或真實 API）
	defaultExpiration   time.Duration
	usStockExpiration   time.Duration // 美股專用快取時間（較長，避免 API 限制）
}

// NewCachedPriceService 建立帶快取的價格服務
func NewCachedPriceService(redisCache *cache.RedisCache, fallback PriceService, cacheExpiration time.Duration) PriceService {
	return &cachedPriceService{
		cache:             redisCache,
		fallback:          fallback,
		defaultExpiration: cacheExpiration,
		usStockExpiration: 1 * time.Hour, // 美股快取 1 小時（避免 Alpha Vantage API 限制）
	}
}

// getCacheExpiration 根據資產類型取得快取過期時間
func (s *cachedPriceService) getCacheExpiration(assetType models.AssetType) time.Duration {
	if assetType == models.AssetTypeUSStock {
		return s.usStockExpiration
	}
	return s.defaultExpiration
}

// getCacheKey 產生快取 key
func (s *cachedPriceService) getCacheKey(symbol string, assetType models.AssetType) string {
	return fmt.Sprintf("price:%s:%s", assetType, symbol)
}

// GetPrice 取得單一標的價格（優先從快取取得）
func (s *cachedPriceService) GetPrice(symbol string, assetType models.AssetType) (*models.Price, error) {
	cacheKey := s.getCacheKey(symbol, assetType)

	// 1. 嘗試從快取取得
	var cachedPrice models.PriceCache
	cacheErr := s.cache.Get(cacheKey, &cachedPrice)
	if cacheErr == nil {
		// 快取命中，返回快取資料
		return &models.Price{
			Symbol:    cachedPrice.Symbol,
			AssetType: cachedPrice.AssetType,
			Price:     cachedPrice.Price,
			Currency:  cachedPrice.Currency,
			Source:    "cache",
			UpdatedAt: cachedPrice.CachedAt,
		}, nil
	}

	// 2. 快取未命中，從 fallback 服務取得
	price, err := s.fallback.GetPrice(symbol, assetType)
	if err != nil {
		// 檢查是否為 API rate limit 錯誤
		errMsg := strings.ToLower(err.Error())
		isRateLimit := strings.Contains(errMsg, "rate limit") ||
		               strings.Contains(errMsg, "api rate limit") ||
		               strings.Contains(errMsg, "standard api rate limit")

		// 如果是 rate limit 且有舊的快取資料，返回舊快取
		if isRateLimit && cacheErr == nil {
			fmt.Printf("Warning: API rate limit for %s, using stale cache\n", symbol)
			return &models.Price{
				Symbol:      cachedPrice.Symbol,
				AssetType:   cachedPrice.AssetType,
				Price:       cachedPrice.Price,
				Currency:    cachedPrice.Currency,
				Source:      "stale-cache",
				UpdatedAt:   cachedPrice.CachedAt,
				IsStale:     true,
				StaleReason: "API rate limit exceeded",
			}, nil
		}

		// 其他錯誤或沒有快取資料，返回錯誤
		return nil, fmt.Errorf("failed to get price from fallback: %w", err)
	}

	// 3. 儲存到快取（根據資產類型使用不同的過期時間）
	cacheData := models.PriceCache{
		Symbol:    price.Symbol,
		AssetType: price.AssetType,
		Price:     price.Price,
		Currency:  price.Currency,
		CachedAt:  time.Now(),
	}

	expiration := s.getCacheExpiration(assetType)
	if err := s.cache.Set(cacheKey, cacheData, expiration); err != nil {
		// 快取失敗不影響返回結果，只記錄錯誤
		fmt.Printf("Warning: failed to cache price for %s: %v\n", symbol, err)
	}

	// 更新 source 為 API
	price.Source = "api"
	return price, nil
}

// GetPrices 批次取得多個標的價格
func (s *cachedPriceService) GetPrices(symbols []string, assetTypes map[string]models.AssetType) (map[string]*models.Price, error) {
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

// RefreshPrice 手動更新價格（清除快取並重新取得）
func (s *cachedPriceService) RefreshPrice(symbol string, assetType models.AssetType) (*models.Price, error) {
	cacheKey := s.getCacheKey(symbol, assetType)

	// 1. 清除快取
	if err := s.cache.Delete(cacheKey); err != nil {
		fmt.Printf("Warning: failed to delete cache for %s: %v\n", symbol, err)
	}

	// 2. 從 fallback 服務取得最新價格
	price, err := s.fallback.GetPrice(symbol, assetType)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh price: %w", err)
	}

	// 3. 儲存到快取（根據資產類型使用不同的過期時間）
	cacheData := models.PriceCache{
		Symbol:    price.Symbol,
		AssetType: price.AssetType,
		Price:     price.Price,
		Currency:  price.Currency,
		CachedAt:  time.Now(),
	}

	expiration := s.getCacheExpiration(assetType)
	if err := s.cache.Set(cacheKey, cacheData, expiration); err != nil {
		fmt.Printf("Warning: failed to cache refreshed price for %s: %v\n", symbol, err)
	}

	price.Source = "api"
	return price, nil
}
