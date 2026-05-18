package audit

import "database/sql"

type PostgresAuditRepository struct {
	db *sql.DB
}

func NewPostgresAuditRepository(db *sql.DB) *PostgresAuditRepository {
	return &PostgresAuditRepository{db: db}
}

func (r *PostgresAuditRepository) Record(e Entry) error {
	args := e.Arguments
	if len(args) == 0 {
		args = []byte("{}")
	}
	var errText any
	if e.Error != "" {
		errText = e.Error
	}
	var summary any
	if len(e.ResultSummary) > 0 {
		summary = []byte(e.ResultSummary)
	}
	_, err := r.db.Exec(
		`INSERT INTO agent_audit_log (tool, arguments, status, error, result_summary, duration_ms)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		e.Tool, []byte(args), e.Status, errText, summary, e.DurationMS,
	)
	return err
}
