# MCP Server

An in-process Model Context Protocol server that lets AI agents read and modify
portfolio data. It speaks MCP over **stdio** and reuses the same database and
service layer as the API (no HTTP, no API keys). Phase 1 added read tools;
Phase 2 adds write tools behind a dry-run/confirm guardrail with fail-closed
auditing.

## Write tools (Phase 2)

`create_transaction`, `update_transaction`, `delete_transaction`,
`reconcile_holding`. Every write tool is **two-step**:

- Without `confirm:true` (default) the tool returns a **dry-run preview** of
  the effect and persists nothing.
- With `confirm:true` it executes through the existing service layer (FIFO,
  realized-profit and reconciliation row-locks preserved).

Writes are **fail-closed audited**: a `pending` row is written to
`agent_audit_log` *before* the mutation; if that insert fails the write is
aborted and nothing is persisted. After execution the row is finalized to
`success`/`error` with `executed_at` set. (Read tools remain best-effort
audited, per Phase 1.)

## Tools

| Tool | Arguments | Description |
|------|-----------|-------------|
| `get_holdings` | — | All holdings: quantity, avg cost, market value, unrealized P&L |
| `get_holding` | `symbol` | A single holding by symbol |
| `list_transactions` | `from?` `to?` (YYYY-MM-DD), `type?`, `symbol?`, `limit?` | Filtered transactions |
| `get_transaction` | `id` (UUID) | A single transaction |
| `get_analytics` | `range?` (`week`\|`month`\|`quarter`\|`year`\|`all`, default `all`) | Realized/unrealized P&L summary + top assets |
| `get_allocation` | — | Allocation by asset type |
| `get_performance_trend` | `days?` (positive int, default `30`) | Asset-value trend for the latest N days |

All tools return raw JSON (numeric values + currency codes, not preformatted
strings). Invalid arguments and not-found lookups return a structured tool
error; the server never panics.

## Auditing

Every tool call (success or error) writes one row to the `agent_audit_log`
table: tool name, arguments, status, error, a lightweight result summary, and
duration. Logging is **best-effort** in Phase 1 — if the audit insert fails the
server logs to stderr but still returns the result (a logging hiccup must not
block reads). Fail-closed behavior is deferred to Phase 2.

## Running

Requires `.env.local` in `backend/` (same DB config as the API). Build/run:

```bash
make mcp-build   # -> backend/bin/mcp
# or
make mcp-run
```

## Configuring an MCP client

Register the built binary as a stdio MCP server. Example client config:

```json
{
  "mcpServers": {
    "asset-manager": {
      "command": "/absolute/path/to/asset-manager/backend/bin/mcp"
    }
  }
}
```

The working directory must contain `.env.local` (or run the binary from
`backend/`) so the database connection resolves.
