package service

import (
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
)

// AssetSnapshotService 資產快照服務介面
type AssetSnapshotService interface {
	// CreateSnapshot 建立資產快照
	CreateSnapshot(input *models.CreateAssetSnapshotInput) (*models.AssetSnapshot, error)

	// GetSnapshotByDate 根據日期和資產類型取得快照
	GetSnapshotByDate(date time.Time, assetType models.SnapshotAssetType) (*models.AssetSnapshot, error)

	// GetSnapshotsByDateRange 根據日期範圍取得快照列表
	GetSnapshotsByDateRange(startDate, endDate time.Time, assetType models.SnapshotAssetType) ([]*models.AssetSnapshot, error)

	// GetLatestSnapshot 取得最新的快照
	GetLatestSnapshot(assetType models.SnapshotAssetType) (*models.AssetSnapshot, error)

	// UpdateSnapshot 更新快照
	UpdateSnapshot(date time.Time, assetType models.SnapshotAssetType, valueTWD float64) (*models.AssetSnapshot, error)

	// DeleteSnapshot 刪除快照
	DeleteSnapshot(date time.Time, assetType models.SnapshotAssetType) error
}

// assetSnapshotService 資產快照服務實作
type assetSnapshotService struct {
	repo repository.AssetSnapshotRepository
}

// NewAssetSnapshotService 建立資產快照服務
func NewAssetSnapshotService(repo repository.AssetSnapshotRepository) AssetSnapshotService {
	return &assetSnapshotService{
		repo: repo,
	}
}

// CreateSnapshot 建立資產快照
func (s *assetSnapshotService) CreateSnapshot(input *models.CreateAssetSnapshotInput) (*models.AssetSnapshot, error) {
	// 驗證輸入
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// 建立快照
	snapshot, err := s.repo.Create(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	return snapshot, nil
}

// GetSnapshotByDate 根據日期和資產類型取得快照
func (s *assetSnapshotService) GetSnapshotByDate(date time.Time, assetType models.SnapshotAssetType) (*models.AssetSnapshot, error) {
	// 將日期標準化為當天的 00:00:00
	normalizedDate := date.Truncate(24 * time.Hour)

	snapshot, err := s.repo.GetByDateAndType(normalizedDate, assetType)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	return snapshot, nil
}

// GetSnapshotsByDateRange 根據日期範圍取得快照列表
func (s *assetSnapshotService) GetSnapshotsByDateRange(startDate, endDate time.Time, assetType models.SnapshotAssetType) ([]*models.AssetSnapshot, error) {
	// 驗證日期範圍
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date must be before or equal to end date")
	}

	// 將日期標準化
	normalizedStartDate := startDate.Truncate(24 * time.Hour)
	normalizedEndDate := endDate.Truncate(24 * time.Hour)

	assetTypePtr := assetType
	filters := models.AssetSnapshotFilters{
		StartDate: &normalizedStartDate,
		EndDate:   &normalizedEndDate,
		AssetType: &assetTypePtr,
	}

	snapshots, err := s.repo.GetByDateRange(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots: %w", err)
	}

	return snapshots, nil
}

// GetLatestSnapshot 取得最新的快照
func (s *assetSnapshotService) GetLatestSnapshot(assetType models.SnapshotAssetType) (*models.AssetSnapshot, error) {
	// 使用日期範圍查詢，取得最近 30 天的資料
	endDate := time.Now().Truncate(24 * time.Hour)
	startDate := endDate.Add(-30 * 24 * time.Hour)

	assetTypePtr := assetType
	filters := models.AssetSnapshotFilters{
		StartDate: &startDate,
		EndDate:   &endDate,
		AssetType: &assetTypePtr,
	}

	snapshots, err := s.repo.GetByDateRange(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots: %w", err)
	}

	if len(snapshots) == 0 {
		return nil, fmt.Errorf("no snapshot found")
	}

	// 返回最新的快照（列表已按日期降序排列）
	return snapshots[0], nil
}

// UpdateSnapshot 更新快照
func (s *assetSnapshotService) UpdateSnapshot(date time.Time, assetType models.SnapshotAssetType, valueTWD float64) (*models.AssetSnapshot, error) {
	// 驗證金額
	if valueTWD < 0 {
		return nil, fmt.Errorf("value_twd must be non-negative")
	}

	// 將日期標準化
	normalizedDate := date.Truncate(24 * time.Hour)

	snapshot, err := s.repo.Update(normalizedDate, assetType, valueTWD)
	if err != nil {
		return nil, fmt.Errorf("failed to update snapshot: %w", err)
	}

	return snapshot, nil
}

// DeleteSnapshot 刪除快照
func (s *assetSnapshotService) DeleteSnapshot(date time.Time, assetType models.SnapshotAssetType) error {
	// 將日期標準化
	normalizedDate := date.Truncate(24 * time.Hour)

	if err := s.repo.Delete(normalizedDate, assetType); err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	return nil
}

