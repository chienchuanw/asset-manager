package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PerformanceSnapshotRepository 績效快照 Repository 介面
type PerformanceSnapshotRepository interface {
	// Create 建立每日績效快照（包含明細）
	Create(input *models.CreateDailyPerformanceSnapshotInput) (*models.DailyPerformanceSnapshot, error)
	// GetByDate 取得指定日期的績效快照
	GetByDate(date time.Time) (*models.DailyPerformanceSnapshot, error)
	// GetByDateRange 取得日期範圍內的績效快照
	GetByDateRange(startDate, endDate time.Time) ([]*models.DailyPerformanceSnapshot, error)
	// GetLatest 取得最新的 N 筆績效快照
	GetLatest(limit int) ([]*models.DailyPerformanceSnapshot, error)
	// GetDetailsBySnapshotID 取得指定快照的明細
	GetDetailsBySnapshotID(snapshotID uuid.UUID) ([]*models.DailyPerformanceSnapshotDetail, error)
	// GetDetailsByDateRange 取得日期範圍內的所有快照明細
	GetDetailsByDateRange(startDate, endDate time.Time) (map[uuid.UUID][]*models.DailyPerformanceSnapshotDetail, error)
}

type performanceSnapshotRepository struct {
	db *sqlx.DB
}

// NewPerformanceSnapshotRepository 建立績效快照 Repository
func NewPerformanceSnapshotRepository(db *sqlx.DB) PerformanceSnapshotRepository {
	return &performanceSnapshotRepository{db: db}
}

// Create 建立每日績效快照（包含明細）
func (r *performanceSnapshotRepository) Create(input *models.CreateDailyPerformanceSnapshotInput) (*models.DailyPerformanceSnapshot, error) {
	// 開始交易
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 建立主快照
	snapshot := &models.DailyPerformanceSnapshot{
		ID:                 uuid.New(),
		SnapshotDate:       input.SnapshotDate,
		TotalMarketValue:   input.TotalMarketValue,
		TotalCost:          input.TotalCost,
		TotalUnrealizedPL:  input.TotalUnrealizedPL,
		TotalUnrealizedPct: input.TotalUnrealizedPct,
		TotalRealizedPL:    input.TotalRealizedPL,
		TotalRealizedPct:   input.TotalRealizedPct,
		HoldingCount:       input.HoldingCount,
		Currency:           input.Currency,
	}

	query := `
		INSERT INTO daily_performance_snapshots (
			id, snapshot_date, total_market_value, total_cost,
			total_unrealized_pl, total_unrealized_pct,
			total_realized_pl, total_realized_pct,
			holding_count, currency
		) VALUES (
			:id, :snapshot_date, :total_market_value, :total_cost,
			:total_unrealized_pl, :total_unrealized_pct,
			:total_realized_pl, :total_realized_pct,
			:holding_count, :currency
		)
		ON CONFLICT (snapshot_date) DO UPDATE SET
			total_market_value = EXCLUDED.total_market_value,
			total_cost = EXCLUDED.total_cost,
			total_unrealized_pl = EXCLUDED.total_unrealized_pl,
			total_unrealized_pct = EXCLUDED.total_unrealized_pct,
			total_realized_pl = EXCLUDED.total_realized_pl,
			total_realized_pct = EXCLUDED.total_realized_pct,
			holding_count = EXCLUDED.holding_count,
			currency = EXCLUDED.currency
		RETURNING id, snapshot_date, total_market_value, total_cost,
			total_unrealized_pl, total_unrealized_pct,
			total_realized_pl, total_realized_pct,
			holding_count, currency, created_at, updated_at
	`

	rows, err := tx.NamedQuery(query, snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(snapshot); err != nil {
			return nil, fmt.Errorf("failed to scan snapshot: %w", err)
		}
	}
	rows.Close()

	// 建立明細（使用 CASCADE DELETE，當主快照被更新時會自動刪除舊明細）
	if len(input.Details) > 0 {
		// 先刪除舊的明細
		_, err = tx.Exec("DELETE FROM daily_performance_snapshot_details WHERE snapshot_id = $1", snapshot.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete old details: %w", err)
		}

		// 插入新的明細
		detailQuery := `
			INSERT INTO daily_performance_snapshot_details (
				id, snapshot_id, asset_type, market_value, cost,
				unrealized_pl, unrealized_pct, realized_pl, realized_pct,
				holding_count
			) VALUES (
				:id, :snapshot_id, :asset_type, :market_value, :cost,
				:unrealized_pl, :unrealized_pct, :realized_pl, :realized_pct,
				:holding_count
			)
		`

		for _, detailInput := range input.Details {
			detail := &models.DailyPerformanceSnapshotDetail{
				ID:            uuid.New(),
				SnapshotID:    snapshot.ID,
				AssetType:     detailInput.AssetType,
				MarketValue:   detailInput.MarketValue,
				Cost:          detailInput.Cost,
				UnrealizedPL:  detailInput.UnrealizedPL,
				UnrealizedPct: detailInput.UnrealizedPct,
				RealizedPL:    detailInput.RealizedPL,
				RealizedPct:   detailInput.RealizedPct,
				HoldingCount:  detailInput.HoldingCount,
			}

			if _, err := tx.NamedExec(detailQuery, detail); err != nil {
				return nil, fmt.Errorf("failed to create detail: %w", err)
			}
		}
	}

	// 提交交易
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return snapshot, nil
}

// GetByDate 取得指定日期的績效快照
func (r *performanceSnapshotRepository) GetByDate(date time.Time) (*models.DailyPerformanceSnapshot, error) {
	query := `
		SELECT id, snapshot_date, total_market_value, total_cost,
			total_unrealized_pl, total_unrealized_pct,
			total_realized_pl, total_realized_pct,
			holding_count, currency, created_at, updated_at
		FROM daily_performance_snapshots
		WHERE snapshot_date = $1
	`

	var snapshot models.DailyPerformanceSnapshot
	if err := r.db.Get(&snapshot, query, date); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get snapshot by date: %w", err)
	}

	return &snapshot, nil
}

// GetByDateRange 取得日期範圍內的績效快照
func (r *performanceSnapshotRepository) GetByDateRange(startDate, endDate time.Time) ([]*models.DailyPerformanceSnapshot, error) {
	query := `
		SELECT id, snapshot_date, total_market_value, total_cost,
			total_unrealized_pl, total_unrealized_pct,
			total_realized_pl, total_realized_pct,
			holding_count, currency, created_at, updated_at
		FROM daily_performance_snapshots
		WHERE snapshot_date >= $1 AND snapshot_date <= $2
		ORDER BY snapshot_date ASC
	`

	var snapshots []*models.DailyPerformanceSnapshot
	if err := r.db.Select(&snapshots, query, startDate, endDate); err != nil {
		return nil, fmt.Errorf("failed to get snapshots by date range: %w", err)
	}

	return snapshots, nil
}

// GetLatest 取得最新的 N 筆績效快照
func (r *performanceSnapshotRepository) GetLatest(limit int) ([]*models.DailyPerformanceSnapshot, error) {
	query := `
		SELECT id, snapshot_date, total_market_value, total_cost,
			total_unrealized_pl, total_unrealized_pct,
			total_realized_pl, total_realized_pct,
			holding_count, currency, created_at, updated_at
		FROM daily_performance_snapshots
		ORDER BY snapshot_date DESC
		LIMIT $1
	`

	var snapshots []*models.DailyPerformanceSnapshot
	if err := r.db.Select(&snapshots, query, limit); err != nil {
		return nil, fmt.Errorf("failed to get latest snapshots: %w", err)
	}

	return snapshots, nil
}

// GetDetailsBySnapshotID 取得指定快照的明細
func (r *performanceSnapshotRepository) GetDetailsBySnapshotID(snapshotID uuid.UUID) ([]*models.DailyPerformanceSnapshotDetail, error) {
	query := `
		SELECT id, snapshot_id, asset_type, market_value, cost,
			unrealized_pl, unrealized_pct, realized_pl, realized_pct,
			holding_count, created_at
		FROM daily_performance_snapshot_details
		WHERE snapshot_id = $1
		ORDER BY asset_type
	`

	var details []*models.DailyPerformanceSnapshotDetail
	if err := r.db.Select(&details, query, snapshotID); err != nil {
		return nil, fmt.Errorf("failed to get details by snapshot ID: %w", err)
	}

	return details, nil
}

// GetDetailsByDateRange 取得日期範圍內的所有快照明細
func (r *performanceSnapshotRepository) GetDetailsByDateRange(startDate, endDate time.Time) (map[uuid.UUID][]*models.DailyPerformanceSnapshotDetail, error) {
	query := `
		SELECT d.id, d.snapshot_id, d.asset_type, d.market_value, d.cost,
			d.unrealized_pl, d.unrealized_pct, d.realized_pl, d.realized_pct,
			d.holding_count, d.created_at
		FROM daily_performance_snapshot_details d
		INNER JOIN daily_performance_snapshots s ON d.snapshot_id = s.id
		WHERE s.snapshot_date >= $1 AND s.snapshot_date <= $2
		ORDER BY s.snapshot_date, d.asset_type
	`

	var details []*models.DailyPerformanceSnapshotDetail
	if err := r.db.Select(&details, query, startDate, endDate); err != nil {
		return nil, fmt.Errorf("failed to get details by date range: %w", err)
	}

	// 按 snapshot_id 分組
	result := make(map[uuid.UUID][]*models.DailyPerformanceSnapshotDetail)
	for _, detail := range details {
		result[detail.SnapshotID] = append(result[detail.SnapshotID], detail)
	}

	return result, nil
}

