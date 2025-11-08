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
  CreditCardGroup,
  CreditCardGroupWithCards,
  CreateCreditCardGroupInput,
  UpdateCreditCardGroupInput,
  AddCardsToGroupInput,
  RemoveCardsFromGroupInput,
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
 * 信用卡群組 API 端點
 */
const CREDIT_CARD_GROUP_ENDPOINTS = {
  CREDIT_CARD_GROUPS: "/api/credit-card-groups",
  CREDIT_CARD_GROUP_BY_ID: (id: string) => `/api/credit-card-groups/${id}`,
  ADD_CARDS: (id: string) => `/api/credit-card-groups/${id}/cards`,
  REMOVE_CARDS: (id: string) => `/api/credit-card-groups/${id}/cards`,
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
    return apiClient.get<BankAccount[]>(BANK_ACCOUNT_ENDPOINTS.BANK_ACCOUNTS, {
      params: filters,
    });
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
    return apiClient.get<CreditCard[]>(CREDIT_CARD_ENDPOINTS.UPCOMING_BILLING, {
      params,
    });
  },

  /**
   * 取得即將到來的繳款截止日信用卡
   * @param params 查詢參數
   * @returns 信用卡陣列
   */
  getUpcomingPayment: async (
    params?: CreditCardQueryParams
  ): Promise<CreditCard[]> => {
    return apiClient.get<CreditCard[]>(CREDIT_CARD_ENDPOINTS.UPCOMING_PAYMENT, {
      params,
    });
  },

  /**
   * 建立信用卡
   * @param data 信用卡資料
   * @returns 建立的信用卡
   */
  create: async (data: CreateCreditCardInput): Promise<CreditCard> => {
    return apiClient.post<CreditCard>(CREDIT_CARD_ENDPOINTS.CREDIT_CARDS, data);
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
    return apiClient.delete<void>(CREDIT_CARD_ENDPOINTS.CREDIT_CARD_BY_ID(id));
  },
};
/**
 * 信用卡群組 API
 */
export const creditCardGroupsAPI = {
  /**
   * 取得所有信用卡群組
   * @returns 信用卡群組陣列（包含卡片列表）
   */
  getAll: async (): Promise<CreditCardGroupWithCards[]> => {
    return apiClient.get<CreditCardGroupWithCards[]>(
      CREDIT_CARD_GROUP_ENDPOINTS.CREDIT_CARD_GROUPS
    );
  },

  /**
   * 取得單筆信用卡群組
   * @param id 信用卡群組 ID
   * @returns 信用卡群組（包含卡片列表）
   */
  getById: async (id: string): Promise<CreditCardGroupWithCards> => {
    return apiClient.get<CreditCardGroupWithCards>(
      CREDIT_CARD_GROUP_ENDPOINTS.CREDIT_CARD_GROUP_BY_ID(id)
    );
  },

  /**
   * 建立信用卡群組
   * @param data 信用卡群組資料
   * @returns 建立的信用卡群組（包含卡片列表）
   */
  create: async (
    data: CreateCreditCardGroupInput
  ): Promise<CreditCardGroupWithCards> => {
    return apiClient.post<CreditCardGroupWithCards>(
      CREDIT_CARD_GROUP_ENDPOINTS.CREDIT_CARD_GROUPS,
      data
    );
  },

  /**
   * 更新信用卡群組
   * @param id 信用卡群組 ID
   * @param data 更新的信用卡群組資料
   * @returns 更新後的信用卡群組
   */
  update: async (
    id: string,
    data: UpdateCreditCardGroupInput
  ): Promise<CreditCardGroup> => {
    return apiClient.put<CreditCardGroup>(
      CREDIT_CARD_GROUP_ENDPOINTS.CREDIT_CARD_GROUP_BY_ID(id),
      data
    );
  },

  /**
   * 刪除信用卡群組
   * @param id 信用卡群組 ID
   * @returns void
   */
  delete: async (id: string): Promise<void> => {
    return apiClient.delete<void>(
      CREDIT_CARD_GROUP_ENDPOINTS.CREDIT_CARD_GROUP_BY_ID(id)
    );
  },

  /**
   * 新增卡片到群組
   * @param id 信用卡群組 ID
   * @param data 要新增的卡片 ID 列表
   * @returns void
   */
  addCards: async (id: string, data: AddCardsToGroupInput): Promise<void> => {
    return apiClient.post<void>(
      CREDIT_CARD_GROUP_ENDPOINTS.ADD_CARDS(id),
      data
    );
  },

  /**
   * 從群組移除卡片
   * @param id 信用卡群組 ID
   * @param data 要移除的卡片 ID 列表
   * @returns void
   */
  removeCards: async (
    id: string,
    data: RemoveCardsFromGroupInput
  ): Promise<void> => {
    return apiClient.delete<void>(
      CREDIT_CARD_GROUP_ENDPOINTS.REMOVE_CARDS(id),
      data
    );
  },
};
