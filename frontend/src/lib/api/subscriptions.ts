import { apiClient } from "./client";
import type {
  Subscription,
  CreateSubscriptionInput,
  UpdateSubscriptionInput,
  CancelSubscriptionInput,
  SubscriptionFilters,
} from "@/types/subscription";

/**
 * 訂閱 API 端點
 */
const ENDPOINTS = {
  SUBSCRIPTIONS: "/api/subscriptions",
  SUBSCRIPTION_BY_ID: (id: string) => `/api/subscriptions/${id}`,
  CANCEL_SUBSCRIPTION: (id: string) => `/api/subscriptions/${id}/cancel`,
} as const;

/**
 * 訂閱 API
 */
export const subscriptionsAPI = {
  /**
   * 取得所有訂閱
   * @param filters 篩選條件
   * @returns 訂閱陣列
   */
  getAll: async (filters?: SubscriptionFilters): Promise<Subscription[]> => {
    return apiClient.get<Subscription[]>(ENDPOINTS.SUBSCRIPTIONS, {
      params: filters,
    });
  },

  /**
   * 取得單筆訂閱
   * @param id 訂閱 ID
   * @returns 訂閱
   */
  getById: async (id: string): Promise<Subscription> => {
    return apiClient.get<Subscription>(ENDPOINTS.SUBSCRIPTION_BY_ID(id));
  },

  /**
   * 建立訂閱
   * @param data 訂閱資料
   * @returns 建立的訂閱
   */
  create: async (data: CreateSubscriptionInput): Promise<Subscription> => {
    return apiClient.post<Subscription>(ENDPOINTS.SUBSCRIPTIONS, data);
  },

  /**
   * 更新訂閱
   * @param id 訂閱 ID
   * @param data 更新的訂閱資料
   * @returns 更新後的訂閱
   */
  update: async (
    id: string,
    data: UpdateSubscriptionInput
  ): Promise<Subscription> => {
    return apiClient.put<Subscription>(ENDPOINTS.SUBSCRIPTION_BY_ID(id), data);
  },

  /**
   * 刪除訂閱
   * @param id 訂閱 ID
   * @returns void
   */
  delete: async (id: string): Promise<void> => {
    return apiClient.delete<void>(ENDPOINTS.SUBSCRIPTION_BY_ID(id));
  },

  /**
   * 取消訂閱
   * @param id 訂閱 ID
   * @param data 取消訂閱資料（結束日期）
   * @returns 更新後的訂閱
   */
  cancel: async (
    id: string,
    data: CancelSubscriptionInput
  ): Promise<Subscription> => {
    return apiClient.post<Subscription>(
      ENDPOINTS.CANCEL_SUBSCRIPTION(id),
      data
    );
  },
};

