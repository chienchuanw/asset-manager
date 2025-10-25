-- 新增訂閱分期相關的通知設定
INSERT INTO settings (key, value, description) VALUES
    ('notification_daily_billing', 'true', 'Enable daily billing notification for subscriptions and installments'),
    ('notification_subscription_expiry', 'true', 'Enable subscription expiry notification'),
    ('notification_installment_completion', 'true', 'Enable installment completion notification'),
    ('notification_expiry_days', '7', 'Days before expiry to send notification')
ON CONFLICT (key) DO NOTHING;

