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
