-- 還原為僅允許 TWD。
-- 警告：執行前必須確保現有資料不存在 USD 列，否則會失敗。
ALTER TABLE cash_flows DROP CONSTRAINT IF EXISTS cash_flows_currency_check;
ALTER TABLE cash_flows ADD CONSTRAINT cash_flows_currency_check
    CHECK (currency = 'TWD');
