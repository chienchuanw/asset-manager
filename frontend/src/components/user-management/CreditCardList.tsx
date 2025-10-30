/**
 * 信用卡列表元件
 * 顯示所有信用卡的列表
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
import type { CreditCard } from "@/types/user-management";

interface CreditCardListProps {
  creditCards?: CreditCard[];
  isLoading?: boolean;
  onEdit?: (creditCard: CreditCard) => void;
  onDelete?: (id: string) => void;
}

/**
 * 計算信用卡使用率
 */
function calculateUtilization(card: CreditCard): number {
  if (card.credit_limit === 0) return 0;
  return (card.used_credit / card.credit_limit) * 100;
}

/**
 * 取得使用率顏色
 */
function getUtilizationColor(utilization: number): string {
  if (utilization >= 80) return "text-destructive";
  if (utilization >= 50) return "text-yellow-600";
  return "text-green-600";
}

/**
 * 計算可用額度
 */
function getAvailableCredit(card: CreditCard): number {
  return card.credit_limit - card.used_credit;
}

export function CreditCardList({
  creditCards = [],
  isLoading = false,
  onEdit,
  onDelete,
}: CreditCardListProps) {
  return (
    <div className="rounded-md border mx-4">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>發卡銀行</TableHead>
            <TableHead>卡片名稱</TableHead>
            <TableHead>卡號後四碼</TableHead>
            <TableHead>帳單日</TableHead>
            <TableHead>繳款日</TableHead>
            <TableHead className="text-right">信用額度</TableHead>
            <TableHead>使用率</TableHead>
            <TableHead>備註</TableHead>
            <TableHead className="w-[50px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading ? (
            <TableRow>
              <TableCell colSpan={9} className="h-24 text-center">
                <div className="flex justify-center">
                  <div className="space-y-2">
                    {[...Array(3)].map((_, i) => (
                      <Skeleton key={i} className="h-12 w-[1000px]" />
                    ))}
                  </div>
                </div>
              </TableCell>
            </TableRow>
          ) : creditCards.length === 0 ? (
            <TableRow>
              <TableCell
                colSpan={9}
                className="h-24 text-center text-muted-foreground"
              >
                尚無信用卡記錄
              </TableCell>
            </TableRow>
          ) : (
            creditCards.map((card) => {
              const utilization = calculateUtilization(card);
              const availableCredit = getAvailableCredit(card);

              return (
                <TableRow key={card.id}>
                  <TableCell className="font-medium">
                    {card.issuing_bank}
                  </TableCell>
                  <TableCell>{card.card_name}</TableCell>
                  <TableCell className="font-mono">
                    ****{card.card_number_last4}
                  </TableCell>
                  <TableCell>每月 {card.billing_day} 日</TableCell>
                  <TableCell>每月 {card.payment_due_day} 日</TableCell>
                  <TableCell className="text-right">
                    <div className="space-y-1">
                      <div className="tabular-nums">
                        {card.credit_limit.toLocaleString("zh-TW", {
                          minimumFractionDigits: 0,
                          maximumFractionDigits: 0,
                        })}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        可用:{" "}
                        {availableCredit.toLocaleString("zh-TW", {
                          minimumFractionDigits: 0,
                          maximumFractionDigits: 0,
                        })}
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="space-y-2 min-w-[120px]">
                      <div className="flex items-center justify-between">
                        <span
                          className={`text-sm font-medium ${getUtilizationColor(
                            utilization
                          )}`}
                        >
                          {utilization.toFixed(1)}%
                        </span>
                      </div>
                      <Progress value={utilization} className="h-2" />
                    </div>
                  </TableCell>
                  <TableCell className="max-w-[150px] truncate text-muted-foreground">
                    {card.note || "-"}
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
                          <DropdownMenuItem onClick={() => onEdit(card)}>
                            <PencilIcon className="mr-2 h-4 w-4" />
                            編輯
                          </DropdownMenuItem>
                        )}
                        {onDelete && (
                          <DropdownMenuItem
                            onClick={() => onDelete(card.id)}
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
            })
          )}
        </TableBody>
      </Table>
    </div>
  );
}
