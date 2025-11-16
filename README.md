# Asset Manager

A comprehensive personal finance management system supporting investment portfolio tracking, cash flow management, subscription/installment billing, and financial analytics with real-time valuations and detailed reporting.

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Tech Stack](#️-tech-stack)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Quick Start with Docker Compose](#quick-start-with-docker-compose)
  - [Manual Setup](#manual-setup)
- [Development](#-development)
  - [Running Tests](#running-tests)
  - [Development Commands](#development-commands)
  - [Code Standards](#code-standards)
- [Architecture](#️-architecture)
- [Development Progress](#development-progress)
- [API Endpoints](#api-endpoints)
- [Documentation](#documentation)
- [Contributing](#-contributing)
- [License](#license)
- [Support](#support)
- [Acknowledgments](#acknowledgments)

---

## 🎯 Overview

Asset Manager is a full-stack application designed to help users manage their personal finances comprehensively. The system supports:

### Investment Portfolio Tracking

- **Taiwan Stocks** (台股)
- **U.S. Stocks** (美股)
- **Cryptocurrencies** (加密貨幣)
- **Cash Holdings** (現金)

### Financial Management

- **Cash Flow Tracking** - Income and expense management with categorization
- **Subscription Management** - Track recurring subscriptions and auto-billing
- **Installment Tracking** - Manage payment plans with interest calculations
- **Bank & Credit Card Management** - Multi-account support with grouping

### Analytics & Reporting

- **Holdings Calculation** - FIFO (First-In, First-Out) cost basis calculation
- **Performance Analytics** - Realized/unrealized P&L tracking
- **Asset Allocation** - Portfolio composition visualization
- **Daily Snapshots** - Historical valuation tracking
- **Discord Integration** - Automated daily reports and alerts

---

## ✨ Features

### Investment & Holdings Management ✅

- ✅ Transaction management (Buy, Sell, Dividend, Fee, Tax)
- ✅ Multi-asset support (Taiwan stocks, U.S. stocks, cryptocurrencies, cash)
- ✅ FIFO cost basis calculation
- ✅ Holdings tracking with real-time valuations
- ✅ Asset allocation by type and individual assets
- ✅ CSV import/export for transactions
- ✅ Batch transaction creation

### Financial Analytics ✅

- ✅ Realized profit/loss calculation
- ✅ Unrealized profit/loss tracking
- ✅ Performance trends with daily snapshots
- ✅ Time-range based analytics (week, month, quarter, year, all)
- ✅ Top performing/underperforming assets
- ✅ Asset allocation visualization

### Cash Flow Management ✅

- ✅ Income and expense tracking with categorization
- ✅ Predefined and custom categories
- ✅ Monthly/yearly cash flow reports
- ✅ Summary statistics and trends
- ✅ Discord integration for reports

### Subscription & Installment Management ✅

- ✅ Subscription creation and management
- ✅ Automatic daily billing
- ✅ Installment tracking with interest calculations
- ✅ Payment progress visualization
- ✅ Expiration reminders and alerts
- ✅ Auto-renewal settings

### Financial Account Management ✅

- ✅ Bank account tracking
- ✅ Credit card management
- ✅ Credit card grouping
- ✅ Multi-account support

### System Features ✅

- ✅ JWT authentication
- ✅ Role-based access control
- ✅ Settings management (notifications, preferences)
- ✅ Discord webhook integration
- ✅ Scheduled tasks (daily snapshots, billing, reports)
- ✅ Exchange rate management
- ✅ Graceful API degradation with caching
- ✅ Comprehensive error handling

---

## 🛠️ Tech Stack

### Frontend

- **Framework**: Next.js 16 (App Router)
- **Language**: TypeScript 5
- **UI Library**: shadcn/ui (Tailwind CSS 4)
- **State Management**: TanStack Query 5 (React Query)
- **Form Management**: react-hook-form 7 + zod 4
- **Charts**: Recharts 2
- **Date Handling**: date-fns 4
- **Notifications**: Sonner
- **Package Manager**: pnpm
- **Runtime**: React 19, Node.js 18+

### Backend

- **Language**: Go 1.24
- **Web Framework**: Gin 1.11
- **Database**: PostgreSQL 12+
- **Cache/Queue**: Redis 9
- **Authentication**: JWT (golang-jwt/jwt v5)
- **Migration**: golang-migrate
- **Testing**: testify 1.11
- **Task Scheduling**: robfig/cron v3
- **HTTP Client**: Standard library + custom clients

### DevOps & Deployment

- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Reverse Proxy**: Nginx
- **Deployment Target**: AWS EC2
- **CI/CD**: GitHub Actions (ready for setup)

---

## 📁 Project Structure

```bash
asset-manager/
├── frontend/                      # Next.js frontend application
│   ├── src/
│   │   ├── app/                  # Next.js App Router pages
│   │   │   ├── dashboard/        # Dashboard page
│   │   │   ├── transactions/     # Transaction management
│   │   │   ├── holdings/         # Holdings tracking
│   │   │   ├── cash-flows/       # Cash flow management
│   │   │   ├── recurring/        # Subscriptions & installments
│   │   │   ├── analytics/        # Performance analytics
│   │   │   ├── settings/         # User settings
│   │   │   └── user-management/  # User management
│   │   ├── components/           # React components
│   │   ├── hooks/                # Custom React hooks
│   │   ├── lib/                  # Utilities and API clients
│   │   ├── types/                # TypeScript type definitions
│   │   └── providers/            # Context providers
│   ├── doc/                      # Documentation
│   └── package.json
│
├── backend/                       # Go backend application
│   ├── cmd/
│   │   ├── api/                  # Main API server
│   │   ├── snapshot/             # Snapshot utility
│   │   └── [other utilities]/    # Various CLI tools
│   ├── internal/
│   │   ├── api/                  # HTTP handlers
│   │   ├── models/               # Data models and DTOs
│   │   ├── repository/           # Data access layer
│   │   ├── service/              # Business logic layer
│   │   ├── db/                   # Database connection
│   │   ├── auth/                 # Authentication
│   │   ├── middleware/           # HTTP middleware
│   │   ├── cache/                # Caching layer
│   │   ├── external/             # External API clients
│   │   ├── scheduler/            # Task scheduling
│   │   └── client/               # HTTP clients
│   ├── migrations/               # Database migrations (23 files)
│   ├── mock/                     # Mock data
│   ├── test/                     # Integration tests
│   ├── scripts/                  # Utility scripts
│   ├── doc/                      # Documentation
│   └── go.mod
│
├── scripts/                       # Project-level scripts
├── docker-compose.yml            # Docker Compose configuration
├── nginx.conf                     # Nginx configuration
└── README.md                      # This file
```

---

## 🚀 Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.24 or higher ([Download](https://golang.org/dl/))
- **PostgreSQL** 12 or higher ([Download](https://www.postgresql.org/download/))
- **Redis** 6+ (optional, for caching and scheduling)
- **Node.js** 18+ and **pnpm** ([Installation Guide](https://pnpm.io/))
- **Docker** & **Docker Compose** (optional, for containerized setup)

### Quick Start with Docker Compose

The easiest way to get started is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/chienchuanw/asset-manager.git
cd asset-manager

# Start all services (PostgreSQL, Redis, Backend, Frontend)
docker-compose up -d

# Backend API: http://localhost:8080
# Frontend: http://localhost:3000
```

### Manual Setup

#### 1. Backend Setup

```bash
cd backend

# Install Go dependencies
go mod download

# Create environment file
cp .env.example .env
# Edit .env with your database credentials

# Create databases
psql -U postgres -c "CREATE DATABASE asset_manager;"
psql -U postgres -c "CREATE DATABASE asset_manager_test;"

# Run migrations
make migrate-up

# Run tests
make test

# Start the server
make run
```

The API server will start at `http://localhost:8080`.

#### 2. Frontend Setup

```bash
cd frontend

# Install dependencies
pnpm install

# Create environment file
cp .env.example .env.local
# Edit .env.local with API URL (default: http://localhost:8080)

# Start development server
pnpm dev
```

The frontend will start at `http://localhost:3000`.

---

## 💻 Development

### Running Tests

#### All Tests

```bash
cd backend
make test
```

#### Unit Tests Only (no database required)

```bash
make test-unit
```

#### Integration Tests Only (requires database)

```bash
# Set test database environment variables
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=your_password
export TEST_DB_NAME=asset_manager_test

make test-integration
```

### Development Commands

```bash
# Backend
cd backend

# Install dependencies
make install

# Run linter
make lint

# Format code
make fmt

# Run migrations
make migrate-up
make migrate-down

# Seed test data
make seed

# Frontend
cd frontend

# Install dependencies
pnpm install

# Type check
pnpm tsc --noEmit

# Build
pnpm build

# Start dev server
pnpm dev
```

### Code Standards

- **Backend**: Follow `.augment/rules/coding-standards.md`

  - Use `gofmt` and `goimports` for formatting
  - Use `golangci-lint` for linting
  - Write tests using TDD approach
  - Use meaningful variable names and comments in Chinese

- **Frontend**: Follow `.augment/rules/coding-standards.md`
  - Use Prettier for formatting
  - Use ESLint with TypeScript rules
  - Follow React best practices
  - Use TypeScript strictly

---

## 🏗️ Architecture

### Backend Architecture

The backend follows a **clean architecture** pattern with clear separation of concerns:

```text
HTTP Request
    ↓
Middleware (Auth, CORS, Logging)
    ↓
API Handler Layer (HTTP handling, request/response formatting)
    ↓
Service Layer (Business logic, validation, orchestration)
    ↓
Repository Layer (Data access, SQL queries)
    ↓
Database (PostgreSQL) / Cache (Redis)
```

#### Core Layers

1. **API Handler Layer** (`internal/api/`)

   - HTTP request/response handling
   - Input validation and parsing
   - Error handling and status codes
   - 20+ handlers for different features

2. **Service Layer** (`internal/service/`)

   - Business logic implementation
   - Data validation and transformation
   - Orchestration of multiple repositories
   - External API integration
   - 20+ services for different domains

3. **Repository Layer** (`internal/repository/`)

   - Database CRUD operations
   - SQL query construction
   - Data mapping and transactions
   - 15+ repositories for different entities

4. **Models Layer** (`internal/models/`)
   - Data structure definitions
   - Input/output DTOs
   - Validation methods
   - 30+ model types

#### Supporting Layers

- **Middleware** (`internal/middleware/`) - Authentication, CORS, logging
- **Auth** (`internal/auth/`) - JWT token generation and validation
- **Cache** (`internal/cache/`) - Redis caching layer
- **External** (`internal/external/`) - External API clients
- **Scheduler** (`internal/scheduler/`) - Task scheduling and management
- **Database** (`internal/db/`) - Database connection and initialization

### Frontend Architecture

```text
User Interface (Next.js Pages)
    ↓
React Components (Presentational)
    ↓
Custom Hooks (Data fetching, state management)
    ↓
API Client Layer (HTTP requests)
    ↓
Backend API
```

#### Frontend Structure

- **Pages** (`src/app/`) - Next.js App Router pages
- **Components** (`src/components/`) - Reusable React components
- **Hooks** (`src/hooks/`) - Custom React hooks for data fetching
- **API Client** (`src/lib/`) - API communication layer
- **Types** (`src/types/`) - TypeScript type definitions
- **Providers** (`src/providers/`) - Context providers and configuration

For detailed architecture documentation, see [`backend/doc/ARCHITECTURE.md`](backend/doc/ARCHITECTURE.md).

---

## Development Progress

### Phase 1: Backend Transaction API ✅ COMPLETED

**Status**: Fully implemented and tested

- ✅ Database schema and migrations
- ✅ Transaction CRUD API
- ✅ Multi-asset support (stocks, crypto, cash)
- ✅ CSV import/export functionality
- ✅ Comprehensive test coverage
- ✅ Complete documentation

### Phase 2: Frontend Integration ✅ COMPLETED

**Status**: Fully implemented with all core pages

- ✅ Dashboard with overview statistics
- ✅ Transaction management page
- ✅ Holdings tracking page
- ✅ API client layer with React Query
- ✅ Form validation with react-hook-form + zod
- ✅ Authentication system

### Phase 3: Holdings & Analytics ✅ COMPLETED

**Status**: Fully implemented with FIFO calculations

- ✅ FIFO cost basis calculation
- ✅ Holdings API endpoints
- ✅ Realized/unrealized P&L tracking
- ✅ Asset allocation calculations
- ✅ Performance trends with daily snapshots
- ✅ Analytics dashboard with charts

### Phase 4: Cash Flow Management ✅ COMPLETED

**Status**: Fully implemented with categorization

- ✅ Income/expense tracking
- ✅ Category management
- ✅ Monthly/yearly reports
- ✅ Summary statistics
- ✅ Discord integration for reports

### Phase 5: Subscription & Installment Management ✅ COMPLETED

**Status**: Fully implemented with auto-billing

- ✅ Subscription creation and management
- ✅ Installment tracking with interest
- ✅ Automatic daily billing
- ✅ Payment progress visualization
- ✅ Expiration reminders

### Phase 6: Account Management ✅ COMPLETED

**Status**: Fully implemented

- ✅ Bank account tracking
- ✅ Credit card management
- ✅ Credit card grouping
- ✅ Multi-account support

### Phase 7: System Features & Integration ✅ COMPLETED

**Status**: Fully implemented

- ✅ JWT authentication
- ✅ Settings management
- ✅ Discord webhook integration
- ✅ Scheduled tasks (snapshots, billing, reports)
- ✅ Exchange rate management
- ✅ Graceful API degradation with caching
- ✅ Docker containerization
- ✅ Nginx reverse proxy configuration

---

## 🗺️ Future Roadmap

### Phase 8: Advanced Analytics (Planned)

- [ ] Tax reporting and export
- [ ] Portfolio rebalancing recommendations
- [ ] Risk analysis and metrics
- [ ] Benchmark comparison
- [ ] Custom report generation

### Phase 9: Mobile & Notifications (Planned)

- [ ] Mobile app (React Native)
- [ ] Push notifications
- [ ] Email notifications
- [ ] SMS alerts
- [ ] Webhook support for custom integrations

### Phase 10: Multi-User & Collaboration (Planned)

- [ ] Multi-user support
- [ ] Role-based access control (RBAC)
- [ ] Shared portfolios
- [ ] Audit logging
- [ ] User activity tracking

### Phase 11: Advanced Features (Planned)

- [ ] Machine learning for predictions
- [ ] Automated trading signals
- [ ] Portfolio optimization
- [ ] Tax-loss harvesting recommendations
- [ ] Integration with brokers (API)

### Phase 12: Enterprise Features (Planned)

- [ ] Multi-currency support (beyond TWD/USD)
- [ ] Advanced reporting (PDF/Excel export)
- [ ] Data backup and recovery
- [ ] API rate limiting and quotas
- [ ] White-label support

---

## API Endpoints

### Authentication

- `POST /auth/login` - User login
- `POST /auth/register` - User registration
- `POST /auth/refresh` - Refresh JWT token

### Transactions

- `POST /api/transactions` - Create transaction
- `POST /api/transactions/batch` - Batch create transactions
- `GET /api/transactions` - List transactions (with filters)
- `GET /api/transactions/:id` - Get transaction by ID
- `PUT /api/transactions/:id` - Update transaction
- `DELETE /api/transactions/:id` - Delete transaction
- `GET /api/transactions/template` - Download CSV template
- `POST /api/transactions/parse-csv` - Parse CSV file

### Holdings

- `GET /api/holdings` - Get all holdings
- `GET /api/holdings/:symbol` - Get holding by symbol

### Analytics

- `GET /api/analytics/summary` - Get analytics summary
- `GET /api/analytics/performance` - Get performance data
- `GET /api/analytics/top-assets` - Get top performing assets
- `GET /api/analytics/unrealized` - Get unrealized P&L

### Asset Allocation

- `GET /api/allocation/current` - Get current allocation
- `GET /api/allocation/by-type` - Get allocation by asset type
- `GET /api/allocation/by-asset` - Get allocation by individual asset

### Cash Flows

- `POST /api/cash-flows` - Create cash flow record
- `GET /api/cash-flows` - List cash flows (with filters)
- `GET /api/cash-flows/:id` - Get cash flow by ID
- `PUT /api/cash-flows/:id` - Update cash flow
- `DELETE /api/cash-flows/:id` - Delete cash flow
- `GET /api/cash-flows/summary` - Get cash flow summary

### Categories

- `POST /api/categories` - Create category
- `GET /api/categories` - List categories
- `PUT /api/categories/:id` - Update category
- `DELETE /api/categories/:id` - Delete category

### Subscriptions

- `POST /api/subscriptions` - Create subscription
- `GET /api/subscriptions` - List subscriptions
- `GET /api/subscriptions/expiring-soon` - Get expiring subscriptions
- `PUT /api/subscriptions/:id` - Update subscription
- `DELETE /api/subscriptions/:id` - Delete subscription

### Installments

- `POST /api/installments` - Create installment
- `GET /api/installments` - List installments
- `GET /api/installments/completing-soon` - Get completing installments
- `PUT /api/installments/:id` - Update installment
- `DELETE /api/installments/:id` - Delete installment

### Settings

- `GET /api/settings` - Get user settings
- `PUT /api/settings` - Update user settings

### Discord Integration

- `POST /api/discord/test` - Test Discord webhook
- `POST /api/discord/daily-report` - Send daily report

For complete API documentation, see the backend documentation files in `backend/doc/`.

---

## 🤝 Contributing

This is a personal project, but suggestions and feedback are welcome!

### Development Workflow

1. Follow TDD (Test-Driven Development) approach
2. Write tests before implementation
3. Ensure all tests pass before committing
4. Follow the coding standards in `.augment/rules/`
5. Use meaningful commit messages in English

### Coding Standards

- **Backend (Go)**

  - Use `gofmt` and `goimports` for formatting
  - Follow clean architecture principles
  - Write comprehensive tests
  - Use meaningful variable names and Chinese comments
  - Handle all errors explicitly

- **Frontend (TypeScript/React)**
  - Use Prettier for formatting
  - Follow React best practices
  - Use TypeScript strictly
  - Component-based architecture
  - Write Chinese comments for complex logic

For detailed coding standards, see [`.augment/rules/coding-standards.md`](.augment/rules/coding-standards.md).

---

## Documentation

### Backend Documentation

- [`backend/doc/ARCHITECTURE.md`](backend/doc/ARCHITECTURE.md) - System architecture
- [`backend/doc/QUICK_START.md`](backend/doc/QUICK_START.md) - Quick start guide
- [`backend/doc/TESTING_GUIDE.md`](backend/doc/TESTING_GUIDE.md) - Testing guide
- [`backend/doc/DEPLOYMENT.md`](backend/doc/DEPLOYMENT.md) - Deployment guide
- [`backend/doc/ANALYTICS_COMPLETE_SUMMARY.md`](backend/doc/ANALYTICS_COMPLETE_SUMMARY.md) - Analytics feature documentation

### Frontend Documentation

- [`frontend/doc/PHASE_6_SUMMARY.md`](frontend/doc/PHASE_6_SUMMARY.md) - Frontend implementation summary
- [`frontend/README.md`](frontend/README.md) - Frontend setup guide

---

## License

This project is for personal use.

---

## Support

For questions or issues:

1. Check the documentation in `backend/doc/` and `frontend/doc/` directories
2. Review the [Quick Start Guide](backend/doc/QUICK_START.md)
3. Check the [Architecture Documentation](backend/doc/ARCHITECTURE.md)
4. Review the [Testing Guide](backend/doc/TESTING_GUIDE.md)

---

## Acknowledgments

### Backend Technologies

- [Go](https://golang.org/) - Programming language
- [Gin](https://gin-gonic.com/) - Web framework
- [PostgreSQL](https://www.postgresql.org/) - Database
- [Redis](https://redis.io/) - Cache and message broker
- [testify](https://github.com/stretchr/testify) - Testing framework
- [golang-jwt](https://github.com/golang-jwt/jwt) - JWT authentication

### Frontend Technologies

- [Next.js](https://nextjs.org/) - React framework
- [React](https://react.dev/) - UI library
- [TypeScript](https://www.typescriptlang.org/) - Type safety
- [shadcn/ui](https://ui.shadcn.com/) - Component library
- [Tailwind CSS](https://tailwindcss.com/) - Styling
- [TanStack Query](https://tanstack.com/query/) - Data fetching
- [Recharts](https://recharts.org/) - Charting library

### DevOps

- [Docker](https://www.docker.com/) - Containerization
- [Nginx](https://nginx.org/) - Reverse proxy

---

**Last Updated**: 2025-11-16
**Current Version**: Phase 7 Complete (All Core Features Implemented)
**Status**: Production Ready
