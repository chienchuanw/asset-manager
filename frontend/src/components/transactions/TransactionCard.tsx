/**
 * 交易記錄卡片元件
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
  getAssetTypeLabel,
  getTransactionTypeLabel,
  type Transaction,
  TransactionType,
  AssetType,
} from "@/types/transaction";
import { Edit, MoreVertical, Trash2 } from "lucide-react";

interface TransactionCardProps {
  transaction: Transaction;
  onEdit: (transaction: Transaction) => void;
  onDelete: (id: string) => void;
  isDeleting?: boolean;
}

/**
 * 取得交易類型的顏色（台灣習慣：紅漲綠跌）
 */
const getTransactionTypeColor = (type: TransactionType) => {
  switch (type) {
    case TransactionType.BUY:
      return "bg-red-100 text-red-800";
    case TransactionType.SELL:
      return "bg-green-100 text-green-800";
    case TransactionType.DIVIDEND:
      return "bg-blue-100 text-blue-800";
    case TransactionType.FEE:
      return "bg-gray-100 text-gray-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
};

/**
 * 取得資產類型的顏色
 */
const getAssetTypeColor = (assetType: AssetType) => {
  switch (assetType) {
    case AssetType.TW_STOCK:
      return "bg-purple-100 text-purple-800";
    case AssetType.US_STOCK:
      return "bg-indigo-100 text-indigo-800";
    case AssetType.CRYPTO:
      return "bg-orange-100 text-orange-800";
    case AssetType.CASH:
      return "bg-emerald-100 text-emerald-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
};

export function TransactionCard({
  transaction,
  onEdit,
  onDelete,
  isDeleting = false,
}: TransactionCardProps) {
  return (
    <Card className="overflow-hidden">
      <CardContent className="p-4">
        {/* 頂部：日期和操作按鈕 */}
        <div className="flex items-start justify-between mb-3">
          <div className="text-sm font-medium text-muted-foreground">
            {new Date(transaction.date).toLocaleDateString("zh-TW", {
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
              <DropdownMenuItem onClick={() => onEdit(transaction)}>
                <Edit className="mr-2 h-4 w-4" />
                編輯
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => onDelete(transaction.id)}
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
            className={getTransactionTypeColor(transaction.type)}
          >
            {getTransactionTypeLabel(transaction.type)}
          </Badge>
          <Badge
            variant="outline"
            className={getAssetTypeColor(transaction.asset_type)}
          >
            {getAssetTypeLabel(transaction.asset_type)}
          </Badge>
        </div>

        {/* 資產資訊 */}
        <div className="mb-3">
          <div className="font-semibold text-lg">{transaction.symbol}</div>
          <div className="text-sm text-muted-foreground">
            {transaction.name}
          </div>
        </div>

        {/* 交易詳情 */}
        <div className="grid grid-cols-2 gap-3 mb-3">
          <div>
            <div className="text-xs text-muted-foreground mb-1">數量</div>
            <div className="font-medium">
              {transaction.quantity.toLocaleString()}
            </div>
          </div>
          <div>
            <div className="text-xs text-muted-foreground mb-1">單價</div>
            <div className="font-medium">
              {transaction.price.toLocaleString()}
            </div>
          </div>
        </div>

        {/* 總金額 */}
        <div className="bg-gray-50 rounded-lg p-3 mb-3">
          <div className="flex items-center justify-between">
            <div className="text-sm text-muted-foreground">總金額</div>
            <div className="flex items-center gap-2">
              <Badge
                variant="outline"
                className="bg-amber-100 text-amber-800 text-xs"
              >
                {transaction.currency}
              </Badge>
              <div className="text-lg font-bold">
                {transaction.amount.toLocaleString()}
              </div>
            </div>
          </div>
        </div>

        {/* 費用資訊（如果有） */}
        {(transaction.fee || transaction.tax) && (
          <div className="grid grid-cols-2 gap-3 mb-3 text-sm">
            {transaction.fee && (
              <div>
                <div className="text-xs text-muted-foreground mb-1">手續費</div>
                <div className="font-medium">
                  {transaction.fee.toLocaleString()}
                </div>
              </div>
            )}
            {transaction.tax && (
              <div>
                <div className="text-xs text-muted-foreground mb-1">交易稅</div>
                <div className="font-medium">
                  {transaction.tax.toLocaleString()}
                </div>
              </div>
            )}
          </div>
        )}

        {/* 備註（如果有） */}
        {transaction.note && (
          <div className="text-sm text-muted-foreground border-t pt-3">
            <div className="text-xs mb-1">備註</div>
            <div>{transaction.note}</div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

