package mcpserver

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/chienchuanw/asset-manager/internal/audit"
	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type auditSpy struct {
	entries []spyEntry
}
type spyEntry struct {
	tool, status string
}

func (s *auditSpy) Record(tool, status string, _, _ json.RawMessage, _ int, _ string) {
	s.entries = append(s.entries, spyEntry{tool, status})
}

func TestDispatch_AuditsSuccessAndError(t *testing.T) {
	spy := &auditSpy{}
	d := NewDispatcher(
		[]Tool{NewGetHoldingsTool(fakeHoldingsPort{list: []*models.Holding{{Symbol: "X"}}})},
		spy,
	)

	payload, err := d.Dispatch("get_holdings", json.RawMessage(`{}`))
	require.NoError(t, err)
	assert.Contains(t, string(payload), "X")

	_, err = d.Dispatch("does_not_exist", json.RawMessage(`{}`))
	assert.Error(t, err)

	require.Len(t, spy.entries, 2)
	assert.Equal(t, spyEntry{"get_holdings", "success"}, spy.entries[0])
	assert.Equal(t, spyEntry{"does_not_exist", "error"}, spy.entries[1])
}

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	get := func(k, d string) string {
		if v := os.Getenv(k); v != "" {
			return v
		}
		return d
	}
	dsn := "host=" + get("TEST_DB_HOST", "localhost") +
		" port=" + get("TEST_DB_PORT", "5432") +
		" user=" + get("TEST_DB_USER", "postgres") +
		" password=" + get("TEST_DB_PASSWORD", "postgres") +
		" dbname=" + get("TEST_DB_NAME", "asset_manager_test") +
		" sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	return db
}

// End-to-end: real audit repository, real dispatcher, fake adapter data.
func TestDispatch_WritesAuditRowToDB(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	_, err := db.Exec("DELETE FROM agent_audit_log")
	require.NoError(t, err)

	sink := NewLoggerSink(audit.NewLogger(audit.NewPostgresAuditRepository(db)))
	d := NewDispatcher(
		[]Tool{NewGetHoldingsTool(fakeHoldingsPort{list: []*models.Holding{{Symbol: "2330"}}})},
		sink,
	)

	payload, err := d.Dispatch("get_holdings", json.RawMessage(`{}`))
	require.NoError(t, err)
	assert.Contains(t, string(payload), "2330")

	var tool, status string
	row := db.QueryRow("SELECT tool, status FROM agent_audit_log LIMIT 1")
	require.NoError(t, row.Scan(&tool, &status))
	assert.Equal(t, "get_holdings", tool)
	assert.Equal(t, "success", status)
}
