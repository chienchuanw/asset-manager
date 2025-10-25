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
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useCashFlows, useDeleteCashFlow } from "@/hooks";
import {
  type CashFlow,
  type CashFlowFilters,
  getCashFlowTypeLabel,
  getCashFlowTypeColor,
  getCashFlowTypeBgColor,
  formatAmount,
  formatDate,
} from "@/types/cash-flow";
import { MoreHorizontal, Trash2, Loader2 } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

interface CashFlowListProps {
  filters?: CashFlowFilters;
  onRefresh?: () => void;
}

/**
 * 現金流列表元件
 *
 * 顯示現金流記錄的表格，支援刪除操作
 */
export function CashFlowList({ filters, onRefresh }: CashFlowListProps) {
  const [deletingId, setDeletingId] = useState<string | null>(null);

  // 取得現金流列表
  const { data: cashFlows, isLoading } = useCashFlows(filters);

  // 刪除現金流 mutation
  const deleteMutation = useDeleteCashFlow({
    onSuccess: () => {
      toast.success("記錄刪除成功");
      setDeletingId(null);
      onRefresh?.();
    },
    onError: (error) => {
      toast.error(error.message || "刪除失敗");
      setDeletingId(null);
    },
  });

  const handleDelete = (id: string) => {
    if (confirm("確定要刪除這筆記錄嗎？")) {
      setDeletingId(id);
      deleteMutation.mutate(id);
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3, 4, 5].map((i) => (
          <Skeleton key={i} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  if (!cashFlows || cashFlows.length === 0) {
    return (
      <div className="text-center py-12 text-muted-foreground">
        <p>尚無記錄</p>
        <p className="text-sm mt-2">點擊「新增記錄」開始記錄您的現金流</p>
      </div>
    );
  }

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>日期</TableHead>
            <TableHead>類型</TableHead>
            <TableHead>分類</TableHead>
            <TableHead>描述</TableHead>
            <TableHead className="text-right">金額</TableHead>
            <TableHead className="w-[50px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {cashFlows.map((cashFlow) => (
            <TableRow key={cashFlow.id}>
              <TableCell className="font-medium">
                {formatDate(cashFlow.date)}
              </TableCell>
              <TableCell>
                <Badge
                  variant="outline"
                  className={`${getCashFlowTypeBgColor(
                    cashFlow.type
                  )} ${getCashFlowTypeColor(cashFlow.type)} border-0`}
                >
                  {getCashFlowTypeLabel(cashFlow.type)}
                </Badge>
              </TableCell>
              <TableCell>
                {cashFlow.category?.name || "-"}
                {cashFlow.category?.is_system && (
                  <span className="ml-1 text-xs text-muted-foreground">
                    (系統)
                  </span>
                )}
              </TableCell>
              <TableCell>
                <div>
                  <div className="font-medium">{cashFlow.description}</div>
                  {cashFlow.note && (
                    <div className="text-sm text-muted-foreground mt-1">
                      {cashFlow.note}
                    </div>
                  )}
                </div>
              </TableCell>
              <TableCell
                className={`text-right font-semibold ${getCashFlowTypeColor(
                  cashFlow.type
                )}`}
              >
                {cashFlow.type === "income" ? "+" : "-"}$
                {formatAmount(cashFlow.amount)}
              </TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-8 w-8 p-0"
                      disabled={deletingId === cashFlow.id}
                    >
                      {deletingId === cashFlow.id ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <MoreHorizontal className="h-4 w-4" />
                      )}
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem
                      className="text-red-600"
                      onClick={() => handleDelete(cashFlow.id)}
                    >
                      <Trash2 className="mr-2 h-4 w-4" />
                      刪除
                    </DropdownMenuItem>
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

