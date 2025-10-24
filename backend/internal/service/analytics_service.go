package service

import (
	"fmt"
	"sort"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
)

// AnalyticsService 分析服務介面
type AnalyticsService interface {
	// GetSummary 取得分析摘要
	GetSummary(timeRange models.TimeRange) (*models.AnalyticsSummary, error)

	// GetPerformance 取得各資產類型績效
	GetPerformance(timeRange models.TimeRange) ([]*models.PerformanceData, error)

	// GetTopAssets 取得最佳/最差表現資產
	GetTopAssets(timeRange models.TimeRange, limit int) ([]*models.TopAsset, error)
}

// analyticsService 分析服務實作
type analyticsService struct {
	realizedProfitRepo repository.RealizedProfitRepository
}

// NewAnalyticsService 建立新的分析服務
func NewAnalyticsService(realizedProfitRepo repository.RealizedProfitRepository) AnalyticsService {
	return &analyticsService{
		realizedProfitRepo: realizedProfitRepo,
	}
}

// GetSummary 取得分析摘要
func (s *analyticsService) GetSummary(timeRange models.TimeRange) (*models.AnalyticsSummary, error) {
	// 驗證時間範圍
	if !timeRange.Validate() {
		return nil, fmt.Errorf("invalid time range: %s", timeRange)
	}

	// 取得時間範圍的起始和結束日期
	startDate, endDate := timeRange.GetDateRange()

	// 查詢已實現損益記錄
	filters := models.RealizedProfitFilters{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	records, err := s.realizedProfitRepo.GetAll(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get realized profits: %w", err)
	}

	// 計算摘要資料
	summary := &models.AnalyticsSummary{
		TimeRange:        string(timeRange),
		StartDate:        startDate.Format("2006-01-02"),
		EndDate:          endDate.Format("2006-01-02"),
		Currency:         "TWD", // 預設使用 TWD
		TransactionCount: len(records),
	}

	for _, record := range records {
		summary.TotalRealizedPL += record.RealizedPL
		summary.TotalCostBasis += record.CostBasis
		summary.TotalSellAmount += record.SellAmount
		summary.TotalSellFee += record.SellFee
	}

	// 計算總已實現損益百分比
	if summary.TotalCostBasis > 0 {
		summary.TotalRealizedPLPct = (summary.TotalRealizedPL / summary.TotalCostBasis) * 100
	}

	return summary, nil
}

// GetPerformance 取得各資產類型績效
func (s *analyticsService) GetPerformance(timeRange models.TimeRange) ([]*models.PerformanceData, error) {
	// 驗證時間範圍
	if !timeRange.Validate() {
		return nil, fmt.Errorf("invalid time range: %s", timeRange)
	}

	// 取得時間範圍的起始和結束日期
	startDate, endDate := timeRange.GetDateRange()

	// 查詢已實現損益記錄
	filters := models.RealizedProfitFilters{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	records, err := s.realizedProfitRepo.GetAll(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get realized profits: %w", err)
	}

	// 按資產類型分組統計
	performanceMap := make(map[models.AssetType]*models.PerformanceData)

	for _, record := range records {
		if _, exists := performanceMap[record.AssetType]; !exists {
			performanceMap[record.AssetType] = &models.PerformanceData{
				AssetType: record.AssetType,
				Name:      models.GetAssetTypeName(record.AssetType),
			}
		}

		perf := performanceMap[record.AssetType]
		perf.RealizedPL += record.RealizedPL
		perf.CostBasis += record.CostBasis
		perf.SellAmount += record.SellAmount
		perf.TransactionCount++
	}

	// 計算各資產類型的已實現損益百分比
	for _, perf := range performanceMap {
		if perf.CostBasis > 0 {
			perf.RealizedPLPct = (perf.RealizedPL / perf.CostBasis) * 100
		}
	}

	// 轉換為陣列並排序（按已實現損益由高到低）
	performance := make([]*models.PerformanceData, 0, len(performanceMap))
	for _, perf := range performanceMap {
		performance = append(performance, perf)
	}

	sort.Slice(performance, func(i, j int) bool {
		return performance[i].RealizedPL > performance[j].RealizedPL
	})

	return performance, nil
}

// GetTopAssets 取得最佳/最差表現資產
func (s *analyticsService) GetTopAssets(timeRange models.TimeRange, limit int) ([]*models.TopAsset, error) {
	// 驗證時間範圍
	if !timeRange.Validate() {
		return nil, fmt.Errorf("invalid time range: %s", timeRange)
	}

	// 驗證 limit
	if limit <= 0 {
		limit = 5 // 預設 5 筆
	}

	// 取得時間範圍的起始和結束日期
	startDate, endDate := timeRange.GetDateRange()

	// 查詢已實現損益記錄
	filters := models.RealizedProfitFilters{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	records, err := s.realizedProfitRepo.GetAll(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get realized profits: %w", err)
	}

	// 按標的分組統計
	assetMap := make(map[string]*models.TopAsset)

	for _, record := range records {
		if _, exists := assetMap[record.Symbol]; !exists {
			assetMap[record.Symbol] = &models.TopAsset{
				Symbol:    record.Symbol,
				Name:      record.Symbol, // 暫時使用 Symbol 作為 Name
				AssetType: record.AssetType,
			}
		}

		asset := assetMap[record.Symbol]
		asset.RealizedPL += record.RealizedPL
		asset.CostBasis += record.CostBasis
		asset.SellAmount += record.SellAmount
	}

	// 計算各標的的已實現損益百分比
	for _, asset := range assetMap {
		if asset.CostBasis > 0 {
			asset.RealizedPLPct = (asset.RealizedPL / asset.CostBasis) * 100
		}
	}

	// 轉換為陣列並排序（按已實現損益由高到低）
	topAssets := make([]*models.TopAsset, 0, len(assetMap))
	for _, asset := range assetMap {
		topAssets = append(topAssets, asset)
	}

	sort.Slice(topAssets, func(i, j int) bool {
		return topAssets[i].RealizedPL > topAssets[j].RealizedPL
	})

	// 限制回傳數量
	if len(topAssets) > limit {
		topAssets = topAssets[:limit]
	}

	return topAssets, nil
}

