package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSchedulerManager 模擬 SchedulerManager
type MockSchedulerManager struct {
	mock.Mock
}

func (m *MockSchedulerManager) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSchedulerManager) Stop() {
	m.Called()
}

func (m *MockSchedulerManager) RunSnapshotNow() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSchedulerManager) TriggerSnapshot() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSchedulerManager) TriggerDiscordReport() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSchedulerManager) TriggerDailyBilling() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSchedulerManager) RunDiscordReportNow() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSchedulerManager) ReloadDiscordSchedule() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSchedulerManager) GetStatus() scheduler.SchedulerStatus {
	args := m.Called()
	return args.Get(0).(scheduler.SchedulerStatus)
}

func (m *MockSchedulerManager) GetTaskSummaries() ([]models.SchedulerLogSummary, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SchedulerLogSummary), args.Error(1)
}

// TestSchedulerHandler_GetTaskSummaries_Success 測試成功取得任務摘要
func TestSchedulerHandler_GetTaskSummaries_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockManager := new(MockSchedulerManager)
	handler := NewSchedulerHandler(mockManager)

	now := time.Now()
	summaries := []models.SchedulerLogSummary{
		{
			TaskName:      "daily_snapshot",
			LastRunStatus: "success",
			LastRunTime:   &now,
		},
		{
			TaskName:      "discord_report",
			LastRunStatus: "success",
			LastRunTime:   &now,
		},
	}

	mockManager.On("GetTaskSummaries").Return(summaries, nil)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/scheduler/summaries", nil)

	// Act
	handler.GetTaskSummaries(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 2)

	mockManager.AssertExpectations(t)
}

// TestSchedulerHandler_GetTaskSummaries_Error 測試取得任務摘要失敗
func TestSchedulerHandler_GetTaskSummaries_Error(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockManager := new(MockSchedulerManager)
	handler := NewSchedulerHandler(mockManager)

	mockManager.On("GetTaskSummaries").Return(nil, errors.New("repository not available"))

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/scheduler/summaries", nil)

	// Act
	handler.GetTaskSummaries(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	errorObj, ok := response["error"].(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, errorObj["message"], "Failed to get task summaries")

	mockManager.AssertExpectations(t)
}

// TestSchedulerHandler_GetTaskSummaries_EmptyResult 測試空結果
func TestSchedulerHandler_GetTaskSummaries_EmptyResult(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockManager := new(MockSchedulerManager)
	handler := NewSchedulerHandler(mockManager)

	mockManager.On("GetTaskSummaries").Return([]models.SchedulerLogSummary{}, nil)

	// 建立測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/scheduler/summaries", nil)

	// Act
	handler.GetTaskSummaries(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 0)

	mockManager.AssertExpectations(t)
}

