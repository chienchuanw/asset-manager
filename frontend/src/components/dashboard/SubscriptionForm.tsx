/**
 * 訂閱表單元件
 * 用於建立和編輯訂閱
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
  BillingCycle,
} from "@/types/subscription";
import type { CashFlowCategory } from "@/types/cash-flow";

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
  const form = useForm<CreateSubscriptionInput>({
    defaultValues: {
      name: "",
      amount: 0,
      currency: "TWD",
      billing_cycle: "monthly",
      billing_day: 1,
      category_id: "",
      start_date: new Date().toISOString().split("T")[0],
      auto_renew: true,
      note: "",
    },
  });

  // 如果是編輯模式，填入現有資料
  useEffect(() => {
    if (subscription) {
      form.reset({
        name: subscription.name,
        amount: subscription.amount,
        currency: subscription.currency,
        billing_cycle: subscription.billing_cycle,
        billing_day: subscription.billing_day,
        category_id: subscription.category_id,
        start_date: subscription.start_date.split("T")[0],
        end_date: subscription.end_date?.split("T")[0],
        auto_renew: subscription.auto_renew,
        note: subscription.note || "",
      });
    }
  }, [subscription, form]);

  // 篩選支出類別
  const expenseCategories = categories.filter((c) => c.type === "expense");

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        {/* 名稱 */}
        <FormField
          control={form.control}
          name="name"
          rules={{ required: "請輸入訂閱名稱" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>訂閱名稱</FormLabel>
              <FormControl>
                <Input placeholder="例如：Netflix" {...field} />
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
            required: "請輸入金額",
            min: { value: 0.01, message: "金額必須大於 0" },
          }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>金額</FormLabel>
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
              <FormLabel>計費週期</FormLabel>
              <Select onValueChange={field.onChange} value={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="選擇計費週期" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="monthly">每月</SelectItem>
                  <SelectItem value="quarterly">每季</SelectItem>
                  <SelectItem value="yearly">每年</SelectItem>
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

        {/* 結束日期（可選） */}
        <FormField
          control={form.control}
          name="end_date"
          render={({ field }) => (
            <FormItem>
              <FormLabel>結束日期（可選）</FormLabel>
              <FormControl>
                <Input type="date" {...field} />
              </FormControl>
              <FormDescription>
                如果不設定，訂閱將持續到手動取消
              </FormDescription>
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
                <FormLabel className="text-base">自動續約</FormLabel>
                <FormDescription>到期後自動續約訂閱</FormDescription>
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
            {isSubmitting ? "儲存中..." : subscription ? "更新" : "建立"}
          </Button>
        </div>
      </form>
    </Form>
  );
}
