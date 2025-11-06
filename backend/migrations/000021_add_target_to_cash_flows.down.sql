-- 移除 target_type 和 target_id 欄位

DROP INDEX IF EXISTS idx_cash_flows_target;

ALTER TABLE cash_flows 
DROP COLUMN IF EXISTS target_type,
DROP COLUMN IF EXISTS target_id;

