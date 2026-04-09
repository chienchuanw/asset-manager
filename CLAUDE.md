# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Asset Manager is a full-stack personal finance system for tracking investment portfolios, cash flows, subscriptions, installments, and financial analytics. It includes a Discord bot for natural language bookkeeping powered by Google Gemini.

Primary language of comments and UI: Traditional Chinese (zh-TW), with English as fallback.

## Architecture

**Backend (Go/Gin):** Clean layered architecture with constructor-based dependency injection.
```
HTTP Request → Middleware (Auth/CORS) → API Handler → Service → Repository → PostgreSQL/Redis
```
- Entry point: `backend/cmd/api/main.go` — wires all dependencies
- `internal/api/` — HTTP handlers (Gin)
- `internal/service/` — business logic, validation, external API calls
- `internal/repository/` — SQL queries via jmoiron/sqlx
- `internal/models/` — data structures and DTOs
- `internal/discord/` — Discord bot with Gemini-powered NLP parser, adapters for bookkeeping/payments/queries
- `internal/external/` — price API clients (FinMind, CoinGecko, Alpha Vantage)
- `internal/cache/` — Redis caching for price data
- `internal/scheduler/` — cron jobs (daily snapshots, subscription billing)

**Frontend (Next.js 16 App Router):**
- `src/app/` — pages (dashboard, holdings, transactions, cash-flows, analytics, recurring, settings)
- `src/components/` — React components using shadcn/ui + Radix UI
- `src/hooks/` — TanStack Query v5 data-fetching hooks
- `src/lib/api.ts` — backend API client
- `messages/` — i18n translation files (en, zh-TW) via next-intl
- Forms use react-hook-form + zod validation
- Charts via Recharts

## Common Commands

### Backend (from `backend/`)
```bash
make run                    # Start API server (loads .env.local)
make test                   # All tests (seeds test categories first)
make test-unit              # Service + handler tests only (no DB needed)
make test-integration       # Repository tests (needs test DB)
make test-coverage          # Generate coverage.html
make test-watch             # Auto-rerun on file changes
make build                  # Compile to bin/api

make db-create              # Create dev + test databases
make migrate-up-env         # Load .env.local and run migrations
make migrate-test-up-env    # Load .env.test and run migrations
make seed                   # Seed mock data (preserves existing)
make seed-clean             # Clean DB then seed
```

### Frontend (from `frontend/`)
```bash
pnpm dev                    # Dev server on port 3001
pnpm build                  # Production build
pnpm tsc --noEmit           # Type check
```

### Docker (from project root)
```bash
make build                  # Build all Docker images
make up                     # Start all containers (uses .env.production)
make down                   # Stop containers
make health                 # Check API/frontend/nginx health
make logs-backend           # Backend container logs
```

## Testing

- Tests use `testify` (assertions, mocks, suites) and `gotestsum` for output formatting
- Unit tests (`internal/service/`, `internal/api/`) don't need a database
- Integration tests (`internal/repository/`) require a running PostgreSQL test database configured via `.env.test`
- Test DB is seeded with system categories via `scripts/seed_test_categories.sql` before each `make test` run
- Run a single test: `cd backend && go test -run TestName ./internal/service/...`

## Environment Configuration

- Backend uses `.env.local` for dev, `.env.test` for test DB config
- Frontend uses `.env.local` with `NEXT_PUBLIC_API_URL`
- See `.env.template` for all available variables
- Key external API keys: `FINMIND_API_KEY`, `COINGECKO_API_KEY`, `ALPHA_VANTAGE_API_KEY`, `GEMINI_API_KEY`
- Redis is optional — price service degrades gracefully without it

## Database

- PostgreSQL with 28+ sequential migrations in `backend/migrations/`
- Uses `golang-migrate` for schema management
- FIFO cost basis calculation for investment P&L
- Key tables: transactions, cash_flows, categories, subscriptions, installments, bank_accounts, credit_cards, asset_snapshots, performance_snapshots

## Deployment

- Docker Compose orchestration (PostgreSQL, Redis, backend, frontend, Nginx, Certbot)
- CI/CD via GitHub Actions: test → build → deploy to AWS EC2
- Nginx reverse proxy with Let's Encrypt SSL
- Discord webhook notifications on deploy status
