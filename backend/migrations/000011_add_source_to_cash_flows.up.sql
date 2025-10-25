-- 新增來源類型和來源 ID 欄位到現金流表
ALTER TABLE cash_flows 
ADD COLUMN source_type VARCHAR(20) CHECK (source_type IN ('manual', 'subscription', 'installment')),
ADD COLUMN source_id UUID;

-- 建立索引以提升查詢效能
CREATE INDEX idx_cash_flows_source_type ON cash_flows(source_type);
CREATE INDEX idx_cash_flows_source_id ON cash_flows(source_id);
CREATE INDEX idx_cash_flows_source ON cash_flows(source_type, source_id);

-- 將現有的現金流記錄標記為手動建立
UPDATE cash_flows SET source_type = 'manual' WHERE source_type IS NULL;

-- 註解說明
COMMENT ON COLUMN cash_flows.source_type IS '來源類型（manual: 手動建立, subscription: 訂閱自動產生, installment: 分期自動產生）';
COMMENT ON COLUMN cash_flows.source_id IS '來源 ID（關聯到訂閱或分期的 ID）';

