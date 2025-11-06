-- 新增 target_type 和 target_id 欄位到 cash_flows 表
-- 用於記錄轉帳目標 (例如:繳信用卡費時的目標信用卡)

ALTER TABLE cash_flows 
ADD COLUMN target_type VARCHAR(50),
ADD COLUMN target_id UUID;

-- 新增索引以提升查詢效能
CREATE INDEX idx_cash_flows_target ON cash_flows(target_type, target_id);

-- 新增欄位註解
COMMENT ON COLUMN cash_flows.target_type IS '轉帳目標類型 (用於 transfer_in/transfer_out,例如: credit_card, bank_account)';
COMMENT ON COLUMN cash_flows.target_id IS '轉帳目標ID (用於 transfer_in/transfer_out)';

