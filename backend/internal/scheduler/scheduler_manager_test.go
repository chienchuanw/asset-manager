package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSchedulerLogRepository 模擬 SchedulerLogRepository
type MockSchedulerLogRepository struct {
	mock.Mock
}

func (m *MockSchedulerLogRepository) Create(log *models.SchedulerLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockSchedulerLogRepository) Update(log *models.SchedulerLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockSchedulerLogRepository) GetByID(id int) (*models.SchedulerLog, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SchedulerLog), args.Error(1)
}

func (m *MockSchedulerLogRepository) List(filters models.SchedulerLogFilters) ([]models.SchedulerLog, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SchedulerLog), args.Error(1)
}

func (m *MockSchedulerLogRepository) GetLatestByTaskName(taskName string) (*models.SchedulerLog, error) {
	args := m.Called(taskName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SchedulerLog), args.Error(1)
}

func (m *MockSchedulerLogRepository) GetSummaryByTaskName(taskName string) (*models.SchedulerLogSummary, error) {
	args := m.Called(taskName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SchedulerLogSummary), args.Error(1)
}

func (m *MockSchedulerLogRepository) DeleteOldLogs(olderThan time.Time) error {
	args := m.Called(olderThan)
	return args.Error(0)
}

// 確保 MockSchedulerLogRepository 實作 SchedulerLogRepository 介面
var _ repository.SchedulerLogRepository = (*MockSchedulerLogRepository)(nil)

// TestSchedulerManager_logTaskExecution_Success 測試記錄成功的任務執行
func TestSchedulerManager_logTaskExecution_Success(t *testing.T) {
	// Arrange
	mockSnapshotService := new(MockAssetSnapshotService)
	mockDiscordService := new(MockDiscordService)
	mockSettingsService := new(MockSettingsService)
	mockHoldingService := new(MockHoldingService)
	mockRebalanceService := new(MockRebalanceService)
	mockBillingService := new(MockBillingService)
	mockExchangeRateService := new(MockExchangeRateService)
	mockCreditCardService := new(MockCreditCardService)
	mockSchedulerLogRepo := new(MockSchedulerLogRepository)

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
		mockCreditCardService,
		mockSchedulerLogRepo,
		config,
	)

	// Mock repository Create 方法
	mockSchedulerLogRepo.On("Create", mock.MatchedBy(func(log *models.SchedulerLog) bool {
		return log.TaskName == "test_task" &&
			log.Status == "success" &&
			log.ErrorMessage == nil
	})).Return(nil)

	// Act
	startTime := time.Now()
	manager.logTaskExecution("test_task", startTime, nil)

	// Assert
	mockSchedulerLogRepo.AssertExpectations(t)
}

// TestSchedulerManager_logTaskExecution_Failed 測試記錄失敗的任務執行
func TestSchedulerManager_logTaskExecution_Failed(t *testing.T) {
	// Arrange
	mockSnapshotService := new(MockAssetSnapshotService)
	mockDiscordService := new(MockDiscordService)
	mockSettingsService := new(MockSettingsService)
	mockHoldingService := new(MockHoldingService)
	mockRebalanceService := new(MockRebalanceService)
	mockBillingService := new(MockBillingService)
	mockExchangeRateService := new(MockExchangeRateService)
	mockCreditCardService := new(MockCreditCardService)
	mockSchedulerLogRepo := new(MockSchedulerLogRepository)

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
		mockCreditCardService,
		mockSchedulerLogRepo,
		config,
	)

	// Mock repository Create 方法
	mockSchedulerLogRepo.On("Create", mock.MatchedBy(func(log *models.SchedulerLog) bool {
		return log.TaskName == "test_task" &&
			log.Status == "failed" &&
			log.ErrorMessage != nil &&
			*log.ErrorMessage == "test error"
	})).Return(nil)

	// Act
	startTime := time.Now()
	taskErr := errors.New("test error")
	manager.logTaskExecution("test_task", startTime, taskErr)

	// Assert
	mockSchedulerLogRepo.AssertExpectations(t)
}

// TestSchedulerManager_logTaskExecution_NoRepository 測試沒有 repository 時不會出錯
func TestSchedulerManager_logTaskExecution_NoRepository(t *testing.T) {
	// Arrange
	mockSnapshotService := new(MockAssetSnapshotService)
	mockDiscordService := new(MockDiscordService)
	mockSettingsService := new(MockSettingsService)
	mockHoldingService := new(MockHoldingService)
	mockRebalanceService := new(MockRebalanceService)
	mockBillingService := new(MockBillingService)
	mockExchangeRateService := new(MockExchangeRateService)
	mockCreditCardService := new(MockCreditCardService)

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
		mockCreditCardService,
		nil, // 沒有 repository
		config,
	)

	// Act - 應該不會 panic
	startTime := time.Now()
	manager.logTaskExecution("test_task", startTime, nil)

	// Assert - 沒有 panic 就算成功
	assert.NotNil(t, manager)
}

// TestSchedulerManager_sendFailureNotification_Success 測試發送失敗通知
func TestSchedulerManager_sendFailureNotification_Success(t *testing.T) {
	// Arrange
	mockSnapshotService := new(MockAssetSnapshotService)
	mockDiscordService := new(MockDiscordService)
	mockSettingsService := new(MockSettingsService)
	mockHoldingService := new(MockHoldingService)
	mockRebalanceService := new(MockRebalanceService)
	mockBillingService := new(MockBillingService)
	mockExchangeRateService := new(MockExchangeRateService)
	mockCreditCardService := new(MockCreditCardService)
	mockSchedulerLogRepo := new(MockSchedulerLogRepository)

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
		mockCreditCardService,
		mockSchedulerLogRepo,
		config,
	)

	// Mock settings service
	settings := &models.SettingsGroup{
		Discord: models.DiscordSettings{
			Enabled:    true,
			WebhookURL: "https://discord.com/webhook/test",
		},
	}
	mockSettingsService.On("GetSettings").Return(settings, nil)

	// Mock Discord service
	mockDiscordService.On("SendMessage", "https://discord.com/webhook/test", mock.AnythingOfType("*models.DiscordMessage")).Return(nil)

	// Act
	taskErr := errors.New("test error")
	manager.sendFailureNotification("測試任務", taskErr)

	// Assert
	mockSettingsService.AssertExpectations(t)
	mockDiscordService.AssertExpectations(t)
}

// TestSchedulerManager_sendFailureNotification_DiscordDisabled 測試 Discord 停用時不發送通知
func TestSchedulerManager_sendFailureNotification_DiscordDisabled(t *testing.T) {
	// Arrange
	mockSnapshotService := new(MockAssetSnapshotService)
	mockDiscordService := new(MockDiscordService)
	mockSettingsService := new(MockSettingsService)
	mockHoldingService := new(MockHoldingService)
	mockRebalanceService := new(MockRebalanceService)
	mockBillingService := new(MockBillingService)
	mockExchangeRateService := new(MockExchangeRateService)
	mockCreditCardService := new(MockCreditCardService)
	mockSchedulerLogRepo := new(MockSchedulerLogRepository)

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
		mockCreditCardService,
		mockSchedulerLogRepo,
		config,
	)

	// Mock settings service - Discord 停用
	settings := &models.SettingsGroup{
		Discord: models.DiscordSettings{
			Enabled:    false,
			WebhookURL: "",
		},
	}
	mockSettingsService.On("GetSettings").Return(settings, nil)

	// Act
	taskErr := errors.New("test error")
	manager.sendFailureNotification("測試任務", taskErr)

	// Assert
	mockSettingsService.AssertExpectations(t)
	// Discord service 不應該被呼叫
	mockDiscordService.AssertNotCalled(t, "SendMessage")
}

// TestSchedulerManager_GetTaskSummaries_Success 測試取得任務摘要
func TestSchedulerManager_GetTaskSummaries_Success(t *testing.T) {
	// Arrange
	mockSnapshotService := new(MockAssetSnapshotService)
	mockDiscordService := new(MockDiscordService)
	mockSettingsService := new(MockSettingsService)
	mockHoldingService := new(MockHoldingService)
	mockRebalanceService := new(MockRebalanceService)
	mockBillingService := new(MockBillingService)
	mockExchangeRateService := new(MockExchangeRateService)
	mockCreditCardService := new(MockCreditCardService)
	mockSchedulerLogRepo := new(MockSchedulerLogRepository)

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
		mockCreditCardService,
		mockSchedulerLogRepo,
		config,
	)

	// Mock repository GetSummaryByTaskName 方法
	now := time.Now()
	mockSchedulerLogRepo.On("GetSummaryByTaskName", "daily_snapshot").Return(&models.SchedulerLogSummary{
		TaskName:      "daily_snapshot",
		LastRunStatus: "success",
		LastRunTime:   &now,
	}, nil)

	mockSchedulerLogRepo.On("GetSummaryByTaskName", "discord_report").Return(&models.SchedulerLogSummary{
		TaskName:      "discord_report",
		LastRunStatus: "success",
		LastRunTime:   &now,
	}, nil)

	mockSchedulerLogRepo.On("GetSummaryByTaskName", "daily_billing").Return(&models.SchedulerLogSummary{
		TaskName:      "daily_billing",
		LastRunStatus: "success",
		LastRunTime:   &now,
	}, nil)

	// Act
	summaries, err := manager.GetTaskSummaries()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, summaries, 3)
	mockSchedulerLogRepo.AssertExpectations(t)
}

// TestSchedulerManager_GetTaskSummaries_NoRepository 測試沒有 repository 時返回錯誤
func TestSchedulerManager_GetTaskSummaries_NoRepository(t *testing.T) {
	// Arrange
	mockSnapshotService := new(MockAssetSnapshotService)
	mockDiscordService := new(MockDiscordService)
	mockSettingsService := new(MockSettingsService)
	mockHoldingService := new(MockHoldingService)
	mockRebalanceService := new(MockRebalanceService)
	mockBillingService := new(MockBillingService)
	mockExchangeRateService := new(MockExchangeRateService)
	mockCreditCardService := new(MockCreditCardService)

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
		mockCreditCardService,
		nil, // 沒有 repository
		config,
	)

	// Act
	summaries, err := manager.GetTaskSummaries()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, summaries)
	assert.Contains(t, err.Error(), "scheduler log repository not available")
}

