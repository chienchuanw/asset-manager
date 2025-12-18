// 訂閱相關型別定義

/**
 * 計費週期類型
 */
export type BillingCycle = "monthly" | "quarterly" | "yearly";

/**
 * 訂閱狀態
 */
export type SubscriptionStatus = "active" | "cancelled";

/**
 * 付款方式類型
 */
export type PaymentMethod = "cash" | "bank_account" | "credit_card";

/**
 * 訂閱資料結構
 */
export interface Subscription {
  id: string;
  name: string;
  amount: number;
  currency: string;
  billing_cycle: BillingCycle;
  billing_day: number;
  category_id: string;
  category_name?: string;
  category_type?: string;
  payment_method: PaymentMethod;
  account_id?: string;
  start_date: string;
  end_date?: string;
  auto_renew: boolean;
  status: SubscriptionStatus;
  note?: string;
  created_at: string;
  updated_at: string;
}

/**
 * 建立訂閱的輸入資料
 */
export interface CreateSubscriptionInput {
  name: string;
  amount: number;
  currency: string;
  billing_cycle: BillingCycle;
  billing_day: number;
  category_id: string;
  payment_method: PaymentMethod;
  account_id?: string;
  start_date: string;
  end_date?: string;
  auto_renew: boolean;
  note?: string;
}

/**
 * 更新訂閱的輸入資料
 */
export interface UpdateSubscriptionInput {
  name?: string;
  amount?: number;
  billing_cycle?: BillingCycle;
  billing_day?: number;
  category_id?: string;
  payment_method?: PaymentMethod;
  account_id?: string;
  end_date?: string;
  auto_renew?: boolean;
  note?: string;
}

/**
 * 取消訂閱的輸入資料
 */
export interface CancelSubscriptionInput {
  end_date: string;
}

/**
 * 訂閱列表篩選條件
 */
export interface SubscriptionFilters
  extends Record<string, string | number | boolean | undefined | null> {
  status?: SubscriptionStatus;
  category_id?: string;
  limit?: number;
  offset?: number;
}

/**
 * 訂閱統計資料
 */
export interface SubscriptionStats {
  total_count: number;
  active_count: number;
  cancelled_count: number;
  monthly_total: number;
  quarterly_total: number;
  yearly_total: number;
  total_monthly_cost: number; // 換算成每月成本
}
