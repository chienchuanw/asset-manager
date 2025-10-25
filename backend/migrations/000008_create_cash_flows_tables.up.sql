-- 建立現金流分類表
CREATE TABLE IF NOT EXISTS cash_flow_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense')),
    is_system BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, type)
);

-- 建立現金流記錄表
CREATE TABLE IF NOT EXISTS cash_flows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense')),
    category_id UUID NOT NULL REFERENCES cash_flow_categories(id) ON DELETE RESTRICT,
    amount DECIMAL(20, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'TWD' CHECK (currency = 'TWD'),
    description VARCHAR(500) NOT NULL,
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 建立索引以提升查詢效能
CREATE INDEX idx_cash_flows_date ON cash_flows(date DESC);
CREATE INDEX idx_cash_flows_type ON cash_flows(type);
CREATE INDEX idx_cash_flows_category_id ON cash_flows(category_id);
CREATE INDEX idx_cash_flow_categories_type ON cash_flow_categories(type);

-- 建立更新時間的觸發器
CREATE TRIGGER update_cash_flows_updated_at
    BEFORE UPDATE ON cash_flows
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cash_flow_categories_updated_at
    BEFORE UPDATE ON cash_flow_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 插入系統預設分類
-- 收入分類
INSERT INTO cash_flow_categories (name, type, is_system) VALUES
    ('薪資', 'income', true),
    ('獎金', 'income', true),
    ('利息', 'income', true),
    ('其他收入', 'income', true);

-- 支出分類
INSERT INTO cash_flow_categories (name, type, is_system) VALUES
    ('飲食', 'expense', true),
    ('交通', 'expense', true),
    ('娛樂', 'expense', true),
    ('醫療', 'expense', true),
    ('房租', 'expense', true),
    ('水電', 'expense', true),
    ('保險', 'expense', true),
    ('其他支出', 'expense', true);

-- 註解說明
COMMENT ON TABLE cash_flow_categories IS '現金流分類表';
COMMENT ON TABLE cash_flows IS '現金流記錄表';
COMMENT ON COLUMN cash_flow_categories.name IS '分類名稱';
COMMENT ON COLUMN cash_flow_categories.type IS '分類類型（income: 收入, expense: 支出）';
COMMENT ON COLUMN cash_flow_categories.is_system IS '是否為系統預設分類';
COMMENT ON COLUMN cash_flows.date IS '交易日期';
COMMENT ON COLUMN cash_flows.type IS '現金流類型（income: 收入, expense: 支出）';
COMMENT ON COLUMN cash_flows.category_id IS '分類 ID';
COMMENT ON COLUMN cash_flows.amount IS '金額（正數）';
COMMENT ON COLUMN cash_flows.currency IS '幣別（目前僅支援 TWD）';
COMMENT ON COLUMN cash_flows.description IS '描述';
COMMENT ON COLUMN cash_flows.note IS '備註';

