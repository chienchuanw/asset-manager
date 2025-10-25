import {
  DailyPerformanceSnapshot,
  PerformanceTrendPoint,
  PerformanceTrendSummary,
} from "@/types/analytics";
import { apiClient } from "./client";

/**
 * 建立每日績效快照
 */
export async function createDailySnapshot(): Promise<DailyPerformanceSnapshot> {
  return await apiClient.post<DailyPerformanceSnapshot>(
    "/performance-trends/snapshot"
  );
}

/**
 * 取得日期範圍內的績效趨勢
 */
export async function getTrendByDateRange(
  startDate: string,
  endDate: string
): Promise<PerformanceTrendSummary> {
  return await apiClient.get<PerformanceTrendSummary>(
    "/performance-trends/range",
    {
      params: {
        start_date: startDate,
        end_date: endDate,
      },
    }
  );
}

/**
 * 取得最新的 N 天績效趨勢
 */
export async function getLatestTrend(
  days: number = 30
): Promise<PerformanceTrendPoint[]> {
  return await apiClient.get<PerformanceTrendPoint[]>(
    "/performance-trends/latest",
    {
      params: {
        days,
      },
    }
  );
}
