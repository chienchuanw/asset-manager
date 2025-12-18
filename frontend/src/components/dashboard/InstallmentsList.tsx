/**
 * 分期列表元件
 * 顯示所有分期付款的列表
 */

"use client";

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
import { Progress } from "@/components/ui/progress";
import { MoreHorizontalIcon, PencilIcon, TrashIcon } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import type { Installment } from "@/types/installment";
import type { PaymentMethod } from "@/types/subscription";

interface InstallmentsListProps {
  installments?: Installment[];
  isLoading?: boolean;
  onEdit?: (installment: Installment) => void;
  onDelete?: (id: string) => void;
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
 * 計算進度百分比
 */
function calculateProgress(paidCount: number, totalCount: number): number {
  return Math.round((paidCount / totalCount) * 100);
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

export function InstallmentsList({
  installments = [],
  isLoading = false,
  onEdit,
  onDelete,
}: InstallmentsListProps) {
  if (isLoading) {
    return (
      <div className="space-y-2">
        {[...Array(5)].map((_, i) => (
          <Skeleton key={i} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  if (installments.length === 0) {
    return (
      <div className="text-center py-12 text-muted-foreground">
        <p>尚無分期記錄</p>
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
            <TableHead>總金額</TableHead>
            <TableHead>每期金額</TableHead>
            <TableHead>付款方式</TableHead>
            <TableHead>進度</TableHead>
            <TableHead>狀態</TableHead>
            <TableHead className="w-[50px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {installments.map((installment) => {
            const progress = calculateProgress(
              installment.paid_count,
              installment.installment_count
            );
            const remaining =
              installment.installment_count - installment.paid_count;

            return (
              <TableRow key={installment.id}>
                <TableCell className="font-medium">
                  {installment.name}
                </TableCell>
                <TableCell>{installment.category?.name || "-"}</TableCell>
                <TableCell className="tabular-nums">
                  <div>
                    <p>
                      {installment.currency}{" "}
                      {installment.total_amount.toLocaleString("zh-TW")}
                    </p>
                    {installment.total_interest > 0 && (
                      <p className="text-xs text-muted-foreground">
                        利息: {installment.currency}{" "}
                        {installment.total_interest.toLocaleString("zh-TW")}
                      </p>
                    )}
                  </div>
                </TableCell>
                <TableCell className="tabular-nums">
                  {installment.currency}{" "}
                  {installment.installment_amount?.toLocaleString("zh-TW") ??
                    "0"}
                </TableCell>
                <TableCell>
                  {formatPaymentMethod(installment.payment_method)}
                </TableCell>
                <TableCell>
                  <div className="space-y-1">
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-muted-foreground">
                        {installment.paid_count} /{" "}
                        {installment.installment_count} 期
                      </span>
                      <span className="font-medium">{progress}%</span>
                    </div>
                    <Progress value={progress} className="h-2" />
                  </div>
                </TableCell>
                <TableCell>
                  <Badge
                    variant={
                      installment.status === "active"
                        ? "default"
                        : installment.status === "completed"
                        ? "secondary"
                        : "outline"
                    }
                  >
                    {installment.status === "active"
                      ? "進行中"
                      : installment.status === "completed"
                      ? "已完成"
                      : "已取消"}
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
                        <DropdownMenuItem onClick={() => onEdit(installment)}>
                          <PencilIcon className="mr-2 h-4 w-4" />
                          編輯
                        </DropdownMenuItem>
                      )}
                      {onDelete && (
                        <DropdownMenuItem
                          onClick={() => onDelete(installment.id)}
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
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}
