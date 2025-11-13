import { useQuery, type UseQueryOptions } from "@tanstack/react-query";
import { holdingsAPI } from "@/lib/api/holdings";
import type { Holding, HoldingFilters } from "@/types/holding";
import { APIError, type APIResponseWithWarnings } from "@/lib/api/client";
import type { APIWarning } from "@/types/transaction";

/**
 * Query Keys
 * 用於識別和管理快取
 */
export const holdingKeys = {
  all: ["holdings"] as const,
  lists: () => [...holdingKeys.all, "list"] as const,
  list: (filters?: HoldingFilters) =>
    [...holdingKeys.lists(), filters] as const,
  details: () => [...holdingKeys.all, "detail"] as const,
  detail: (symbol: string) => [...holdingKeys.details(), symbol] as const,
};

/**
 * 取得持倉列表（包含 warnings）
 *
 * @param filters 篩選條件
 * @param options React Query 選項
 * @returns 持倉列表查詢結果（包含 warnings）
 *
 * @example
 * ```tsx
 * // 取得所有持倉
 * const { data, isLoading, error } = useHoldings();
 * // data.data 是持倉陣列
 * // data.warnings 是警告陣列
 *
 * // 只取得台股持倉
 * const { data } = useHoldings({ asset_type: "tw-stock" });
 *
 * // 自訂 refetch 間隔（每 30 秒更新一次價格）
 * const { data } = useHoldings(undefined, {
 *   refetchInterval: 30000,
 * });
 * ```
 */
export function useHoldings(
  filters?: HoldingFilters,
  options?: Omit<
    UseQueryOptions<APIResponseWithWarnings<Holding[]>, APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<APIResponseWithWarnings<Holding[]>, APIError>({
    queryKey: holdingKeys.list(filters),
    queryFn: () => holdingsAPI.getAll(filters),
    // 預設每 5 分鐘自動更新一次（配合後端 Redis 快取）
    refetchInterval: 5 * 60 * 1000,
    // 視窗重新獲得焦點時自動更新
    refetchOnWindowFocus: true,
    // 快取時間 10 分鐘
    staleTime: 10 * 60 * 1000,
    ...options,
  });
}

/**
 * 取得單一持倉
 *
 * @param symbol 標的代碼
 * @param options React Query 選項
 * @returns 單一持倉查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useHolding("2330");
 * ```
 */
export function useHolding(
  symbol: string,
  options?: Omit<UseQueryOptions<Holding, APIError>, "queryKey" | "queryFn">
) {
  return useQuery<Holding, APIError>({
    queryKey: holdingKeys.detail(symbol),
    queryFn: () => holdingsAPI.getBySymbol(symbol),
    enabled: !!symbol, // 只有當 symbol 存在時才執行查詢
    // 預設每 5 分鐘自動更新一次
    refetchInterval: 5 * 60 * 1000,
    // 視窗重新獲得焦點時自動更新
    refetchOnWindowFocus: true,
    // 快取時間 10 分鐘
    staleTime: 10 * 60 * 1000,
    ...options,
  });
}

/**
 * 取得台股持倉
 *
 * @param options React Query 選項
 * @returns 台股持倉列表
 *
 * @example
 * ```tsx
 * const { data: twStocks } = useTWStockHoldings();
 * ```
 */
export function useTWStockHoldings(
  options?: Omit<UseQueryOptions<Holding[], APIError>, "queryKey" | "queryFn">
) {
  return useHoldings({ asset_type: "tw-stock" }, options);
}

/**
 * 取得美股持倉
 *
 * @param options React Query 選項
 * @returns 美股持倉列表
 *
 * @example
 * ```tsx
 * const { data: usStocks } = useUSStockHoldings();
 * ```
 */
export function useUSStockHoldings(
  options?: Omit<UseQueryOptions<Holding[], APIError>, "queryKey" | "queryFn">
) {
  return useHoldings({ asset_type: "us-stock" }, options);
}

/**
 * 取得加密貨幣持倉
 *
 * @param options React Query 選項
 * @returns 加密貨幣持倉列表
 *
 * @example
 * ```tsx
 * const { data: cryptos } = useCryptoHoldings();
 * ```
 */
export function useCryptoHoldings(
  options?: Omit<UseQueryOptions<Holding[], APIError>, "queryKey" | "queryFn">
) {
  return useHoldings({ asset_type: "crypto" }, options);
}
