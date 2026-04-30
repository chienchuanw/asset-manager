/**
 * 批次持倉對帳 dialog：預載入所有現有 holdings，使用者可一次編輯多筆，並可新增首次建倉的標的。
 * 兩段式流程：編輯 → 預覽 (含 action 標籤、noop 灰階) → 確認寫入。
 */

"use client";

import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Loader2, AlertCircle, Plus, X } from "lucide-react";
import { toast } from "sonner";
import {
  useReconcileHoldings,
  useReconcileHoldingsPreview,
} from "@/hooks/useReconcileHoldings";
import type {
  HoldingReconcileItem,
  HoldingReconcilePreviewItem,
  ReconcileAction,
} from "@/lib/api/holdings";
import type { Holding } from "@/types/holding";

type Row = HoldingReconcileItem & { __existing: boolean };

interface Props {
  open: boolean;
  onOpenChange: (v: boolean) => void;
  currentHoldings: Holding[];
}

export function BatchReconcileHoldingsDialog({
  open,
  onOpenChange,
  currentHoldings,
}: Props) {
  const t = useTranslations("holdings.reconcile");
  const [rows, setRows] = useState<Row[]>([]);
  const [preview, setPreview] = useState<HoldingReconcilePreviewItem[] | null>(
    null
  );
  const [error, setError] = useState<string | null>(null);

  const previewMut = useReconcileHoldingsPreview();
  const executeMut = useReconcileHoldings();
  const submitting = previewMut.isPending || executeMut.isPending;

  // 開啟時用當前 holdings 預載入列；只考慮非現金資產
  useEffect(() => {
    if (open) {
      setRows(
        currentHoldings
          .filter((h) => h.asset_type !== "cash")
          .map((h) => ({
            asset_type: h.asset_type as HoldingReconcileItem["asset_type"],
            symbol: h.symbol,
            name: h.name,
            currency: h.currency as HoldingReconcileItem["currency"],
            target_quantity: h.quantity,
            target_avg_cost: h.avg_cost_original,
            __existing: true,
          }))
      );
      setPreview(null);
      setError(null);
    }
  }, [open, currentHoldings]);

  const updateRow = (i: number, patch: Partial<Row>) =>
    setRows((rs) => rs.map((r, idx) => (idx === i ? { ...r, ...patch } : r)));

  const addRow = () =>
    setRows((rs) => [
      ...rs,
      {
        asset_type: "us-stock",
        symbol: "",
        name: "",
        currency: "USD",
        target_quantity: 0,
        target_avg_cost: 0,
        __existing: false,
      },
    ]);

  const removeRow = (i: number) =>
    setRows((rs) => rs.filter((_, idx) => idx !== i));

  const validate = (): HoldingReconcileItem[] | null => {
    const items: HoldingReconcileItem[] = [];
    for (const r of rows) {
      const q = Number(r.target_quantity);
      const a = Number(r.target_avg_cost);
      if (!r.symbol) {
        setError(`${t("field.symbol")} required`);
        return null;
      }
      if (Number.isNaN(q) || q < 0) {
        setError(t("error.invalidQuantity"));
        return null;
      }
      if (Number.isNaN(a) || a < 0 || (a === 0 && q > 0)) {
        setError(t("error.invalidAvgCost"));
        return null;
      }
      if (!r.__existing && !r.name) {
        setError(t("error.nameRequired"));
        return null;
      }
      items.push({
        asset_type: r.asset_type,
        symbol: r.symbol,
        name: r.name,
        currency: r.currency,
        target_quantity: q,
        target_avg_cost: a,
      });
    }
    return items;
  };

  const onPreview = async () => {
    setError(null);
    const items = validate();
    if (!items) return;
    try {
      const res = await previewMut.mutateAsync(items);
      setPreview(res.items);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  };

  const onConfirm = async () => {
    if (!preview) return;
    setError(null);
    const items = preview.filter((p) => p.action !== "noop").map((p) => p.item);
    if (items.length === 0) {
      onOpenChange(false);
      return;
    }
    try {
      await executeMut.mutateAsync(items);
      toast.success(t("toast.success"));
      onOpenChange(false);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="flex max-h-[90vh] max-w-3xl flex-col">
        <DialogHeader>
          <DialogTitle>{t("batchDialogTitle")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>

        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <div className="flex-1 space-y-2 overflow-y-auto py-2">
          {!preview ? (
            <>
              <div className="grid grid-cols-12 gap-2 text-xs text-muted-foreground">
                <div className="col-span-2">{t("field.symbol")}</div>
                <div className="col-span-3">{t("field.name")}</div>
                <div className="col-span-2">{t("field.currency")}</div>
                <div className="col-span-2">{t("field.quantity")}</div>
                <div className="col-span-2">{t("field.avgCost")}</div>
                <div className="col-span-1" />
              </div>
              {rows.map((r, i) => (
                <div key={i} className="grid grid-cols-12 items-center gap-2">
                  <Input
                    className="col-span-2"
                    disabled={r.__existing}
                    value={r.symbol}
                    onChange={(e) => updateRow(i, { symbol: e.target.value })}
                  />
                  <Input
                    className="col-span-3"
                    disabled={r.__existing}
                    value={r.name ?? ""}
                    onChange={(e) => updateRow(i, { name: e.target.value })}
                  />
                  <Select
                    disabled={r.__existing}
                    value={r.currency}
                    onValueChange={(v) =>
                      updateRow(i, {
                        currency: v as HoldingReconcileItem["currency"],
                      })
                    }
                  >
                    <SelectTrigger className="col-span-2">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="USD">USD</SelectItem>
                      <SelectItem value="TWD">TWD</SelectItem>
                    </SelectContent>
                  </Select>
                  <Input
                    className="col-span-2"
                    type="number"
                    step="0.00000001"
                    value={r.target_quantity}
                    onChange={(e) =>
                      updateRow(i, { target_quantity: Number(e.target.value) })
                    }
                  />
                  <Input
                    className="col-span-2"
                    type="number"
                    step="0.00000001"
                    value={r.target_avg_cost}
                    onChange={(e) =>
                      updateRow(i, { target_avg_cost: Number(e.target.value) })
                    }
                  />
                  {!r.__existing && (
                    <Button
                      type="button"
                      size="icon"
                      variant="ghost"
                      className="col-span-1"
                      onClick={() => removeRow(i)}
                      aria-label={t("removeRow")}
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              ))}
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={addRow}
                data-testid="reconcile-batch-add-row"
              >
                <Plus className="mr-1 h-4 w-4" />
                {t("addRow")}
              </Button>
            </>
          ) : (
            <PreviewList preview={preview} t={t} />
          )}
        </div>

        <DialogFooter className="border-t pt-3">
          {!preview ? (
            <>
              <Button
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={submitting}
              >
                {t("cancel")}
              </Button>
              <Button
                onClick={onPreview}
                disabled={submitting || rows.length === 0}
                data-testid="reconcile-batch-preview-btn"
              >
                {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {t("preview.title")}
              </Button>
            </>
          ) : (
            <>
              <Button
                variant="outline"
                onClick={() => setPreview(null)}
                disabled={submitting}
              >
                {t("cancel")}
              </Button>
              <Button
                onClick={onConfirm}
                disabled={submitting}
                data-testid="reconcile-batch-confirm-btn"
              >
                {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {t("confirm")}
              </Button>
            </>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function PreviewList({
  preview,
  t,
}: {
  preview: HoldingReconcilePreviewItem[];
  t: (k: string) => string;
}) {
  if (preview.length === 0) {
    return (
      <p className="text-center text-sm text-muted-foreground">
        {t("preview.noChanges")}
      </p>
    );
  }
  return (
    <div className="space-y-2 text-sm">
      {preview.map((p, i) => (
        <div
          key={i}
          className={`rounded-md border p-2 ${
            p.action === "noop" ? "opacity-50" : ""
          }`}
        >
          <div className="font-medium">
            {p.item.symbol} —{" "}
            <ActionLabel action={p.action} t={t} />
          </div>
          <div className="text-xs text-muted-foreground">
            {t("preview.delta.qty")}: {p.prev_quantity} →{" "}
            {p.item.target_quantity}
          </div>
          <div className="text-xs text-muted-foreground">
            {t("preview.delta.avgCost")}: {p.prev_avg_cost} →{" "}
            {p.item.target_avg_cost}
          </div>
        </div>
      ))}
    </div>
  );
}

function ActionLabel({
  action,
  t,
}: {
  action: ReconcileAction;
  t: (k: string) => string;
}) {
  const color: Record<ReconcileAction, string> = {
    create: "text-green-600",
    update: "text-blue-600",
    liquidate: "text-orange-600",
    noop: "text-muted-foreground",
  };
  return <span className={color[action]}>{t(`preview.action.${action}`)}</span>;
}
