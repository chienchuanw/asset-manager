-- 刪除訂閱分期相關的通知設定
DELETE FROM settings WHERE key IN (
    'notification_daily_billing',
    'notification_subscription_expiry',
    'notification_installment_completion',
    'notification_expiry_days'
);

-- 刪除訂閱分期相關的支出分類
DELETE FROM cash_flow_categories WHERE name IN (
    '訂閱 - 娛樂',
    '訂閱 - 工具',
    '訂閱 - 學習',
    '訂閱 - 其他',
    '分期 - 3C產品',
    '分期 - 家電',
    '分期 - 其他'
) AND is_system = true;

