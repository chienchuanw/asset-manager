"use client";

import { useState, useEffect } from "react";
import { AppLayout } from "@/components/layout/AppLayout";
import { useSettings, useUpdateSettings } from "@/hooks/useSettings";
import { useTestDiscord, useSendDailyReport } from "@/hooks/useDiscord";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
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
import { Loader2, Send } from "lucide-react";
import type {
  DiscordSettings,
  AllocationSettings,
  NotificationSettings,
} from "@/types/analytics";

export default function SettingsPage() {
  const { data: settings, isLoading } = useSettings();
  const updateSettingsMutation = useUpdateSettings();
  const testDiscordMutation = useTestDiscord();
  const sendDailyReportMutation = useSendDailyReport();

  // Discord 設定狀態
  const [discordSettings, setDiscordSettings] = useState<DiscordSettings>({
    webhook_url: "",
    enabled: false,
    report_time: "09:00",
  });

  // 資產配置設定狀態
  const [allocationSettings, setAllocationSettings] =
    useState<AllocationSettings>({
      tw_stock: 40,
      us_stock: 40,
      crypto: 20,
      rebalance_threshold: 5,
    });

  // 通知設定狀態
  const [notificationSettings, setNotificationSettings] =
    useState<NotificationSettings>({
      daily_billing: true,
      subscription_expiry: true,
      installment_completion: true,
      expiry_days: 7,
    });

  // 當設定載入完成時，更新狀態
  useEffect(() => {
    if (settings) {
      setDiscordSettings(settings.discord);
      setAllocationSettings(settings.allocation);
      // 如果後端沒有返回 notification 設定，使用預設值
      if (settings.notification) {
        setNotificationSettings(settings.notification);
      }
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
      toast.error("驗證失敗", {
        description: `資產配置總和必須為 100%，目前為 ${total.toFixed(2)}%`,
      });
      return;
    }

    try {
      await updateSettingsMutation.mutateAsync({
        discord: discordSettings,
        allocation: allocationSettings,
        notification: notificationSettings,
      });

      toast.success("儲存成功", {
        description: "設定已成功更新",
      });
    } catch (error) {
      toast.error("儲存失敗", {
        description: error instanceof Error ? error.message : "未知錯誤",
      });
    }
  };

  // 處理重置
  const handleReset = () => {
    if (settings) {
      setDiscordSettings(settings.discord);
      setAllocationSettings(settings.allocation);
      // 如果後端沒有返回 notification 設定，使用預設值
      if (settings.notification) {
        setNotificationSettings(settings.notification);
      }
      toast.info("已重置", {
        description: "設定已重置為上次儲存的值",
      });
    }
  };

  // 處理測試 Discord
  const handleTestDiscord = async () => {
    // 檢查 Webhook URL 是否已設定
    if (!discordSettings.webhook_url) {
      toast.error("測試失敗", {
        description: "請先設定 Discord Webhook URL",
      });
      return;
    }

    try {
      await testDiscordMutation.mutateAsync({
        message: "📢（測試）資產管理系統的測試訊息！",
      });
    } catch (error) {
      // 錯誤已在 mutation 的 onError 中處理
    }
  };

  // 處理發送每日報告
  const handleSendDailyReport = async () => {
    // 檢查 Webhook URL 是否已設定
    if (!discordSettings.webhook_url) {
      toast.error("發送失敗", {
        description: "請先設定 Discord Webhook URL",
      });
      return;
    }

    try {
      await sendDailyReportMutation.mutateAsync();
    } catch (error) {
      // 錯誤已在 mutation 的 onError 中處理
    }
  };

  if (isLoading) {
    return (
      <AppLayout title="設定" description="管理系統設定和偏好">
        <main className="flex-1 p-4 md:p-6 bg-gray-50">
          <Loading variant="page" size="lg" text="載入設定中..." />
        </main>
      </AppLayout>
    );
  }

  return (
    <AppLayout title="設定" description="管理系統設定和偏好">
      {/* Main Content */}
      <main className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* Discord 和資產配置設定 - 並排顯示 */}
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            {/* Discord 設定 */}
            <Card>
              <CardHeader>
                <CardTitle>Discord 通知設定</CardTitle>
                <CardDescription>
                  設定 Discord Webhook 以接收每日報告
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Discord Webhook URL */}
                <div className="space-y-2">
                  <Label htmlFor="webhook-url">Webhook URL</Label>
                  <Input
                    id="webhook-url"
                    type="url"
                    placeholder="https://discord.com/api/webhooks/..."
                    value={discordSettings.webhook_url}
                    onChange={(e) =>
                      setDiscordSettings({
                        ...discordSettings,
                        webhook_url: e.target.value,
                      })
                    }
                  />
                  <p className="text-sm text-muted-foreground">
                    在 Discord 伺服器設定中建立 Webhook 並貼上 URL
                  </p>
                </div>

                {/* 啟用 Discord 通知 */}
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label htmlFor="discord-enabled">啟用每日報告</Label>
                    <p className="text-sm text-muted-foreground">
                      每天自動發送投資組合報告到 Discord
                    </p>
                  </div>
                  <Switch
                    id="discord-enabled"
                    checked={discordSettings.enabled}
                    onCheckedChange={(checked) =>
                      setDiscordSettings({
                        ...discordSettings,
                        enabled: checked,
                      })
                    }
                  />
                </div>

                {/* 報告時間 */}
                <div className="space-y-2">
                  <Label htmlFor="report-time">報告時間</Label>
                  <Input
                    id="report-time"
                    type="time"
                    value={discordSettings.report_time}
                    onChange={(e) =>
                      setDiscordSettings({
                        ...discordSettings,
                        report_time: e.target.value,
                      })
                    }
                  />
                  <p className="text-sm text-muted-foreground">
                    每日報告發送時間（24 小時制）
                  </p>
                </div>

                <Separator />

                {/* Discord 測試按鈕 */}
                <div className="space-y-3">
                  <Label>測試 Discord 功能</Label>
                  <div className="flex gap-3">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleTestDiscord}
                      disabled={
                        !discordSettings.webhook_url ||
                        testDiscordMutation.isPending
                      }
                      className="flex-1"
                    >
                      {testDiscordMutation.isPending ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          發送中...
                        </>
                      ) : (
                        <>
                          <Send className="mr-2 h-4 w-4" />
                          發送測試訊息
                        </>
                      )}
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleSendDailyReport}
                      disabled={
                        !discordSettings.webhook_url ||
                        sendDailyReportMutation.isPending
                      }
                      className="flex-1"
                    >
                      {sendDailyReportMutation.isPending ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          發送中...
                        </>
                      ) : (
                        <>
                          <Send className="mr-2 h-4 w-4" />
                          發送每日報告
                        </>
                      )}
                    </Button>
                  </div>
                  <p className="text-sm text-muted-foreground">
                    測試 Discord Webhook 是否正常運作
                  </p>
                </div>
              </CardContent>
            </Card>

            {/* 資產配置設定 */}
            <Card>
              <CardHeader>
                <CardTitle>目標資產配置</CardTitle>
                <CardDescription>
                  設定各資產類型的目標配置百分比（總和必須為 100%）
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* 台股配置 */}
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <Label htmlFor="tw-stock">台股</Label>
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
                    <Label htmlFor="us-stock">美股</Label>
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
                    <Label htmlFor="crypto">加密貨幣</Label>
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
                  <span>總和</span>
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
                    <Label htmlFor="rebalance-threshold">再平衡閾值</Label>
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
                    當實際配置與目標配置偏差超過此百分比時，系統會發出再平衡提醒
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* 通知設定 - 獨立一行 */}
          <Card>
            <CardHeader>
              <CardTitle>訂閱分期通知設定</CardTitle>
              <CardDescription>設定訂閱和分期相關的通知選項</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* 每日扣款通知 */}
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label htmlFor="notification-daily-billing">
                    每日扣款通知
                  </Label>
                  <p className="text-sm text-muted-foreground">
                    每天自動扣款後發送通知到 Discord
                  </p>
                </div>
                <Switch
                  id="notification-daily-billing"
                  checked={notificationSettings.daily_billing}
                  onCheckedChange={(checked) =>
                    setNotificationSettings({
                      ...notificationSettings,
                      daily_billing: checked,
                    })
                  }
                />
              </div>

              <Separator />

              {/* 訂閱到期通知 */}
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label htmlFor="notification-subscription-expiry">
                    訂閱到期通知
                  </Label>
                  <p className="text-sm text-muted-foreground">
                    訂閱即將到期時發送提醒通知
                  </p>
                </div>
                <Switch
                  id="notification-subscription-expiry"
                  checked={notificationSettings.subscription_expiry}
                  onCheckedChange={(checked) =>
                    setNotificationSettings({
                      ...notificationSettings,
                      subscription_expiry: checked,
                    })
                  }
                />
              </div>

              <Separator />

              {/* 分期完成通知 */}
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label htmlFor="notification-installment-completion">
                    分期完成通知
                  </Label>
                  <p className="text-sm text-muted-foreground">
                    分期即將完成時發送提醒通知
                  </p>
                </div>
                <Switch
                  id="notification-installment-completion"
                  checked={notificationSettings.installment_completion}
                  onCheckedChange={(checked) =>
                    setNotificationSettings({
                      ...notificationSettings,
                      installment_completion: checked,
                    })
                  }
                />
              </div>

              <Separator />

              {/* 到期提醒天數 */}
              <div className="space-y-2">
                <Label htmlFor="notification-expiry-days">到期提醒天數</Label>
                <Input
                  id="notification-expiry-days"
                  type="number"
                  min="1"
                  max="30"
                  value={notificationSettings.expiry_days}
                  onChange={(e) =>
                    setNotificationSettings({
                      ...notificationSettings,
                      expiry_days: parseInt(e.target.value) || 7,
                    })
                  }
                />
                <p className="text-sm text-muted-foreground">
                  提前幾天發送到期提醒（1-30 天）
                </p>
              </div>
            </CardContent>
          </Card>

          {/* 操作按鈕 */}
          <div className="flex justify-end gap-4">
            <Button
              variant="outline"
              onClick={handleReset}
              disabled={updateSettingsMutation.isPending}
            >
              重置
            </Button>
            <Button
              onClick={handleSave}
              disabled={updateSettingsMutation.isPending}
            >
              {updateSettingsMutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              儲存設定
            </Button>
          </div>
        </div>
      </main>
    </AppLayout>
  );
}
