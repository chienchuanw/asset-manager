package repository

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupSchedulerLogTestDB 設定測試資料庫並返回 sqlx.DB
func setupSchedulerLogTestDB(t *testing.T) *sqlx.DB {
	db, err := setupTestDB()
	require.NoError(t, err)
	
	dbx := sqlx.NewDb(db, "postgres")
	
	// 清理測試資料
	_, err = db.Exec("DELETE FROM scheduler_logs")
	require.NoError(t, err)
	
	return dbx
}

// TestSchedulerLogRepository_Create 測試建立排程記錄
func TestSchedulerLogRepository_Create(t *testing.T) {
	dbx := setupSchedulerLogTestDB(t)
	defer dbx.Close()

	repo := NewSchedulerLogRepository(dbx)

	// 準備測試資料
	startedAt := time.Now()
	completedAt := startedAt.Add(5 * time.Second)
	duration := 5.234
	errorMsg := "test error"

	log := &models.SchedulerLog{
		TaskName:        "daily_snapshot",
		Status:          "success",
		ErrorMessage:    nil,
		StartedAt:       startedAt,
		CompletedAt:     &completedAt,
		DurationSeconds: &duration,
	}

	// 執行建立
	err := repo.Create(log)
	require.NoError(t, err)
	assert.NotZero(t, log.ID)
	assert.NotZero(t, log.CreatedAt)

	// 測試建立失敗記錄
	logFailed := &models.SchedulerLog{
		TaskName:        "discord_report",
		Status:          "failed",
		ErrorMessage:    &errorMsg,
		StartedAt:       startedAt,
		CompletedAt:     &completedAt,
		DurationSeconds: &duration,
	}

	err = repo.Create(logFailed)
	require.NoError(t, err)
	assert.NotZero(t, logFailed.ID)
	assert.Equal(t, "test error", *logFailed.ErrorMessage)
}

// TestSchedulerLogRepository_Update 測試更新排程記錄
func TestSchedulerLogRepository_Update(t *testing.T) {
	dbx := setupSchedulerLogTestDB(t)
	defer dbx.Close()

	repo := NewSchedulerLogRepository(dbx)

	// 先建立一筆記錄
	startedAt := time.Now()
	log := &models.SchedulerLog{
		TaskName:  "daily_snapshot",
		Status:    "running",
		StartedAt: startedAt,
	}

	err := repo.Create(log)
	require.NoError(t, err)

	// 更新為成功狀態
	completedAt := startedAt.Add(3 * time.Second)
	duration := 3.456
	log.Status = "success"
	log.CompletedAt = &completedAt
	log.DurationSeconds = &duration

	err = repo.Update(log)
	require.NoError(t, err)

	// 驗證更新結果
	updated, err := repo.GetByID(log.ID)
	require.NoError(t, err)
	assert.Equal(t, "success", updated.Status)
	assert.NotNil(t, updated.CompletedAt)
	assert.NotNil(t, updated.DurationSeconds)
	assert.InDelta(t, 3.456, *updated.DurationSeconds, 0.001)
}

// TestSchedulerLogRepository_GetByID 測試根據 ID 取得記錄
func TestSchedulerLogRepository_GetByID(t *testing.T) {
	dbx := setupSchedulerLogTestDB(t)
	defer dbx.Close()

	repo := NewSchedulerLogRepository(dbx)

	// 建立測試資料
	log := &models.SchedulerLog{
		TaskName:  "daily_snapshot",
		Status:    "success",
		StartedAt: time.Now(),
	}

	err := repo.Create(log)
	require.NoError(t, err)

	// 取得記錄
	retrieved, err := repo.GetByID(log.ID)
	require.NoError(t, err)
	assert.Equal(t, log.ID, retrieved.ID)
	assert.Equal(t, "daily_snapshot", retrieved.TaskName)
	assert.Equal(t, "success", retrieved.Status)

	// 測試不存在的 ID
	_, err = repo.GetByID(99999)
	assert.Error(t, err)
}

// TestSchedulerLogRepository_List 測試列出排程記錄
func TestSchedulerLogRepository_List(t *testing.T) {
	dbx := setupSchedulerLogTestDB(t)
	defer dbx.Close()

	repo := NewSchedulerLogRepository(dbx)

	// 建立多筆測試資料
	now := time.Now()
	logs := []*models.SchedulerLog{
		{TaskName: "daily_snapshot", Status: "success", StartedAt: now.Add(-3 * time.Hour)},
		{TaskName: "daily_snapshot", Status: "failed", StartedAt: now.Add(-2 * time.Hour)},
		{TaskName: "discord_report", Status: "success", StartedAt: now.Add(-1 * time.Hour)},
		{TaskName: "daily_billing", Status: "success", StartedAt: now},
	}

	for _, log := range logs {
		err := repo.Create(log)
		require.NoError(t, err)
	}

	// 測試不帶過濾條件
	results, err := repo.List(models.SchedulerLogFilters{})
	require.NoError(t, err)
	assert.Len(t, results, 4)

	// 測試按任務名稱過濾
	results, err = repo.List(models.SchedulerLogFilters{TaskName: "daily_snapshot"})
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// 測試按狀態過濾
	results, err = repo.List(models.SchedulerLogFilters{Status: "success"})
	require.NoError(t, err)
	assert.Len(t, results, 3)

	// 測試限制筆數
	results, err = repo.List(models.SchedulerLogFilters{Limit: 2})
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// 驗證排序（最新的在前）
	assert.True(t, results[0].StartedAt.After(results[1].StartedAt) || results[0].StartedAt.Equal(results[1].StartedAt))
}

// TestSchedulerLogRepository_GetLatestByTaskName 測試取得最新記錄
func TestSchedulerLogRepository_GetLatestByTaskName(t *testing.T) {
	dbx := setupSchedulerLogTestDB(t)
	defer dbx.Close()

	repo := NewSchedulerLogRepository(dbx)

	// 建立多筆同任務的記錄
	now := time.Now()
	logs := []*models.SchedulerLog{
		{TaskName: "daily_snapshot", Status: "success", StartedAt: now.Add(-2 * time.Hour)},
		{TaskName: "daily_snapshot", Status: "failed", StartedAt: now.Add(-1 * time.Hour)},
		{TaskName: "daily_snapshot", Status: "success", StartedAt: now},
	}

	for _, log := range logs {
		err := repo.Create(log)
		require.NoError(t, err)
	}

	// 取得最新記錄
	latest, err := repo.GetLatestByTaskName("daily_snapshot")
	require.NoError(t, err)
	assert.NotNil(t, latest)
	assert.Equal(t, "success", latest.Status)
	assert.True(t, latest.StartedAt.After(now.Add(-1 * time.Second)))

	// 測試不存在的任務
	notFound, err := repo.GetLatestByTaskName("non_existent_task")
	require.NoError(t, err)
	assert.Nil(t, notFound)
}

// TestSchedulerLogRepository_GetSummaryByTaskName 測試取得任務摘要
func TestSchedulerLogRepository_GetSummaryByTaskName(t *testing.T) {
	dbx := setupSchedulerLogTestDB(t)
	defer dbx.Close()

	repo := NewSchedulerLogRepository(dbx)

	// 建立測試資料
	now := time.Now()
	errorMsg := "connection timeout"
	logs := []*models.SchedulerLog{
		{TaskName: "daily_snapshot", Status: "success", StartedAt: now.Add(-3 * time.Hour)},
		{TaskName: "daily_snapshot", Status: "success", StartedAt: now.Add(-2 * time.Hour)},
		{TaskName: "daily_snapshot", Status: "failed", StartedAt: now.Add(-1 * time.Hour), ErrorMessage: &errorMsg},
	}

	for _, log := range logs {
		err := repo.Create(log)
		require.NoError(t, err)
	}

	// 取得摘要
	summary, err := repo.GetSummaryByTaskName("daily_snapshot")
	require.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, "daily_snapshot", summary.TaskName)
	assert.Equal(t, "failed", summary.LastRunStatus)
	assert.NotNil(t, summary.LastRunTime)
	assert.NotNil(t, summary.LastSuccessTime)
	assert.NotNil(t, summary.LastErrorMsg)
	assert.Equal(t, "connection timeout", *summary.LastErrorMsg)

	// 測試不存在的任務
	emptySummary, err := repo.GetSummaryByTaskName("non_existent_task")
	require.NoError(t, err)
	assert.NotNil(t, emptySummary)
	assert.Equal(t, "non_existent_task", emptySummary.TaskName)
	assert.Empty(t, emptySummary.LastRunStatus)
}

// TestSchedulerLogRepository_DeleteOldLogs 測試刪除舊記錄
func TestSchedulerLogRepository_DeleteOldLogs(t *testing.T) {
	dbx := setupSchedulerLogTestDB(t)
	defer dbx.Close()

	repo := NewSchedulerLogRepository(dbx)

	// 建立測試資料
	now := time.Now()
	logs := []*models.SchedulerLog{
		{TaskName: "daily_snapshot", Status: "success", StartedAt: now.Add(-48 * time.Hour)},
		{TaskName: "daily_snapshot", Status: "success", StartedAt: now.Add(-24 * time.Hour)},
		{TaskName: "daily_snapshot", Status: "success", StartedAt: now.Add(-1 * time.Hour)},
	}

	for _, log := range logs {
		err := repo.Create(log)
		require.NoError(t, err)
	}

	// 刪除 30 小時前的記錄
	err := repo.DeleteOldLogs(now.Add(-30 * time.Hour))
	require.NoError(t, err)

	// 驗證結果
	results, err := repo.List(models.SchedulerLogFilters{})
	require.NoError(t, err)
	assert.Len(t, results, 2) // 應該剩下 2 筆
}

