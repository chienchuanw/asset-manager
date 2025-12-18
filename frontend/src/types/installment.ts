// 分期相關型別定義

import { PaymentMethod } from "./subscription";

/**
 * 分期狀態
 */
export type InstallmentStatus = "active" | "completed" | "cancelled";

/**
 * 分期資料結構
 */
export interface Installment {
  id: string;
  name: string;
  total_amount: number;
  currency: string;
  installment_count: number;
  paid_count: number;
  installment_amount: number; // 每期金額（後端欄位名稱）
  interest_rate: number;
  total_interest: number;
  category_id: string;
  category_name?: string;
  category_type?: string;
  payment_method: PaymentMethod;
  account_id?: string;
  start_date: string;
  billing_day: number;
  status: InstallmentStatus;
  note?: string;
  created_at: string;
  updated_at: string;
}

/**
 * 建立分期的輸入資料
 */
export interface CreateInstallmentInput {
  name: string;
  total_amount: number;
  currency: string;
  installment_count: number;
  interest_rate: number;
  category_id: string;
  payment_method: PaymentMethod;
  account_id?: string;
  start_date: string;
  billing_day: number;
  note?: string;
}

/**
 * 更新分期的輸入資料
 */
export interface UpdateInstallmentInput {
  name?: string;
  billing_day?: number;
  payment_method?: PaymentMethod;
  account_id?: string;
  note?: string;
}

/**
 * 分期列表篩選條件
 */
export interface InstallmentFilters
  extends Record<string, string | number | boolean | undefined | null> {
  status?: InstallmentStatus;
  category_id?: string;
  limit?: number;
  offset?: number;
}

/**
 * 分期統計資料
 */
export interface InstallmentStats {
  total_count: number;
  active_count: number;
  completed_count: number;
  cancelled_count: number;
  total_principal: number;
  total_interest: number;
  total_paid: number;
  total_remaining: number;
  monthly_payment: number; // 每月總付款金額
}

/**
 * 即將完成的分期查詢參數
 */
export interface CompletingSoonParams
  extends Record<string, string | number | boolean | undefined | null> {
  months?: number;
}
