import {
  useQuery,
  type UseQueryOptions,
} from "@tanstack/react-query";
import { analyticsAPI } from "@/lib/api/analytics";
import type {
  AnalyticsSummary,
  PerformanceData,
  TopAsset,
  TimeRange,
} from "@/types/analytics";
import { APIError } from "@/lib/api/client";

/**
 * Query Keys
 * 用於識別和管理快取
 */
export const analyticsKeys = {
  all: ["analytics"] as const,
  summary: (timeRange: TimeRange) => [...analyticsKeys.all, "summary", timeRange] as const,
  performance: (timeRange: TimeRange) => [...analyticsKeys.all, "performance", timeRange] as const,
  topAssets: (timeRange: TimeRange, limit: number) =>
    [...analyticsKeys.all, "top-assets", timeRange, limit] as const,
};

/**
 * 取得分析摘要
 *
 * @param timeRange 時間範圍
 * @param options React Query 選項
 * @returns 分析摘要查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useAnalyticsSummary("month");
 * ```
 */
export function useAnalyticsSummary(
  timeRange: TimeRange = "month",
  options?: Omit<
    UseQueryOptions<AnalyticsSummary, APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<AnalyticsSummary, APIError>({
    queryKey: analyticsKeys.summary(timeRange),
    queryFn: () => analyticsAPI.getSummary(timeRange),
    staleTime: 1000 * 60 * 5, // 5 分鐘內不重新取得
    ...options,
  });
}

/**
 * 取得各資產類型績效
 *
 * @param timeRange 時間範圍
 * @param options React Query 選項
 * @returns 績效資料查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useAnalyticsPerformance("month");
 * ```
 */
export function useAnalyticsPerformance(
  timeRange: TimeRange = "month",
  options?: Omit<
    UseQueryOptions<PerformanceData[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<PerformanceData[], APIError>({
    queryKey: analyticsKeys.performance(timeRange),
    queryFn: () => analyticsAPI.getPerformance(timeRange),
    staleTime: 1000 * 60 * 5, // 5 分鐘內不重新取得
    ...options,
  });
}

/**
 * 取得最佳/最差表現資產
 *
 * @param timeRange 時間範圍
 * @param limit 回傳數量限制
 * @param options React Query 選項
 * @returns 最佳表現資產查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useAnalyticsTopAssets("month", 10);
 * ```
 */
export function useAnalyticsTopAssets(
  timeRange: TimeRange = "month",
  limit: number = 5,
  options?: Omit<
    UseQueryOptions<TopAsset[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<TopAsset[], APIError>({
    queryKey: analyticsKeys.topAssets(timeRange, limit),
    queryFn: () => analyticsAPI.getTopAssets(timeRange, limit),
    staleTime: 1000 * 60 * 5, // 5 分鐘內不重新取得
    ...options,
  });
}

/**
 * 一次取得所有分析資料
 *
 * @param timeRange 時間範圍
 * @param topAssetsLimit 最佳表現資產數量限制
 * @returns 所有分析資料的查詢結果
 *
 * @example
 * ```tsx
 * const { summary, performance, topAssets } = useAnalytics("month");
 * ```
 */
export function useAnalytics(
  timeRange: TimeRange = "month",
  topAssetsLimit: number = 5
) {
  const summary = useAnalyticsSummary(timeRange);
  const performance = useAnalyticsPerformance(timeRange);
  const topAssets = useAnalyticsTopAssets(timeRange, topAssetsLimit);

  return {
    summary,
    performance,
    topAssets,
    isLoading: summary.isLoading || performance.isLoading || topAssets.isLoading,
    isError: summary.isError || performance.isError || topAssets.isError,
    error: summary.error || performance.error || topAssets.error,
  };
}

