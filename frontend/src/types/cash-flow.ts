import { z } from "zod";
import { Currency } from "./transaction";

// ==================== 列舉型別 ====================

/**
 * 現金流類型
 */
export const CashFlowType = {
  INCOME: "income",
  EXPENSE: "expense",
} as const;

export type CashFlowType = (typeof CashFlowType)[keyof typeof CashFlowType];

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
 * 現金流列表的篩選條件
 */
export interface CashFlowFilters {
  type?: CashFlowType;
  category_id?: string;
  start_date?: string; // ISO 8601 格式
  end_date?: string; // ISO 8601 格式
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
]);

/**
 * 建立現金流記錄的表單 Schema
 */
export const createCashFlowSchema = z.object({
  date: z.string().min(1, "日期為必填"),
  type: cashFlowTypeSchema,
  category_id: z.string().min(1, "分類為必填"),
  amount: z
    .number({ message: "金額必須為數字" })
    .positive("金額必須大於 0"),
  description: z
    .string()
    .min(1, "描述為必填")
    .max(500, "描述不可超過 500 字元"),
  note: z.string().max(1000, "備註不可超過 1000 字元").nullable().optional(),
});

/**
 * 建立現金流記錄的表單資料型別
 */
export type CreateCashFlowFormData = z.infer<typeof createCashFlowSchema>;

/**
 * 更新現金流記錄的表單 Schema
 */
export const updateCashFlowSchema = createCashFlowSchema.partial();

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
  };
  return labels[cashFlowType];
}

/**
 * 取得現金流類型的顏色（用於 UI 顯示）
 */
export function getCashFlowTypeColor(cashFlowType: CashFlowType): string {
  const colors: Record<CashFlowType, string> = {
    [CashFlowType.INCOME]: "text-green-600",
    [CashFlowType.EXPENSE]: "text-red-600",
  };
  return colors[cashFlowType];
}

/**
 * 取得現金流類型的背景顏色（用於 UI 顯示）
 */
export function getCashFlowTypeBgColor(cashFlowType: CashFlowType): string {
  const colors: Record<CashFlowType, string> = {
    [CashFlowType.INCOME]: "bg-green-100",
    [CashFlowType.EXPENSE]: "bg-red-100",
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
    expense: ["飲食", "交通", "娛樂", "醫療", "房租", "水電", "保險", "其他支出"],
  };
}

