-- 建立已實現損益表
CREATE TABLE IF NOT EXISTS realized_profits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    asset_type VARCHAR(20) NOT NULL CHECK (asset_type IN ('cash', 'tw-stock', 'us-stock', 'crypto')),
    sell_date DATE NOT NULL,
    quantity DECIMAL(20, 8) NOT NULL CHECK (quantity > 0),
    sell_price DECIMAL(20, 8) NOT NULL CHECK (sell_price >= 0),
    sell_amount DECIMAL(20, 8) NOT NULL CHECK (sell_amount >= 0),
    sell_fee DECIMAL(20, 8) NOT NULL DEFAULT 0 CHECK (sell_fee >= 0),
    cost_basis DECIMAL(20, 8) NOT NULL CHECK (cost_basis >= 0),
    realized_pl DECIMAL(20, 8) NOT NULL,
    realized_pl_pct DECIMAL(10, 4) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'TWD' CHECK (currency IN ('TWD', 'USD')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 建立索引以提升查詢效能
CREATE INDEX idx_realized_profits_symbol ON realized_profits(symbol);
CREATE INDEX idx_realized_profits_asset_type ON realized_profits(asset_type);
CREATE INDEX idx_realized_profits_sell_date ON realized_profits(sell_date DESC);
CREATE INDEX idx_realized_profits_transaction_id ON realized_profits(transaction_id);

-- 建立複合索引（用於時間範圍查詢）
CREATE INDEX idx_realized_profits_date_asset ON realized_profits(sell_date, asset_type);
CREATE INDEX idx_realized_profits_date_symbol ON realized_profits(sell_date, symbol);

-- 建立更新時間的觸發器
CREATE TRIGGER update_realized_profits_updated_at
    BEFORE UPDATE ON realized_profits
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 加入表格和欄位註解
COMMENT ON TABLE realized_profits IS '已實現損益記錄表 - 記錄每筆賣出交易的損益資訊';
COMMENT ON COLUMN realized_profits.id IS '主鍵 UUID';
COMMENT ON COLUMN realized_profits.transaction_id IS '關聯的賣出交易 ID';
COMMENT ON COLUMN realized_profits.symbol IS '標的代碼';
COMMENT ON COLUMN realized_profits.asset_type IS '資產類型 (cash, tw-stock, us-stock, crypto)';
COMMENT ON COLUMN realized_profits.sell_date IS '賣出日期';
COMMENT ON COLUMN realized_profits.quantity IS '賣出數量';
COMMENT ON COLUMN realized_profits.sell_price IS '賣出價格';
COMMENT ON COLUMN realized_profits.sell_amount IS '賣出金額';
COMMENT ON COLUMN realized_profits.sell_fee IS '賣出手續費';
COMMENT ON COLUMN realized_profits.cost_basis IS 'FIFO 計算的成本基礎（含買入手續費）';
COMMENT ON COLUMN realized_profits.realized_pl IS '已實現損益 = (賣出金額 - 賣出手續費) - 成本基礎';
COMMENT ON COLUMN realized_profits.realized_pl_pct IS '已實現損益百分比 = (已實現損益 / 成本基礎) × 100';
COMMENT ON COLUMN realized_profits.currency IS '幣別 (TWD, USD)';
COMMENT ON COLUMN realized_profits.created_at IS '建立時間';
COMMENT ON COLUMN realized_profits.updated_at IS '更新時間';

