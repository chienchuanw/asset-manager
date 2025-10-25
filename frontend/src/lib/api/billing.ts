import { apiClient } from "./client";
import type {
  BillingResult,
  DailyBillingResult,
  ProcessBillingInput,
} from "@/types/billing";

/**
 * 扣款 API 端點
 */
const ENDPOINTS = {
  PROCESS_DAILY: "/api/billing/process-daily",
  PROCESS_SUBSCRIPTIONS: "/api/billing/process-subscriptions",
  PROCESS_INSTALLMENTS: "/api/billing/process-installments",
} as const;

/**
 * 扣款 API
 */
export const billingAPI = {
  /**
   * 處理每日扣款（訂閱 + 分期）
   * @param data 處理扣款的輸入參數（可選日期）
   * @returns 每日扣款結果
   */
  processDaily: async (
    data?: ProcessBillingInput
  ): Promise<DailyBillingResult> => {
    return apiClient.post<DailyBillingResult>(ENDPOINTS.PROCESS_DAILY, data);
  },

  /**
   * 處理訂閱扣款
   * @param data 處理扣款的輸入參數（可選日期）
   * @returns 扣款結果
   */
  processSubscriptions: async (
    data?: ProcessBillingInput
  ): Promise<BillingResult> => {
    return apiClient.post<BillingResult>(
      ENDPOINTS.PROCESS_SUBSCRIPTIONS,
      data
    );
  },

  /**
   * 處理分期扣款
   * @param data 處理扣款的輸入參數（可選日期）
   * @returns 扣款結果
   */
  processInstallments: async (
    data?: ProcessBillingInput
  ): Promise<BillingResult> => {
    return apiClient.post<BillingResult>(ENDPOINTS.PROCESS_INSTALLMENTS, data);
  },
};

