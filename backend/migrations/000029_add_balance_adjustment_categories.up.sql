-- 新增「餘額調整」系統分類
-- 用於銀行帳戶 / 信用卡校準時自動產生的對帳 cash flow
INSERT INTO cash_flow_categories (name, type, is_system, sort_order)
VALUES
    ('餘額調整-收入', 'income', true, 9000),
    ('餘額調整-支出', 'expense', true, 9000)
ON CONFLICT (name, type) DO NOTHING;
