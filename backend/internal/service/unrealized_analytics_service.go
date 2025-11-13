package service

import (
	"sort"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// UnrealizedAnalyticsService 未實現損益分析服務介面
type UnrealizedAnalyticsService interface {
	// GetSummary 取得未實現損益摘要
	GetSummary() (*models.UnrealizedSummary, error)

	// GetPerformance 取得各資產類型未實現績效
	GetPerformance() ([]models.UnrealizedPerformance, error)

	// GetTopAssets 取得 Top 未實現損益資產
	GetTopAssets(limit int) ([]models.UnrealizedTopAsset, error)
}

// unrealizedAnalyticsService 未實現損益分析服務實作
type unrealizedAnalyticsService struct {
	holdingService HoldingService
}

// NewUnrealizedAnalyticsService 建立新的未實現損益分析服務
func NewUnrealizedAnalyticsService(holdingService HoldingService) UnrealizedAnalyticsService {
	return &unrealizedAnalyticsService{
		holdingService: holdingService,
	}
}

// GetSummary 取得未實現損益摘要
func (s *unrealizedAnalyticsService) GetSummary() (*models.UnrealizedSummary, error) {
	// 取得所有持倉
	result, err := s.holdingService.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		return nil, err
	}

	// 計算總計
	var totalCost, totalMarketValue, totalUnrealizedPL float64

	for _, h := range result.Holdings {
		totalCost += h.TotalCost
		totalMarketValue += h.MarketValue
		totalUnrealizedPL += h.UnrealizedPL
	}

	// 計算總報酬率
	totalUnrealizedPct := 0.0
	if totalCost > 0 {
		totalUnrealizedPct = (totalUnrealizedPL / totalCost) * 100
	}

	return &models.UnrealizedSummary{
		TotalCost:          totalCost,
		TotalMarketValue:   totalMarketValue,
		TotalUnrealizedPL:  totalUnrealizedPL,
		TotalUnrealizedPct: totalUnrealizedPct,
		HoldingCount:       len(result.Holdings),
		Currency:           "TWD",
	}, nil
}

// GetPerformance 取得各資產類型未實現績效
func (s *unrealizedAnalyticsService) GetPerformance() ([]models.UnrealizedPerformance, error) {
	// 取得所有持倉
	result, err := s.holdingService.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		return nil, err
	}

	// 按資產類型分組計算
	perfMap := make(map[models.AssetType]*models.UnrealizedPerformance)

	for _, h := range result.Holdings {
		// 如果該資產類型還沒有記錄，建立新的
		if _, exists := perfMap[h.AssetType]; !exists {
			perfMap[h.AssetType] = &models.UnrealizedPerformance{
				AssetType: h.AssetType,
				Name:      models.GetAssetTypeName(h.AssetType),
			}
		}

		// 累加數值
		perf := perfMap[h.AssetType]
		perf.Cost += h.TotalCost
		perf.MarketValue += h.MarketValue
		perf.UnrealizedPL += h.UnrealizedPL
		perf.HoldingCount++
	}

	// 計算百分比並轉換為陣列
	performances := make([]models.UnrealizedPerformance, 0, len(perfMap))
	for _, perf := range perfMap {
		if perf.Cost > 0 {
			perf.UnrealizedPct = (perf.UnrealizedPL / perf.Cost) * 100
		}
		performances = append(performances, *perf)
	}

	// 按資產類型排序（可選，保持一致性）
	sort.Slice(performances, func(i, j int) bool {
		return performances[i].AssetType < performances[j].AssetType
	})

	return performances, nil
}

// GetTopAssets 取得 Top 未實現損益資產
func (s *unrealizedAnalyticsService) GetTopAssets(limit int) ([]models.UnrealizedTopAsset, error) {
	// 取得所有持倉
	result, err := s.holdingService.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		return nil, err
	}

	// 轉換為 TopAsset 格式
	topAssets := make([]models.UnrealizedTopAsset, 0, len(result.Holdings))
	for _, h := range result.Holdings {
		topAssets = append(topAssets, models.UnrealizedTopAsset{
			Symbol:        h.Symbol,
			Name:          h.Name,
			AssetType:     h.AssetType,
			Quantity:      h.Quantity,
			AvgCost:       h.AvgCost,
			CurrentPrice:  h.CurrentPriceTWD,
			Cost:          h.TotalCost,
			MarketValue:   h.MarketValue,
			UnrealizedPL:  h.UnrealizedPL,
			UnrealizedPct: h.UnrealizedPLPct,
		})
	}

	// 按未實現損益降序排序
	sort.Slice(topAssets, func(i, j int) bool {
		return topAssets[i].UnrealizedPL > topAssets[j].UnrealizedPL
	})

	// 限制回傳數量
	if limit > 0 && limit < len(topAssets) {
		topAssets = topAssets[:limit]
	}

	return topAssets, nil
}

