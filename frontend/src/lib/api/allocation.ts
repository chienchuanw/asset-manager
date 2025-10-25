import { apiClient } from "./client";
import type {
  AllocationSummary,
  AllocationByType,
  AllocationByAsset,
} from "@/types/analytics";

/**
 * Allocation API
 * 提供資產配置相關的 API 呼叫
 */
export const allocationAPI = {
  /**
   * 取得當前資產配置摘要
   *
   * @returns 資產配置摘要資料
   *
   * @example
   * ```ts
   * const summary = await allocationAPI.getCurrentAllocation();
   * ```
   */
  getCurrentAllocation: async (): Promise<AllocationSummary> => {
    return apiClient.get<AllocationSummary>("/api/allocation/current");
  },

  /**
   * 取得按資產類型的配置
   *
   * @returns 按資產類型分類的配置資料陣列
   *
   * @example
   * ```ts
   * const byType = await allocationAPI.getAllocationByType();
   * ```
   */
  getAllocationByType: async (): Promise<AllocationByType[]> => {
    return apiClient.get<AllocationByType[]>("/api/allocation/by-type");
  },

  /**
   * 取得按個別資產的配置
   *
   * @param limit 回傳數量限制（預設 20）
   * @returns 按個別資產分類的配置資料陣列
   *
   * @example
   * ```ts
   * const byAsset = await allocationAPI.getAllocationByAsset(10);
   * ```
   */
  getAllocationByAsset: async (
    limit: number = 20
  ): Promise<AllocationByAsset[]> => {
    return apiClient.get<AllocationByAsset[]>("/api/allocation/by-asset", {
      params: { limit },
    });
  },
};

