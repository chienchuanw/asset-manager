-- 新增訂閱分期相關的通知設定
INSERT INTO settings (key, value, description) VALUES
    ('notification_daily_billing', 'true', 'Enable daily billing notification for subscriptions and installments'),
    ('notification_subscription_expiry', 'true', 'Enable subscription expiry notification'),
    ('notification_installment_completion', 'true', 'Enable installment completion notification'),
    ('notification_expiry_days', '7', 'Days before expiry to send notification')
ON CONFLICT (key) DO NOTHING;

-- 新增訂閱分期相關的支出分類
INSERT INTO cash_flow_categories (name, type, is_system) VALUES
    ('訂閱 - 娛樂', 'expense', true),
    ('訂閱 - 工具', 'expense', true),
    ('訂閱 - 學習', 'expense', true),
    ('訂閱 - 其他', 'expense', true),
    ('分期 - 3C產品', 'expense', true),
    ('分期 - 家電', 'expense', true),
    ('分期 - 其他', 'expense', true)
ON CONFLICT (name, type) DO NOTHING;

-- 註解說明
COMMENT ON COLUMN settings.key IS '設定鍵值';
COMMENT ON COLUMN settings.value IS '設定值';
COMMENT ON COLUMN settings.description IS '設定說明';

