import { apiClient } from "./client";
import type { FireProjectionResult, FireProjectionParams } from "@/types/fire";

/**
 * FIRE API
 * 提供 FIRE 投影相關的 API 呼叫
 */
export const fireAPI = {
  /**
   * 取得 FIRE 投影
   *
   * 未提供的參數由後端從應用資料自動推導（最新資產快照淨值、
   * 近 12 個月現金流推導的年支出與年儲蓄）。
   *
   * @param params 可選的覆寫參數
   * @returns FIRE 投影結果
   */
  getProjection: async (
    params: FireProjectionParams = {}
  ): Promise<FireProjectionResult> => {
    return apiClient.get<FireProjectionResult>("/api/fire/projection", {
      params: params as Record<string, number | undefined>,
    });
  },
};
