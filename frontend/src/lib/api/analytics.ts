import { apiClient } from "./client";
import type {
  AnalyticsSummary,
  PerformanceData,
  TopAsset,
  TimeRange,
} from "@/types/analytics";

/**
 * Analytics API
 * 提供分析報表相關的 API 呼叫
 */
export const analyticsAPI = {
  /**
   * 取得分析摘要
   *
   * @param timeRange 時間範圍
   * @returns 分析摘要資料
   *
   * @example
   * ```ts
   * const summary = await analyticsAPI.getSummary("month");
   * ```
   */
  getSummary: async (timeRange: TimeRange = "month"): Promise<AnalyticsSummary> => {
    return apiClient.get<AnalyticsSummary>("/api/analytics/summary", {
      params: { time_range: timeRange },
    });
  },

  /**
   * 取得各資產類型績效
   *
   * @param timeRange 時間範圍
   * @returns 績效資料陣列
   *
   * @example
   * ```ts
   * const performance = await analyticsAPI.getPerformance("month");
   * ```
   */
  getPerformance: async (timeRange: TimeRange = "month"): Promise<PerformanceData[]> => {
    return apiClient.get<PerformanceData[]>("/api/analytics/performance", {
      params: { time_range: timeRange },
    });
  },

  /**
   * 取得最佳/最差表現資產
   *
   * @param timeRange 時間範圍
   * @param limit 回傳數量限制（預設 5）
   * @returns 最佳表現資產陣列
   *
   * @example
   * ```ts
   * const topAssets = await analyticsAPI.getTopAssets("month", 10);
   * ```
   */
  getTopAssets: async (
    timeRange: TimeRange = "month",
    limit: number = 5
  ): Promise<TopAsset[]> => {
    return apiClient.get<TopAsset[]>("/api/analytics/top-assets", {
      params: {
        time_range: timeRange,
        limit,
      },
    });
  },
};

