-- 回滾：移除轉帳相關的預設分類

-- 移除存入類型的移轉分類
DELETE FROM cash_flow_categories 
WHERE name = '移轉' AND type = 'transfer_in';

-- 移除轉出類型的移轉分類
DELETE FROM cash_flow_categories 
WHERE name = '移轉' AND type = 'transfer_out';
