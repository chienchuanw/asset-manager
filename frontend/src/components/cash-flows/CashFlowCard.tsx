/**
 * 現金流記錄卡片元件
 * 用於手機友善的卡片式顯示,取代傳統的 Table
 */

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  CashFlowType,
  getCashFlowTypeLabel,
  type CashFlow,
} from "@/types/cash-flow";
import { Edit, MoreVertical, Trash2 } from "lucide-react";

interface CashFlowCardProps {
  cashFlow: CashFlow;
  onEdit: (cashFlow: CashFlow) => void;
  onDelete: (id: string) => void;
  isDeleting?: boolean;
}

/**
 * 取得現金流類型的顏色
 */
const getCashFlowTypeColor = (type: CashFlowType) => {
  switch (type) {
    case CashFlowType.INCOME:
      return "bg-green-100 text-green-800";
    case CashFlowType.EXPENSE:
      return "bg-red-100 text-red-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
};

/**
 * 取得現金流類型的圖示
 */
const getCashFlowTypeIcon = (type: CashFlowType) => {
  switch (type) {
    case CashFlowType.INCOME:
      return "↑";
    case CashFlowType.EXPENSE:
      return "↓";
    default:
      return "•";
  }
};

export function CashFlowCard({
  cashFlow,
  onEdit,
  onDelete,
  isDeleting = false,
}: CashFlowCardProps) {
  return (
    <Card className="overflow-hidden">
      <CardContent className="p-4">
        {/* 頂部：日期和操作按鈕 */}
        <div className="flex items-start justify-between mb-3">
          <div className="text-sm font-medium text-muted-foreground">
            {new Date(cashFlow.date).toLocaleDateString("zh-TW", {
              year: "numeric",
              month: "2-digit",
              day: "2-digit",
            })}
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0"
                disabled={isDeleting}
              >
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => onEdit(cashFlow)}>
                <Edit className="mr-2 h-4 w-4" />
                編輯
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => onDelete(cashFlow.id)}
                className="text-red-600"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                刪除
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        {/* 類型標籤 */}
        <div className="flex gap-2 mb-3">
          <Badge
            variant="outline"
            className={getCashFlowTypeColor(cashFlow.type)}
          >
            <span className="mr-1">{getCashFlowTypeIcon(cashFlow.type)}</span>
            {getCashFlowTypeLabel(cashFlow.type)}
          </Badge>
          {cashFlow.category && (
            <Badge variant="outline" className="bg-blue-100 text-blue-800">
              {cashFlow.category.name}
            </Badge>
          )}
        </div>

        {/* 描述 */}
        <div className="mb-3">
          <div className="font-semibold text-lg">{cashFlow.description}</div>
        </div>

        {/* 金額 */}
        <div className="bg-gray-50 rounded-lg p-3 mb-3">
          <div className="flex items-center justify-between">
            <div className="text-sm text-muted-foreground">金額</div>
            <div className="flex items-center gap-2">
              <Badge
                variant="outline"
                className="bg-amber-100 text-amber-800 text-xs"
              >
                {cashFlow.currency}
              </Badge>
              <div
                className={`text-lg font-bold ${
                  cashFlow.type === CashFlowType.INCOME
                    ? "text-green-600"
                    : "text-red-600"
                }`}
              >
                {cashFlow.type === CashFlowType.INCOME ? "+" : "-"}
                {cashFlow.amount.toLocaleString()}
              </div>
            </div>
          </div>
        </div>

        {/* 備註（如果有） */}
        {cashFlow.note && (
          <div className="text-sm text-muted-foreground border-t pt-3">
            <div className="text-xs mb-1">備註</div>
            <div>{cashFlow.note}</div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

