-- 新增 exchange_rate_id 外鍵到 transactions 表
ALTER TABLE transactions
ADD COLUMN exchange_rate_id INTEGER REFERENCES exchange_rates(id) ON DELETE SET NULL;

-- 為 exchange_rate_id 欄位建立索引
CREATE INDEX idx_transactions_exchange_rate_id ON transactions(exchange_rate_id);

-- 註解說明
COMMENT ON COLUMN transactions.exchange_rate_id IS '關聯的匯率記錄 ID（僅用於非 TWD 交易，TWD 交易此欄位為 NULL）';

