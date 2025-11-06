import { z } from "zod";
import { Currency } from "./transaction";

// ==================== 列舉型別 ====================

/**
 * 現金流類型
 */
export const CashFlowType = {
  INCOME: "income",
  EXPENSE: "expense",
  TRANSFER_IN: "transfer_in", // 存入帳戶
  TRANSFER_OUT: "transfer_out", // 從帳戶轉出
} as const;

export type CashFlowType = (typeof CashFlowType)[keyof typeof CashFlowType];

/**
 * 付款來源類型
 */
export const SourceType = {
  MANUAL: "manual", // 手動建立（現金交易）
  SUBSCRIPTION: "subscription", // 訂閱自動產生
  INSTALLMENT: "installment", // 分期自動產生
  BANK_ACCOUNT: "bank_account", // 銀行帳戶交易
  CREDIT_CARD: "credit_card", // 信用卡交易
} as const;

export type SourceType = (typeof SourceType)[keyof typeof SourceType];

/**
 * 付款方式類型（用於前端 UI）
 */
export const PaymentMethodType = {
  CASH: "cash", // 現金
  BANK_ACCOUNT: "bank_account", // 銀行帳戶
  CREDIT_CARD: "credit_card", // 信用卡
} as const;

export type PaymentMethodType =
  (typeof PaymentMethodType)[keyof typeof PaymentMethodType];

// ==================== 資料模型 ====================

/**
 * 現金流分類
 */
export interface CashFlowCategory {
  id: string;
  name: string;
  type: CashFlowType;
  is_system: boolean;
  created_at: string; // ISO 8601 格式
  updated_at: string; // ISO 8601 格式
}

/**
 * 現金流記錄
 */
export interface CashFlow {
  id: string;
  date: string; // ISO 8601 格式
  type: CashFlowType;
  category_id: string;
  amount: number;
  currency: Currency;
  description: string;
  note: string | null;
  source_type?: SourceType | null; // 付款來源類型
  source_id?: string | null; // 付款來源 ID（銀行帳戶或信用卡 ID）
  target_type?: SourceType | null; // 轉帳目標類型（用於 transfer_out）
  target_id?: string | null; // 轉帳目標 ID（用於 transfer_out，例如信用卡 ID）
  created_at: string; // ISO 8601 格式
  updated_at: string; // ISO 8601 格式
  category?: CashFlowCategory; // 關聯的分類資料（可選）
}

/**
 * 現金流摘要
 */
export interface CashFlowSummary {
  total_income: number;
  total_expense: number;
  net_cash_flow: number;
}

// ==================== 輸入資料型別 ====================

/**
 * 建立現金流記錄的輸入資料
 */
export interface CreateCashFlowInput {
  date: string; // ISO 8601 格式
  type: CashFlowType;
  category_id: string;
  amount: number;
  description: string;
  note?: string | null;
  source_type?: SourceType | null; // 付款來源類型
  source_id?: string | null; // 付款來源 ID
  target_type?: SourceType | null; // 轉帳目標類型（用於 transfer_out）
  target_id?: string | null; // 轉帳目標 ID（用於 transfer_out）
}

/**
 * 更新現金流記錄的輸入資料
 */
export interface UpdateCashFlowInput {
  date?: string; // ISO 8601 格式
  type?: CashFlowType;
  category_id?: string;
  amount?: number;
  description?: string;
  note?: string | null;
  source_type?: SourceType | null; // 付款來源類型
  source_id?: string | null; // 付款來源 ID
  target_type?: SourceType | null; // 轉帳目標類型（用於 transfer_out）
  target_id?: string | null; // 轉帳目標 ID（用於 transfer_out）
}

/**
 * 建立分類的輸入資料
 */
export interface CreateCategoryInput {
  name: string;
  type: CashFlowType;
}

/**
 * 更新分類的輸入資料
 */
export interface UpdateCategoryInput {
  name: string;
}

/**
 * 付款方式選項
 */
export interface PaymentMethodOption {
  type: PaymentMethodType;
  account_id?: string; // 當 type 為 bank_account 或 credit_card 時必填
}

/**
 * 現金流列表的篩選條件
 */
export interface CashFlowFilters {
  type?: CashFlowType;
  category_id?: string;
  start_date?: string; // ISO 8601 格式
  end_date?: string; // ISO 8601 格式
  source_type?: SourceType; // 付款來源類型篩選
  limit?: number;
  offset?: number;
  [key: string]: string | number | undefined;
}

// ==================== API 回應格式 ====================

/**
 * API 錯誤回應
 */
export interface APIError {
  code: string;
  message: string;
}

/**
 * API 回應（泛型）
 */
export interface APIResponse<T> {
  data: T | null;
  error: APIError | null;
}

// ==================== Zod Schema（用於表單驗證）====================

/**
 * 現金流類型 Schema
 */
export const cashFlowTypeSchema = z.enum([
  CashFlowType.INCOME,
  CashFlowType.EXPENSE,
  CashFlowType.TRANSFER_IN,
  CashFlowType.TRANSFER_OUT,
]);

/**
 * 付款來源類型 Schema
 */
export const sourceTypeSchema = z.enum([
  SourceType.MANUAL,
  SourceType.SUBSCRIPTION,
  SourceType.INSTALLMENT,
  SourceType.BANK_ACCOUNT,
  SourceType.CREDIT_CARD,
]);

/**
 * 付款方式類型 Schema
 */
export const paymentMethodTypeSchema = z.enum([
  PaymentMethodType.CASH,
  PaymentMethodType.BANK_ACCOUNT,
  PaymentMethodType.CREDIT_CARD,
]);

/**
 * 建立現金流記錄的表單 Schema
 */
export const createCashFlowSchema = z
  .object({
    date: z.string().min(1, "日期為必填"),
    type: cashFlowTypeSchema,
    category_id: z.string().min(1, "分類為必填"),
    amount: z.number({ message: "金額必須為數字" }).positive("金額必須大於 0"),
    description: z
      .string()
      .min(1, "描述為必填")
      .max(500, "描述不可超過 500 字元"),
    note: z.string().max(1000, "備註不可超過 1000 字元").nullable().optional(),
    // 付款方式相關欄位
    payment_method: paymentMethodTypeSchema,
    account_id: z.string().optional(), // 當付款方式為銀行帳戶或信用卡時必填
    // 轉帳目標相關欄位（用於 transfer_out）
    target_payment_method: paymentMethodTypeSchema.optional(), // 轉帳目標付款方式
    target_account_id: z.string().optional(), // 轉帳目標帳戶 ID
  })
  .refine(
    (data) => {
      // 轉帳類型必須選擇銀行帳戶
      if (
        data.type === CashFlowType.TRANSFER_IN ||
        data.type === CashFlowType.TRANSFER_OUT
      ) {
        return (
          data.payment_method === PaymentMethodType.BANK_ACCOUNT &&
          data.account_id &&
          data.account_id.length > 0
        );
      }

      // 當付款方式為銀行帳戶或信用卡時，account_id 為必填
      if (
        data.payment_method === PaymentMethodType.BANK_ACCOUNT ||
        data.payment_method === PaymentMethodType.CREDIT_CARD
      ) {
        return data.account_id && data.account_id.length > 0;
      }
      return true;
    },
    {
      message: "轉帳類型必須選擇銀行帳戶",
      path: ["account_id"],
    }
  )
  .refine(
    (data) => {
      // transfer_out 類型必須選擇轉帳目標
      if (data.type === CashFlowType.TRANSFER_OUT) {
        return (
          data.target_payment_method &&
          data.target_account_id &&
          data.target_account_id.length > 0
        );
      }
      return true;
    },
    {
      message: "轉出類型必須選擇轉帳目標",
      path: ["target_account_id"],
    }
  );

/**
 * 建立現金流記錄的表單資料型別
 */
export type CreateCashFlowFormData = z.infer<typeof createCashFlowSchema>;

/**
 * 更新現金流記錄的表單 Schema
 */
export const updateCashFlowSchema = z
  .object({
    date: z.string().min(1, "日期為必填").optional(),
    type: cashFlowTypeSchema.optional(),
    category_id: z.string().min(1, "分類為必填").optional(),
    amount: z
      .number({ message: "金額必須為數字" })
      .positive("金額必須大於 0")
      .optional(),
    description: z
      .string()
      .min(1, "描述為必填")
      .max(500, "描述不可超過 500 字元")
      .optional(),
    note: z.string().max(1000, "備註不可超過 1000 字元").nullable().optional(),
    // 付款方式相關欄位
    payment_method: paymentMethodTypeSchema.optional(),
    account_id: z.string().optional(),
  })
  .refine(
    (data) => {
      // 當付款方式為銀行帳戶或信用卡時，account_id 為必填
      if (
        data.payment_method === PaymentMethodType.BANK_ACCOUNT ||
        data.payment_method === PaymentMethodType.CREDIT_CARD
      ) {
        return data.account_id && data.account_id.length > 0;
      }
      return true;
    },
    {
      message: "請選擇帳戶",
      path: ["account_id"],
    }
  );

/**
 * 更新現金流記錄的表單資料型別
 */
export type UpdateCashFlowFormData = z.infer<typeof updateCashFlowSchema>;

/**
 * 建立分類的表單 Schema
 */
export const createCategorySchema = z.object({
  name: z
    .string()
    .min(1, "分類名稱為必填")
    .max(100, "分類名稱不可超過 100 字元"),
  type: cashFlowTypeSchema,
});

/**
 * 建立分類的表單資料型別
 */
export type CreateCategoryFormData = z.infer<typeof createCategorySchema>;

/**
 * 更新分類的表單 Schema
 */
export const updateCategorySchema = z.object({
  name: z
    .string()
    .min(1, "分類名稱為必填")
    .max(100, "分類名稱不可超過 100 字元"),
});

/**
 * 更新分類的表單資料型別
 */
export type UpdateCategoryFormData = z.infer<typeof updateCategorySchema>;

// ==================== 輔助函式 ====================

/**
 * 取得現金流類型的顯示名稱
 */
export function getCashFlowTypeLabel(cashFlowType: CashFlowType): string {
  const labels: Record<CashFlowType, string> = {
    [CashFlowType.INCOME]: "收入",
    [CashFlowType.EXPENSE]: "支出",
    [CashFlowType.TRANSFER_IN]: "存入",
    [CashFlowType.TRANSFER_OUT]: "轉出",
  };
  return labels[cashFlowType];
}

/**
 * 取得現金流類型的顏色（用於 UI 顯示）
 * 台灣習慣：紅漲綠跌 - 收入為紅色，支出為綠色
 */
export function getCashFlowTypeColor(cashFlowType: CashFlowType): string {
  const colors: Record<CashFlowType, string> = {
    [CashFlowType.INCOME]: "text-red-600",
    [CashFlowType.EXPENSE]: "text-green-600",
    [CashFlowType.TRANSFER_IN]: "text-gray-600",
    [CashFlowType.TRANSFER_OUT]: "text-gray-600",
  };
  return colors[cashFlowType];
}

/**
 * 取得現金流類型的背景顏色（用於 UI 顯示）
 * 台灣習慣：紅漲綠跌 - 收入為紅色，支出為綠色，轉帳為灰色
 */
export function getCashFlowTypeBgColor(cashFlowType: CashFlowType): string {
  const colors: Record<CashFlowType, string> = {
    [CashFlowType.INCOME]: "bg-red-100",
    [CashFlowType.EXPENSE]: "bg-green-100",
    [CashFlowType.TRANSFER_IN]: "bg-gray-100",
    [CashFlowType.TRANSFER_OUT]: "bg-gray-100",
  };
  return colors[cashFlowType];
}

/**
 * 取得所有現金流類型選項
 */
export function getCashFlowTypeOptions() {
  return Object.values(CashFlowType).map((value) => ({
    value,
    label: getCashFlowTypeLabel(value),
  }));
}

/**
 * 取得付款方式類型的顯示名稱
 */
export function getPaymentMethodTypeLabel(
  paymentMethodType: PaymentMethodType
): string {
  const labels: Record<PaymentMethodType, string> = {
    [PaymentMethodType.CASH]: "現金",
    [PaymentMethodType.BANK_ACCOUNT]: "銀行帳戶",
    [PaymentMethodType.CREDIT_CARD]: "信用卡",
  };
  return labels[paymentMethodType];
}

/**
 * 取得所有付款方式類型選項
 */
export function getPaymentMethodTypeOptions() {
  return Object.values(PaymentMethodType).map((value) => ({
    value,
    label: getPaymentMethodTypeLabel(value),
  }));
}

/**
 * 將付款方式類型轉換為來源類型
 */
export function paymentMethodTypeToSourceType(
  paymentMethodType: PaymentMethodType
): SourceType {
  const mapping: Record<PaymentMethodType, SourceType> = {
    [PaymentMethodType.CASH]: SourceType.MANUAL,
    [PaymentMethodType.BANK_ACCOUNT]: SourceType.BANK_ACCOUNT,
    [PaymentMethodType.CREDIT_CARD]: SourceType.CREDIT_CARD,
  };
  return mapping[paymentMethodType];
}

/**
 * 將來源類型轉換為付款方式類型
 */
export function sourceTypeToPaymentMethodType(
  sourceType: SourceType
): PaymentMethodType {
  const mapping: Record<SourceType, PaymentMethodType> = {
    [SourceType.MANUAL]: PaymentMethodType.CASH,
    [SourceType.SUBSCRIPTION]: PaymentMethodType.CASH, // 訂閱預設為現金
    [SourceType.INSTALLMENT]: PaymentMethodType.CASH, // 分期預設為現金
    [SourceType.BANK_ACCOUNT]: PaymentMethodType.BANK_ACCOUNT,
    [SourceType.CREDIT_CARD]: PaymentMethodType.CREDIT_CARD,
  };
  return mapping[sourceType];
}

/**
 * 格式化金額顯示（加上千分位逗號）
 */
export function formatAmount(amount: number): string {
  return new Intl.NumberFormat("zh-TW", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(amount);
}

/**
 * 格式化日期顯示
 */
export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat("zh-TW", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(date);
}

/**
 * 取得系統預設分類列表
 */
export function getSystemCategories(): {
  income: string[];
  expense: string[];
} {
  return {
    income: ["薪資", "獎金", "利息", "其他收入"],
    expense: [
      "飲食",
      "交通",
      "娛樂",
      "醫療",
      "房租",
      "水電",
      "保險",
      "其他支出",
    ],
  };
}
