# Phase 4 實作指南

## 📋 你需要做的事

Phase 4 的程式碼已經全部完成，現在你需要：

### 1. 安裝必要的套件

```bash
cd frontend
pnpm add @radix-ui/react-dialog @radix-ui/react-label @radix-ui/react-slot class-variance-authority
```

這些是 shadcn/ui 元件所需的依賴套件。

---

### 2. 確認後端 API 正在執行

```bash
cd backend
make run
```

後端應該在 `http://localhost:8080` 執行。

---

### 3. 啟動前端開發伺服器

```bash
cd frontend
pnpm dev
```

前端應該在 `http://localhost:3000` 執行。

---

### 4. 訪問交易列表頁面

開啟瀏覽器，訪問：
```
http://localhost:3000/transactions
```

---

## 🎯 測試功能

### 測試 1: 新增交易

1. 點擊「新增交易」按鈕
2. 填寫表單：
   - 日期：選擇今天
   - 資產類型：台股
   - 代碼：2330
   - 名稱：台積電
   - 交易類型：買入
   - 數量：10
   - 價格：620
   - 金額：6200（應該自動計算）
3. 點擊「建立交易」
4. 確認對話框關閉，列表顯示新交易

### 測試 2: 篩選交易

1. 在搜尋框輸入「2330」
2. 確認只顯示台積電的交易
3. 清空搜尋框
4. 使用「交易類型」下拉選單選擇「買入」
5. 確認只顯示買入交易
6. 使用「資產類別」下拉選單選擇「台股」
7. 確認只顯示台股交易

### 測試 3: 刪除交易

1. 點擊任一交易列的刪除按鈕（垃圾桶圖示）
2. 確認刪除提示
3. 確認交易從列表中消失

---

## 🎨 已實作的功能

### ✅ UI 元件
- Dialog（對話框）
- Form（表單）
- Label（標籤）
- Textarea（文字區域）

### ✅ 功能元件
- AddTransactionDialog（新增交易對話框）
  - 表單驗證
  - 自動計算金額
  - 錯誤處理

### ✅ 頁面功能
- 交易列表顯示
- 載入狀態（Skeleton）
- 錯誤處理
- 統計摘要卡片
- 篩選功能（搜尋、交易類型、資產類別）
- 新增交易
- 刪除交易
- 響應式設計

---

## 📚 學習重點

### 1. React Hook Form + Zod

表單驗證的最佳實踐：

```tsx
const form = useForm<CreateTransactionFormData>({
  resolver: zodResolver(createTransactionSchema),
  defaultValues: { ... },
});
```

### 2. React Query Hooks

資料管理的最佳實踐：

```tsx
const { data, isLoading, error, refetch } = useTransactions();
const createMutation = useCreateTransaction({
  onSuccess: () => refetch(),
});
```

### 3. useMemo 效能優化

避免不必要的重新計算：

```tsx
const filteredTransactions = useMemo(() => {
  // 篩選邏輯
}, [transactions, filterType, filterAssetType, searchQuery]);
```

### 4. 條件渲染

根據狀態顯示不同內容：

```tsx
{isLoading ? (
  <Skeleton />
) : error ? (
  <ErrorMessage />
) : data.length === 0 ? (
  <EmptyState />
) : (
  <DataList />
)}
```

---

## 🐛 可能遇到的問題

### 問題 1: 套件安裝失敗

**解決方法**：
```bash
# 清除快取
pnpm store prune

# 重新安裝
pnpm install
```

### 問題 2: CORS 錯誤

**解決方法**：
確認後端有設定 CORS middleware。

### 問題 3: 新增交易後列表沒有更新

**解決方法**：
檢查 `onSuccess` 回調是否正確呼叫 `refetch()`。

---

## 🎯 下一步

Phase 4 完成後，你可以：

1. **測試所有功能**：確保新增、篩選、刪除都正常運作
2. **新增更多測試資料**：測試不同的資產類型和交易類型
3. **學習程式碼**：理解每個元件的實作方式
4. **考慮 Phase 5**：實作編輯功能、匯出功能等進階功能

---

## 📖 相關文件

- `doc/PHASE2_SETUP.md` - Phase 2 基礎建設
- `doc/PHASE3_HOOKS.md` - Phase 3 React Query Hooks 詳細指南
- `doc/PHASE3_COMPLETE.md` - Phase 3 完成總結
- `doc/PHASE4_COMPLETE.md` - Phase 4 完成總結

---

**準備好了嗎？開始測試你的交易列表頁面吧！** 🚀

