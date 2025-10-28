-- 移除 tax 欄位
ALTER TABLE transactions
DROP COLUMN IF EXISTS tax;

