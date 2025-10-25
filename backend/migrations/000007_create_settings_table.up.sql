-- 建立 settings 表
CREATE TABLE IF NOT EXISTS settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(255) NOT NULL UNIQUE,
    value TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 建立索引
CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(key);

-- 建立 updated_at 自動更新觸發器
CREATE OR REPLACE FUNCTION update_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_settings_updated_at
    BEFORE UPDATE ON settings
    FOR EACH ROW
    EXECUTE FUNCTION update_settings_updated_at();

-- 插入預設設定
INSERT INTO settings (key, value, description) VALUES
    ('discord_webhook_url', '', 'Discord Webhook URL for daily reports'),
    ('discord_enabled', 'false', 'Enable Discord daily reports'),
    ('discord_report_time', '09:00', 'Daily report time (HH:MM)'),
    ('target_allocation_tw_stock', '40', 'Target allocation percentage for Taiwan stocks'),
    ('target_allocation_us_stock', '40', 'Target allocation percentage for US stocks'),
    ('target_allocation_crypto', '20', 'Target allocation percentage for cryptocurrencies'),
    ('rebalance_threshold', '5', 'Rebalancing threshold percentage')
ON CONFLICT (key) DO NOTHING;

