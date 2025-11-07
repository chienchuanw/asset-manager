# Makefile for Asset Manager Backend

# 使用 zsh 作為 shell（確保能找到 go 指令）
SHELL := /bin/zsh

# 載入測試環境變數
ifneq (,$(wildcard .env.test))
    include .env.test
    export
endif

# 顏色定義
GREEN  := \033[0;32m
RED    := \033[0;31m
YELLOW := \033[0;33m
BLUE   := \033[0;34m
NC     := \033[0m # No Color

# Go 工具路徑
GOTESTSUM := ./scripts/run-tests.sh

.PHONY: help install test test-verbose test-unit test-integration test-coverage test-watch migrate-up migrate-down migrate-create run build clean db-create db-drop seed seed-clean

# 顯示幫助訊息
help:
	@echo "$(BLUE)Asset Manager Backend - Available Commands:$(NC)"
	@echo ""
	@echo "$(YELLOW)Testing:$(NC)"
	@echo "  $(GREEN)make test$(NC)             - 執行所有測試（使用 gotestsum）"
	@echo "  $(GREEN)make test-verbose$(NC)     - 執行所有測試（詳細模式）"
	@echo "  $(GREEN)make test-unit$(NC)        - 執行單元測試（Service + API Handler）"
	@echo "  $(GREEN)make test-integration$(NC) - 執行整合測試（Repository）"
	@echo "  $(GREEN)make test-coverage$(NC)    - 執行測試並產生覆蓋率報告"
	@echo "  $(GREEN)make test-watch$(NC)       - 執行測試（watch 模式）"
	@echo ""
	@echo "$(YELLOW)Database:$(NC)"
	@echo "  $(GREEN)make db-create$(NC)        - 建立開發和測試資料庫"
	@echo "  $(GREEN)make db-drop$(NC)          - 刪除開發和測試資料庫"
	@echo "  $(GREEN)make migrate-up$(NC)       - 執行資料庫 migration（開發環境）"
	@echo "  $(GREEN)make migrate-down$(NC)     - 回滾資料庫 migration（開發環境）"
	@echo "  $(GREEN)make migrate-test-up$(NC)  - 執行資料庫 migration（測試環境）"
	@echo "  $(GREEN)make migrate-up-env$(NC)   - 載入 .env.local 並執行 migration"
	@echo "  $(GREEN)make migrate-test-up-env$(NC) - 載入 .env.test 並執行 migration"
	@echo ""
	@echo "$(YELLOW)Development:$(NC)"
	@echo "  $(GREEN)make install$(NC)          - 安裝所有依賴套件"
	@echo "  $(GREEN)make run$(NC)              - 啟動 API 伺服器"
	@echo "  $(GREEN)make build$(NC)            - 編譯應用程式"
	@echo "  $(GREEN)make clean$(NC)            - 清理編譯產物"
	@echo ""
	@echo "$(YELLOW)Data Seeding:$(NC)"
	@echo "  $(GREEN)make seed$(NC)             - 匯入 Mock 資料（保留現有資料）"
	@echo "  $(GREEN)make seed-clean$(NC)       - 清空資料庫並匯入 Mock 資料"
	@echo ""

# 安裝依賴套件
install:
	@echo "Installing dependencies..."
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/suite
	go get github.com/stretchr/testify/mock
	go get github.com/google/uuid
	go mod tidy
	@echo "Dependencies installed successfully!"

# 執行所有測試（使用 gotestsum，彩色輸出）
test:
	@printf "\033[0;34mSeeding test database with system categories...\033[0m\n"
	@PGPASSWORD=$(TEST_DB_PASSWORD) psql -h $(TEST_DB_HOST) -p $(TEST_DB_PORT) -U $(TEST_DB_USER) -d $(TEST_DB_NAME) -f scripts/seed_test_categories.sql > /dev/null 2>&1 || true
	@printf "\033[0;34mRunning all tests...\033[0m\n"
	@$(GOTESTSUM) --format testname -- -cover ./...

# 執行所有測試（詳細模式）
test-verbose:
	@printf "\033[0;34mRunning all tests (verbose)...\033[0m\n"
	@$(GOTESTSUM) --format standard-verbose -- -cover ./...

# 執行單元測試（不需要資料庫）
test-unit:
	@printf "\033[0;34mRunning unit tests...\033[0m\n"
	@$(GOTESTSUM) --format testname -- -cover ./internal/service/... ./internal/api/...

# 執行整合測試（需要資料庫）
test-integration:
	@printf "\033[0;34mRunning integration tests...\033[0m\n"
	@printf "\033[1;33mMake sure test database is running and TEST_DB_* environment variables are set\033[0m\n"
	@$(GOTESTSUM) --format testname -- -cover ./internal/repository/...

# 執行測試並產生覆蓋率報告
test-coverage:
	@printf "\033[0;34mRunning tests with coverage report...\033[0m\n"
	@$(GOTESTSUM) --format testname -- -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@printf "\033[0;32mCoverage report generated: coverage.html\033[0m\n"

# 執行測試（watch 模式，檔案變更時自動重新執行）
test-watch:
	@printf "\033[0;34mRunning tests in watch mode...\033[0m\n"
	@$(GOTESTSUM) --watch --format testname -- -cover ./...

# 建立資料庫
db-create:
	@echo "$(BLUE)Creating databases...$(NC)"
	@set -a && . ./.env.local && set +a && \
	psql -U $$DB_USER -h $$DB_HOST -p $$DB_PORT -c "CREATE DATABASE $$DB_NAME;" 2>/dev/null || echo "$(YELLOW)Database $$DB_NAME already exists$(NC)"
	@set -a && . ./.env.test && set +a && \
	psql -U $$TEST_DB_USER -h $$TEST_DB_HOST -p $$TEST_DB_PORT -c "CREATE DATABASE $$TEST_DB_NAME;" 2>/dev/null || echo "$(YELLOW)Database $$TEST_DB_NAME already exists$(NC)"
	@echo "$(GREEN)Databases created successfully!$(NC)"

# 刪除資料庫
db-drop:
	@echo "$(RED)Dropping databases...$(NC)"
	@set -a && . ./.env.local && set +a && \
	psql -U $$DB_USER -h $$DB_HOST -p $$DB_PORT -c "DROP DATABASE IF EXISTS $$DB_NAME;"
	@set -a && . ./.env.test && set +a && \
	psql -U $$TEST_DB_USER -h $$TEST_DB_HOST -p $$TEST_DB_PORT -c "DROP DATABASE IF EXISTS $$TEST_DB_NAME;"
	@echo "$(GREEN)Databases dropped successfully!$(NC)"

# 執行 migration（開發環境）
migrate-up:
	@echo "Running migrations (development)..."
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

# 回滾 migration（開發環境）
migrate-down:
	@echo "Rolling back migrations (development)..."
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# 執行 migration（測試環境）
migrate-test-up:
	@echo "Running migrations (test)..."
	migrate -path migrations -database "postgresql://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@$(TEST_DB_HOST):$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable" up

# 建立新的 migration
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# 啟動 API 伺服器
run:
	@echo "Starting API server..."
	go run cmd/api/main.go

# 編譯應用程式
build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go
	@echo "Build complete! Binary: bin/api"

# 清理編譯產物
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Clean complete!"

# 載入環境變數並執行 migration
migrate-up-env:
	@echo "Loading .env and running migrations..."
	@set -a && . ./.env.local && set +a && \
	migrate -path migrations -database "postgresql://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" up

# 載入測試環境變數並執行 migration
migrate-test-up-env:
	@echo "Loading .env.test and running migrations..."
	@set -a && . ./.env.test && set +a && \
	migrate -path migrations -database "postgresql://$$TEST_DB_USER:$$TEST_DB_PASSWORD@$$TEST_DB_HOST:$$TEST_DB_PORT/$$TEST_DB_NAME?sslmode=disable" up

# 匯入 Mock 資料（保留現有資料）
seed:
	@echo "$(BLUE)Seeding database with mock data...$(NC)"
	@go run cmd/seed/main.go
	@echo "$(GREEN)Seeding completed!$(NC)"

# 清空資料庫並匯入 Mock 資料
seed-clean:
	@echo "$(RED)Cleaning database and seeding with mock data...$(NC)"
	@go run cmd/seed/main.go --clean
	@echo "$(GREEN)Database cleaned and seeded successfully!$(NC)"

