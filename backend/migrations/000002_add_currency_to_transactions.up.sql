-- 新增 currency 欄位到 transactions 表
ALTER TABLE transactions
ADD COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'TWD'
CHECK (currency IN ('TWD', 'USD'));

-- 為 currency 欄位建立索引
CREATE INDEX idx_transactions_currency ON transactions(currency);

-- 註解說明
COMMENT ON COLUMN transactions.currency IS '交易幣別（TWD: 新台幣, USD: 美金）';

