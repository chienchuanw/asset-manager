package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
)

// CashFlowReportLogRepository 現金流報告記錄資料存取介面
type CashFlowReportLogRepository interface {
	Create(input *models.CreateCashFlowReportLogInput) (*models.CashFlowReportLog, error)
	GetByID(id uuid.UUID) (*models.CashFlowReportLog, error)
	GetLatestByType(reportType models.CashFlowReportType, year, month int) (*models.CashFlowReportLog, error)
	UpdateStatus(id uuid.UUID, input *models.UpdateCashFlowReportLogInput) error
	GetPendingRetries(reportType models.CashFlowReportType) ([]*models.CashFlowReportLog, error)
}

// cashFlowReportLogRepository 現金流報告記錄資料存取實作
type cashFlowReportLogRepository struct {
	db *sql.DB
}

// NewCashFlowReportLogRepository 建立新的現金流報告記錄 repository
func NewCashFlowReportLogRepository(db *sql.DB) CashFlowReportLogRepository {
	return &cashFlowReportLogRepository{db: db}
}

// Create 建立新的報告記錄
func (r *cashFlowReportLogRepository) Create(input *models.CreateCashFlowReportLogInput) (*models.CashFlowReportLog, error) {
	query := `
		INSERT INTO cash_flow_report_logs (
			report_type, year, month, sent_at, success, error_message, retry_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, report_type, year, month, sent_at, success, error_message, retry_count, created_at, updated_at
	`

	log := &models.CashFlowReportLog{}
	err := r.db.QueryRow(
		query,
		input.ReportType,
		input.Year,
		input.Month,
		input.SentAt,
		input.Success,
		input.ErrorMessage,
		input.RetryCount,
	).Scan(
		&log.ID,
		&log.ReportType,
		&log.Year,
		&log.Month,
		&log.SentAt,
		&log.Success,
		&log.ErrorMessage,
		&log.RetryCount,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create cash flow report log: %w", err)
	}

	return log, nil
}

// GetByID 根據 ID 取得報告記錄
func (r *cashFlowReportLogRepository) GetByID(id uuid.UUID) (*models.CashFlowReportLog, error) {
	query := `
		SELECT id, report_type, year, month, sent_at, success, error_message, retry_count, created_at, updated_at
		FROM cash_flow_report_logs
		WHERE id = $1
	`

	log := &models.CashFlowReportLog{}
	err := r.db.QueryRow(query, id).Scan(
		&log.ID,
		&log.ReportType,
		&log.Year,
		&log.Month,
		&log.SentAt,
		&log.Success,
		&log.ErrorMessage,
		&log.RetryCount,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cash flow report log not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get cash flow report log: %w", err)
	}

	return log, nil
}

// GetLatestByType 取得指定類型和期間的最新報告記錄
func (r *cashFlowReportLogRepository) GetLatestByType(reportType models.CashFlowReportType, year, month int) (*models.CashFlowReportLog, error) {
	query := `
		SELECT id, report_type, year, month, sent_at, success, error_message, retry_count, created_at, updated_at
		FROM cash_flow_report_logs
		WHERE report_type = $1 AND year = $2 AND month = $3
		ORDER BY created_at DESC
		LIMIT 1
	`

	log := &models.CashFlowReportLog{}
	err := r.db.QueryRow(query, reportType, year, month).Scan(
		&log.ID,
		&log.ReportType,
		&log.Year,
		&log.Month,
		&log.SentAt,
		&log.Success,
		&log.ErrorMessage,
		&log.RetryCount,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // 沒有記錄不算錯誤
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get latest cash flow report log: %w", err)
	}

	return log, nil
}

// UpdateStatus 更新報告記錄狀態
func (r *cashFlowReportLogRepository) UpdateStatus(id uuid.UUID, input *models.UpdateCashFlowReportLogInput) error {
	query := `
		UPDATE cash_flow_report_logs
		SET success = $1, error_message = $2, retry_count = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(
		query,
		input.Success,
		input.ErrorMessage,
		input.RetryCount,
		time.Now(),
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update cash flow report log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cash flow report log not found")
	}

	return nil
}

// GetPendingRetries 取得需要重試的報告記錄
func (r *cashFlowReportLogRepository) GetPendingRetries(reportType models.CashFlowReportType) ([]*models.CashFlowReportLog, error) {
	query := `
		SELECT id, report_type, year, month, sent_at, success, error_message, retry_count, created_at, updated_at
		FROM cash_flow_report_logs
		WHERE report_type = $1 AND success = false AND retry_count < 3
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, reportType)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending retries: %w", err)
	}
	defer rows.Close()

	var logs []*models.CashFlowReportLog
	for rows.Next() {
		log := &models.CashFlowReportLog{}
		err := rows.Scan(
			&log.ID,
			&log.ReportType,
			&log.Year,
			&log.Month,
			&log.SentAt,
			&log.Success,
			&log.ErrorMessage,
			&log.RetryCount,
			&log.CreatedAt,
			&log.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cash flow report log: %w", err)
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cash flow report log rows: %w", err)
	}

	return logs, nil
}

