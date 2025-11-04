"use client";

import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { useBankAccounts, useCreditCards } from "@/hooks";
import { type CashFlow, SourceType } from "@/types/cash-flow";
import { Banknote, CreditCard, Landmark, AlertCircle } from "lucide-react";

interface PaymentMethodDisplayProps {
  cashFlow: CashFlow;
}

/**
 * 付款方式顯示組件
 *
 * 根據現金流記錄的 source_type 和 source_id 顯示對應的付款方式資訊
 */
export function PaymentMethodDisplay({ cashFlow }: PaymentMethodDisplayProps) {
  const { data: bankAccounts } = useBankAccounts();
  const { data: creditCards } = useCreditCards();

  return (
    <TooltipProvider>
      {/* 如果沒有 source_type，顯示現金 */}
      {(!cashFlow.source_type ||
        cashFlow.source_type === SourceType.MANUAL) && (
        <Tooltip>
          <TooltipTrigger asChild>
            <Badge variant="outline" className="flex items-center gap-1">
              <Banknote className="h-3 w-3" />
              現金
            </Badge>
          </TooltipTrigger>
          <TooltipContent>
            <p>現金交易</p>
          </TooltipContent>
        </Tooltip>
      )}

      {/* 銀行帳戶交易 */}
      {cashFlow.source_type === SourceType.BANK_ACCOUNT &&
        cashFlow.source_id &&
        (() => {
          const bankAccount = bankAccounts?.find(
            (account) => account.id === cashFlow.source_id
          );

          if (!bankAccount) {
            return (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Badge
                    variant="destructive"
                    className="flex items-center gap-1"
                  >
                    <AlertCircle className="h-3 w-3" />
                    帳戶不存在
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>
                  <p>關聯的銀行帳戶已被刪除</p>
                </TooltipContent>
              </Tooltip>
            );
          }

          return (
            <Tooltip>
              <TooltipTrigger asChild>
                <Badge variant="outline" className="flex items-center gap-1">
                  <Landmark className="h-3 w-3" />
                  {bankAccount.bank_name}
                </Badge>
              </TooltipTrigger>
              <TooltipContent>
                <div className="space-y-1">
                  <p className="font-medium">{bankAccount.bank_name}</p>
                  <p className="text-sm text-muted-foreground">
                    {bankAccount.account_type} (****
                    {bankAccount.account_number_last4})
                  </p>
                  <p className="text-sm">
                    餘額: ${bankAccount.balance.toLocaleString()}{" "}
                    {bankAccount.currency}
                  </p>
                </div>
              </TooltipContent>
            </Tooltip>
          );
        })()}

      {/* 信用卡交易 */}
      {cashFlow.source_type === SourceType.CREDIT_CARD &&
        cashFlow.source_id &&
        (() => {
          const creditCard = creditCards?.find(
            (card) => card.id === cashFlow.source_id
          );

          if (!creditCard) {
            return (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Badge
                    variant="destructive"
                    className="flex items-center gap-1"
                  >
                    <AlertCircle className="h-3 w-3" />
                    信用卡不存在
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>
                  <p>關聯的信用卡已被刪除</p>
                </TooltipContent>
              </Tooltip>
            );
          }

          const availableCredit =
            creditCard.credit_limit - creditCard.used_credit;
          const isOverLimit = availableCredit < 0;

          return (
            <Tooltip>
              <TooltipTrigger asChild>
                <Badge variant="outline" className="flex items-center gap-1">
                  <CreditCard className="h-3 w-3" />
                  {creditCard.issuing_bank}
                </Badge>
              </TooltipTrigger>
              <TooltipContent>
                <div className="space-y-1">
                  <p className="font-medium">{creditCard.issuing_bank}</p>
                  <p className="text-sm text-muted-foreground">
                    {creditCard.card_name} (****{creditCard.card_number_last4})
                  </p>
                  <p className="text-sm">
                    已使用: ${creditCard.used_credit.toLocaleString()}
                  </p>
                  <p className={`text-sm ${isOverLimit ? "text-red-600" : ""}`}>
                    可用額度: ${availableCredit.toLocaleString()}
                  </p>
                </div>
              </TooltipContent>
            </Tooltip>
          );
        })()}

      {/* 訂閱或分期自動產生 */}
      {(cashFlow.source_type === SourceType.SUBSCRIPTION ||
        cashFlow.source_type === SourceType.INSTALLMENT) &&
        (() => {
          const isSubscription =
            cashFlow.source_type === SourceType.SUBSCRIPTION;

          return (
            <Tooltip>
              <TooltipTrigger asChild>
                <Badge variant="secondary" className="flex items-center gap-1">
                  <AlertCircle className="h-3 w-3" />
                  {isSubscription ? "訂閱" : "分期"}
                </Badge>
              </TooltipTrigger>
              <TooltipContent>
                <p>{isSubscription ? "訂閱自動產生" : "分期自動產生"}</p>
              </TooltipContent>
            </Tooltip>
          );
        })()}

      {/* 未知類型 */}
      {cashFlow.source_type &&
        cashFlow.source_type !== SourceType.MANUAL &&
        cashFlow.source_type !== SourceType.BANK_ACCOUNT &&
        cashFlow.source_type !== SourceType.CREDIT_CARD &&
        cashFlow.source_type !== SourceType.SUBSCRIPTION &&
        cashFlow.source_type !== SourceType.INSTALLMENT && (
          <Tooltip>
            <TooltipTrigger asChild>
              <Badge variant="outline" className="flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                未知
              </Badge>
            </TooltipTrigger>
            <TooltipContent>
              <p>未知的付款方式</p>
            </TooltipContent>
          </Tooltip>
        )}
    </TooltipProvider>
  );
}
