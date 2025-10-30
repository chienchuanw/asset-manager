/**
 * 銀行帳戶列表元件
 * 顯示所有銀行帳戶的列表
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
import { MoreHorizontalIcon, PencilIcon, TrashIcon } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import type { BankAccount } from "@/types/user-management";

interface BankAccountListProps {
  bankAccounts?: BankAccount[];
  isLoading?: boolean;
  onEdit?: (bankAccount: BankAccount) => void;
  onDelete?: (id: string) => void;
}

/**
 * 取得幣別顏色
 */
function getCurrencyColor(currency: string): string {
  const colorMap: Record<string, string> = {
    TWD: "default",
    USD: "secondary",
    JPY: "outline",
    EUR: "outline",
    CNY: "outline",
  };
  return colorMap[currency] || "default";
}

export function BankAccountList({
  bankAccounts = [],
  isLoading = false,
  onEdit,
  onDelete,
}: BankAccountListProps) {
  return (
    <div className="rounded-md border mx-4">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>銀行名稱</TableHead>
            <TableHead>帳戶類型</TableHead>
            <TableHead>帳號後四碼</TableHead>
            <TableHead>幣別</TableHead>
            <TableHead className="text-right">餘額</TableHead>
            <TableHead>備註</TableHead>
            <TableHead className="w-[50px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading ? (
            <TableRow>
              <TableCell colSpan={7} className="h-24 text-center">
                <div className="flex justify-center">
                  <div className="space-y-2">
                    {[...Array(3)].map((_, i) => (
                      <Skeleton key={i} className="h-12 w-[800px]" />
                    ))}
                  </div>
                </div>
              </TableCell>
            </TableRow>
          ) : bankAccounts.length === 0 ? (
            <TableRow>
              <TableCell
                colSpan={7}
                className="h-24 text-center text-muted-foreground"
              >
                尚無銀行帳戶記錄
              </TableCell>
            </TableRow>
          ) : (
            bankAccounts.map((account) => (
              <TableRow key={account.id}>
                <TableCell className="font-medium">
                  {account.bank_name}
                </TableCell>
                <TableCell>{account.account_type}</TableCell>
                <TableCell className="font-mono">
                  ****{account.account_number_last4}
                </TableCell>
                <TableCell>
                  <Badge variant={getCurrencyColor(account.currency) as any}>
                    {account.currency}
                  </Badge>
                </TableCell>
                <TableCell className="text-right tabular-nums">
                  {account.balance.toLocaleString("zh-TW", {
                    minimumFractionDigits: 2,
                    maximumFractionDigits: 2,
                  })}
                </TableCell>
                <TableCell className="max-w-[200px] truncate text-muted-foreground">
                  {account.note || "-"}
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
                        <DropdownMenuItem onClick={() => onEdit(account)}>
                          <PencilIcon className="mr-2 h-4 w-4" />
                          編輯
                        </DropdownMenuItem>
                      )}
                      {onDelete && (
                        <DropdownMenuItem
                          onClick={() => onDelete(account.id)}
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
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}
