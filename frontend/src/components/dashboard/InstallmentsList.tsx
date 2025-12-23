/**
 * 分期列表元件
 * 顯示所有分期付款的列表
 */

"use client";

import { useTranslations } from "next-intl";
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
 * 計算進度百分比
 */
function calculateProgress(paidCount: number, totalCount: number): number {
  return Math.round((paidCount / totalCount) * 100);
}

export function InstallmentsList({
  installments = [],
  isLoading = false,
  onEdit,
  onDelete,
}: InstallmentsListProps) {
  const t = useTranslations("recurring");
  const tCommon = useTranslations("common");

  // 格式化付款方式
  const formatPaymentMethod = (method: PaymentMethod): string => {
    const methodMap: Record<PaymentMethod, string> = {
      cash: t("cash"),
      bank_account: t("bankAccount"),
      credit_card: t("creditCard"),
    };
    return methodMap[method] || method;
  };

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
        <p>{t("noInstallments")}</p>
      </div>
    );
  }

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{tCommon("name")}</TableHead>
            <TableHead>{t("category")}</TableHead>
            <TableHead>{t("totalAmount")}</TableHead>
            <TableHead>{t("monthlyPayment")}</TableHead>
            <TableHead>{t("paymentMethod")}</TableHead>
            <TableHead>{t("progress")}</TableHead>
            <TableHead>{tCommon("status")}</TableHead>
            <TableHead className="w-[50px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {installments.map((installment) => {
            const progress = calculateProgress(
              installment.paid_count,
              installment.installment_count
            );

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
                        {t("interest")}: {installment.currency}{" "}
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
                        {installment.installment_count} {t("periods")}
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
                      ? t("active")
                      : installment.status === "completed"
                      ? t("completed")
                      : t("cancelled")}
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
                          {tCommon("edit")}
                        </DropdownMenuItem>
                      )}
                      {onDelete && (
                        <DropdownMenuItem
                          onClick={() => onDelete(installment.id)}
                          className="text-destructive"
                        >
                          <TrashIcon className="mr-2 h-4 w-4" />
                          {tCommon("delete")}
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
