/**
 * 訂閱分期統計卡片元件
 * 顯示訂閱和分期的統計資訊
 */

"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import {
  CalendarIcon,
  CreditCardIcon,
  TrendingUpIcon,
  AlertCircleIcon,
} from "lucide-react";
import type { Subscription } from "@/types/subscription";
import type { Installment } from "@/types/installment";

interface RecurringStatsCardProps {
  subscriptions?: Subscription[];
  installments?: Installment[];
  isLoading?: boolean;
}

/**
 * 計算訂閱的每月總成本
 */
function calculateMonthlySubscriptionCost(
  subscriptions: Subscription[]
): number {
  return subscriptions.reduce((total, sub) => {
    if (sub.status !== "active") return total;

    // 根據計費週期換算成每月成本
    switch (sub.billing_cycle) {
      case "monthly":
        return total + sub.amount;
      case "quarterly":
        return total + sub.amount / 3;
      case "yearly":
        return total + sub.amount / 12;
      default:
        return total;
    }
  }, 0);
}

/**
 * 計算分期的每月總付款
 */
function calculateMonthlyInstallmentPayment(
  installments: Installment[]
): number {
  return installments.reduce((total, inst) => {
    if (inst.status !== "active") return total;
    return total + (inst.installment_amount ?? 0);
  }, 0);
}

export function RecurringStatsCard({
  subscriptions = [],
  installments = [],
  isLoading = false,
}: RecurringStatsCardProps) {
  const t = useTranslations("recurring");

  // 計算統計資料
  const activeSubscriptions = subscriptions.filter(
    (s) => s.status === "active"
  );
  const activeInstallments = installments.filter((i) => i.status === "active");

  const monthlySubscriptionCost =
    calculateMonthlySubscriptionCost(activeSubscriptions);
  const monthlyInstallmentPayment =
    calculateMonthlyInstallmentPayment(activeInstallments);
  const totalMonthlyCost = monthlySubscriptionCost + monthlyInstallmentPayment;

  // 計算即將到期的訂閱（30天內）
  const expiringSubscriptions = activeSubscriptions.filter((sub) => {
    if (!sub.end_date) return false;
    const endDate = new Date(sub.end_date);
    const now = new Date();
    const daysUntilExpiry = Math.ceil(
      (endDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)
    );
    return daysUntilExpiry > 0 && daysUntilExpiry <= 30;
  });

  // 計算即將完成的分期（3個月內）
  const completingInstallments = activeInstallments.filter((inst) => {
    const remainingMonths = inst.installment_count - inst.paid_count;
    return remainingMonths > 0 && remainingMonths <= 3;
  });

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <CalendarIcon className="h-5 w-5" />
            {t("title")}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-20 w-full" />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <CalendarIcon className="h-5 w-5" />
          {t("title")}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* 每月總支出 */}
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium text-muted-foreground">
              {t("monthlyTotal")}
            </p>
            <Badge variant="outline" className="gap-1">
              <TrendingUpIcon className="h-3 w-3" />
              {t("fixedExpense")}
            </Badge>
          </div>
          <div className="text-3xl font-bold tabular-nums">
            NT${" "}
            {totalMonthlyCost.toLocaleString("zh-TW", {
              maximumFractionDigits: 0,
            })}
          </div>
          <div className="flex gap-4 text-xs text-muted-foreground">
            <span>
              {t("subscriptions")}: NT${" "}
              {monthlySubscriptionCost.toLocaleString("zh-TW", {
                maximumFractionDigits: 0,
              })}
            </span>
            <span>
              {t("installments")}: NT${" "}
              {monthlyInstallmentPayment.toLocaleString("zh-TW", {
                maximumFractionDigits: 0,
              })}
            </span>
          </div>
        </div>

        {/* 訂閱統計 */}
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <CreditCardIcon className="h-4 w-4 text-muted-foreground" />
            <p className="text-sm font-medium">{t("subscriptionServices")}</p>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-2xl font-semibold tabular-nums">
                {activeSubscriptions.length}
              </p>
              <p className="text-xs text-muted-foreground">{t("active")}</p>
            </div>
            <div>
              <p className="text-2xl font-semibold tabular-nums">
                {expiringSubscriptions.length}
              </p>
              <p className="text-xs text-muted-foreground">
                {t("expiringSoon")}
              </p>
            </div>
          </div>
        </div>

        {/* 分期統計 */}
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <CreditCardIcon className="h-4 w-4 text-muted-foreground" />
            <p className="text-sm font-medium">{t("installmentPayments")}</p>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-2xl font-semibold tabular-nums">
                {activeInstallments.length}
              </p>
              <p className="text-xs text-muted-foreground">{t("active")}</p>
            </div>
            <div>
              <p className="text-2xl font-semibold tabular-nums">
                {completingInstallments.length}
              </p>
              <p className="text-xs text-muted-foreground">
                {t("completingSoon")}
              </p>
            </div>
          </div>
        </div>

        {/* 提醒 */}
        {(expiringSubscriptions.length > 0 ||
          completingInstallments.length > 0) && (
          <div className="flex items-start gap-2 rounded-lg bg-amber-50 p-3 text-amber-900">
            <AlertCircleIcon className="h-4 w-4 mt-0.5 shrink-0" />
            <div className="text-xs space-y-1">
              {expiringSubscriptions.length > 0 && (
                <p>
                  {t("subscriptionsExpiringAlert", {
                    count: expiringSubscriptions.length,
                  })}
                </p>
              )}
              {completingInstallments.length > 0 && (
                <p>
                  {t("installmentsCompletingAlert", {
                    count: completingInstallments.length,
                  })}
                </p>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
