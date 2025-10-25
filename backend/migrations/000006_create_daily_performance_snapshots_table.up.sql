-- 建立每日績效快照表
CREATE TABLE IF NOT EXISTS daily_performance_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snapshot_date DATE NOT NULL,
    
    -- 總體績效指標
    total_market_value DECIMAL(20, 2) NOT NULL DEFAULT 0,
    total_cost DECIMAL(20, 2) NOT NULL DEFAULT 0,
    total_unrealized_pl DECIMAL(20, 2) NOT NULL DEFAULT 0,
    total_unrealized_pct DECIMAL(10, 4) NOT NULL DEFAULT 0,
    total_realized_pl DECIMAL(20, 2) NOT NULL DEFAULT 0,
    total_realized_pct DECIMAL(10, 4) NOT NULL DEFAULT 0,
    
    -- 持倉統計
    holding_count INT NOT NULL DEFAULT 0,
    
    -- 幣別
    currency VARCHAR(3) NOT NULL DEFAULT 'TWD' CHECK (currency IN ('TWD', 'USD')),
    
    -- 時間戳記
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 確保每天只有一筆記錄
    UNIQUE(snapshot_date)
);

-- 建立索引
CREATE INDEX idx_daily_performance_snapshots_date ON daily_performance_snapshots(snapshot_date DESC);

-- 建立更新時間觸發器
CREATE OR REPLACE FUNCTION update_daily_performance_snapshots_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_daily_performance_snapshots_updated_at
    BEFORE UPDATE ON daily_performance_snapshots
    FOR EACH ROW
    EXECUTE FUNCTION update_daily_performance_snapshots_updated_at();

-- 建立每日績效快照明細表（按資產類型）
CREATE TABLE IF NOT EXISTS daily_performance_snapshot_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snapshot_id UUID NOT NULL REFERENCES daily_performance_snapshots(id) ON DELETE CASCADE,
    asset_type VARCHAR(20) NOT NULL CHECK (asset_type IN ('tw-stock', 'us-stock', 'crypto')),
    
    -- 績效指標
    market_value DECIMAL(20, 2) NOT NULL DEFAULT 0,
    cost DECIMAL(20, 2) NOT NULL DEFAULT 0,
    unrealized_pl DECIMAL(20, 2) NOT NULL DEFAULT 0,
    unrealized_pct DECIMAL(10, 4) NOT NULL DEFAULT 0,
    realized_pl DECIMAL(20, 2) NOT NULL DEFAULT 0,
    realized_pct DECIMAL(10, 4) NOT NULL DEFAULT 0,
    
    -- 持倉統計
    holding_count INT NOT NULL DEFAULT 0,
    
    -- 時間戳記
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 確保每個快照的每個資產類型只有一筆記錄
    UNIQUE(snapshot_id, asset_type)
);

-- 建立索引
CREATE INDEX idx_daily_performance_snapshot_details_snapshot_id ON daily_performance_snapshot_details(snapshot_id);
CREATE INDEX idx_daily_performance_snapshot_details_asset_type ON daily_performance_snapshot_details(asset_type);

-- 新增註解
COMMENT ON TABLE daily_performance_snapshots IS '每日績效快照表，記錄每天的總體績效指標';
COMMENT ON TABLE daily_performance_snapshot_details IS '每日績效快照明細表，記錄每天各資產類型的績效指標';
COMMENT ON COLUMN daily_performance_snapshots.snapshot_date IS '快照日期';
COMMENT ON COLUMN daily_performance_snapshots.total_market_value IS '總市值（TWD）';
COMMENT ON COLUMN daily_performance_snapshots.total_cost IS '總成本（TWD）';
COMMENT ON COLUMN daily_performance_snapshots.total_unrealized_pl IS '總未實現損益（TWD）';
COMMENT ON COLUMN daily_performance_snapshots.total_unrealized_pct IS '總未實現報酬率（%）';
COMMENT ON COLUMN daily_performance_snapshots.total_realized_pl IS '總已實現損益（TWD）';
COMMENT ON COLUMN daily_performance_snapshots.total_realized_pct IS '總已實現報酬率（%）';
COMMENT ON COLUMN daily_performance_snapshots.holding_count IS '持倉數量';

