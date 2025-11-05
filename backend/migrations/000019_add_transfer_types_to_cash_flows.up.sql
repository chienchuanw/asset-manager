-- 擴展現金流類型，新增轉帳相關類型
-- 修改 cash_flows 表的 type 欄位約束，新增 transfer_in 和 transfer_out

-- 移除現有的 type 約束
ALTER TABLE cash_flows DROP CONSTRAINT IF EXISTS cash_flows_type_check;

-- 新增包含轉帳類型的約束
ALTER TABLE cash_flows ADD CONSTRAINT cash_flows_type_check 
    CHECK (type IN ('income', 'expense', 'transfer_in', 'transfer_out'));

-- 新增註解說明各類型用途
COMMENT ON COLUMN cash_flows.type IS '現金流類型: income(收入), expense(支出), transfer_in(存入帳戶), transfer_out(從帳戶轉出)';
