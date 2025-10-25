-- 刪除每日績效快照明細表
DROP TABLE IF EXISTS daily_performance_snapshot_details;

-- 刪除每日績效快照表
DROP TABLE IF EXISTS daily_performance_snapshots;

-- 刪除觸發器函式
DROP FUNCTION IF EXISTS update_daily_performance_snapshots_updated_at();

