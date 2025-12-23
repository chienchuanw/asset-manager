"use client";

import { useState } from "react";
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
import { useCreateTransaction } from "@/hooks";
import {
  createTransactionSchema,
  type CreateTransactionFormData,
  getAssetTypeOptions,
  getTransactionTypeOptions,
  AssetType,
  TransactionType,
  Currency,
} from "@/types/transaction";
import { Plus, Loader2 } from "lucide-react";
import { InstrumentCombobox } from "./InstrumentCombobox";

interface AddTransactionDialogProps {
  onSuccess?: () => void;
}

/**
 * 新增交易對話框
 *
 * 使用 react-hook-form + zod 進行表單驗證
 * 使用 useCreateTransaction hook 送出資料
 */
export function AddTransactionDialog({ onSuccess }: AddTransactionDialogProps) {
  const t = useTranslations("transactions");
  const tCommon = useTranslations("common");
  const [open, setOpen] = useState(false);

  // 建立交易 mutation
  const createMutation = useCreateTransaction({
    onSuccess: () => {
      setOpen(false);
      form.reset();
      onSuccess?.();
    },
  });

  // 表單設定
  const form = useForm<CreateTransactionFormData>({
    resolver: zodResolver(createTransactionSchema),
    defaultValues: {
      date: (() => {
        // 使用本地時間格式化日期，避免時區轉換問題
        const now = new Date();
        const year = now.getFullYear();
        const month = String(now.getMonth() + 1).padStart(2, "0");
        const day = String(now.getDate()).padStart(2, "0");
        return `${year}-${month}-${day}`;
      })(),
      asset_type: AssetType.TW_STOCK,
      symbol: "",
      name: "",
      type: TransactionType.BUY,
      quantity: 0,
      price: 0,
      amount: 0,
      fee: null,
      tax: null,
      currency: Currency.TWD,
      note: null,
    },
  });

  // 監聽數量和價格變化，自動計算金額
  const quantity = form.watch("quantity");
  const price = form.watch("price");

  // 當數量或價格改變時，自動更新金額
  const handleQuantityOrPriceChange = () => {
    const calculatedAmount = quantity * price;
    if (!isNaN(calculatedAmount)) {
      form.setValue("amount", calculatedAmount);
    }
  };

  // 送出表單
  const onSubmit = (data: CreateTransactionFormData) => {
    // 將日期字串轉換為 ISO 8601 格式，使用本地時間避免時區問題
    const [year, month, day] = data.date.split("-").map(Number);
    const localDate = new Date(year, month - 1, day, 12, 0, 0); // 設定為中午避免夏令時問題
    const isoDate = localDate.toISOString();

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
          {t("addTransaction")}
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{t("addTransactionTitle")}</DialogTitle>
          <DialogDescription>{t("addTransactionDesc")}</DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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

            {/* 資產類型 */}
            <FormField
              control={form.control}
              name="asset_type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("assetTypeLabel")}</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder={t("selectAssetType")} />
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
                    <FormLabel>{t("symbolLabel")}</FormLabel>
                    <FormControl>
                      <InstrumentCombobox
                        value={field.value}
                        assetType={form.watch("asset_type")}
                        onChange={(symbol) => {
                          // 手動輸入時更新代碼
                          form.setValue("symbol", symbol);
                        }}
                        onSelect={(instrument) => {
                          // 從清單選擇時，設定代碼並自動帶入名稱
                          form.setValue("symbol", instrument.symbol);
                          form.setValue("name", instrument.name);
                        }}
                        placeholder={t("symbolPlaceholder")}
                        searchPlaceholder={t("symbolSearchPlaceholder")}
                      />
                    </FormControl>
                    <FormDescription>{t("symbolDescription")}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("nameLabel")}</FormLabel>
                    <FormControl>
                      <Input placeholder={t("namePlaceholder")} {...field} />
                    </FormControl>
                    <FormDescription>{t("nameDescription")}</FormDescription>
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
                  <FormLabel>{t("typeLabel")}</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder={t("selectType")} />
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
                render={({ field }) => {
                  // 根據資產類型決定 step：台股/美股為整數，加密貨幣為小數
                  const assetType = form.watch("asset_type");
                  const step =
                    assetType === AssetType.TW_STOCK ||
                    assetType === AssetType.US_STOCK
                      ? "1"
                      : "0.01";

                  return (
                    <FormItem>
                      <FormLabel>{t("quantityLabel")}</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          step={step}
                          placeholder="0"
                          {...field}
                          onChange={(e) => {
                            field.onChange(parseFloat(e.target.value) || 0);
                            handleQuantityOrPriceChange();
                          }}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  );
                }}
              />

              <FormField
                control={form.control}
                name="price"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("priceLabel")}</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        step="0.01"
                        placeholder="0"
                        {...field}
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
                    <FormLabel>{t("amountLabel")}</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        step="0.01"
                        placeholder="0"
                        {...field}
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

            {/* 手續費、交易稅和幣別 */}
            <div className="grid grid-cols-3 gap-4">
              <FormField
                control={form.control}
                name="fee"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("feeOptional")}</FormLabel>
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
                name="tax"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("taxOptional")}</FormLabel>
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
                    <FormLabel>{t("currencyLabel")}</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder={t("selectCurrency")} />
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
                  <FormLabel>{t("noteLabel")}</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder={t("notePlaceholder")}
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
                onClick={() => setOpen(false)}
                disabled={createMutation.isPending}
              >
                {tCommon("cancel")}
              </Button>
              <Button type="submit" disabled={createMutation.isPending}>
                {createMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                {t("createTransaction")}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
