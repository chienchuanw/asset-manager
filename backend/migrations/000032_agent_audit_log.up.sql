CREATE TABLE agent_audit_log (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tool           TEXT NOT NULL,
    arguments      JSONB NOT NULL DEFAULT '{}'::jsonb,
    status         TEXT NOT NULL CHECK (status IN ('success', 'error')),
    error          TEXT,
    result_summary JSONB,
    duration_ms    INTEGER NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_agent_audit_log_created_at ON agent_audit_log (created_at DESC);
CREATE INDEX idx_agent_audit_log_tool ON agent_audit_log (tool);
