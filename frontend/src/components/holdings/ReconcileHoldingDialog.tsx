/**
 * 單筆持倉對帳 dialog：使用者填入目標數量與平均成本，先預覽 diff，再確認寫入。
 * 此操作不結算實現損益，也不影響現金。
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
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Loader2, AlertCircle } from "lucide-react";
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

interface Props {
  open: boolean;
  onOpenChange: (v: boolean) => void;
  holding: Holding;
}

export function ReconcileHoldingDialog({ open, onOpenChange, holding }: Props) {
  const t = useTranslations("holdings.reconcile");
  const [quantity, setQuantity] = useState(String(holding.quantity));
  const [avgCost, setAvgCost] = useState(String(holding.avg_cost_original));
  const [reason, setReason] = useState("");
  const [previewItem, setPreviewItem] =
    useState<HoldingReconcilePreviewItem | null>(null);
  const [error, setError] = useState<string | null>(null);

  const previewMut = useReconcileHoldingsPreview();
  const executeMut = useReconcileHoldings();
  const submitting = previewMut.isPending || executeMut.isPending;

  // 對話框開啟時重設欄位至當前值
  useEffect(() => {
    if (open) {
      setQuantity(String(holding.quantity));
      setAvgCost(String(holding.avg_cost_original));
      setReason("");
      setPreviewItem(null);
      setError(null);
    }
  }, [open, holding.quantity, holding.avg_cost_original]);

  const buildItem = (): HoldingReconcileItem | null => {
    const q = Number(quantity);
    const a = Number(avgCost);
    if (Number.isNaN(q) || q < 0) {
      setError(t("error.invalidQuantity"));
      return null;
    }
    if (Number.isNaN(a) || a < 0 || (a === 0 && q > 0)) {
      setError(t("error.invalidAvgCost"));
      return null;
    }
    return {
      asset_type: holding.asset_type as HoldingReconcileItem["asset_type"],
      symbol: holding.symbol,
      name: holding.name,
      currency: holding.currency as HoldingReconcileItem["currency"],
      target_quantity: q,
      target_avg_cost: a,
      reason: reason || undefined,
    };
  };

  const onPreview = async () => {
    setError(null);
    const item = buildItem();
    if (!item) return;
    try {
      const res = await previewMut.mutateAsync([item]);
      setPreviewItem(res.items[0] ?? null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  };

  const onConfirm = async () => {
    if (!previewItem) return;
    setError(null);
    try {
      await executeMut.mutateAsync([previewItem.item]);
      toast.success(t("toast.success"));
      onOpenChange(false);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("dialogTitle")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>

        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {!previewItem ? (
          <div className="space-y-3 py-2">
            <div className="grid grid-cols-2 gap-3">
              <div>
                <Label className="mb-1 block text-xs text-muted-foreground">
                  {t("field.symbol")}
                </Label>
                <Input value={holding.symbol} disabled />
              </div>
              <div>
                <Label className="mb-1 block text-xs text-muted-foreground">
                  {t("field.currency")}
                </Label>
                <Input value={holding.currency} disabled />
              </div>
            </div>
            <div>
              <Label className="mb-1 block text-xs text-muted-foreground">
                {t("field.quantity")}
              </Label>
              <Input
                type="number"
                step="0.00000001"
                value={quantity}
                onChange={(e) => setQuantity(e.target.value)}
                data-testid="reconcile-quantity"
              />
            </div>
            <div>
              <Label className="mb-1 block text-xs text-muted-foreground">
                {t("field.avgCost")}
              </Label>
              <Input
                type="number"
                step="0.00000001"
                value={avgCost}
                onChange={(e) => setAvgCost(e.target.value)}
                data-testid="reconcile-avg-cost"
              />
            </div>
            <div>
              <Label className="mb-1 block text-xs text-muted-foreground">
                {t("field.reason")}
              </Label>
              <Input
                value={reason}
                onChange={(e) => setReason(e.target.value)}
              />
            </div>
          </div>
        ) : (
          <div className="space-y-2 py-2 text-sm">
            <PreviewRow item={previewItem} t={t} />
          </div>
        )}

        <DialogFooter>
          {!previewItem ? (
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
                disabled={submitting}
                data-testid="reconcile-preview-btn"
              >
                {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {t("preview.title")}
              </Button>
            </>
          ) : (
            <>
              <Button
                variant="outline"
                onClick={() => setPreviewItem(null)}
                disabled={submitting}
              >
                {t("cancel")}
              </Button>
              <Button
                onClick={onConfirm}
                disabled={submitting}
                data-testid="reconcile-confirm-btn"
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

function PreviewRow({
  item,
  t,
}: {
  item: HoldingReconcilePreviewItem;
  t: (k: string) => string;
}) {
  return (
    <div className="rounded-md border p-3">
      <div className="mb-1 font-medium">
        {item.item.symbol} —{" "}
        <ActionBadge action={item.action} t={t} />
      </div>
      <div className="text-xs text-muted-foreground">
        {t("preview.delta.qty")}: {item.prev_quantity} →{" "}
        {item.item.target_quantity} ({signed(item.quantity_delta)})
      </div>
      <div className="text-xs text-muted-foreground">
        {t("preview.delta.avgCost")}: {item.prev_avg_cost} →{" "}
        {item.item.target_avg_cost} ({signed(item.avg_cost_delta)})
      </div>
    </div>
  );
}

function ActionBadge({
  action,
  t,
}: {
  action: ReconcileAction;
  t: (k: string) => string;
}) {
  const colorByAction: Record<ReconcileAction, string> = {
    create: "text-green-600",
    update: "text-blue-600",
    liquidate: "text-orange-600",
    noop: "text-muted-foreground",
  };
  return (
    <span className={colorByAction[action]}>
      {t(`preview.action.${action}`)}
    </span>
  );
}

function signed(n: number): string {
  return n >= 0 ? `+${n}` : String(n);
}
