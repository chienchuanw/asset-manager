"use client";

import { useState, useEffect } from "react";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Plus,
  Copy,
  Trash2,
  ChevronDown,
  ChevronUp,
  Loader2,
} from "lucide-react";

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
}

/**
 * 批次新增交易對話框
 *
 * 使用表格式介面讓使用者一次新增多筆交易
 * 支援動態新增/刪除/複製列
 */
export function BatchAddTransactionDialog({
  onSuccess,
}: BatchAddTransactionDialogProps) {
  const [open, setOpen] = useState(false);
  const [expandedRows, setExpandedRows] = useState<Set<number>>(new Set());

  // 建立批次交易 mutation
  const createBatchMutation = useCreateTransactionsBatch({
    onSuccess: () => {
      setOpen(false);
      form.reset();
      setExpandedRows(new Set());
      onSuccess?.();
    },
  });

  // 表單設定
  const form = useForm<BatchCreateTransactionsFormData>({
    resolver: zodResolver(batchCreateTransactionsSchema),
    defaultValues: {
      asset_type: AssetType.TW_STOCK,
      currency: Currency.TWD,
      transactions: Array(5)
        .fill(null)
        .map(() => ({
          date: new Date().toISOString().split("T")[0],
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

  // 使用 useFieldArray 管理動態列
  const { fields, append, remove, insert } = useFieldArray({
    control: form.control,
    name: "transactions",
  });

  // 監聽資產類型變化，自動更新幣別
  const assetType = form.watch("asset_type");
  useEffect(() => {
    if (assetType === AssetType.TW_STOCK) {
      form.setValue("currency", Currency.TWD);
    } else if (assetType === AssetType.US_STOCK) {
      form.setValue("currency", Currency.USD);
    }
    // 加密貨幣可以是 TWD 或 USD，不自動變更
  }, [assetType, form]);

  // 新增一列
  const handleAddRow = () => {
    append({
      date: new Date().toISOString().split("T")[0],
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
      // 如果刪除的列是展開的，從展開集合中移除
      setExpandedRows((prev) => {
        const newSet = new Set(prev);
        newSet.delete(index);
        return newSet;
      });
    }
  };

  // 切換列的展開狀態
  const toggleRow = (index: number) => {
    setExpandedRows((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(index)) {
        newSet.delete(index);
      } else {
        newSet.add(index);
      }
      return newSet;
    });
  };

  // 送出表單
  const onSubmit = (data: BatchCreateTransactionsFormData) => {
    // 轉換資料格式
    const transactions = data.transactions.map((tx) => ({
      date: new Date(tx.date).toISOString(),
      asset_type: data.asset_type,
      symbol: tx.symbol,
      name: tx.name,
      type: tx.type,
      quantity: tx.quantity,
      price: tx.price,
      amount: tx.amount,
      fee: tx.fee,
      tax: tx.tax,
      currency: data.currency,
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

            {/* 資產類型和幣別選擇 */}
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="asset_type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>資產類型</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="選擇資產類型" />
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
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="選擇幣別" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value={Currency.TWD}>TWD</SelectItem>
                        <SelectItem value={Currency.USD}>USD</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* 交易列表表格 */}
            <div className="border rounded-lg overflow-hidden">
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-12">#</TableHead>
                      <TableHead className="min-w-[140px]">日期</TableHead>
                      <TableHead className="min-w-[100px]">代碼</TableHead>
                      <TableHead className="min-w-[120px]">名稱</TableHead>
                      <TableHead className="min-w-[100px]">類型</TableHead>
                      <TableHead className="min-w-[100px]">數量</TableHead>
                      <TableHead className="min-w-[100px]">價格</TableHead>
                      <TableHead className="min-w-[120px]">金額</TableHead>
                      <TableHead className="min-w-[100px]">手續費</TableHead>
                      <TableHead className="w-32">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {fields.map((field, index) => (
                      <TransactionRow
                        key={field.id}
                        index={index}
                        form={form}
                        isExpanded={expandedRows.has(index)}
                        onToggleExpand={() => toggleRow(index)}
                        onCopy={() => handleCopyRow(index)}
                        onDelete={() => handleDeleteRow(index)}
                        canDelete={fields.length > 1}
                        assetType={assetType}
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
  isExpanded: boolean;
  onToggleExpand: () => void;
  onCopy: () => void;
  onDelete: () => void;
  canDelete: boolean;
  assetType: AssetType;
}

function TransactionRow({
  index,
  form,
  isExpanded,
  onToggleExpand,
  onCopy,
  onDelete,
  canDelete,
  assetType,
}: TransactionRowProps) {
  // 監聽數量和價格變化，自動計算金額
  const quantity = form.watch(`transactions.${index}.quantity`);
  const price = form.watch(`transactions.${index}.price`);

  useEffect(() => {
    const amount = quantity * price;
    if (!isNaN(amount) && amount >= 0) {
      form.setValue(`transactions.${index}.amount`, amount);
    }
  }, [quantity, price, index, form]);

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
              <FormItem>
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

        {/* 操作按鈕 */}
        <TableCell>
          <div className="flex gap-1">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={onToggleExpand}
              title={isExpanded ? "收合" : "展開更多欄位"}
            >
              {isExpanded ? (
                <ChevronUp className="h-4 w-4" />
              ) : (
                <ChevronDown className="h-4 w-4" />
              )}
            </Button>
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

      {/* 展開的選填欄位 */}
      {isExpanded && (
        <TableRow className="bg-muted/30">
          <TableCell colSpan={10} className="p-0">
            <div className="px-6 py-4 space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* 交易稅 */}
                <FormField
                  control={form.control}
                  name={`transactions.${index}.tax`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-sm font-medium">
                        交易稅（選填）
                      </FormLabel>
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
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* 備註 */}
                <FormField
                  control={form.control}
                  name={`transactions.${index}.note`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-sm font-medium">
                        備註（選填）
                      </FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="輸入備註..."
                          {...field}
                          value={field.value ?? ""}
                          onChange={(e) =>
                            field.onChange(e.target.value || null)
                          }
                          rows={2}
                          className="resize-none"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </div>
          </TableCell>
        </TableRow>
      )}
    </>
  );
}
