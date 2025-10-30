package service

import (
	"errors"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExchangeRateRepository 匯率 repository 的 mock
type MockExchangeRateRepository struct {
	mock.Mock
}

func (m *MockExchangeRateRepository) Create(input *models.ExchangeRateInput) (*models.ExchangeRate, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExchangeRate), args.Error(1)
}

func (m *MockExchangeRateRepository) GetByDate(fromCurrency, toCurrency models.Currency, date time.Time) (*models.ExchangeRate, error) {
	args := m.Called(fromCurrency, toCurrency, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExchangeRate), args.Error(1)
}

func (m *MockExchangeRateRepository) GetLatest(fromCurrency, toCurrency models.Currency) (*models.ExchangeRate, error) {
	args := m.Called(fromCurrency, toCurrency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExchangeRate), args.Error(1)
}

func (m *MockExchangeRateRepository) Upsert(input *models.ExchangeRateInput) (*models.ExchangeRate, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExchangeRate), args.Error(1)
}

// MockExchangeRateAPIClient ExchangeRate-API 客戶端的 mock
type MockExchangeRateAPIClient struct {
	mock.Mock
}

// 確保 MockExchangeRateAPIClient 實作 ExchangeRateAPIClient 介面
var _ ExchangeRateAPIClient = (*MockExchangeRateAPIClient)(nil)

func (m *MockExchangeRateAPIClient) GetUSDToTWDRate() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

// TestRefreshTodayRate_Success 測試成功更新今日匯率
func TestRefreshTodayRate_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockExchangeRateRepository)
	mockBankClient := new(MockExchangeRateAPIClient)

	// Mock ExchangeRate-API 回傳匯率
	expectedRate := 30.6
	mockBankClient.On("GetUSDToTWDRate").Return(expectedRate, nil)

	// Mock Repository Upsert 成功
	today := time.Now().Truncate(24 * time.Hour)
	expectedInput := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         expectedRate,
		Date:         today,
	}
	expectedResult := &models.ExchangeRate{
		ID:           1,
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         expectedRate,
		Date:         today,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	mockRepo.On("Upsert", mock.MatchedBy(func(input *models.ExchangeRateInput) bool {
		return input.FromCurrency == expectedInput.FromCurrency &&
			input.ToCurrency == expectedInput.ToCurrency &&
			input.Rate == expectedInput.Rate &&
			input.Date.Format("2006-01-02") == expectedInput.Date.Format("2006-01-02")
	})).Return(expectedResult, nil)

	// 建立 service（不使用 Redis）
	service := NewExchangeRateService(mockRepo, mockBankClient, nil)

	// Act
	err := service.RefreshTodayRate()

	// Assert
	assert.NoError(t, err)
	mockBankClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestRefreshTodayRate_BankAPIError 測試 ExchangeRate-API 錯誤
func TestRefreshTodayRate_BankAPIError(t *testing.T) {
	// Arrange
	mockRepo := new(MockExchangeRateRepository)
	mockBankClient := new(MockExchangeRateAPIClient)

	// Mock ExchangeRate-API 回傳錯誤
	mockBankClient.On("GetUSDToTWDRate").Return(0.0, errors.New("API error"))

	// 建立 service
	service := NewExchangeRateService(mockRepo, mockBankClient, nil)

	// Act
	err := service.RefreshTodayRate()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch USD/TWD rate")
	mockBankClient.AssertExpectations(t)
	// Repository 不應該被呼叫
	mockRepo.AssertNotCalled(t, "Upsert")
}

// TestRefreshTodayRate_RepositoryError 測試 Repository 錯誤
func TestRefreshTodayRate_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockExchangeRateRepository)
	mockBankClient := new(MockExchangeRateAPIClient)

	// Mock ExchangeRate-API 成功
	mockBankClient.On("GetUSDToTWDRate").Return(30.6, nil)

	// Mock Repository Upsert 失敗
	mockRepo.On("Upsert", mock.Anything).Return(nil, errors.New("database error"))

	// 建立 service
	service := NewExchangeRateService(mockRepo, mockBankClient, nil)

	// Act
	err := service.RefreshTodayRate()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save exchange rate")
	mockBankClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestGetTodayRate_FromDatabase 測試從資料庫取得今日匯率
func TestGetTodayRate_FromDatabase(t *testing.T) {
	// Arrange
	mockRepo := new(MockExchangeRateRepository)
	mockBankClient := new(MockExchangeRateAPIClient)

	today := time.Now().Truncate(24 * time.Hour)
	expectedRate := &models.ExchangeRate{
		ID:           1,
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         31.5,
		Date:         today,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Mock Repository 回傳今日匯率
	mockRepo.On("GetByDate", models.CurrencyUSD, models.CurrencyTWD, mock.MatchedBy(func(date time.Time) bool {
		return date.Format("2006-01-02") == today.Format("2006-01-02")
	})).Return(expectedRate, nil)

	// 建立 service
	service := NewExchangeRateService(mockRepo, mockBankClient, nil)

	// Act
	rate, err := service.GetTodayRate(models.CurrencyUSD, models.CurrencyTWD)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedRate.Rate, rate)
	mockRepo.AssertExpectations(t)
	// 不應該呼叫 Bank API
	mockBankClient.AssertNotCalled(t, "GetUSDToTWDRate")
}

// TestGetTodayRate_RefreshWhenNotFound 測試資料庫沒有今日匯率時自動更新
func TestGetTodayRate_RefreshWhenNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockExchangeRateRepository)
	mockBankClient := new(MockExchangeRateAPIClient)

	today := time.Now().Truncate(24 * time.Hour)

	// Mock Repository 第一次查詢回傳 nil（沒有資料）
	mockRepo.On("GetByDate", models.CurrencyUSD, models.CurrencyTWD, mock.MatchedBy(func(date time.Time) bool {
		return date.Format("2006-01-02") == today.Format("2006-01-02")
	})).Return(nil, nil).Once()

	// Mock ExchangeRate-API 回傳匯率
	expectedRate := 30.6
	mockBankClient.On("GetUSDToTWDRate").Return(expectedRate, nil)

	// Mock Repository Upsert 成功
	upsertResult := &models.ExchangeRate{
		ID:           1,
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         expectedRate,
		Date:         today,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	mockRepo.On("Upsert", mock.Anything).Return(upsertResult, nil)

	// Mock Repository 第二次查詢回傳更新後的匯率
	mockRepo.On("GetByDate", models.CurrencyUSD, models.CurrencyTWD, mock.MatchedBy(func(date time.Time) bool {
		return date.Format("2006-01-02") == today.Format("2006-01-02")
	})).Return(upsertResult, nil).Once()

	// 建立 service
	service := NewExchangeRateService(mockRepo, mockBankClient, nil)

	// Act
	rate, err := service.GetTodayRate(models.CurrencyUSD, models.CurrencyTWD)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedRate, rate)
	mockRepo.AssertExpectations(t)
	mockBankClient.AssertExpectations(t)
}

