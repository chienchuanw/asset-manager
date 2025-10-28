import { Instrument, InstrumentSearchResult } from "@/types/instrument";
import { AssetType, Transaction } from "@/types/transaction";
import { transactionsAPI } from "./transactions";
import { getCommonInstrumentsByType } from "@/lib/data/common-instruments";

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
 * 合併歷史交易標的與常用標的清單
 * 
 * @param assetType 資產類型
 * @returns 合併後的標的清單（去重，歷史記錄優先）
 * 
 * @example
 * ```typescript
 * const instruments = await getMergedInstruments(AssetType.TW_STOCK);
 * // 回傳合併後的台股標的清單
 * ```
 */
export async function getMergedInstruments(
  assetType: AssetType
): Promise<InstrumentSearchResult[]> {
  // 取得歷史交易標的
  const fromHistory = await getInstrumentsFromTransactions(assetType);

  // 取得常用標的
  const fromCommon = getCommonInstrumentsByType(assetType);

  // 使用 Map 去重（歷史記錄優先）
  const mergedMap = new Map<string, InstrumentSearchResult>();

  // 先加入歷史記錄
  fromHistory.forEach((instrument) => {
    mergedMap.set(instrument.symbol, instrument);
  });

  // 再加入常用標的（如果不存在）
  fromCommon.forEach((instrument) => {
    if (!mergedMap.has(instrument.symbol)) {
      mergedMap.set(instrument.symbol, {
        ...instrument,
        from_history: false,
        usage_count: 0,
      });
    }
  });

  // 轉換為陣列
  // 排序規則：
  // 1. 歷史記錄優先（from_history = true）
  // 2. 使用次數多的優先
  // 3. symbol 字母順序
  return Array.from(mergedMap.values()).sort((a, b) => {
    // 歷史記錄優先
    if (a.from_history && !b.from_history) return -1;
    if (!a.from_history && b.from_history) return 1;

    // 使用次數排序
    const usageA = a.usage_count || 0;
    const usageB = b.usage_count || 0;
    if (usageA !== usageB) return usageB - usageA;

    // symbol 字母順序
    return a.symbol.localeCompare(b.symbol);
  });
}

/**
 * 搜尋標的（包含歷史記錄與常用標的）
 * 
 * @param query 搜尋關鍵字
 * @param assetType 資產類型
 * @returns 符合條件的標的清單
 * 
 * @example
 * ```typescript
 * const results = await searchInstruments("台積", AssetType.TW_STOCK);
 * // 回傳包含「台積」的標的
 * ```
 */
export async function searchInstruments(
  query: string,
  assetType: AssetType
): Promise<InstrumentSearchResult[]> {
  // 取得合併後的標的清單
  const instruments = await getMergedInstruments(assetType);

  // 如果沒有搜尋關鍵字，回傳全部
  if (!query || query.trim() === "") {
    return instruments;
  }

  const normalizedQuery = query.toLowerCase().trim();

  // 模糊搜尋 symbol 或 name
  return instruments.filter(
    (instrument) =>
      instrument.symbol.toLowerCase().includes(normalizedQuery) ||
      instrument.name.toLowerCase().includes(normalizedQuery)
  );
}

