# Asset Manager — Frontend

The web client for Asset Manager, a personal finance system for tracking investment portfolios, cash flows, subscriptions, installments, and financial analytics. Built with the Next.js App Router, it talks to the Go/Gin backend API and renders dashboards, holdings, analytics, and rebalancing views.

Primary UI language is Traditional Chinese (zh-TW) with English as a fallback.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Getting Started](#getting-started)
- [Scripts](#scripts)
- [Configuration](#configuration)
- [Design System](#design-system)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [Contributing](#contributing)

## Overview

This package is the presentation layer of the Asset Manager monorepo. Data is fetched from the backend via the API client in `src/lib/` using TanStack Query; forms use react-hook-form with Zod validation; charts are rendered with Recharts. Internationalisation is handled by next-intl with message catalogs in `messages/`.

## Features

- **Dashboard** — portfolio value, asset trend, and allocation at a glance
- **Holdings** — positions by asset type with FIFO cost-basis P&L
- **Analytics** — realized and unrealized P&L breakdowns over selectable time ranges
- **Cash flows** — income/expense tracking with daily, weekly, and monthly views
- **Rebalance** — allocation drift detection against target weights
- **Recurring** — subscription and installment management
- **Zen Mode** — one-tap toggle that masks all absolute monetary amounts (percentages stay visible) for screenshots and privacy
- **Light/Dark theme** — system-aware theme switching via next-themes
- **i18n** — full zh-TW and en translations via next-intl

## Getting Started

Prerequisites: Node.js 20+ and pnpm.

```bash
pnpm install
pnpm dev
```

The dev server runs on [http://localhost:3001](http://localhost:3001). Point it at a running backend via `NEXT_PUBLIC_API_URL` (see [Configuration](#configuration)).

## Scripts

| Command | Purpose |
|---------|---------|
| `pnpm dev` | Start the dev server on port 3001 |
| `pnpm build` | Production build (also type-checks app code) |
| `pnpm start` | Serve the production build |
| `pnpm test` | Run the Vitest unit/component suite |
| `pnpm test:watch` | Vitest in watch mode |
| `pnpm test:coverage` | Vitest with coverage report |
| `pnpm test:e2e` | Playwright BDD end-to-end tests |
| `pnpm tsc --noEmit` | Type-check |

## Configuration

Frontend config lives in `.env.local`:

| Variable | Description |
|----------|-------------|
| `NEXT_PUBLIC_API_URL` | Base URL of the backend API |

See the repository root `.env.template` for the full set of backend variables.

## Design System

The design foundation (sub-project A) provides a reusable, token-driven base:

- **Semantic design tokens** — `src/app/globals.css` defines color/space/typography tokens including finance-semantic `gain`/`loss`/`neutral` and `surface`, with full light and dark value sets (OKLCH, Tailwind 4 `@theme`).
- **Chart theme** — `src/lib/chartTheme.ts` derives all chart colors, grid, and tooltip styles from the tokens, so charts stay consistent across themes.
- **Theming** — `src/providers/ThemeProvider.tsx` wires next-themes; `ThemeToggle` in the app shell switches light/dark.
- **Zen Mode** — `src/providers/ZenModeProvider.tsx` holds a localStorage-persisted flag; the `<Money>` component (`src/components/common/Money.tsx`) renders currency and masks it when Zen is active. `ZenModeToggle` lives in the shell.
- **Shared state components** — `src/components/ui/states/` provides `EmptyState`, `LoadingState`, and `ErrorState` (with retry), backed by the `states` i18n namespace.
- **Unified formatting** — `src/lib/format.ts` is the single `formatCurrency` implementation (glyph prefix for TWD/USD, currency-code fallback otherwise).

Subsequent phases (returns engine, multi-dimensional allocation, holdings detail pages, X-ray rules, FIRE projection) build on this foundation.

## Project Structure

```text
frontend/
├── src/
│   ├── app/              # App Router pages (dashboard, holdings, analytics, ...)
│   ├── components/       # UI components (shadcn/ui + Radix), charts, dialogs
│   │   ├── common/       # Money, ThemeToggle, ZenModeToggle, LanguageSwitcher
│   │   ├── ui/states/    # Shared Empty/Loading/Error states
│   │   └── layout/       # App shell
│   ├── providers/        # Theme, ZenMode, Query, Auth, Locale providers
│   ├── hooks/            # TanStack Query data hooks
│   ├── lib/              # API client, format, chartTheme
│   └── test/             # Vitest setup and provider-aware render utils
├── messages/             # next-intl catalogs (en.json, zh-TW.json)
└── package.json
```

## Testing

Unit and component tests use Vitest with jsdom and MSW for API mocking. `src/test/utils.tsx` exposes `renderWithProviders` (wraps Query, next-intl, Theme, and Zen providers); `src/test/setup.ts` registers jest-dom matchers and a global `matchMedia` mock. End-to-end tests use Playwright BDD (`pnpm test:e2e`).

```bash
pnpm test
```

## Contributing

Branches follow the `issues/<N>` convention and target `dev` via pull request. Run `pnpm test` and `pnpm build` before opening a PR. Commit messages use a `type(scope): summary` style with no co-author trailers, matching the existing history.
