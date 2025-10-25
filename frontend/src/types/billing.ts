// 扣款相關型別定義

/**
 * 扣款結果
 */
export interface BillingResult {
  count: number;
  total_amount: number;
  cash_flow_ids: string[];
}

/**
 * 每日扣款結果
 */
export interface DailyBillingResult {
  subscription_count: number;
  installment_count: number;
  total_amount: number;
  subscription_amount: number;
  installment_amount: number;
  subscription_cash_flow_ids: string[];
  installment_cash_flow_ids: string[];
}

/**
 * 處理扣款的輸入參數
 */
export interface ProcessBillingInput {
  date?: string; // ISO 8601 格式，例如 "2024-01-15"
}

