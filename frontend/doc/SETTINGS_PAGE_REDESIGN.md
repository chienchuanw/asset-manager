# Settings 頁面重新設計

本文件記錄 Settings 頁面的排版統一修改，使其與整個專案的前端風格保持一致。

---

## 🎯 修改目標

將 Settings 頁面的排版風格統一為與其他頁面（Dashboard、Holdings、Transactions、Analytics）一致的設計。

---

## 📊 專案前端風格規範

### 1. 主容器結構

```tsx
<AppLayout>
  <main className="flex-1 p-4 md:p-6 bg-gray-50">
    <div className="flex flex-col gap-6">
      {/* 頁面內容 */}
    </div>
  </main>
</AppLayout>
```

### 2. 頁面標題

使用 **Card 組件** 包裹標題，而非直接使用 `<h1>` 標籤：

```tsx
<Card>
  <CardHeader>
    <CardTitle className="text-2xl">頁面標題</CardTitle>
    <CardDescription>頁面描述</CardDescription>
  </CardHeader>
</Card>
```

### 3. Loading 狀態

使用統一的 `<Loading />` 組件：

```tsx
if (isLoading) {
  return (
    <AppLayout>
      <main className="flex-1 p-4 md:p-6 bg-gray-50">
        <Loading variant="page" size="lg" text="載入設定中..." />
      </main>
    </AppLayout>
  );
}
```

### 4. 內容間距

- 主要區塊間距：`gap-6`
- 卡片內部間距：`space-y-4`
- 表單欄位間距：`space-y-2`

---

## 🔄 修改內容

### 修改 1：更新 Loading 狀態

**修改前：**

```tsx
if (isLoading) {
  return (
    <AppLayout>
      <div className="flex items-center justify-center h-full">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    </AppLayout>
  );
}
```

**修改後：**

```tsx
if (isLoading) {
  return (
    <AppLayout>
      <main className="flex-1 p-4 md:p-6 bg-gray-50">
        <Loading variant="page" size="lg" text="載入設定中..." />
      </main>
    </AppLayout>
  );
}
```

**改進：**
- ✅ 使用統一的 `<Loading />` 組件
- ✅ 加入 `main` 容器和背景色
- ✅ 提供載入文字說明

---

### 修改 2：統一主容器結構

**修改前：**

```tsx
<AppLayout>
  <div className="flex-1 overflow-auto">
    <div className="container mx-auto py-8 px-4 lg:px-8 max-w-5xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">設定</h1>
        <p className="text-muted-foreground mt-2">管理系統設定和偏好</p>
      </div>
      <div className="space-y-6">
        {/* 內容 */}
      </div>
    </div>
  </div>
</AppLayout>
```

**修改後：**

```tsx
<AppLayout>
  <main className="flex-1 p-4 md:p-6 bg-gray-50">
    <div className="flex flex-col gap-6">
      {/* 頁面標題卡片 */}
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">設定</CardTitle>
          <CardDescription>管理系統設定和偏好</CardDescription>
        </CardHeader>
      </Card>
      
      {/* Discord 設定 */}
      <Card>
        {/* ... */}
      </Card>
      
      {/* 資產配置設定 */}
      <Card>
        {/* ... */}
      </Card>
      
      {/* 操作按鈕 */}
      <div className="flex justify-end gap-4">
        {/* ... */}
      </div>
    </div>
  </main>
</AppLayout>
```

**改進：**
- ✅ 使用 `<main>` 標籤替代 `<div>`
- ✅ 統一背景色 `bg-gray-50`
- ✅ 統一內外邊距 `p-4 md:p-6`
- ✅ 使用 Card 組件包裹標題
- ✅ 移除不必要的嵌套 `div`
- ✅ 統一間距為 `gap-6`

---

### 修改 3：新增 Loading 組件 import

**修改前：**

```tsx
import { Separator } from "@/components/ui/separator";
import { toast } from "sonner";
import { Loader2, Send } from "lucide-react";
```

**修改後：**

```tsx
import { Separator } from "@/components/ui/separator";
import { Loading } from "@/components/ui/loading";
import { toast } from "sonner";
import { Loader2, Send } from "lucide-react";
```

---

## ✅ 修改後的優點

### 1. **視覺一致性**
- 所有頁面使用相同的背景色（`bg-gray-50`）
- 統一的內外邊距和間距
- 一致的 Card 組件使用方式

### 2. **結構清晰**
- 使用語義化的 `<main>` 標籤
- 減少不必要的嵌套層級
- 更清晰的區塊劃分

### 3. **響應式設計**
- 使用 `p-4 md:p-6` 確保在不同螢幕尺寸下的適當間距
- 與其他頁面保持一致的響應式行為

### 4. **Loading 體驗**
- 使用統一的 Loading 組件
- 提供明確的載入狀態文字
- 更好的使用者體驗

---

## 📝 頁面結構對比

### 其他頁面（Dashboard、Holdings、Transactions、Analytics）

```
AppLayout
└── main.flex-1.p-4.md:p-6.bg-gray-50
    └── div.flex.flex-col.gap-6
        ├── Card (標題或統計)
        ├── Card (內容區塊 1)
        ├── Card (內容區塊 2)
        └── ...
```

### Settings 頁面（修改後）

```
AppLayout
└── main.flex-1.p-4.md:p-6.bg-gray-50
    └── div.flex.flex-col.gap-6
        ├── Card (頁面標題)
        ├── Card (Discord 設定)
        ├── Card (資產配置設定)
        └── div (操作按鈕)
```

**✅ 結構完全一致！**

---

## 🎨 設計規範總結

| 項目 | 規範 | Settings 頁面 |
|------|------|--------------|
| 主容器標籤 | `<main>` | ✅ 已修改 |
| 背景色 | `bg-gray-50` | ✅ 已修改 |
| 內外邊距 | `p-4 md:p-6` | ✅ 已修改 |
| 內容容器 | `flex flex-col gap-6` | ✅ 已修改 |
| 頁面標題 | Card 組件 | ✅ 已修改 |
| Loading 狀態 | `<Loading />` 組件 | ✅ 已修改 |
| 區塊間距 | `gap-6` | ✅ 已修改 |

---

## 📸 視覺效果改進

### 修改前
- 使用 `container mx-auto` 導致內容寬度受限
- 使用 `<h1>` 標籤，與其他頁面風格不一致
- Loading 狀態過於簡單，缺少說明文字

### 修改後
- 使用全寬佈局，與其他頁面一致
- 使用 Card 組件包裹標題，視覺更統一
- Loading 狀態提供明確的文字說明

---

## 🚀 後續建議

1. **保持一致性**：未來新增頁面時，請參考此設計規範
2. **組件化**：考慮將頁面標題卡片抽取為共用組件
3. **文檔更新**：將此設計規範加入專案的 UI 設計指南

---

## 📅 更新日誌

- **2025-10-25**: 完成 Settings 頁面排版統一修改
  - 更新主容器結構
  - 統一 Loading 狀態
  - 使用 Card 組件包裹標題
  - 調整間距和佈局

