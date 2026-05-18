package mcpserver

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeHoldingsPort struct {
	list []*models.Holding
	one  *models.Holding
	err  error
}

func (f fakeHoldingsPort) List() ([]*models.Holding, error)         { return f.list, f.err }
func (f fakeHoldingsPort) Get(string) (*models.Holding, error)      { return f.one, f.err }

func TestGetHoldingsTool_Success(t *testing.T) {
	tool := NewGetHoldingsTool(fakeHoldingsPort{list: []*models.Holding{{Symbol: "2330"}}})
	payload, summary, err := tool.Run(json.RawMessage(`{}`))
	require.NoError(t, err)
	assert.Contains(t, string(payload), "2330")
	assert.JSONEq(t, `{"rows":1}`, string(summary))
}

func TestGetHoldingTool_MissingSymbol(t *testing.T) {
	tool := NewGetHoldingTool(fakeHoldingsPort{})
	_, _, err := tool.Run(json.RawMessage(`{}`))
	assert.ErrorContains(t, err, "symbol")
}

func TestGetHoldingTool_NotFound(t *testing.T) {
	tool := NewGetHoldingTool(fakeHoldingsPort{err: errors.New("nope")})
	_, _, err := tool.Run(json.RawMessage(`{"symbol":"ZZZ"}`))
	assert.ErrorContains(t, err, "not found")
}

type fakeTxPort struct {
	list []*models.Transaction
	one  *models.Transaction
	err  error
}

func (f fakeTxPort) List(repository.TransactionFilters) ([]*models.Transaction, error) {
	return f.list, f.err
}
func (f fakeTxPort) Get(string) (*models.Transaction, error) { return f.one, f.err }

func TestListTransactionsTool_BadDate(t *testing.T) {
	tool := NewListTransactionsTool(fakeTxPort{})
	_, _, err := tool.Run(json.RawMessage(`{"from":"18-05-2026"}`))
	assert.ErrorContains(t, err, "from")
}

func TestListTransactionsTool_Success(t *testing.T) {
	tool := NewListTransactionsTool(fakeTxPort{list: []*models.Transaction{{Symbol: "AAPL"}}})
	payload, summary, err := tool.Run(json.RawMessage(`{"symbol":"AAPL"}`))
	require.NoError(t, err)
	assert.Contains(t, string(payload), "AAPL")
	assert.JSONEq(t, `{"rows":1}`, string(summary))
}

func TestGetTransactionTool_MissingID(t *testing.T) {
	tool := NewGetTransactionTool(fakeTxPort{})
	_, _, err := tool.Run(json.RawMessage(`{}`))
	assert.ErrorContains(t, err, "id")
}

type fakeAnalyticsPort struct{ res *AnalyticsResult }

func (f fakeAnalyticsPort) Get(models.TimeRange) (*AnalyticsResult, error) { return f.res, nil }

func TestGetAnalyticsTool_InvalidRange(t *testing.T) {
	tool := NewGetAnalyticsTool(fakeAnalyticsPort{})
	_, _, err := tool.Run(json.RawMessage(`{"range":"decade"}`))
	assert.ErrorContains(t, err, "invalid range")
}

func TestGetAnalyticsTool_DefaultRange(t *testing.T) {
	tool := NewGetAnalyticsTool(fakeAnalyticsPort{res: &AnalyticsResult{}})
	_, summary, err := tool.Run(json.RawMessage(`{}`))
	require.NoError(t, err)
	assert.JSONEq(t, `{"ok":true}`, string(summary))
}

type fakeTrendPort struct{ rows []models.PerformanceTrendPoint }

func (f fakeTrendPort) Latest(int) ([]models.PerformanceTrendPoint, error) { return f.rows, nil }

func TestGetPerformanceTrendTool_RejectsNonPositiveDays(t *testing.T) {
	tool := NewGetPerformanceTrendTool(fakeTrendPort{})
	_, _, err := tool.Run(json.RawMessage(`{"days":0}`))
	assert.ErrorContains(t, err, "positive")
}

func TestGetPerformanceTrendTool_DefaultDays(t *testing.T) {
	tool := NewGetPerformanceTrendTool(fakeTrendPort{rows: []models.PerformanceTrendPoint{{}}})
	_, summary, err := tool.Run(json.RawMessage(`{}`))
	require.NoError(t, err)
	assert.JSONEq(t, `{"rows":1}`, string(summary))
}
