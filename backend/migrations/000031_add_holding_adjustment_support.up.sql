-- 放寬 transaction_type CHECK constraint，新增 'adjustment' 用於持倉對帳
ALTER TABLE transactions DROP CONSTRAINT transactions_transaction_type_check;
ALTER TABLE transactions ADD CONSTRAINT transactions_transaction_type_check
    CHECK (transaction_type IN ('buy', 'sell', 'dividend', 'fee', 'adjustment'));

-- 新增對帳稽核欄位（僅 adjustment 類型會填值）
ALTER TABLE transactions
    ADD COLUMN adjustment_prev_quantity DECIMAL(20, 8),
    ADD COLUMN adjustment_prev_avg_cost DECIMAL(20, 8),
    ADD COLUMN adjustment_reason TEXT,
    ADD COLUMN broker_account_id UUID;

COMMENT ON COLUMN transactions.adjustment_prev_quantity IS '對帳前數量（僅 adjustment 類型）';
COMMENT ON COLUMN transactions.adjustment_prev_avg_cost IS '對帳前平均成本（僅 adjustment 類型）';
COMMENT ON COLUMN transactions.adjustment_reason IS '對帳原因（僅 adjustment 類型，可選）';
COMMENT ON COLUMN transactions.broker_account_id IS '預留欄位，未來券商帳戶維度';
