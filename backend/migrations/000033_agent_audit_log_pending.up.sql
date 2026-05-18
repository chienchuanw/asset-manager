ALTER TABLE agent_audit_log DROP CONSTRAINT IF EXISTS agent_audit_log_status_check;
ALTER TABLE agent_audit_log
    ADD CONSTRAINT agent_audit_log_status_check
    CHECK (status IN ('success', 'error', 'pending'));
ALTER TABLE agent_audit_log ADD COLUMN IF NOT EXISTS executed_at TIMESTAMPTZ;
