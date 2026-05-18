package mcpserver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
)

// Tool is the SDK-independent contract: parse raw args, return the JSON
// payload and a lightweight summary for the audit log.
type Tool interface {
	Name() string
	Run(args json.RawMessage) (payload json.RawMessage, summary json.RawMessage, err error)
}

func okSummary() json.RawMessage          { return json.RawMessage(`{"ok":true}`) }
func rowsSummary(n int) json.RawMessage   { b, _ := json.Marshal(map[string]int{"rows": n}); return b }
func marshalPayload(v any) (json.RawMessage, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encode result: %w", err)
	}
	return b, nil
}

// --- get_holdings ---

type holdingsPort interface {
	List() ([]*models.Holding, error)
	Get(symbol string) (*models.Holding, error)
}

type getHoldingsTool struct{ port holdingsPort }

func NewGetHoldingsTool(p holdingsPort) Tool { return &getHoldingsTool{p} }
func (t *getHoldingsTool) Name() string      { return "get_holdings" }
func (t *getHoldingsTool) Run(json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	rows, err := t.port.List()
	if err != nil {
		return nil, nil, err
	}
	p, err := marshalPayload(rows)
	if err != nil {
		return nil, nil, err
	}
	return p, rowsSummary(len(rows)), nil
}

// --- get_holding ---

type getHoldingArgs struct {
	Symbol string `json:"symbol"`
}
type getHoldingTool struct{ port holdingsPort }

func NewGetHoldingTool(p holdingsPort) Tool { return &getHoldingTool{p} }
func (t *getHoldingTool) Name() string      { return "get_holding" }
func (t *getHoldingTool) Run(raw json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	var a getHoldingArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return nil, nil, fmt.Errorf("invalid arguments: %w", err)
	}
	if a.Symbol == "" {
		return nil, nil, fmt.Errorf("symbol is required")
	}
	h, err := t.port.Get(a.Symbol)
	if err != nil {
		return nil, nil, fmt.Errorf("holding %q not found: %w", a.Symbol, err)
	}
	p, err := marshalPayload(h)
	if err != nil {
		return nil, nil, err
	}
	return p, okSummary(), nil
}

// --- list_transactions ---

type txPort interface {
	List(repository.TransactionFilters) ([]*models.Transaction, error)
	Get(id string) (*models.Transaction, error)
}

type listTxArgs struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
	Limit  int    `json:"limit"`
}
type listTransactionsTool struct{ port txPort }

func NewListTransactionsTool(p txPort) Tool { return &listTransactionsTool{p} }
func (t *listTransactionsTool) Name() string { return "list_transactions" }
func (t *listTransactionsTool) Run(raw json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	var a listTxArgs
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &a); err != nil {
			return nil, nil, fmt.Errorf("invalid arguments: %w", err)
		}
	}
	var f repository.TransactionFilters
	if a.Symbol != "" {
		f.Symbol = &a.Symbol
	}
	if a.Type != "" {
		tt := models.TransactionType(a.Type)
		f.TransactionType = &tt
	}
	if a.From != "" {
		d, err := time.Parse("2006-01-02", a.From)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid 'from' date (want YYYY-MM-DD): %w", err)
		}
		f.StartDate = &d
	}
	if a.To != "" {
		d, err := time.Parse("2006-01-02", a.To)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid 'to' date (want YYYY-MM-DD): %w", err)
		}
		f.EndDate = &d
	}
	if a.Limit > 0 {
		f.Limit = a.Limit
	}
	rows, err := t.port.List(f)
	if err != nil {
		return nil, nil, err
	}
	p, err := marshalPayload(rows)
	if err != nil {
		return nil, nil, err
	}
	return p, rowsSummary(len(rows)), nil
}

// --- get_transaction ---

type getTxArgs struct {
	ID string `json:"id"`
}
type getTransactionTool struct{ port txPort }

func NewGetTransactionTool(p txPort) Tool { return &getTransactionTool{p} }
func (t *getTransactionTool) Name() string { return "get_transaction" }
func (t *getTransactionTool) Run(raw json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	var a getTxArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return nil, nil, fmt.Errorf("invalid arguments: %w", err)
	}
	if a.ID == "" {
		return nil, nil, fmt.Errorf("id is required")
	}
	tx, err := t.port.Get(a.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("transaction %q not found: %w", a.ID, err)
	}
	p, err := marshalPayload(tx)
	if err != nil {
		return nil, nil, err
	}
	return p, okSummary(), nil
}

// --- get_analytics ---

type analyticsPort interface {
	Get(models.TimeRange) (*AnalyticsResult, error)
}
type rangeArgs struct {
	Range string `json:"range"`
}
type getAnalyticsTool struct{ port analyticsPort }

func NewGetAnalyticsTool(p analyticsPort) Tool { return &getAnalyticsTool{p} }
func (t *getAnalyticsTool) Name() string       { return "get_analytics" }

func parseRange(v string) (models.TimeRange, error) {
	if v == "" {
		return models.TimeRangeAll, nil
	}
	tr := models.TimeRange(v)
	switch tr {
	case models.TimeRangeWeek, models.TimeRangeMonth, models.TimeRangeQuarter,
		models.TimeRangeYear, models.TimeRangeAll:
		return tr, nil
	}
	return "", fmt.Errorf("invalid range %q (want week|month|quarter|year|all)", v)
}

func (t *getAnalyticsTool) Run(raw json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	var a rangeArgs
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &a); err != nil {
			return nil, nil, fmt.Errorf("invalid arguments: %w", err)
		}
	}
	tr, err := parseRange(a.Range)
	if err != nil {
		return nil, nil, err
	}
	res, err := t.port.Get(tr)
	if err != nil {
		return nil, nil, err
	}
	p, err := marshalPayload(res)
	if err != nil {
		return nil, nil, err
	}
	return p, okSummary(), nil
}

// --- get_allocation ---

type allocPort interface {
	ByType() ([]models.AllocationByType, error)
}
type getAllocationTool struct{ port allocPort }

func NewGetAllocationTool(p allocPort) Tool { return &getAllocationTool{p} }
func (t *getAllocationTool) Name() string    { return "get_allocation" }
func (t *getAllocationTool) Run(json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	rows, err := t.port.ByType()
	if err != nil {
		return nil, nil, err
	}
	p, err := marshalPayload(rows)
	if err != nil {
		return nil, nil, err
	}
	return p, rowsSummary(len(rows)), nil
}

// --- get_performance_trend ---

type trendPort interface {
	Latest(days int) ([]models.PerformanceTrendPoint, error)
}
type trendArgs struct {
	Days int `json:"days"`
}
type getPerformanceTrendTool struct{ port trendPort }

func NewGetPerformanceTrendTool(p trendPort) Tool { return &getPerformanceTrendTool{p} }
func (t *getPerformanceTrendTool) Name() string    { return "get_performance_trend" }
func (t *getPerformanceTrendTool) Run(raw json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	a := trendArgs{Days: 30}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &a); err != nil {
			return nil, nil, fmt.Errorf("invalid arguments: %w", err)
		}
	}
	if a.Days <= 0 {
		return nil, nil, fmt.Errorf("days must be a positive integer")
	}
	rows, err := t.port.Latest(a.Days)
	if err != nil {
		return nil, nil, err
	}
	p, err := marshalPayload(rows)
	if err != nil {
		return nil, nil, err
	}
	return p, rowsSummary(len(rows)), nil
}
