"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useCashFlowSummary } from "@/hooks";
import { formatAmount } from "@/types/cash-flow";
import { TrendingUp, TrendingDown } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";

interface DailySummaryCardsProps {
  date: string; // ISO 8601 格式 (YYYY-MM-DD)
}

/**
 * 「今日」專用的摘要卡片元件
 *
 * 只顯示「總收入」和「總支出」兩張卡片
 */
export function DailySummaryCards({ date }: DailySummaryCardsProps) {
  const t = useTranslations("common");

  // 取得該日期的現金流摘要（startDate 和 endDate 都是同一天）
  const { data: summary, isLoading } = useCashFlowSummary(date, date, {
    // 確保資料總是最新的（不使用快取）
    staleTime: 0,
  });

  if (isLoading) {
    return (
      <div className="grid grid-cols-2 gap-4">
        {[1, 2].map((i) => (
          <Card key={i}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <Skeleton className="h-4 w-20" />
              <Skeleton className="h-4 w-4 rounded-full" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-32 mb-2" />
              <Skeleton className="h-3 w-24" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (!summary) {
    return null;
  }

  return (
    <div className="grid grid-cols-2 gap-4">
      {/* 總收入 */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            {t("totalIncome")}
          </CardTitle>
          <TrendingUp className="h-4 w-4 text-green-600" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-green-600">
            ${formatAmount(summary.total_income)}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            {t("dailyIncomeDesc")}
          </p>
        </CardContent>
      </Card>

      {/* 總支出 */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            {t("totalExpense")}
          </CardTitle>
          <TrendingDown className="h-4 w-4 text-red-600" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-red-600">
            ${formatAmount(summary.total_expense)}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            {t("dailyExpenseDesc")}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
