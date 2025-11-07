-- 插入測試資料庫的系統預設分類
-- 使用 ON CONFLICT DO NOTHING 避免重複插入錯誤

-- 收入分類
INSERT INTO cash_flow_categories (name, type, is_system) VALUES
    ('薪資', 'income', true),
    ('獎金', 'income', true),
    ('利息', 'income', true),
    ('其他收入', 'income', true)
ON CONFLICT (name, type) DO NOTHING;

-- 支出分類
INSERT INTO cash_flow_categories (name, type, is_system) VALUES
    ('飲食', 'expense', true),
    ('交通', 'expense', true),
    ('娛樂', 'expense', true),
    ('醫療', 'expense', true),
    ('房租', 'expense', true),
    ('水電', 'expense', true),
    ('保險', 'expense', true),
    ('其他支出', 'expense', true)
ON CONFLICT (name, type) DO NOTHING;

-- 轉帳分類
INSERT INTO cash_flow_categories (name, type, is_system) VALUES
    ('移轉', 'transfer_in', true),
    ('移轉', 'transfer_out', true)
ON CONFLICT (name, type) DO NOTHING;

