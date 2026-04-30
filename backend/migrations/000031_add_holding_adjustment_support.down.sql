-- 先清掉所有 adjustment 交易，否則 CHECK constraint 還原會失敗
DELETE FROM transactions WHERE transaction_type = 'adjustment';

-- 移除對帳稽核欄位
ALTER TABLE transactions
    DROP COLUMN adjustment_reason,
    DROP COLUMN adjustment_prev_avg_cost,
    DROP COLUMN adjustment_prev_quantity;

-- 還原 transaction_type CHECK constraint
ALTER TABLE transactions DROP CONSTRAINT transactions_transaction_type_check;
ALTER TABLE transactions ADD CONSTRAINT transactions_transaction_type_check
    CHECK (transaction_type IN ('buy', 'sell', 'dividend', 'fee'));
