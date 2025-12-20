import { apiClient } from "./client";
import type {
  CashFlow,
  CashFlowCategory,
  CashFlowSummary,
  CreateCashFlowInput,
  UpdateCashFlowInput,
  CreateCategoryInput,
  UpdateCategoryInput,
  ReorderCategoryInput,
  CashFlowFilters,
  CashFlowType,
} from "@/types/cash-flow";

/**
 * 現金流 API 端點
 */
const ENDPOINTS = {
  CASH_FLOWS: "/api/cash-flows",
  CASH_FLOW_BY_ID: (id: string) => `/api/cash-flows/${id}`,
  CASH_FLOW_SUMMARY: "/api/cash-flows/summary",
  CATEGORIES: "/api/categories",
  CATEGORY_BY_ID: (id: string) => `/api/categories/${id}`,
  CATEGORIES_REORDER: "/api/categories/reorder",
} as const;

/**
 * 現金流 API
 */
export const cashFlowsAPI = {
  /**
   * 取得所有現金流記錄
   * @param filters 篩選條件
   * @returns 現金流記錄陣列
   */
  getAll: async (filters?: CashFlowFilters): Promise<CashFlow[]> => {
    return apiClient.get<CashFlow[]>(ENDPOINTS.CASH_FLOWS, {
      params: filters,
    });
  },

  /**
   * 取得單筆現金流記錄
   * @param id 現金流記錄 ID
   * @returns 現金流記錄
   */
  getById: async (id: string): Promise<CashFlow> => {
    return apiClient.get<CashFlow>(ENDPOINTS.CASH_FLOW_BY_ID(id));
  },

  /**
   * 建立現金流記錄
   * @param data 現金流記錄資料
   * @returns 建立的現金流記錄
   */
  create: async (data: CreateCashFlowInput): Promise<CashFlow> => {
    return apiClient.post<CashFlow>(ENDPOINTS.CASH_FLOWS, data);
  },

  /**
   * 更新現金流記錄
   * @param id 現金流記錄 ID
   * @param data 更新的現金流記錄資料
   * @returns 更新後的現金流記錄
   */
  update: async (id: string, data: UpdateCashFlowInput): Promise<CashFlow> => {
    return apiClient.put<CashFlow>(ENDPOINTS.CASH_FLOW_BY_ID(id), data);
  },

  /**
   * 刪除現金流記錄
   * @param id 現金流記錄 ID
   * @returns void
   */
  delete: async (id: string): Promise<void> => {
    return apiClient.delete<void>(ENDPOINTS.CASH_FLOW_BY_ID(id));
  },

  /**
   * 取得現金流摘要
   * @param startDate 開始日期 (YYYY-MM-DD)
   * @param endDate 結束日期 (YYYY-MM-DD)
   * @returns 現金流摘要
   */
  getSummary: async (
    startDate: string,
    endDate: string
  ): Promise<CashFlowSummary> => {
    return apiClient.get<CashFlowSummary>(ENDPOINTS.CASH_FLOW_SUMMARY, {
      params: {
        start_date: startDate,
        end_date: endDate,
      },
    });
  },
};

/**
 * 分類 API
 */
export const categoriesAPI = {
  /**
   * 取得所有分類
   * @param type 現金流類型（可選）
   * @returns 分類陣列
   */
  getAll: async (type?: CashFlowType): Promise<CashFlowCategory[]> => {
    return apiClient.get<CashFlowCategory[]>(ENDPOINTS.CATEGORIES, {
      params: type ? { type } : undefined,
    });
  },

  /**
   * 取得單筆分類
   * @param id 分類 ID
   * @returns 分類
   */
  getById: async (id: string): Promise<CashFlowCategory> => {
    return apiClient.get<CashFlowCategory>(ENDPOINTS.CATEGORY_BY_ID(id));
  },

  /**
   * 建立分類
   * @param data 分類資料
   * @returns 建立的分類
   */
  create: async (data: CreateCategoryInput): Promise<CashFlowCategory> => {
    return apiClient.post<CashFlowCategory>(ENDPOINTS.CATEGORIES, data);
  },

  /**
   * 更新分類
   * @param id 分類 ID
   * @param data 更新的分類資料
   * @returns 更新後的分類
   */
  update: async (
    id: string,
    data: UpdateCategoryInput
  ): Promise<CashFlowCategory> => {
    return apiClient.put<CashFlowCategory>(ENDPOINTS.CATEGORY_BY_ID(id), data);
  },

  /**
   * 刪除分類
   * @param id 分類 ID
   * @returns void
   */
  delete: async (id: string): Promise<void> => {
    return apiClient.delete<void>(ENDPOINTS.CATEGORY_BY_ID(id));
  },

  /**
   * 重新排序分類
   * @param data 排序資料
   * @returns void
   */
  reorder: async (data: ReorderCategoryInput): Promise<void> => {
    return apiClient.put<void>(ENDPOINTS.CATEGORIES_REORDER, data);
  },
};
