-- 回滾「餘額調整」系統分類
DELETE FROM cash_flow_categories
WHERE is_system = true
  AND name IN ('餘額調整-收入', '餘額調整-支出');
