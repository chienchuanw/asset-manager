# 多階段建置 Dockerfile for Go Backend
# Stage 1: 建置階段
FROM golang:1.21-alpine AS builder

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
FROM alpine:latest

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

