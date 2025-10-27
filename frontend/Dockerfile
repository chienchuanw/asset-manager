# 多階段建置 Dockerfile for Next.js Frontend
# Stage 1: 依賴安裝階段
FROM node:18-alpine AS deps

# 安裝 pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

# 複製 package 檔案
COPY package.json pnpm-lock.yaml ./

# 安裝依賴套件
RUN pnpm install --frozen-lockfile

# Stage 2: 建置階段
FROM node:18-alpine AS builder

# 安裝 pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

# 從 deps 階段複製 node_modules
COPY --from=deps /app/node_modules ./node_modules

# 複製所有原始碼
COPY . .

# 設定環境變數 (build time)
# NEXT_PUBLIC_ 開頭的變數會被打包進前端程式碼
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

# 建置 Next.js 應用程式
# 使用 standalone 模式減少映像檔大小
RUN pnpm build

# Stage 3: 執行階段
FROM node:18-alpine AS runner

WORKDIR /app

# 設定環境變數
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# 建立非 root 使用者
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# 複製必要檔案
# Next.js standalone 模式會產生最小化的執行環境
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# 切換到非 root 使用者
USER nextjs

# 暴露 port
EXPOSE 3000

# 設定環境變數
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD node -e "require('http').get('http://localhost:3000/api/health', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"

# 啟動應用程式
CMD ["node", "server.js"]

