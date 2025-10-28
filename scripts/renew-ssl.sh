#!/bin/bash

# SSL 憑證自動更新腳本
# 此腳本會使用 Docker Compose 中的 Certbot 容器來更新 Let's Encrypt SSL 憑證
# 並在更新後重新載入 Nginx 設定

set -e  # 遇到錯誤立即停止

# 設定變數
PROJECT_DIR="/home/ubuntu/asset-manager"
LOG_FILE="/home/ubuntu/scripts/renew-ssl.log"
ENV_FILE=".env.production"

# 記錄開始時間
echo "========================================" >> "$LOG_FILE"
echo "SSL Certificate Renewal Started at $(date)" >> "$LOG_FILE"
echo "========================================" >> "$LOG_FILE"

# 切換到專案目錄
cd "$PROJECT_DIR" || {
    echo "ERROR: Failed to change directory to $PROJECT_DIR" >> "$LOG_FILE"
    exit 1
}

# 執行 Certbot 更新憑證
echo "Running Certbot renew..." >> "$LOG_FILE"
if docker-compose --env-file "$ENV_FILE" run --rm certbot renew >> "$LOG_FILE" 2>&1; then
    echo "Certbot renew completed successfully" >> "$LOG_FILE"
else
    echo "ERROR: Certbot renew failed" >> "$LOG_FILE"
    exit 1
fi

# 重新載入 Nginx 設定
echo "Reloading Nginx configuration..." >> "$LOG_FILE"
if docker-compose --env-file "$ENV_FILE" exec -T nginx nginx -s reload >> "$LOG_FILE" 2>&1; then
    echo "Nginx reloaded successfully" >> "$LOG_FILE"
else
    echo "WARNING: Nginx reload failed, trying to restart..." >> "$LOG_FILE"
    if docker-compose --env-file "$ENV_FILE" restart nginx >> "$LOG_FILE" 2>&1; then
        echo "Nginx restarted successfully" >> "$LOG_FILE"
    else
        echo "ERROR: Nginx restart failed" >> "$LOG_FILE"
        exit 1
    fi
fi

# 記錄完成時間
echo "SSL Certificate Renewal Completed at $(date)" >> "$LOG_FILE"
echo "========================================" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# 顯示憑證資訊
echo "Certificate information:" >> "$LOG_FILE"
docker-compose --env-file "$ENV_FILE" run --rm certbot certificates >> "$LOG_FILE" 2>&1

exit 0

