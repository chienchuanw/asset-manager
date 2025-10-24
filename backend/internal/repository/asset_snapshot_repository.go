package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// AssetSnapshotRepository 資產快照資料存取介面
type AssetSnapshotRepository interface {
	Create(input *models.CreateAssetSnapshotInput) (*models.AssetSnapshot, error)
	GetByDateAndType(date time.Time, assetType models.SnapshotAssetType) (*models.AssetSnapshot, error)
	GetByDateRange(filters models.AssetSnapshotFilters) ([]*models.AssetSnapshot, error)
	Update(date time.Time, assetType models.SnapshotAssetType, valueTWD float64) (*models.AssetSnapshot, error)
	Delete(date time.Time, assetType models.SnapshotAssetType) error
}

// assetSnapshotRepository 資產快照資料存取實作
type assetSnapshotRepository struct {
	db *sql.DB
}

// NewAssetSnapshotRepository 建立新的資產快照 repository
func NewAssetSnapshotRepository(db *sql.DB) AssetSnapshotRepository {
	return &assetSnapshotRepository{db: db}
}

// Create 建立新的資產快照
func (r *assetSnapshotRepository) Create(input *models.CreateAssetSnapshotInput) (*models.AssetSnapshot, error) {
	query := `
		INSERT INTO asset_snapshots (snapshot_date, asset_type, value_twd)
		VALUES ($1, $2, $3)
		RETURNING id, snapshot_date, asset_type, value_twd, created_at, updated_at
	`

	snapshot := &models.AssetSnapshot{}
	err := r.db.QueryRow(
		query,
		input.SnapshotDate,
		input.AssetType,
		input.ValueTWD,
	).Scan(
		&snapshot.ID,
		&snapshot.SnapshotDate,
		&snapshot.AssetType,
		&snapshot.ValueTWD,
		&snapshot.CreatedAt,
		&snapshot.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create asset snapshot: %w", err)
	}

	return snapshot, nil
}

// GetByDateAndType 根據日期和資產類型取得快照
func (r *assetSnapshotRepository) GetByDateAndType(date time.Time, assetType models.SnapshotAssetType) (*models.AssetSnapshot, error) {
	query := `
		SELECT id, snapshot_date, asset_type, value_twd, created_at, updated_at
		FROM asset_snapshots
		WHERE snapshot_date = $1 AND asset_type = $2
	`

	snapshot := &models.AssetSnapshot{}
	err := r.db.QueryRow(query, date, assetType).Scan(
		&snapshot.ID,
		&snapshot.SnapshotDate,
		&snapshot.AssetType,
		&snapshot.ValueTWD,
		&snapshot.CreatedAt,
		&snapshot.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("asset snapshot not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get asset snapshot: %w", err)
	}

	return snapshot, nil
}

// GetByDateRange 根據日期範圍和篩選條件取得快照列表
func (r *assetSnapshotRepository) GetByDateRange(filters models.AssetSnapshotFilters) ([]*models.AssetSnapshot, error) {
	query := `
		SELECT id, snapshot_date, asset_type, value_twd, created_at, updated_at
		FROM asset_snapshots
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	// 建立動態查詢條件
	var conditions []string

	if filters.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("snapshot_date >= $%d", argIndex))
		args = append(args, *filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("snapshot_date <= $%d", argIndex))
		args = append(args, *filters.EndDate)
		argIndex++
	}

	if filters.AssetType != nil {
		conditions = append(conditions, fmt.Sprintf("asset_type = $%d", argIndex))
		args = append(args, *filters.AssetType)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// 按日期降序排列
	query += " ORDER BY snapshot_date DESC, asset_type ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query asset snapshots: %w", err)
	}
	defer rows.Close()

	snapshots := []*models.AssetSnapshot{}
	for rows.Next() {
		snapshot := &models.AssetSnapshot{}
		err := rows.Scan(
			&snapshot.ID,
			&snapshot.SnapshotDate,
			&snapshot.AssetType,
			&snapshot.ValueTWD,
			&snapshot.CreatedAt,
			&snapshot.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset snapshot: %w", err)
		}
		snapshots = append(snapshots, snapshot)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating asset snapshots: %w", err)
	}

	return snapshots, nil
}

// Update 更新資產快照的價值
func (r *assetSnapshotRepository) Update(date time.Time, assetType models.SnapshotAssetType, valueTWD float64) (*models.AssetSnapshot, error) {
	query := `
		UPDATE asset_snapshots
		SET value_twd = $1, updated_at = CURRENT_TIMESTAMP
		WHERE snapshot_date = $2 AND asset_type = $3
		RETURNING id, snapshot_date, asset_type, value_twd, created_at, updated_at
	`

	snapshot := &models.AssetSnapshot{}
	err := r.db.QueryRow(query, valueTWD, date, assetType).Scan(
		&snapshot.ID,
		&snapshot.SnapshotDate,
		&snapshot.AssetType,
		&snapshot.ValueTWD,
		&snapshot.CreatedAt,
		&snapshot.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("asset snapshot not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update asset snapshot: %w", err)
	}

	return snapshot, nil
}

// Delete 刪除資產快照
func (r *assetSnapshotRepository) Delete(date time.Time, assetType models.SnapshotAssetType) error {
	query := `
		DELETE FROM asset_snapshots
		WHERE snapshot_date = $1 AND asset_type = $2
	`

	result, err := r.db.Exec(query, date, assetType)
	if err != nil {
		return fmt.Errorf("failed to delete asset snapshot: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset snapshot not found")
	}

	return nil
}

