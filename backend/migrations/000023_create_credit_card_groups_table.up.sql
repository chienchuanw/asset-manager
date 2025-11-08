-- 建立信用卡群組表
CREATE TABLE IF NOT EXISTS credit_card_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    issuing_bank VARCHAR(255) NOT NULL,
    shared_credit_limit DECIMAL(20, 2) NOT NULL CHECK (shared_credit_limit > 0),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 為 credit_cards 表新增 group_id 欄位
ALTER TABLE credit_cards 
ADD COLUMN group_id UUID REFERENCES credit_card_groups(id) ON DELETE SET NULL;

-- 建立索引以提升查詢效能
CREATE INDEX idx_credit_cards_group_id ON credit_cards(group_id);

-- 新增註解說明
COMMENT ON TABLE credit_card_groups IS '信用卡群組表,用於管理共享額度的信用卡群組';
COMMENT ON COLUMN credit_card_groups.name IS '群組名稱';
COMMENT ON COLUMN credit_card_groups.issuing_bank IS '發卡銀行(群組內所有卡片必須來自同一銀行)';
COMMENT ON COLUMN credit_card_groups.shared_credit_limit IS '共享信用額度';
COMMENT ON COLUMN credit_cards.group_id IS '所屬群組 ID,若為 NULL 則為獨立卡片';

