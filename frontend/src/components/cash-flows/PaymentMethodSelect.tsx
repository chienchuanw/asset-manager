"use client";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  PaymentMethodType,
  getPaymentMethodTypeOptions,
  getPaymentMethodTypeLabel,
} from "@/types/cash-flow";
import { useBankAccounts, useCreditCards } from "@/hooks";
import { CreditCard, Wallet, Building2 } from "lucide-react";

interface PaymentMethodSelectProps {
  value?: PaymentMethodType;
  onValueChange: (value: PaymentMethodType) => void;
  placeholder?: string;
  disabled?: boolean;
  excludeCash?: boolean; // 是否排除現金選項
}

/**
 * 付款方式選擇器元件
 *
 * 提供現金、銀行帳戶、信用卡三種付款方式選擇
 * 如果使用者沒有建立任何帳戶或信用卡，則不提供對應選項
 */
export function PaymentMethodSelect({
  value,
  onValueChange,
  placeholder = "選擇付款方式",
  disabled = false,
  excludeCash = false,
}: PaymentMethodSelectProps) {
  // 取得銀行帳戶和信用卡列表以判斷是否有可用選項
  const { data: bankAccounts } = useBankAccounts();
  const { data: creditCards } = useCreditCards();

  // 根據使用者是否有建立帳戶來決定可用選項
  const availableOptions = getPaymentMethodTypeOptions().filter((option) => {
    if (option.value === PaymentMethodType.CASH) {
      return !excludeCash; // 如果 excludeCash 為 true，則排除現金選項
    }
    if (option.value === PaymentMethodType.BANK_ACCOUNT) {
      return bankAccounts && bankAccounts.length > 0;
    }
    if (option.value === PaymentMethodType.CREDIT_CARD) {
      return creditCards && creditCards.length > 0;
    }
    return false;
  });

  // 取得付款方式的圖示
  const getPaymentMethodIcon = (paymentMethodType: PaymentMethodType) => {
    switch (paymentMethodType) {
      case PaymentMethodType.CASH:
        return <Wallet className="h-4 w-4" />;
      case PaymentMethodType.BANK_ACCOUNT:
        return <Building2 className="h-4 w-4" />;
      case PaymentMethodType.CREDIT_CARD:
        return <CreditCard className="h-4 w-4" />;
      default:
        return null;
    }
  };

  return (
    <Select value={value} onValueChange={onValueChange} disabled={disabled}>
      <SelectTrigger>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {availableOptions.map((option) => (
          <SelectItem key={option.value} value={option.value}>
            <div className="flex items-center gap-2">
              {getPaymentMethodIcon(option.value)}
              <span>{option.label}</span>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
