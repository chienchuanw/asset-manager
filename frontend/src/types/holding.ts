import { z } from "zod";
import { AssetType, assetTypeSchema } from "./transaction";

// ==================== 資料模型 ====================

/**
 * 持倉明細
 * 對應後端 Holding 結構
 */
export interface Holding {
  symbol: string;
  name: string;
  asset_type: AssetType;
  quantity: number;
  avg_cost: number;
  total_cost: number;
  current_price: number;
  currency: string;
  current_price_twd: number; // TWD 轉換後的價格
  market_value: number;
  unrealized_pl: number;
  unrealized_pl_pct: number;
  price_source?: string; // 價格來源 (cache, api, stale-cache)
  is_price_stale?: boolean; // 價格是否過期
  price_stale_reason?: string; // 價格過期原因
}

/**
 * 持倉篩選條件
 */
export interface HoldingFilters {
  asset_type?: AssetType;
  symbol?: string;
}

// ==================== Zod Schema（用於驗證）====================

/**
 * 持倉 Schema
 */
export const holdingSchema = z.object({
  symbol: z.string(),
  name: z.string(),
  asset_type: assetTypeSchema,
  quantity: z.number(),
  avg_cost: z.number(),
  total_cost: z.number(),
  current_price: z.number(),
  currency: z.string(),
  current_price_twd: z.number(),
  market_value: z.number(),
  unrealized_pl: z.number(),
  unrealized_pl_pct: z.number(),
  price_source: z.string().optional(),
  is_price_stale: z.boolean().optional(),
  price_stale_reason: z.string().optional(),
});

/**
 * 持倉篩選條件 Schema
 */
export const holdingFiltersSchema = z.object({
  asset_type: assetTypeSchema.optional(),
  symbol: z.string().optional(),
});

// ==================== 輔助函式 ====================

/**
 * 格式化貨幣顯示
 * @param value 金額
 * @param currency 貨幣代碼
 * @returns 格式化後的字串
 */
export function formatCurrency(value: number, currency: string): string {
  const prefix = currency === "TWD" ? "NT$" : currency === "USD" ? "$" : "";
  return `${prefix} ${value.toLocaleString("zh-TW", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  })}`;
}

/**
 * 格式化百分比顯示
 * @param value 百分比數值
 * @returns 格式化後的字串
 */
export function formatPercentage(value: number): string {
  const sign = value >= 0 ? "+" : "";
  return `${sign}${value.toFixed(2)}%`;
}

/**
 * 取得損益顏色類別（台灣股市習慣：紅漲綠跌）
 * @param value 損益數值
 * @returns Tailwind CSS 類別名稱
 */
export function getProfitLossColor(value: number): string {
  return value >= 0 ? "text-red-600" : "text-green-600";
}

/**
 * 計算總市值
 * @param holdings 持倉陣列
 * @returns 總市值
 */
export function calculateTotalMarketValue(holdings: Holding[]): number {
  return holdings.reduce((sum, h) => sum + h.market_value, 0);
}

/**
 * 計算總成本
 * @param holdings 持倉陣列
 * @returns 總成本
 */
export function calculateTotalCost(holdings: Holding[]): number {
  return holdings.reduce((sum, h) => sum + h.total_cost, 0);
}

/**
 * 計算總損益
 * @param holdings 持倉陣列
 * @returns 總損益
 */
export function calculateTotalProfitLoss(holdings: Holding[]): number {
  return holdings.reduce((sum, h) => sum + h.unrealized_pl, 0);
}

/**
 * 計算總損益百分比
 * @param holdings 持倉陣列
 * @returns 總損益百分比
 */
export function calculateTotalProfitLossPct(holdings: Holding[]): number {
  const totalCost = calculateTotalCost(holdings);
  const totalProfitLoss = calculateTotalProfitLoss(holdings);
  return totalCost > 0 ? (totalProfitLoss / totalCost) * 100 : 0;
}

/**
 * 按資產類型分組持倉
 * @param holdings 持倉陣列
 * @returns 分組後的持倉
 */
export function groupHoldingsByAssetType(
  holdings: Holding[]
): Record<AssetType, Holding[]> {
  return holdings.reduce((groups, holding) => {
    const type = holding.asset_type;
    if (!groups[type]) {
      groups[type] = [];
    }
    groups[type].push(holding);
    return groups;
  }, {} as Record<AssetType, Holding[]>);
}

/**
 * 排序持倉
 * @param holdings 持倉陣列
 * @param sortBy 排序欄位
 * @param order 排序順序
 * @returns 排序後的持倉陣列
 */
export function sortHoldings(
  holdings: Holding[],
  sortBy: "market_value" | "unrealized_pl" | "quantity" | "symbol",
  order: "asc" | "desc" = "desc"
): Holding[] {
  const sorted = [...holdings].sort((a, b) => {
    let compareValue = 0;
    switch (sortBy) {
      case "market_value":
        compareValue = a.market_value - b.market_value;
        break;
      case "unrealized_pl":
        compareValue = a.unrealized_pl - b.unrealized_pl;
        break;
      case "quantity":
        compareValue = a.quantity - b.quantity;
        break;
      case "symbol":
        compareValue = a.symbol.localeCompare(b.symbol);
        break;
    }
    return order === "asc" ? compareValue : -compareValue;
  });
  return sorted;
}

/**
 * 搜尋持倉
 * @param holdings 持倉陣列
 * @param query 搜尋關鍵字
 * @returns 符合條件的持倉陣列
 */
export function searchHoldings(holdings: Holding[], query: string): Holding[] {
  if (!query) return holdings;
  const lowerQuery = query.toLowerCase();
  return holdings.filter(
    (h) =>
      h.symbol.toLowerCase().includes(lowerQuery) ||
      h.name.toLowerCase().includes(lowerQuery)
  );
}
