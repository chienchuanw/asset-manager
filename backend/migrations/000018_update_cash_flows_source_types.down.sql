-- 回復現金流來源類型約束到原始狀態
-- 移除更新後的約束
ALTER TABLE cash_flows DROP CONSTRAINT IF EXISTS cash_flows_source_type_check;

-- 恢復原始約束
ALTER TABLE cash_flows 
ADD CONSTRAINT cash_flows_source_type_check 
CHECK (source_type IN ('manual', 'subscription', 'installment'));

-- 恢復原始註解說明
COMMENT ON COLUMN cash_flows.source_type IS '來源類型（manual: 手動建立, subscription: 訂閱自動產生, installment: 分期自動產生）';
COMMENT ON COLUMN cash_flows.source_id IS '來源 ID（關聯到訂閱或分期的 ID）';
