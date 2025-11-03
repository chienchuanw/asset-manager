-- 建立排程任務執行記錄表
CREATE TABLE IF NOT EXISTS scheduler_logs (
    id SERIAL PRIMARY KEY,
    task_name VARCHAR(100) NOT NULL,           -- 任務名稱（例如：daily_snapshot, discord_report）
    status VARCHAR(20) NOT NULL,               -- 執行狀態：success, failed, running
    error_message TEXT,                        -- 錯誤訊息（如果失敗）
    started_at TIMESTAMP NOT NULL,             -- 開始時間
    completed_at TIMESTAMP,                    -- 完成時間
    duration_seconds DECIMAL(10, 3),           -- 執行時長（秒）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 建立索引以加速查詢
CREATE INDEX idx_scheduler_logs_task_name ON scheduler_logs(task_name);
CREATE INDEX idx_scheduler_logs_started_at ON scheduler_logs(started_at DESC);
CREATE INDEX idx_scheduler_logs_status ON scheduler_logs(status);

