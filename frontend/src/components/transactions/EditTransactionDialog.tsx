"use client";

import { useEffect } from "react";
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
  FormDescription,
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
import { useUpdateTransaction } from "@/hooks";
import {
  updateTransactionSchema,
  type UpdateTransactionFormData,
  type Transaction,
  getAssetTypeOptions,
  getTransactionTypeOptions,
  Currency,
} from "@/types/transaction";
import { Loader2 } from "lucide-react";

interface EditTransactionDialogProps {
  transaction: Transaction;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

/**
 * 編輯交易對話框
 *
 * 使用 react-hook-form + zod 進行表單驗證
 * 使用 useUpdateTransaction hook 更新資料
 */
export function EditTransactionDialog({
  transaction,
  open,
  onOpenChange,
  onSuccess,
}: EditTransactionDialogProps) {
  // 更新交易 mutation
  const updateMutation = useUpdateTransaction({
    onSuccess: () => {
      onOpenChange(false);
      onSuccess?.();
    },
  });

  // 表單設定
  const form = useForm<UpdateTransactionFormData>({
    resolver: zodResolver(updateTransactionSchema),
    defaultValues: {
      date: transaction.date.split("T")[0], // 轉換為 YYYY-MM-DD 格式
      asset_type: transaction.asset_type,
      symbol: transaction.symbol,
      name: transaction.name,
      type: transaction.type,
      quantity: transaction.quantity,
      price: transaction.price,
      amount: transaction.amount,
      fee: transaction.fee,
      currency: transaction.currency,
      note: transaction.note,
    },
  });

  // 當 transaction 改變時，重新設定表單預設值
  useEffect(() => {
    form.reset({
      date: transaction.date.split("T")[0],
      asset_type: transaction.asset_type,
      symbol: transaction.symbol,
      name: transaction.name,
      type: transaction.type,
      quantity: transaction.quantity,
      price: transaction.price,
      amount: transaction.amount,
      fee: transaction.fee,
      currency: transaction.currency,
      note: transaction.note,
    });
  }, [transaction, form]);

  // 監聽數量和價格變化，自動計算金額
  const quantity = form.watch("quantity");
  const price = form.watch("price");

  // 當數量或價格改變時，自動更新金額
  const handleQuantityOrPriceChange = () => {
    if (quantity !== undefined && price !== undefined) {
      const calculatedAmount = quantity * price;
      if (!isNaN(calculatedAmount)) {
        form.setValue("amount", calculatedAmount);
      }
    }
  };

  // 送出表單
  const onSubmit = (data: UpdateTransactionFormData) => {
    // 將日期轉換為 ISO 8601 格式（如果有提供日期）
    const submitData = {
      ...data,
      date: data.date ? new Date(data.date).toISOString() : undefined,
    };

    updateMutation.mutate({
      id: transaction.id,
      data: submitData,
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>編輯交易記錄</DialogTitle>
          <DialogDescription>修改交易資訊並儲存變更</DialogDescription>
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
                    <Input type="date" {...field} value={field.value ?? ""} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* 資產類型 */}
            <FormField
              control={form.control}
              name="asset_type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>資產類型</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇資產類型" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {getAssetTypeOptions().map((option) => (
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

            {/* 代碼和名稱 */}
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="symbol"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>代碼</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="例如: 2330"
                        {...field}
                        value={field.value ?? ""}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>名稱</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="例如: 台積電"
                        {...field}
                        value={field.value ?? ""}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* 交易類型 */}
            <FormField
              control={form.control}
              name="type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>交易類型</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇交易類型" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {getTransactionTypeOptions().map((option) => (
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

            {/* 數量、價格、金額 */}
            <div className="grid grid-cols-3 gap-4">
              <FormField
                control={form.control}
                name="quantity"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>數量</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        step="0.01"
                        placeholder="0"
                        {...field}
                        value={field.value ?? ""}
                        onChange={(e) => {
                          field.onChange(parseFloat(e.target.value) || 0);
                          handleQuantityOrPriceChange();
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="price"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>價格</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        step="0.01"
                        placeholder="0"
                        {...field}
                        value={field.value ?? ""}
                        onChange={(e) => {
                          field.onChange(parseFloat(e.target.value) || 0);
                          handleQuantityOrPriceChange();
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="amount"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>金額</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        step="0.01"
                        placeholder="0"
                        {...field}
                        value={field.value ?? ""}
                        onChange={(e) =>
                          field.onChange(parseFloat(e.target.value) || 0)
                        }
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      自動計算或手動輸入
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* 手續費和幣別 */}
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="fee"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>手續費（選填）</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        step="0.01"
                        placeholder="0"
                        {...field}
                        value={field.value ?? ""}
                        onChange={(e) =>
                          field.onChange(
                            e.target.value ? parseFloat(e.target.value) : null
                          )
                        }
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="currency"
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
                        <SelectItem value={Currency.TWD}>
                          新台幣 (TWD)
                        </SelectItem>
                        <SelectItem value={Currency.USD}>美金 (USD)</SelectItem>
                      </SelectContent>
                    </Select>
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
                      {...field}
                      value={field.value ?? ""}
                      onChange={(e) => field.onChange(e.target.value || null)}
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
                disabled={updateMutation.isPending}
              >
                取消
              </Button>
              <Button type="submit" disabled={updateMutation.isPending}>
                {updateMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                儲存變更
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
