package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAssetSnapshotService 模擬 AssetSnapshotService
type MockAssetSnapshotService struct {
	mock.Mock
}

func (m *MockAssetSnapshotService) CreateSnapshot(input *models.CreateAssetSnapshotInput) (*models.AssetSnapshot, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssetSnapshot), args.Error(1)
}

func (m *MockAssetSnapshotService) GetSnapshotByDate(date time.Time, assetType models.SnapshotAssetType) (*models.AssetSnapshot, error) {
	args := m.Called(date, assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssetSnapshot), args.Error(1)
}

func (m *MockAssetSnapshotService) GetSnapshotsByDateRange(startDate, endDate time.Time, assetType models.SnapshotAssetType) ([]*models.AssetSnapshot, error) {
	args := m.Called(startDate, endDate, assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.AssetSnapshot), args.Error(1)
}

func (m *MockAssetSnapshotService) GetLatestSnapshot(assetType models.SnapshotAssetType) (*models.AssetSnapshot, error) {
	args := m.Called(assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssetSnapshot), args.Error(1)
}

func (m *MockAssetSnapshotService) UpdateSnapshot(date time.Time, assetType models.SnapshotAssetType, valueTWD float64) (*models.AssetSnapshot, error) {
	args := m.Called(date, assetType, valueTWD)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssetSnapshot), args.Error(1)
}

func (m *MockAssetSnapshotService) DeleteSnapshot(date time.Time, assetType models.SnapshotAssetType) error {
	args := m.Called(date, assetType)
	return args.Error(0)
}

func (m *MockAssetSnapshotService) CreateDailySnapshots() error {
	args := m.Called()
	return args.Error(0)
}

// TestNewSnapshotScheduler 測試建立排程器
func TestNewSnapshotScheduler(t *testing.T) {
	mockService := new(MockAssetSnapshotService)
	config := SnapshotSchedulerConfig{
		Enabled:           true,
		DailySnapshotTime: "23:59",
	}

	scheduler := NewSnapshotScheduler(mockService, config)

	assert.NotNil(t, scheduler)
	assert.Equal(t, true, scheduler.enabled)
	assert.Equal(t, "23:59", scheduler.dailySnapshotTime)
}

// TestSnapshotScheduler_Start_Disabled 測試停用的排程器
func TestSnapshotScheduler_Start_Disabled(t *testing.T) {
	mockService := new(MockAssetSnapshotService)
	config := SnapshotSchedulerConfig{
		Enabled:           false,
		DailySnapshotTime: "23:59",
	}

	scheduler := NewSnapshotScheduler(mockService, config)
	err := scheduler.Start()

	assert.NoError(t, err)
}

// TestSnapshotScheduler_Start_InvalidTime 測試無效的時間格式
func TestSnapshotScheduler_Start_InvalidTime(t *testing.T) {
	mockService := new(MockAssetSnapshotService)
	config := SnapshotSchedulerConfig{
		Enabled:           true,
		DailySnapshotTime: "invalid",
	}

	scheduler := NewSnapshotScheduler(mockService, config)
	err := scheduler.Start()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid daily snapshot time")
}

// TestSnapshotScheduler_RunNow_Success 測試立即執行成功
func TestSnapshotScheduler_RunNow_Success(t *testing.T) {
	mockService := new(MockAssetSnapshotService)
	config := SnapshotSchedulerConfig{
		Enabled:           true,
		DailySnapshotTime: "23:59",
	}

	mockService.On("CreateDailySnapshots").Return(nil)

	scheduler := NewSnapshotScheduler(mockService, config)
	err := scheduler.RunNow()

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

// TestSnapshotScheduler_RunNow_Error 測試立即執行失敗
func TestSnapshotScheduler_RunNow_Error(t *testing.T) {
	mockService := new(MockAssetSnapshotService)
	config := SnapshotSchedulerConfig{
		Enabled:           true,
		DailySnapshotTime: "23:59",
	}

	mockService.On("CreateDailySnapshots").Return(errors.New("database error"))

	scheduler := NewSnapshotScheduler(mockService, config)
	err := scheduler.RunNow()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create snapshots")
	mockService.AssertExpectations(t)
}

// TestParseTime 測試時間解析
func TestParseTime(t *testing.T) {
	tests := []struct {
		name        string
		timeStr     string
		wantHour    int
		wantMinute  int
		wantErr     bool
	}{
		{
			name:        "有效時間 - 23:59",
			timeStr:     "23:59",
			wantHour:    23,
			wantMinute:  59,
			wantErr:     false,
		},
		{
			name:        "有效時間 - 00:00",
			timeStr:     "00:00",
			wantHour:    0,
			wantMinute:  0,
			wantErr:     false,
		},
		{
			name:        "有效時間 - 12:30",
			timeStr:     "12:30",
			wantHour:    12,
			wantMinute:  30,
			wantErr:     false,
		},
		{
			name:        "無效時間 - 格式錯誤",
			timeStr:     "invalid",
			wantErr:     true,
		},
		{
			name:        "無效時間 - 超出範圍",
			timeStr:     "25:00",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hour, minute, err := parseTime(tt.timeStr)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantHour, hour)
				assert.Equal(t, tt.wantMinute, minute)
			}
		})
	}
}

// MockExchangeRateService 模擬 ExchangeRateService
type MockExchangeRateService struct {
	mock.Mock
}

func (m *MockExchangeRateService) GetRate(fromCurrency, toCurrency models.Currency, date time.Time) (float64, error) {
	args := m.Called(fromCurrency, toCurrency, date)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExchangeRateService) GetRateRecord(fromCurrency, toCurrency models.Currency, date time.Time) (*models.ExchangeRate, error) {
	args := m.Called(fromCurrency, toCurrency, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExchangeRate), args.Error(1)
}

func (m *MockExchangeRateService) GetTodayRate(fromCurrency, toCurrency models.Currency) (float64, error) {
	args := m.Called(fromCurrency, toCurrency)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExchangeRateService) RefreshTodayRate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockExchangeRateService) ConvertToTWD(amount float64, currency models.Currency, date time.Time) (float64, error) {
	args := m.Called(amount, currency, date)
	return args.Get(0).(float64), args.Error(1)
}

// MockDiscordService 模擬 DiscordService
type MockDiscordService struct {
	mock.Mock
}

func (m *MockDiscordService) SendMessage(webhookURL string, message *models.DiscordMessage) error {
	args := m.Called(webhookURL, message)
	return args.Error(0)
}

func (m *MockDiscordService) FormatDailyReport(data *models.DailyReportData) *models.DiscordMessage {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.DiscordMessage)
}

func (m *MockDiscordService) SendDailyBillingNotification(webhookURL string, result *service.DailyBillingResult) error {
	args := m.Called(webhookURL, result)
	return args.Error(0)
}

func (m *MockDiscordService) SendSubscriptionExpiryNotification(webhookURL string, subscriptions []*models.Subscription, days int) error {
	args := m.Called(webhookURL, subscriptions, days)
	return args.Error(0)
}

func (m *MockDiscordService) SendInstallmentCompletionNotification(webhookURL string, installments []*models.Installment, remainingCount int) error {
	args := m.Called(webhookURL, installments, remainingCount)
	return args.Error(0)
}

// MockSettingsService 模擬 SettingsService
type MockSettingsService struct {
	mock.Mock
}

func (m *MockSettingsService) GetSettings() (*models.SettingsGroup, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SettingsGroup), args.Error(1)
}

func (m *MockSettingsService) UpdateSettings(input *models.UpdateSettingsGroupInput) (*models.SettingsGroup, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SettingsGroup), args.Error(1)
}

// MockHoldingService 模擬 HoldingService
type MockHoldingService struct {
	mock.Mock
}

func (m *MockHoldingService) GetAllHoldings(filters models.HoldingFilters) ([]*models.Holding, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Holding), args.Error(1)
}

func (m *MockHoldingService) GetHoldingBySymbol(symbol string) (*models.Holding, error) {
	args := m.Called(symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Holding), args.Error(1)
}

// MockRebalanceService 模擬 RebalanceService
type MockRebalanceService struct {
	mock.Mock
}

func (m *MockRebalanceService) CheckRebalance() (*models.RebalanceCheck, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RebalanceCheck), args.Error(1)
}

// MockBillingService 模擬 BillingService
type MockBillingService struct {
	mock.Mock
}

func (m *MockBillingService) ProcessSubscriptionBilling(date time.Time) (*service.BillingResult, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BillingResult), args.Error(1)
}

func (m *MockBillingService) ProcessInstallmentBilling(date time.Time) (*service.BillingResult, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BillingResult), args.Error(1)
}

func (m *MockBillingService) ProcessDailyBilling(date time.Time) (*service.DailyBillingResult, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.DailyBillingResult), args.Error(1)
}

// TestSchedulerManager_RunSnapshotNow_WithExchangeRateUpdate 測試手動觸發快照時會更新匯率
func TestSchedulerManager_RunSnapshotNow_WithExchangeRateUpdate(t *testing.T) {
	// Arrange
	mockSnapshotService := new(MockAssetSnapshotService)
	mockDiscordService := new(MockDiscordService)
	mockSettingsService := new(MockSettingsService)
	mockHoldingService := new(MockHoldingService)
	mockRebalanceService := new(MockRebalanceService)
	mockBillingService := new(MockBillingService)
	mockExchangeRateService := new(MockExchangeRateService)

	config := SchedulerManagerConfig{
		Enabled:           true,
		DailySnapshotTime: "23:59",
	}

	manager := NewSchedulerManager(
		mockSnapshotService,
		mockDiscordService,
		mockSettingsService,
		mockHoldingService,
		mockRebalanceService,
		mockBillingService,
		mockExchangeRateService,
		config,
	)

	// Mock 匯率更新成功
	mockExchangeRateService.On("RefreshTodayRate").Return(nil)

	// Mock 快照建立成功
	mockSnapshotService.On("CreateDailySnapshots").Return(nil)

	// Act
	err := manager.RunSnapshotNow()

	// Assert
	assert.NoError(t, err)
	mockExchangeRateService.AssertExpectations(t)
	mockSnapshotService.AssertExpectations(t)
	// 確保匯率更新在快照建立之前被呼叫
	mockExchangeRateService.AssertCalled(t, "RefreshTodayRate")
	mockSnapshotService.AssertCalled(t, "CreateDailySnapshots")
}

// TestSchedulerManager_RunSnapshotNow_ExchangeRateError 測試匯率更新失敗時仍會建立快照
func TestSchedulerManager_RunSnapshotNow_ExchangeRateError(t *testing.T) {
	// Arrange
	mockSnapshotService := new(MockAssetSnapshotService)
	mockDiscordService := new(MockDiscordService)
	mockSettingsService := new(MockSettingsService)
	mockHoldingService := new(MockHoldingService)
	mockRebalanceService := new(MockRebalanceService)
	mockBillingService := new(MockBillingService)
	mockExchangeRateService := new(MockExchangeRateService)

	config := SchedulerManagerConfig{
		Enabled:           true,
		DailySnapshotTime: "23:59",
	}

	manager := NewSchedulerManager(
		mockSnapshotService,
		mockDiscordService,
		mockSettingsService,
		mockHoldingService,
		mockRebalanceService,
		mockBillingService,
		mockExchangeRateService,
		config,
	)

	// Mock 匯率更新失敗
	mockExchangeRateService.On("RefreshTodayRate").Return(errors.New("API error"))

	// Mock 快照建立成功（即使匯率更新失敗，快照仍應建立）
	mockSnapshotService.On("CreateDailySnapshots").Return(nil)

	// Act
	err := manager.RunSnapshotNow()

	// Assert
	assert.NoError(t, err) // 整體應該成功
	mockExchangeRateService.AssertExpectations(t)
	mockSnapshotService.AssertExpectations(t)
	// 確保即使匯率更新失敗，快照仍會建立
	mockSnapshotService.AssertCalled(t, "CreateDailySnapshots")
}

