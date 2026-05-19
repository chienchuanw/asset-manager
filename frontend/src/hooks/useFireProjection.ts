import { useQuery, type UseQueryOptions } from "@tanstack/react-query";
import { fireAPI } from "@/lib/api/fire";
import type { FireProjectionResult, FireProjectionParams } from "@/types/fire";
import { APIError } from "@/lib/api/client";

/**
 * Query Keys
 */
export const fireKeys = {
  all: ["fire"] as const,
  projection: (params: FireProjectionParams) =>
    [...fireKeys.all, "projection", params] as const,
};

/**
 * 取得 FIRE 投影
 *
 * 省略的覆寫參數由後端自動推導。覆寫值變動時 query key 隨之改變並重新查詢。
 *
 * @param params 可選的覆寫參數
 * @param options React Query 選項
 */
export function useFireProjection(
  params: FireProjectionParams = {},
  options?: Omit<
    UseQueryOptions<FireProjectionResult, APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<FireProjectionResult, APIError>({
    queryKey: fireKeys.projection(params),
    queryFn: () => fireAPI.getProjection(params),
    staleTime: 1000 * 60 * 5,
    ...options,
  });
}
