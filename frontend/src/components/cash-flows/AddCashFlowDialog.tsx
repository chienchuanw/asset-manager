"use client";

import React, { useState } from "react";
import { useTranslations } from "next-intl";
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
  PaymentMethodType,
  paymentMethodTypeToSourceType,
} from "@/types/cash-flow";
import { Plus, Loader2 } from "lucide-react";
import { CategorySelect } from "./CategorySelect";
import { PaymentMethodSelect } from "./PaymentMethodSelect";
import { AccountSelect } from "./AccountSelect";
import { toast } from "sonner";
import { useCategories } from "@/hooks";

/**
 * 將日期字串轉換為 ISO 格式，避免時區問題
 * @param dateString YYYY-MM-DD 格式的日期字串
 * @returns ISO 8601 格式的日期字串
 */
const formatDateToISO = (dateString: string): string => {
  // 解析日期字串為年、月、日
  const [year, month, day] = dateString.split("-").map(Number);
  // 使用本地時間建立 Date 物件，避免時區轉換
  const date = new Date(year, month - 1, day, 12, 0, 0); // 設定為中午避免夏令時問題
  return date.toISOString();
};

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
  const t = useTranslations("cashFlows");
  const tCommon = useTranslations("common");
  const [open, setOpen] = useState(false);

  // 建立現金流 mutation
  const createMutation = useCreateCashFlow({
    onSuccess: () => {
      toast.success(t("successMessage"));
      setOpen(false);
      form.reset();
      onSuccess?.();
    },
    onError: (error) => {
      toast.error(error.message || t("errorMessage"));
    },
  });

  // 表單設定
  const form = useForm<CreateCashFlowFormData>({
    resolver: zodResolver(createCashFlowSchema),
    defaultValues: {
      date: (() => {
        // 使用本地時間格式化日期，避免時區轉換問題
        const now = new Date();
        const year = now.getFullYear();
        const month = String(now.getMonth() + 1).padStart(2, "0");
        const day = String(now.getDate()).padStart(2, "0");
        return `${year}-${month}-${day}`;
      })(),
      type: CashFlowType.EXPENSE,
      category_id: "",
      amount: undefined, // 不設定預設值
      description: "",
      note: null,
      payment_method: PaymentMethodType.CASH, // 預設為現金
      account_id: "", // 帳戶 ID
      target_payment_method: undefined, // 轉帳目標付款方式
      target_account_id: "", // 轉帳目標帳戶 ID
    },
  });

  // 監聽類型變化，重置分類選擇
  const cashFlowType = form.watch("type");
  // 監聽付款方式變化，重置帳戶選擇
  const paymentMethod = form.watch("payment_method");
  // 監聽轉帳目標付款方式變化
  const targetPaymentMethod = form.watch("target_payment_method");
  // 監聽分類變化
  const categoryId = form.watch("category_id");

  // 取得分類列表（用於判斷是否為「提領」分類）
  const { data: categories } = useCategories(cashFlowType);

  // 當類型變為轉帳時，自動設定付款方式為銀行帳戶（分類交由子元件自動帶入）
  React.useEffect(() => {
    const isTransfer =
      cashFlowType === CashFlowType.TRANSFER_IN ||
      cashFlowType === CashFlowType.TRANSFER_OUT;

    if (isTransfer) {
      // 設定付款方式為銀行帳戶
      form.setValue("payment_method", PaymentMethodType.BANK_ACCOUNT);
      form.setValue("account_id", "");
    }
  }, [cashFlowType, form]);

  // 判斷當前選中的分類是否為「提領」
  const selectedCategory = React.useMemo(() => {
    if (!categories || !categoryId) return null;
    return categories.find((c) => c.id === categoryId);
  }, [categories, categoryId]);

  const isCashWithdrawal =
    cashFlowType === CashFlowType.TRANSFER_OUT &&
    selectedCategory?.name === "提領";

  // 當選擇「提領」分類時，自動設定 target_payment_method 為 cash
  React.useEffect(() => {
    if (isCashWithdrawal) {
      // 自動設定目標付款方式為現金
      form.setValue("target_payment_method", PaymentMethodType.CASH);
      // 清空目標帳戶 ID（現金不需要帳戶）
      form.setValue("target_account_id", "");
    }
  }, [isCashWithdrawal, form]);

  // 送出表單
  const onSubmit = (data: CreateCashFlowFormData) => {
    // 將日期轉換為 ISO 8601 格式，避免時區問題
    const isoDate = formatDateToISO(data.date);

    // 準備提交資料
    const submitData: any = {
      date: isoDate,
      type: data.type,
      category_id: data.category_id,
      amount: data.amount,
      description: data.description,
      note: data.note,
    };

    // 根據付款方式設定 source_type 和 source_id
    if (data.payment_method !== PaymentMethodType.CASH) {
      submitData.source_type = paymentMethodTypeToSourceType(
        data.payment_method
      );
      submitData.source_id = data.account_id;
    }

    // 如果是 transfer_out 類型，設定 target_type 和 target_id
    if (data.type === CashFlowType.TRANSFER_OUT && data.target_payment_method) {
      submitData.target_type = paymentMethodTypeToSourceType(
        data.target_payment_method
      );
      // 如果目標是現金，target_id 為 null；否則使用 target_account_id
      if (data.target_payment_method === PaymentMethodType.CASH) {
        submitData.target_id = null;
      } else if (data.target_account_id) {
        submitData.target_id = data.target_account_id;
      }
    }

    createMutation.mutate(submitData);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <Plus className="h-4 w-4 mr-2" />
          {t("addRecord")}
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{t("addCashFlowTitle")}</DialogTitle>
          <DialogDescription>{t("addCashFlowDesc")}</DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {/* 金額顯示區域（計算機風格） */}
            <FormField
              control={form.control}
              name="amount"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-sm text-muted-foreground">
                    {t("amountCurrency")}
                  </FormLabel>
                  <FormControl>
                    <div className="relative">
                      {/* 計算機風格的金額顯示 - 可輸入 */}
                      <input
                        type="number"
                        inputMode="numeric"
                        step="1"
                        min="0"
                        placeholder="0"
                        className="w-full bg-slate-100 dark:bg-slate-900 rounded-lg p-6 min-h-20 text-4xl font-bold tabular-nums text-right border-0 focus:outline-none focus:ring-2 focus:ring-ring"
                        {...field}
                        value={field.value || ""}
                        onChange={(e) => {
                          const value = e.target.value;
                          field.onChange(value ? parseFloat(value) : undefined);
                        }}
                      />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* 類型與分類並排 */}
            <div className="grid grid-cols-2 gap-4">
              {/* 類型 */}
              <FormField
                control={form.control}
                name="type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("type")}</FormLabel>
                    <Select
                      onValueChange={(value) => {
                        field.onChange(value);
                        // 重置分類選擇（轉帳類型會由 useEffect 自動設定，所以這裡不重置）
                        const isTransfer =
                          value === CashFlowType.TRANSFER_IN ||
                          value === CashFlowType.TRANSFER_OUT;
                        if (!isTransfer) {
                          form.setValue("category_id", "");
                        }
                      }}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder={t("selectType")} />
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
                    <FormLabel>{t("category")}</FormLabel>
                    <FormControl>
                      <CategorySelect
                        value={field.value}
                        onValueChange={field.onChange}
                        type={cashFlowType}
                        placeholder={t("selectCategory")}
                        autoSelectName={
                          cashFlowType === CashFlowType.TRANSFER_IN ||
                          cashFlowType === CashFlowType.TRANSFER_OUT
                            ? t("transfer")
                            : undefined
                        }
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* 日期 */}
            <FormField
              control={form.control}
              name="date"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("date")}</FormLabel>
                  <FormControl>
                    <Input type="date" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* 付款方式 */}
            <FormField
              control={form.control}
              name="payment_method"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("paymentMethod")}</FormLabel>
                  <FormControl>
                    <PaymentMethodSelect
                      value={field.value}
                      onValueChange={(value) => {
                        field.onChange(value);
                        // 重置帳戶選擇
                        form.setValue("account_id", "");
                      }}
                      placeholder={t("selectPaymentMethod")}
                      disabled={
                        cashFlowType === CashFlowType.TRANSFER_IN ||
                        cashFlowType === CashFlowType.TRANSFER_OUT
                      }
                    />
                  </FormControl>
                  <FormMessage />
                  {(cashFlowType === CashFlowType.TRANSFER_IN ||
                    cashFlowType === CashFlowType.TRANSFER_OUT) && (
                    <p className="text-sm text-muted-foreground">
                      {t("transferOnlyBankAccount")}
                    </p>
                  )}
                </FormItem>
              )}
            />

            {/* 銀行帳戶選擇 */}
            {paymentMethod !== PaymentMethodType.CASH && (
              <FormField
                control={form.control}
                name="account_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      {paymentMethod === PaymentMethodType.BANK_ACCOUNT
                        ? t("bankAccount")
                        : t("creditCard")}
                    </FormLabel>
                    <FormControl>
                      <AccountSelect
                        value={field.value || ""}
                        onValueChange={field.onChange}
                        paymentMethodType={paymentMethod}
                        placeholder={
                          paymentMethod === PaymentMethodType.BANK_ACCOUNT
                            ? t("selectBankAccount")
                            : t("selectCreditCard")
                        }
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            {/* 轉帳目標選擇（僅在 transfer_out 且非提領時顯示） */}
            {cashFlowType === CashFlowType.TRANSFER_OUT &&
              !isCashWithdrawal && (
                <>
                  {/* 轉帳目標付款方式 */}
                  <FormField
                    control={form.control}
                    name="target_payment_method"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t("transferTargetType")}</FormLabel>
                        <FormControl>
                          <PaymentMethodSelect
                            value={field.value}
                            onValueChange={(value) => {
                              field.onChange(value);
                              // 重置目標帳戶選擇
                              form.setValue("target_account_id", "");
                            }}
                            placeholder={t("selectTransferTargetType")}
                            excludeCash={true}
                          />
                        </FormControl>
                        <FormMessage />
                        <p className="text-sm text-muted-foreground">
                          {t("selectTransferTarget")}
                        </p>
                      </FormItem>
                    )}
                  />

                  {/* 轉帳目標帳戶選擇 */}
                  {targetPaymentMethod && (
                    <FormField
                      control={form.control}
                      name="target_account_id"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            {targetPaymentMethod ===
                            PaymentMethodType.BANK_ACCOUNT
                              ? t("targetBankAccount")
                              : t("targetCreditCard")}
                          </FormLabel>
                          <FormControl>
                            <AccountSelect
                              value={field.value || ""}
                              onValueChange={field.onChange}
                              paymentMethodType={targetPaymentMethod}
                              placeholder={
                                targetPaymentMethod ===
                                PaymentMethodType.BANK_ACCOUNT
                                  ? t("selectBankAccount")
                                  : t("selectCreditCard")
                              }
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  )}
                </>
              )}

            {/* 描述 */}
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("description")}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t("descriptionPlaceholder")}
                      {...field}
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
                      placeholder={t("notePlaceholder")}
                      className="resize-none"
                      {...field}
                      value={field.value || ""}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter className="flex flex-row gap-2 sm:gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
                disabled={createMutation.isPending}
                className="flex-1"
              >
                {tCommon("cancel")}
              </Button>
              <Button
                type="submit"
                disabled={createMutation.isPending}
                className="flex-1"
              >
                {createMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                {tCommon("create")}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
