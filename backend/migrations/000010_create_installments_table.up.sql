-- 建立分期表
CREATE TABLE IF NOT EXISTS installments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    total_amount DECIMAL(20, 2) NOT NULL CHECK (total_amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'TWD' CHECK (currency = 'TWD'),
    installment_count INT NOT NULL CHECK (installment_count > 0),
    installment_amount DECIMAL(20, 2) NOT NULL CHECK (installment_amount > 0),
    interest_rate DECIMAL(5, 2) NOT NULL DEFAULT 0 CHECK (interest_rate >= 0),
    total_interest DECIMAL(20, 2) NOT NULL DEFAULT 0 CHECK (total_interest >= 0),
    paid_count INT NOT NULL DEFAULT 0 CHECK (paid_count >= 0),
    billing_day INT NOT NULL CHECK (billing_day BETWEEN 1 AND 31),
    category_id UUID NOT NULL REFERENCES cash_flow_categories(id) ON DELETE RESTRICT,
    start_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'cancelled')),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_paid_count_not_exceed_total CHECK (paid_count <= installment_count)
);

-- 建立索引以提升查詢效能
CREATE INDEX idx_installments_status ON installments(status);
CREATE INDEX idx_installments_billing_day ON installments(billing_day);
CREATE INDEX idx_installments_category_id ON installments(category_id);
CREATE INDEX idx_installments_start_date ON installments(start_date);

-- 建立更新時間的觸發器
CREATE TRIGGER update_installments_updated_at
    BEFORE UPDATE ON installments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 註解說明
COMMENT ON TABLE installments IS '分期表';
COMMENT ON COLUMN installments.name IS '商品名稱';
COMMENT ON COLUMN installments.total_amount IS '總金額（本金）';
COMMENT ON COLUMN installments.currency IS '幣別（目前僅支援 TWD）';
COMMENT ON COLUMN installments.installment_count IS '總期數';
COMMENT ON COLUMN installments.installment_amount IS '每期金額';
COMMENT ON COLUMN installments.interest_rate IS '利率（%）';
COMMENT ON COLUMN installments.total_interest IS '總利息';
COMMENT ON COLUMN installments.paid_count IS '已付期數';
COMMENT ON COLUMN installments.billing_day IS '扣款日（1-31）';
COMMENT ON COLUMN installments.category_id IS '分類 ID';
COMMENT ON COLUMN installments.start_date IS '開始日期';
COMMENT ON COLUMN installments.status IS '狀態（active: 進行中, completed: 已完成, cancelled: 已取消）';
COMMENT ON COLUMN installments.note IS '備註';

