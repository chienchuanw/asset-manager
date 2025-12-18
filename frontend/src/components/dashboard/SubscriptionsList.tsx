/**
 * 訂閱列表元件
 * 顯示所有訂閱服務的列表
 */

"use client";

import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import {
  MoreHorizontalIcon,
  PencilIcon,
  TrashIcon,
  XCircleIcon,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import type {
  Subscription,
  BillingCycle,
  PaymentMethod,
} from "@/types/subscription";

interface SubscriptionsListProps {
  subscriptions?: Subscription[];
  isLoading?: boolean;
  onEdit?: (subscription: Subscription) => void;
  onDelete?: (id: string) => void;
  onCancel?: (subscription: Subscription) => void;
}

/**
 * 格式化計費週期
 */
function formatBillingCycle(cycle: BillingCycle): string {
  const cycleMap: Record<BillingCycle, string> = {
    monthly: "每月",
    quarterly: "每季",
    yearly: "每年",
  };
  return cycleMap[cycle] || cycle;
}

/**
 * 格式化付款方式
 */
function formatPaymentMethod(method: PaymentMethod): string {
  const methodMap: Record<PaymentMethod, string> = {
    cash: "現金",
    bank_account: "銀行帳戶",
    credit_card: "信用卡",
  };
  return methodMap[method] || method;
}

/**
 * 格式化日期
 */
function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString("zh-TW", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  });
}

/**
 * 計算下次扣款日期
 */
function getNextBillingDate(subscription: Subscription): string {
  const now = new Date();
  const currentYear = now.getFullYear();
  const currentMonth = now.getMonth();
  const billingDay = subscription.billing_day;

  let nextDate = new Date(currentYear, currentMonth, billingDay);

  // 如果這個月的扣款日已過，計算下個週期
  if (nextDate <= now) {
    switch (subscription.billing_cycle) {
      case "monthly":
        nextDate = new Date(currentYear, currentMonth + 1, billingDay);
        break;
      case "quarterly":
        nextDate = new Date(currentYear, currentMonth + 3, billingDay);
        break;
      case "yearly":
        nextDate = new Date(currentYear + 1, currentMonth, billingDay);
        break;
    }
  }

  return formatDate(nextDate.toISOString());
}

export function SubscriptionsList({
  subscriptions = [],
  isLoading = false,
  onEdit,
  onDelete,
  onCancel,
}: SubscriptionsListProps) {
  if (isLoading) {
    return (
      <div className="space-y-2">
        {[...Array(5)].map((_, i) => (
          <Skeleton key={i} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  if (subscriptions.length === 0) {
    return (
      <div className="text-center py-12 text-muted-foreground">
        <p>尚無訂閱記錄</p>
      </div>
    );
  }

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>名稱</TableHead>
            <TableHead>分類</TableHead>
            <TableHead>金額</TableHead>
            <TableHead>計費週期</TableHead>
            <TableHead>付款方式</TableHead>
            <TableHead>下次扣款</TableHead>
            <TableHead>狀態</TableHead>
            <TableHead className="w-[50px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {subscriptions.map((subscription) => (
            <TableRow key={subscription.id}>
              <TableCell className="font-medium">{subscription.name}</TableCell>
              <TableCell>{subscription.category?.name || "-"}</TableCell>
              <TableCell className="tabular-nums">
                {subscription.currency}{" "}
                {subscription.amount.toLocaleString("zh-TW")}
              </TableCell>
              <TableCell>
                {formatBillingCycle(subscription.billing_cycle)}
              </TableCell>
              <TableCell>
                {formatPaymentMethod(subscription.payment_method)}
              </TableCell>
              <TableCell className="tabular-nums">
                {subscription.status === "active"
                  ? getNextBillingDate(subscription)
                  : "-"}
              </TableCell>
              <TableCell>
                <Badge
                  variant={
                    subscription.status === "active" ? "default" : "secondary"
                  }
                >
                  {subscription.status === "active" ? "進行中" : "已取消"}
                </Badge>
              </TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon">
                      <MoreHorizontalIcon className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    {onEdit && (
                      <DropdownMenuItem onClick={() => onEdit(subscription)}>
                        <PencilIcon className="mr-2 h-4 w-4" />
                        編輯
                      </DropdownMenuItem>
                    )}
                    {onCancel && subscription.status === "active" && (
                      <DropdownMenuItem onClick={() => onCancel(subscription)}>
                        <XCircleIcon className="mr-2 h-4 w-4" />
                        取消訂閱
                      </DropdownMenuItem>
                    )}
                    {onDelete && (
                      <DropdownMenuItem
                        onClick={() => onDelete(subscription.id)}
                        className="text-destructive"
                      >
                        <TrashIcon className="mr-2 h-4 w-4" />
                        刪除
                      </DropdownMenuItem>
                    )}
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
