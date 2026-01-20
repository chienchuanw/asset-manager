-- 恢復信用卡已使用額度不能超過信用額度的約束
ALTER TABLE credit_cards ADD CONSTRAINT check_used_credit_limit CHECK (used_credit <= credit_limit);
