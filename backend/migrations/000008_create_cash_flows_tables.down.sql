-- 刪除觸發器
DROP TRIGGER IF EXISTS update_cash_flows_updated_at ON cash_flows;
DROP TRIGGER IF EXISTS update_cash_flow_categories_updated_at ON cash_flow_categories;

-- 刪除資料表（注意順序，先刪除有外鍵的表）
DROP TABLE IF EXISTS cash_flows;
DROP TABLE IF EXISTS cash_flow_categories;

