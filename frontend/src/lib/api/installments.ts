import { apiClient } from "./client";
import type {
  Installment,
  CreateInstallmentInput,
  UpdateInstallmentInput,
  InstallmentFilters,
  CompletingSoonParams,
} from "@/types/installment";

/**
 * 分期 API 端點
 */
const ENDPOINTS = {
  INSTALLMENTS: "/api/installments",
  INSTALLMENT_BY_ID: (id: string) => `/api/installments/${id}`,
  COMPLETING_SOON: "/api/installments/completing-soon",
} as const;

/**
 * 分期 API
 */
export const installmentsAPI = {
  /**
   * 取得所有分期
   * @param filters 篩選條件
   * @returns 分期陣列
   */
  getAll: async (filters?: InstallmentFilters): Promise<Installment[]> => {
    return apiClient.get<Installment[]>(ENDPOINTS.INSTALLMENTS, {
      params: filters,
    });
  },

  /**
   * 取得單筆分期
   * @param id 分期 ID
   * @returns 分期
   */
  getById: async (id: string): Promise<Installment> => {
    return apiClient.get<Installment>(ENDPOINTS.INSTALLMENT_BY_ID(id));
  },

  /**
   * 建立分期
   * @param data 分期資料
   * @returns 建立的分期
   */
  create: async (data: CreateInstallmentInput): Promise<Installment> => {
    return apiClient.post<Installment>(ENDPOINTS.INSTALLMENTS, data);
  },

  /**
   * 更新分期
   * @param id 分期 ID
   * @param data 更新的分期資料
   * @returns 更新後的分期
   */
  update: async (
    id: string,
    data: UpdateInstallmentInput
  ): Promise<Installment> => {
    return apiClient.put<Installment>(ENDPOINTS.INSTALLMENT_BY_ID(id), data);
  },

  /**
   * 刪除分期
   * @param id 分期 ID
   * @returns void
   */
  delete: async (id: string): Promise<void> => {
    return apiClient.delete<void>(ENDPOINTS.INSTALLMENT_BY_ID(id));
  },

  /**
   * 取得即將完成的分期
   * @param params 查詢參數
   * @returns 即將完成的分期陣列
   */
  getCompletingSoon: async (
    params?: CompletingSoonParams
  ): Promise<Installment[]> => {
    return apiClient.get<Installment[]>(ENDPOINTS.COMPLETING_SOON, {
      params,
    });
  },
};

