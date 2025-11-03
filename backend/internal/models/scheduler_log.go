package models

import "time"

// SchedulerLog 排程任務執行記錄
type SchedulerLog struct {
	ID              int        `db:"id" json:"id"`
	TaskName        string     `db:"task_name" json:"task_name"`
	Status          string     `db:"status" json:"status"` // success, failed, running
	ErrorMessage    *string    `db:"error_message" json:"error_message,omitempty"`
	StartedAt       time.Time  `db:"started_at" json:"started_at"`
	CompletedAt     *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	DurationSeconds *float64   `db:"duration_seconds" json:"duration_seconds,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
}

// SchedulerLogFilters 排程記錄查詢過濾條件
type SchedulerLogFilters struct {
	TaskName string // 任務名稱過濾
	Status   string // 狀態過濾
	Limit    int    // 限制筆數（預設 50）
}

// SchedulerLogSummary 排程任務執行摘要（用於狀態查詢）
type SchedulerLogSummary struct {
	TaskName        string     `json:"task_name"`
	LastRunStatus   string     `json:"last_run_status"`
	LastRunTime     *time.Time `json:"last_run_time,omitempty"`
	LastSuccessTime *time.Time `json:"last_success_time,omitempty"`
	LastErrorMsg    *string    `json:"last_error_message,omitempty"`
}

