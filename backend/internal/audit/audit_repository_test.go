package audit

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestPostgresAuditRepository_Record(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	_, err := db.Exec("DELETE FROM agent_audit_log")
	require.NoError(t, err)

	repo := NewPostgresAuditRepository(db)
	err = repo.Record(Entry{
		Tool:          "get_holdings",
		Arguments:     json.RawMessage(`{"a":1}`),
		Status:        "success",
		ResultSummary: json.RawMessage(`{"rows":3}`),
		DurationMS:    12,
	})
	assert.NoError(t, err)

	var tool, status string
	var dur int
	row := db.QueryRow("SELECT tool, status, duration_ms FROM agent_audit_log LIMIT 1")
	require.NoError(t, row.Scan(&tool, &status, &dur))
	assert.Equal(t, "get_holdings", tool)
	assert.Equal(t, "success", status)
	assert.Equal(t, 12, dur)
}

func TestPostgresAuditRepository_RecordError(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	_, err := db.Exec("DELETE FROM agent_audit_log")
	require.NoError(t, err)

	repo := NewPostgresAuditRepository(db)
	err = repo.Record(Entry{
		Tool:       "get_holding",
		Status:     "error",
		Error:      "symbol not found",
		DurationMS: 4,
	})
	assert.NoError(t, err)

	var status, errText string
	row := db.QueryRow("SELECT status, error FROM agent_audit_log LIMIT 1")
	require.NoError(t, row.Scan(&status, &errText))
	assert.Equal(t, "error", status)
	assert.Equal(t, "symbol not found", errText)
}
