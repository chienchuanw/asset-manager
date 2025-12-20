-- 移除 sort_order 相關變更

-- 移除索引
DROP INDEX IF EXISTS idx_cash_flow_categories_sort_order;

-- 移除 sort_order 欄位
ALTER TABLE cash_flow_categories
DROP COLUMN IF EXISTS sort_order;

