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

// Begin inserts a 'pending' row for a mutating tool and returns its id.
// A failure here must abort the write (fail-closed), so the error is returned.
func (r *PostgresAuditRepository) Begin(tool string, args []byte) (string, error) {
	if len(args) == 0 {
		args = []byte("{}")
	}
	var id string
	err := r.db.QueryRow(
		`INSERT INTO agent_audit_log (tool, arguments, status, duration_ms)
		 VALUES ($1, $2, 'pending', 0) RETURNING id`,
		tool, args,
	).Scan(&id)
	return id, err
}

// Finish updates a previously-Begun row with the outcome.
func (r *PostgresAuditRepository) Finish(id, status, errMsg string, summary []byte, durationMS int) error {
	var errText any
	if errMsg != "" {
		errText = errMsg
	}
	var sum any
	if len(summary) > 0 {
		sum = summary
	}
	_, err := r.db.Exec(
		`UPDATE agent_audit_log
		 SET status=$2, error=$3, result_summary=$4, duration_ms=$5, executed_at=now()
		 WHERE id=$1`,
		id, status, errText, sum, durationMS,
	)
	return err
}
