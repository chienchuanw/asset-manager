# Phase 3: React Query Hooks - 實作指南

## ✅ 已完成的工作

### Step 3.1: 建立 React Query Hooks ✅

- ✅ 建立 `src/hooks/useTransactions.ts`
  - `useTransactions()` - 取得交易列表
  - `useTransaction()` - 取得單筆交易
  - `useCreateTransaction()` - 建立交易
  - `useUpdateTransaction()` - 更新交易
  - `useDeleteTransaction()` - 刪除交易
  - `useCreateTransactionOptimistic()` - 建立交易（樂觀更新）
  - `transactionKeys` - Query Keys 管理

- ✅ 建立 `src/hooks/index.ts`
  - 統一匯出所有 hooks

---

## 📁 建立的檔案

```
frontend/src/hooks/
├── useTransactions.ts    # 交易相關 hooks
└── index.ts              # 統一匯出
```

---

## 🎯 Hooks 使用指南

### 1. useTransactions - 取得交易列表

用於顯示交易列表頁面。

```tsx
import { useTransactions } from "@/hooks";

function TransactionsPage() {
  const { data, isLoading, error } = useTransactions();

  if (isLoading) return <div>載入中...</div>;
  if (error) return <div>錯誤: {error.message}</div>;

  return (
    <div>
      {data?.map((transaction) => (
        <div key={transaction.id}>{transaction.name}</div>
      ))}
    </div>
  );
}
```

#### 使用篩選條件

```tsx
const { data } = useTransactions({
  asset_type: "tw-stock",
  type: "buy",
  limit: 10,
  offset: 0,
});
```

#### 自訂選項

```tsx
const { data } = useTransactions(
  { asset_type: "tw-stock" },
  {
    staleTime: 10 * 60 * 1000, // 10 分鐘
    refetchInterval: 30 * 1000, // 每 30 秒自動重新獲取
  }
);
```

---

### 2. useTransaction - 取得單筆交易

用於顯示交易詳情頁面。

```tsx
import { useTransaction } from "@/hooks";

function TransactionDetailPage({ id }: { id: string }) {
  const { data, isLoading, error } = useTransaction(id);

  if (isLoading) return <div>載入中...</div>;
  if (error) return <div>錯誤: {error.message}</div>;
  if (!data) return <div>找不到交易</div>;

  return (
    <div>
      <h1>{data.name}</h1>
      <p>數量: {data.quantity}</p>
      <p>價格: {data.price}</p>
    </div>
  );
}
```

---

### 3. useCreateTransaction - 建立交易

用於新增交易表單。

```tsx
import { useCreateTransaction } from "@/hooks";
import { toast } from "sonner"; // 或其他 toast 套件

function AddTransactionForm() {
  const createMutation = useCreateTransaction({
    onSuccess: () => {
      toast.success("交易建立成功");
      // 關閉對話框或重置表單
    },
    onError: (error) => {
      toast.error(`建立失敗: ${error.message}`);
    },
  });

  const handleSubmit = (data: CreateTransactionInput) => {
    createMutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* 表單欄位 */}
      <button
        type="submit"
        disabled={createMutation.isPending}
      >
        {createMutation.isPending ? "建立中..." : "建立交易"}
      </button>
    </form>
  );
}
```

---

### 4. useUpdateTransaction - 更新交易

用於編輯交易表單。

```tsx
import { useUpdateTransaction } from "@/hooks";

function EditTransactionForm({ id }: { id: string }) {
  const updateMutation = useUpdateTransaction({
    onSuccess: () => {
      toast.success("交易更新成功");
    },
    onError: (error) => {
      toast.error(`更新失敗: ${error.message}`);
    },
  });

  const handleSubmit = (data: UpdateTransactionInput) => {
    updateMutation.mutate({ id, data });
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* 表單欄位 */}
      <button
        type="submit"
        disabled={updateMutation.isPending}
      >
        {updateMutation.isPending ? "更新中..." : "更新交易"}
      </button>
    </form>
  );
}
```

---

### 5. useDeleteTransaction - 刪除交易

用於刪除交易功能。

```tsx
import { useDeleteTransaction } from "@/hooks";

function DeleteTransactionButton({ id }: { id: string }) {
  const deleteMutation = useDeleteTransaction({
    onSuccess: () => {
      toast.success("交易刪除成功");
    },
    onError: (error) => {
      toast.error(`刪除失敗: ${error.message}`);
    },
  });

  const handleDelete = () => {
    if (confirm("確定要刪除這筆交易嗎？")) {
      deleteMutation.mutate(id);
    }
  };

  return (
    <button
      onClick={handleDelete}
      disabled={deleteMutation.isPending}
    >
      {deleteMutation.isPending ? "刪除中..." : "刪除"}
    </button>
  );
}
```

---

### 6. useCreateTransactionOptimistic - 樂觀更新

提供更好的使用者體驗，在伺服器回應之前先更新 UI。

```tsx
import { useCreateTransactionOptimistic } from "@/hooks";

function AddTransactionFormOptimistic() {
  const createMutation = useCreateTransactionOptimistic({
    onSuccess: () => {
      toast.success("交易建立成功");
    },
    onError: (error) => {
      toast.error(`建立失敗: ${error.message}`);
    },
  });

  const handleSubmit = (data: CreateTransactionInput) => {
    // UI 會立即更新，不需要等待伺服器回應
    createMutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* 表單欄位 */}
      <button type="submit">建立交易</button>
    </form>
  );
}
```

---

## 🔑 Query Keys 管理

`transactionKeys` 提供統一的 query key 管理：

```tsx
import { transactionKeys } from "@/hooks";

// 所有交易相關的 queries
transactionKeys.all; // ["transactions"]

// 所有交易列表
transactionKeys.lists(); // ["transactions", "list"]

// 特定篩選條件的交易列表
transactionKeys.list({ asset_type: "tw-stock" });
// ["transactions", "list", { asset_type: "tw-stock" }]

// 所有交易詳情
transactionKeys.details(); // ["transactions", "detail"]

// 特定交易的詳情
transactionKeys.detail("transaction-id");
// ["transactions", "detail", "transaction-id"]
```

### 手動操作快取

```tsx
import { useQueryClient } from "@tanstack/react-query";
import { transactionKeys } from "@/hooks";

function SomeComponent() {
  const queryClient = useQueryClient();

  // 使所有交易列表失效
  const invalidateAllLists = () => {
    queryClient.invalidateQueries({
      queryKey: transactionKeys.lists(),
    });
  };

  // 使特定交易失效
  const invalidateTransaction = (id: string) => {
    queryClient.invalidateQueries({
      queryKey: transactionKeys.detail(id),
    });
  };

  // 手動設定快取資料
  const setTransactionData = (id: string, data: Transaction) => {
    queryClient.setQueryData(transactionKeys.detail(id), data);
  };

  // 取得快取資料
  const getTransactionData = (id: string) => {
    return queryClient.getQueryData<Transaction>(
      transactionKeys.detail(id)
    );
  };

  return <div>...</div>;
}
```

---

## 🎨 進階使用

### 1. 組合多個 Hooks

```tsx
function TransactionManager({ id }: { id: string }) {
  const { data: transaction } = useTransaction(id);
  const updateMutation = useUpdateTransaction();
  const deleteMutation = useDeleteTransaction();

  const handleUpdate = (data: UpdateTransactionInput) => {
    updateMutation.mutate({ id, data });
  };

  const handleDelete = () => {
    deleteMutation.mutate(id);
  };

  return (
    <div>
      <h1>{transaction?.name}</h1>
      <button onClick={() => handleUpdate({ quantity: 20 })}>
        更新數量
      </button>
      <button onClick={handleDelete}>刪除</button>
    </div>
  );
}
```

---

### 2. 使用 React Hook Form

```tsx
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { createTransactionSchema } from "@/types/transaction";
import { useCreateTransaction } from "@/hooks";

function AddTransactionForm() {
  const createMutation = useCreateTransaction();

  const form = useForm({
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
    },
  });

  const onSubmit = form.handleSubmit((data) => {
    createMutation.mutate(data, {
      onSuccess: () => {
        form.reset();
      },
    });
  });

  return (
    <form onSubmit={onSubmit}>
      {/* 表單欄位 */}
    </form>
  );
}
```

---

### 3. 條件式查詢

```tsx
function ConditionalQuery({ shouldFetch }: { shouldFetch: boolean }) {
  const { data } = useTransactions(
    undefined,
    {
      enabled: shouldFetch, // 只有當 shouldFetch 為 true 時才執行查詢
    }
  );

  return <div>...</div>;
}
```

---

### 4. 分頁

```tsx
import { useState } from "react";

function PaginatedTransactions() {
  const [page, setPage] = useState(0);
  const limit = 10;

  const { data, isLoading } = useTransactions({
    limit,
    offset: page * limit,
  });

  return (
    <div>
      {/* 顯示交易列表 */}
      <button onClick={() => setPage(page - 1)} disabled={page === 0}>
        上一頁
      </button>
      <button onClick={() => setPage(page + 1)}>下一頁</button>
    </div>
  );
}
```

---

## 🔄 自動快取管理

所有 mutation hooks 都會自動管理快取：

### 建立交易後
- ✅ 自動使所有交易列表失效
- ✅ 觸發重新獲取資料

### 更新交易後
- ✅ 自動使所有交易列表失效
- ✅ 自動使該筆交易的詳情失效
- ✅ 觸發重新獲取資料

### 刪除交易後
- ✅ 自動使所有交易列表失效
- ✅ 自動移除該筆交易的快取
- ✅ 觸發重新獲取資料

---

## 🐛 錯誤處理

所有 hooks 都使用 `APIError` 型別：

```tsx
import { APIError } from "@/lib/api/client";

function TransactionsPage() {
  const { data, error } = useTransactions();

  if (error) {
    // error 是 APIError 型別
    console.error("錯誤代碼:", error.code);
    console.error("錯誤訊息:", error.message);
    console.error("HTTP 狀態:", error.status);

    // 根據錯誤代碼顯示不同訊息
    if (error.code === "NETWORK_ERROR") {
      return <div>網路連線失敗，請檢查網路設定</div>;
    }

    return <div>發生錯誤: {error.message}</div>;
  }

  return <div>...</div>;
}
```

---

## 📊 載入狀態

React Query 提供多種載入狀態：

```tsx
function TransactionsPage() {
  const {
    data,
    isLoading,      // 首次載入
    isFetching,     // 任何時候的資料獲取（包含背景重新獲取）
    isRefetching,   // 背景重新獲取
    isError,        // 是否有錯誤
    error,          // 錯誤物件
  } = useTransactions();

  if (isLoading) return <div>首次載入中...</div>;
  if (isError) return <div>錯誤: {error.message}</div>;

  return (
    <div>
      {isFetching && <div>更新中...</div>}
      {/* 顯示資料 */}
    </div>
  );
}
```

---

## 🎯 最佳實踐

### 1. 使用 onSuccess 和 onError

```tsx
const createMutation = useCreateTransaction({
  onSuccess: (data) => {
    toast.success("建立成功");
    console.log("新建立的交易:", data);
  },
  onError: (error) => {
    toast.error(error.message);
    console.error("建立失敗:", error);
  },
});
```

### 2. 檢查 isPending 狀態

```tsx
<button
  onClick={() => createMutation.mutate(data)}
  disabled={createMutation.isPending}
>
  {createMutation.isPending ? "建立中..." : "建立交易"}
</button>
```

### 3. 使用樂觀更新提升體驗

對於不太可能失敗的操作（如建立交易），使用樂觀更新：

```tsx
const createMutation = useCreateTransactionOptimistic();
```

### 4. 適當設定 staleTime

```tsx
// 資料不常變動，可以設定較長的 staleTime
const { data } = useTransactions(undefined, {
  staleTime: 10 * 60 * 1000, // 10 分鐘
});
```

---

## ✅ Phase 3 完成檢查清單

- [x] 建立 `useTransactions` hook
- [x] 建立 `useTransaction` hook
- [x] 建立 `useCreateTransaction` hook
- [x] 建立 `useUpdateTransaction` hook
- [x] 建立 `useDeleteTransaction` hook
- [x] 建立 `useCreateTransactionOptimistic` hook
- [x] 建立 `transactionKeys` 管理
- [x] 建立 hooks index 檔案
- [x] 自動快取管理
- [x] 錯誤處理
- [x] TypeScript 型別安全

---

## 🚀 下一步：Phase 4

Phase 3 完成後，接下來是 Phase 4：

### Phase 4: 更新交易列表頁面

- [ ] 更新 `src/app/transactions/page.tsx`
- [ ] 使用 `useTransactions` hook 取得資料
- [ ] 實作交易列表 UI
- [ ] 實作篩選功能
- [ ] 實作搜尋功能
- [ ] 實作新增交易對話框

---

**Phase 3 完成！所有 React Query Hooks 已建立完成！** 🎉

