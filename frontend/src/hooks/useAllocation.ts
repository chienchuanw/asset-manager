import { useQuery, type UseQueryOptions } from "@tanstack/react-query";
import { allocationAPI } from "@/lib/api/allocation";
import type {
  AllocationSummary,
  AllocationByType,
  AllocationByAsset,
} from "@/types/analytics";
import { APIError } from "@/lib/api/client";

/**
 * Query Keys
 * 用於識別和管理快取
 */
export const allocationKeys = {
  all: ["allocation"] as const,
  current: () => [...allocationKeys.all, "current"] as const,
  byType: () => [...allocationKeys.all, "by-type"] as const,
  byAsset: (limit: number) =>
    [...allocationKeys.all, "by-asset", limit] as const,
};

/**
 * 取得當前資產配置摘要
 *
 * @param options React Query 選項
 * @returns 資產配置摘要查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useCurrentAllocation();
 * ```
 */
export function useCurrentAllocation(
  options?: Omit<
    UseQueryOptions<AllocationSummary, APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<AllocationSummary, APIError>({
    queryKey: allocationKeys.current(),
    queryFn: () => allocationAPI.getCurrentAllocation(),
    staleTime: 1000 * 60 * 5, // 5 分鐘內不重新取得
    ...options,
  });
}

/**
 * 取得按資產類型的配置
 *
 * @param options React Query 選項
 * @returns 按資產類型分類的配置查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useAllocationByType();
 * ```
 */
export function useAllocationByType(
  options?: Omit<
    UseQueryOptions<AllocationByType[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<AllocationByType[], APIError>({
    queryKey: allocationKeys.byType(),
    queryFn: () => allocationAPI.getAllocationByType(),
    staleTime: 1000 * 60 * 5, // 5 分鐘內不重新取得
    ...options,
  });
}

/**
 * 取得按個別資產的配置
 *
 * @param limit 回傳數量限制（預設 20）
 * @param options React Query 選項
 * @returns 按個別資產分類的配置查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useAllocationByAsset(10);
 * ```
 */
export function useAllocationByAsset(
  limit: number = 20,
  options?: Omit<
    UseQueryOptions<AllocationByAsset[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<AllocationByAsset[], APIError>({
    queryKey: allocationKeys.byAsset(limit),
    queryFn: () => allocationAPI.getAllocationByAsset(limit),
    staleTime: 1000 * 60 * 5, // 5 分鐘內不重新取得
    ...options,
  });
}

/**
 * 一次取得所有資產配置資料
 *
 * @param assetLimit 個別資產數量限制
 * @returns 所有資產配置資料的查詢結果
 *
 * @example
 * ```tsx
 * const { current, byType, byAsset, isLoading, isError } = useAllocation(10);
 * ```
 */
export function useAllocation(assetLimit: number = 20) {
  const current = useCurrentAllocation();
  const byType = useAllocationByType();
  const byAsset = useAllocationByAsset(assetLimit);

  return {
    current,
    byType,
    byAsset,
    isLoading: current.isLoading || byType.isLoading || byAsset.isLoading,
    isError: current.isError || byType.isError || byAsset.isError,
    error: current.error || byType.error || byAsset.error,
  };
}

