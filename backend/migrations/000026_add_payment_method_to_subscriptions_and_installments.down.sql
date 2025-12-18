-- 回滾訂閱和分期表的付款方式欄位

-- 移除 installments 表的索引和欄位
DROP INDEX IF EXISTS idx_installments_account_id;
DROP INDEX IF EXISTS idx_installments_payment_method;
ALTER TABLE installments DROP COLUMN IF EXISTS account_id;
ALTER TABLE installments DROP COLUMN IF EXISTS payment_method;

-- 移除 subscriptions 表的索引和欄位
DROP INDEX IF EXISTS idx_subscriptions_account_id;
DROP INDEX IF EXISTS idx_subscriptions_payment_method;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS account_id;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS payment_method;

