# Phase 2: 前端基礎建設 - 設定指南

## ✅ 已完成的工作

### Step 2.1: 環境變數設定 ✅
- ✅ 建立 `.env.local`
- ✅ 建立 `.env.example`
- ✅ 設定 `NEXT_PUBLIC_API_URL=http://localhost:8080`

### Step 2.2: 型別定義 ✅
- ✅ 建立 `src/types/transaction.ts`
  - 定義 `AssetType` 和 `TransactionType` 列舉
  - 定義 `Transaction` 介面（與後端 API 對應）
  - 定義 `CreateTransactionInput` 和 `UpdateTransactionInput`
  - 定義 `TransactionFilters`
  - 定義 `APIResponse` 和 `APIError`
  - 建立 Zod Schema（`createTransactionSchema`, `updateTransactionSchema`）
  - 提供輔助函式（`getAssetTypeLabel`, `getTransactionTypeLabel` 等）

### Step 2.3: API Client ✅
- ✅ 建立 `src/lib/api/client.ts`
  - 實作基礎 fetch wrapper
  - 處理 API 錯誤
  - 支援查詢參數
  - 提供 `apiClient.get()`, `apiClient.post()`, `apiClient.put()`, `apiClient.delete()`
  
- ✅ 建立 `src/lib/api/transactions.ts`
  - 實作 `transactionsAPI.getAll()`
  - 實作 `transactionsAPI.getById()`
  - 實作 `transactionsAPI.create()`
  - 實作 `transactionsAPI.update()`
  - 實作 `transactionsAPI.delete()`

### Step 2.4: React Query 設定 ✅
- ✅ 建立 `src/providers/QueryProvider.tsx`
  - 設定 QueryClient
  - 設定預設選項（staleTime, gcTime, retry 等）
  - 加入 React Query Devtools
  
- ✅ 更新 `src/app/layout.tsx`
  - 加入 `QueryProvider`
  - 更新 metadata（標題和描述）
  - 更新語言為 `zh-TW`

---

## 📦 需要安裝的套件

請執行以下指令安裝必要的套件：

```bash
cd frontend
pnpm add @tanstack/react-query @tanstack/react-query-devtools react-hook-form zod @hookform/resolvers
```

### 套件說明

| 套件 | 版本 | 用途 |
|------|------|------|
| `@tanstack/react-query` | latest | 資料獲取和狀態管理 |
| `@tanstack/react-query-devtools` | latest | React Query 開發工具 |
| `react-hook-form` | latest | 表單管理 |
| `zod` | latest | 資料驗證 |
| `@hookform/resolvers` | latest | react-hook-form 與 zod 的整合 |

---

## 📁 建立的檔案

### 環境變數
- `frontend/.env.local` - 本地環境變數
- `frontend/.env.example` - 環境變數範例

### 型別定義
- `frontend/src/types/transaction.ts` - 交易相關型別和 Zod Schema

### API Client
- `frontend/src/lib/api/client.ts` - 基礎 API Client
- `frontend/src/lib/api/transactions.ts` - 交易 API

### Providers
- `frontend/src/providers/QueryProvider.tsx` - React Query Provider

### 更新的檔案
- `frontend/src/app/layout.tsx` - 加入 QueryProvider

---

## 🎯 檔案結構

```
frontend/
├── .env.local                    # 環境變數
├── .env.example                  # 環境變數範例
├── src/
│   ├── app/
│   │   └── layout.tsx           # 更新：加入 QueryProvider
│   ├── types/
│   │   └── transaction.ts       # 交易型別定義
│   ├── lib/
│   │   └── api/
│   │       ├── client.ts        # 基礎 API Client
│   │       └── transactions.ts  # 交易 API
│   └── providers/
│       └── QueryProvider.tsx    # React Query Provider
```

---

## 🔧 使用方式

### 1. 環境變數

`.env.local` 檔案已設定：

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

如果後端 API 位址不同，請修改此檔案。

### 2. 型別定義

使用交易型別：

```typescript
import type { Transaction, CreateTransactionInput } from "@/types/transaction";
import { AssetType, TransactionType } from "@/types/transaction";

// 建立交易資料
const newTransaction: CreateTransactionInput = {
  date: "2025-10-23T00:00:00Z",
  asset_type: AssetType.TW_STOCK,
  symbol: "2330",
  name: "TSMC",
  type: TransactionType.BUY,
  quantity: 10,
  price: 620,
  amount: 6200,
  fee: 28,
  note: "定期定額",
};
```

### 3. API 呼叫

使用交易 API：

```typescript
import { transactionsAPI } from "@/lib/api/transactions";

// 取得所有交易
const transactions = await transactionsAPI.getAll();

// 取得單筆交易
const transaction = await transactionsAPI.getById("transaction-id");

// 建立交易
const newTransaction = await transactionsAPI.create({
  date: "2025-10-23T00:00:00Z",
  asset_type: "tw-stock",
  symbol: "2330",
  name: "TSMC",
  type: "buy",
  quantity: 10,
  price: 620,
  amount: 6200,
});

// 更新交易
const updatedTransaction = await transactionsAPI.update("transaction-id", {
  quantity: 20,
});

// 刪除交易
await transactionsAPI.delete("transaction-id");
```

### 4. React Query（下一步會實作）

在元件中使用 React Query：

```typescript
import { useQuery, useMutation } from "@tanstack/react-query";
import { transactionsAPI } from "@/lib/api/transactions";

// 取得交易列表
const { data, isLoading, error } = useQuery({
  queryKey: ["transactions"],
  queryFn: () => transactionsAPI.getAll(),
});

// 建立交易
const createMutation = useMutation({
  mutationFn: transactionsAPI.create,
  onSuccess: () => {
    // 重新獲取交易列表
    queryClient.invalidateQueries({ queryKey: ["transactions"] });
  },
});
```

---

## 🎨 Zod Schema 驗證

使用 Zod Schema 進行表單驗證：

```typescript
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { createTransactionSchema } from "@/types/transaction";
import type { CreateTransactionFormData } from "@/types/transaction";

const form = useForm<CreateTransactionFormData>({
  resolver: zodResolver(createTransactionSchema),
  defaultValues: {
    date: new Date().toISOString(),
    asset_type: "tw-stock",
    symbol: "",
    name: "",
    type: "buy",
    quantity: 0,
    price: 0,
    amount: 0,
    fee: null,
    note: null,
  },
});
```

---

## 🐛 錯誤處理

API Client 會自動處理錯誤並拋出 `APIError`：

```typescript
import { APIError } from "@/lib/api/client";

try {
  const transactions = await transactionsAPI.getAll();
} catch (error) {
  if (error instanceof APIError) {
    console.error("API 錯誤:", error.code, error.message);
    // 可以根據 error.code 顯示不同的錯誤訊息
  } else {
    console.error("未知錯誤:", error);
  }
}
```

---

## 📊 React Query Devtools

React Query Devtools 已經加入，在開發環境中會自動顯示。

你可以：
- 查看所有 query 的狀態
- 手動觸發 refetch
- 查看快取資料
- 除錯 query 和 mutation

---

## ✅ 驗證設定

### 1. 檢查環境變數

```bash
# 確認 .env.local 存在
cat frontend/.env.local
```

應該顯示：
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 2. 檢查套件安裝

安裝完套件後，執行：

```bash
cd frontend
pnpm list @tanstack/react-query react-hook-form zod @hookform/resolvers
```

應該顯示已安裝的版本。

### 3. 啟動開發伺服器

```bash
cd frontend
pnpm dev
```

應該可以正常啟動，沒有錯誤。

---

## 🎯 下一步：Phase 3

Phase 2 完成後，接下來是 Phase 3：

### Phase 3: 實作 React Query Hooks

- [ ] 建立 `useTransactions` hook（取得交易列表）
- [ ] 建立 `useTransaction` hook（取得單筆交易）
- [ ] 建立 `useCreateTransaction` hook（建立交易）
- [ ] 建立 `useUpdateTransaction` hook（更新交易）
- [ ] 建立 `useDeleteTransaction` hook（刪除交易）

### Phase 4: 實作交易列表頁面

- [ ] 更新 `src/app/transactions/page.tsx`
- [ ] 使用 `useTransactions` hook 取得資料
- [ ] 實作篩選功能
- [ ] 實作搜尋功能
- [ ] 實作分頁功能

### Phase 5: 實作新增交易功能

- [ ] 建立 `AddTransactionDialog` 元件
- [ ] 使用 `react-hook-form` + `zod` 建立表單
- [ ] 使用 `useCreateTransaction` hook 送出資料
- [ ] 實作表單驗證和錯誤顯示

---

## 📞 需要幫助？

如果遇到問題：

1. 確認所有套件都已安裝
2. 確認 `.env.local` 檔案存在且設定正確
3. 確認後端 API 正在執行（`http://localhost:8080`）
4. 檢查瀏覽器 console 是否有錯誤訊息

---

**Phase 2 基礎建設完成！準備進入 Phase 3！** 🚀

