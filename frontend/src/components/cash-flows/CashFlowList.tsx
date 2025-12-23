"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { useCashFlows, useDeleteCashFlow } from "@/hooks";
import {
  type CashFlowFilters,
  type CashFlow,
  getCashFlowTypeLabel,
  formatAmount,
} from "@/types/cash-flow";
import {
  MoreVertical,
  Pencil,
  Trash2,
  TrendingUp,
  TrendingDown,
  ArrowDownToLine,
  ArrowUpFromLine,
} from "lucide-react";
import { toast } from "sonner";
import { PaymentMethodDisplay } from "./PaymentMethodDisplay";
import { EditCashFlowDialog } from "./EditCashFlowDialog";

interface CashFlowListProps {
  filters?: CashFlowFilters;
}

/**
 * 現金流列表元件
 *
 * 顯示現金流記錄的表格列表，支援編輯和刪除操作
 */
export function CashFlowList({ filters }: CashFlowListProps) {
  const t = useTranslations("cashFlows");
  const tCommon = useTranslations("common");

  // 編輯狀態
  const [editingCashFlow, setEditingCashFlow] = useState<CashFlow | null>(null);

  // 取得現金流列表
  const { data: cashFlows, isLoading } = useCashFlows(filters, {
    // 確保資料總是最新的
    staleTime: 0,
  });

  // 刪除現金流 mutation
  const deleteMutation = useDeleteCashFlow({
    onError: (error) => {
      toast.error(error.message || t("errorMessage"));
    },
  });

  const handleDelete = (id: string) => {
    if (confirm(t("confirmDelete"))) {
      deleteMutation.mutate(id);
    }
  };

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{t("type")}</TableHead>
            <TableHead>{t("category")}</TableHead>
            <TableHead>{t("paymentMethod")}</TableHead>
            <TableHead>{t("description")}</TableHead>
            <TableHead className="text-right">{t("amount")}</TableHead>
            <TableHead className="w-[50px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading ? (
            // 載入中骨架屏
            Array.from({ length: 5 }).map((_, index) => (
              <TableRow key={index}>
                <TableCell>
                  <Skeleton className="h-4 w-16" />
                </TableCell>
                <TableCell>
                  <Skeleton className="h-4 w-24" />
                </TableCell>
                <TableCell>
                  <Skeleton className="h-4 w-20" />
                </TableCell>
                <TableCell>
                  <Skeleton className="h-4 w-32" />
                </TableCell>
                <TableCell>
                  <Skeleton className="h-4 w-20" />
                </TableCell>
                <TableCell>
                  <Skeleton className="h-4 w-8" />
                </TableCell>
              </TableRow>
            ))
          ) : !cashFlows || cashFlows.length === 0 ? (
            // 無資料
            <TableRow>
              <TableCell colSpan={7} className="h-24 text-center">
                <p className="text-muted-foreground">{t("noRecords")}</p>
              </TableCell>
            </TableRow>
          ) : (
            // 現金流記錄列表
            cashFlows.map((cashFlow) => {
              // 根據類型決定顏色和圖示
              let typeColor: string;
              let TypeIcon: any;

              switch (cashFlow.type) {
                case "income":
                  typeColor = "text-green-600";
                  TypeIcon = TrendingUp;
                  break;
                case "expense":
                  typeColor = "text-red-600";
                  TypeIcon = TrendingDown;
                  break;
                case "transfer_in":
                  typeColor = "text-gray-600";
                  TypeIcon = ArrowDownToLine;
                  break;
                case "transfer_out":
                  typeColor = "text-gray-600";
                  TypeIcon = ArrowUpFromLine;
                  break;
                default:
                  typeColor = "text-gray-600";
                  TypeIcon = TrendingUp;
              }

              return (
                <TableRow key={cashFlow.id}>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <TypeIcon className={`h-4 w-4 ${typeColor}`} />
                    </div>
                  </TableCell>
                  <TableCell>
                    {cashFlow.category ? (
                      <Badge variant="secondary">
                        {cashFlow.category.name}
                      </Badge>
                    ) : (
                      <span className="text-muted-foreground text-sm">-</span>
                    )}
                  </TableCell>
                  <TableCell>
                    <PaymentMethodDisplay cashFlow={cashFlow} />
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col">
                      <span className="text-sm">{cashFlow.description}</span>
                      {cashFlow.note && (
                        <span className="text-xs text-muted-foreground">
                          {cashFlow.note}
                        </span>
                      )}
                    </div>
                  </TableCell>
                  <TableCell
                    className={`text-right tabular-nums text-sm font-medium ${typeColor}`}
                  >
                    {cashFlow.type === "income" ||
                    cashFlow.type === "transfer_in"
                      ? "+"
                      : "-"}
                    ${formatAmount(cashFlow.amount)}
                  </TableCell>
                  <TableCell>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="sm">
                          <MoreVertical className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem
                          onClick={() => setEditingCashFlow(cashFlow)}
                        >
                          <Pencil className="h-4 w-4 mr-2" />
                          {tCommon("edit")}
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() => handleDelete(cashFlow.id)}
                          className="text-red-600"
                          disabled={deleteMutation.isPending}
                        >
                          <Trash2 className="h-4 w-4 mr-2" />
                          {tCommon("delete")}
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              );
            })
          )}
        </TableBody>
      </Table>

      {/* 編輯現金流對話框 */}
      {editingCashFlow && (
        <EditCashFlowDialog
          cashFlow={editingCashFlow}
          open={!!editingCashFlow}
          onOpenChange={(open) => {
            if (!open) setEditingCashFlow(null);
          }}
          onSuccess={() => {
            // React Query 的自動失效機制會處理資料更新
            setEditingCashFlow(null);
          }}
        />
      )}
    </div>
  );
}
