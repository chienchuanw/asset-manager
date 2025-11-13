/**
 * 不足數量修復對話框
 * 當 FIFO 計算發現賣出數量超過買入數量時，顯示此對話框讓使用者輸入當前實際持有股數
 */

"use client";

import { useState } from "react";
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
import { AlertCircle, Loader2 } from "lucide-react";
import type { APIWarning } from "@/types/transaction";

/**
 * Dialog Props
 */
interface InsufficientQuantityDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  warning: APIWarning | null;
  onFix: (data: {
    symbol: string;
    currentHolding: number;
    estimatedCost?: number;
  }) => Promise<void>;
}

export function InsufficientQuantityDialog({
  open,
  onOpenChange,
  warning,
  onFix,
}: InsufficientQuantityDialogProps) {
  const [currentHolding, setCurrentHolding] = useState("");
  const [estimatedCost, setEstimatedCost] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 重置表單
  const resetForm = () => {
    setCurrentHolding("");
    setEstimatedCost("");
    setError(null);
  };

  // 處理關閉
  const handleClose = () => {
    if (!isSubmitting) {
      resetForm();
      onOpenChange(false);
    }
  };

  // 處理提交
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!warning) return;

    // 驗證輸入
    const currentHoldingNum = parseFloat(currentHolding);
    if (isNaN(currentHoldingNum) || currentHoldingNum <= 0) {
      setError("請輸入有效的持有股數（必須大於 0）");
      return;
    }

    const estimatedCostNum = estimatedCost
      ? parseFloat(estimatedCost)
      : undefined;
    if (estimatedCost && (isNaN(estimatedCostNum!) || estimatedCostNum! <= 0)) {
      setError("請輸入有效的估計成本（必須大於 0）");
      return;
    }

    try {
      setIsSubmitting(true);
      await onFix({
        symbol: warning.symbol,
        currentHolding: currentHoldingNum,
        estimatedCost: estimatedCostNum,
      });
      resetForm();
      onOpenChange(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "修復失敗，請稍後再試");
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!warning) return null;

  const details = warning.details;

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>修復資料不一致問題</DialogTitle>
            <DialogDescription>
              系統偵測到 <strong>{warning.symbol}</strong> 的交易記錄不一致
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            {/* 警告訊息 */}
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                <div className="text-sm">
                  <p className="font-medium mb-2">{warning.message}</p>
                  {details && (
                    <div className="text-xs text-muted-foreground space-y-1">
                      <p>需要數量：{details.required}</p>
                      <p>可用數量：{details.available}</p>
                      <p>缺少數量：{details.missing}</p>
                    </div>
                  )}
                </div>
              </AlertDescription>
            </Alert>

            {/* 當前持有股數輸入 */}
            <div className="grid gap-2">
              <Label htmlFor="current-holding">
                當前實際持有股數 <span className="text-red-500">*</span>
              </Label>
              <Input
                id="current-holding"
                type="number"
                step="0.0001"
                placeholder="請輸入您目前實際持有的股數"
                value={currentHolding}
                onChange={(e) => setCurrentHolding(e.target.value)}
                disabled={isSubmitting}
                required
              />
              <p className="text-xs text-muted-foreground">
                請輸入您在券商帳戶中實際持有的股數（可包含小數）
              </p>
            </div>
          </div>

          {/* 錯誤訊息 */}
          {error && (
            <Alert variant="destructive" className="mb-4">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={isSubmitting}
            >
              取消
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              確認修復
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

