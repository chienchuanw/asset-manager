import { useQuery, type UseQueryOptions } from "@tanstack/react-query";
import { unrealizedAnalyticsAPI } from "@/lib/api/unrealized-analytics";
import type {
  UnrealizedSummary,
  UnrealizedPerformance,
  UnrealizedTopAsset,
} from "@/types/analytics";
import { APIError } from "@/lib/api/client";

/**
 * Query Keys
 * 用於識別和管理快取
 */
export const unrealizedAnalyticsKeys = {
  all: ["unrealized-analytics"] as const,
  summary: () => [...unrealizedAnalyticsKeys.all, "summary"] as const,
  performance: () =>
    [...unrealizedAnalyticsKeys.all, "performance"] as const,
  topAssets: (limit: number) =>
    [...unrealizedAnalyticsKeys.all, "top-assets", limit] as const,
};

/**
 * 取得未實現損益摘要
 *
 * @param options React Query 選項
 * @returns 未實現損益摘要查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useUnrealizedSummary();
 * ```
 */
export function useUnrealizedSummary(
  options?: Omit<
    UseQueryOptions<UnrealizedSummary, APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<UnrealizedSummary, APIError>({
    queryKey: unrealizedAnalyticsKeys.summary(),
    queryFn: () => unrealizedAnalyticsAPI.getSummary(),
    staleTime: 1000 * 60 * 5, // 5 分鐘內不重新取得
    ...options,
  });
}

/**
 * 取得各資產類型未實現績效
 *
 * @param options React Query 選項
 * @returns 未實現績效資料查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useUnrealizedPerformance();
 * ```
 */
export function useUnrealizedPerformance(
  options?: Omit<
    UseQueryOptions<UnrealizedPerformance[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<UnrealizedPerformance[], APIError>({
    queryKey: unrealizedAnalyticsKeys.performance(),
    queryFn: () => unrealizedAnalyticsAPI.getPerformance(),
    staleTime: 1000 * 60 * 5, // 5 分鐘內不重新取得
    ...options,
  });
}

/**
 * 取得 Top 未實現損益資產
 *
 * @param limit 回傳數量限制（預設 10）
 * @param options React Query 選項
 * @returns Top 未實現損益資產查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useUnrealizedTopAssets(10);
 * ```
 */
export function useUnrealizedTopAssets(
  limit: number = 10,
  options?: Omit<
    UseQueryOptions<UnrealizedTopAsset[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<UnrealizedTopAsset[], APIError>({
    queryKey: unrealizedAnalyticsKeys.topAssets(limit),
    queryFn: () => unrealizedAnalyticsAPI.getTopAssets(limit),
    staleTime: 1000 * 60 * 5, // 5 分鐘內不重新取得
    ...options,
  });
}

/**
 * 一次取得所有未實現分析資料
 *
 * @param topAssetsLimit 最佳表現資產數量限制
 * @returns 所有未實現分析資料的查詢結果
 *
 * @example
 * ```tsx
 * const { summary, performance, topAssets, isLoading, isError } = useUnrealizedAnalytics(10);
 * ```
 */
export function useUnrealizedAnalytics(topAssetsLimit: number = 10) {
  const summary = useUnrealizedSummary();
  const performance = useUnrealizedPerformance();
  const topAssets = useUnrealizedTopAssets(topAssetsLimit);

  return {
    summary,
    performance,
    topAssets,
    isLoading:
      summary.isLoading || performance.isLoading || topAssets.isLoading,
    isError: summary.isError || performance.isError || topAssets.isError,
    error: summary.error || performance.error || topAssets.error,
  };
}

