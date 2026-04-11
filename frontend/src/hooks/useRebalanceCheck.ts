import { useQuery, type UseQueryOptions } from "@tanstack/react-query";
import { rebalanceAPI } from "@/lib/api/rebalance";
import type { RebalanceCheck } from "@/types/rebalance";
import { APIError } from "@/lib/api/client";

/**
 * Query Keys
 */
export const rebalanceKeys = {
  all: ["rebalance"] as const,
  check: () => [...rebalanceKeys.all, "check"] as const,
};

/**
 * 取得再平衡檢查結果
 *
 * @param options React Query 選項
 * @returns 再平衡檢查查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useRebalanceCheck();
 * if (data?.needs_rebalance) {
 *   // 顯示再平衡建議
 * }
 * ```
 */
export function useRebalanceCheck(
  options?: Omit<
    UseQueryOptions<RebalanceCheck, APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<RebalanceCheck, APIError>({
    queryKey: rebalanceKeys.check(),
    queryFn: () => rebalanceAPI.check(),
    // 每 10 分鐘自動更新
    staleTime: 10 * 60 * 1000,
    refetchOnWindowFocus: true,
    ...options,
  });
}
