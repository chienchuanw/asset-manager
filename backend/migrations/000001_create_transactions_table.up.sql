-- 建立交易記錄資料表
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL,
    asset_type VARCHAR(20) NOT NULL CHECK (asset_type IN ('cash', 'tw-stock', 'us-stock', 'crypto')),
    symbol VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('buy', 'sell', 'dividend', 'fee')),
    quantity DECIMAL(20, 8) NOT NULL DEFAULT 0,
    price DECIMAL(20, 8) NOT NULL DEFAULT 0,
    amount DECIMAL(20, 2) NOT NULL,
    fee DECIMAL(20, 2) DEFAULT 0,
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 建立索引以提升查詢效能
CREATE INDEX idx_transactions_date ON transactions(date DESC);
CREATE INDEX idx_transactions_asset_type ON transactions(asset_type);
CREATE INDEX idx_transactions_symbol ON transactions(symbol);
CREATE INDEX idx_transactions_type ON transactions(transaction_type);

-- 建立更新時間的觸發器
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

