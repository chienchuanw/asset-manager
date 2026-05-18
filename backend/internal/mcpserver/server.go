package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/audit"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// auditSink is the narrow logging contract the dispatcher needs (kept
// SDK-independent so dispatch logic is unit-testable without a DB).
type auditSink interface {
	Record(tool, status string, args, summary json.RawMessage, durationMS int, errMsg string)
}

// loggerSink adapts *audit.Logger to auditSink.
type loggerSink struct{ l *audit.Logger }

func NewLoggerSink(l *audit.Logger) auditSink { return loggerSink{l: l} }

func (s loggerSink) Record(tool, status string, args, summary json.RawMessage, durationMS int, errMsg string) {
	s.l.Record(audit.Entry{
		Tool:          tool,
		Arguments:     args,
		Status:        status,
		Error:         errMsg,
		ResultSummary: summary,
		DurationMS:    durationMS,
	})
}

// Dispatcher routes a tool name to its Tool and audits every call.
type Dispatcher struct {
	tools map[string]Tool
	audit auditSink
}

func NewDispatcher(tools []Tool, a auditSink) *Dispatcher {
	m := make(map[string]Tool, len(tools))
	for _, t := range tools {
		m[t.Name()] = t
	}
	return &Dispatcher{tools: m, audit: a}
}

// SelfAudited marks tools that own their audit lifecycle (the fail-closed
// write tools). The Dispatcher must NOT also best-effort log them, or every
// mutation would produce two audit rows.
type SelfAudited interface{ selfAudited() }

func (d *Dispatcher) Dispatch(name string, args json.RawMessage) (json.RawMessage, error) {
	start := time.Now()
	tool, ok := d.tools[name]
	if !ok {
		d.audit.Record(name, "error", args, nil, ms(start), "unknown tool")
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	_, selfAudits := tool.(SelfAudited)
	payload, summary, err := tool.Run(args)
	if err != nil {
		if !selfAudits {
			d.audit.Record(name, "error", args, nil, ms(start), err.Error())
		}
		return nil, err
	}
	if !selfAudits {
		d.audit.Record(name, "success", args, summary, ms(start), "")
	}
	return payload, nil
}

func ms(start time.Time) int { return int(time.Since(start).Milliseconds()) }

// BuildServer registers every tool with the MCP SDK server. The SDK is
// isolated here; tool/dispatch logic above is testable without it, and this
// function is reusable by both stdio serving and in-memory transport tests.
func BuildServer(d *Dispatcher, tools []Tool) *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "asset-manager",
		Version: "0.1.0",
	}, nil)

	for _, t := range tools {
		name := t.Name()
		handler := func(_ context.Context, _ *mcp.CallToolRequest, in map[string]any) (*mcp.CallToolResult, any, error) {
			raw, mErr := json.Marshal(in)
			if mErr != nil {
				raw = []byte("{}")
			}
			payload, err := d.Dispatch(name, raw)
			if err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
				}, nil, nil
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: string(payload)}},
			}, nil, nil
		}
		mcp.AddTool(srv, &mcp.Tool{Name: name, Description: descriptionFor(name)}, handler)
	}
	return srv
}

// Serve runs the MCP server over stdio.
func Serve(ctx context.Context, d *Dispatcher, tools []Tool) error {
	return BuildServer(d, tools).Run(ctx, &mcp.StdioTransport{})
}

func descriptionFor(name string) string {
	switch name {
	case "get_holdings":
		return "List all current holdings with quantity, average cost, market value and unrealized P&L."
	case "get_holding":
		return "Get a single holding by symbol."
	case "list_transactions":
		return "List transactions, optionally filtered by from/to (YYYY-MM-DD), type, symbol, limit."
	case "get_transaction":
		return "Get a single transaction by id (UUID)."
	case "get_analytics":
		return "Get realized/unrealized P&L summary and top assets for a range (week|month|quarter|year|all)."
	case "get_allocation":
		return "Get current allocation broken down by asset type."
	case "get_performance_trend":
		return "Get asset-value performance trend for the latest N days (default 30)."
	}
	return name
}
