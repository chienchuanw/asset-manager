package mcpserver

import (
	"context"
	"errors"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeTxWriter struct {
	created *models.Transaction
	err     error
	deleted bool
}

func (f *fakeTxWriter) CreateTransaction(*models.CreateTransactionInput) (*models.Transaction, error) {
	return f.created, f.err
}
func (f *fakeTxWriter) UpdateTransaction(uuid.UUID, *models.UpdateTransactionInput) (*models.Transaction, error) {
	return f.created, f.err
}
func (f *fakeTxWriter) DeleteTransaction(uuid.UUID) error { f.deleted = true; return f.err }
func (f *fakeTxWriter) GetTransaction(uuid.UUID) (*models.Transaction, error) {
	return f.created, f.err
}

func TestTxWriteAdapter_Create(t *testing.T) {
	a := NewTxWriteAdapter(&fakeTxWriter{created: &models.Transaction{Symbol: "2330"}})
	tx, err := a.Create(&models.CreateTransactionInput{Symbol: "2330"})
	require.NoError(t, err)
	assert.Equal(t, "2330", tx.Symbol)
}

func TestTxWriteAdapter_DeleteInvalidUUID(t *testing.T) {
	f := &fakeTxWriter{}
	a := NewTxWriteAdapter(f)
	err := a.Delete("not-a-uuid")
	assert.Error(t, err)
	assert.False(t, f.deleted, "service must not be called on bad uuid")
}

func TestTxWriteAdapter_CreateError(t *testing.T) {
	a := NewTxWriteAdapter(&fakeTxWriter{err: errors.New("fifo violation")})
	_, err := a.Create(&models.CreateTransactionInput{})
	assert.ErrorContains(t, err, "fifo")
}

type fakeReconciler struct {
	gotDryRun bool
	called    bool
}

func (f *fakeReconciler) ReconcileHoldings(_ context.Context, _ []models.HoldingReconcileItem, dryRun bool) (*models.HoldingReconcilePreview, error) {
	f.called = true
	f.gotDryRun = dryRun
	return &models.HoldingReconcilePreview{}, nil
}

func TestReconcileAdapter_PassesDryRunFlag(t *testing.T) {
	f := &fakeReconciler{}
	a := NewReconcileAdapter(f)
	_, err := a.Reconcile(nil, true)
	require.NoError(t, err)
	assert.True(t, f.called)
	assert.True(t, f.gotDryRun)
}
