"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useCashFlowSummary } from "@/hooks";
import { formatAmount } from "@/types/cash-flow";
import { TrendingUp, TrendingDown } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";

interface CashFlowSummaryCardProps {
  startDate: string;
  endDate: string;
  totalRecords: number;
  incomeRecords: number;
  expenseRecords: number;
}

/**
 * 現金流摘要卡片
 *
 * 顯示指定期間的收入、支出、淨現金流以及記錄統計
 * 採用 2 欄 3 列的佈局設計
 */
export function CashFlowSummaryCard({
  startDate,
  endDate,
  totalRecords,
  incomeRecords,
  expenseRecords,
}: CashFlowSummaryCardProps) {
  const { data: summary, isLoading } = useCashFlowSummary(startDate, endDate, {
    // 確保資料總是最新的
    staleTime: 0,
  });

  if (isLoading) {
    return (
      <div className="grid grid-cols-2 gap-6">
        {[1, 2, 3, 4, 5, 6].map((i) => (
          <Card key={i} className="hover:shadow-lg transition-shadow">
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
    <div className="grid grid-cols-2 gap-6">
      {/* 左欄第一列：總收入 */}
      <Card className="hover:shadow-lg transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">總收入</CardTitle>
          <TrendingUp className="h-4 w-4 text-green-600" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-green-600">
            ${formatAmount(summary.total_income)}
          </div>
          <p className="text-xs text-muted-foreground mt-1">期間內的所有收入</p>
        </CardContent>
      </Card>

      {/* 右欄第一列：總記錄數 */}
      <Card className="hover:shadow-lg transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">總記錄數</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-gray-600">{totalRecords}</div>
          <p className="text-xs text-muted-foreground mt-1">
            期間內的所有交易記錄
          </p>
        </CardContent>
      </Card>

      {/* 左欄第二列：總支出 */}
      <Card className="hover:shadow-lg transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">總支出</CardTitle>
          <TrendingDown className="h-4 w-4 text-red-600" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-red-600">
            ${formatAmount(summary.total_expense)}
          </div>
          <p className="text-xs text-muted-foreground mt-1">期間內的所有支出</p>
        </CardContent>
      </Card>

      {/* 右欄第二列：收入記錄 */}
      <Card className="hover:shadow-lg transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">收入記錄</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-gray-600">
            {incomeRecords}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            期間內的收入交易筆數
          </p>
        </CardContent>
      </Card>

      {/* 左欄第三列：淨現金流 */}
      <Card className="hover:shadow-lg transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">淨現金流</CardTitle>
        </CardHeader>
        <CardContent>
          <div
            className={`text-2xl font-bold ${
              netCashFlowIsPositive ? "text-green-600" : "text-red-600"
            }`}
          >
            {netCashFlowIsPositive ? "+" : ""}$
            {formatAmount(summary.net_cash_flow)}
          </div>
          <p className="text-xs text-muted-foreground mt-1">收入減去支出</p>
        </CardContent>
      </Card>

      {/* 右欄第三列：支出記錄 */}
      <Card className="hover:shadow-lg transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">支出記錄</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-gray-600">
            {expenseRecords}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            期間內的支出交易筆數
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
