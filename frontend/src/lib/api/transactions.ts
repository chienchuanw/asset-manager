import { apiClient } from "./client";
import type {
  Transaction,
  CreateTransactionInput,
  UpdateTransactionInput,
  TransactionFilters,
  BatchCreateTransactionsInput,
} from "@/types/transaction";

/**
 * 交易 API 端點
 */
const ENDPOINTS = {
  TRANSACTIONS: "/api/transactions",
  TRANSACTIONS_BATCH: "/api/transactions/batch",
  TRANSACTION_BY_ID: (id: string) => `/api/transactions/${id}`,
} as const;

/**
 * 交易 API
 */
export const transactionsAPI = {
  /**
   * 取得所有交易記錄
   * @param filters 篩選條件
   * @returns 交易記錄陣列
   */
  getAll: async (filters?: TransactionFilters): Promise<Transaction[]> => {
    return apiClient.get<Transaction[]>(ENDPOINTS.TRANSACTIONS, {
      params: filters,
    });
  },

  /**
   * 取得單筆交易記錄
   * @param id 交易 ID
   * @returns 交易記錄
   */
  getById: async (id: string): Promise<Transaction> => {
    return apiClient.get<Transaction>(ENDPOINTS.TRANSACTION_BY_ID(id));
  },

  /**
   * 建立交易記錄
   * @param data 交易資料
   * @returns 建立的交易記錄
   */
  create: async (data: CreateTransactionInput): Promise<Transaction> => {
    return apiClient.post<Transaction>(ENDPOINTS.TRANSACTIONS, data);
  },

  /**
   * 更新交易記錄
   * @param id 交易 ID
   * @param data 更新的交易資料
   * @returns 更新後的交易記錄
   */
  update: async (
    id: string,
    data: UpdateTransactionInput
  ): Promise<Transaction> => {
    return apiClient.put<Transaction>(ENDPOINTS.TRANSACTION_BY_ID(id), data);
  },

  /**
   * 批次建立交易記錄
   * @param data 批次交易資料
   * @returns 建立的交易記錄陣列
   */
  createBatch: async (
    data: BatchCreateTransactionsInput
  ): Promise<Transaction[]> => {
    return apiClient.post<Transaction[]>(ENDPOINTS.TRANSACTIONS_BATCH, data);
  },

  /**
   * 刪除交易記錄
   * @param id 交易 ID
   * @returns void
   */
  delete: async (id: string): Promise<void> => {
    return apiClient.delete<void>(ENDPOINTS.TRANSACTION_BY_ID(id));
  },
};
