package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/jmoiron/sqlx"
)

// SchedulerLogRepository 排程記錄資料庫操作介面
type SchedulerLogRepository interface {
	Create(log *models.SchedulerLog) error
	Update(log *models.SchedulerLog) error
	GetByID(id int) (*models.SchedulerLog, error)
	List(filters models.SchedulerLogFilters) ([]models.SchedulerLog, error)
	GetLatestByTaskName(taskName string) (*models.SchedulerLog, error)
	GetSummaryByTaskName(taskName string) (*models.SchedulerLogSummary, error)
	DeleteOldLogs(olderThan time.Time) error
}

type schedulerLogRepository struct {
	db *sqlx.DB
}

// NewSchedulerLogRepository 建立排程記錄 repository
func NewSchedulerLogRepository(db *sqlx.DB) SchedulerLogRepository {
	return &schedulerLogRepository{db: db}
}

// Create 建立新的排程執行記錄
func (r *schedulerLogRepository) Create(log *models.SchedulerLog) error {
	query := `
		INSERT INTO scheduler_logs (task_name, status, error_message, started_at, completed_at, duration_seconds)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		log.TaskName,
		log.Status,
		log.ErrorMessage,
		log.StartedAt,
		log.CompletedAt,
		log.DurationSeconds,
	).Scan(&log.ID, &log.CreatedAt)
}

// Update 更新排程執行記錄
func (r *schedulerLogRepository) Update(log *models.SchedulerLog) error {
	query := `
		UPDATE scheduler_logs
		SET status = $1, error_message = $2, completed_at = $3, duration_seconds = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(
		query,
		log.Status,
		log.ErrorMessage,
		log.CompletedAt,
		log.DurationSeconds,
		log.ID,
	)
	return err
}

// GetByID 根據 ID 取得排程記錄
func (r *schedulerLogRepository) GetByID(id int) (*models.SchedulerLog, error) {
	var log models.SchedulerLog
	query := `SELECT * FROM scheduler_logs WHERE id = $1`
	err := r.db.Get(&log, query, id)
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// List 取得排程記錄列表
func (r *schedulerLogRepository) List(filters models.SchedulerLogFilters) ([]models.SchedulerLog, error) {
	query := `SELECT * FROM scheduler_logs WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	// 任務名稱過濾
	if filters.TaskName != "" {
		query += fmt.Sprintf(" AND task_name = $%d", argCount)
		args = append(args, filters.TaskName)
		argCount++
	}

	// 狀態過濾
	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
		argCount++
	}

	// 排序：最新的在前
	query += " ORDER BY started_at DESC"

	// 限制筆數
	limit := filters.Limit
	if limit <= 0 {
		limit = 50 // 預設 50 筆
	}
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, limit)

	var logs []models.SchedulerLog
	err := r.db.Select(&logs, query, args...)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetLatestByTaskName 取得指定任務的最新執行記錄
func (r *schedulerLogRepository) GetLatestByTaskName(taskName string) (*models.SchedulerLog, error) {
	var log models.SchedulerLog
	query := `
		SELECT * FROM scheduler_logs
		WHERE task_name = $1
		ORDER BY started_at DESC
		LIMIT 1
	`
	err := r.db.Get(&log, query, taskName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 沒有記錄不算錯誤
		}
		return nil, err
	}
	return &log, nil
}

// GetSummaryByTaskName 取得指定任務的執行摘要
func (r *schedulerLogRepository) GetSummaryByTaskName(taskName string) (*models.SchedulerLogSummary, error) {
	summary := &models.SchedulerLogSummary{
		TaskName: taskName,
	}

	// 取得最後一次執行記錄
	latestLog, err := r.GetLatestByTaskName(taskName)
	if err != nil {
		return nil, err
	}

	if latestLog != nil {
		summary.LastRunStatus = latestLog.Status
		summary.LastRunTime = &latestLog.StartedAt
		summary.LastErrorMsg = latestLog.ErrorMessage
	}

	// 取得最後一次成功執行的時間
	var lastSuccessTime time.Time
	query := `
		SELECT started_at FROM scheduler_logs
		WHERE task_name = $1 AND status = 'success'
		ORDER BY started_at DESC
		LIMIT 1
	`
	err = r.db.Get(&lastSuccessTime, query, taskName)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == nil {
		summary.LastSuccessTime = &lastSuccessTime
	}

	return summary, nil
}

// DeleteOldLogs 刪除舊的排程記錄（資料清理）
func (r *schedulerLogRepository) DeleteOldLogs(olderThan time.Time) error {
	query := `DELETE FROM scheduler_logs WHERE started_at < $1`
	_, err := r.db.Exec(query, olderThan)
	return err
}

