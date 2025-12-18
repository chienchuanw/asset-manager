-- 為訂閱表新增付款方式相關欄位
-- payment_method: 付款方式 (cash, bank_account, credit_card)
-- account_id: 帳戶 ID (銀行帳戶或信用卡 ID)

-- 新增 payment_method 欄位到 subscriptions 表
ALTER TABLE subscriptions
ADD COLUMN payment_method VARCHAR(20) NOT NULL DEFAULT 'cash'
    CHECK (payment_method IN ('cash', 'bank_account', 'credit_card'));

-- 新增 account_id 欄位到 subscriptions 表
ALTER TABLE subscriptions
ADD COLUMN account_id UUID;

-- 為 subscriptions 表建立索引
CREATE INDEX idx_subscriptions_payment_method ON subscriptions(payment_method);
CREATE INDEX idx_subscriptions_account_id ON subscriptions(account_id);

-- 新增 payment_method 欄位到 installments 表
ALTER TABLE installments
ADD COLUMN payment_method VARCHAR(20) NOT NULL DEFAULT 'cash'
    CHECK (payment_method IN ('cash', 'bank_account', 'credit_card'));

-- 新增 account_id 欄位到 installments 表
ALTER TABLE installments
ADD COLUMN account_id UUID;

-- 為 installments 表建立索引
CREATE INDEX idx_installments_payment_method ON installments(payment_method);
CREATE INDEX idx_installments_account_id ON installments(account_id);

-- 註解說明
COMMENT ON COLUMN subscriptions.payment_method IS '付款方式（cash: 現金, bank_account: 銀行帳戶, credit_card: 信用卡）';
COMMENT ON COLUMN subscriptions.account_id IS '帳戶 ID（銀行帳戶或信用卡的 ID，現金付款時為 NULL）';
COMMENT ON COLUMN installments.payment_method IS '付款方式（cash: 現金, bank_account: 銀行帳戶, credit_card: 信用卡）';
COMMENT ON COLUMN installments.account_id IS '帳戶 ID（銀行帳戶或信用卡的 ID，現金付款時為 NULL）';

