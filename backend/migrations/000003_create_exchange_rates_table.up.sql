-- 建立匯率表
CREATE TABLE IF NOT EXISTS exchange_rates (
    id SERIAL PRIMARY KEY,
    from_currency VARCHAR(3) NOT NULL CHECK (from_currency IN ('TWD', 'USD')),
    to_currency VARCHAR(3) NOT NULL CHECK (to_currency IN ('TWD', 'USD')),
    rate DECIMAL(20, 6) NOT NULL CHECK (rate > 0),
    date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(from_currency, to_currency, date)
);

-- 建立索引
CREATE INDEX idx_exchange_rates_date ON exchange_rates(date);
CREATE INDEX idx_exchange_rates_currencies ON exchange_rates(from_currency, to_currency);

-- 插入初始資料（USD to TWD 的預設匯率）
INSERT INTO exchange_rates (from_currency, to_currency, rate, date)
VALUES ('USD', 'TWD', 31.5, CURRENT_DATE)
ON CONFLICT (from_currency, to_currency, date) DO NOTHING;

