# 多階段建置 Dockerfile for Go Backend
# Stage 1: 建置階段
# Base images 以 digest 釘選，避免浮動 tag 在重新建置時拉到不同內容
# Dependabot (.github/dependabot.yml) 會自動開 PR 升級這些 digest
FROM golang:1.25-alpine@sha256:8d22e29d960bc50cd025d93d5b7c7d220b1ee9aa7a239b3c8f55a57e987e8d45 AS builder

# 安裝必要工具
RUN apk add --no-cache git make

# 設定工作目錄
WORKDIR /app

# 複製 go mod 檔案
COPY go.mod go.sum ./

# 下載依賴套件
RUN go mod download

# 複製所有原始碼
COPY . .

# 編譯應用程式
# CGO_ENABLED=0 確保編譯出靜態連結的執行檔
# -ldflags="-w -s" 減少執行檔大小
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/api cmd/api/main.go

# Stage 2: 執行階段
# alpine:3.20 + digest 釘選（之前用 alpine:latest，浮動 tag 是穩定性與安全風險）
FROM alpine:3.23@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

# 安裝必要工具
# ca-certificates: HTTPS 連線需要
# postgresql-client: 執行 migrations 需要
# tzdata: 時區設定
RUN apk add --no-cache ca-certificates postgresql-client tzdata curl

# 設定時區為台北
ENV TZ=Asia/Taipei

# 建立非 root 使用者
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 設定工作目錄
WORKDIR /app

# 從 builder 階段複製編譯好的執行檔
COPY --from=builder /app/bin/api /app/api

# 複製 migrations 檔案
COPY --from=builder /app/migrations /app/migrations

# 安裝 golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

# 建立啟動腳本
RUN echo '#!/bin/sh' > /app/entrypoint.sh && \
    echo 'set -e' >> /app/entrypoint.sh && \
    echo '' >> /app/entrypoint.sh && \
    echo 'echo "Waiting for PostgreSQL to be ready..."' >> /app/entrypoint.sh && \
    echo 'until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do' >> /app/entrypoint.sh && \
    echo '  echo "PostgreSQL is unavailable - sleeping"' >> /app/entrypoint.sh && \
    echo '  sleep 2' >> /app/entrypoint.sh && \
    echo 'done' >> /app/entrypoint.sh && \
    echo '' >> /app/entrypoint.sh && \
    echo 'echo "PostgreSQL is ready!"' >> /app/entrypoint.sh && \
    echo '' >> /app/entrypoint.sh && \
    echo 'echo "Running database migrations..."' >> /app/entrypoint.sh && \
    echo 'migrate -path /app/migrations -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up' >> /app/entrypoint.sh && \
    echo '' >> /app/entrypoint.sh && \
    echo 'echo "Starting API server..."' >> /app/entrypoint.sh && \
    echo 'exec /app/api' >> /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

# 切換到非 root 使用者
USER appuser

# 暴露 API port
EXPOSE 8080

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 執行啟動腳本
ENTRYPOINT ["/app/entrypoint.sh"]

