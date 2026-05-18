package mcpserver

import (
	"encoding/json"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/audit"
	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// create_transaction with confirm:true, real Postgres audit repo, fake tx
// service: asserts the fail-closed lifecycle wrote a pending row that was
// finalized to success with executed_at set.
func TestCreateTransaction_ConfirmWritesAuditRowToDB(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	_, err := db.Exec("DELETE FROM agent_audit_log")
	require.NoError(t, err)

	w := &recordingTxWriter{tx: &models.Transaction{ID: uuid.New(), Symbol: "2330"}}
	fc := audit.NewFailClosed(audit.NewPostgresAuditRepository(db))
	tool := NewCreateTransactionTool(NewTxWriteAdapter(w), fc)

	_, _, err = tool.Run(json.RawMessage(
		`{"confirm":true,"date":"2026-05-18T00:00:00Z","asset_type":"tw-stock","symbol":"2330","name":"TSMC","type":"buy","quantity":1,"price":1,"amount":1,"currency":"TWD"}`))
	require.NoError(t, err)
	assert.True(t, w.createCalled)

	var tool_, status string
	var executedAtValid bool
	require.NoError(t, db.QueryRow(
		`SELECT tool, status, executed_at IS NOT NULL FROM agent_audit_log WHERE tool='create_transaction' LIMIT 1`,
	).Scan(&tool_, &status, &executedAtValid))
	assert.Equal(t, "create_transaction", tool_)
	assert.Equal(t, "success", status)
	assert.True(t, executedAtValid, "executed_at set after commit")
}

// Dry run must leave the audit table empty (no mutation, no audit row).
func TestCreateTransaction_DryRunWritesNoAuditRow(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	_, err := db.Exec("DELETE FROM agent_audit_log")
	require.NoError(t, err)

	w := &recordingTxWriter{tx: &models.Transaction{ID: uuid.New()}}
	fc := audit.NewFailClosed(audit.NewPostgresAuditRepository(db))
	tool := NewCreateTransactionTool(NewTxWriteAdapter(w), fc)

	_, _, err = tool.Run(json.RawMessage(
		`{"date":"2026-05-18T00:00:00Z","asset_type":"tw-stock","symbol":"2330","name":"TSMC","type":"buy","quantity":1,"price":1,"amount":1,"currency":"TWD"}`))
	require.NoError(t, err)

	var n int
	require.NoError(t, db.QueryRow("SELECT count(*) FROM agent_audit_log").Scan(&n))
	assert.Equal(t, 0, n, "dry run must not write an audit row")
	assert.False(t, w.createCalled)
}
