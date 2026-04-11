import { apiClient } from "./client";
import type { RebalanceCheck } from "@/types/rebalance";

/**
 * Rebalance API 端點
 */
const ENDPOINTS = {
  CHECK: "/api/rebalance/check",
} as const;

/**
 * Rebalance API
 */
export const rebalanceAPI = {
  /**
   * 檢查是否需要再平衡
   * @returns 再平衡檢查結果
   */
  check: async (): Promise<RebalanceCheck> => {
    return apiClient.get<RebalanceCheck>(ENDPOINTS.CHECK);
  },
};
