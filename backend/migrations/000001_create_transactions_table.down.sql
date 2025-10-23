-- 刪除觸發器
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;

-- 刪除函式
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 刪除資料表
DROP TABLE IF EXISTS transactions;

