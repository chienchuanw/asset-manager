/**
 * 資產再平衡頁面
 * 顯示當前資產配置偏差和再平衡建議
 */

"use client";

import { useTranslations } from "next-intl";
import { AppLayout } from "@/components/layout/AppLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Loading } from "@/components/ui/loading";
import { Badge } from "@/components/ui/badge";
import { RefreshCw } from "lucide-react";
import { useRebalanceCheck } from "@/hooks";
import { AllocationTable } from "@/components/rebalance/AllocationTable";
import { SuggestionCards } from "@/components/rebalance/SuggestionCards";

export default function RebalancePage() {
  const t = useTranslations("rebalance");
  const tCommon = useTranslations("common");
  const tErrors = useTranslations("errors");

  const { data, isLoading, error, refetch, isFetching } = useRebalanceCheck();

  // Loading 狀態
  if (isLoading) {
    return (
      <AppLayout title={t("title")} description={t("description")}>
        <div className="flex-1 p-4 md:p-6 bg-gray-50">
          <div className="flex items-center justify-center h-96">
            <Loading variant="page" size="lg" text={tCommon("loading")} />
          </div>
        </div>
      </AppLayout>
    );
  }

  // 錯誤狀態
  if (error) {
    return (
      <AppLayout title={t("title")} description={t("description")}>
        <div className="flex-1 p-4 md:p-6 bg-gray-50">
          <div className="flex items-center justify-center h-96">
            <Card className="w-full max-w-md">
              <CardHeader>
                <CardTitle className="text-red-600">
                  {tErrors("loadFailed")}
                </CardTitle>
                <CardDescription>{error.message}</CardDescription>
              </CardHeader>
              <CardContent>
                <Button
                  onClick={() => refetch()}
                  variant="outline"
                  className="w-full"
                >
                  <RefreshCw className="mr-2 h-4 w-4" />
                  {tCommon("reload")}
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout title={t("title")} description={t("description")}>
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* 摘要資訊 */}
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>{t("totalPortfolioValue")}</CardDescription>
                <CardTitle className="text-2xl tabular-nums">
                  NT$ {(data?.current_total ?? 0).toLocaleString("zh-TW")}
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>{t("threshold")}</CardDescription>
                <CardTitle className="text-2xl tabular-nums">
                  {data?.threshold ?? 0}%
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>{t("needsRebalance")}</CardDescription>
                <CardTitle className="text-2xl">
                  {data?.needs_rebalance ? (
                    <Badge variant="destructive">{t("needsRebalance")}</Badge>
                  ) : (
                    <Badge variant="secondary">{t("balanced")}</Badge>
                  )}
                </CardTitle>
              </CardHeader>
            </Card>
          </div>

          {/* 工具列 */}
          <div className="flex justify-end">
            <Button
              variant="outline"
              size="sm"
              onClick={() => refetch()}
              disabled={isFetching}
            >
              <RefreshCw
                className={`mr-2 h-4 w-4 ${isFetching ? "animate-spin" : ""}`}
              />
              {tCommon("refresh")}
            </Button>
          </div>

          {/* 配置偏差表格 */}
          <AllocationTable
            deviations={data?.deviations ?? []}
            threshold={data?.threshold ?? 0}
          />

          {/* 再平衡建議 */}
          <SuggestionCards suggestions={data?.suggestions ?? []} />
        </div>
      </div>
    </AppLayout>
  );
}
