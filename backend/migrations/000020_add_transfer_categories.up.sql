-- 新增轉帳相關的預設分類
-- 為 transfer_in 和 transfer_out 類型建立「移轉」分類

-- 新增存入類型的移轉分類
INSERT INTO cash_flow_categories (name, type) 
VALUES ('移轉', 'transfer_in')
ON CONFLICT (name, type) DO NOTHING;

-- 新增轉出類型的移轉分類
INSERT INTO cash_flow_categories (name, type) 
VALUES ('移轉', 'transfer_out')
ON CONFLICT (name, type) DO NOTHING;
