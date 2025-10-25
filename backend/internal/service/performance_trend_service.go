package service

import (
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// PerformanceTrendService 績效趨勢服務介面
type PerformanceTrendService interface {
	// CreateDailySnapshot 建立每日績效快照
	CreateDailySnapshot() (*models.DailyPerformanceSnapshot, error)
	// GetTrendByDateRange 取得日期範圍內的績效趨勢
	GetTrendByDateRange(startDate, endDate time.Time) (*models.PerformanceTrendSummary, error)
	// GetLatestTrend 取得最新的 N 天績效趨勢
	GetLatestTrend(days int) ([]models.PerformanceTrendPoint, error)
}

type performanceTrendService struct {
	repo                repository.PerformanceSnapshotRepository
	unrealizedService   UnrealizedAnalyticsService
	analyticsService    AnalyticsService
}

// NewPerformanceTrendService 建立績效趨勢服務
func NewPerformanceTrendService(
	repo repository.PerformanceSnapshotRepository,
	unrealizedService UnrealizedAnalyticsService,
	analyticsService AnalyticsService,
) PerformanceTrendService {
	return &performanceTrendService{
		repo:              repo,
		unrealizedService: unrealizedService,
		analyticsService:  analyticsService,
	}
}

// CreateDailySnapshot 建立每日績效快照
func (s *performanceTrendService) CreateDailySnapshot() (*models.DailyPerformanceSnapshot, error) {
	// 取得未實現損益摘要
	unrealizedSummary, err := s.unrealizedService.GetSummary()
	if err != nil {
		return nil, fmt.Errorf("failed to get unrealized summary: %w", err)
	}

	// 取得已實現損益摘要（全部時間範圍）
	analyticsSummary, err := s.analyticsService.GetSummary(models.TimeRangeAll)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics summary: %w", err)
	}

	// 取得未實現損益績效（按資產類型）
	unrealizedPerformance, err := s.unrealizedService.GetPerformance()
	if err != nil {
		return nil, fmt.Errorf("failed to get unrealized performance: %w", err)
	}

	// 取得已實現損益績效（按資產類型）
	analyticsPerformancePtr, err := s.analyticsService.GetPerformance(models.TimeRangeAll)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics performance: %w", err)
	}

	// 轉換為非指標陣列
	analyticsPerformance := make([]models.PerformanceData, len(analyticsPerformancePtr))
	for i, p := range analyticsPerformancePtr {
		analyticsPerformance[i] = *p
	}

	// 建立明細資料
	details := s.buildSnapshotDetails(unrealizedPerformance, analyticsPerformance)

	// 計算總已實現報酬率
	totalRealizedPct := 0.0
	if analyticsSummary.TotalCostBasis > 0 {
		totalRealizedPct = (analyticsSummary.TotalRealizedPL / analyticsSummary.TotalCostBasis) * 100
	}

	// 建立快照輸入
	input := &models.CreateDailyPerformanceSnapshotInput{
		SnapshotDate:       time.Now().Truncate(24 * time.Hour),
		TotalMarketValue:   unrealizedSummary.TotalMarketValue,
		TotalCost:          unrealizedSummary.TotalCost,
		TotalUnrealizedPL:  unrealizedSummary.TotalUnrealizedPL,
		TotalUnrealizedPct: unrealizedSummary.TotalUnrealizedPct,
		TotalRealizedPL:    analyticsSummary.TotalRealizedPL,
		TotalRealizedPct:   totalRealizedPct,
		HoldingCount:       unrealizedSummary.HoldingCount,
		Currency:           "TWD",
		Details:            details,
	}

	// 建立快照
	snapshot, err := s.repo.Create(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	return snapshot, nil
}

// buildSnapshotDetails 建立快照明細
func (s *performanceTrendService) buildSnapshotDetails(
	unrealizedPerformance []models.UnrealizedPerformance,
	analyticsPerformance []models.PerformanceData,
) []models.CreateDailyPerformanceSnapshotDetailInput {
	// 建立已實現損益的 map
	realizedMap := make(map[models.AssetType]models.PerformanceData)
	for _, perf := range analyticsPerformance {
		realizedMap[perf.AssetType] = perf
	}

	// 合併未實現和已實現資料
	details := make([]models.CreateDailyPerformanceSnapshotDetailInput, 0, len(unrealizedPerformance))
	for _, unrealized := range unrealizedPerformance {
		realized, exists := realizedMap[unrealized.AssetType]
		
		realizedPL := 0.0
		realizedPct := 0.0
		if exists {
			realizedPL = realized.RealizedPL
			realizedPct = realized.RealizedPLPct
		}

		detail := models.CreateDailyPerformanceSnapshotDetailInput{
			AssetType:     unrealized.AssetType,
			MarketValue:   unrealized.MarketValue,
			Cost:          unrealized.Cost,
			UnrealizedPL:  unrealized.UnrealizedPL,
			UnrealizedPct: unrealized.UnrealizedPct,
			RealizedPL:    realizedPL,
			RealizedPct:   realizedPct,
			HoldingCount:  unrealized.HoldingCount,
		}
		details = append(details, detail)
	}

	return details
}

// GetTrendByDateRange 取得日期範圍內的績效趨勢
func (s *performanceTrendService) GetTrendByDateRange(startDate, endDate time.Time) (*models.PerformanceTrendSummary, error) {
	// 取得快照資料
	snapshots, err := s.repo.GetByDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots: %w", err)
	}

	if len(snapshots) == 0 {
		return &models.PerformanceTrendSummary{
			StartDate:      startDate,
			EndDate:        endDate,
			TotalData:      []models.PerformanceTrendPoint{},
			ByType:         []models.PerformanceTrendByType{},
			Currency:       "TWD",
			DataPointCount: 0,
		}, nil
	}

	// 取得明細資料
	detailsMap, err := s.repo.GetDetailsByDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get details: %w", err)
	}

	// 建立總體趨勢資料
	totalData := make([]models.PerformanceTrendPoint, 0, len(snapshots))
	for _, snapshot := range snapshots {
		point := models.PerformanceTrendPoint{
			Date:           snapshot.SnapshotDate,
			MarketValue:    snapshot.TotalMarketValue,
			Cost:           snapshot.TotalCost,
			UnrealizedPL:   snapshot.TotalUnrealizedPL,
			UnrealizedPct:  snapshot.TotalUnrealizedPct,
			RealizedPL:     snapshot.TotalRealizedPL,
			RealizedPct:    snapshot.TotalRealizedPct,
			TotalPL:        snapshot.TotalUnrealizedPL + snapshot.TotalRealizedPL,
			TotalPct:       snapshot.TotalUnrealizedPct + snapshot.TotalRealizedPct,
			HoldingCount:   snapshot.HoldingCount,
		}
		totalData = append(totalData, point)
	}

	// 建立按資產類型的趨勢資料
	byType := s.buildTrendByType(snapshots, detailsMap)

	return &models.PerformanceTrendSummary{
		StartDate:      startDate,
		EndDate:        endDate,
		TotalData:      totalData,
		ByType:         byType,
		Currency:       snapshots[0].Currency,
		DataPointCount: len(snapshots),
	}, nil
}

// buildTrendByType 建立按資產類型的趨勢資料
func (s *performanceTrendService) buildTrendByType(
	snapshots []*models.DailyPerformanceSnapshot,
	detailsMap map[uuid.UUID][]*models.DailyPerformanceSnapshotDetail,
) []models.PerformanceTrendByType {
	// 收集所有資產類型
	assetTypes := []models.AssetType{
		models.AssetTypeTWStock,
		models.AssetTypeUSStock,
		models.AssetTypeCrypto,
	}

	result := make([]models.PerformanceTrendByType, 0, len(assetTypes))
	for _, assetType := range assetTypes {
		data := make([]models.PerformanceTrendPoint, 0, len(snapshots))
		
		for _, snapshot := range snapshots {
			// 找到對應的明細
			details := detailsMap[snapshot.ID]
			var detail *models.DailyPerformanceSnapshotDetail
			for _, d := range details {
				if d.AssetType == assetType {
					detail = d
					break
				}
			}

			// 如果沒有明細，使用零值
			if detail == nil {
				point := models.PerformanceTrendPoint{
					Date:           snapshot.SnapshotDate,
					MarketValue:    0,
					Cost:           0,
					UnrealizedPL:   0,
					UnrealizedPct:  0,
					RealizedPL:     0,
					RealizedPct:    0,
					TotalPL:        0,
					TotalPct:       0,
					HoldingCount:   0,
				}
				data = append(data, point)
			} else {
				point := models.PerformanceTrendPoint{
					Date:           snapshot.SnapshotDate,
					MarketValue:    detail.MarketValue,
					Cost:           detail.Cost,
					UnrealizedPL:   detail.UnrealizedPL,
					UnrealizedPct:  detail.UnrealizedPct,
					RealizedPL:     detail.RealizedPL,
					RealizedPct:    detail.RealizedPct,
					TotalPL:        detail.UnrealizedPL + detail.RealizedPL,
					TotalPct:       detail.UnrealizedPct + detail.RealizedPct,
					HoldingCount:   detail.HoldingCount,
				}
				data = append(data, point)
			}
		}

		result = append(result, models.PerformanceTrendByType{
			AssetType: assetType,
			Name:      models.AssetTypeNameMap[assetType],
			Data:      data,
		})
	}

	return result
}

// GetLatestTrend 取得最新的 N 天績效趨勢
func (s *performanceTrendService) GetLatestTrend(days int) ([]models.PerformanceTrendPoint, error) {
	// 取得最新的快照
	snapshots, err := s.repo.GetLatest(days)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest snapshots: %w", err)
	}

	// 轉換為趨勢資料點
	data := make([]models.PerformanceTrendPoint, 0, len(snapshots))
	for _, snapshot := range snapshots {
		point := models.PerformanceTrendPoint{
			Date:           snapshot.SnapshotDate,
			MarketValue:    snapshot.TotalMarketValue,
			Cost:           snapshot.TotalCost,
			UnrealizedPL:   snapshot.TotalUnrealizedPL,
			UnrealizedPct:  snapshot.TotalUnrealizedPct,
			RealizedPL:     snapshot.TotalRealizedPL,
			RealizedPct:    snapshot.TotalRealizedPct,
			TotalPL:        snapshot.TotalUnrealizedPL + snapshot.TotalRealizedPL,
			TotalPct:       snapshot.TotalUnrealizedPct + snapshot.TotalRealizedPct,
			HoldingCount:   snapshot.HoldingCount,
		}
		data = append(data, point)
	}

	return data, nil
}

