package service

import (
	"fmt"
	"math"
	"sort"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// RebalanceService 再平衡服務介面
type RebalanceService interface {
	// CheckRebalance 檢查是否需要再平衡
	CheckRebalance() (*models.RebalanceCheck, error)
}

// rebalanceService 再平衡服務實作
type rebalanceService struct {
	settingsService SettingsService
	holdingService  HoldingService
}

// NewRebalanceService 建立再平衡服務
func NewRebalanceService(settingsService SettingsService, holdingService HoldingService) RebalanceService {
	return &rebalanceService{
		settingsService: settingsService,
		holdingService:  holdingService,
	}
}

// CheckRebalance 檢查是否需要再平衡
func (s *rebalanceService) CheckRebalance() (*models.RebalanceCheck, error) {
	// 1. 取得設定
	settings, err := s.settingsService.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	// 2. 取得所有持倉
	holdings, err := s.holdingService.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to get holdings: %w", err)
	}

	// 3. 如果沒有持倉，返回空結果
	if len(holdings) == 0 {
		return &models.RebalanceCheck{
			NeedsRebalance: false,
			Threshold:      settings.Allocation.RebalanceThreshold,
			Deviations:     []models.AssetTypeDeviation{},
			Suggestions:    []models.RebalanceSuggestion{},
			CurrentTotal:   0,
		}, nil
	}

	// 4. 計算當前配置
	currentAllocation := s.calculateCurrentAllocation(holdings)

	// 5. 計算偏離情況
	deviations := s.calculateDeviations(settings.Allocation, currentAllocation)

	// 6. 檢查是否需要再平衡
	needsRebalance := s.checkIfRebalanceNeeded(deviations, settings.Allocation.RebalanceThreshold)

	// 7. 如果需要再平衡，生成建議
	suggestions := []models.RebalanceSuggestion{}
	if needsRebalance {
		suggestions = s.generateSuggestions(deviations, currentAllocation.TotalValue)
	}

	return &models.RebalanceCheck{
		NeedsRebalance: needsRebalance,
		Threshold:      settings.Allocation.RebalanceThreshold,
		Deviations:     deviations,
		Suggestions:    suggestions,
		CurrentTotal:   currentAllocation.TotalValue,
	}, nil
}

// currentAllocationData 當前配置資料
type currentAllocationData struct {
	TotalValue float64
	ByType     map[string]float64 // 資產類型 -> 市值
}

// calculateCurrentAllocation 計算當前配置
func (s *rebalanceService) calculateCurrentAllocation(holdings []*models.Holding) currentAllocationData {
	result := currentAllocationData{
		TotalValue: 0,
		ByType:     make(map[string]float64),
	}

	for _, holding := range holdings {
		assetTypeStr := string(holding.AssetType)
		result.ByType[assetTypeStr] += holding.MarketValue
		result.TotalValue += holding.MarketValue
	}

	return result
}

// calculateDeviations 計算偏離情況
func (s *rebalanceService) calculateDeviations(
	targetAllocation models.AllocationSettings,
	currentAllocation currentAllocationData,
) []models.AssetTypeDeviation {
	deviations := []models.AssetTypeDeviation{}

	// 定義所有資產類型和目標配置
	assetTypes := map[string]float64{
		"tw-stock": targetAllocation.TWStock,
		"us-stock": targetAllocation.USStock,
		"crypto":   targetAllocation.Crypto,
	}

	// 計算每個資產類型的偏離
	for assetType, targetPercent := range assetTypes {
		currentValue := currentAllocation.ByType[assetType]
		currentPercent := 0.0
		if currentAllocation.TotalValue > 0 {
			currentPercent = (currentValue / currentAllocation.TotalValue) * 100
		}

		targetValue := (targetPercent / 100) * currentAllocation.TotalValue
		deviation := currentPercent - targetPercent
		deviationAbs := math.Abs(deviation)

		deviations = append(deviations, models.AssetTypeDeviation{
			AssetType:        assetType,
			TargetPercent:    targetPercent,
			CurrentPercent:   currentPercent,
			Deviation:        deviation,
			DeviationAbs:     deviationAbs,
			ExceedsThreshold: false, // 稍後會設定
			CurrentValue:     currentValue,
			TargetValue:      targetValue,
		})
	}

	// 按資產類型排序（保持一致性）
	sort.Slice(deviations, func(i, j int) bool {
		return deviations[i].AssetType < deviations[j].AssetType
	})

	return deviations
}

// checkIfRebalanceNeeded 檢查是否需要再平衡
func (s *rebalanceService) checkIfRebalanceNeeded(
	deviations []models.AssetTypeDeviation,
	threshold float64,
) bool {
	for i := range deviations {
		if deviations[i].DeviationAbs > threshold {
			deviations[i].ExceedsThreshold = true
			return true
		}
	}
	return false
}

// generateSuggestions 生成再平衡建議
func (s *rebalanceService) generateSuggestions(
	deviations []models.AssetTypeDeviation,
	totalValue float64,
) []models.RebalanceSuggestion {
	suggestions := []models.RebalanceSuggestion{}

	for _, deviation := range deviations {
		if !deviation.ExceedsThreshold {
			continue
		}

		// 計算需要調整的金額
		adjustAmount := math.Abs(deviation.CurrentValue - deviation.TargetValue)

		// 判斷動作
		action := "buy"
		reason := fmt.Sprintf("當前配置 %.2f%% 低於目標 %.2f%%，建議買入", deviation.CurrentPercent, deviation.TargetPercent)
		if deviation.Deviation > 0 {
			action = "sell"
			reason = fmt.Sprintf("當前配置 %.2f%% 高於目標 %.2f%%，建議賣出", deviation.CurrentPercent, deviation.TargetPercent)
		}

		suggestions = append(suggestions, models.RebalanceSuggestion{
			AssetType: deviation.AssetType,
			Action:    action,
			Amount:    adjustAmount,
			Reason:    reason,
		})
	}

	// 按金額降序排序（優先處理偏離最大的）
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Amount > suggestions[j].Amount
	})

	return suggestions
}

