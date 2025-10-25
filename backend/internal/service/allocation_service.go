package service

import (
	"sort"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// AllocationService 資產配置服務介面
type AllocationService interface {
	// GetCurrentAllocation 取得當前資產配置摘要
	GetCurrentAllocation() (*models.AllocationSummary, error)

	// GetAllocationByType 取得按資產類型的配置
	GetAllocationByType() ([]models.AllocationByType, error)

	// GetAllocationByAsset 取得按個別資產的配置
	GetAllocationByAsset(limit int) ([]models.AllocationByAsset, error)
}

// allocationService 資產配置服務實作
type allocationService struct {
	holdingService HoldingService
}

// NewAllocationService 建立資產配置服務
func NewAllocationService(holdingService HoldingService) AllocationService {
	return &allocationService{
		holdingService: holdingService,
	}
}

// GetCurrentAllocation 取得當前資產配置摘要
func (s *allocationService) GetCurrentAllocation() (*models.AllocationSummary, error) {
	// 取得所有持倉
	holdings, err := s.holdingService.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		return nil, err
	}

	// 計算總市值
	var totalMarketValue float64
	for _, holding := range holdings {
		totalMarketValue += holding.MarketValue
	}

	// 計算按資產類型的配置
	byType, err := s.calculateAllocationByType(holdings, totalMarketValue)
	if err != nil {
		return nil, err
	}

	// 計算按個別資產的配置
	byAsset, err := s.calculateAllocationByAsset(holdings, totalMarketValue, 0)
	if err != nil {
		return nil, err
	}

	return &models.AllocationSummary{
		TotalMarketValue: totalMarketValue,
		ByType:           byType,
		ByAsset:          byAsset,
		Currency:         "TWD",
		AsOfDate:         time.Now(),
	}, nil
}

// GetAllocationByType 取得按資產類型的配置
func (s *allocationService) GetAllocationByType() ([]models.AllocationByType, error) {
	// 取得所有持倉
	holdings, err := s.holdingService.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		return nil, err
	}

	// 計算總市值
	var totalMarketValue float64
	for _, holding := range holdings {
		totalMarketValue += holding.MarketValue
	}

	return s.calculateAllocationByType(holdings, totalMarketValue)
}

// GetAllocationByAsset 取得按個別資產的配置
func (s *allocationService) GetAllocationByAsset(limit int) ([]models.AllocationByAsset, error) {
	// 取得所有持倉
	holdings, err := s.holdingService.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		return nil, err
	}

	// 計算總市值
	var totalMarketValue float64
	for _, holding := range holdings {
		totalMarketValue += holding.MarketValue
	}

	return s.calculateAllocationByAsset(holdings, totalMarketValue, limit)
}

// calculateAllocationByType 計算按資產類型的配置
func (s *allocationService) calculateAllocationByType(holdings []*models.Holding, totalMarketValue float64) ([]models.AllocationByType, error) {
	// 按資產類型分組
	typeMap := make(map[models.AssetType]*models.AllocationByType)

	for _, holding := range holdings {
		if _, exists := typeMap[holding.AssetType]; !exists {
			typeMap[holding.AssetType] = &models.AllocationByType{
				AssetType:   holding.AssetType,
				Name:        getAssetTypeName(holding.AssetType),
				MarketValue: 0,
				Percentage:  0,
				Count:       0,
			}
		}

		typeMap[holding.AssetType].MarketValue += holding.MarketValue
		typeMap[holding.AssetType].Count++
	}

	// 計算百分比
	result := make([]models.AllocationByType, 0, len(typeMap))
	for _, allocation := range typeMap {
		if totalMarketValue > 0 {
			allocation.Percentage = (allocation.MarketValue / totalMarketValue) * 100
		}
		result = append(result, *allocation)
	}

	// 按市值降序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].MarketValue > result[j].MarketValue
	})

	return result, nil
}

// calculateAllocationByAsset 計算按個別資產的配置
func (s *allocationService) calculateAllocationByAsset(holdings []*models.Holding, totalMarketValue float64, limit int) ([]models.AllocationByAsset, error) {
	result := make([]models.AllocationByAsset, 0, len(holdings))

	for _, holding := range holdings {
		percentage := 0.0
		if totalMarketValue > 0 {
			percentage = (holding.MarketValue / totalMarketValue) * 100
		}

		result = append(result, models.AllocationByAsset{
			Symbol:      holding.Symbol,
			Name:        holding.Name,
			AssetType:   holding.AssetType,
			MarketValue: holding.MarketValue,
			Percentage:  percentage,
			Quantity:    holding.Quantity,
		})
	}

	// 按市值降序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].MarketValue > result[j].MarketValue
	})

	// 限制回傳數量
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// getAssetTypeName 取得資產類型名稱
func getAssetTypeName(assetType models.AssetType) string {
	switch assetType {
	case models.AssetTypeTWStock:
		return "台股"
	case models.AssetTypeUSStock:
		return "美股"
	case models.AssetTypeCrypto:
		return "加密貨幣"
	default:
		return string(assetType)
	}
}

