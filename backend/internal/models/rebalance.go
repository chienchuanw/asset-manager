package models

// RebalanceCheck 再平衡檢查結果
type RebalanceCheck struct {
	NeedsRebalance bool                      `json:"needs_rebalance"` // 是否需要再平衡
	Threshold      float64                   `json:"threshold"`       // 閾值百分比
	Deviations     []AssetTypeDeviation      `json:"deviations"`      // 各資產類型的偏離情況
	Suggestions    []RebalanceSuggestion     `json:"suggestions"`     // 再平衡建議
	CurrentTotal   float64                   `json:"current_total"`   // 當前總資產
}

// AssetTypeDeviation 資產類型偏離情況
type AssetTypeDeviation struct {
	AssetType         string  `json:"asset_type"`          // 資產類型
	TargetPercent     float64 `json:"target_percent"`      // 目標配置百分比
	CurrentPercent    float64 `json:"current_percent"`     // 當前配置百分比
	Deviation         float64 `json:"deviation"`           // 偏離百分比（當前 - 目標）
	DeviationAbs      float64 `json:"deviation_abs"`       // 絕對偏離百分比
	ExceedsThreshold  bool    `json:"exceeds_threshold"`   // 是否超過閾值
	CurrentValue      float64 `json:"current_value"`       // 當前市值
	TargetValue       float64 `json:"target_value"`        // 目標市值
}

// RebalanceSuggestion 再平衡建議
type RebalanceSuggestion struct {
	AssetType string  `json:"asset_type"` // 資產類型
	Action    string  `json:"action"`     // 動作：buy（買入）或 sell（賣出）
	Amount    float64 `json:"amount"`     // 金額（TWD）
	Reason    string  `json:"reason"`     // 原因說明
}

