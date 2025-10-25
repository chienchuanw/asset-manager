package service

import (
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSettingsRepository 用於測試的 Mock Repository
type MockSettingsRepository struct {
	mock.Mock
}

func (m *MockSettingsRepository) GetByKey(key string) (*models.Setting, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Setting), args.Error(1)
}

func (m *MockSettingsRepository) GetAll() ([]*models.Setting, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Setting), args.Error(1)
}

func (m *MockSettingsRepository) Update(key string, input *models.UpdateSettingInput) (*models.Setting, error) {
	args := m.Called(key, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Setting), args.Error(1)
}

// TestSettingsService_GetSettings 測試取得設定群組
func TestSettingsService_GetSettings(t *testing.T) {
	mockRepo := new(MockSettingsRepository)
	service := NewSettingsService(mockRepo)

	// Mock 資料
	settings := []*models.Setting{
		{Key: "discord_webhook_url", Value: "https://discord.com/api/webhooks/test"},
		{Key: "discord_enabled", Value: "true"},
		{Key: "discord_report_time", Value: "09:00"},
		{Key: "target_allocation_tw_stock", Value: "40"},
		{Key: "target_allocation_us_stock", Value: "40"},
		{Key: "target_allocation_crypto", Value: "20"},
		{Key: "rebalance_threshold", Value: "5"},
	}

	mockRepo.On("GetAll").Return(settings, nil)

	// 執行
	result, err := service.GetSettings()

	// 驗證
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "https://discord.com/api/webhooks/test", result.Discord.WebhookURL)
	assert.Equal(t, true, result.Discord.Enabled)
	assert.Equal(t, "09:00", result.Discord.ReportTime)
	assert.Equal(t, 40.0, result.Allocation.TWStock)
	assert.Equal(t, 40.0, result.Allocation.USStock)
	assert.Equal(t, 20.0, result.Allocation.Crypto)
	assert.Equal(t, 5.0, result.Allocation.RebalanceThreshold)

	mockRepo.AssertExpectations(t)
}

// TestSettingsService_UpdateSettings 測試更新設定群組
func TestSettingsService_UpdateSettings(t *testing.T) {
	mockRepo := new(MockSettingsRepository)
	service := NewSettingsService(mockRepo)

	// 輸入
	input := &models.UpdateSettingsGroupInput{
		Discord: &models.DiscordSettings{
			WebhookURL: "https://discord.com/api/webhooks/new",
			Enabled:    true,
			ReportTime: "10:00",
		},
		Allocation: &models.AllocationSettings{
			TWStock:            50,
			USStock:            30,
			Crypto:             20,
			RebalanceThreshold: 10,
		},
	}

	// Mock 更新
	mockRepo.On("Update", "discord_webhook_url", mock.Anything).Return(&models.Setting{Key: "discord_webhook_url", Value: "https://discord.com/api/webhooks/new"}, nil)
	mockRepo.On("Update", "discord_enabled", mock.Anything).Return(&models.Setting{Key: "discord_enabled", Value: "true"}, nil)
	mockRepo.On("Update", "discord_report_time", mock.Anything).Return(&models.Setting{Key: "discord_report_time", Value: "10:00"}, nil)
	mockRepo.On("Update", "target_allocation_tw_stock", mock.Anything).Return(&models.Setting{Key: "target_allocation_tw_stock", Value: "50"}, nil)
	mockRepo.On("Update", "target_allocation_us_stock", mock.Anything).Return(&models.Setting{Key: "target_allocation_us_stock", Value: "30"}, nil)
	mockRepo.On("Update", "target_allocation_crypto", mock.Anything).Return(&models.Setting{Key: "target_allocation_crypto", Value: "20"}, nil)
	mockRepo.On("Update", "rebalance_threshold", mock.Anything).Return(&models.Setting{Key: "rebalance_threshold", Value: "10"}, nil)

	// Mock GetAll（更新後取得）
	updatedSettings := []*models.Setting{
		{Key: "discord_webhook_url", Value: "https://discord.com/api/webhooks/new"},
		{Key: "discord_enabled", Value: "true"},
		{Key: "discord_report_time", Value: "10:00"},
		{Key: "target_allocation_tw_stock", Value: "50"},
		{Key: "target_allocation_us_stock", Value: "30"},
		{Key: "target_allocation_crypto", Value: "20"},
		{Key: "rebalance_threshold", Value: "10"},
	}
	mockRepo.On("GetAll").Return(updatedSettings, nil)

	// 執行
	result, err := service.UpdateSettings(input)

	// 驗證
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "https://discord.com/api/webhooks/new", result.Discord.WebhookURL)
	assert.Equal(t, true, result.Discord.Enabled)
	assert.Equal(t, "10:00", result.Discord.ReportTime)
	assert.Equal(t, 50.0, result.Allocation.TWStock)
	assert.Equal(t, 30.0, result.Allocation.USStock)
	assert.Equal(t, 20.0, result.Allocation.Crypto)
	assert.Equal(t, 10.0, result.Allocation.RebalanceThreshold)

	mockRepo.AssertExpectations(t)
}

