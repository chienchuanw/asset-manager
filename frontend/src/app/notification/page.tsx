"use client";

import { useState, useEffect } from "react";
import { AppLayout } from "@/components/layout/AppLayout";
import { useSettings, useUpdateSettings } from "@/hooks/useSettings";
import {
  useTestDiscord,
  useSendDailyReport,
  useSendMonthlyReport,
  useSendYearlyReport,
} from "@/hooks/useDiscord";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
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
import type { DiscordSettings, NotificationSettings } from "@/types/analytics";

export default function NotificationPage() {
  const { data: settings, isLoading } = useSettings();
  const updateSettingsMutation = useUpdateSettings();
  const testDiscordMutation = useTestDiscord();
  const sendDailyReportMutation = useSendDailyReport();
  const sendMonthlyReportMutation = useSendMonthlyReport();
  const sendYearlyReportMutation = useSendYearlyReport();

  // Discord 設定狀態
  const [discordSettings, setDiscordSettings] = useState<DiscordSettings>({
    webhook_url: "",
    enabled: false,
    report_time: "09:00",
    monthly_report_enabled: false,
    monthly_report_day: 1,
    yearly_report_enabled: false,
    yearly_report_month: 1,
    yearly_report_day: 1,
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
      if (settings.notification) {
        setNotificationSettings(settings.notification);
      }
    }
  }, [settings]);

  // 處理儲存
  const handleSave = async () => {
    try {
      await updateSettingsMutation.mutateAsync({
        discord: discordSettings,
        notification: notificationSettings,
      });
      toast.success("儲存成功", { description: "通知設定已成功更新" });
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
      if (settings.notification) {
        setNotificationSettings(settings.notification);
      }
      toast.info("已重置", { description: "設定已重置為上次儲存的值" });
    }
  };

  // 處理測試 Discord
  const handleTestDiscord = async () => {
    if (!discordSettings.webhook_url) {
      toast.error("測試失敗", { description: "請先設定 Discord Webhook URL" });
      return;
    }
    try {
      await testDiscordMutation.mutateAsync({
        message: "（測試）資產管理系統的測試訊息!",
      });
    } catch (error) {
      // 錯誤已在 mutation 的 onError 中處理
    }
  };

  // 處理發送每日報告
  const handleSendDailyReport = async () => {
    if (!discordSettings.webhook_url) {
      toast.error("發送失敗", { description: "請先設定 Discord Webhook URL" });
      return;
    }
    try {
      await sendDailyReportMutation.mutateAsync();
    } catch (error) {
      // 錯誤已在 mutation 的 onError 中處理
    }
  };

  // 處理發送月度報告
  const handleSendMonthlyReport = async () => {
    if (!discordSettings.webhook_url) {
      toast.error("發送失敗", { description: "請先設定 Discord Webhook URL" });
      return;
    }
    const now = new Date();
    const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1);
    try {
      await sendMonthlyReportMutation.mutateAsync({
        year: lastMonth.getFullYear(),
        month: lastMonth.getMonth() + 1,
        webhook_url: discordSettings.webhook_url,
      });
    } catch (error) {
      // 錯誤已在 mutation 的 onError 中處理
    }
  };

  // 處理發送年度報告
  const handleSendYearlyReport = async () => {
    if (!discordSettings.webhook_url) {
      toast.error("發送失敗", { description: "請先設定 Discord Webhook URL" });
      return;
    }
    try {
      await sendYearlyReportMutation.mutateAsync({
        year: new Date().getFullYear() - 1,
        webhook_url: discordSettings.webhook_url,
      });
    } catch (error) {
      // 錯誤已在 mutation 的 onError 中處理
    }
  };

  if (isLoading) {
    return (
      <AppLayout title="通知管理" description="管理 Discord 通知和訂閱分期通知設定">
        <main className="flex-1 p-4 md:p-6 bg-gray-50">
          <Loading variant="page" size="lg" text="載入設定中..." />
        </main>
      </AppLayout>
    );
  }

  return (
    <AppLayout title="通知管理" description="管理 Discord 通知和訂閱分期通知設定">
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* Discord 設定和訂閱分期通知設定 - 並排顯示 */}
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

                {/* 月度現金流報告設定 */}
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div className="space-y-0.5">
                      <Label htmlFor="monthly-report-enabled">
                        啟用月度現金流報告
                      </Label>
                      <p className="text-sm text-muted-foreground">
                        每月自動發送現金流摘要報告到 Discord
                      </p>
                    </div>
                    <Switch
                      id="monthly-report-enabled"
                      checked={discordSettings.monthly_report_enabled}
                      onCheckedChange={(checked) =>
                        setDiscordSettings({
                          ...discordSettings,
                          monthly_report_enabled: checked,
                        })
                      }
                    />
                  </div>

                  {discordSettings.monthly_report_enabled && (
                    <div className="space-y-2 pl-4 border-l-2 border-muted">
                      <Label htmlFor="monthly-report-day">每月發送日期</Label>
                      <Select
                        value={String(discordSettings.monthly_report_day || 1)}
                        onValueChange={(value) =>
                          setDiscordSettings({
                            ...discordSettings,
                            monthly_report_day: parseInt(value),
                          })
                        }
                      >
                        <SelectTrigger
                          id="monthly-report-day"
                          className="w-full"
                        >
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {Array.from({ length: 10 }, (_, i) => i + 1).map(
                            (day) => (
                              <SelectItem key={day} value={String(day)}>
                                每月 {day} 號
                              </SelectItem>
                            )
                          )}
                        </SelectContent>
                      </Select>
                      <p className="text-sm text-muted-foreground">
                        報告將於每月指定日期的 09:00 發送
                      </p>
                    </div>
                  )}
                </div>

                <Separator />

                {/* 年度現金流報告設定 */}
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div className="space-y-0.5">
                      <Label htmlFor="yearly-report-enabled">
                        啟用年度現金流報告
                      </Label>
                      <p className="text-sm text-muted-foreground">
                        每年自動發送現金流摘要報告到 Discord
                      </p>
                    </div>
                    <Switch
                      id="yearly-report-enabled"
                      checked={discordSettings.yearly_report_enabled}
                      onCheckedChange={(checked) =>
                        setDiscordSettings({
                          ...discordSettings,
                          yearly_report_enabled: checked,
                        })
                      }
                    />
                  </div>

                  {discordSettings.yearly_report_enabled && (
                    <div className="space-y-4 pl-4 border-l-2 border-muted">
                      <div className="space-y-2">
                        <Label htmlFor="yearly-report-month">發送月份</Label>
                        <Select
                          value={String(
                            discordSettings.yearly_report_month || 1
                          )}
                          onValueChange={(value) =>
                            setDiscordSettings({
                              ...discordSettings,
                              yearly_report_month: parseInt(value),
                            })
                          }
                        >
                          <SelectTrigger
                            id="yearly-report-month"
                            className="w-full"
                          >
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            {Array.from({ length: 12 }, (_, i) => i + 1).map(
                              (month) => (
                                <SelectItem key={month} value={String(month)}>
                                  {month} 月
                                </SelectItem>
                              )
                            )}
                          </SelectContent>
                        </Select>
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="yearly-report-day">發送日期</Label>
                        <Select
                          value={String(discordSettings.yearly_report_day || 1)}
                          onValueChange={(value) =>
                            setDiscordSettings({
                              ...discordSettings,
                              yearly_report_day: parseInt(value),
                            })
                          }
                        >
                          <SelectTrigger
                            id="yearly-report-day"
                            className="w-full"
                          >
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            {Array.from({ length: 10 }, (_, i) => i + 1).map(
                              (day) => (
                                <SelectItem key={day} value={String(day)}>
                                  {day} 號
                                </SelectItem>
                              )
                            )}
                          </SelectContent>
                        </Select>
                      </div>

                      <p className="text-sm text-muted-foreground">
                        報告將於每年{discordSettings.yearly_report_month || 1}月
                        {discordSettings.yearly_report_day || 1}號的 09:00 發送
                      </p>
                    </div>
                  )}
                </div>

                <Separator />

                {/* Discord 測試按鈕 */}
                <div className="space-y-3">
                  <Label>測試 Discord 功能</Label>
                  <div className="grid grid-cols-2 gap-3">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleTestDiscord}
                      disabled={
                        !discordSettings.webhook_url ||
                        testDiscordMutation.isPending
                      }
                    >
                      {testDiscordMutation.isPending ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          發送中...
                        </>
                      ) : (
                        <>
                          <Send className="mr-2 h-4 w-4" />
                          測試訊息
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
                    >
                      {sendDailyReportMutation.isPending ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          發送中...
                        </>
                      ) : (
                        <>
                          <Send className="mr-2 h-4 w-4" />
                          每日報告
                        </>
                      )}
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleSendMonthlyReport}
                      disabled={
                        !discordSettings.webhook_url ||
                        sendMonthlyReportMutation.isPending
                      }
                    >
                      {sendMonthlyReportMutation.isPending ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          發送中...
                        </>
                      ) : (
                        <>
                          <Send className="mr-2 h-4 w-4" />
                          月度報告
                        </>
                      )}
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleSendYearlyReport}
                      disabled={
                        !discordSettings.webhook_url ||
                        sendYearlyReportMutation.isPending
                      }
                    >
                      {sendYearlyReportMutation.isPending ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          發送中...
                        </>
                      ) : (
                        <>
                          <Send className="mr-2 h-4 w-4" />
                          年度報告
                        </>
                      )}
                    </Button>
                  </div>
                  <p className="text-sm text-muted-foreground">
                    測試 Discord Webhook
                    是否正常運作（月度/年度報告將發送上一期間的資料）
                  </p>
                </div>
              </CardContent>
            </Card>

            {/* 訂閱分期通知設定 */}
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
          </div>

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
      </div>
    </AppLayout>
  );
}
