/**
 * 信用卡表單元件
 * 用於建立和編輯信用卡
 */

"use client";

import { useEffect } from "react";
import { useForm } from "react-hook-form";
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
  const t = useTranslations("userManagement");
  const tCommon = useTranslations("common");

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
          rules={{ required: t("issuingBankRequired") }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("issuingBank")}</FormLabel>
              <FormControl>
                <Input placeholder={t("issuingBankPlaceholder")} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 卡片名稱 */}
        <FormField
          control={form.control}
          name="card_name"
          rules={{ required: t("cardNameRequired") }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("cardName")}</FormLabel>
              <FormControl>
                <Input placeholder={t("cardNamePlaceholder")} {...field} />
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
            required: t("cardLast4Required"),
            pattern: {
              value: /^\d{4}$/,
              message: t("cardLast4Invalid"),
            },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("cardLast4")}</FormLabel>
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
              required: t("billingDayRequired"),
              min: { value: 1, message: t("billingDayRange") },
              max: { value: 31, message: t("billingDayRange") },
            }}
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("billingDay")}</FormLabel>
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
                <FormDescription>{t("billingDayDesc")}</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          {/* 繳款截止日 */}
          <FormField
            control={form.control}
            name="payment_due_day"
            rules={{
              required: t("paymentDueDayRequired"),
              min: { value: 1, message: t("paymentDueDayRange") },
              max: { value: 31, message: t("paymentDueDayRange") },
            }}
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("paymentDueDay")}</FormLabel>
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
                <FormDescription>{t("paymentDueDayDesc")}</FormDescription>
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
              required: t("creditLimitRequired"),
              min: { value: 0.01, message: t("creditLimitPositive") },
            }}
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("creditLimit")}</FormLabel>
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

          {/* 已使用額度（可手動編輯） */}
          {creditCard && (
            <FormField
              control={form.control}
              name="used_credit"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("usedCredit")}</FormLabel>
                  <FormControl>
                    <Input type="number" step="0.01" {...field} />
                  </FormControl>
                  <FormDescription>{t("usedCreditDesc")}</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}
        </div>

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
              : creditCard
              ? tCommon("update")
              : tCommon("create")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
