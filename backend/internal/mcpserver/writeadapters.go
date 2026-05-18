package mcpserver

import (
	"context"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
)

// --- transaction writes ---

type txWriter interface {
	CreateTransaction(*models.CreateTransactionInput) (*models.Transaction, error)
	UpdateTransaction(uuid.UUID, *models.UpdateTransactionInput) (*models.Transaction, error)
	DeleteTransaction(uuid.UUID) error
	GetTransaction(uuid.UUID) (*models.Transaction, error)
}

type TxWriteAdapter struct{ svc txWriter }

func NewTxWriteAdapter(svc txWriter) *TxWriteAdapter { return &TxWriteAdapter{svc: svc} }

func (a *TxWriteAdapter) Create(in *models.CreateTransactionInput) (*models.Transaction, error) {
	return a.svc.CreateTransaction(in)
}

func (a *TxWriteAdapter) Update(id string, in *models.UpdateTransactionInput) (*models.Transaction, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return a.svc.UpdateTransaction(uid, in)
}

func (a *TxWriteAdapter) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return a.svc.DeleteTransaction(uid)
}

// Get is used to build a delete/update preview (what would change).
func (a *TxWriteAdapter) Get(id string) (*models.Transaction, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return a.svc.GetTransaction(uid)
}

// --- holding reconciliation ---

type reconciler interface {
	ReconcileHoldings(context.Context, []models.HoldingReconcileItem, bool) (*models.HoldingReconcilePreview, error)
}

type ReconcileAdapter struct{ svc reconciler }

func NewReconcileAdapter(svc reconciler) *ReconcileAdapter { return &ReconcileAdapter{svc: svc} }

// Reconcile delegates to the service whose own dryRun flag returns a preview
// (true) or writes adjustments inside a single DB transaction (false).
func (a *ReconcileAdapter) Reconcile(items []models.HoldingReconcileItem, dryRun bool) (*models.HoldingReconcilePreview, error) {
	return a.svc.ReconcileHoldings(context.Background(), items, dryRun)
}
