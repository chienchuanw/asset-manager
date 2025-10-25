import { AssetType } from "./transaction";

// ==================== 列舉型別 ====================

/**
 * 時間範圍
 */
export const TimeRange = {
  WEEK: "week",
  MONTH: "month",
  QUARTER: "quarter",
  YEAR: "year",
  ALL: "all",
} as const;

export type TimeRange = (typeof TimeRange)[keyof typeof TimeRange];

// ==================== 資料模型 ====================

/**
 * 分析摘要
 */
export interface AnalyticsSummary {
  total_realized_pl: number;
  total_realized_pl_pct: number;
  total_cost_basis: number;
  total_sell_amount: number;
  total_sell_fee: number;
  transaction_count: number;
  currency: string;
  time_range: string;
  start_date: string;
  end_date: string;
}

/**
 * 績效資料
 */
export interface PerformanceData {
  asset_type: AssetType;
  name: string;
  realized_pl: number;
  realized_pl_pct: number;
  cost_basis: number;
  sell_amount: number;
  transaction_count: number;
}

/**
 * 最佳/最差表現資產
 */
export interface TopAsset {
  symbol: string;
  name: string;
  asset_type: AssetType;
  realized_pl: number;
  realized_pl_pct: number;
  cost_basis: number;
  sell_amount: number;
}

// ==================== 輔助函式 ====================

/**
 * 取得時間範圍的顯示名稱
 */
export function getTimeRangeLabel(timeRange: TimeRange): string {
  const labels: Record<TimeRange, string> = {
    [TimeRange.WEEK]: "本週",
    [TimeRange.MONTH]: "本月",
    [TimeRange.QUARTER]: "本季",
    [TimeRange.YEAR]: "本年",
    [TimeRange.ALL]: "全部",
  };
  return labels[timeRange];
}

/**
 * 取得所有時間範圍選項
 */
export function getTimeRangeOptions() {
  return Object.values(TimeRange).map((value) => ({
    value,
    label: getTimeRangeLabel(value),
  }));
}

/**
 * 格式化金額（加上千分位）
 */
export function formatCurrency(
  amount: number,
  currency: string = "TWD"
): string {
  return `${currency} ${amount.toLocaleString("zh-TW", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  })}`;
}

/**
 * 格式化百分比
 */
export function formatPercentage(value: number): string {
  const sign = value >= 0 ? "+" : "";
  return `${sign}${value.toFixed(2)}%`;
}

/**
 * 判斷是否為正值
 */
export function isPositive(value: number): boolean {
  return value >= 0;
}

// ==================== 未實現損益分析 ====================

/**
 * 未實現損益摘要
 */
export interface UnrealizedSummary {
  total_cost: number;
  total_market_value: number;
  total_unrealized_pl: number;
  total_unrealized_pct: number;
  holding_count: number;
  currency: string;
}

/**
 * 未實現績效資料
 */
export interface UnrealizedPerformance {
  asset_type: AssetType;
  name: string;
  cost: number;
  market_value: number;
  unrealized_pl: number;
  unrealized_pct: number;
  holding_count: number;
}

/**
 * Top 未實現損益資產
 */
export interface UnrealizedTopAsset {
  symbol: string;
  name: string;
  asset_type: AssetType;
  quantity: number;
  avg_cost: number;
  current_price: number;
  cost: number;
  market_value: number;
  unrealized_pl: number;
  unrealized_pct: number;
}

// ==================== 資產配置分析 ====================

/**
 * 按資產類型的配置
 */
export interface AllocationByType {
  asset_type: AssetType;
  name: string;
  market_value: number;
  percentage: number;
  count: number;
}

/**
 * 按個別資產的配置
 */
export interface AllocationByAsset {
  symbol: string;
  name: string;
  asset_type: AssetType;
  market_value: number;
  percentage: number;
  quantity: number;
}

/**
 * 資產配置摘要
 */
export interface AllocationSummary {
  total_market_value: number;
  by_type: AllocationByType[];
  by_asset: AllocationByAsset[];
  currency: string;
  as_of_date: string;
}

// ==================== 績效趨勢分析 ====================

/**
 * 績效趨勢資料點
 */
export interface PerformanceTrendPoint {
  date: string;
  market_value: number;
  cost: number;
  unrealized_pl: number;
  unrealized_pct: number;
  realized_pl: number;
  realized_pct: number;
  total_pl: number;
  total_pct: number;
  holding_count: number;
}

/**
 * 按資產類型的績效趨勢
 */
export interface PerformanceTrendByType {
  asset_type: AssetType;
  name: string;
  data: PerformanceTrendPoint[];
}

/**
 * 績效趨勢摘要
 */
export interface PerformanceTrendSummary {
  start_date: string;
  end_date: string;
  total_data: PerformanceTrendPoint[];
  by_type: PerformanceTrendByType[];
  currency: string;
  data_point_count: number;
}

/**
 * 每日績效快照
 */
export interface DailyPerformanceSnapshot {
  id: string;
  snapshot_date: string;
  total_market_value: number;
  total_cost: number;
  total_unrealized_pl: number;
  total_unrealized_pct: number;
  total_realized_pl: number;
  total_realized_pct: number;
  holding_count: number;
  currency: string;
  created_at: string;
  updated_at: string;
}
