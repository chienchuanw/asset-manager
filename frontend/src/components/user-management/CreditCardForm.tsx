/**
 * 信用卡表單元件
 * 用於建立和編輯信用卡
 */

"use client";

import { useEffect } from "react";
import { useForm } from "react-hook-form";
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
import { Textarea } from "@/components/ui/textarea";
import type {
  CreditCard,
  CreateCreditCardInput,
} from "@/types/user-management";

interface CreditCardFormProps {
  creditCard?: CreditCard;
  onSubmit: (data: CreateCreditCardInput) => void;
  onCancel?: () => void;
  isSubmitting?: boolean;
}

export function CreditCardForm({
  creditCard,
  onSubmit,
  onCancel,
  isSubmitting = false,
}: CreditCardFormProps) {
  const form = useForm<CreateCreditCardInput>({
    defaultValues: {
      issuing_bank: "",
      card_name: "",
      card_number_last4: "",
      billing_day: 1,
      payment_due_day: 1,
      credit_limit: 0,
      used_credit: 0,
      note: "",
    },
  });

  // 如果是編輯模式，填入現有資料
  useEffect(() => {
    if (creditCard) {
      form.reset({
        issuing_bank: creditCard.issuing_bank,
        card_name: creditCard.card_name,
        card_number_last4: creditCard.card_number_last4,
        billing_day: creditCard.billing_day,
        payment_due_day: creditCard.payment_due_day,
        credit_limit: creditCard.credit_limit,
        used_credit: creditCard.used_credit,
        note: creditCard.note || "",
      });
    }
  }, [creditCard, form]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        {/* 發卡銀行 */}
        <FormField
          control={form.control}
          name="issuing_bank"
          rules={{ required: "請輸入發卡銀行" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>發卡銀行</FormLabel>
              <FormControl>
                <Input placeholder="例如：玉山銀行" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 卡片名稱 */}
        <FormField
          control={form.control}
          name="card_name"
          rules={{ required: "請輸入卡片名稱" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>卡片名稱</FormLabel>
              <FormControl>
                <Input placeholder="例如：Pi 拍錢包信用卡" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 卡號後四碼 */}
        <FormField
          control={form.control}
          name="card_number_last4"
          rules={{
            required: "請輸入卡號後四碼",
            pattern: {
              value: /^\d{4}$/,
              message: "請輸入正確的四位數字",
            },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>卡號後四碼</FormLabel>
              <FormControl>
                <Input
                  placeholder="5678"
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

        <div className="grid grid-cols-2 gap-4">
          {/* 帳單日 */}
          <FormField
            control={form.control}
            name="billing_day"
            rules={{
              required: "請輸入帳單日",
              min: { value: 1, message: "帳單日必須在 1-31 之間" },
              max: { value: 31, message: "帳單日必須在 1-31 之間" },
            }}
            render={({ field }) => (
              <FormItem>
                <FormLabel>帳單日</FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    min={1}
                    max={31}
                    placeholder="15"
                    {...field}
                    onChange={(e) => {
                      const value = parseInt(e.target.value);
                      field.onChange(isNaN(value) ? "" : value);
                    }}
                  />
                </FormControl>
                <FormDescription>每月帳單結算日（1-31）</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          {/* 繳款截止日 */}
          <FormField
            control={form.control}
            name="payment_due_day"
            rules={{
              required: "請輸入繳款截止日",
              min: { value: 1, message: "繳款截止日必須在 1-31 之間" },
              max: { value: 31, message: "繳款截止日必須在 1-31 之間" },
            }}
            render={({ field }) => (
              <FormItem>
                <FormLabel>繳款截止日</FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    min={1}
                    max={31}
                    placeholder="5"
                    {...field}
                    onChange={(e) => {
                      const value = parseInt(e.target.value);
                      field.onChange(isNaN(value) ? "" : value);
                    }}
                  />
                </FormControl>
                <FormDescription>每月繳款截止日（1-31）</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          {/* 信用額度 */}
          <FormField
            control={form.control}
            name="credit_limit"
            rules={{
              required: "請輸入信用額度",
              min: { value: 0.01, message: "信用額度必須大於 0" },
            }}
            render={({ field }) => (
              <FormItem>
                <FormLabel>信用額度</FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    step="0.01"
                    placeholder="100000.00"
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

          {/* 已使用額度 */}
          <FormField
            control={form.control}
            name="used_credit"
            rules={{
              required: "請輸入已使用額度",
              min: { value: 0, message: "已使用額度不能為負數" },
              validate: (value) => {
                const creditLimit = form.getValues("credit_limit");
                if (value > creditLimit) {
                  return "已使用額度不能超過信用額度";
                }
                return true;
              },
            }}
            render={({ field }) => (
              <FormItem>
                <FormLabel>已使用額度</FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    step="0.01"
                    placeholder="25000.00"
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
        </div>

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
            {isSubmitting ? "儲存中..." : creditCard ? "更新" : "新增"}
          </Button>
        </div>
      </form>
    </Form>
  );
}

