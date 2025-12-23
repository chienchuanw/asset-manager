"use client";

import { useTranslations } from "next-intl";
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
 * 採用 3 欄 2 列的佈局設計
 */
export function CashFlowSummaryCard({
  startDate,
  endDate,
  totalRecords,
  incomeRecords,
  expenseRecords,
}: CashFlowSummaryCardProps) {
  const t = useTranslations("cashFlows");

  const { data: summary, isLoading } = useCashFlowSummary(startDate, endDate, {
    // 確保資料總是最新的
    staleTime: 0,
  });

  if (isLoading) {
    return (
      <div className="grid grid-cols-3 gap-6">
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
    <div className="grid grid-cols-3 gap-6">
      {/* 第一列第一欄：總收入 */}
      <Card className="hover:shadow-lg transition-shadow">
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
            {t("totalIncomeDesc")}
          </p>
        </CardContent>
      </Card>

      {/* 第一列第二欄：總支出 */}
      <Card className="hover:shadow-lg transition-shadow">
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
            {t("totalExpenseDesc")}
          </p>
        </CardContent>
      </Card>

      {/* 第一列第三欄：淨現金流 */}
      <Card className="hover:shadow-lg transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            {t("netCashFlow")}
          </CardTitle>
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
          <p className="text-xs text-muted-foreground mt-1">
            {t("netCashFlowDesc")}
          </p>
        </CardContent>
      </Card>

      {/* 第二列第一欄：收入記錄 */}
      <Card className="hover:shadow-lg transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            {t("incomeRecords")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-gray-600">
            {incomeRecords}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            {t("incomeRecordsDesc")}
          </p>
        </CardContent>
      </Card>

      {/* 第二列第二欄：支出記錄 */}
      <Card className="hover:shadow-lg transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            {t("expenseRecords")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-gray-600">
            {expenseRecords}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            {t("expenseRecordsDesc")}
          </p>
        </CardContent>
      </Card>

      {/* 第二列第三欄：總記錄數 */}
      <Card className="hover:shadow-lg transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            {t("totalRecords")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-gray-600">{totalRecords}</div>
          <p className="text-xs text-muted-foreground mt-1">
            {t("totalRecordsDesc")}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
