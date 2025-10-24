# Asset Manager

A personal asset-tracking system supporting Taiwan stocks, U.S. stocks, and cryptocurrencies with manual transaction import, end-of-day valuations, and holdings visualization.

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Tech Stack](#️-tech-stack)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Quick Start](#quick-start)
  - [Manual Setup](#manual-setup)
- [Development](#-development)
  - [Running Tests](#running-tests)
  - [API Documentation](#api-documentation)
- [Architecture](#️-architecture)
- [Phase 1 Status](#-phase-1-status)
- [Roadmap](#️-roadmap)
- [Contributing](#-contributing)

---

## 🎯 Overview

Asset Manager is a full-stack application designed to help users track their investment portfolio across multiple asset classes. The system supports:

- **Taiwan Stocks** (台股)
- **U.S. Stocks** (美股)
- **Cryptocurrencies** (加密貨幣)

The application calculates holdings using the **FIFO (First-In, First-Out)** method and provides end-of-day valuations in TWD or USD.

---

## ✨ Features

### Current Features (Phase 1 - Backend)

- ✅ Transaction management (Create, Read, Update, Delete)
- ✅ Support for multiple asset types (Taiwan stocks, U.S. stocks, crypto)
- ✅ Support for multiple transaction types (buy, sell, dividend, fee)
- ✅ RESTful API with comprehensive error handling
- ✅ PostgreSQL database with migrations
- ✅ Comprehensive test coverage (21 test cases)
- ✅ TDD (Test-Driven Development) approach

### Planned Features (Phase 2+)

- 🔄 Frontend dashboard with React/Next.js
- 🔄 Holdings calculation with FIFO cost basis
- 🔄 Real-time price integration
- 🔄 Profit & Loss (P&L) calculation
- 🔄 Asset allocation visualization
- 🔄 Discord notifications for rebalancing alerts

---

## 🛠️ Tech Stack

### Frontend

- **Framework**: Next.js 16 (App Router)
- **Language**: TypeScript
- **UI Library**: shadcn/ui (Tailwind CSS)
- **State Management**: TanStack Query (React Query)
- **Form Management**: react-hook-form + zod
- **Package Manager**: pnpm

### Backend

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL 12+
- **Cache**: Redis (planned)
- **Migration**: golang-migrate
- **Testing**: testify

### DevOps

- **Deployment**: AWS EC2 (planned)
- **CI/CD**: GitHub Actions (planned)

---

## 📁 Project Structure

```bash
asset-manager/
├── frontend/                 # Next.js frontend application
│   ├── src/
│   │   ├── app/             # Next.js App Router pages
│   │   ├── components/      # React components
│   │   ├── lib/             # Utilities and API clients
│   │   └── types/           # TypeScript type definitions
│   └── package.json
│
├── backend/                  # Go backend application
│   ├── cmd/
│   │   └── api/             # Main application entry point
│   ├── internal/
│   │   ├── models/          # Data models
│   │   ├── repository/      # Data access layer
│   │   ├── service/         # Business logic layer
│   │   ├── api/             # HTTP handlers
│   │   └── db/              # Database connection
│   ├── migrations/          # Database migrations
│   ├── scripts/             # Utility scripts
│   └── go.mod
│
└── README.md                # This file
```

---

## 🚀 Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.21 or higher ([Download](https://golang.org/dl/))
- **PostgreSQL** 12 or higher ([Download](https://www.postgresql.org/download/))
- **Node.js** 18+ and **pnpm** (for frontend)
- **golang-migrate** CLI ([Installation Guide](https://github.com/golang-migrate/migrate))

### Quick Start

#### Backend Setup (Automated)

```bash
cd backend
chmod +x scripts/setup.sh
./scripts/setup.sh
```

This script will:

- ✅ Check Go and PostgreSQL installation
- ✅ Install golang-migrate (if needed)
- ✅ Install Go dependencies
- ✅ Create `.env` file
- ✅ Create databases (optional)
- ✅ Run migrations (optional)

#### Start the Backend Server

```bash
cd backend
make run
```

The API server will start at `http://localhost:8080`.

#### Test the API

```bash
cd backend
chmod +x scripts/test-api.sh
./scripts/test-api.sh
```

---

### Manual Setup

If you prefer manual setup, follow these steps:

#### 1. Backend Setup

```bash
# Navigate to backend directory
cd backend

# Install Go dependencies
make install

# Create environment file
cp .env.example .env
# Edit .env with your database credentials

# Create databases
psql -U postgres
CREATE DATABASE asset_manager;
CREATE DATABASE asset_manager_test;
\q

# Run migrations
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=asset_manager

make migrate-up

# Run tests
make test

# Start the server
make run
```

#### 2. Frontend Setup (Coming in Phase 2)

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
pnpm install

# Create environment file
cp .env.example .env.local
# Edit .env.local with API URL

# Start development server
pnpm dev
```

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

### API Documentation

#### Endpoints

| Method | Path                    | Description                          |
| ------ | ----------------------- | ------------------------------------ |
| GET    | `/health`               | Health check                         |
| POST   | `/api/transactions`     | Create a transaction                 |
| GET    | `/api/transactions`     | List all transactions (with filters) |
| GET    | `/api/transactions/:id` | Get a transaction by ID              |
| PUT    | `/api/transactions/:id` | Update a transaction                 |
| DELETE | `/api/transactions/:id` | Delete a transaction                 |

#### Example: Create a Transaction

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2025-10-22T00:00:00Z",
    "asset_type": "tw-stock",
    "symbol": "2330",
    "name": "TSMC",
    "type": "buy",
    "quantity": 10,
    "price": 620,
    "amount": 6200,
    "fee": 28,
    "note": "Regular investment"
  }'
```

#### Example: List Transactions with Filters

```bash
# Filter by asset type
curl "http://localhost:8080/api/transactions?asset_type=tw-stock"

# Filter by date range
curl "http://localhost:8080/api/transactions?start_date=2025-10-01&end_date=2025-10-31"

# Pagination
curl "http://localhost:8080/api/transactions?limit=10&offset=0"
```

For more examples, see [`backend/README_PHASE1.md`](backend/README_PHASE1.md).

---

## 🏗️ Architecture

### Backend Architecture

The backend follows a **clean architecture** pattern with clear separation of concerns:

```bash
Client Request
      ↓
API Handler Layer (HTTP handling, request/response formatting)
      ↓
Service Layer (Business logic, validation)
      ↓
Repository Layer (Data access, SQL queries)
      ↓
Database (PostgreSQL)
```

#### Layers

1. **API Handler Layer** (`internal/api/`)

   - HTTP request/response handling
   - Input validation and parsing
   - Error handling and status codes

2. **Service Layer** (`internal/service/`)

   - Business logic implementation
   - Data validation
   - Orchestration of multiple repositories

3. **Repository Layer** (`internal/repository/`)

   - Database CRUD operations
   - SQL query construction
   - Data mapping

4. **Models Layer** (`internal/models/`)
   - Data structure definitions
   - Input/output DTOs
   - Validation methods

For detailed architecture documentation, see [`backend/ARCHITECTURE.md`](backend/ARCHITECTURE.md).

---

## ✅ Phase 1 Status

**Phase 1: Backend Transaction API** - ✅ **COMPLETED**

### Completed Work

- ✅ Database schema design and migrations
- ✅ Transaction model with validation
- ✅ Repository layer with CRUD operations
- ✅ Service layer with business logic
- ✅ API handlers with RESTful endpoints
- ✅ Comprehensive test coverage (21 test cases)
- ✅ Development tools (Makefile, scripts)
- ✅ Complete documentation

### Test Coverage

- **Repository Layer**: 7 integration tests
- **Service Layer**: 8 unit tests (with mocks)
- **API Handler Layer**: 6 unit tests (with mocks)
- **Total**: 21 test cases

### Documentation

- [`backend/QUICK_START.md`](backend/QUICK_START.md) - Quick start guide (5 minutes)
- [`backend/README_PHASE1.md`](backend/README_PHASE1.md) - Detailed implementation guide
- [`backend/ARCHITECTURE.md`](backend/ARCHITECTURE.md) - Architecture documentation
- [`backend/PHASE1_SUMMARY.md`](backend/PHASE1_SUMMARY.md) - Phase 1 completion summary
- [`backend/FILES_CREATED.md`](backend/FILES_CREATED.md) - File listing

---

## 🗺️ Roadmap

### Phase 1: Backend Transaction API ✅ COMPLETED

- ✅ Database design and migrations
- ✅ Transaction CRUD API
- ✅ Comprehensive testing
- ✅ Documentation

### Phase 2: Frontend Integration 🔄 IN PROGRESS

- [ ] Install frontend dependencies (React Query, react-hook-form, zod)
- [ ] Create API client layer
- [ ] Implement transaction list page
- [ ] Implement add transaction dialog
- [ ] Implement edit/delete functionality

### Phase 3: Holdings Calculation

- [ ] Implement FIFO cost calculation
- [ ] Holdings API endpoints
- [ ] Holdings dashboard page
- [ ] Real-time price integration

### Phase 4: Analytics & Reporting

- [ ] Asset allocation calculation
- [ ] P&L calculation (realized/unrealized)
- [ ] Performance analytics
- [ ] Charts and visualizations

### Phase 5: Advanced Features

- [ ] Discord notifications
- [ ] Rebalancing alerts
- [ ] Multi-currency support
- [ ] Export functionality (CSV/PDF)

---

## 🤝 Contributing

This is a personal project, but suggestions and feedback are welcome!

### Development Workflow

1. Follow TDD (Test-Driven Development) approach
2. Write tests before implementation
3. Ensure all tests pass before committing
4. Follow the coding standards in `.augment/rules/`

### Coding Standards

- **Backend (Go)**

  - Use `gofmt` and `goimports` for formatting
  - Follow clean architecture principles
  - Write comprehensive tests
  - Use meaningful variable names

- **Frontend (TypeScript/React)**
  - Use Prettier for formatting
  - Follow React best practices
  - Use TypeScript strictly
  - Component-based architecture

For detailed coding standards, see [`.augment/rules/coding-standards.md`](.augment/rules/coding-standards.md).

---

## 📄 License

This project is for personal use.

---

## 📞 Support

For questions or issues:

1. Check the documentation in `backend/` directory
2. Review the [Quick Start Guide](backend/QUICK_START.md)
3. Check the [Architecture Documentation](backend/ARCHITECTURE.md)

---

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/), [Gin](https://gin-gonic.com/), and [PostgreSQL](https://www.postgresql.org/)
- Frontend powered by [Next.js](https://nextjs.org/) and [shadcn/ui](https://ui.shadcn.com/)
- Testing with [testify](https://github.com/stretchr/testify)

---

**Last Updated**: 2025-10-23
**Current Version**: Phase 1 Complete
