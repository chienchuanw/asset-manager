-- 建立現金流報告記錄表
CREATE TABLE cash_flow_report_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_type VARCHAR(20) NOT NULL CHECK (report_type IN ('monthly', 'yearly')),
    year INTEGER NOT NULL,
    month INTEGER,
    sent_at TIMESTAMP NOT NULL,
    success BOOLEAN NOT NULL DEFAULT false,
    error_msg TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 建立索引以提升查詢效能
CREATE INDEX idx_cash_flow_report_logs_type_year_month ON cash_flow_report_logs(report_type, year, month);
CREATE INDEX idx_cash_flow_report_logs_success ON cash_flow_report_logs(success);
CREATE INDEX idx_cash_flow_report_logs_created_at ON cash_flow_report_logs(created_at DESC);

-- 新增註解
COMMENT ON TABLE cash_flow_report_logs IS '現金流報告發送記錄表';
COMMENT ON COLUMN cash_flow_report_logs.report_type IS '報告類型：monthly（月度）或 yearly（年度）';
COMMENT ON COLUMN cash_flow_report_logs.year IS '報告年份';
COMMENT ON COLUMN cash_flow_report_logs.month IS '報告月份（年度報告時為 NULL）';
COMMENT ON COLUMN cash_flow_report_logs.sent_at IS '發送時間';
COMMENT ON COLUMN cash_flow_report_logs.success IS '是否發送成功';
COMMENT ON COLUMN cash_flow_report_logs.error_msg IS '錯誤訊息（如果發送失敗）';
COMMENT ON COLUMN cash_flow_report_logs.retry_count IS '重試次數';

