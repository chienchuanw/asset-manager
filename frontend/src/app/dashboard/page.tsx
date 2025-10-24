/**
 * Dashboard 主頁面
 * 整合所有 Dashboard 元件,顯示資產概況
 */

"use client";

import { useMemo } from "react";
import { AppLayout } from "@/components/layout/AppLayout";
import { StatCard } from "@/components/dashboard/StatCard";
import { AssetTrendChart } from "@/components/dashboard/AssetTrendChart";
import { HoldingsTable } from "@/components/dashboard/HoldingsTable";
import { AssetAllocationChart } from "@/components/dashboard/AssetAllocationChart";
import { RecentTransactions } from "@/components/dashboard/RecentTransactions";
import { useHoldings, useTransactions } from "@/hooks";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { AlertCircle } from "lucide-react";

export default function DashboardPage() {
  // 取得持倉資料
  const {
    data: holdings,
    isLoading: holdingsLoading,
    error: holdingsError,
  } = useHoldings();

  // 取得交易資料
  const {
    data: transactions,
    isLoading: transactionsLoading,
    error: transactionsError,
  } = useTransactions();

  // 計算統計資料
  const stats = useMemo(() => {
    if (!holdings || !transactions) {
      return {
        totalValue: 0,
        totalCost: 0,
        totalPL: 0,
        totalPLPct: 0,
        holdingsCount: 0,
      };
    }

    const totalValue = holdings.reduce((sum, h) => sum + h.market_value, 0);
    const totalCost = holdings.reduce((sum, h) => sum + h.total_cost, 0);
    const totalPL = totalValue - totalCost;
    const totalPLPct = totalCost > 0 ? (totalPL / totalCost) * 100 : 0;

    return {
      totalValue,
      totalCost,
      totalPL,
      totalPLPct,
      holdingsCount: holdings.length,
    };
  }, [holdings, transactions]);

  // 計算資產配置資料
  const assetAllocation = useMemo(() => {
    if (!holdings) return [];

    const allocation = holdings.reduce((acc, holding) => {
      const existing = acc.find((item) => item.name === holding.asset_type);
      if (existing) {
        existing.value += holding.market_value;
      } else {
        acc.push({
          name: holding.asset_type,
          value: holding.market_value,
        });
      }
      return acc;
    }, [] as { name: string; value: number }[]);

    return allocation;
  }, [holdings]);

  // 取得最近 5 筆交易
  const recentTransactions = useMemo(() => {
    if (!transactions) return [];
    return transactions.slice(0, 5);
  }, [transactions]);

  // Loading 狀態
  if (holdingsLoading || transactionsLoading) {
    return (
      <AppLayout>
        <main className="flex-1 p-4 md:p-6 bg-gray-50">
          <div className="@container/main flex flex-1 flex-col gap-4 md:gap-6">
            {/* 統計卡片 Loading */}
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
              {Array.from({ length: 4 }).map((_, i) => (
                <Card key={i}>
                  <CardHeader className="pb-2">
                    <Skeleton className="h-4 w-24" />
                  </CardHeader>
                  <CardContent>
                    <Skeleton className="h-8 w-32 mb-2" />
                    <Skeleton className="h-3 w-20" />
                  </CardContent>
                </Card>
              ))}
            </div>

            {/* 圖表 Loading */}
            <div className="grid grid-cols-1 gap-4 lg:grid-cols-7 md:gap-6">
              <div className="lg:col-span-4">
                <Card>
                  <CardHeader>
                    <Skeleton className="h-6 w-32" />
                  </CardHeader>
                  <CardContent>
                    <Skeleton className="h-64 w-full" />
                  </CardContent>
                </Card>
              </div>
              <div className="lg:col-span-3">
                <Card>
                  <CardHeader>
                    <Skeleton className="h-6 w-32" />
                  </CardHeader>
                  <CardContent>
                    <Skeleton className="h-64 w-full" />
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        </main>
      </AppLayout>
    );
  }

  // Error 狀態
  if (holdingsError || transactionsError) {
    return (
      <AppLayout>
        <main className="flex-1 p-4 md:p-6 bg-gray-50">
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-2 text-red-600">
                <AlertCircle className="h-5 w-5" />
                <p>
                  載入資料失敗：
                  {holdingsError?.message || transactionsError?.message}
                </p>
              </div>
            </CardContent>
          </Card>
        </main>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      {/* 內容區域 */}
      <main className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="@container/main flex flex-1 flex-col gap-4 md:gap-6">
          {/* 統計卡片區 - 響應式網格 */}
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <StatCard
              title="總資產價值"
              value={`TWD ${stats.totalValue.toLocaleString("zh-TW", {
                maximumFractionDigits: 0,
              })}`}
              change={stats.totalPLPct}
              description="所有持倉市值總和"
            />
            <StatCard
              title="總成本"
              value={`TWD ${stats.totalCost.toLocaleString("zh-TW", {
                maximumFractionDigits: 0,
              })}`}
              change={0}
              description="所有持倉成本總和"
            />
            <StatCard
              title="未實現損益"
              value={`TWD ${stats.totalPL.toLocaleString("zh-TW", {
                maximumFractionDigits: 0,
              })}`}
              change={stats.totalPLPct}
              description={`${
                stats.totalPLPct >= 0 ? "+" : ""
              }${stats.totalPLPct.toFixed(1)}%`}
            />
            <StatCard
              title="持倉數量"
              value={stats.holdingsCount.toString()}
              change={0}
              description="目前持有標的數量"
            />
          </div>

          {/* 主要內容區 - 響應式佈局 */}
          <div className="grid grid-cols-1 gap-4 lg:grid-cols-7 md:gap-6">
            {/* 左側：資產趨勢圖表 */}
            <div className="lg:col-span-4">
              <AssetTrendChart />
            </div>

            {/* 右側：資產配置圓餅圖 */}
            <div className="lg:col-span-3">
              <AssetAllocationChart data={assetAllocation} />
            </div>
          </div>

          {/* 底部區域 - 響應式佈局 */}
          <div className="grid grid-cols-1 gap-4 lg:grid-cols-7 md:gap-6">
            {/* 左側：持倉明細表格 */}
            <div className="lg:col-span-4">
              <HoldingsTable holdings={holdings || []} />
            </div>

            {/* 右側：近期交易 */}
            <div className="lg:col-span-3">
              <RecentTransactions transactions={recentTransactions} />
            </div>
          </div>
        </div>
      </main>
    </AppLayout>
  );
}
