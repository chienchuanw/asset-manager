-- 刪除觸發器
DROP TRIGGER IF EXISTS trigger_update_settings_updated_at ON settings;

-- 刪除函式
DROP FUNCTION IF EXISTS update_settings_updated_at();

-- 刪除索引
DROP INDEX IF EXISTS idx_settings_key;

-- 刪除表
DROP TABLE IF EXISTS settings;

