package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chienchuanw/asset-manager/internal/client"
	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/redis/go-redis/v9"
)

// ExchangeRateService 匯率服務介面
type ExchangeRateService interface {
	// GetRate 取得指定日期的匯率（優先從快取取得）
	GetRate(fromCurrency, toCurrency models.Currency, date time.Time) (float64, error)
	// GetTodayRate 取得今日匯率（優先從快取取得）
	GetTodayRate(fromCurrency, toCurrency models.Currency) (float64, error)
	// RefreshTodayRate 從 API 更新今日匯率
	RefreshTodayRate() error
	// ConvertToTWD 將金額轉換為 TWD
	ConvertToTWD(amount float64, currency models.Currency, date time.Time) (float64, error)
}

// exchangeRateService 匯率服務實作
type exchangeRateService struct {
	repo        repository.ExchangeRateRepository
	bankClient  *client.TaiwanBankClient
	redisClient *redis.Client
	ctx         context.Context
}

// NewExchangeRateService 建立新的匯率服務
func NewExchangeRateService(
	repo repository.ExchangeRateRepository,
	bankClient *client.TaiwanBankClient,
	redisClient *redis.Client,
) ExchangeRateService {
	return &exchangeRateService{
		repo:        repo,
		bankClient:  bankClient,
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// GetRate 取得指定日期的匯率
func (s *exchangeRateService) GetRate(fromCurrency, toCurrency models.Currency, date time.Time) (float64, error) {
	// 如果是相同幣別，匯率為 1
	if fromCurrency == toCurrency {
		return 1.0, nil
	}

	// 只支援 USD <-> TWD 轉換
	if !isValidCurrencyPair(fromCurrency, toCurrency) {
		return 0, fmt.Errorf("unsupported currency pair: %s -> %s", fromCurrency, toCurrency)
	}

	// 標準化為 USD -> TWD
	normalizedFrom, normalizedTo := normalizeCurrencyPair(fromCurrency, toCurrency)

	// 1. 嘗試從 Redis 快取取得
	cacheKey := fmt.Sprintf("exchange_rate:%s:%s:%s", normalizedFrom, normalizedTo, date.Format("2006-01-02"))
	if s.redisClient != nil {
		cachedRate, err := s.redisClient.Get(s.ctx, cacheKey).Result()
		if err == nil {
			var rate float64
			if err := json.Unmarshal([]byte(cachedRate), &rate); err == nil {
				log.Printf("Exchange rate cache hit: %s -> %s on %s = %.4f", fromCurrency, toCurrency, date.Format("2006-01-02"), rate)
				return s.adjustRateForPair(rate, fromCurrency, toCurrency), nil
			}
		}
	}

	// 2. 從資料庫取得
	dbRate, err := s.repo.GetByDate(normalizedFrom, normalizedTo, date)
	if err != nil {
		return 0, fmt.Errorf("failed to get exchange rate from database: %w", err)
	}

	// 如果資料庫中有資料，快取並回傳
	if dbRate != nil {
		if s.redisClient != nil {
			rateJSON, _ := json.Marshal(dbRate.Rate)
			s.redisClient.Set(s.ctx, cacheKey, rateJSON, 24*time.Hour)
		}
		return s.adjustRateForPair(dbRate.Rate, fromCurrency, toCurrency), nil
	}

	// 3. 如果是今日，從 API 取得
	today := time.Now().Truncate(24 * time.Hour)
	if date.Truncate(24*time.Hour).Equal(today) {
		if err := s.RefreshTodayRate(); err != nil {
			return 0, fmt.Errorf("failed to refresh today's rate: %w", err)
		}
		// 重新從資料庫取得
		dbRate, err = s.repo.GetByDate(normalizedFrom, normalizedTo, date)
		if err != nil {
			return 0, fmt.Errorf("failed to get exchange rate after refresh: %w", err)
		}
		if dbRate != nil {
			return s.adjustRateForPair(dbRate.Rate, fromCurrency, toCurrency), nil
		}
	}

	// 4. 如果都沒有，使用最新的匯率
	latestRate, err := s.repo.GetLatest(normalizedFrom, normalizedTo)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest exchange rate: %w", err)
	}
	if latestRate != nil {
		log.Printf("Using latest exchange rate for %s: %.4f (from %s)", date.Format("2006-01-02"), latestRate.Rate, latestRate.Date.Format("2006-01-02"))
		return s.adjustRateForPair(latestRate.Rate, fromCurrency, toCurrency), nil
	}

	return 0, fmt.Errorf("no exchange rate found for %s -> %s", fromCurrency, toCurrency)
}

// GetTodayRate 取得今日匯率
func (s *exchangeRateService) GetTodayRate(fromCurrency, toCurrency models.Currency) (float64, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return s.GetRate(fromCurrency, toCurrency, today)
}

// RefreshTodayRate 從 API 更新今日匯率
func (s *exchangeRateService) RefreshTodayRate() error {
	// 從台灣銀行 API 取得匯率
	rate, err := s.bankClient.GetUSDToTWDRate()
	if err != nil {
		return fmt.Errorf("failed to fetch USD/TWD rate from Taiwan Bank: %w", err)
	}

	// 儲存到資料庫
	today := time.Now().Truncate(24 * time.Hour)
	input := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         rate,
		Date:         today,
	}

	dbRate, err := s.repo.Upsert(input)
	if err != nil {
		return fmt.Errorf("failed to save exchange rate: %w", err)
	}

	// 更新 Redis 快取
	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("exchange_rate:USD:TWD:%s", today.Format("2006-01-02"))
		rateJSON, _ := json.Marshal(dbRate.Rate)
		s.redisClient.Set(s.ctx, cacheKey, rateJSON, 24*time.Hour)
	}

	log.Printf("Refreshed today's USD/TWD rate: %.4f", rate)
	return nil
}

// ConvertToTWD 將金額轉換為 TWD
func (s *exchangeRateService) ConvertToTWD(amount float64, currency models.Currency, date time.Time) (float64, error) {
	if currency == models.CurrencyTWD {
		return amount, nil
	}

	rate, err := s.GetRate(currency, models.CurrencyTWD, date)
	if err != nil {
		return 0, err
	}

	return amount * rate, nil
}

// isValidCurrencyPair 檢查是否為有效的幣別組合
func isValidCurrencyPair(from, to models.Currency) bool {
	return (from == models.CurrencyUSD && to == models.CurrencyTWD) ||
		(from == models.CurrencyTWD && to == models.CurrencyUSD)
}

// normalizeCurrencyPair 標準化幣別組合為 USD -> TWD
func normalizeCurrencyPair(from, to models.Currency) (models.Currency, models.Currency) {
	if from == models.CurrencyTWD && to == models.CurrencyUSD {
		return models.CurrencyUSD, models.CurrencyTWD
	}
	return from, to
}

// adjustRateForPair 根據幣別組合調整匯率
func (s *exchangeRateService) adjustRateForPair(rate float64, from, to models.Currency) float64 {
	// 如果是 TWD -> USD，需要取倒數
	if from == models.CurrencyTWD && to == models.CurrencyUSD {
		return 1.0 / rate
	}
	return rate
}

