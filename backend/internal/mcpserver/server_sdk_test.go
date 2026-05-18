package mcpserver

import (
	"context"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/audit"
	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// End-to-end through the real SDK over in-memory transports: a client lists
// the 7 tools and calls get_holdings; a real audit row is written.
func TestSDKServer_ListAndCallToolsWithAudit(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	_, err := db.Exec("DELETE FROM agent_audit_log")
	require.NoError(t, err)

	holdings := fakeHoldingsPort{list: []*models.Holding{{Symbol: "2330", Name: "TSMC"}}}
	tools := []Tool{
		NewGetHoldingsTool(holdings),
		NewGetHoldingTool(holdings),
		NewListTransactionsTool(fakeTxPort{}),
		NewGetTransactionTool(fakeTxPort{}),
		NewGetAnalyticsTool(fakeAnalyticsPort{res: &AnalyticsResult{}}),
		NewGetAllocationTool(fakeAllocPortE2E{}),
		NewGetPerformanceTrendTool(fakeTrendPort{}),
	}
	sink := NewLoggerSink(audit.NewLogger(audit.NewPostgresAuditRepository(db)))
	srv := BuildServer(NewDispatcher(tools, sink), tools)

	ctx := context.Background()
	clientT, serverT := mcp.NewInMemoryTransports()
	ss, err := srv.Connect(ctx, serverT, nil)
	require.NoError(t, err)
	defer ss.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	cs, err := client.Connect(ctx, clientT, nil)
	require.NoError(t, err)
	defer cs.Close()

	listed, err := cs.ListTools(ctx, nil)
	require.NoError(t, err)
	names := map[string]bool{}
	for _, tl := range listed.Tools {
		names[tl.Name] = true
	}
	for _, want := range []string{
		"get_holdings", "get_holding", "list_transactions", "get_transaction",
		"get_analytics", "get_allocation", "get_performance_trend",
	} {
		assert.True(t, names[want], "tool %s should be listed", want)
	}
	assert.Len(t, listed.Tools, 7)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "get_holdings"})
	require.NoError(t, err)
	require.False(t, res.IsError)
	require.NotEmpty(t, res.Content)
	tc, ok := res.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, tc.Text, "2330")

	var tool, status string
	row := db.QueryRow("SELECT tool, status FROM agent_audit_log WHERE tool='get_holdings' LIMIT 1")
	require.NoError(t, row.Scan(&tool, &status))
	assert.Equal(t, "success", status)
}

type fakeAllocPortE2E struct{}

func (fakeAllocPortE2E) ByType() ([]models.AllocationByType, error) {
	return []models.AllocationByType{}, nil
}
