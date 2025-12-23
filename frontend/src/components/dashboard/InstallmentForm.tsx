/**
 * 分期表單元件
 * 用於建立和編輯分期
 */

"use client";

import { useEffect, useState } from "react";
import { useForm, useWatch } from "react-hook-form";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import type { Installment, CreateInstallmentInput } from "@/types/installment";
import type { CashFlowCategory } from "@/types/cash-flow";
import type { PaymentMethod } from "@/types/subscription";
import { useBankAccounts } from "@/hooks/useBankAccounts";
import { useCreditCards } from "@/hooks/useCreditCards";

interface InstallmentFormProps {
  installment?: Installment;
  categories?: CashFlowCategory[];
  onSubmit: (data: CreateInstallmentInput) => void;
  onCancel?: () => void;
  isSubmitting?: boolean;
}

export function InstallmentForm({
  installment,
  categories = [],
  onSubmit,
  onCancel,
  isSubmitting = false,
}: InstallmentFormProps) {
  const t = useTranslations("recurring");
  const tCommon = useTranslations("common");

  // 取得銀行帳戶和信用卡列表
  const { data: bankAccounts = [] } = useBankAccounts();
  const { data: creditCards = [] } = useCreditCards();

  // 根據是否為編輯模式設定初始值
  const getDefaultValues = (): CreateInstallmentInput => {
    if (installment) {
      return {
        name: installment.name,
        total_amount: installment.total_amount,
        currency: installment.currency,
        installment_count: installment.installment_count,
        interest_rate: installment.interest_rate,
        category_id: installment.category_id,
        payment_method: installment.payment_method || "cash",
        account_id: installment.account_id,
        start_date: installment.start_date.split("T")[0],
        billing_day: installment.billing_day,
        note: installment.note || "",
      };
    }
    // 使用本地時間格式化日期，避免時區轉換問題
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, "0");
    const day = String(now.getDate()).padStart(2, "0");
    return {
      name: "",
      total_amount: 0,
      currency: "TWD",
      installment_count: 12,
      interest_rate: 0,
      category_id: "",
      payment_method: "cash",
      account_id: undefined,
      start_date: `${year}-${month}-${day}`,
      billing_day: 1,
      note: "",
    };
  };

  const form = useForm<CreateInstallmentInput>({
    defaultValues: getDefaultValues(),
  });

  // 監聽付款方式變化
  const paymentMethod = useWatch({
    control: form.control,
    name: "payment_method",
  });

  // 監聽表單變化以計算每期金額和總利息
  const [calculatedValues, setCalculatedValues] = useState({
    amountPerInstallment: 0,
    totalInterest: 0,
    totalWithInterest: 0,
  });

  const totalAmount = form.watch("total_amount");
  const installmentCount = form.watch("installment_count");
  const interestRate = form.watch("interest_rate");

  useEffect(() => {
    if (totalAmount > 0 && installmentCount > 0) {
      // 計算總利息（簡單利息）
      const totalInterest = totalAmount * (interestRate / 100);
      const totalWithInterest = totalAmount + totalInterest;
      const amountPerInstallment = totalWithInterest / installmentCount;

      setCalculatedValues({
        amountPerInstallment,
        totalInterest,
        totalWithInterest,
      });
    }
  }, [totalAmount, installmentCount, interestRate]);

  // 篩選支出類別
  const expenseCategories = categories.filter((c) => c.type === "expense");

  // 根據付款方式取得帳戶選項
  const getAccountOptions = () => {
    if (paymentMethod === "bank_account") {
      return bankAccounts.map((account) => ({
        id: account.id,
        name: `${account.bank_name} - ${account.account_type} (${account.account_number_last4})`,
      }));
    }
    if (paymentMethod === "credit_card") {
      return creditCards.map((card) => ({
        id: card.id,
        name: `${card.issuing_bank} - ${card.card_name}`,
      }));
    }
    return [];
  };

  const accountOptions = getAccountOptions();

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        {/* 名稱 */}
        <FormField
          control={form.control}
          name="name"
          rules={{ required: t("installmentNameRequired") }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("installmentName")}</FormLabel>
              <FormControl>
                <Input
                  placeholder={t("installmentNamePlaceholder")}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 總金額 */}
        <FormField
          control={form.control}
          name="total_amount"
          rules={{
            required: t("totalAmountRequired"),
            min: { value: 0.01, message: t("totalAmountPositive") },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("totalAmountLabel")}</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  step="0.01"
                  placeholder="0.00"
                  {...field}
                  onChange={(e) => {
                    const value = parseFloat(e.target.value);
                    field.onChange(isNaN(value) ? "" : value);
                  }}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 分期期數 */}
        <FormField
          control={form.control}
          name="installment_count"
          rules={{
            required: t("installmentCountRequired"),
            min: { value: 2, message: t("installmentCountMin") },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("installmentCount")}</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  min="2"
                  placeholder="12"
                  {...field}
                  onChange={(e) => {
                    const value = parseInt(e.target.value);
                    field.onChange(isNaN(value) ? "" : value);
                  }}
                />
              </FormControl>
              <FormDescription>{t("installmentCountDesc")}</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 利率 */}
        <FormField
          control={form.control}
          name="interest_rate"
          rules={{
            min: { value: 0, message: t("interestRateNonNegative") },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("annualInterestRate")}</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  step="0.01"
                  placeholder="0.00"
                  {...field}
                  onChange={(e) => {
                    const value = parseFloat(e.target.value);
                    field.onChange(isNaN(value) ? "" : value);
                  }}
                />
              </FormControl>
              <FormDescription>{t("interestRateDesc")}</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 計算結果顯示 */}
        {calculatedValues.amountPerInstallment > 0 && (
          <div className="rounded-lg border p-4 space-y-2 bg-muted/50">
            <p className="text-sm font-medium">{t("calculationResult")}</p>
            <div className="grid grid-cols-2 gap-2 text-sm">
              <div>
                <p className="text-muted-foreground">{t("monthlyPayment")}</p>
                <p className="font-semibold tabular-nums">
                  TWD{" "}
                  {calculatedValues.amountPerInstallment.toLocaleString(
                    "zh-TW",
                    { maximumFractionDigits: 2 }
                  )}
                </p>
              </div>
              <div>
                <p className="text-muted-foreground">{t("totalInterest")}</p>
                <p className="font-semibold tabular-nums">
                  TWD{" "}
                  {calculatedValues.totalInterest.toLocaleString("zh-TW", {
                    maximumFractionDigits: 2,
                  })}
                </p>
              </div>
              <div className="col-span-2">
                <p className="text-muted-foreground">
                  {t("totalPaymentAmount")}
                </p>
                <p className="font-semibold tabular-nums">
                  TWD{" "}
                  {calculatedValues.totalWithInterest.toLocaleString("zh-TW", {
                    maximumFractionDigits: 2,
                  })}
                </p>
              </div>
            </div>
          </div>
        )}

        {/* 分類 */}
        <FormField
          control={form.control}
          name="category_id"
          rules={{ required: t("categoryRequired") }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("category")}</FormLabel>
              <Select onValueChange={field.onChange} value={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder={t("selectCategory")} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {expenseCategories.map((category) => (
                    <SelectItem key={category.id} value={category.id}>
                      {category.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 付款方式 */}
        <FormField
          control={form.control}
          name="payment_method"
          rules={{ required: t("paymentMethodRequired") }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("paymentMethod")}</FormLabel>
              <Select
                onValueChange={(value: PaymentMethod) => {
                  field.onChange(value);
                  // 切換付款方式時清除帳戶選擇
                  if (value === "cash") {
                    form.setValue("account_id", undefined);
                  }
                }}
                value={field.value}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder={t("selectPaymentMethod")} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="cash">{t("cash")}</SelectItem>
                  <SelectItem value="bank_account">
                    {t("bankAccount")}
                  </SelectItem>
                  <SelectItem value="credit_card">{t("creditCard")}</SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 帳戶選擇（當付款方式為銀行帳戶或信用卡時顯示） */}
        {paymentMethod !== "cash" && accountOptions.length > 0 && (
          <FormField
            control={form.control}
            name="account_id"
            rules={{
              required:
                paymentMethod === "bank_account"
                  ? t("bankAccountRequired")
                  : t("creditCardRequired"),
            }}
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  {paymentMethod === "bank_account"
                    ? t("bankAccount")
                    : t("creditCard")}
                </FormLabel>
                <Select
                  onValueChange={field.onChange}
                  value={field.value || ""}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue
                        placeholder={
                          paymentMethod === "bank_account"
                            ? t("selectBankAccount")
                            : t("selectCreditCard")
                        }
                      />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {accountOptions.map((account) => (
                      <SelectItem key={account.id} value={account.id}>
                        {account.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
        )}

        {/* 開始日期 */}
        <FormField
          control={form.control}
          name="start_date"
          rules={{ required: t("startDateRequired") }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("startDate")}</FormLabel>
              <FormControl>
                <Input type="date" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 扣款日 */}
        <FormField
          control={form.control}
          name="billing_day"
          rules={{
            required: t("billingDayRequired"),
            min: { value: 1, message: t("billingDayRange") },
            max: { value: 31, message: t("billingDayRange") },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("billingDayLabel")}</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  min="1"
                  max="31"
                  placeholder="1"
                  {...field}
                  onChange={(e) => {
                    const value = parseInt(e.target.value);
                    field.onChange(isNaN(value) ? "" : value);
                  }}
                />
              </FormControl>
              <FormDescription>{t("billingDayDesc")}</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 備註 */}
        <FormField
          control={form.control}
          name="note"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{tCommon("noteOptional")}</FormLabel>
              <FormControl>
                <Textarea
                  placeholder={tCommon("enterNote")}
                  className="resize-none"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 按鈕 */}
        <div className="flex gap-2 justify-end">
          {onCancel && (
            <Button type="button" variant="outline" onClick={onCancel}>
              {tCommon("cancel")}
            </Button>
          )}
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting
              ? tCommon("saving")
              : installment
              ? tCommon("update")
              : tCommon("create")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
