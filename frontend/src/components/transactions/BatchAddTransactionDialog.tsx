"use client";

import { useState, useEffect } from "react";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Copy, Trash2, Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { InstrumentCombobox } from "./InstrumentCombobox";

import { useCreateTransactionsBatch } from "@/hooks/useTransactions";
import {
  AssetType,
  TransactionType,
  Currency,
  batchCreateTransactionsSchema,
  type BatchCreateTransactionsFormData,
  getAssetTypeLabel,
  getTransactionTypeLabel,
} from "@/types/transaction";

interface BatchAddTransactionDialogProps {
  onSuccess?: () => void;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  initialTransactions?: any[];
}

/**
 * 批次新增交易對話框
 *
 * 使用表格式介面讓使用者一次新增多筆交易
 * 支援動態新增/刪除/複製列
 */
export function BatchAddTransactionDialog({
  onSuccess,
  open: controlledOpen,
  onOpenChange,
  initialTransactions,
}: BatchAddTransactionDialogProps) {
  const [internalOpen, setInternalOpen] = useState(false);

  // 如果有傳入 open prop，使用受控模式；否則使用內部狀態
  const open = controlledOpen !== undefined ? controlledOpen : internalOpen;
  const setOpen = onOpenChange || setInternalOpen;

  // 建立批次交易 mutation
  const createBatchMutation = useCreateTransactionsBatch({
    onSuccess: () => {
      setOpen(false);
      form.reset();
      onSuccess?.();
    },
  });

  // 表單設定
  const form = useForm<BatchCreateTransactionsFormData>({
    resolver: zodResolver(batchCreateTransactionsSchema),
    defaultValues: {
      transactions: Array(5)
        .fill(null)
        .map(() => ({
          date: new Date().toISOString().split("T")[0],
          asset_type: AssetType.TW_STOCK,
          currency: Currency.TWD,
          symbol: "",
          name: "",
          type: TransactionType.BUY,
          quantity: 0,
          price: 0,
          amount: 0,
          fee: null,
          tax: null,
          note: null,
        })),
    },
  });

  // 當有 initialTransactions 時，更新表單資料
  useEffect(() => {
    if (initialTransactions && initialTransactions.length > 0) {
      // 設定所有交易（每筆交易都有自己的 asset_type 和 currency）
      form.setValue(
        "transactions",
        initialTransactions.map((t: any) => ({
          date: t.date.split("T")[0], // 確保日期格式正確
          asset_type: t.asset_type,
          currency: t.currency,
          symbol: t.symbol,
          name: t.name,
          type: t.type, // 修正：後端回傳的欄位名稱是 "type" 而非 "transaction_type"
          quantity: t.quantity,
          price: t.price,
          amount: t.amount,
          fee: t.fee || null,
          tax: t.tax || null,
          note: t.note || null,
        }))
      );
    }
  }, [initialTransactions, form]);

  // 使用 useFieldArray 管理動態列
  const { fields, append, remove, insert } = useFieldArray({
    control: form.control,
    name: "transactions",
  });

  // 新增一列
  const handleAddRow = () => {
    append({
      date: new Date().toISOString().split("T")[0],
      asset_type: AssetType.TW_STOCK,
      currency: Currency.TWD,
      symbol: "",
      name: "",
      type: TransactionType.BUY,
      quantity: 0,
      price: 0,
      amount: 0,
      fee: null,
      tax: null,
      note: null,
    });
  };

  // 複製一列
  const handleCopyRow = (index: number) => {
    const rowData = form.getValues(`transactions.${index}`);
    insert(index + 1, { ...rowData });
  };

  // 刪除一列（至少保留 1 列）
  const handleDeleteRow = (index: number) => {
    if (fields.length > 1) {
      remove(index);
    }
  };

  // 送出表單
  const onSubmit = (data: BatchCreateTransactionsFormData) => {
    // 轉換資料格式（每筆交易都有自己的 asset_type 和 currency）
    const transactions = data.transactions.map((tx) => ({
      date: new Date(tx.date).toISOString(),
      asset_type: tx.asset_type,
      symbol: tx.symbol,
      name: tx.name,
      type: tx.type,
      quantity: tx.quantity,
      price: tx.price,
      amount: tx.amount,
      fee: tx.fee,
      tax: tx.tax,
      currency: tx.currency,
      note: tx.note,
    }));

    createBatchMutation.mutate({ transactions });
  };

  // 收集錯誤摘要
  const errors = form.formState.errors;
  const transactionErrors = errors.transactions as any[] | undefined;
  const errorSummary = (transactionErrors || [])
    .map((err: any, index: number) => {
      if (!err) return null;
      const errorMessages = Object.entries(err)
        .filter(
          ([_, value]) =>
            value && typeof value === "object" && "message" in value
        )
        .map(([key, value]) => `${key}: ${(value as any).message}`);
      return errorMessages.length > 0
        ? { row: index + 1, errors: errorMessages }
        : null;
    })
    .filter((item: any) => item !== null);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm" variant="outline">
          <Plus className="h-4 w-4 mr-2" />
          批次新增交易
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-[95vw] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>批次新增交易記錄</DialogTitle>
          <DialogDescription>
            一次新增多筆交易記錄，填寫完成後點擊「建立交易」
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {/* 錯誤摘要 */}
            {errorSummary.length > 0 && (
              <Alert variant="destructive">
                <AlertDescription>
                  <div className="font-semibold mb-2">請修正以下錯誤：</div>
                  <ul className="list-disc list-inside space-y-1">
                    {errorSummary.map((item, idx) => (
                      <li key={idx}>
                        第 {item!.row} 列：{item!.errors.join(", ")}
                      </li>
                    ))}
                  </ul>
                </AlertDescription>
              </Alert>
            )}

            {/* 交易列表表格 */}
            <div className="border rounded-lg overflow-hidden">
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-12">#</TableHead>
                      <TableHead className="min-w-[140px]">日期</TableHead>
                      <TableHead className="min-w-[120px]">資產類型</TableHead>
                      <TableHead className="min-w-20">幣別</TableHead>
                      <TableHead className="min-w-[100px]">代碼</TableHead>
                      <TableHead className="min-w-[120px]">名稱</TableHead>
                      <TableHead className="min-w-[100px]">類型</TableHead>
                      <TableHead className="min-w-[100px]">數量</TableHead>
                      <TableHead className="min-w-[100px]">價格</TableHead>
                      <TableHead className="min-w-[120px]">金額</TableHead>
                      <TableHead className="min-w-[100px]">手續費</TableHead>
                      <TableHead className="min-w-[100px]">交易稅</TableHead>
                      <TableHead className="min-w-[150px]">備註</TableHead>
                      <TableHead className="w-24">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {fields.map((field, index) => (
                      <TransactionRow
                        key={field.id}
                        index={index}
                        form={form}
                        onCopy={() => handleCopyRow(index)}
                        onDelete={() => handleDeleteRow(index)}
                        canDelete={fields.length > 1}
                      />
                    ))}
                  </TableBody>
                </Table>
              </div>
            </div>

            {/* 新增列按鈕 */}
            <div className="flex justify-start">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={handleAddRow}
              >
                <Plus className="h-4 w-4 mr-2" />
                新增一列
              </Button>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
                disabled={createBatchMutation.isPending}
              >
                取消
              </Button>
              <Button type="submit" disabled={createBatchMutation.isPending}>
                {createBatchMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                建立交易
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

/**
 * 單列交易表單元件
 */
interface TransactionRowProps {
  index: number;
  form: any;
  onCopy: () => void;
  onDelete: () => void;
  canDelete: boolean;
}

function TransactionRow({
  index,
  form,
  onCopy,
  onDelete,
  canDelete,
}: TransactionRowProps) {
  // 監聽數量和價格變化，自動計算金額
  const quantity = form.watch(`transactions.${index}.quantity`);
  const price = form.watch(`transactions.${index}.price`);
  const assetType = form.watch(`transactions.${index}.asset_type`);

  useEffect(() => {
    const amount = quantity * price;
    if (!isNaN(amount) && amount >= 0) {
      form.setValue(`transactions.${index}.amount`, amount);
    }
  }, [quantity, price, index, form]);

  // 監聽資產類型變化，自動更新幣別
  useEffect(() => {
    if (assetType === AssetType.TW_STOCK) {
      form.setValue(`transactions.${index}.currency`, Currency.TWD);
    } else if (assetType === AssetType.US_STOCK) {
      form.setValue(`transactions.${index}.currency`, Currency.USD);
    }
    // 加密貨幣可以是 TWD 或 USD，不自動變更
  }, [assetType, index, form]);

  return (
    <>
      <TableRow>
        {/* 序號 */}
        <TableCell className="font-medium">{index + 1}</TableCell>

        {/* 日期 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.date`}
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input type="date" {...field} className="h-9" />
                </FormControl>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 資產類型 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.asset_type`}
            render={({ field }) => (
              <FormItem className="space-y-0">
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger className="h-9">
                      <SelectValue />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {Object.values(AssetType)
                      .filter((type) => type !== AssetType.CASH)
                      .map((type) => (
                        <SelectItem key={type} value={type}>
                          {getAssetTypeLabel(type)}
                        </SelectItem>
                      ))}
                  </SelectContent>
                </Select>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 幣別 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.currency`}
            render={({ field }) => (
              <FormItem className="space-y-0">
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger className="h-9">
                      <SelectValue />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value={Currency.TWD}>TWD</SelectItem>
                    <SelectItem value={Currency.USD}>USD</SelectItem>
                  </SelectContent>
                </Select>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 代碼 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.symbol`}
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <InstrumentCombobox
                    value={field.value}
                    assetType={assetType}
                    onChange={field.onChange}
                    onSelect={(instrument) => {
                      form.setValue(
                        `transactions.${index}.symbol`,
                        instrument.symbol
                      );
                      form.setValue(
                        `transactions.${index}.name`,
                        instrument.name
                      );
                    }}
                    placeholder="2330"
                  />
                </FormControl>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 名稱 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.name`}
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input placeholder="台積電" {...field} className="h-9" />
                </FormControl>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 交易類型 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.type`}
            render={({ field }) => (
              <FormItem className="space-y-0">
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger className="h-9">
                      <SelectValue />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {Object.values(TransactionType).map((type) => (
                      <SelectItem key={type} value={type}>
                        {getTransactionTypeLabel(type)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 數量 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.quantity`}
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input
                    type="number"
                    step="any"
                    placeholder="0"
                    {...field}
                    onChange={(e) =>
                      field.onChange(parseFloat(e.target.value) || 0)
                    }
                    className="h-9"
                  />
                </FormControl>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 價格 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.price`}
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input
                    type="number"
                    step="any"
                    placeholder="0"
                    {...field}
                    onChange={(e) =>
                      field.onChange(parseFloat(e.target.value) || 0)
                    }
                    className="h-9"
                  />
                </FormControl>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 金額（自動計算） */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.amount`}
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input
                    type="number"
                    step="any"
                    placeholder="0"
                    {...field}
                    onChange={(e) =>
                      field.onChange(parseFloat(e.target.value) || 0)
                    }
                    className="h-9 bg-muted"
                  />
                </FormControl>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 手續費 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.fee`}
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input
                    type="number"
                    step="any"
                    placeholder="0"
                    {...field}
                    value={field.value ?? ""}
                    onChange={(e) =>
                      field.onChange(
                        e.target.value === ""
                          ? null
                          : parseFloat(e.target.value)
                      )
                    }
                    className="h-9"
                  />
                </FormControl>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 交易稅 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.tax`}
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input
                    type="number"
                    step="any"
                    placeholder="0"
                    {...field}
                    value={field.value ?? ""}
                    onChange={(e) =>
                      field.onChange(
                        e.target.value === ""
                          ? null
                          : parseFloat(e.target.value)
                      )
                    }
                    className="h-9"
                  />
                </FormControl>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 備註 */}
        <TableCell>
          <FormField
            control={form.control}
            name={`transactions.${index}.note`}
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input
                    placeholder="輸入備註..."
                    {...field}
                    value={field.value ?? ""}
                    onChange={(e) => field.onChange(e.target.value || null)}
                    className="h-9"
                  />
                </FormControl>
              </FormItem>
            )}
          />
        </TableCell>

        {/* 操作按鈕 */}
        <TableCell>
          <div className="flex gap-1">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={onCopy}
              title="複製此列"
            >
              <Copy className="h-4 w-4" />
            </Button>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={onDelete}
              disabled={!canDelete}
              title="刪除此列"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </TableCell>
      </TableRow>
    </>
  );
}
