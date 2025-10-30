-- 建立銀行帳戶表
CREATE TABLE IF NOT EXISTS bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bank_name VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) NOT NULL,
    account_number_last4 VARCHAR(4) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'TWD' CHECK (currency IN ('TWD', 'USD')),
    balance DECIMAL(20, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 建立信用卡表
CREATE TABLE IF NOT EXISTS credit_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issuing_bank VARCHAR(255) NOT NULL,
    card_name VARCHAR(255) NOT NULL,
    card_number_last4 VARCHAR(4) NOT NULL,
    billing_day INT NOT NULL CHECK (billing_day >= 1 AND billing_day <= 31),
    payment_due_day INT NOT NULL CHECK (payment_due_day >= 1 AND payment_due_day <= 31),
    credit_limit DECIMAL(20, 2) NOT NULL CHECK (credit_limit > 0),
    used_credit DECIMAL(20, 2) NOT NULL DEFAULT 0 CHECK (used_credit >= 0),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_used_credit_limit CHECK (used_credit <= credit_limit)
);

-- 建立索引以提升查詢效能
CREATE INDEX idx_bank_accounts_currency ON bank_accounts(currency);
CREATE INDEX idx_credit_cards_billing_day ON credit_cards(billing_day);
CREATE INDEX idx_credit_cards_payment_due_day ON credit_cards(payment_due_day);

-- 建立更新時間的觸發器
CREATE TRIGGER update_bank_accounts_updated_at
    BEFORE UPDATE ON bank_accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_credit_cards_updated_at
    BEFORE UPDATE ON credit_cards
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

