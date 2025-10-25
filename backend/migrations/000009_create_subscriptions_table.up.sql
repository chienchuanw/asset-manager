-- 建立訂閱表
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(20, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'TWD' CHECK (currency = 'TWD'),
    billing_cycle VARCHAR(20) NOT NULL CHECK (billing_cycle IN ('monthly', 'quarterly', 'yearly')),
    billing_day INT NOT NULL CHECK (billing_day BETWEEN 1 AND 31),
    category_id UUID NOT NULL REFERENCES cash_flow_categories(id) ON DELETE RESTRICT,
    start_date DATE NOT NULL,
    end_date DATE,
    auto_renew BOOLEAN NOT NULL DEFAULT true,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'cancelled')),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 建立索引以提升查詢效能
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_billing_day ON subscriptions(billing_day);
CREATE INDEX idx_subscriptions_category_id ON subscriptions(category_id);
CREATE INDEX idx_subscriptions_start_date ON subscriptions(start_date);
CREATE INDEX idx_subscriptions_end_date ON subscriptions(end_date);

-- 建立更新時間的觸發器
CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 註解說明
COMMENT ON TABLE subscriptions IS '訂閱表';
COMMENT ON COLUMN subscriptions.name IS '訂閱名稱';
COMMENT ON COLUMN subscriptions.amount IS '訂閱金額';
COMMENT ON COLUMN subscriptions.currency IS '幣別（目前僅支援 TWD）';
COMMENT ON COLUMN subscriptions.billing_cycle IS '計費週期（monthly: 月繳, quarterly: 季繳, yearly: 年繳）';
COMMENT ON COLUMN subscriptions.billing_day IS '扣款日（1-31）';
COMMENT ON COLUMN subscriptions.category_id IS '分類 ID';
COMMENT ON COLUMN subscriptions.start_date IS '開始日期';
COMMENT ON COLUMN subscriptions.end_date IS '結束日期（可選）';
COMMENT ON COLUMN subscriptions.auto_renew IS '自動續約';
COMMENT ON COLUMN subscriptions.status IS '狀態（active: 進行中, cancelled: 已取消）';
COMMENT ON COLUMN subscriptions.note IS '備註';

