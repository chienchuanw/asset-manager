import { apiClient } from "./client";
import type {
  UnrealizedSummary,
  UnrealizedPerformance,
  UnrealizedTopAsset,
} from "@/types/analytics";

/**
 * Unrealized Analytics API
 * 提供未實現損益分析相關的 API 呼叫
 */
export const unrealizedAnalyticsAPI = {
  /**
   * 取得未實現損益摘要
   *
   * @returns 未實現損益摘要資料
   *
   * @example
   * ```ts
   * const summary = await unrealizedAnalyticsAPI.getSummary();
   * ```
   */
  getSummary: async (): Promise<UnrealizedSummary> => {
    return apiClient.get<UnrealizedSummary>(
      "/api/analytics/unrealized/summary"
    );
  },

  /**
   * 取得各資產類型未實現績效
   *
   * @returns 未實現績效資料陣列
   *
   * @example
   * ```ts
   * const performance = await unrealizedAnalyticsAPI.getPerformance();
   * ```
   */
  getPerformance: async (): Promise<UnrealizedPerformance[]> => {
    return apiClient.get<UnrealizedPerformance[]>(
      "/api/analytics/unrealized/performance"
    );
  },

  /**
   * 取得 Top 未實現損益資產
   *
   * @param limit 回傳數量限制（預設 10）
   * @returns Top 未實現損益資產陣列
   *
   * @example
   * ```ts
   * const topAssets = await unrealizedAnalyticsAPI.getTopAssets(10);
   * ```
   */
  getTopAssets: async (limit: number = 10): Promise<UnrealizedTopAsset[]> => {
    return apiClient.get<UnrealizedTopAsset[]>(
      "/api/analytics/unrealized/top-assets",
      {
        params: { limit },
      }
    );
  },
};

