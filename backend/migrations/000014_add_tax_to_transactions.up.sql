-- 新增 tax 欄位到 transactions 表
ALTER TABLE transactions
ADD COLUMN tax DECIMAL(20, 2) DEFAULT 0 CHECK (tax >= 0);

-- 註解說明
COMMENT ON COLUMN transactions.tax IS '交易稅（選填，預設為 0）';

