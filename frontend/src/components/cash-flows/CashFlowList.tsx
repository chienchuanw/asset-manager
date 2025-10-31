"use client";

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
  getCashFlowTypeLabel,
  formatAmount,
} from "@/types/cash-flow";
import {
  MoreVertical,
  Pencil,
  Trash2,
  TrendingUp,
  TrendingDown,
} from "lucide-react";
import { toast } from "sonner";

interface CashFlowListProps {
  filters?: CashFlowFilters;
  onRefresh?: () => void;
}

/**
 * 現金流列表元件
 *
 * 顯示現金流記錄的卡片列表，支援刪除操作
 */
export function CashFlowList({ filters, onRefresh }: CashFlowListProps) {
  // 取得現金流列表
  const { data: cashFlows, isLoading } = useCashFlows(filters);

  // 刪除現金流 mutation
  const deleteMutation = useDeleteCashFlow({
    onSuccess: () => {
      toast.success("記錄刪除成功");
      onRefresh?.();
    },
    onError: (error) => {
      toast.error(error.message || "刪除失敗");
    },
  });

  const handleDelete = (id: string) => {
    if (confirm("確定要刪除這筆記錄嗎？")) {
      deleteMutation.mutate(id);
    }
  };

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
          {isLoading ? (
            // 載入中骨架屏
            Array.from({ length: 5 }).map((_, index) => (
              <TableRow key={index}>
                <TableCell>
                  <Skeleton className="h-4 w-20" />
                </TableCell>
                <TableCell>
                  <Skeleton className="h-4 w-16" />
                </TableCell>
                <TableCell>
                  <Skeleton className="h-4 w-24" />
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
              <TableCell colSpan={6} className="h-24 text-center">
                <p className="text-muted-foreground">尚無現金流記錄</p>
              </TableCell>
            </TableRow>
          ) : (
            // 現金流記錄列表
            cashFlows.map((cashFlow) => {
              const isIncome = cashFlow.type === "income";
              const typeColor = isIncome ? "text-green-600" : "text-red-600";
              const TypeIcon = isIncome ? TrendingUp : TrendingDown;

              return (
                <TableRow key={cashFlow.id}>
                  <TableCell className="text-sm">
                    {new Date(cashFlow.date).toLocaleDateString("zh-TW")}
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <TypeIcon className={`h-4 w-4 ${typeColor}`} />
                      <Badge
                        variant="outline"
                        className={`${typeColor} border-current`}
                      >
                        {getCashFlowTypeLabel(cashFlow.type)}
                      </Badge>
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
                    {isIncome ? "+" : "-"}${formatAmount(cashFlow.amount)}
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
                          onClick={() => {
                            // TODO: 實作編輯功能
                            console.log("Edit cash flow:", cashFlow.id);
                          }}
                        >
                          <Pencil className="h-4 w-4 mr-2" />
                          編輯
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() => handleDelete(cashFlow.id)}
                          className="text-red-600"
                          disabled={deleteMutation.isPending}
                        >
                          <Trash2 className="h-4 w-4 mr-2" />
                          刪除
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
    </div>
  );
}
