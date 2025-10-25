# Header 標題重構

本文件記錄將頁面標題從頁面內容移至 AppLayout Header 的重構工作。

---

## 🎯 重構目標

將所有頁面的標題統一顯示在 AppLayout 的 Header 區域，而非在頁面內容中使用獨立的標題卡片。

### 優點

1. **視覺一致性**：所有頁面的標題位置統一
2. **節省空間**：移除重複的標題卡片，釋放更多內容空間
3. **更好的導航體驗**：標題固定在 Header，滾動時始終可見
4. **符合現代設計**：與主流 Web App 的設計模式一致
5. **簡化頁面結構**：減少不必要的嵌套層級

---

## 📝 修改內容

### 1. AppLayout 組件修改

**檔案：** `frontend/src/components/layout/AppLayout.tsx`

#### 1.1 新增 Props

```typescript
interface AppLayoutProps {
  children: React.ReactNode;
  title?: string;        // 新增：頁面標題
  description?: string;  // 新增：頁面描述
}

export function AppLayout({ children, title, description }: AppLayoutProps) {
  // ...
}
```

#### 1.2 修改 Header 區域

**修改前：**

```tsx
<header className="sticky top-0 z-10 flex h-14 items-center gap-4 border-b bg-background/95 backdrop-blur px-4 lg:px-6">
  <SidebarTrigger />
  <Separator orientation="vertical" className="h-6" />
  <div className="flex items-center gap-2">
    <span className="font-semibold">Asset Manager</span>
  </div>
</header>
```

**修改後：**

```tsx
<header className="sticky top-0 z-10 flex h-14 items-center gap-4 border-b bg-background/95 backdrop-blur px-4 lg:px-6">
  <SidebarTrigger />
  <Separator orientation="vertical" className="h-6" />
  
  {/* 動態頁面標題 */}
  {title ? (
    <div className="flex flex-col">
      <span className="font-semibold text-base">{title}</span>
      {description && (
        <span className="text-xs text-muted-foreground hidden sm:block">
          {description}
        </span>
      )}
    </div>
  ) : (
    <div className="flex items-center gap-2">
      <span className="font-semibold">Asset Manager</span>
    </div>
  )}
</header>
```

**改進：**
- ✅ 支援動態顯示頁面標題和描述
- ✅ 描述文字在小螢幕上隱藏（`hidden sm:block`）
- ✅ 未提供標題時顯示預設的 "Asset Manager"

---

### 2. 各頁面修改

#### 2.1 Dashboard 頁面

**檔案：** `frontend/src/app/dashboard/page.tsx`

**修改：**

```tsx
// 所有 <AppLayout> 都加上 title 和 description
<AppLayout title="首頁" description="資產概況總覽">
  {/* 頁面內容 */}
</AppLayout>
```

**影響範圍：**
- Loading 狀態
- Error 狀態
- 正常顯示狀態

---

#### 2.2 Holdings 頁面

**檔案：** `frontend/src/app/holdings/page.tsx`

**修改：**

```tsx
<AppLayout title="持倉明細" description="查看所有資產的詳細持倉資訊">
  {/* 頁面內容 */}
</AppLayout>
```

**影響範圍：**
- Loading 狀態
- Error 狀態
- 正常顯示狀態

---

#### 2.3 Transactions 頁面

**檔案：** `frontend/src/app/transactions/page.tsx`

**修改：**

```tsx
<AppLayout title="交易記錄" description="管理和查看所有交易記錄">
  {/* 頁面內容 */}
</AppLayout>
```

---

#### 2.4 Analytics 頁面

**檔案：** `frontend/src/app/analytics/page.tsx`

**修改：**

```tsx
<AppLayout title="分析報表" description="查看投資績效分析">
  {/* 頁面內容 */}
</AppLayout>
```

---

#### 2.5 Settings 頁面

**檔案：** `frontend/src/app/settings/page.tsx`

**修改前：**

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
      
      {/* 設定卡片 */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* ... */}
      </div>
    </div>
  </main>
</AppLayout>
```

**修改後：**

```tsx
<AppLayout title="設定" description="管理系統設定和偏好">
  <main className="flex-1 p-4 md:p-6 bg-gray-50">
    <div className="flex flex-col gap-6">
      {/* 直接從設定卡片開始，移除標題卡片 */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Discord 設定 */}
        <Card>...</Card>
        
        {/* 資產配置設定 */}
        <Card>...</Card>
      </div>
      
      {/* 操作按鈕 */}
      <div className="flex justify-end gap-4">...</div>
    </div>
  </main>
</AppLayout>
```

**改進：**
- ✅ 移除獨立的標題卡片
- ✅ 標題顯示在 Header 區域
- ✅ 節省垂直空間
- ✅ 視覺更簡潔

---

## 📊 頁面標題對照表

| 頁面路由 | 標題 | 描述 |
|---------|------|------|
| `/dashboard` | 首頁 | 資產概況總覽 |
| `/holdings` | 持倉明細 | 查看所有資產的詳細持倉資訊 |
| `/transactions` | 交易記錄 | 管理和查看所有交易記錄 |
| `/analytics` | 分析報表 | 查看投資績效分析 |
| `/settings` | 設定 | 管理系統設定和偏好 |

---

## 🎨 視覺效果對比

### 修改前

```
┌─────────────────────────────────────┐
│ [☰] Asset Manager                  │ ← Header（固定）
├─────────────────────────────────────┤
│ ┌─────────────────────────────────┐ │
│ │ 設定                            │ │ ← 標題卡片（重複）
│ │ 管理系統設定和偏好              │ │
│ └─────────────────────────────────┘ │
│                                     │
│ ┌───────────┐ ┌───────────┐        │
│ │ Discord   │ │ 資產配置  │        │
│ └───────────┘ └───────────┘        │
└─────────────────────────────────────┘
```

### 修改後

```
┌─────────────────────────────────────┐
│ [☰] 設定                            │ ← Header（包含標題）
│     管理系統設定和偏好              │
├─────────────────────────────────────┤
│ ┌───────────┐ ┌───────────┐        │
│ │ Discord   │ │ 資產配置  │        │ ← 直接顯示內容
│ └───────────┘ └───────────┘        │
│                                     │
│ [重置] [儲存設定]                   │
└─────────────────────────────────────┘
```

---

## ✨ 響應式設計

### 大螢幕（≥ 640px）

- 標題：`font-semibold text-base`
- 描述：顯示（`text-xs text-muted-foreground`）

### 小螢幕（< 640px）

- 標題：`font-semibold text-base`
- 描述：隱藏（`hidden sm:block`）

---

## 🔧 技術細節

### Header 高度

- 固定高度：`h-14`（56px）
- 足夠容納標題和描述的雙行顯示

### 樣式類別

```tsx
// 標題容器
<div className="flex flex-col">
  
// 標題文字
<span className="font-semibold text-base">{title}</span>

// 描述文字（響應式隱藏）
<span className="text-xs text-muted-foreground hidden sm:block">
  {description}
</span>
```

### Sticky Header

```tsx
className="sticky top-0 z-10 ... bg-background/95 backdrop-blur"
```

- `sticky top-0`：固定在頂部
- `z-10`：確保在其他內容之上
- `bg-background/95`：95% 不透明度背景
- `backdrop-blur`：背景模糊效果

---

## ✅ 驗證清單

- [x] AppLayout 組件支援 `title` 和 `description` props
- [x] Header 區域正確顯示動態標題
- [x] Dashboard 頁面已更新
- [x] Holdings 頁面已更新
- [x] Transactions 頁面已更新
- [x] Analytics 頁面已更新
- [x] Settings 頁面已更新並移除標題卡片
- [x] 所有頁面的 Loading 狀態已更新
- [x] 所有頁面的 Error 狀態已更新
- [x] 響應式設計正常運作
- [x] 無 TypeScript 錯誤

---

## 🚀 後續建議

1. **測試各頁面**：
   - 檢查標題和描述是否正確顯示
   - 測試響應式行為（小螢幕隱藏描述）
   - 確認 Sticky Header 滾動效果

2. **考慮未來擴展**：
   - 可在 Header 加入麵包屑導航
   - 可在 Header 加入頁面操作按鈕
   - 可支援自訂 Header 右側內容

3. **保持一致性**：
   - 未來新增頁面時，記得使用相同模式
   - 標題和描述文字保持簡潔明瞭

---

## 📅 更新日誌

- **2025-10-25**: 完成 Header 標題重構
  - 修改 AppLayout 組件支援動態標題
  - 更新所有 5 個頁面使用新的標題模式
  - Settings 頁面移除獨立標題卡片
  - 實現響應式標題顯示

