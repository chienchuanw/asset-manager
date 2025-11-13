/**
 * 分期表單元件
 * 用於建立和編輯分期
 */

"use client";

import { useEffect, useState } from "react";
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
  const form = useForm<CreateInstallmentInput>({
    defaultValues: {
      name: "",
      total_amount: 0,
      currency: "TWD",
      installment_count: 12,
      interest_rate: 0,
      category_id: "",
      start_date: (() => {
        // 使用本地時間格式化日期，避免時區轉換問題
        const now = new Date();
        const year = now.getFullYear();
        const month = String(now.getMonth() + 1).padStart(2, "0");
        const day = String(now.getDate()).padStart(2, "0");
        return `${year}-${month}-${day}`;
      })(),
      billing_day: 1,
      note: "",
    },
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

  // 如果是編輯模式，填入現有資料
  useEffect(() => {
    if (installment) {
      form.reset({
        name: installment.name,
        total_amount: installment.total_amount,
        currency: installment.currency,
        installment_count: installment.installment_count,
        interest_rate: installment.interest_rate,
        category_id: installment.category_id,
        start_date: installment.start_date.split("T")[0],
        billing_day: installment.billing_day,
        note: installment.note || "",
      });
    }
  }, [installment, form]);

  // 篩選支出類別
  const expenseCategories = categories.filter((c) => c.type === "expense");

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        {/* 名稱 */}
        <FormField
          control={form.control}
          name="name"
          rules={{ required: "請輸入分期名稱" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>分期名稱</FormLabel>
              <FormControl>
                <Input placeholder="例如：iPhone 15 Pro" {...field} />
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
            required: "請輸入總金額",
            min: { value: 0.01, message: "總金額必須大於 0" },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>總金額（本金）</FormLabel>
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
            required: "請輸入分期期數",
            min: { value: 2, message: "分期期數至少為 2 期" },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>分期期數</FormLabel>
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
              <FormDescription>分幾期付款</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 利率 */}
        <FormField
          control={form.control}
          name="interest_rate"
          rules={{
            min: { value: 0, message: "利率不能為負數" },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>年利率 (%)</FormLabel>
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
              <FormDescription>如果是無息分期，請輸入 0</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 計算結果顯示 */}
        {calculatedValues.amountPerInstallment > 0 && (
          <div className="rounded-lg border p-4 space-y-2 bg-muted/50">
            <p className="text-sm font-medium">計算結果</p>
            <div className="grid grid-cols-2 gap-2 text-sm">
              <div>
                <p className="text-muted-foreground">每期金額</p>
                <p className="font-semibold tabular-nums">
                  TWD{" "}
                  {calculatedValues.amountPerInstallment.toLocaleString(
                    "zh-TW",
                    { maximumFractionDigits: 2 }
                  )}
                </p>
              </div>
              <div>
                <p className="text-muted-foreground">總利息</p>
                <p className="font-semibold tabular-nums">
                  TWD{" "}
                  {calculatedValues.totalInterest.toLocaleString("zh-TW", {
                    maximumFractionDigits: 2,
                  })}
                </p>
              </div>
              <div className="col-span-2">
                <p className="text-muted-foreground">總付款金額</p>
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
          rules={{ required: "請選擇分類" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>分類</FormLabel>
              <Select onValueChange={field.onChange} value={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="選擇分類" />
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

        {/* 開始日期 */}
        <FormField
          control={form.control}
          name="start_date"
          rules={{ required: "請選擇開始日期" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>開始日期</FormLabel>
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
            required: "請輸入扣款日",
            min: { value: 1, message: "扣款日必須在 1-31 之間" },
            max: { value: 31, message: "扣款日必須在 1-31 之間" },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>每月扣款日</FormLabel>
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
              <FormDescription>每月的第幾天扣款（1-31）</FormDescription>
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
              <FormLabel>備註（可選）</FormLabel>
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
        <div className="flex gap-2 justify-end">
          {onCancel && (
            <Button type="button" variant="outline" onClick={onCancel}>
              取消
            </Button>
          )}
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "儲存中..." : installment ? "更新" : "建立"}
          </Button>
        </div>
      </form>
    </Form>
  );
}
