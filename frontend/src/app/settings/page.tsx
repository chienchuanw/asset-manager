"use client";

import { useState, useEffect } from "react";
import { useTranslations } from "next-intl";
import { AppLayout } from "@/components/layout/AppLayout";
import { useSettings, useUpdateSettings } from "@/hooks/useSettings";
import { useRefreshExchangeRate } from "@/hooks/useExchangeRates";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Loading } from "@/components/ui/loading";
import { toast } from "sonner";
import { Loader2, RefreshCw } from "lucide-react";
import type { AllocationSettings } from "@/types/analytics";

export default function SettingsPage() {
  const t = useTranslations("settings");
  const tCommon = useTranslations("common");
  const tAssetTypes = useTranslations("assetTypes");
  const { data: settings, isLoading } = useSettings();
  const updateSettingsMutation = useUpdateSettings();
  const refreshExchangeRateMutation = useRefreshExchangeRate();

  // 資產配置設定狀態
  const [allocationSettings, setAllocationSettings] =
    useState<AllocationSettings>({
      tw_stock: 40,
      us_stock: 40,
      crypto: 20,
      rebalance_threshold: 5,
    });

  // 匯率資訊狀態
  const [exchangeRateInfo, setExchangeRateInfo] = useState<{
    rate: number;
    updatedAt: string;
  } | null>(null);

  // 當設定載入完成時，更新狀態
  useEffect(() => {
    if (settings) {
      setAllocationSettings(settings.allocation);
    }
  }, [settings]);

  // 處理儲存
  const handleSave = async () => {
    // 驗證資產配置總和是否為 100%
    const total =
      allocationSettings.tw_stock +
      allocationSettings.us_stock +
      allocationSettings.crypto;
    if (Math.abs(total - 100) > 0.01) {
      toast.error(t("validationFailed"), {
        description: t("allocationSumError", { total: total.toFixed(2) }),
      });
      return;
    }

    try {
      await updateSettingsMutation.mutateAsync({
        allocation: allocationSettings,
      });

      toast.success(t("saveSuccess"), {
        description: t("settingsUpdated"),
      });
    } catch (error) {
      toast.error(t("saveFailed"), {
        description: error instanceof Error ? error.message : t("unknownError"),
      });
    }
  };

  // 處理重置
  const handleReset = () => {
    if (settings) {
      setAllocationSettings(settings.allocation);
      toast.info(t("resetDone"), {
        description: t("resetToLastSaved"),
      });
    }
  };

  // 處理更新匯率
  const handleRefreshExchangeRate = async () => {
    try {
      const result = await refreshExchangeRateMutation.mutateAsync();
      // 更新本地狀態以顯示最新資訊
      setExchangeRateInfo({
        rate: result.rate,
        updatedAt: result.updated_at,
      });
    } catch (error) {
      // 錯誤已在 mutation 的 onError 中處理
    }
  };

  if (isLoading) {
    return (
      <AppLayout title={t("title")} description={t("description")}>
        <main className="flex-1 p-4 md:p-6 bg-gray-50">
          <Loading variant="page" size="lg" text={t("loadingSettings")} />
        </main>
      </AppLayout>
    );
  }

  return (
    <AppLayout title={t("title")} description={t("description")}>
      {/* Main Content */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* 資產配置設定和匯率設定 - 並排顯示 */}
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            {/* 資產配置設定 */}
            <Card>
              <CardHeader>
                <CardTitle>{t("targetAllocation")}</CardTitle>
                <CardDescription>{t("targetAllocationDesc")}</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* 台股配置 */}
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <Label htmlFor="tw-stock">{tAssetTypes("twStock")}</Label>
                    <span className="text-sm text-muted-foreground">
                      {allocationSettings.tw_stock}%
                    </span>
                  </div>
                  <Input
                    id="tw-stock"
                    type="number"
                    min="0"
                    max="100"
                    step="0.1"
                    value={allocationSettings.tw_stock}
                    onChange={(e) =>
                      setAllocationSettings({
                        ...allocationSettings,
                        tw_stock: parseFloat(e.target.value) || 0,
                      })
                    }
                  />
                </div>

                {/* 美股配置 */}
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <Label htmlFor="us-stock">{tAssetTypes("usStock")}</Label>
                    <span className="text-sm text-muted-foreground">
                      {allocationSettings.us_stock}%
                    </span>
                  </div>
                  <Input
                    id="us-stock"
                    type="number"
                    min="0"
                    max="100"
                    step="0.1"
                    value={allocationSettings.us_stock}
                    onChange={(e) =>
                      setAllocationSettings({
                        ...allocationSettings,
                        us_stock: parseFloat(e.target.value) || 0,
                      })
                    }
                  />
                </div>

                {/* 加密貨幣配置 */}
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <Label htmlFor="crypto">{tAssetTypes("crypto")}</Label>
                    <span className="text-sm text-muted-foreground">
                      {allocationSettings.crypto}%
                    </span>
                  </div>
                  <Input
                    id="crypto"
                    type="number"
                    min="0"
                    max="100"
                    step="0.1"
                    value={allocationSettings.crypto}
                    onChange={(e) =>
                      setAllocationSettings({
                        ...allocationSettings,
                        crypto: parseFloat(e.target.value) || 0,
                      })
                    }
                  />
                </div>

                <Separator />

                {/* 總和顯示 */}
                <div className="flex items-center justify-between font-medium">
                  <span>{t("total")}</span>
                  <span
                    className={
                      Math.abs(
                        allocationSettings.tw_stock +
                          allocationSettings.us_stock +
                          allocationSettings.crypto -
                          100
                      ) < 0.01
                        ? "text-green-600"
                        : "text-red-600"
                    }
                  >
                    {(
                      allocationSettings.tw_stock +
                      allocationSettings.us_stock +
                      allocationSettings.crypto
                    ).toFixed(2)}
                    %
                  </span>
                </div>

                <Separator />

                {/* 再平衡閾值 */}
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <Label htmlFor="rebalance-threshold">
                      {t("rebalanceThreshold")}
                    </Label>
                    <span className="text-sm text-muted-foreground">
                      {allocationSettings.rebalance_threshold}%
                    </span>
                  </div>
                  <Input
                    id="rebalance-threshold"
                    type="number"
                    min="0"
                    max="50"
                    step="0.1"
                    value={allocationSettings.rebalance_threshold}
                    onChange={(e) =>
                      setAllocationSettings({
                        ...allocationSettings,
                        rebalance_threshold: parseFloat(e.target.value) || 0,
                      })
                    }
                  />
                  <p className="text-sm text-muted-foreground">
                    {t("rebalanceThresholdDesc")}
                  </p>
                </div>
              </CardContent>
            </Card>

            {/* 匯率設定 */}
            <Card>
              <CardHeader>
                <CardTitle>{t("exchangeRateSettings")}</CardTitle>
                <CardDescription>
                  {t("exchangeRateSettingsDesc")}
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* 當前匯率資訊 */}
                <div className="space-y-2">
                  <Label>{t("currentExchangeRate")}</Label>
                  <div className="flex items-center gap-2">
                    <div className="text-2xl font-bold">
                      {exchangeRateInfo
                        ? `USD/TWD: ${exchangeRateInfo.rate.toFixed(4)}`
                        : "USD/TWD: --"}
                    </div>
                    <Badge variant="secondary">ExchangeRate-API</Badge>
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {t("exchangeRateSource")}
                  </p>
                </div>

                <Separator />

                {/* 最後更新時間 */}
                <div className="space-y-2">
                  <Label>{t("lastUpdated")}</Label>
                  <div className="text-sm text-muted-foreground">
                    {exchangeRateInfo
                      ? new Date(exchangeRateInfo.updatedAt).toLocaleString(
                          "zh-TW",
                          {
                            year: "numeric",
                            month: "2-digit",
                            day: "2-digit",
                            hour: "2-digit",
                            minute: "2-digit",
                            second: "2-digit",
                          }
                        )
                      : t("notYetUpdated")}
                  </div>
                </div>

                <Separator />

                {/* 更新按鈕 */}
                <div className="flex justify-start">
                  <Button
                    onClick={handleRefreshExchangeRate}
                    disabled={refreshExchangeRateMutation.isPending}
                    variant="default"
                  >
                    {refreshExchangeRateMutation.isPending ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        {t("updating")}
                      </>
                    ) : (
                      <>
                        <RefreshCw className="mr-2 h-4 w-4" />
                        {t("updateExchangeRate")}
                      </>
                    )}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* 操作按鈕 */}
          <div className="flex justify-end gap-4">
            <Button
              variant="outline"
              onClick={handleReset}
              disabled={updateSettingsMutation.isPending}
            >
              {tCommon("reset")}
            </Button>
            <Button
              onClick={handleSave}
              disabled={updateSettingsMutation.isPending}
            >
              {updateSettingsMutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              {t("saveSettings")}
            </Button>
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
