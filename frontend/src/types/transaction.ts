import { z } from "zod";

// ==================== 列舉型別 ====================

/**
 * 資產類型
 */
export const AssetType = {
  CASH: "cash",
  TW_STOCK: "tw-stock",
  US_STOCK: "us-stock",
  CRYPTO: "crypto",
} as const;

export type AssetType = (typeof AssetType)[keyof typeof AssetType];

/**
 * 交易類型
 */
export const TransactionType = {
  BUY: "buy",
  SELL: "sell",
  DIVIDEND: "dividend",
  FEE: "fee",
} as const;

export type TransactionType =
  (typeof TransactionType)[keyof typeof TransactionType];

/**
 * 幣別
 */
export const Currency = {
  TWD: "TWD",
  USD: "USD",
} as const;

export type Currency = (typeof Currency)[keyof typeof Currency];

// ==================== 資料模型 ====================

/**
 * 交易記錄
 */
export interface Transaction {
  id: string;
  date: string; // ISO 8601 格式
  asset_type: AssetType;
  symbol: string;
  name: string;
  type: TransactionType;
  quantity: number;
  price: number;
  amount: number;
  fee: number | null;
  tax: number | null;
  currency: Currency;
  exchange_rate_id: number | null; // 關聯的匯率記錄 ID（僅用於非 TWD 交易）
  note: string | null;
  created_at: string; // ISO 8601 格式
  updated_at: string; // ISO 8601 格式
}

/**
 * 建立交易的輸入資料
 */
export interface CreateTransactionInput {
  date: string; // ISO 8601 格式
  asset_type: AssetType;
  symbol: string;
  name: string;
  type: TransactionType;
  quantity: number;
  price: number;
  amount: number;
  fee?: number | null;
  tax?: number | null;
  currency: Currency;
  note?: string | null;
}

/**
 * 更新交易的輸入資料
 */
export interface UpdateTransactionInput {
  date?: string; // ISO 8601 格式
  asset_type?: AssetType;
  symbol?: string;
  name?: string;
  type?: TransactionType;
  quantity?: number;
  price?: number;
  amount?: number;
  fee?: number | null;
  tax?: number | null;
  currency?: Currency;
  note?: string | null;
}

/**
 * 交易列表的篩選條件
 */
export interface TransactionFilters {
  asset_type?: AssetType;
  symbol?: string;
  type?: TransactionType;
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
 * 資產類型 Schema
 */
export const assetTypeSchema = z.enum([
  AssetType.CASH,
  AssetType.TW_STOCK,
  AssetType.US_STOCK,
  AssetType.CRYPTO,
]);

/**
 * 交易類型 Schema
 */
export const transactionTypeSchema = z.enum([
  TransactionType.BUY,
  TransactionType.SELL,
  TransactionType.DIVIDEND,
  TransactionType.FEE,
]);

/**
 * 幣別 Schema
 */
export const currencySchema = z.enum([Currency.TWD, Currency.USD]);

/**
 * 建立交易的表單 Schema
 *
 * 數量驗證規則：
 * - 台股/美股：必須為正整數（股數）
 * - 加密貨幣：可以是小數（數量）
 */
export const createTransactionSchema = z
  .object({
    date: z.string().min(1, "日期為必填"),
    asset_type: assetTypeSchema,
    symbol: z.string().min(1, "代碼為必填"),
    name: z.string().min(1, "名稱為必填"),
    type: transactionTypeSchema,
    quantity: z
      .number({ message: "數量必須為數字" })
      .positive("數量必須大於 0"),
    price: z
      .number({ message: "價格必須為數字" })
      .nonnegative("價格不可為負數"),
    amount: z
      .number({ message: "金額必須為數字" })
      .nonnegative("金額不可為負數"),
    fee: z
      .number({ message: "手續費必須為數字" })
      .nonnegative("手續費不可為負數")
      .nullable()
      .optional(),
    tax: z
      .number({ message: "交易稅必須為數字" })
      .nonnegative("交易稅不可為負數")
      .nullable()
      .optional(),
    currency: currencySchema,
    note: z.string().nullable().optional(),
  })
  .refine(
    (data) => {
      // 台股和美股的數量必須為整數
      if (
        data.asset_type === AssetType.TW_STOCK ||
        data.asset_type === AssetType.US_STOCK
      ) {
        return Number.isInteger(data.quantity);
      }
      return true;
    },
    {
      message: "股票數量必須為整數",
      path: ["quantity"],
    }
  );

/**
 * 建立交易的表單資料型別
 */
export type CreateTransactionFormData = z.infer<typeof createTransactionSchema>;

/**
 * 更新交易的表單 Schema
 */
export const updateTransactionSchema = createTransactionSchema.partial();

/**
 * 更新交易的表單資料型別
 */
export type UpdateTransactionFormData = z.infer<typeof updateTransactionSchema>;

// ==================== 輔助函式 ====================

/**
 * 取得資產類型的顯示名稱
 */
export function getAssetTypeLabel(assetType: AssetType): string {
  const labels: Record<AssetType, string> = {
    [AssetType.CASH]: "現金",
    [AssetType.TW_STOCK]: "台股",
    [AssetType.US_STOCK]: "美股",
    [AssetType.CRYPTO]: "加密貨幣",
  };
  return labels[assetType];
}

/**
 * 取得交易類型的顯示名稱
 */
export function getTransactionTypeLabel(
  transactionType: TransactionType
): string {
  const labels: Record<TransactionType, string> = {
    [TransactionType.BUY]: "買入",
    [TransactionType.SELL]: "賣出",
    [TransactionType.DIVIDEND]: "股息",
    [TransactionType.FEE]: "手續費",
  };
  return labels[transactionType];
}

/**
 * 取得所有資產類型選項
 */
export function getAssetTypeOptions() {
  return Object.values(AssetType).map((value) => ({
    value,
    label: getAssetTypeLabel(value),
  }));
}

/**
 * 取得所有交易類型選項
 */
export function getTransactionTypeOptions() {
  return Object.values(TransactionType).map((value) => ({
    value,
    label: getTransactionTypeLabel(value),
  }));
}
