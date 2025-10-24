package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
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

