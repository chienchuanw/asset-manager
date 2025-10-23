-- 刪除索引
DROP INDEX IF EXISTS idx_exchange_rates_currencies;
DROP INDEX IF EXISTS idx_exchange_rates_date;

-- 刪除匯率表
DROP TABLE IF EXISTS exchange_rates;

