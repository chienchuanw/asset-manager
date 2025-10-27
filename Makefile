# Makefile for Asset Manager Project
# 簡化常用的 Docker 和部署指令

.PHONY: help build up down restart logs ps clean backup restore deploy

# 顏色定義
GREEN  := \033[0;32m
YELLOW := \033[0;33m
BLUE   := \033[0;34m
NC     := \033[0m # No Color

# 預設目標
.DEFAULT_GOAL := help

# 顯示幫助訊息
help:
	@echo "$(BLUE)Asset Manager - Available Commands:$(NC)"
	@echo ""
	@echo "$(YELLOW)Docker 指令:$(NC)"
	@echo "  $(GREEN)make build$(NC)          - 建置 Docker 映像檔"
	@echo "  $(GREEN)make up$(NC)             - 啟動所有容器"
	@echo "  $(GREEN)make down$(NC)           - 停止並移除所有容器"
	@echo "  $(GREEN)make restart$(NC)        - 重啟所有容器"
	@echo "  $(GREEN)make logs$(NC)           - 查看所有容器日誌"
	@echo "  $(GREEN)make logs-backend$(NC)   - 查看 Backend 日誌"
	@echo "  $(GREEN)make logs-frontend$(NC)  - 查看 Frontend 日誌"
	@echo "  $(GREEN)make ps$(NC)             - 查看容器狀態"
	@echo "  $(GREEN)make clean$(NC)          - 清理未使用的 Docker 資源"
	@echo ""
	@echo "$(YELLOW)備份與還原:$(NC)"
	@echo "  $(GREEN)make backup$(NC)         - 備份資料庫"
	@echo "  $(GREEN)make restore$(NC)        - 還原資料庫 (需指定檔案)"
	@echo "  $(GREEN)make backup-s3$(NC)      - 備份並上傳到 S3"
	@echo ""
	@echo "$(YELLOW)部署:$(NC)"
	@echo "  $(GREEN)make deploy$(NC)         - 執行完整部署流程"
	@echo "  $(GREEN)make health$(NC)         - 檢查服務健康狀態"
	@echo ""
	@echo "$(YELLOW)開發:$(NC)"
	@echo "  $(GREEN)make dev$(NC)            - 啟動開發環境"
	@echo "  $(GREEN)make test$(NC)           - 執行測試"
	@echo ""

# 建置 Docker 映像檔
build:
	@echo "$(BLUE)建置 Docker 映像檔...$(NC)"
	docker-compose build --no-cache

# 啟動所有容器
up:
	@echo "$(BLUE)啟動所有容器...$(NC)"
	docker-compose --env-file .env.production up -d
	@echo "$(GREEN)✓ 容器已啟動$(NC)"
	@echo ""
	@echo "查看日誌: make logs"
	@echo "查看狀態: make ps"

# 停止並移除所有容器
down:
	@echo "$(BLUE)停止所有容器...$(NC)"
	docker-compose down
	@echo "$(GREEN)✓ 容器已停止$(NC)"

# 重啟所有容器
restart:
	@echo "$(BLUE)重啟所有容器...$(NC)"
	docker-compose restart
	@echo "$(GREEN)✓ 容器已重啟$(NC)"

# 查看所有容器日誌
logs:
	docker-compose logs -f

# 查看 Backend 日誌
logs-backend:
	docker-compose logs -f backend

# 查看 Frontend 日誌
logs-frontend:
	docker-compose logs -f frontend

# 查看容器狀態
ps:
	docker-compose ps

# 清理未使用的 Docker 資源
clean:
	@echo "$(BLUE)清理未使用的 Docker 資源...$(NC)"
	docker image prune -f
	docker container prune -f
	@echo "$(GREEN)✓ 清理完成$(NC)"

# 備份資料庫
backup:
	@echo "$(BLUE)備份資料庫...$(NC)"
	bash scripts/backup-db.sh

# 還原資料庫
restore:
	@echo "$(YELLOW)請指定備份檔案:$(NC)"
	@echo "  make restore FILE=/path/to/backup.sql.gz"
	@echo ""
	@echo "可用的備份檔案:"
	@ls -lh /home/ubuntu/backups/*.sql.gz 2>/dev/null || echo "  無備份檔案"

# 備份並上傳到 S3
backup-s3:
	@echo "$(BLUE)備份並上傳到 S3...$(NC)"
	bash scripts/backup-to-s3.sh

# 執行完整部署流程
deploy:
	@echo "$(BLUE)執行部署...$(NC)"
	bash scripts/deploy.sh

# 檢查服務健康狀態
health:
	@echo "$(BLUE)檢查服務健康狀態...$(NC)"
	@echo ""
	@echo "Backend API:"
	@curl -f http://localhost:8080/health && echo " $(GREEN)✓$(NC)" || echo " $(RED)✗$(NC)"
	@echo ""
	@echo "Frontend:"
	@curl -f http://localhost:3000 > /dev/null 2>&1 && echo " $(GREEN)✓$(NC)" || echo " $(RED)✗$(NC)"
	@echo ""
	@echo "Nginx:"
	@curl -f http://localhost/health && echo " $(GREEN)✓$(NC)" || echo " $(RED)✗$(NC)"
	@echo ""

# 啟動開發環境
dev:
	@echo "$(BLUE)啟動開發環境...$(NC)"
	@echo ""
	@echo "Backend:"
	@echo "  cd backend && make run"
	@echo ""
	@echo "Frontend:"
	@echo "  cd frontend && pnpm dev"

# 執行測試
test:
	@echo "$(BLUE)執行測試...$(NC)"
	@echo ""
	@echo "Backend 測試:"
	cd backend && make test
	@echo ""
	@echo "Frontend 測試:"
	@echo "  (尚未設定)"

