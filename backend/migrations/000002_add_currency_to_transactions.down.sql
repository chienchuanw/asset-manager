-- 移除 currency 欄位的索引
DROP INDEX IF EXISTS idx_transactions_currency;

-- 移除 currency 欄位
ALTER TABLE transactions
DROP COLUMN IF EXISTS currency;

