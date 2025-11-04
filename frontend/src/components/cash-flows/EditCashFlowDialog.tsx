"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { useUpdateCashFlow } from "@/hooks";
import {
  createCashFlowSchema,
  type CreateCashFlowFormData,
  type CashFlow,
  getCashFlowTypeOptions,
  CashFlowType,
  PaymentMethodType,
  paymentMethodTypeToSourceType,
} from "@/types/cash-flow";
import { Loader2 } from "lucide-react";
import { CategorySelect } from "./CategorySelect";
import { PaymentMethodSelect } from "./PaymentMethodSelect";
import { AccountSelect } from "./AccountSelect";
import { toast } from "sonner";

interface EditCashFlowDialogProps {
  cashFlow: CashFlow;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

/**
 * 編輯現金流對話框
 * 
 * 提供編輯現金流記錄的表單介面
 */
export function EditCashFlowDialog({
  cashFlow,
  open,
  onOpenChange,
  onSuccess,
}: EditCashFlowDialogProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);

  // 表單設定
  const form = useForm<CreateCashFlowFormData>({
    resolver: zodResolver(createCashFlowSchema),
    defaultValues: {
      date: "",
      type: CashFlowType.EXPENSE,
      category_id: "",
      description: "",
      amount: 0,
      payment_method_type: PaymentMethodType.CASH,
      payment_method_id: "",
      note: "",
    },
  });

  // 當 cashFlow 變更時，更新表單預設值
  useEffect(() => {
    if (cashFlow && open) {
      // 格式化日期為 YYYY-MM-DD 格式
      const formattedDate = new Date(cashFlow.date).toISOString().split('T')[0];
      
      form.reset({
        date: formattedDate,
        type: cashFlow.type,
        category_id: cashFlow.category?.id || "",
        description: cashFlow.description,
        amount: cashFlow.amount,
        payment_method_type: cashFlow.payment_method_type,
        payment_method_id: cashFlow.payment_method_id || "",
        note: cashFlow.note || "",
      });
    }
  }, [cashFlow, open, form]);

  // 更新現金流 mutation
  const updateMutation = useUpdateCashFlow({
    onSuccess: () => {
      toast.success("記錄更新成功");
      onOpenChange(false);
      onSuccess?.();
      form.reset();
    },
    onError: (error) => {
      toast.error(error.message || "更新失敗");
    },
    onSettled: () => {
      setIsSubmitting(false);
    },
  });

  // 表單提交處理
  const onSubmit = async (data: CreateCashFlowFormData) => {
    setIsSubmitting(true);

    // 準備更新資料
    const updateData = {
      ...data,
      // 轉換日期格式為 ISO 字串
      date: new Date(data.date + "T00:00:00").toISOString(),
      // 根據付款方式類型決定是否包含 payment_method_id
      payment_method_id: data.payment_method_type === PaymentMethodType.CASH 
        ? null 
        : data.payment_method_id,
    };

    updateMutation.mutate({
      id: cashFlow.id,
      data: updateData,
    });
  };

  // 監聽付款方式類型變更
  const paymentMethodType = form.watch("payment_method_type");

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>編輯現金流記錄</DialogTitle>
          <DialogDescription>
            修改現金流記錄的詳細資訊
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* 日期 */}
              <FormField
                control={form.control}
                name="date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>日期 *</FormLabel>
                    <FormControl>
                      <Input type="date" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* 類型 */}
              <FormField
                control={form.control}
                name="type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>類型 *</FormLabel>
                    <Select onValueChange={field.onChange} value={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="選擇類型" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {getCashFlowTypeOptions().map((option) => (
                          <SelectItem key={option.value} value={option.value}>
                            {option.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* 分類 */}
              <FormField
                control={form.control}
                name="category_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>分類</FormLabel>
                    <CategorySelect
                      value={field.value}
                      onValueChange={field.onChange}
                      type={form.watch("type")}
                    />
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* 金額 */}
              <FormField
                control={form.control}
                name="amount"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>金額 *</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        step="0.01"
                        min="0"
                        placeholder="0.00"
                        {...field}
                        onChange={(e) => field.onChange(Number(e.target.value))}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* 描述 */}
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>描述 *</FormLabel>
                  <FormControl>
                    <Input placeholder="請輸入描述" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* 付款方式 */}
            <FormField
              control={form.control}
              name="payment_method_type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>付款方式 *</FormLabel>
                  <PaymentMethodSelect
                    value={field.value}
                    onValueChange={field.onChange}
                  />
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* 帳戶選擇（當付款方式不是現金時顯示） */}
            {paymentMethodType !== PaymentMethodType.CASH && (
              <FormField
                control={form.control}
                name="payment_method_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      {paymentMethodType === PaymentMethodType.BANK_ACCOUNT
                        ? "銀行帳戶"
                        : "信用卡"}{" "}
                      *
                    </FormLabel>
                    <AccountSelect
                      value={field.value}
                      onValueChange={field.onChange}
                      sourceType={paymentMethodTypeToSourceType(paymentMethodType)}
                    />
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            {/* 備註 */}
            <FormField
              control={form.control}
              name="note"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>備註</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="選填備註資訊"
                      className="resize-none"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isSubmitting}
              >
                取消
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                更新記錄
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
