/**
 * 資產再平衡相關型別定義
 */

/**
 * 資產類別偏差
 */
export interface AssetTypeDeviation {
  asset_type: string;
  target_percent: number;
  current_percent: number;
  deviation: number;
  deviation_abs: number;
  exceeds_threshold: boolean;
  current_value: number;
  target_value: number;
}

/**
 * 再平衡建議
 */
export interface RebalanceSuggestion {
  asset_type: string;
  action: string;
  amount: number;
  reason: string;
}

/**
 * 再平衡檢查結果
 */
export interface RebalanceCheck {
  needs_rebalance: boolean;
  threshold: number;
  deviations: AssetTypeDeviation[];
  suggestions: RebalanceSuggestion[];
  current_total: number;
}
