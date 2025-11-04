-- 更新現金流來源類型約束，新增銀行帳戶和信用卡選項
-- 移除舊的約束
ALTER TABLE cash_flows DROP CONSTRAINT IF EXISTS cash_flows_source_type_check;

-- 新增更新後的約束，包含銀行帳戶和信用卡
ALTER TABLE cash_flows 
ADD CONSTRAINT cash_flows_source_type_check 
CHECK (source_type IN ('manual', 'subscription', 'installment', 'bank_account', 'credit_card'));

-- 更新註解說明
COMMENT ON COLUMN cash_flows.source_type IS '來源類型（manual: 手動建立/現金交易, subscription: 訂閱自動產生, installment: 分期自動產生, bank_account: 銀行帳戶交易, credit_card: 信用卡交易）';
COMMENT ON COLUMN cash_flows.source_id IS '來源 ID（關聯到訂閱、分期、銀行帳戶或信用卡的 ID）';
