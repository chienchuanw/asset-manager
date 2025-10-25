"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
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
import { useCreateCashFlow } from "@/hooks";
import {
  createCashFlowSchema,
  type CreateCashFlowFormData,
  getCashFlowTypeOptions,
  CashFlowType,
} from "@/types/cash-flow";
import { Plus, Loader2 } from "lucide-react";
import { CategorySelect } from "./CategorySelect";
import { toast } from "sonner";

interface AddCashFlowDialogProps {
  onSuccess?: () => void;
}

/**
 * 新增現金流對話框
 *
 * 使用 react-hook-form + zod 進行表單驗證
 * 使用 useCreateCashFlow hook 送出資料
 */
export function AddCashFlowDialog({ onSuccess }: AddCashFlowDialogProps) {
  const [open, setOpen] = useState(false);

  // 建立現金流 mutation
  const createMutation = useCreateCashFlow({
    onSuccess: () => {
      toast.success("記錄建立成功");
      setOpen(false);
      form.reset();
      onSuccess?.();
    },
    onError: (error) => {
      toast.error(error.message || "建立失敗");
    },
  });

  // 表單設定
  const form = useForm<CreateCashFlowFormData>({
    resolver: zodResolver(createCashFlowSchema),
    defaultValues: {
      date: new Date().toISOString().split("T")[0], // YYYY-MM-DD 格式
      type: CashFlowType.EXPENSE,
      category_id: "",
      amount: 0,
      description: "",
      note: null,
    },
  });

  // 監聽類型變化，重置分類選擇
  const cashFlowType = form.watch("type");

  // 送出表單
  const onSubmit = (data: CreateCashFlowFormData) => {
    // 將日期轉換為 ISO 8601 格式
    const isoDate = new Date(data.date).toISOString();

    createMutation.mutate({
      ...data,
      date: isoDate,
    });
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <Plus className="h-4 w-4 mr-2" />
          新增記錄
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>新增現金流記錄</DialogTitle>
          <DialogDescription>
            記錄您的收入或支出，以便追蹤現金流動
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {/* 日期 */}
            <FormField
              control={form.control}
              name="date"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>日期</FormLabel>
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
                  <FormLabel>類型</FormLabel>
                  <Select
                    onValueChange={(value) => {
                      field.onChange(value);
                      // 重置分類選擇
                      form.setValue("category_id", "");
                    }}
                    defaultValue={field.value}
                  >
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

            {/* 分類 */}
            <FormField
              control={form.control}
              name="category_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>分類</FormLabel>
                  <FormControl>
                    <CategorySelect
                      value={field.value}
                      onValueChange={field.onChange}
                      type={cashFlowType}
                      placeholder="選擇分類"
                    />
                  </FormControl>
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
                  <FormLabel>金額 (TWD)</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      step="0.01"
                      placeholder="0"
                      {...field}
                      onChange={(e) => {
                        field.onChange(parseFloat(e.target.value) || 0);
                      }}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* 描述 */}
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>描述</FormLabel>
                  <FormControl>
                    <Input placeholder="例如: 十月薪資" {...field} />
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
                  <FormLabel>備註（可選）</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="額外說明..."
                      className="resize-none"
                      {...field}
                      value={field.value || ""}
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
                onClick={() => setOpen(false)}
                disabled={createMutation.isPending}
              >
                取消
              </Button>
              <Button type="submit" disabled={createMutation.isPending}>
                {createMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                建立
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

