-- 刪除觸發器
DROP TRIGGER IF EXISTS update_credit_cards_updated_at ON credit_cards;
DROP TRIGGER IF EXISTS update_bank_accounts_updated_at ON bank_accounts;

-- 刪除索引
DROP INDEX IF EXISTS idx_credit_cards_payment_due_day;
DROP INDEX IF EXISTS idx_credit_cards_billing_day;
DROP INDEX IF EXISTS idx_bank_accounts_currency;

-- 刪除資料表
DROP TABLE IF EXISTS credit_cards;
DROP TABLE IF EXISTS bank_accounts;

