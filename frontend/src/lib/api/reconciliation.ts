import { apiClient } from "./client";
import type {
  ReconcileBankAccountInput,
  ReconcileBatchInput,
  ReconcileCreditCardInput,
  ReconcileResult,
} from "@/types/reconciliation";

const RECONCILIATION_ENDPOINTS = {
  BANK_ACCOUNT: (id: string) => `/api/bank-accounts/${id}/reconcile`,
  CREDIT_CARD: (id: string) => `/api/credit-cards/${id}/reconcile`,
  BATCH: "/api/user-management/reconcile-batch",
} as const;

/**
 * Balance reconciliation API
 */
export const reconciliationAPI = {
  /**
   * 校準單一銀行帳戶餘額
   */
  reconcileBankAccount: async (
    id: string,
    data: ReconcileBankAccountInput
  ): Promise<ReconcileResult> => {
    return apiClient.post<ReconcileResult>(
      RECONCILIATION_ENDPOINTS.BANK_ACCOUNT(id),
      data
    );
  },

  /**
   * 校準單一信用卡 used_credit / credit_limit
   */
  reconcileCreditCard: async (
    id: string,
    data: ReconcileCreditCardInput
  ): Promise<ReconcileResult> => {
    return apiClient.post<ReconcileResult>(
      RECONCILIATION_ENDPOINTS.CREDIT_CARD(id),
      data
    );
  },

  /**
   * 批次校準（在後端單一 transaction 內執行）
   */
  reconcileBatch: async (
    data: ReconcileBatchInput
  ): Promise<ReconcileResult[]> => {
    return apiClient.post<ReconcileResult[]>(
      RECONCILIATION_ENDPOINTS.BATCH,
      data
    );
  },
};
