# ✅ Phase 3: React Query Hooks - 完成！

## 🎉 恭喜！Phase 3 已完成

所有 React Query Hooks 都已建立完成，現在你可以在元件中輕鬆使用這些 hooks 來管理交易資料。

---

## 📦 建立的檔案

### Hooks
- ✅ `src/hooks/useTransactions.ts` - 所有交易相關的 hooks
- ✅ `src/hooks/index.ts` - 統一匯出

### 範例
- ✅ `src/components/examples/TransactionExample.tsx` - 使用範例

### 文件
- ✅ `doc/PHASE3_HOOKS.md` - 完整使用指南

---

## 🎯 可用的 Hooks

| Hook | 用途 | 回傳型別 |
|------|------|---------|
| `useTransactions()` | 取得交易列表 | `UseQueryResult<Transaction[]>` |
| `useTransaction(id)` | 取得單筆交易 | `UseQueryResult<Transaction>` |
| `useCreateTransaction()` | 建立交易 | `UseMutationResult<Transaction>` |
| `useUpdateTransaction()` | 更新交易 | `UseMutationResult<Transaction>` |
| `useDeleteTransaction()` | 刪除交易 | `UseMutationResult<void>` |
| `useCreateTransactionOptimistic()` | 建立交易（樂觀更新）| `UseMutationResult<Transaction>` |

---

## 🚀 快速開始

### 1. 在元件中使用

```tsx
import { useTransactions } from "@/hooks";

export function MyComponent() {
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

### 2. 建立交易

```tsx
import { useCreateTransaction } from "@/hooks";

export function AddTransactionButton() {
  const createMutation = useCreateTransaction({
    onSuccess: () => alert("建立成功"),
  });

  const handleClick = () => {
    createMutation.mutate({
      date: new Date().toISOString(),
      asset_type: "tw-stock",
      symbol: "2330",
      name: "台積電",
      type: "buy",
      quantity: 10,
      price: 620,
      amount: 6200,
    });
  };

  return (
    <button onClick={handleClick} disabled={createMutation.isPending}>
      {createMutation.isPending ? "建立中..." : "建立交易"}
    </button>
  );
}
```

---

## 🎨 特色功能

### 1. 自動快取管理 ✅
- 建立/更新/刪除後自動更新快取
- 不需要手動重新獲取資料

### 2. 樂觀更新 ✅
- `useCreateTransactionOptimistic` 提供即時 UI 更新
- 錯誤時自動回滾

### 3. 型別安全 ✅
- 完整的 TypeScript 型別支援
- 自動補全和型別檢查

### 4. 錯誤處理 ✅
- 統一的 `APIError` 型別
- 清楚的錯誤訊息

### 5. Query Keys 管理 ✅
- `transactionKeys` 提供統一的 key 管理
- 方便手動操作快取

---

## 📚 詳細文件

請參考 `doc/PHASE3_HOOKS.md` 獲取：
- 完整的使用範例
- 進階功能說明
- 最佳實踐建議
- 錯誤處理指南

---

## 🧪 測試範例

我已經建立了一個測試範例元件：`src/components/examples/TransactionExample.tsx`

你可以在任何頁面中使用它來測試 hooks：

```tsx
import { TransactionExample } from "@/components/examples/TransactionExample";

export default function TestPage() {
  return <TransactionExample />;
}
```

---

## ✅ Phase 3 檢查清單

- [x] 建立 `useTransactions` hook（取得列表）
- [x] 建立 `useTransaction` hook（取得單筆）
- [x] 建立 `useCreateTransaction` hook（建立）
- [x] 建立 `useUpdateTransaction` hook（更新）
- [x] 建立 `useDeleteTransaction` hook（刪除）
- [x] 建立 `useCreateTransactionOptimistic` hook（樂觀更新）
- [x] 建立 `transactionKeys` 管理
- [x] 自動快取管理
- [x] 錯誤處理
- [x] TypeScript 型別安全
- [x] 建立使用範例
- [x] 建立完整文件

---

## 🎯 下一步：Phase 4

現在基礎建設都完成了，接下來我們要實作真實的 UI：

### Phase 4: 更新交易列表頁面

**目標**：將 `src/app/transactions/page.tsx` 改造成功能完整的交易列表頁面

**功能**：
- [ ] 使用 `useTransactions` hook 顯示交易列表
- [ ] 實作交易列表 Table UI
- [ ] 實作篩選功能（資產類型、交易類型）
- [ ] 實作搜尋功能（代碼、名稱）
- [ ] 實作新增交易對話框
- [ ] 實作編輯交易功能
- [ ] 實作刪除交易功能
- [ ] 實作載入狀態和錯誤處理

---

## 💡 提示

### 確認後端 API 正在執行

在測試前端功能之前，請確保後端 API 正在執行：

```bash
cd backend
make run
```

後端應該在 `http://localhost:8080` 執行。

### 啟動前端開發伺服器

```bash
cd frontend
pnpm dev
```

前端應該在 `http://localhost:3000` 執行。

### 測試 API 連線

你可以使用測試範例元件來測試 API 連線是否正常。

---

## 🐛 常見問題

### Q: 為什麼會出現 CORS 錯誤？

**A**: 確保後端 API 有設定 CORS。在 Go 後端加入 CORS middleware：

```go
import "github.com/gin-contrib/cors"

router.Use(cors.Default())
```

### Q: 為什麼資料沒有自動更新？

**A**: 檢查是否正確使用了 mutation hooks 的 `onSuccess` 回調，或確認 `invalidateQueries` 有被呼叫。

### Q: 如何除錯 React Query？

**A**: 使用 React Query Devtools（已經在 `QueryProvider` 中加入）：
- 開啟瀏覽器開發者工具
- 點擊右下角的 React Query 圖示
- 查看所有 queries 和 mutations 的狀態

---

**Phase 3 完成！準備進入 Phase 4！** 🚀

