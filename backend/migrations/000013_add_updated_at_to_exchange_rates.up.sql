-- 為 exchange_rates 資料表新增 updated_at 欄位
ALTER TABLE exchange_rates 
ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- 為現有資料設定 updated_at = created_at
UPDATE exchange_rates 
SET updated_at = created_at 
WHERE updated_at IS NULL;

