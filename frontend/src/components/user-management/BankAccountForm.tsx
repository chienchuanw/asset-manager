/**
 * 銀行帳戶表單元件
 * 用於建立和編輯銀行帳戶
 */

"use client";

import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
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
import type {
  BankAccount,
  CreateBankAccountInput,
} from "@/types/user-management";

interface BankAccountFormProps {
  bankAccount?: BankAccount;
  onSubmit: (data: CreateBankAccountInput) => void;
  onCancel?: () => void;
  isSubmitting?: boolean;
}

export function BankAccountForm({
  bankAccount,
  onSubmit,
  onCancel,
  isSubmitting = false,
}: BankAccountFormProps) {
  const t = useTranslations("userManagement");
  const tCommon = useTranslations("common");

  const form = useForm<CreateBankAccountInput>({
    defaultValues: {
      bank_name: "",
      account_type: "",
      account_number_last4: "",
      currency: "TWD",
      balance: 0,
      note: "",
    },
  });

  // 如果是編輯模式，填入現有資料
  useEffect(() => {
    if (bankAccount) {
      form.reset({
        bank_name: bankAccount.bank_name,
        account_type: bankAccount.account_type,
        account_number_last4: bankAccount.account_number_last4,
        currency: bankAccount.currency,
        balance: bankAccount.balance,
        note: bankAccount.note || "",
      });
    }
  }, [bankAccount, form]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        {/* 銀行名稱 */}
        <FormField
          control={form.control}
          name="bank_name"
          rules={{ required: t("bankNameRequired") }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("bankName")}</FormLabel>
              <FormControl>
                <Input placeholder={t("bankNamePlaceholder")} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 帳戶類型 */}
        <FormField
          control={form.control}
          name="account_type"
          rules={{ required: t("accountTypeRequired") }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("accountType")}</FormLabel>
              <FormControl>
                <Input placeholder={t("accountTypePlaceholder")} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 帳號後四碼 */}
        <FormField
          control={form.control}
          name="account_number_last4"
          rules={{
            required: t("accountLast4Required"),
            pattern: {
              value: /^\d{4}$/,
              message: t("accountLast4Invalid"),
            },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("accountLast4")}</FormLabel>
              <FormControl>
                <Input
                  placeholder="1234"
                  maxLength={4}
                  {...field}
                  onChange={(e) => {
                    const value = e.target.value.replace(/\D/g, "");
                    field.onChange(value);
                  }}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 幣別 */}
        <FormField
          control={form.control}
          name="currency"
          rules={{ required: t("currencyRequired") }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("currency")}</FormLabel>
              <Select onValueChange={field.onChange} value={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder={t("selectCurrency")} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="TWD">TWD (新台幣)</SelectItem>
                  <SelectItem value="USD">USD (美元)</SelectItem>
                  <SelectItem value="JPY">JPY (日圓)</SelectItem>
                  <SelectItem value="EUR">EUR (歐元)</SelectItem>
                  <SelectItem value="CNY">CNY (人民幣)</SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 餘額 */}
        <FormField
          control={form.control}
          name="balance"
          rules={{
            required: t("balanceRequired"),
            min: { value: 0, message: t("balanceNonNegative") },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("balance")}</FormLabel>
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
        <div className="flex justify-end gap-3">
          {onCancel && (
            <Button type="button" variant="outline" onClick={onCancel}>
              {tCommon("cancel")}
            </Button>
          )}
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting
              ? tCommon("saving")
              : bankAccount
              ? tCommon("update")
              : tCommon("create")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
