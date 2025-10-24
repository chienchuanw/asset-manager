-- 建立資產快照表
CREATE TABLE IF NOT EXISTS asset_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snapshot_date DATE NOT NULL,
    asset_type VARCHAR(20) NOT NULL CHECK (asset_type IN ('tw-stock', 'us-stock', 'crypto', 'total')),
    value_twd DECIMAL(20, 2) NOT NULL CHECK (value_twd >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(snapshot_date, asset_type)
);

-- 建立索引以提升查詢效能
CREATE INDEX idx_asset_snapshots_date ON asset_snapshots(snapshot_date DESC);
CREATE INDEX idx_asset_snapshots_asset_type ON asset_snapshots(asset_type);
CREATE INDEX idx_asset_snapshots_date_type ON asset_snapshots(snapshot_date, asset_type);

-- 建立更新時間的觸發器
CREATE TRIGGER update_asset_snapshots_updated_at
    BEFORE UPDATE ON asset_snapshots
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 加入表格和欄位註解
COMMENT ON TABLE asset_snapshots IS '資產快照表 - 記錄每日各類資產的總價值';
COMMENT ON COLUMN asset_snapshots.id IS '主鍵 UUID';
COMMENT ON COLUMN asset_snapshots.snapshot_date IS '快照日期';
COMMENT ON COLUMN asset_snapshots.asset_type IS '資產類型 (tw-stock, us-stock, crypto, total)';
COMMENT ON COLUMN asset_snapshots.value_twd IS '資產價值（新台幣）';
COMMENT ON COLUMN asset_snapshots.created_at IS '建立時間';
COMMENT ON COLUMN asset_snapshots.updated_at IS '更新時間';

