-- 移除 exchange_rates 資料表的 updated_at 欄位
ALTER TABLE exchange_rates 
DROP COLUMN IF EXISTS updated_at;

