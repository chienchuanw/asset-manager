/**
 * 信用卡列表元件
 * 顯示所有信用卡的列表，支援群組化顯示
 */

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
import { Skeleton } from "@/components/ui/skeleton";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";

import {
  MoreHorizontalIcon,
  PencilIcon,
  TrashIcon,
  ChevronDownIcon,
  ChevronRightIcon,
  FolderIcon,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import type {
  CreditCard,
  CreditCardGroupWithCards,
} from "@/types/user-management";

interface CreditCardListProps {
  creditCards?: CreditCard[];
  groups?: CreditCardGroupWithCards[];
  isLoading?: boolean;
  onEdit?: (creditCard: CreditCard) => void;
  onDelete?: (id: string) => void;
  onEditGroup?: (group: CreditCardGroupWithCards) => void;
  onDeleteGroup?: (id: string) => void;
  onRemoveCardFromGroup?: (groupId: string, cardId: string) => void;
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
  groups = [],
  isLoading = false,
  onEdit,
  onDelete,
  onEditGroup,
  onDeleteGroup,
  onRemoveCardFromGroup,
}: CreditCardListProps) {
  // 取得所有在群組中的卡片 ID
  const cardsInGroups = new Set(
    groups.flatMap((group) => group.cards.map((card) => card.id))
  );

  // 過濾出不在群組中的獨立卡片
  const independentCards = creditCards.filter(
    (card) => !cardsInGroups.has(card.id)
  );

  return (
    <div className="rounded-md border mx-4">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[50px]"></TableHead>
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
              <TableCell colSpan={10} className="h-24 text-center">
                <div className="flex justify-center">
                  <div className="space-y-2">
                    {[...Array(3)].map((_, i) => (
                      <Skeleton key={i} className="h-12 w-[1000px]" />
                    ))}
                  </div>
                </div>
              </TableCell>
            </TableRow>
          ) : creditCards.length === 0 && groups.length === 0 ? (
            <TableRow>
              <TableCell
                colSpan={10}
                className="h-24 text-center text-muted-foreground"
              >
                尚無信用卡記錄
              </TableCell>
            </TableRow>
          ) : (
            <>
              {/* 顯示群組 */}
              {groups.map((group) => (
                <CreditCardGroupRow
                  key={group.id}
                  group={group}
                  onEdit={onEdit}
                  onDelete={onDelete}
                  onEditGroup={onEditGroup}
                  onDeleteGroup={onDeleteGroup}
                  onRemoveCardFromGroup={onRemoveCardFromGroup}
                />
              ))}

              {/* 顯示獨立卡片 */}
              {independentCards.map((card) => {
                return (
                  <CreditCardRow
                    key={card.id}
                    card={card}
                    onEdit={onEdit}
                    onDelete={onDelete}
                  />
                );
              })}
            </>
          )}
        </TableBody>
      </Table>
    </div>
  );
}

/**
 * 信用卡群組行元件
 */
interface CreditCardGroupRowProps {
  group: CreditCardGroupWithCards;
  onEdit?: (creditCard: CreditCard) => void;
  onDelete?: (id: string) => void;
  onEditGroup?: (group: CreditCardGroupWithCards) => void;
  onDeleteGroup?: (id: string) => void;
  onRemoveCardFromGroup?: (groupId: string, cardId: string) => void;
}

function CreditCardGroupRow({
  group,
  onEdit,
  onDelete,
  onEditGroup,
  onDeleteGroup,
  onRemoveCardFromGroup,
}: CreditCardGroupRowProps) {
  const [isOpen, setIsOpen] = useState(true);

  const groupUtilization =
    group.shared_credit_limit > 0
      ? (group.total_used_credit / group.shared_credit_limit) * 100
      : 0;
  const availableCredit = group.shared_credit_limit - group.total_used_credit;

  return (
    <>
      {/* 群組標題行 */}
      <TableRow className="bg-muted/50 hover:bg-muted/70">
        <TableCell>
          <Button
            variant="ghost"
            size="icon"
            className="h-6 w-6"
            onClick={() => setIsOpen(!isOpen)}
          >
            {isOpen ? (
              <ChevronDownIcon className="h-4 w-4" />
            ) : (
              <ChevronRightIcon className="h-4 w-4" />
            )}
          </Button>
        </TableCell>
        <TableCell colSpan={2} className="font-medium">
          <div className="flex items-center gap-2">
            <FolderIcon className="h-4 w-4 text-muted-foreground" />
            <span>{group.name}</span>
            <Badge variant="secondary" className="text-xs">
              {group.cards.length} 張卡片
            </Badge>
          </div>
        </TableCell>
        <TableCell className="text-muted-foreground">
          {group.issuing_bank}
        </TableCell>
        <TableCell colSpan={2} className="text-muted-foreground text-sm">
          共同額度群組
        </TableCell>
        <TableCell className="text-right">
          <div className="space-y-1">
            <div className="tabular-nums">
              {group.shared_credit_limit.toLocaleString("zh-TW", {
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
                  groupUtilization
                )}`}
              >
                {groupUtilization.toFixed(1)}%
              </span>
            </div>
            <Progress value={groupUtilization} className="h-2" />
          </div>
        </TableCell>
        <TableCell className="max-w-[150px] truncate text-muted-foreground">
          {group.note || "-"}
        </TableCell>
        <TableCell>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon">
                <MoreHorizontalIcon className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {onEditGroup && (
                <DropdownMenuItem onClick={() => onEditGroup(group)}>
                  <PencilIcon className="mr-2 h-4 w-4" />
                  編輯群組
                </DropdownMenuItem>
              )}
              {onDeleteGroup && (
                <DropdownMenuItem
                  onClick={() => onDeleteGroup(group.id)}
                  className="text-destructive"
                >
                  <TrashIcon className="mr-2 h-4 w-4" />
                  解散群組
                </DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </TableCell>
      </TableRow>

      {/* 群組內的卡片 */}
      {isOpen &&
        group.cards.map((card) => (
          <CreditCardRow
            key={card.id}
            card={card}
            isInGroup={true}
            groupId={group.id}
            onEdit={onEdit}
            onDelete={onDelete}
            onRemoveFromGroup={onRemoveCardFromGroup}
          />
        ))}
    </>
  );
}

/**
 * 信用卡行元件
 */
interface CreditCardRowProps {
  card: CreditCard;
  isInGroup?: boolean;
  groupId?: string;
  onEdit?: (creditCard: CreditCard) => void;
  onDelete?: (id: string) => void;
  onRemoveFromGroup?: (groupId: string, cardId: string) => void;
}

function CreditCardRow({
  card,
  isInGroup = false,
  groupId,
  onEdit,
  onDelete,
  onRemoveFromGroup,
}: CreditCardRowProps) {
  const utilization = calculateUtilization(card);
  const availableCredit = getAvailableCredit(card);

  return (
    <TableRow className={isInGroup ? "bg-muted/20" : ""}>
      <TableCell>{isInGroup && <div className="ml-6" />}</TableCell>
      <TableCell className="font-medium">{card.issuing_bank}</TableCell>
      <TableCell>{card.card_name}</TableCell>
      <TableCell className="font-mono">****{card.card_number_last4}</TableCell>
      <TableCell>每月 {card.billing_day} 日</TableCell>
      <TableCell>每月 {card.payment_due_day} 日</TableCell>
      <TableCell className="text-right">
        {isInGroup ? (
          <div className="text-sm text-muted-foreground">
            已用:{" "}
            {card.used_credit.toLocaleString("zh-TW", {
              minimumFractionDigits: 0,
              maximumFractionDigits: 0,
            })}
          </div>
        ) : (
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
        )}
      </TableCell>
      <TableCell>
        {!isInGroup && (
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
        )}
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
            {isInGroup && onRemoveFromGroup && groupId && (
              <DropdownMenuItem
                onClick={() => onRemoveFromGroup(groupId, card.id)}
              >
                <TrashIcon className="mr-2 h-4 w-4" />
                從群組移除
              </DropdownMenuItem>
            )}
            {!isInGroup && onDelete && (
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
}
