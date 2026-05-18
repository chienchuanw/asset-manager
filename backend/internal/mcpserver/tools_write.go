package mcpserver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/audit"
	"github.com/chienchuanw/asset-manager/internal/models"
)

// previewSummary marks an audit-free dry-run result.
func previewSummary() json.RawMessage { return json.RawMessage(`{"mode":"preview"}`) }

// commitAudited runs a mutation under the fail-closed audit lifecycle:
// Begin a pending row first — if that fails the mutation is aborted and
// nothing is persisted — then execute, then Finish with the outcome.
func commitAudited(
	fc *audit.FailClosed,
	tool string,
	raw json.RawMessage,
	exec func() (payload json.RawMessage, summary json.RawMessage, err error),
) (json.RawMessage, json.RawMessage, error) {
	start := time.Now()
	id, err := fc.Begin(tool, []byte(raw))
	if err != nil {
		return nil, nil, fmt.Errorf("audit unavailable, write aborted: %w", err)
	}
	payload, summary, execErr := exec()
	dur := int(time.Since(start).Milliseconds())
	if execErr != nil {
		fc.Finish(id, "error", execErr.Error(), nil, dur)
		return nil, nil, execErr
	}
	fc.Finish(id, "success", "", summary, dur)
	return payload, summary, nil
}

// --- create_transaction ---

type createTxArgs struct {
	models.CreateTransactionInput
	Confirm bool `json:"confirm"`
}

type createTransactionTool struct {
	port *TxWriteAdapter
	fc   *audit.FailClosed
}

func NewCreateTransactionTool(p *TxWriteAdapter, fc *audit.FailClosed) Tool {
	return &createTransactionTool{p, fc}
}
func (t *createTransactionTool) Name() string  { return "create_transaction" }
func (t *createTransactionTool) selfAudited()  {}

func (t *createTransactionTool) Run(raw json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	var a createTxArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return nil, nil, fmt.Errorf("invalid arguments: %w", err)
	}
	if a.Symbol == "" || a.AssetType == "" || a.TransactionType == "" || a.Currency == "" {
		return nil, nil, fmt.Errorf("symbol, asset_type, type and currency are required")
	}
	if a.Date.IsZero() {
		return nil, nil, fmt.Errorf("date is required (RFC3339)")
	}
	in := a.CreateTransactionInput
	if !a.Confirm {
		p, err := json.Marshal(map[string]any{
			"would_create": in,
			"note":         "dry run — pass confirm:true to persist",
		})
		if err != nil {
			return nil, nil, err
		}
		return p, previewSummary(), nil
	}
	return commitAudited(t.fc, t.Name(), raw, func() (json.RawMessage, json.RawMessage, error) {
		tx, err := t.port.Create(&in)
		if err != nil {
			return nil, nil, err
		}
		p, mErr := json.Marshal(tx)
		if mErr != nil {
			return nil, nil, mErr
		}
		return p, json.RawMessage(fmt.Sprintf(`{"created":%q}`, tx.ID.String())), nil
	})
}

// --- update_transaction ---

type updateTxArgs struct {
	ID      string `json:"id"`
	Confirm bool   `json:"confirm"`
	models.UpdateTransactionInput
}

type updateTransactionTool struct {
	port *TxWriteAdapter
	fc   *audit.FailClosed
}

func NewUpdateTransactionTool(p *TxWriteAdapter, fc *audit.FailClosed) Tool {
	return &updateTransactionTool{p, fc}
}
func (t *updateTransactionTool) Name() string { return "update_transaction" }
func (t *updateTransactionTool) selfAudited() {}

func (t *updateTransactionTool) Run(raw json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	var a updateTxArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return nil, nil, fmt.Errorf("invalid arguments: %w", err)
	}
	if a.ID == "" {
		return nil, nil, fmt.Errorf("id is required")
	}
	in := a.UpdateTransactionInput
	if !a.Confirm {
		current, err := t.port.Get(a.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("transaction %q not found: %w", a.ID, err)
		}
		p, mErr := json.Marshal(map[string]any{
			"current":          current,
			"proposed_changes": in,
			"note":             "dry run — pass confirm:true to apply",
		})
		if mErr != nil {
			return nil, nil, mErr
		}
		return p, previewSummary(), nil
	}
	return commitAudited(t.fc, t.Name(), raw, func() (json.RawMessage, json.RawMessage, error) {
		tx, err := t.port.Update(a.ID, &in)
		if err != nil {
			return nil, nil, err
		}
		p, mErr := json.Marshal(tx)
		if mErr != nil {
			return nil, nil, mErr
		}
		return p, json.RawMessage(fmt.Sprintf(`{"updated":%q}`, a.ID)), nil
	})
}

// --- delete_transaction ---

type deleteTxArgs struct {
	ID      string `json:"id"`
	Confirm bool   `json:"confirm"`
}

type deleteTransactionTool struct {
	port *TxWriteAdapter
	fc   *audit.FailClosed
}

func NewDeleteTransactionTool(p *TxWriteAdapter, fc *audit.FailClosed) Tool {
	return &deleteTransactionTool{p, fc}
}
func (t *deleteTransactionTool) Name() string { return "delete_transaction" }
func (t *deleteTransactionTool) selfAudited() {}

func (t *deleteTransactionTool) Run(raw json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	var a deleteTxArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return nil, nil, fmt.Errorf("invalid arguments: %w", err)
	}
	if a.ID == "" {
		return nil, nil, fmt.Errorf("id is required")
	}
	if !a.Confirm {
		current, err := t.port.Get(a.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("transaction %q not found: %w", a.ID, err)
		}
		p, mErr := json.Marshal(map[string]any{
			"would_delete": current,
			"note":         "dry run — pass confirm:true to delete",
		})
		if mErr != nil {
			return nil, nil, mErr
		}
		return p, previewSummary(), nil
	}
	return commitAudited(t.fc, t.Name(), raw, func() (json.RawMessage, json.RawMessage, error) {
		if err := t.port.Delete(a.ID); err != nil {
			return nil, nil, err
		}
		return json.RawMessage(fmt.Sprintf(`{"deleted":%q}`, a.ID)),
			json.RawMessage(fmt.Sprintf(`{"deleted":%q}`, a.ID)), nil
	})
}

// --- reconcile_holding ---

type reconcileArgs struct {
	Items   []models.HoldingReconcileItem `json:"items"`
	Confirm bool                          `json:"confirm"`
}

type reconcileHoldingTool struct {
	port *ReconcileAdapter
	fc   *audit.FailClosed
}

func NewReconcileHoldingTool(p *ReconcileAdapter, fc *audit.FailClosed) Tool {
	return &reconcileHoldingTool{p, fc}
}
func (t *reconcileHoldingTool) Name() string { return "reconcile_holding" }
func (t *reconcileHoldingTool) selfAudited() {}

func (t *reconcileHoldingTool) Run(raw json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	var a reconcileArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return nil, nil, fmt.Errorf("invalid arguments: %w", err)
	}
	if len(a.Items) == 0 {
		return nil, nil, fmt.Errorf("items is required and must be non-empty")
	}
	if !a.Confirm {
		// The service's own dryRun computes the diff WITHOUT writing.
		prev, err := t.port.Reconcile(a.Items, true)
		if err != nil {
			return nil, nil, err
		}
		p, mErr := json.Marshal(map[string]any{
			"preview": prev,
			"note":    "dry run — pass confirm:true to write adjustments",
		})
		if mErr != nil {
			return nil, nil, mErr
		}
		return p, previewSummary(), nil
	}
	return commitAudited(t.fc, t.Name(), raw, func() (json.RawMessage, json.RawMessage, error) {
		prev, err := t.port.Reconcile(a.Items, false)
		if err != nil {
			return nil, nil, err
		}
		p, mErr := json.Marshal(prev)
		if mErr != nil {
			return nil, nil, mErr
		}
		return p, json.RawMessage(fmt.Sprintf(`{"reconciled_items":%d}`, len(a.Items))), nil
	})
}
