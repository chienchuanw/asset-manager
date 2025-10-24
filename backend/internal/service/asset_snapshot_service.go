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

	// CreateDailySnapshots 建立每日資產快照（所有類型）
	CreateDailySnapshots() error
}

// assetSnapshotService 資產快照服務實作
type assetSnapshotService struct {
	repo           repository.AssetSnapshotRepository
	holdingService HoldingService // 用於計算持倉價值
}

// NewAssetSnapshotService 建立資產快照服務
func NewAssetSnapshotService(repo repository.AssetSnapshotRepository) AssetSnapshotService {
	return &assetSnapshotService{
		repo: repo,
	}
}

// NewAssetSnapshotServiceWithDeps 建立資產快照服務（包含依賴）
func NewAssetSnapshotServiceWithDeps(repo repository.AssetSnapshotRepository, holdingService HoldingService) AssetSnapshotService {
	return &assetSnapshotService{
		repo:           repo,
		holdingService: holdingService,
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

// CreateDailySnapshots 建立每日資產快照（所有類型）
func (s *assetSnapshotService) CreateDailySnapshots() error {
	// 檢查是否有 HoldingService 依賴
	if s.holdingService == nil {
		return fmt.Errorf("holding service is required for creating daily snapshots")
	}

	today := time.Now().Truncate(24 * time.Hour)

	// 取得所有持倉（不使用篩選條件）
	holdings, err := s.holdingService.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		return fmt.Errorf("failed to get holdings: %w", err)
	}

	// 計算各類型資產的總價值
	var totalValueTWD float64
	var twStockValueTWD float64
	var usStockValueTWD float64
	var cryptoValueTWD float64

	for _, holding := range holdings {
		// 將所有價值轉換為 TWD
		valueTWD := holding.MarketValue
		if holding.Currency == "USD" {
			// 這裡應該使用匯率服務轉換,暫時使用固定匯率 31.5
			valueTWD = holding.MarketValue * 31.5
		}

		totalValueTWD += valueTWD

		// 根據資產類型累加
		switch holding.AssetType {
		case models.AssetTypeTWStock:
			twStockValueTWD += valueTWD
		case models.AssetTypeUSStock:
			usStockValueTWD += valueTWD
		case models.AssetTypeCrypto:
			cryptoValueTWD += valueTWD
		}
	}

	// 建立各類型快照
	snapshots := []struct {
		assetType models.SnapshotAssetType
		value     float64
	}{
		{models.SnapshotAssetTypeTotal, totalValueTWD},
		{models.SnapshotAssetTypeTWStock, twStockValueTWD},
		{models.SnapshotAssetTypeUSStock, usStockValueTWD},
		{models.SnapshotAssetTypeCrypto, cryptoValueTWD},
	}

	// 建立或更新快照
	for _, snapshot := range snapshots {
		// 檢查是否已存在今日快照
		existing, err := s.repo.GetByDateAndType(today, snapshot.assetType)
		if err == nil && existing != nil {
			// 已存在,更新
			_, err = s.repo.Update(today, snapshot.assetType, snapshot.value)
			if err != nil {
				return fmt.Errorf("failed to update snapshot for %s: %w", snapshot.assetType, err)
			}
		} else {
			// 不存在,建立新的
			input := &models.CreateAssetSnapshotInput{
				SnapshotDate: today,
				AssetType:    snapshot.assetType,
				ValueTWD:     snapshot.value,
			}
			_, err = s.repo.Create(input)
			if err != nil {
				return fmt.Errorf("failed to create snapshot for %s: %w", snapshot.assetType, err)
			}
		}
	}

	return nil
}

