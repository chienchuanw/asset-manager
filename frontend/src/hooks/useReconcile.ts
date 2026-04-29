import { useMutation, useQueryClient } from "@tanstack/react-query";
import { reconciliationAPI } from "@/lib/api/reconciliation";
import type {
  ReconcileBankAccountInput,
  ReconcileBatchInput,
  ReconcileCreditCardInput,
  ReconcileResult,
} from "@/types/reconciliation";
import type { APIError } from "@/lib/api/client";

/**
 * 校準銀行帳戶餘額
 */
export function useReconcileBankAccount() {
  const qc = useQueryClient();
  return useMutation<
    ReconcileResult,
    APIError,
    { id: string; input: ReconcileBankAccountInput }
  >({
    mutationFn: ({ id, input }) =>
      reconciliationAPI.reconcileBankAccount(id, input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["bankAccounts"] });
      qc.invalidateQueries({ queryKey: ["cashFlows"] });
      qc.invalidateQueries({ queryKey: ["analytics"] });
    },
  });
}

/**
 * 校準信用卡 used_credit / credit_limit
 */
export function useReconcileCreditCard() {
  const qc = useQueryClient();
  return useMutation<
    ReconcileResult,
    APIError,
    { id: string; input: ReconcileCreditCardInput }
  >({
    mutationFn: ({ id, input }) =>
      reconciliationAPI.reconcileCreditCard(id, input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["creditCards"] });
      qc.invalidateQueries({ queryKey: ["cashFlows"] });
      qc.invalidateQueries({ queryKey: ["analytics"] });
    },
  });
}

/**
 * 批次校準銀行帳戶 / 信用卡
 */
export function useReconcileBatch() {
  const qc = useQueryClient();
  return useMutation<ReconcileResult[], APIError, ReconcileBatchInput>({
    mutationFn: (input) => reconciliationAPI.reconcileBatch(input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["bankAccounts"] });
      qc.invalidateQueries({ queryKey: ["creditCards"] });
      qc.invalidateQueries({ queryKey: ["cashFlows"] });
      qc.invalidateQueries({ queryKey: ["analytics"] });
    },
  });
}
