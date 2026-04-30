import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  holdingsAPI,
  type HoldingReconcileItem,
  type HoldingReconcilePreview,
} from "@/lib/api/holdings";
import { APIError } from "@/lib/api/client";

/**
 * 持倉對帳 dry-run（preview）— 計算 diff，不寫入
 */
export function useReconcileHoldingsPreview() {
  return useMutation<HoldingReconcilePreview, APIError, HoldingReconcileItem[]>(
    {
      mutationFn: (items) => holdingsAPI.reconcilePreview(items),
    }
  );
}

/**
 * 持倉對帳實際執行；成功後 invalidate holdings / transactions / asset-snapshots queries
 */
export function useReconcileHoldings() {
  const qc = useQueryClient();
  return useMutation<HoldingReconcilePreview, APIError, HoldingReconcileItem[]>(
    {
      mutationFn: (items) => holdingsAPI.reconcile(items),
      onSuccess: () => {
        qc.invalidateQueries({ queryKey: ["holdings"] });
        qc.invalidateQueries({ queryKey: ["transactions"] });
        qc.invalidateQueries({ queryKey: ["asset-snapshots"] });
      },
    }
  );
}
