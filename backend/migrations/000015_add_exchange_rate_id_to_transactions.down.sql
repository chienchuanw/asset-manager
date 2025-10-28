-- 移除 exchange_rate_id 索引
DROP INDEX IF EXISTS idx_transactions_exchange_rate_id;

-- 移除 exchange_rate_id 欄位
ALTER TABLE transactions
DROP COLUMN IF EXISTS exchange_rate_id;

