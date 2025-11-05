-- 回滾：移除轉帳類型，恢復原始的現金流類型約束

-- 移除包含轉帳類型的約束
ALTER TABLE cash_flows DROP CONSTRAINT IF EXISTS cash_flows_type_check;

-- 恢復原始的 type 約束（只包含 income 和 expense）
ALTER TABLE cash_flows ADD CONSTRAINT cash_flows_type_check 
    CHECK (type IN ('income', 'expense'));

-- 移除註解
COMMENT ON COLUMN cash_flows.type IS NULL;
