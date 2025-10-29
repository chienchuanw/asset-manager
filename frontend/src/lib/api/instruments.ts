import { InstrumentSearchResult } from "@/types/instrument";
import { AssetType, Transaction } from "@/types/transaction";
import { transactionsAPI } from "./transactions";

/**
 * 從交易記錄中提取唯一的標的清單
 *
 * @param assetType 資產類型（可選）
 * @returns 標的清單（包含使用次數）
 *
 * @example
 * ```typescript
 * const instruments = await getInstrumentsFromTransactions(AssetType.TW_STOCK);
 * // 回傳台股的所有歷史交易標的
 * ```
 */
export async function getInstrumentsFromTransactions(
  assetType?: AssetType
): Promise<InstrumentSearchResult[]> {
  try {
    // 取得所有交易記錄
    const transactions = await transactionsAPI.getAll(
      assetType ? { asset_type: assetType } : undefined
    );

    // 使用 Map 來統計每個標的的使用次數
    const instrumentMap = new Map<string, InstrumentSearchResult>();

    transactions.forEach((transaction: Transaction) => {
      const key = `${transaction.asset_type}-${transaction.symbol}`;

      if (instrumentMap.has(key)) {
        // 如果已存在，增加使用次數
        const existing = instrumentMap.get(key)!;
        existing.usage_count = (existing.usage_count || 0) + 1;
      } else {
        // 新增標的
        instrumentMap.set(key, {
          symbol: transaction.symbol,
          name: transaction.name,
          asset_type: transaction.asset_type,
          from_history: true,
          usage_count: 1,
        });
      }
    });

    // 轉換為陣列並按使用次數排序（降序）
    return Array.from(instrumentMap.values()).sort(
      (a, b) => (b.usage_count || 0) - (a.usage_count || 0)
    );
  } catch (error) {
    console.error("Failed to fetch instruments from transactions:", error);
    // 如果 API 失敗，回傳空陣列
    return [];
  }
}

/**
 * 取得標的清單（僅從交易記錄提取）
 *
 * @param assetType 資產類型
 * @returns 標的清單（按使用次數排序）
 *
 * @example
 * ```typescript
 * const instruments = await getInstruments(AssetType.TW_STOCK);
 * // 回傳台股的所有歷史交易標的
 * ```
 */
export async function getInstruments(
  assetType: AssetType
): Promise<InstrumentSearchResult[]> {
  return getInstrumentsFromTransactions(assetType);
}
