import { apiClient, type APIResponseWithWarnings } from "./client";
import type { Holding, HoldingFilters } from "@/types/holding";

/**
 * Holdings API 端點
 */
const ENDPOINTS = {
  HOLDINGS: "/api/holdings",
  HOLDING_BY_SYMBOL: (symbol: string) => `/api/holdings/${symbol}`,
  FIX_INSUFFICIENT_QUANTITY: "/api/holdings/fix-insufficient-quantity",
  RECONCILE_PREVIEW: "/api/holdings/reconcile/preview",
  RECONCILE: "/api/holdings/reconcile",
} as const;

/** 對帳項目（請求） */
export interface HoldingReconcileItem {
  asset_type: "tw-stock" | "us-stock" | "crypto";
  symbol: string;
  name?: string;
  currency: "TWD" | "USD";
  target_quantity: number;
  target_avg_cost: number;
  reason?: string;
}

export type ReconcileAction = "create" | "update" | "liquidate" | "noop";

/** 對帳預覽單筆結果 */
export interface HoldingReconcilePreviewItem {
  item: HoldingReconcileItem;
  prev_quantity: number;
  prev_avg_cost: number;
  quantity_delta: number;
  avg_cost_delta: number;
  action: ReconcileAction;
}

export interface HoldingReconcilePreview {
  items: HoldingReconcilePreviewItem[];
}

/**
 * 修復不足數量的輸入
 */
export interface FixInsufficientQuantityInput {
  symbol: string;
  current_holding: number;
  estimated_cost?: number;
}

/**
 * Holdings API
 */
export const holdingsAPI = {
  /**
   * 取得所有持倉（包含 warnings）
   * @param filters 篩選條件
   * @returns 持倉陣列和警告
   *
   * @example
   * ```typescript
   * // 取得所有持倉
   * const { data: holdings, warnings } = await holdingsAPI.getAll();
   *
   * // 只取得台股持倉
   * const { data: twStocks, warnings } = await holdingsAPI.getAll({ asset_type: "tw-stock" });
   * ```
   */
  getAll: async (
    filters?: HoldingFilters
  ): Promise<APIResponseWithWarnings<Holding[]>> => {
    // 將 filters 轉換為符合 apiClient.get 的參數格式
    const params: Record<string, string | undefined> = {};
    if (filters?.asset_type) {
      params.asset_type = filters.asset_type;
    }
    if (filters?.symbol) {
      params.symbol = filters.symbol;
    }

    return apiClient.getWithWarnings<Holding[]>(ENDPOINTS.HOLDINGS, {
      params: Object.keys(params).length > 0 ? params : undefined,
    });
  },

  /**
   * 取得單一持倉
   * @param symbol 標的代碼
   * @returns 持倉資料
   *
   * @example
   * ```typescript
   * const holding = await holdingsAPI.getBySymbol("2330");
   * ```
   */
  getBySymbol: async (symbol: string): Promise<Holding> => {
    return apiClient.get<Holding>(ENDPOINTS.HOLDING_BY_SYMBOL(symbol));
  },

  /**
   * 修復不足數量問題
   * @param input 修復輸入
   * @returns 建立的交易記錄
   *
   * @example
   * ```typescript
   * const transaction = await holdingsAPI.fixInsufficientQuantity({
   *   symbol: "2330",
   *   current_holding: 100,
   *   estimated_cost: 550,
   * });
   * ```
   */
  fixInsufficientQuantity: async (input: FixInsufficientQuantityInput) => {
    return apiClient.post(ENDPOINTS.FIX_INSUFFICIENT_QUANTITY, input);
  },

  /** 預覽持倉對帳（dry-run，不寫入） */
  reconcilePreview: async (
    items: HoldingReconcileItem[]
  ): Promise<HoldingReconcilePreview> => {
    return apiClient.post<HoldingReconcilePreview>(ENDPOINTS.RECONCILE_PREVIEW, {
      items,
    });
  },

  /** 執行持倉對帳（寫入 adjustment 交易） */
  reconcile: async (
    items: HoldingReconcileItem[]
  ): Promise<HoldingReconcilePreview> => {
    return apiClient.post<HoldingReconcilePreview>(ENDPOINTS.RECONCILE, {
      items,
    });
  },
};
