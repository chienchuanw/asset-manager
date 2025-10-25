# Loading 狀態使用指南

本文件說明專案中統一的 Loading 狀態呈現方式。

## 設計原則

- **簡約風格**：使用 `Loader2` icon 搭配 `animate-spin`
- **一致性**：統一使用 shadcn/ui 的設計語言
- **場景適配**：根據不同使用場景選擇合適的 Loading 組件

---

## Loading 組件

### 位置

`frontend/src/components/ui/loading.tsx`

### Props

```typescript
interface LoadingProps {
  variant?: "page" | "inline" | "overlay"  // 顯示變體
  size?: "sm" | "md" | "lg"                // 圖示大小
  text?: string                            // 顯示文字
  className?: string                       // 自訂樣式
}
```

### 使用範例

#### 1. 全頁面 Loading

適用於頁面初次載入時的全屏 Loading 狀態。

```typescript
import { Loading } from "@/components/ui/loading";

// 在頁面組件中
if (isLoading) {
  return (
    <AppLayout>
      <main className="flex-1 p-4 md:p-6 bg-gray-50">
        <Loading variant="page" size="lg" text="載入持倉資料中..." />
      </main>
    </AppLayout>
  );
}
```

**實際使用：**

- `holdings/page.tsx` - 持倉列表頁面
- 其他需要全頁面 Loading 的頁面

#### 2. Inline Loading

適用於區塊內容載入、Tab 切換載入等場景。

```typescript
import { Loading } from "@/components/ui/loading";

// 在組件中
{isLoading && (
  <Loading variant="inline" size="md" text="載入中..." />
)}
```

**實際使用：**

- `analytics/page.tsx` - Tab 內容載入

#### 3. Overlay Loading

適用於需要覆蓋在現有內容上的 Loading 狀態（帶背景模糊效果）。

```typescript
import { Loading } from "@/components/ui/loading";

// 在組件中
<div className="relative">
  {/* 原有內容 */}
  <YourContent />
  
  {/* Loading 覆蓋層 */}
  {isUpdating && (
    <Loading variant="overlay" size="md" text="更新中..." />
  )}
</div>
```

---

## 其他 Loading 場景

### Skeleton Loading

適用於**內容佔位**，提供更好的使用者體驗。

**使用場景：**

- 表格初次載入
- 統計卡片初次載入
- 列表初次載入

**範例：**

```typescript
import { Skeleton } from "@/components/ui/skeleton";

// 統計卡片 Loading
<Card>
  <CardHeader className="pb-2">
    <Skeleton className="h-4 w-24" />
  </CardHeader>
  <CardContent>
    <Skeleton className="h-8 w-32 mb-2" />
    <Skeleton className="h-3 w-20" />
  </CardContent>
</Card>

// 表格 Loading
<TableBody>
  {Array.from({ length: 5 }).map((_, index) => (
    <TableRow key={index}>
      <TableCell>
        <Skeleton className="h-4 w-20" />
      </TableCell>
      {/* 更多 cells... */}
    </TableRow>
  ))}
</TableBody>
```

**實際使用：**

- `dashboard/page.tsx` - 統計卡片、圖表
- `transactions/page.tsx` - 交易記錄表格

### 按鈕 Loading

適用於按鈕操作中的 Loading 狀態。

**範例：**

```typescript
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";

<Button disabled={isPending}>
  {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
  {isPending ? "處理中..." : "提交"}
</Button>
```

**實際使用：**

- `transactions/page.tsx` - 刪除按鈕
- 表單提交按鈕
- 其他操作按鈕

### 重新整理按鈕 Loading

適用於重新整理按鈕的 Loading 狀態。

**範例：**

```typescript
import { RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";

<Button
  variant="outline"
  size="sm"
  onClick={() => refetch()}
  disabled={isFetching}
>
  <RefreshCw
    className={`mr-2 h-4 w-4 ${isFetching ? "animate-spin" : ""}`}
  />
  重新整理
</Button>
```

**實際使用：**

- `holdings/page.tsx` - 重新整理按鈕

---

## 使用決策樹

```
需要顯示 Loading 狀態？
│
├─ 是全頁面初次載入？
│  └─ 是 → 使用 <Loading variant="page" />
│
├─ 是區塊/Tab 內容載入？
│  └─ 是 → 使用 <Loading variant="inline" />
│
├─ 需要覆蓋在現有內容上？
│  └─ 是 → 使用 <Loading variant="overlay" />
│
├─ 是表格/列表/卡片初次載入？
│  └─ 是 → 使用 <Skeleton />
│
├─ 是按鈕操作中？
│  └─ 是 → 使用 <Loader2 className="animate-spin" />
│
└─ 是背景更新/重新整理？
   └─ 是 → 使用 Toast 或按鈕上的小 spinner
```

---

## 統一規範

### ✅ 推薦做法

1. **統一使用 `Loader2` icon** 作為 Loading 圖示
2. **全頁面 Loading** 使用 `<Loading variant="page" />`
3. **內容佔位** 使用 `<Skeleton />`
4. **按鈕操作** 使用 `<Loader2 className="animate-spin" />`
5. **提供文字說明** 讓使用者知道正在載入什麼

### ❌ 避免做法

1. ~~不要使用 `RefreshCw` 作為 Loading 圖示~~（只用於重新整理按鈕）
2. ~~不要混用不同的 Loading 圖示~~
3. ~~不要在全頁面 Loading 時使用 Skeleton~~
4. ~~不要在按鈕操作時使用全頁面 Loading~~

---

## 範例頁面

- ✅ `holdings/page.tsx` - 全頁面 Loading
- ✅ `analytics/page.tsx` - Inline Loading
- ✅ `dashboard/page.tsx` - Skeleton Loading
- ✅ `transactions/page.tsx` - Skeleton + 按鈕 Loading

---

## 更新日誌

- **2025-10-25**: 建立 Loading 組件，統一全專案 Loading 風格

