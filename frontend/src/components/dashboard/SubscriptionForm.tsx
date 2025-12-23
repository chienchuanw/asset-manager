/**
 * 訂閱表單元件
 * 用於建立和編輯訂閱
 */

"use client";

import { useEffect } from "react";
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
import { Switch } from "@/components/ui/switch";
import type {
  Subscription,
  CreateSubscriptionInput,
  PaymentMethod,
} from "@/types/subscription";
import type { CashFlowCategory } from "@/types/cash-flow";
import { useBankAccounts } from "@/hooks/useBankAccounts";
import { useCreditCards } from "@/hooks/useCreditCards";

interface SubscriptionFormProps {
  subscription?: Subscription;
  categories?: CashFlowCategory[];
  onSubmit: (data: CreateSubscriptionInput) => void;
  onCancel?: () => void;
  isSubmitting?: boolean;
}

export function SubscriptionForm({
  subscription,
  categories = [],
  onSubmit,
  onCancel,
  isSubmitting = false,
}: SubscriptionFormProps) {
  const t = useTranslations("recurring");
  const tCommon = useTranslations("common");

  // 取得銀行帳戶和信用卡列表
  const { data: bankAccounts = [] } = useBankAccounts();
  const { data: creditCards = [] } = useCreditCards();

  // 根據是否為編輯模式設定初始值
  const getDefaultValues = (): CreateSubscriptionInput => {
    if (subscription) {
      return {
        name: subscription.name,
        amount: subscription.amount,
        currency: subscription.currency,
        billing_cycle: subscription.billing_cycle,
        billing_day: subscription.billing_day,
        category_id: subscription.category_id,
        payment_method: subscription.payment_method || "cash",
        account_id: subscription.account_id,
        start_date: subscription.start_date.split("T")[0],
        end_date: subscription.end_date?.split("T")[0],
        auto_renew: subscription.auto_renew,
        note: subscription.note || "",
      };
    }
    return {
      name: "",
      amount: 0,
      currency: "TWD",
      billing_cycle: "monthly",
      billing_day: 1,
      category_id: "",
      payment_method: "cash",
      account_id: undefined,
      start_date: new Date().toISOString().split("T")[0],
      auto_renew: true,
      note: "",
    };
  };

  const form = useForm<CreateSubscriptionInput>({
    defaultValues: getDefaultValues(),
  });

  // 監聽付款方式變化
  const paymentMethod = useWatch({
    control: form.control,
    name: "payment_method",
  });

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
          rules={{ required: t("serviceNameRequired") }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("serviceName")}</FormLabel>
              <FormControl>
                <Input placeholder={t("serviceNamePlaceholder")} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 金額 */}
        <FormField
          control={form.control}
          name="amount"
          rules={{
            required: t("amountRequired"),
            min: { value: 0.01, message: t("amountPositive") },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("amount")}</FormLabel>
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

        {/* 計費週期 */}
        <FormField
          control={form.control}
          name="billing_cycle"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("billingCycle")}</FormLabel>
              <Select onValueChange={field.onChange} value={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder={t("selectCycle")} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="monthly">{t("monthly")}</SelectItem>
                  <SelectItem value="quarterly">{t("quarterly")}</SelectItem>
                  <SelectItem value="yearly">{t("yearly")}</SelectItem>
                </SelectContent>
              </Select>
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

        {/* 結束日期（可選） */}
        <FormField
          control={form.control}
          name="end_date"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("endDateOptional")}</FormLabel>
              <FormControl>
                <Input type="date" {...field} />
              </FormControl>
              <FormDescription>{t("endDateDesc")}</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 自動續約 */}
        <FormField
          control={form.control}
          name="auto_renew"
          render={({ field }) => (
            <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
              <div className="space-y-0.5">
                <FormLabel className="text-base">{t("autoRenew")}</FormLabel>
                <FormDescription>{t("autoRenewDesc")}</FormDescription>
              </div>
              <FormControl>
                <Switch
                  checked={field.value}
                  onCheckedChange={field.onChange}
                />
              </FormControl>
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
              : subscription
              ? tCommon("update")
              : tCommon("create")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
