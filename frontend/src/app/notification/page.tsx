"use client";

import { useState, useEffect } from "react";
import { useTranslations } from "next-intl";
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
  const t = useTranslations("notification");
  const tCommon = useTranslations("common");
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
      toast.success(t("saveSuccess"), { description: t("settingsUpdated") });
    } catch (error) {
      toast.error(t("saveFailed"), {
        description:
          error instanceof Error ? error.message : tCommon("unknownError"),
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
      toast.info(t("resetDone"), {
        description: t("settingsResetToLastSaved"),
      });
    }
  };

  // 處理測試 Discord
  const handleTestDiscord = async () => {
    if (!discordSettings.webhook_url) {
      toast.error(t("testFailed"), { description: t("pleaseSetWebhookUrl") });
      return;
    }
    try {
      await testDiscordMutation.mutateAsync({
        message: t("testMessage"),
      });
    } catch (error) {
      // 錯誤已在 mutation 的 onError 中處理
    }
  };

  // 處理發送每日報告
  const handleSendDailyReport = async () => {
    if (!discordSettings.webhook_url) {
      toast.error(t("sendFailed"), { description: t("pleaseSetWebhookUrl") });
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
      toast.error(t("sendFailed"), { description: t("pleaseSetWebhookUrl") });
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
      toast.error(t("sendFailed"), { description: t("pleaseSetWebhookUrl") });
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
      <AppLayout title={t("title")} description={t("description")}>
        <main className="flex-1 p-4 md:p-6 bg-gray-50">
          <Loading variant="page" size="lg" text={tCommon("loading")} />
        </main>
      </AppLayout>
    );
  }

  return (
    <AppLayout title={t("title")} description={t("description")}>
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* Discord 設定和訂閱分期通知設定 - 並排顯示 */}
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            {/* Discord 設定 */}
            <Card>
              <CardHeader>
                <CardTitle>{t("discordSettings")}</CardTitle>
                <CardDescription>{t("discordSettingsDesc")}</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Discord Webhook URL */}
                <div className="space-y-2">
                  <Label htmlFor="webhook-url">{t("webhookUrl")}</Label>
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
                    {t("webhookUrlHint")}
                  </p>
                </div>

                {/* 啟用 Discord 通知 */}
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label htmlFor="discord-enabled">
                      {t("enableDailyReport")}
                    </Label>
                    <p className="text-sm text-muted-foreground">
                      {t("enableDailyReportDesc")}
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
                  <Label htmlFor="report-time">{t("reportTime")}</Label>
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
                    {t("reportTimeHint")}
                  </p>
                </div>

                <Separator />

                {/* 月度現金流報告設定 */}
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div className="space-y-0.5">
                      <Label htmlFor="monthly-report-enabled">
                        {t("enableMonthlyReport")}
                      </Label>
                      <p className="text-sm text-muted-foreground">
                        {t("enableMonthlyReportDesc")}
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
                      <Label htmlFor="monthly-report-day">
                        {t("monthlyReportDay")}
                      </Label>
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
                                {t("dayOfMonth", { day })}
                              </SelectItem>
                            )
                          )}
                        </SelectContent>
                      </Select>
                      <p className="text-sm text-muted-foreground">
                        {t("monthlyReportHint")}
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
                        {t("enableYearlyReport")}
                      </Label>
                      <p className="text-sm text-muted-foreground">
                        {t("enableYearlyReportDesc")}
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
                        <Label htmlFor="yearly-report-month">
                          {t("yearlyReportMonth")}
                        </Label>
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
                                  {t("month", { month })}
                                </SelectItem>
                              )
                            )}
                          </SelectContent>
                        </Select>
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="yearly-report-day">
                          {t("yearlyReportDay")}
                        </Label>
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
                                  {t("day", { day })}
                                </SelectItem>
                              )
                            )}
                          </SelectContent>
                        </Select>
                      </div>

                      <p className="text-sm text-muted-foreground">
                        {t("yearlyReportHint", {
                          month: discordSettings.yearly_report_month || 1,
                          day: discordSettings.yearly_report_day || 1,
                        })}
                      </p>
                    </div>
                  )}
                </div>

                <Separator />

                {/* Discord 測試按鈕 */}
                <div className="space-y-3">
                  <Label>{t("testDiscordFeatures")}</Label>
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
                          {t("sending")}
                        </>
                      ) : (
                        <>
                          <Send className="mr-2 h-4 w-4" />
                          {t("testMessageBtn")}
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
                          {t("sending")}
                        </>
                      ) : (
                        <>
                          <Send className="mr-2 h-4 w-4" />
                          {t("dailyReportBtn")}
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
                          {t("sending")}
                        </>
                      ) : (
                        <>
                          <Send className="mr-2 h-4 w-4" />
                          {t("monthlyReportBtn")}
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
                          {t("sending")}
                        </>
                      ) : (
                        <>
                          <Send className="mr-2 h-4 w-4" />
                          {t("yearlyReportBtn")}
                        </>
                      )}
                    </Button>
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {t("testDiscordHint")}
                  </p>
                </div>
              </CardContent>
            </Card>

            {/* 訂閱分期通知設定 */}
            <Card>
              <CardHeader>
                <CardTitle>{t("subscriptionSettings")}</CardTitle>
                <CardDescription>
                  {t("subscriptionSettingsDesc")}
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* 每日扣款通知 */}
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label htmlFor="notification-daily-billing">
                      {t("dailyBillingNotification")}
                    </Label>
                    <p className="text-sm text-muted-foreground">
                      {t("dailyBillingNotificationDesc")}
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
                      {t("subscriptionExpiryNotification")}
                    </Label>
                    <p className="text-sm text-muted-foreground">
                      {t("subscriptionExpiryNotificationDesc")}
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
                      {t("installmentCompletionNotification")}
                    </Label>
                    <p className="text-sm text-muted-foreground">
                      {t("installmentCompletionNotificationDesc")}
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
                  <Label htmlFor="notification-expiry-days">
                    {t("expiryReminderDays")}
                  </Label>
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
                    {t("expiryReminderDaysHint")}
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
              {tCommon("reset")}
            </Button>
            <Button
              onClick={handleSave}
              disabled={updateSettingsMutation.isPending}
            >
              {updateSettingsMutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              {tCommon("saveSettings")}
            </Button>
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
