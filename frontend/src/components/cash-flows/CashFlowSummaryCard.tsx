"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useCashFlowSummary } from "@/hooks";
import { formatAmount } from "@/types/cash-flow";
import { TrendingUp, TrendingDown, Wallet, Loader2 } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";

interface CashFlowSummaryCardProps {
  startDate: string;
  endDate: string;
}

/**
 * 現金流摘要卡片
 *
 * 顯示指定期間的收入、支出和淨現金流
 */
export function CashFlowSummaryCard({
  startDate,
  endDate,
}: CashFlowSummaryCardProps) {
  const { data: summary, isLoading } = useCashFlowSummary(startDate, endDate);

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-3">
        {[1, 2, 3].map((i) => (
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

  const netCashFlowIsPositive = summary.net_cash_flow >= 0;

  return (
    <div className="grid gap-4 md:grid-cols-3">
      {/* 總收入 */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">總收入</CardTitle>
          <TrendingUp className="h-4 w-4 text-green-600" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-green-600">
            ${formatAmount(summary.total_income)}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            期間內的所有收入
          </p>
        </CardContent>
      </Card>

      {/* 總支出 */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">總支出</CardTitle>
          <TrendingDown className="h-4 w-4 text-red-600" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-red-600">
            ${formatAmount(summary.total_expense)}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            期間內的所有支出
          </p>
        </CardContent>
      </Card>

      {/* 淨現金流 */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">淨現金流</CardTitle>
          <Wallet
            className={`h-4 w-4 ${
              netCashFlowIsPositive ? "text-green-600" : "text-red-600"
            }`}
          />
        </CardHeader>
        <CardContent>
          <div
            className={`text-2xl font-bold ${
              netCashFlowIsPositive ? "text-green-600" : "text-red-600"
            }`}
          >
            {netCashFlowIsPositive ? "+" : ""}${formatAmount(summary.net_cash_flow)}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            收入減去支出
          </p>
        </CardContent>
      </Card>
    </div>
  );
}

