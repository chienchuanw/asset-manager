#!/bin/bash

# 顏色定義
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Asset Manager - Quick Setup${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 載入 .env.local
if [ -f .env.local ]; then
    echo -e "${GREEN}✓ Loading .env.local...${NC}"
    set -a
    source .env.local
    set +a
else
    echo -e "${RED}✗ .env.local not found!${NC}"
    echo -e "${YELLOW}Please create .env.local first${NC}"
    exit 1
fi

# 載入 .env.test
if [ -f .env.test ]; then
    echo -e "${GREEN}✓ Loading .env.test...${NC}"
    set -a
    source .env.test
    set +a
else
    echo -e "${YELLOW}⚠ .env.test not found, skipping test database setup${NC}"
fi

echo ""
echo -e "${BLUE}Step 1: Creating databases...${NC}"

# 建立開發資料庫
echo -e "${YELLOW}Creating development database: ${DB_NAME}${NC}"
psql -U $DB_USER -h $DB_HOST -p $DB_PORT -c "CREATE DATABASE $DB_NAME;" 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Development database created${NC}"
else
    echo -e "${YELLOW}⚠ Development database already exists${NC}"
fi

# 建立測試資料庫
if [ ! -z "$TEST_DB_NAME" ]; then
    echo -e "${YELLOW}Creating test database: ${TEST_DB_NAME}${NC}"
    psql -U $TEST_DB_USER -h $TEST_DB_HOST -p $TEST_DB_PORT -c "CREATE DATABASE $TEST_DB_NAME;" 2>/dev/null
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Test database created${NC}"
    else
        echo -e "${YELLOW}⚠ Test database already exists${NC}"
    fi
fi

echo ""
echo -e "${BLUE}Step 2: Running migrations...${NC}"

# 執行開發資料庫 migration
echo -e "${YELLOW}Running migrations for development database...${NC}"
migrate -path migrations -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Development migrations completed${NC}"
else
    echo -e "${RED}✗ Development migrations failed${NC}"
fi

# 執行測試資料庫 migration
if [ ! -z "$TEST_DB_NAME" ]; then
    echo -e "${YELLOW}Running migrations for test database...${NC}"
    migrate -path migrations -database "postgresql://$TEST_DB_USER:$TEST_DB_PASSWORD@$TEST_DB_HOST:$TEST_DB_PORT/$TEST_DB_NAME?sslmode=disable" up
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Test migrations completed${NC}"
    else
        echo -e "${RED}✗ Test migrations failed${NC}"
    fi
fi

echo ""
echo -e "${BLUE}Step 3: Running tests...${NC}"

# 執行單元測試
echo -e "${YELLOW}Running unit tests...${NC}"
go test ./internal/service/... ./internal/api/... -v -cover

# 執行整合測試
if [ ! -z "$TEST_DB_NAME" ]; then
    echo ""
    echo -e "${YELLOW}Running integration tests...${NC}"
    export TEST_DB_HOST=$TEST_DB_HOST
    export TEST_DB_PORT=$TEST_DB_PORT
    export TEST_DB_USER=$TEST_DB_USER
    export TEST_DB_PASSWORD=$TEST_DB_PASSWORD
    export TEST_DB_NAME=$TEST_DB_NAME
    go test ./internal/repository/... -v -cover
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Setup Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo -e "  1. Run ${GREEN}make run${NC} to start the API server"
echo -e "  2. Run ${GREEN}make test${NC} to run all tests"
echo -e "  3. Run ${GREEN}./scripts/test-api.sh${NC} to test the API"
echo ""

