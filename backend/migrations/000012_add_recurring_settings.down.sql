-- 刪除訂閱分期相關的通知設定
DELETE FROM settings WHERE key IN (
    'notification_daily_billing',
    'notification_subscription_expiry',
    'notification_installment_completion',
    'notification_expiry_days'
);

