package audit

import (
	"encoding/json"
	"log"
)

type Entry struct {
	Tool          string
	Arguments     json.RawMessage
	Status        string // "success" | "error"
	Error         string
	ResultSummary json.RawMessage
	DurationMS    int
}

type Repository interface {
	Record(e Entry) error
}

// Logger is the best-effort façade used by the server. A storage failure is
// logged to stderr and swallowed so a logging hiccup never blocks a read
// (Phase 1 policy; fail-closed is deferred to Phase 2).
type Logger struct {
	repo Repository
}

func NewLogger(repo Repository) *Logger {
	return &Logger{repo: repo}
}

func (l *Logger) Record(e Entry) {
	if err := l.repo.Record(e); err != nil {
		log.Printf("audit: failed to record %s: %v", e.Tool, err)
	}
}

// WriteAuditRepo is the fail-closed lifecycle for mutating tools.
type WriteAuditRepo interface {
	Begin(tool string, args []byte) (id string, err error)
	Finish(id, status, errMsg string, summary []byte, durationMS int) error
}

// FailClosed enforces audit-before-mutate: a write tool calls Begin and must
// abort if it errors (nothing is mutated without an audit row). Finish records
// the outcome; a Finish failure cannot un-do a completed write, so it is
// logged but not fatal — only Begin gates the operation.
type FailClosed struct {
	repo WriteAuditRepo
}

func NewFailClosed(repo WriteAuditRepo) *FailClosed {
	return &FailClosed{repo: repo}
}

func (f *FailClosed) Begin(tool string, args []byte) (string, error) {
	return f.repo.Begin(tool, args)
}

func (f *FailClosed) Finish(id, status, errMsg string, summary []byte, durationMS int) {
	if err := f.repo.Finish(id, status, errMsg, summary, durationMS); err != nil {
		log.Printf("audit: failed to finalize %s: %v", id, err)
	}
}
