"use client";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { PaymentMethodType } from "@/types/cash-flow";
import { useBankAccounts, useCreditCards } from "@/hooks";
import { Loader2, AlertTriangle } from "lucide-react";
import { formatAmount } from "@/types/cash-flow";

interface AccountSelectProps {
  value: string;
  onValueChange: (value: string) => void;
  paymentMethodType: PaymentMethodType;
  placeholder?: string;
  disabled?: boolean;
}

/**
 * 帳戶選擇器元件
 *
 * 根據付款方式類型顯示對應的銀行帳戶或信用卡選項
 * 顯示格式：台灣銀行 活存 (****1234) - 餘額: $50,000
 * 信用卡顯示可用額度，當額度為負數時顯示紅色警告
 */
export function AccountSelect({
  value,
  onValueChange,
  paymentMethodType,
  placeholder = "選擇帳戶",
  disabled = false,
}: AccountSelectProps) {
  // 取得銀行帳戶和信用卡列表
  const { data: bankAccounts, isLoading: isLoadingBankAccounts } = useBankAccounts();
  const { data: creditCards, isLoading: isLoadingCreditCards } = useCreditCards();

  // 根據付款方式類型決定是否顯示載入狀態
  const isLoading = 
    (paymentMethodType === PaymentMethodType.BANK_ACCOUNT && isLoadingBankAccounts) ||
    (paymentMethodType === PaymentMethodType.CREDIT_CARD && isLoadingCreditCards);

  // 如果是現金付款，不顯示選擇器
  if (paymentMethodType === PaymentMethodType.CASH) {
    return null;
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-10 border rounded-md bg-muted">
        <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
      </div>
    );
  }

  // 格式化銀行帳戶顯示
  const formatBankAccountDisplay = (account: any) => {
    return `${account.bank_name} ${account.account_type} (****${account.account_number_last4}) - 餘額: $${formatAmount(account.balance)}`;
  };

  // 格式化信用卡顯示
  const formatCreditCardDisplay = (card: any) => {
    const availableCredit = card.credit_limit - card.used_credit;
    const isOverLimit = availableCredit < 0;
    
    return {
      text: `${card.issuing_bank} ${card.card_name} (****${card.card_number_last4}) - 可用額度: $${formatAmount(availableCredit)}`,
      isOverLimit,
    };
  };

  return (
    <Select value={value} onValueChange={onValueChange} disabled={disabled}>
      <SelectTrigger>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {paymentMethodType === PaymentMethodType.BANK_ACCOUNT && (
          <>
            {bankAccounts?.map((account) => (
              <SelectItem key={account.id} value={account.id}>
                <div className="flex items-center gap-2">
                  <span className="text-sm">
                    {formatBankAccountDisplay(account)}
                  </span>
                </div>
              </SelectItem>
            ))}
            {(!bankAccounts || bankAccounts.length === 0) && (
              <SelectItem value="" disabled>
                <span className="text-muted-foreground">尚未建立銀行帳戶</span>
              </SelectItem>
            )}
          </>
        )}

        {paymentMethodType === PaymentMethodType.CREDIT_CARD && (
          <>
            {creditCards?.map((card) => {
              const displayInfo = formatCreditCardDisplay(card);
              return (
                <SelectItem key={card.id} value={card.id}>
                  <div className="flex items-center gap-2">
                    {displayInfo.isOverLimit && (
                      <AlertTriangle className="h-4 w-4 text-red-500" />
                    )}
                    <span 
                      className={`text-sm ${
                        displayInfo.isOverLimit ? "text-red-600" : ""
                      }`}
                    >
                      {displayInfo.text}
                    </span>
                  </div>
                </SelectItem>
              );
            })}
            {(!creditCards || creditCards.length === 0) && (
              <SelectItem value="" disabled>
                <span className="text-muted-foreground">尚未建立信用卡</span>
              </SelectItem>
            )}
          </>
        )}
      </SelectContent>
    </Select>
  );
}
