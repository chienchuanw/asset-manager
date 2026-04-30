-- 放寬 cash_flows.currency 限制，允許 TWD 與 USD。
-- 動機：bank_accounts 已支援 USD；reconciliation 產生的調整 cash flow 應沿用帳戶幣別。
ALTER TABLE cash_flows DROP CONSTRAINT IF EXISTS cash_flows_currency_check;
ALTER TABLE cash_flows ADD CONSTRAINT cash_flows_currency_check
    CHECK (currency IN ('TWD', 'USD'));
