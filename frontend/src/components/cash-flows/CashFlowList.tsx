"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { CashFlowCard } from "./CashFlowCard";
import { useCashFlows, useDeleteCashFlow } from "@/hooks";
import { type CashFlowFilters } from "@/types/cash-flow";
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

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[1, 2, 3].map((i) => (
          <Card key={i}>
            <CardContent className="p-4">
              <Skeleton className="h-32 w-full" />
            </CardContent>
          </Card>
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
    <div className="space-y-4">
      {cashFlows.map((cashFlow) => (
        <CashFlowCard
          key={cashFlow.id}
          cashFlow={cashFlow}
          onEdit={() => {
            // TODO: 實作編輯功能
            console.log("Edit cash flow:", cashFlow.id);
          }}
          onDelete={handleDelete}
          isDeleting={deleteMutation.isPending}
        />
      ))}
    </div>
  );
}
