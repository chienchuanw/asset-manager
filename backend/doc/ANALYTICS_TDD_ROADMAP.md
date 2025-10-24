# Analytics 功能 TDD 實作路線圖

## 📋 總覽

本文檔記錄 Analytics 功能的完整 TDD 實作計畫，包含 5 個 Phase。

---

## 🎯 功能需求

### 已確認的技術決策

1. **已實現損益記錄方式：** 方案 B - 建立獨立的 `realized_profits` 表
2. **計算時機：** 方案 A - 在建立賣出交易時即時計算
3. **手續費處理：** 方案 A - 手續費計入成本基礎
4. **時間範圍定義：** 交易發生時間（本月 = 本月發生的所有交易）

### API 設計

```
GET /api/analytics/summary?time_range=month
GET /api/analytics/performance?time_range=month
GET /api/analytics/top-assets?time_range=month&limit=5
```

---

## 📊 Phase 進度追蹤

| Phase | 名稱 | 狀態 | 完成時間 |
|-------|------|------|---------|
| Phase 1 | 資料庫 Migration | ✅ 完成 | 2025-10-24 |
| Phase 2 | Model & Repository | 🔄 進行中 | - |
| Phase 3 | FIFO Calculator 增強 | ⏳ 待開始 | - |
| Phase 4 | Transaction Service 整合 | ⏳ 待開始 | - |
| Phase 5 | Analytics Service & API | ⏳ 待開始 | - |

---

## 🗂️ Phase 1: 資料庫 Migration ✅

### 目標
建立 `realized_profits` 表，用於記錄已實現損益

### 完成項目
- ✅ 建立 Migration 檔案
  - `000004_create_realized_profits_table.up.sql`
  - `000004_create_realized_profits_table.down.sql`
- ✅ 執行開發環境 Migration
- ✅ 執行測試環境 Migration
- ✅ 驗證表結構正確

### 詳細文檔
參見：`backend/doc/ANALYTICS_PHASE1_MIGRATION.md`

---

## 🏗️ Phase 2: Model & Repository 🔄

### 目標
建立 `RealizedProfit` Model 和 Repository，並通過測試

### 待完成項目

#### 2.1 建立 Model
- [ ] `backend/internal/models/realized_profit.go`
  - [ ] `RealizedProfit` 結構
  - [ ] `CreateRealizedProfitInput` 結構
  - [ ] `RealizedProfitFilters` 結構

#### 2.2 建立 Repository Interface
- [ ] `backend/internal/repository/realized_profit_repository.go`
  - [ ] `RealizedProfitRepository` interface
  - [ ] `Create()` 方法
  - [ ] `GetByTransactionID()` 方法
  - [ ] `GetAll()` 方法
  - [ ] `Delete()` 方法

#### 2.3 撰寫測試（Red）
- [ ] `backend/internal/repository/realized_profit_repository_test.go`
  - [ ] `TestRealizedProfitRepository_Create`
  - [ ] `TestRealizedProfitRepository_GetByTransactionID`
  - [ ] `TestRealizedProfitRepository_GetByTransactionID_NotFound`
  - [ ] `TestRealizedProfitRepository_GetAll`
  - [ ] `TestRealizedProfitRepository_Delete`

#### 2.4 實作 Repository（Green）
- [ ] 實作 `realizedProfitRepository` 結構
- [ ] 實作所有 CRUD 方法
- [ ] 確保所有測試通過

#### 2.5 執行測試
```bash
cd backend
go test ./internal/repository -v -run TestRealizedProfit
```

### 測試檔案
已建立：`backend/internal/repository/realized_profit_repository_test.go`

---

## 🔧 Phase 3: FIFO Calculator 增強

### 目標
修改 FIFO Calculator，在賣出時返回成本基礎資訊

### 待完成項目

#### 3.1 修改測試（Red）
- [ ] `backend/internal/service/fifo_calculator_test.go`
  - [ ] `TestFIFOCalculator_CalculateWithCostBasis`

#### 3.2 修改 Interface
- [ ] `backend/internal/service/fifo_calculator.go`
  - [ ] 新增 `CalculateCostBasis()` 方法

#### 3.3 實作方法（Green）
- [ ] 實作 `CalculateCostBasis()` 邏輯
- [ ] 使用 FIFO 計算成本基礎
- [ ] 處理部分賣出情況

#### 3.4 執行測試
```bash
cd backend
go test ./internal/service -v -run TestFIFOCalculator_CalculateWithCostBasis
```

---

## 🎯 Phase 4: Transaction Service 整合

### 目標
在建立賣出交易時，自動計算並記錄已實現損益

### 待完成項目

#### 4.1 修改測試（Red）
- [ ] `backend/internal/service/transaction_service_test.go`
  - [ ] `TestTransactionService_CreateSellTransaction_WithRealizedProfit`

#### 4.2 修改 Service（Green）
- [ ] 修改 `TransactionService` 結構
  - [ ] 新增 `realizedProfitRepo` 欄位
  - [ ] 新增 `fifoCalculator` 欄位
- [ ] 修改 `NewTransactionService()` 建構函式
- [ ] 修改 `Create()` 方法
  - [ ] 偵測賣出交易
  - [ ] 呼叫 `createRealizedProfit()`
- [ ] 實作 `createRealizedProfit()` 方法
  - [ ] 取得該標的所有交易
  - [ ] 計算成本基礎
  - [ ] 建立已實現損益記錄

#### 4.3 更新 main.go
- [ ] `backend/cmd/api/main.go`
  - [ ] 初始化 `RealizedProfitRepository`
  - [ ] 更新 `TransactionService` 初始化

#### 4.4 執行測試
```bash
cd backend
go test ./internal/service -v -run TestTransactionService_CreateSellTransaction
```

---

## 📊 Phase 5: Analytics Service & API

### 目標
建立 Analytics Service 和 API Handler，提供分析報表資料

### 待完成項目

#### 5.1 建立 Analytics Service 測試（Red）
- [ ] `backend/internal/service/analytics_service_test.go`
  - [ ] `TestAnalyticsService_GetSummary`
  - [ ] `TestAnalyticsService_GetPerformance`
  - [ ] `TestAnalyticsService_GetTopAssets`

#### 5.2 建立 Analytics Models
- [ ] `backend/internal/models/analytics.go`
  - [ ] `AnalyticsSummary` 結構
  - [ ] `PerformanceData` 結構
  - [ ] `TopAsset` 結構
  - [ ] `TimeRange` 類型

#### 5.3 實作 Analytics Service（Green）
- [ ] `backend/internal/service/analytics_service.go`
  - [ ] `AnalyticsService` interface
  - [ ] `GetSummary()` 方法
  - [ ] `GetPerformance()` 方法
  - [ ] `GetTopAssets()` 方法

#### 5.4 建立 Analytics API Handler 測試（Red）
- [ ] `backend/internal/api/analytics_handler_test.go`
  - [ ] `TestAnalyticsHandler_GetSummary`
  - [ ] `TestAnalyticsHandler_GetPerformance`
  - [ ] `TestAnalyticsHandler_GetTopAssets`

#### 5.5 實作 Analytics API Handler（Green）
- [ ] `backend/internal/api/analytics_handler.go`
  - [ ] `AnalyticsHandler` 結構
  - [ ] `GetSummary()` 方法
  - [ ] `GetPerformance()` 方法
  - [ ] `GetTopAssets()` 方法

#### 5.6 註冊路由
- [ ] `backend/cmd/api/main.go`
  - [ ] 初始化 `AnalyticsService`
  - [ ] 初始化 `AnalyticsHandler`
  - [ ] 註冊 `/api/analytics/*` 路由

#### 5.7 執行測試
```bash
cd backend
go test ./internal/service -v -run TestAnalytics
go test ./internal/api -v -run TestAnalytics
```

---

## 🧪 測試策略

### TDD 循環

每個 Phase 都遵循 **Red → Green → Refactor** 循環：

1. **Red（紅燈）**
   - 先寫測試
   - 執行測試，確認失敗（因為功能尚未實作）

2. **Green（綠燈）**
   - 實作最小可行的程式碼
   - 執行測試，確認通過

3. **Refactor（重構）**
   - 優化程式碼
   - 確保測試仍然通過

### 測試覆蓋率目標

- Repository 層：> 80%
- Service 層：> 80%
- API Handler 層：> 70%

---

## 📝 開發指令

### 執行所有測試
```bash
cd backend
make test
```

### 執行特定測試
```bash
# Repository 測試
go test ./internal/repository -v -run TestRealizedProfit

# Service 測試
go test ./internal/service -v -run TestAnalytics

# API Handler 測試
go test ./internal/api -v -run TestAnalytics
```

### 查看測試覆蓋率
```bash
make test-coverage
open coverage.html
```

---

## 🚀 前端整合（Phase 6-7）

待後端完成後，將進行前端整合：

### Phase 6: 前端 API Client
- [ ] `frontend/src/lib/api/analytics.ts`
- [ ] `frontend/src/types/analytics.ts`

### Phase 7: 前端 Hooks & 頁面
- [ ] `frontend/src/hooks/useAnalytics.ts`
- [ ] 更新 `frontend/src/app/analytics/page.tsx`
- [ ] 移除 Mock 資料依賴
- [ ] 加入 Loading 和錯誤處理

---

## 📚 相關文檔

- [Phase 1 Migration 詳細文檔](./ANALYTICS_PHASE1_MIGRATION.md)
- [專案架構文檔](./ARCHITECTURE.md)
- [測試指南](./TESTING_GUIDE.md)

---

## 🎉 總結

目前進度：**Phase 1 完成 ✅**

下一步：**開始 Phase 2 - Model & Repository 實作**

請按照 TDD 原則，先建立測試檔案，再實作功能。

