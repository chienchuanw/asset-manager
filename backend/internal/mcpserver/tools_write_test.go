package mcpserver

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/audit"
	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- fakes ---

type recordingTxWriter struct {
	createCalled, updateCalled, deleteCalled bool
	tx                                       *models.Transaction
	err                                      error
}

func (f *recordingTxWriter) CreateTransaction(*models.CreateTransactionInput) (*models.Transaction, error) {
	f.createCalled = true
	return f.tx, f.err
}
func (f *recordingTxWriter) UpdateTransaction(uuid.UUID, *models.UpdateTransactionInput) (*models.Transaction, error) {
	f.updateCalled = true
	return f.tx, f.err
}
func (f *recordingTxWriter) DeleteTransaction(uuid.UUID) error {
	f.deleteCalled = true
	return f.err
}
func (f *recordingTxWriter) GetTransaction(uuid.UUID) (*models.Transaction, error) {
	return f.tx, nil
}

type fakeWriteAuditRepo struct {
	beginErr    error
	begun       bool
	finished    bool
	finalStatus string
}

func (f *fakeWriteAuditRepo) Begin(string, []byte) (string, error) {
	if f.beginErr != nil {
		return "", f.beginErr
	}
	f.begun = true
	return "audit-id-1", nil
}
func (f *fakeWriteAuditRepo) Finish(_, status, _ string, _ []byte, _ int) error {
	f.finished = true
	f.finalStatus = status
	return nil
}

func newTx() *models.Transaction {
	return &models.Transaction{ID: uuid.New(), Symbol: "2330"}
}

// --- dry-run persists nothing ---

func TestCreateTransaction_DryRunDoesNotMutateOrAudit(t *testing.T) {
	w := &recordingTxWriter{tx: newTx()}
	ar := &fakeWriteAuditRepo{}
	tool := NewCreateTransactionTool(NewTxWriteAdapter(w), audit.NewFailClosed(ar))

	_, summary, err := tool.Run(json.RawMessage(
		`{"date":"2026-05-18T00:00:00Z","asset_type":"tw-stock","symbol":"2330","name":"TSMC","type":"buy","quantity":1,"price":1,"amount":1,"currency":"TWD"}`))

	require.NoError(t, err)
	assert.JSONEq(t, `{"mode":"preview"}`, string(summary))
	assert.False(t, w.createCalled, "service must not be called on dry run")
	assert.False(t, ar.begun, "no audit row on dry run")
}

// --- confirm path mutates + audits success ---

func TestCreateTransaction_ConfirmMutatesAndAudits(t *testing.T) {
	w := &recordingTxWriter{tx: newTx()}
	ar := &fakeWriteAuditRepo{}
	tool := NewCreateTransactionTool(NewTxWriteAdapter(w), audit.NewFailClosed(ar))

	_, _, err := tool.Run(json.RawMessage(
		`{"confirm":true,"date":"2026-05-18T00:00:00Z","asset_type":"tw-stock","symbol":"2330","name":"TSMC","type":"buy","quantity":1,"price":1,"amount":1,"currency":"TWD"}`))

	require.NoError(t, err)
	assert.True(t, w.createCalled)
	assert.True(t, ar.begun)
	assert.True(t, ar.finished)
	assert.Equal(t, "success", ar.finalStatus)
}

// --- fail-closed: Begin failure aborts before mutating ---

func TestCreateTransaction_AuditBeginFailureAbortsWrite(t *testing.T) {
	w := &recordingTxWriter{tx: newTx()}
	ar := &fakeWriteAuditRepo{beginErr: errors.New("db down")}
	tool := NewCreateTransactionTool(NewTxWriteAdapter(w), audit.NewFailClosed(ar))

	_, _, err := tool.Run(json.RawMessage(
		`{"confirm":true,"date":"2026-05-18T00:00:00Z","asset_type":"tw-stock","symbol":"2330","name":"TSMC","type":"buy","quantity":1,"price":1,"amount":1,"currency":"TWD"}`))

	assert.ErrorContains(t, err, "write aborted")
	assert.False(t, w.createCalled, "service must NOT be called when audit Begin fails")
}

// --- service error is audited as error ---

func TestCreateTransaction_ServiceErrorAuditedAsError(t *testing.T) {
	w := &recordingTxWriter{err: errors.New("fifo violation")}
	ar := &fakeWriteAuditRepo{}
	tool := NewCreateTransactionTool(NewTxWriteAdapter(w), audit.NewFailClosed(ar))

	_, _, err := tool.Run(json.RawMessage(
		`{"confirm":true,"date":"2026-05-18T00:00:00Z","asset_type":"tw-stock","symbol":"2330","name":"TSMC","type":"buy","quantity":1,"price":1,"amount":1,"currency":"TWD"}`))

	assert.ErrorContains(t, err, "fifo")
	assert.True(t, ar.finished)
	assert.Equal(t, "error", ar.finalStatus)
}

func TestDeleteTransaction_DryRunDoesNotDelete(t *testing.T) {
	w := &recordingTxWriter{tx: newTx()}
	ar := &fakeWriteAuditRepo{}
	tool := NewDeleteTransactionTool(NewTxWriteAdapter(w), audit.NewFailClosed(ar))

	id := w.tx.ID.String()
	_, summary, err := tool.Run(json.RawMessage(`{"id":"` + id + `"}`))

	require.NoError(t, err)
	assert.JSONEq(t, `{"mode":"preview"}`, string(summary))
	assert.False(t, w.deleteCalled)
}

func TestReconcileHolding_DryRunUsesServiceDryRun(t *testing.T) {
	rec := &fakeReconciler{}
	ar := &fakeWriteAuditRepo{}
	tool := NewReconcileHoldingTool(NewReconcileAdapter(rec), audit.NewFailClosed(ar))

	_, summary, err := tool.Run(json.RawMessage(
		`{"items":[{"symbol":"2330"}]}`))

	require.NoError(t, err)
	assert.JSONEq(t, `{"mode":"preview"}`, string(summary))
	assert.True(t, rec.called)
	assert.True(t, rec.gotDryRun, "preview must call service with dryRun=true")
	assert.False(t, ar.begun, "no fail-closed audit on preview")
}

func TestReconcileHolding_ConfirmWritesAndAudits(t *testing.T) {
	rec := &fakeReconciler{}
	ar := &fakeWriteAuditRepo{}
	tool := NewReconcileHoldingTool(NewReconcileAdapter(rec), audit.NewFailClosed(ar))

	_, _, err := tool.Run(json.RawMessage(
		`{"confirm":true,"items":[{"symbol":"2330"}]}`))

	require.NoError(t, err)
	assert.True(t, rec.called)
	assert.False(t, rec.gotDryRun, "commit must call service with dryRun=false")
	assert.Equal(t, "success", ar.finalStatus)
}
