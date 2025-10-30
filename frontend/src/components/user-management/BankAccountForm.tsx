/**
 * 銀行帳戶表單元件
 * 用於建立和編輯銀行帳戶
 */

"use client";

import { useEffect } from "react";
import { useForm } from "react-hook-form";
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
          rules={{ required: "請輸入銀行名稱" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>銀行名稱</FormLabel>
              <FormControl>
                <Input placeholder="例如：台灣銀行" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 帳戶類型 */}
        <FormField
          control={form.control}
          name="account_type"
          rules={{ required: "請輸入帳戶類型" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>帳戶類型</FormLabel>
              <FormControl>
                <Input placeholder="例如：活存、定存、外幣帳戶" {...field} />
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
            required: "請輸入帳號後四碼",
            pattern: {
              value: /^\d{4}$/,
              message: "請輸入正確的四位數字",
            },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>帳號後四碼</FormLabel>
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
          rules={{ required: "請選擇幣別" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>幣別</FormLabel>
              <Select onValueChange={field.onChange} value={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="選擇幣別" />
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
            required: "請輸入餘額",
            min: { value: 0, message: "餘額不能為負數" },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>餘額</FormLabel>
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
              <FormLabel>備註（選填）</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="輸入備註..."
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
              取消
            </Button>
          )}
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "儲存中..." : bankAccount ? "更新" : "新增"}
          </Button>
        </div>
      </form>
    </Form>
  );
}

