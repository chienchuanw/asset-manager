import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createDailySnapshot,
  getLatestTrend,
  getTrendByDateRange,
} from "@/lib/api/performance-trends";

/**
 * 建立每日績效快照
 */
export function useCreateDailySnapshot() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createDailySnapshot,
    onSuccess: () => {
      // 快照建立成功後，重新取得趨勢資料
      queryClient.invalidateQueries({ queryKey: ["performance-trends"] });
    },
  });
}

/**
 * 取得日期範圍內的績效趨勢
 */
export function useTrendByDateRange(startDate: string, endDate: string) {
  return useQuery({
    queryKey: ["performance-trends", "range", startDate, endDate],
    queryFn: () => getTrendByDateRange(startDate, endDate),
    enabled: !!startDate && !!endDate,
  });
}

/**
 * 取得最新的 N 天績效趨勢
 */
export function useLatestTrend(days: number = 30) {
  return useQuery({
    queryKey: ["performance-trends", "latest", days],
    queryFn: () => getLatestTrend(days),
  });
}

