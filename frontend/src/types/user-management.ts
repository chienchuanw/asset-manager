// 使用者管理相關型別定義

/**
 * 銀行帳戶資料結構
 */
export interface BankAccount {
  id: string;
  bank_name: string;
  account_type: string;
  account_number_last4: string;
  currency: string;
  balance: number;
  note?: string;
  created_at: string;
  updated_at: string;
}

/**
 * 建立銀行帳戶的輸入資料
 */
export interface CreateBankAccountInput {
  bank_name: string;
  account_type: string;
  account_number_last4: string;
  currency: string;
  balance: number;
  note?: string;
}

/**
 * 更新銀行帳戶的輸入資料
 */
export interface UpdateBankAccountInput {
  bank_name?: string;
  account_type?: string;
  account_number_last4?: string;
  currency?: string;
  balance?: number;
  note?: string;
}

/**
 * 銀行帳戶篩選條件
 */
export interface BankAccountFilters
  extends Record<string, string | number | boolean | undefined | null> {
  currency?: string;
}

/**
 * 信用卡資料結構
 */
export interface CreditCard {
  id: string;
  issuing_bank: string;
  card_name: string;
  card_number_last4: string;
  billing_day: number;
  payment_due_day: number;
  credit_limit: number;
  used_credit: number;
  group_id?: string;
  note?: string;
  created_at: string;
  updated_at: string;
}

/**
 * 建立信用卡的輸入資料
 */
export interface CreateCreditCardInput {
  issuing_bank: string;
  card_name: string;
  card_number_last4: string;
  billing_day: number;
  payment_due_day: number;
  credit_limit: number;
  used_credit: number;
  note?: string;
}

/**
 * 更新信用卡的輸入資料
 */
export interface UpdateCreditCardInput {
  issuing_bank?: string;
  card_name?: string;
  card_number_last4?: string;
  billing_day?: number;
  payment_due_day?: number;
  credit_limit?: number;
  used_credit?: number;
  note?: string;
}

/**
 * 信用卡查詢參數
 */
export interface CreditCardQueryParams
  extends Record<string, string | number | boolean | undefined | null> {
  days_ahead?: number;
}

/**
 * 信用卡群組資料結構
 */
export interface CreditCardGroup {
  id: string;
  name: string;
  issuing_bank: string;
  shared_credit_limit: number;
  note?: string;
  created_at: string;
  updated_at: string;
}

/**
 * 信用卡群組（包含卡片列表）
 */
export interface CreditCardGroupWithCards extends CreditCardGroup {
  cards: CreditCard[];
  total_used_credit: number;
}

/**
 * 建立信用卡群組的輸入資料
 */
export interface CreateCreditCardGroupInput {
  name: string;
  issuing_bank: string;
  shared_credit_limit: number;
  card_ids: string[];
  note?: string;
}

/**
 * 更新信用卡群組的輸入資料
 */
export interface UpdateCreditCardGroupInput {
  name?: string;
  shared_credit_limit?: number;
  note?: string;
}

/**
 * 新增卡片到群組的輸入資料
 */
export interface AddCardsToGroupInput {
  card_ids: string[];
}

/**
 * 從群組移除卡片的輸入資料
 */
export interface RemoveCardsFromGroupInput {
  card_ids: string[];
}
