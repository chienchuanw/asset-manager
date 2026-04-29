/**
 * Balance reconciliation 相關型別。
 * Mirror 後端 internal/models/reconciliation.go 的 DTO。
 */

export type ReconcileTargetType = "bank_account" | "credit_card";

export interface ReconcileBankAccountInput {
  new_balance: number;
  date: string; // ISO 8601
  note?: string;
}

export interface ReconcileCreditCardInput {
  new_used_credit?: number;
  new_credit_limit?: number;
  date: string;
  note?: string;
}

export interface ReconcileBatchItem {
  target_type: ReconcileTargetType;
  target_id: string;
  new_balance?: number;
  new_used_credit?: number;
  new_credit_limit?: number;
}

export interface ReconcileBatchInput {
  date: string;
  note?: string;
  items: ReconcileBatchItem[];
}

export interface ReconcileResult {
  target_type: ReconcileTargetType;
  target_id: string;
  delta: number;
  cash_flow_id?: string;
  new_balance?: number;
  new_used_credit?: number;
  new_credit_limit?: number;
}
