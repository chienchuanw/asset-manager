import { apiClient } from "./client";
import type {
  BankAccount,
  CreateBankAccountInput,
  UpdateBankAccountInput,
  BankAccountFilters,
  CreditCard,
  CreateCreditCardInput,
  UpdateCreditCardInput,
  CreditCardQueryParams,
} from "@/types/user-management";

/**
 * 銀行帳戶 API 端點
 */
const BANK_ACCOUNT_ENDPOINTS = {
  BANK_ACCOUNTS: "/api/bank-accounts",
  BANK_ACCOUNT_BY_ID: (id: string) => `/api/bank-accounts/${id}`,
} as const;

/**
 * 信用卡 API 端點
 */
const CREDIT_CARD_ENDPOINTS = {
  CREDIT_CARDS: "/api/credit-cards",
  CREDIT_CARD_BY_ID: (id: string) => `/api/credit-cards/${id}`,
  UPCOMING_BILLING: "/api/credit-cards/upcoming-billing",
  UPCOMING_PAYMENT: "/api/credit-cards/upcoming-payment",
} as const;

/**
 * 銀行帳戶 API
 */
export const bankAccountsAPI = {
  /**
   * 取得所有銀行帳戶
   * @param filters 篩選條件
   * @returns 銀行帳戶陣列
   */
  getAll: async (filters?: BankAccountFilters): Promise<BankAccount[]> => {
    return apiClient.get<BankAccount[]>(
      BANK_ACCOUNT_ENDPOINTS.BANK_ACCOUNTS,
      {
        params: filters,
      }
    );
  },

  /**
   * 取得單筆銀行帳戶
   * @param id 銀行帳戶 ID
   * @returns 銀行帳戶
   */
  getById: async (id: string): Promise<BankAccount> => {
    return apiClient.get<BankAccount>(
      BANK_ACCOUNT_ENDPOINTS.BANK_ACCOUNT_BY_ID(id)
    );
  },

  /**
   * 建立銀行帳戶
   * @param data 銀行帳戶資料
   * @returns 建立的銀行帳戶
   */
  create: async (data: CreateBankAccountInput): Promise<BankAccount> => {
    return apiClient.post<BankAccount>(
      BANK_ACCOUNT_ENDPOINTS.BANK_ACCOUNTS,
      data
    );
  },

  /**
   * 更新銀行帳戶
   * @param id 銀行帳戶 ID
   * @param data 更新的銀行帳戶資料
   * @returns 更新後的銀行帳戶
   */
  update: async (
    id: string,
    data: UpdateBankAccountInput
  ): Promise<BankAccount> => {
    return apiClient.put<BankAccount>(
      BANK_ACCOUNT_ENDPOINTS.BANK_ACCOUNT_BY_ID(id),
      data
    );
  },

  /**
   * 刪除銀行帳戶
   * @param id 銀行帳戶 ID
   * @returns void
   */
  delete: async (id: string): Promise<void> => {
    return apiClient.delete<void>(
      BANK_ACCOUNT_ENDPOINTS.BANK_ACCOUNT_BY_ID(id)
    );
  },
};

/**
 * 信用卡 API
 */
export const creditCardsAPI = {
  /**
   * 取得所有信用卡
   * @returns 信用卡陣列
   */
  getAll: async (): Promise<CreditCard[]> => {
    return apiClient.get<CreditCard[]>(CREDIT_CARD_ENDPOINTS.CREDIT_CARDS);
  },

  /**
   * 取得單筆信用卡
   * @param id 信用卡 ID
   * @returns 信用卡
   */
  getById: async (id: string): Promise<CreditCard> => {
    return apiClient.get<CreditCard>(
      CREDIT_CARD_ENDPOINTS.CREDIT_CARD_BY_ID(id)
    );
  },

  /**
   * 取得即將到來的帳單日信用卡
   * @param params 查詢參數
   * @returns 信用卡陣列
   */
  getUpcomingBilling: async (
    params?: CreditCardQueryParams
  ): Promise<CreditCard[]> => {
    return apiClient.get<CreditCard[]>(
      CREDIT_CARD_ENDPOINTS.UPCOMING_BILLING,
      {
        params,
      }
    );
  },

  /**
   * 取得即將到來的繳款截止日信用卡
   * @param params 查詢參數
   * @returns 信用卡陣列
   */
  getUpcomingPayment: async (
    params?: CreditCardQueryParams
  ): Promise<CreditCard[]> => {
    return apiClient.get<CreditCard[]>(
      CREDIT_CARD_ENDPOINTS.UPCOMING_PAYMENT,
      {
        params,
      }
    );
  },

  /**
   * 建立信用卡
   * @param data 信用卡資料
   * @returns 建立的信用卡
   */
  create: async (data: CreateCreditCardInput): Promise<CreditCard> => {
    return apiClient.post<CreditCard>(
      CREDIT_CARD_ENDPOINTS.CREDIT_CARDS,
      data
    );
  },

  /**
   * 更新信用卡
   * @param id 信用卡 ID
   * @param data 更新的信用卡資料
   * @returns 更新後的信用卡
   */
  update: async (
    id: string,
    data: UpdateCreditCardInput
  ): Promise<CreditCard> => {
    return apiClient.put<CreditCard>(
      CREDIT_CARD_ENDPOINTS.CREDIT_CARD_BY_ID(id),
      data
    );
  },

  /**
   * 刪除信用卡
   * @param id 信用卡 ID
   * @returns void
   */
  delete: async (id: string): Promise<void> => {
    return apiClient.delete<void>(
      CREDIT_CARD_ENDPOINTS.CREDIT_CARD_BY_ID(id)
    );
  },
};

