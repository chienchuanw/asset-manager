#!/bin/bash

# Asset Manager Backend Setup Script
# 此腳本會協助你設定開發環境

set -e

echo "========================================="
echo "Asset Manager Backend Setup"
echo "========================================="
echo ""

# 檢查 Go 是否已安裝
echo "Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "✅ Go is installed: $GO_VERSION"
echo ""

# 檢查 PostgreSQL 是否已安裝
echo "Checking PostgreSQL installation..."
if ! command -v psql &> /dev/null; then
    echo "⚠️  PostgreSQL client (psql) is not found."
    echo "   Please make sure PostgreSQL is installed and running."
    echo "   macOS: brew install postgresql"
    echo "   Ubuntu: sudo apt-get install postgresql"
else
    PSQL_VERSION=$(psql --version)
    echo "✅ PostgreSQL client is installed: $PSQL_VERSION"
fi
echo ""

# 檢查 migrate CLI 是否已安裝
echo "Checking golang-migrate..."
if ! command -v migrate &> /dev/null; then
    echo "⚠️  golang-migrate is not installed."
    echo "   Installing golang-migrate..."
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install golang-migrate
        else
            echo "   Please install Homebrew first: https://brew.sh/"
            exit 1
        fi
    else
        # Linux or others
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    fi
    
    echo "✅ golang-migrate installed"
else
    MIGRATE_VERSION=$(migrate -version 2>&1 | head -n 1)
    echo "✅ golang-migrate is installed: $MIGRATE_VERSION"
fi
echo ""

# 安裝 Go 依賴套件
echo "Installing Go dependencies..."
cd "$(dirname "$0")/.."
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/suite
go get github.com/stretchr/testify/mock
go get github.com/google/uuid
go mod tidy
echo "✅ Go dependencies installed"
echo ""

# 建立 .env 檔案（如果不存在）
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cp .env.example .env
    echo "✅ .env file created. Please edit it with your database credentials."
    echo ""
    echo "⚠️  IMPORTANT: Please update the following in .env:"
    echo "   - DB_PASSWORD"
    echo "   - Other database settings if needed"
    echo ""
else
    echo "✅ .env file already exists"
    echo ""
fi

# 詢問是否要建立資料庫
read -p "Do you want to create databases now? (y/n) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "Please enter your PostgreSQL credentials:"
    read -p "PostgreSQL username (default: postgres): " PG_USER
    PG_USER=${PG_USER:-postgres}
    
    read -sp "PostgreSQL password: " PG_PASSWORD
    echo ""
    
    read -p "PostgreSQL host (default: localhost): " PG_HOST
    PG_HOST=${PG_HOST:-localhost}
    
    read -p "PostgreSQL port (default: 5432): " PG_PORT
    PG_PORT=${PG_PORT:-5432}
    
    echo ""
    echo "Creating databases..."
    
    # 建立開發資料庫
    PGPASSWORD=$PG_PASSWORD psql -h $PG_HOST -p $PG_PORT -U $PG_USER -tc "SELECT 1 FROM pg_database WHERE datname = 'asset_manager'" | grep -q 1 || \
    PGPASSWORD=$PG_PASSWORD psql -h $PG_HOST -p $PG_PORT -U $PG_USER -c "CREATE DATABASE asset_manager;"
    
    # 建立測試資料庫
    PGPASSWORD=$PG_PASSWORD psql -h $PG_HOST -p $PG_PORT -U $PG_USER -tc "SELECT 1 FROM pg_database WHERE datname = 'asset_manager_test'" | grep -q 1 || \
    PGPASSWORD=$PG_PASSWORD psql -h $PG_HOST -p $PG_PORT -U $PG_USER -c "CREATE DATABASE asset_manager_test;"
    
    echo "✅ Databases created"
    echo ""
    
    # 執行 migration
    read -p "Do you want to run migrations now? (y/n) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Running migrations..."
        
        # 開發資料庫
        migrate -path migrations -database "postgresql://$PG_USER:$PG_PASSWORD@$PG_HOST:$PG_PORT/asset_manager?sslmode=disable" up
        
        # 測試資料庫
        migrate -path migrations -database "postgresql://$PG_USER:$PG_PASSWORD@$PG_HOST:$PG_PORT/asset_manager_test?sslmode=disable" up
        
        echo "✅ Migrations completed"
    fi
fi

echo ""
echo "========================================="
echo "Setup Complete! 🎉"
echo "========================================="
echo ""
echo "Next steps:"
echo "1. Edit .env file with your database credentials"
echo "2. Run 'make test-unit' to run unit tests"
echo "3. Run 'make test-integration' to run integration tests (requires database)"
echo "4. Run 'make run' to start the API server"
echo ""
echo "For more commands, run 'make help'"
echo ""

