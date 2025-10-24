# Analytics Feature - Phase 6-7: Frontend Integration

## 📋 概述

Phase 6-7 完成了前端整合，建立 API Client、Hooks 和更新 Analytics 頁面，實現與後端 API 的完整對接。

## ✅ 完成項目

### Phase 6: 前端 API Client

#### 1. Analytics Types

**檔案：** `frontend/src/types/analytics.ts`

```typescript
// 時間範圍
export const TimeRange = {
  WEEK: "week",
  MONTH: "month",
  QUARTER: "quarter",
  YEAR: "year",
  ALL: "all",
} as const;

// 分析摘要
export interface AnalyticsSummary {
  total_realized_pl: number;
  total_realized_pl_pct: number;
  total_cost_basis: number;
  total_sell_amount: number;
  total_sell_fee: number;
  transaction_count: number;
  currency: string;
  time_range: string;
  start_date: string;
  end_date: string;
}

// 績效資料
export interface PerformanceData {
  asset_type: AssetType;
  name: string;
  realized_pl: number;
  realized_pl_pct: number;
  cost_basis: number;
  sell_amount: number;
  transaction_count: number;
}

// 最佳/最差表現資產
export interface TopAsset {
  symbol: string;
  name: string;
  asset_type: AssetType;
  realized_pl: number;
  realized_pl_pct: number;
  cost_basis: number;
  sell_amount: number;
}
```

**輔助函式：**

- `getTimeRangeLabel()` - 取得時間範圍顯示名稱
- `getTimeRangeOptions()` - 取得所有時間範圍選項
- `formatCurrency()` - 格式化金額
- `formatPercentage()` - 格式化百分比
- `isPositive()` - 判斷是否為正值

#### 2. Analytics API Client

**檔案：** `frontend/src/lib/api/analytics.ts`

```typescript
export const analyticsAPI = {
  // 取得分析摘要
  getSummary: async (
    timeRange: TimeRange = "month"
  ): Promise<AnalyticsSummary> => {
    return apiClient.get<AnalyticsSummary>("/api/analytics/summary", {
      params: { time_range: timeRange },
    });
  },

  // 取得各資產類型績效
  getPerformance: async (
    timeRange: TimeRange = "month"
  ): Promise<PerformanceData[]> => {
    return apiClient.get<PerformanceData[]>("/api/analytics/performance", {
      params: { time_range: timeRange },
    });
  },

  // 取得最佳/最差表現資產
  getTopAssets: async (
    timeRange: TimeRange = "month",
    limit: number = 5
  ): Promise<TopAsset[]> => {
    return apiClient.get<TopAsset[]>("/api/analytics/top-assets", {
      params: {
        time_range: timeRange,
        limit,
      },
    });
  },
};
```

### Phase 7: 前端 Hooks & 頁面

#### 3. Analytics Hooks

**檔案：** `frontend/src/hooks/useAnalytics.ts`

```typescript
// Query Keys
export const analyticsKeys = {
  all: ["analytics"] as const,
  summary: (timeRange: TimeRange) =>
    [...analyticsKeys.all, "summary", timeRange] as const,
  performance: (timeRange: TimeRange) =>
    [...analyticsKeys.all, "performance", timeRange] as const,
  topAssets: (timeRange: TimeRange, limit: number) =>
    [...analyticsKeys.all, "top-assets", timeRange, limit] as const,
};

// 取得分析摘要
export function useAnalyticsSummary(
  timeRange: TimeRange = "month",
  options?: UseQueryOptions<AnalyticsSummary, APIError>
);

// 取得各資產類型績效
export function useAnalyticsPerformance(
  timeRange: TimeRange = "month",
  options?: UseQueryOptions<PerformanceData[], APIError>
);

// 取得最佳/最差表現資產
export function useAnalyticsTopAssets(
  timeRange: TimeRange = "month",
  limit: number = 5,
  options?: UseQueryOptions<TopAsset[], APIError>
);

// 一次取得所有分析資料
export function useAnalytics(
  timeRange: TimeRange = "month",
  topAssetsLimit: number = 5
);
```

**特點：**

- 使用 React Query 進行資料快取
- `staleTime` 設定為 5 分鐘
- 提供統一的 `useAnalytics` Hook

#### 4. Analytics 頁面更新

**檔案：** `frontend/src/app/analytics/page.tsx`

**主要變更：**

1. **移除 Mock 資料依賴**

   - 移除 `mockAssetAllocation`, `mockPerformanceData` 等
   - 使用 `useAnalytics` Hook 取得真實資料

2. **加入 Loading 狀態**

   ```tsx
   {
     isLoading && (
       <div className="flex items-center justify-center py-12">
         <Loader2 className="h-8 w-8 animate-spin" />
         <span>載入中...</span>
       </div>
     );
   }
   ```

3. **加入 Error 處理**

   ```tsx
   {
     isError && (
       <Alert variant="destructive">
         <AlertCircle className="h-4 w-4" />
         <AlertTitle>載入失敗</AlertTitle>
         <AlertDescription>
           {error?.message || "無法載入分析資料，請稍後再試"}
         </AlertDescription>
       </Alert>
     );
   }
   ```

4. **更新摘要卡片**

   - 總成本基礎
   - 總賣出金額
   - 已實現損益
   - 已實現報酬率

5. **更新圖表**

   - 各資產類型績效長條圖（使用 `performance.data`）
   - 各資產類型損益統計

6. **更新 Top 資產表格**
   - 顯示代碼/名稱、類別、成本、賣出金額、損益、報酬率
   - 使用 `topAssets.data`

---

## 🎯 功能特點

### 1. 時間範圍切換

使用者可以選擇不同的時間範圍：

- 本週
- 本月
- 本季
- 本年
- 全部

切換時會自動重新取得資料。

### 2. 即時資料更新

- 使用 React Query 自動管理快取
- 5 分鐘內不重新取得資料
- 可手動重新整理

### 3. 錯誤處理

- 網路錯誤顯示友善訊息
- API 錯誤顯示具體錯誤訊息
- 提供重試機制

### 4. Loading 狀態

- 顯示 Loading 動畫
- 避免閃爍（使用 React Query 的 `staleTime`）

### 5. 空資料處理

- 當沒有資料時顯示友善訊息
- 不會顯示空白頁面

---

## 📊 頁面結構

```bash
Analytics Page
├── 時間範圍選擇 (Tabs)
├── Loading 狀態
├── Error 狀態
└── 資料顯示
    ├── 績效摘要卡片 (4 張卡片)
    │   ├── 總成本基礎
    │   ├── 總賣出金額
    │   ├── 已實現損益
    │   └── 已實現報酬率
    ├── 圖表區域
    │   ├── 各資產類型績效長條圖
    │   └── 各資產類型損益統計
    └── Top 資產表格
```

---

## 🔍 實作細節

### React Query 設定

```typescript
{
  queryKey: analyticsKeys.summary(timeRange),
  queryFn: () => analyticsAPI.getSummary(timeRange),
  staleTime: 1000 * 60 * 5, // 5 分鐘內不重新取得
}
```

### 條件渲染

```typescript
{!isLoading && !isError && summary.data && (
  // 顯示資料
)}
```

### 動態樣式

```typescript
className={`text-2xl font-bold tabular-nums ${
  isPositive(summary.data.total_realized_pl)
    ? "text-green-600"
    : "text-red-600"
}`}
```

---

## 📝 Phase 6-7 學習重點

1. **TypeScript 型別定義**：建立完整的型別系統
2. **API Client 設計**：統一的 API 呼叫介面
3. **React Query**：資料快取和狀態管理
4. **Custom Hooks**：封裝資料取得邏輯
5. **錯誤處理**：完善的錯誤處理機制
6. **Loading 狀態**：良好的使用者體驗
7. **條件渲染**：根據狀態顯示不同內容
8. **動態樣式**：根據資料動態調整樣式

---

## 🎯 下一步建議

1. **測試前端功能**

   - 啟動前端開發伺服器
   - 測試時間範圍切換
   - 測試 Loading 和 Error 狀態

2. **優化使用者體驗**

   - 加入骨架屏（Skeleton）
   - 加入動畫效果
   - 優化行動裝置顯示

3. **加入更多功能**
   - 匯出報表（CSV/PDF）
   - 自訂時間範圍
   - 更多圖表類型

---

**Phase 6-7 完成時間：** 2025-10-24  
**前端整合狀態：** ✅ 完成  
**與後端對接：** ✅ 成功
