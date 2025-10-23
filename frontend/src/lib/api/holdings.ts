import { apiClient } from "./client";
import type { Holding, HoldingFilters } from "@/types/holding";

/**
 * Holdings API 端點
 */
const ENDPOINTS = {
  HOLDINGS: "/api/holdings",
  HOLDING_BY_SYMBOL: (symbol: string) => `/api/holdings/${symbol}`,
} as const;

/**
 * Holdings API
 */
export const holdingsAPI = {
  /**
   * 取得所有持倉
   * @param filters 篩選條件
   * @returns 持倉陣列
   * 
   * @example
   * ```typescript
   * // 取得所有持倉
   * const holdings = await holdingsAPI.getAll();
   * 
   * // 只取得台股持倉
   * const twStocks = await holdingsAPI.getAll({ asset_type: "tw-stock" });
   * ```
   */
  getAll: async (filters?: HoldingFilters): Promise<Holding[]> => {
    return apiClient.get<Holding[]>(ENDPOINTS.HOLDINGS, {
      params: filters,
    });
  },

  /**
   * 取得單一持倉
   * @param symbol 標的代碼
   * @returns 持倉資料
   * 
   * @example
   * ```typescript
   * const holding = await holdingsAPI.getBySymbol("2330");
   * ```
   */
  getBySymbol: async (symbol: string): Promise<Holding> => {
    return apiClient.get<Holding>(ENDPOINTS.HOLDING_BY_SYMBOL(symbol));
  },
};

