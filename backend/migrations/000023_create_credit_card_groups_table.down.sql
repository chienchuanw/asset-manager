-- 移除索引
DROP INDEX IF EXISTS idx_credit_cards_group_id;

-- 移除 credit_cards 表的 group_id 欄位
ALTER TABLE credit_cards 
DROP COLUMN IF EXISTS group_id;

-- 刪除信用卡群組表
DROP TABLE IF EXISTS credit_card_groups;

