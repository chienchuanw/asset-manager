/**
 * Dashboard 主頁面
 * 整合所有 Dashboard 元件,顯示資產概況
 */

"use client";

import { useMemo } from "react";
import { useTranslations } from "next-intl";
import { AppLayout } from "@/components/layout/AppLayout";
import { StatCard } from "@/components/dashboard/StatCard";
import { AssetTrendChart } from "@/components/dashboard/AssetTrendChart";
import { HoldingsByAssetType } from "@/components/dashboard/HoldingsByAssetType";
import { AssetAllocationChart } from "@/components/dashboard/AssetAllocationChart";
import { RecentTransactions } from "@/components/dashboard/RecentTransactions";
import { RecurringStatsCard } from "@/components/dashboard/RecurringStatsCard";
import {
  useHoldings,
  useTransactions,
  useSubscriptions,
  useInstallments,
} from "@/hooks";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, AlertTriangle } from "lucide-react";

export default function DashboardPage() {
  const t = useTranslations("dashboard");
  const tErrors = useTranslations("errors");
  // 取得持倉資料（包含 warnings）
  const {
    data: holdingsResponse,
    isLoading: holdingsLoading,
    error: holdingsError,
  } = useHoldings();

  // 解構 holdings
  const holdings = holdingsResponse?.data || [];

  // 取得交易資料
  const {
    data: transactions,
    isLoading: transactionsLoading,
    error: transactionsError,
  } = useTransactions();

  // 取得訂閱和分期資料
  const { data: subscriptions, isLoading: subscriptionsLoading } =
    useSubscriptions();
  const { data: installments, isLoading: installmentsLoading } =
    useInstallments();

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

  // 檢查是否有過期的價格資料
  const stalePriceInfo = useMemo(() => {
    if (!holdings) return null;

    const staleHoldings = holdings.filter((h) => h.is_price_stale);
    if (staleHoldings.length === 0) return null;

    // 取得過期原因（通常所有過期的持倉原因相同）
    const reason = staleHoldings[0].price_stale_reason || "API 請求失敗";
    const symbols = staleHoldings.map((h) => h.symbol).join(", ");

    return {
      count: staleHoldings.length,
      reason,
      symbols,
    };
  }, [holdings]);

  // Loading 狀態
  if (holdingsLoading || transactionsLoading) {
    return (
      <AppLayout title={t("title")} description={t("description")}>
        <div className="flex-1 p-4 md:p-6 bg-gray-50">
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
        </div>
      </AppLayout>
    );
  }

  // Error 狀態
  if (holdingsError || transactionsError) {
    return (
      <AppLayout title={t("title")} description={t("description")}>
        <div className="flex-1 p-4 md:p-6 bg-gray-50">
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-2 text-red-600">
                <AlertCircle className="h-5 w-5" />
                <p>
                  {tErrors("loadFailed")}:{" "}
                  {holdingsError?.message || transactionsError?.message}
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout title={t("title")} description={t("description")}>
      {/* 內容區域 */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="@container/main flex flex-1 flex-col gap-4 md:gap-6">
          {/* 過期價格警告 */}
          {stalePriceInfo && (
            <Alert variant="default" className="border-amber-200 bg-amber-50">
              <AlertTriangle className="h-4 w-4 text-amber-600" />
              <AlertTitle className="text-amber-900">
                {t("stalePriceWarning")}
              </AlertTitle>
              <AlertDescription className="text-amber-800">
                {t("stalePriceDescription")}
                <br />
                <span className="text-sm">
                  {t("stalePriceReason")}: {stalePriceInfo.reason}
                  <br />
                  {t("stalePriceAffected")}: {stalePriceInfo.symbols}
                </span>
              </AlertDescription>
            </Alert>
          )}

          {/* 統計卡片區 - 響應式網格 */}
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <StatCard
              title={t("totalValue")}
              value={`TWD ${stats.totalValue.toLocaleString("zh-TW", {
                maximumFractionDigits: 0,
              })}`}
              change={stats.totalPLPct}
              description={t("totalValueDesc")}
            />
            <StatCard
              title={t("totalCost")}
              value={`TWD ${stats.totalCost.toLocaleString("zh-TW", {
                maximumFractionDigits: 0,
              })}`}
              change={0}
              description={t("totalCostDesc")}
            />
            <StatCard
              title={t("unrealizedPL")}
              value={`TWD ${stats.totalPL.toLocaleString("zh-TW", {
                maximumFractionDigits: 0,
              })}`}
              change={stats.totalPLPct}
              description={`${
                stats.totalPLPct >= 0 ? "+" : ""
              }${stats.totalPLPct.toFixed(1)}%`}
            />
            <StatCard
              title={t("holdingsCount")}
              value={stats.holdingsCount.toString()}
              change={0}
              description={t("holdingsCountDesc")}
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
            {/* 左側：持倉明細(依資產類別分組) */}
            <div className="lg:col-span-4">
              <HoldingsByAssetType holdings={holdings || []} />
            </div>

            {/* 右側：近期交易和訂閱分期 */}
            <div className="lg:col-span-3 space-y-4">
              <RecentTransactions transactions={recentTransactions} />
              <RecurringStatsCard
                subscriptions={subscriptions}
                installments={installments}
                isLoading={subscriptionsLoading || installmentsLoading}
              />
            </div>
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
