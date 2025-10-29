import { useQuery } from "@tanstack/react-query";
import { AssetType } from "@/types/transaction";
import { InstrumentSearchResult } from "@/types/instrument";
import { getInstruments } from "@/lib/api/instruments";

/**
 * 取得標的清單的 Hook（僅從交易記錄提取）
 *
 * 功能：
 * 1. 從歷史交易記錄中提取標的
 * 2. 按使用次數排序（最常用的在最上方）
 * 3. 使用 React Query 快取結果
 *
 * @param assetType 資產類型
 * @returns React Query 結果
 *
 * @example
 * ```typescript
 * const { data: instruments, isLoading } = useInstruments(AssetType.TW_STOCK);
 * ```
 */
export function useInstruments(assetType: AssetType) {
  return useQuery<InstrumentSearchResult[], Error>({
    queryKey: ["instruments", assetType],
    queryFn: () => getInstruments(assetType),
    // 快取 5 分鐘（避免頻繁請求）
    staleTime: 5 * 60 * 1000,
    // 快取保留 30 分鐘
    gcTime: 30 * 60 * 1000,
    // 失敗時不自動重試（避免消耗 API 配額）
    retry: false,
  });
}
