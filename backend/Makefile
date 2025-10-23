# Makefile for Asset Manager Backend

.PHONY: help install test test-unit test-integration migrate-up migrate-down migrate-create run build clean

# 顯示幫助訊息
help:
	@echo "Asset Manager Backend - Available Commands:"
	@echo ""
	@echo "  make install          - 安裝所有依賴套件"
	@echo "  make test             - 執行所有測試"
	@echo "  make test-unit        - 執行單元測試（Service + API Handler）"
	@echo "  make test-integration - 執行整合測試（Repository）"
	@echo "  make migrate-up       - 執行資料庫 migration（開發環境）"
	@echo "  make migrate-down     - 回滾資料庫 migration（開發環境）"
	@echo "  make migrate-test-up  - 執行資料庫 migration（測試環境）"
	@echo "  make run              - 啟動 API 伺服器"
	@echo "  make build            - 編譯應用程式"
	@echo "  make clean            - 清理編譯產物"
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

# 執行所有測試
test:
	@echo "Running all tests..."
	go test ./... -v -cover

# 執行單元測試（不需要資料庫）
test-unit:
	@echo "Running unit tests..."
	go test ./internal/service/... ./internal/api/... -v -cover

# 執行整合測試（需要資料庫）
test-integration:
	@echo "Running integration tests..."
	@echo "Make sure test database is running and TEST_DB_* environment variables are set"
	go test ./internal/repository/... -v -cover

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

