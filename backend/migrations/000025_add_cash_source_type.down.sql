-- Rollback: Remove 'cash' from source type constraint

-- Remove updated constraint
ALTER TABLE cash_flows DROP CONSTRAINT IF EXISTS cash_flows_source_type_check;

-- Restore previous constraint without 'cash'
ALTER TABLE cash_flows 
ADD CONSTRAINT cash_flows_source_type_check 
CHECK (source_type IN ('manual', 'subscription', 'installment', 'bank_account', 'credit_card'));

-- Restore previous comments
COMMENT ON COLUMN cash_flows.source_type IS '來源類型（manual: 手動建立/現金交易, subscription: 訂閱自動產生, installment: 分期自動產生, bank_account: 銀行帳戶交易, credit_card: 信用卡交易）';
COMMENT ON COLUMN cash_flows.source_id IS '來源 ID（關聯到訂閱、分期、銀行帳戶或信用卡的 ID）';
COMMENT ON COLUMN cash_flows.target_type IS '轉帳目標類型 (用於 transfer_in/transfer_out,例如: credit_card, bank_account)';

