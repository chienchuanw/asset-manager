# Backend Architecture

## 📐 系統架構圖

```
┌─────────────────────────────────────────────────────────────┐
│                         Client                              │
│                    (Frontend / curl)                        │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP Request
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Layer                              │
│                  (Gin HTTP Handlers)                        │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         TransactionHandler                           │  │
│  │  - CreateTransaction()                               │  │
│  │  - GetTransaction()                                  │  │
│  │  - ListTransactions()                                │  │
│  │  - UpdateTransaction()                               │  │
│  │  - DeleteTransaction()                               │  │
│  └──────────────────────┬───────────────────────────────┘  │
└─────────────────────────┼───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    Service Layer                            │
│                  (Business Logic)                           │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         TransactionService                           │  │
│  │  - CreateTransaction()    (+ validation)             │  │
│  │  - GetTransaction()                                  │  │
│  │  - ListTransactions()     (+ validation)             │  │
│  │  - UpdateTransaction()    (+ validation)             │  │
│  │  - DeleteTransaction()                               │  │
│  └──────────────────────┬───────────────────────────────┘  │
└─────────────────────────┼───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                  Repository Layer                           │
│                   (Data Access)                             │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │       TransactionRepository                          │  │
│  │  - Create()                                          │  │
│  │  - GetByID()                                         │  │
│  │  - GetAll()           (+ filters)                    │  │
│  │  - Update()                                          │  │
│  │  - Delete()                                          │  │
│  └──────────────────────┬───────────────────────────────┘  │
└─────────────────────────┼───────────────────────────────────┘
                          │ SQL Queries
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                     Database                                │
│                   (PostgreSQL)                              │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         transactions table                           │  │
│  │  - id (UUID, PK)                                     │  │
│  │  - date                                              │  │
│  │  - asset_type                                        │  │
│  │  - symbol                                            │  │
│  │  - name                                              │  │
│  │  - transaction_type                                  │  │
│  │  - quantity                                          │  │
│  │  - price                                             │  │
│  │  - amount                                            │  │
│  │  - fee                                               │  │
│  │  - note                                              │  │
│  │  - created_at                                        │  │
│  │  - updated_at                                        │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 🏗️ 分層架構說明

### 1. API Layer（API 層）
**職責：**
- 處理 HTTP 請求和回應
- 解析請求參數和 JSON body
- 呼叫 Service 層
- 格式化回應（統一的 APIResponse 格式）
- 錯誤處理和 HTTP 狀態碼

**檔案：**
- `internal/api/transaction_handler.go`

**特點：**
- 使用 Gin 框架
- RESTful API 設計
- 統一的錯誤回應格式

---

### 2. Service Layer（服務層）
**職責：**
- 實作業務邏輯
- 資料驗證
- 呼叫 Repository 層
- 協調多個 Repository（未來可能需要）

**檔案：**
- `internal/service/transaction_service.go`

**特點：**
- 獨立於 HTTP 框架
- 可重用的業務邏輯
- 完整的輸入驗證

---

### 3. Repository Layer（資料存取層）
**職責：**
- 資料庫 CRUD 操作
- SQL 查詢建構
- 資料庫錯誤處理
- 資料映射（DB ↔ Model）

**檔案：**
- `internal/repository/transaction_repository.go`

**特點：**
- 使用 `database/sql` 標準庫
- 支援動態查詢建構
- 支援篩選和分頁

---

### 4. Models Layer（模型層）
**職責：**
- 定義資料結構
- 定義輸入/輸出格式
- 基本的驗證方法

**檔案：**
- `internal/models/transaction.go`

**特點：**
- 清晰的型別定義
- JSON 序列化標籤
- 資料庫欄位映射

---

### 5. Discord Bot 模組
**Discord Bot 架構：**
```
Discord 使用者訊息
    │
    ▼
Handler (handleMessage / handleInteraction)
    │
    ├── Parser (Gemini NLP)
    │   └── 解析自然語言 → ParseResult
    │       action: create | query | cc_payment | unsupported | chat
    │
    ├── action="create" → 記帳流程
    │   └── 選擇付款方式 → 預覽 → 確認 → CashFlowServiceAdapter
    │
    ├── action="query" → 查詢流程
    │   └── CashFlowQuerier / AccountBalanceQuerier
    │
    ├── action="cc_payment" → 信用卡繳款流程
    │   └── 選卡 → 選銀行帳戶 → 預覽 → 確認 → CreditCardPaymentAdapter
    │       └── CashFlowService.CreateCashFlow (transfer_out)
    │
    ├── action="chat" → 友善問候回應
    │
    └── action="unsupported" → 不支援提示 + 功能清單
```

**關鍵檔案：**
- `internal/discord/parser.go` -- Gemini NLP 解析器
- `internal/discord/handler.go` -- Discord 訊息和互動處理
- `internal/discord/adapter.go` -- Service 層橋接（CashFlowServiceAdapter, CreditCardPaymentAdapter）
- `internal/discord/i18n.go` -- 雙語訊息（zh-TW / en）

**設計模式：**
- 無狀態互動：Button custom_id 編碼 `action:payload:authorID`
- Pending entries：多步驟互動透過 `pendingEntry` map 追蹤狀態
- Adapter 模式：Handler 不直接依賴 Service 層，透過介面抽象

---

## 🔄 資料流程

### 建立交易記錄的流程

```
1. Client 發送 POST 請求
   ↓
2. TransactionHandler.CreateTransaction()
   - 解析 JSON body
   - 綁定到 CreateTransactionInput
   ↓
3. TransactionService.CreateTransaction()
   - 驗證 asset_type
   - 驗證 transaction_type
   - 驗證 quantity, price, fee
   ↓
4. TransactionRepository.Create()
   - 建構 SQL INSERT 語句
   - 執行查詢
   - 返回新建立的 Transaction
   ↓
5. 回傳結果
   - Service → Handler
   - Handler 格式化為 APIResponse
   - 返回 HTTP 201 Created
```

---

## 🧪 測試策略

### 1. Repository 層測試
**類型：** 整合測試（Integration Tests）

**特點：**
- 需要真實的測試資料庫
- 測試實際的 SQL 查詢
- 使用 testify/suite 管理測試生命週期

**測試內容：**
- CRUD 操作
- 篩選功能
- 錯誤處理

---

### 2. Service 層測試
**類型：** 單元測試（Unit Tests）

**特點：**
- 使用 Mock Repository
- 不需要資料庫
- 快速執行

**測試內容：**
- 業務邏輯
- 資料驗證
- 錯誤處理

---

### 3. API Handler 層測試
**類型：** 單元測試（Unit Tests）

**特點：**
- 使用 Mock Service
- 使用 httptest 模擬 HTTP 請求
- 不需要資料庫

**測試內容：**
- HTTP 請求處理
- 回應格式
- 錯誤處理
- 狀態碼

---

## 📦 依賴注入

我們使用建構函式注入（Constructor Injection）來管理依賴：

```go
// main.go
func main() {
    // 1. 初始化資料庫
    db := db.InitDB()
    
    // 2. 建立 Repository（注入資料庫）
    repo := repository.NewTransactionRepository(db)
    
    // 3. 建立 Service（注入 Repository）
    service := service.NewTransactionService(repo)
    
    // 4. 建立 Handler（注入 Service）
    handler := api.NewTransactionHandler(service)
    
    // 5. 註冊路由
    router.POST("/api/transactions", handler.CreateTransaction)
    // ...
}
```

**優點：**
- 清晰的依賴關係
- 易於測試（可以注入 Mock）
- 易於維護和擴展

---

## 🔐 錯誤處理

### 統一的錯誤回應格式

```json
{
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message"
  }
}
```

### 錯誤碼定義

| 錯誤碼 | 說明 | HTTP 狀態碼 |
|--------|------|-------------|
| `INVALID_INPUT` | 輸入資料格式錯誤 | 400 |
| `INVALID_ID` | ID 格式錯誤 | 400 |
| `INVALID_DATE` | 日期格式錯誤 | 400 |
| `NOT_FOUND` | 資源不存在 | 404 |
| `CREATE_FAILED` | 建立失敗 | 500 |
| `UPDATE_FAILED` | 更新失敗 | 500 |
| `DELETE_FAILED` | 刪除失敗 | 500 |
| `LIST_FAILED` | 查詢失敗 | 500 |

---

## 🚀 擴展性考量

### 未來可能的擴展

1. **快取層** (已實作)
   - 使用 Redis 快取常用查詢
   - 減少資料庫負載

2. **訊息佇列**
   - 使用 RabbitMQ 或 Kafka
   - 處理非同步任務（如價格更新）

3. **微服務拆分**
   - Transaction Service
   - Holdings Service
   - Analytics Service

4. **API 版本控制**
   - `/api/v1/transactions`
   - `/api/v2/transactions`

5. **認證與授權** (已實作)
   - JWT Token
   - Role-based Access Control (RBAC)

6. **Discord Bot 排程繳款提醒**
7. **最低應繳金額自動計算**

---

## 📚 相關文件

- `README_PHASE1.md` - 詳細實作指南
- `PHASE1_SUMMARY.md` - Phase 1 完成總結
- `QUICK_START.md` - 快速開始指南
- `TESTING_GUIDE.md` - 測試指南（含 Discord Bot 測試）

