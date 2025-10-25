-- 刪除索引
DROP INDEX IF EXISTS idx_cash_flows_source;
DROP INDEX IF EXISTS idx_cash_flows_source_id;
DROP INDEX IF EXISTS idx_cash_flows_source_type;

-- 刪除欄位
ALTER TABLE cash_flows 
DROP COLUMN IF EXISTS source_id,
DROP COLUMN IF EXISTS source_type;

