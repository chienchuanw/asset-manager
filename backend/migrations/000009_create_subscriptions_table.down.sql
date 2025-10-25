-- 刪除觸發器
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;

-- 刪除資料表
DROP TABLE IF EXISTS subscriptions;

