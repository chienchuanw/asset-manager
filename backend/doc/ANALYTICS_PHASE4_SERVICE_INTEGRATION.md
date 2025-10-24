# Analytics Feature - Phase 4: Transaction Service Integration

## 📋 概述

Phase 4 整合了 Transaction Service 與 FIFO Calculator 和 RealizedProfitRepository，實現在建立賣出交易時自動計算並記錄已實現損益。

## ✅ 完成項目

### 1. Transaction Service 修改

**檔案：** `backend/internal/service/transaction_service.go`

#### 新增依賴

```go
type transactionService struct {
    repo               repository.TransactionRepository
    realizedProfitRepo repository.RealizedProfitRepository
    fifoCalculator     FIFOCalculator
}

func NewTransactionService(
    repo repository.TransactionRepository,
    realizedProfitRepo repository.RealizedProfitRepository,
    fifoCalculator FIFOCalculator,
) TransactionService {
    return &transactionService{
        repo:               repo,
        realizedProfitRepo: realizedProfitRepo,
        fifoCalculator:     fifoCalculator,
    }
}
```

#### 修改 CreateTransaction 方法

```go
func (s *transactionService) CreateTransaction(input *models.CreateTransactionInput) (*models.Transaction, error) {
    // ... 驗證邏輯 ...

    // 建立交易記錄
    transaction, err := s.repo.Create(input)
    if err != nil {
        return nil, err
    }

    // 如果是賣出交易，自動計算並記錄已實現損益
    if input.TransactionType == models.TransactionTypeSell {
        if err := s.createRealizedProfit(transaction); err != nil {
            // 記錄錯誤但不影響交易建立
            fmt.Printf("Warning: failed to create realized profit for transaction %s: %v\n", transaction.ID, err)
        }
    }

    return transaction, nil
}
```

#### 新增 createRealizedProfit 方法

```go
func (s *transactionService) createRealizedProfit(sellTransaction *models.Transaction) error {
    // 取得該標的的所有交易記錄
    filters := repository.TransactionFilters{
        Symbol: &sellTransaction.Symbol,
    }
    allTransactions, err := s.repo.GetAll(filters)
    if err != nil {
        return fmt.Errorf("failed to get transactions for symbol %s: %w", sellTransaction.Symbol, err)
    }

    // 使用 FIFO Calculator 計算成本基礎
    costBasis, err := s.fifoCalculator.CalculateCostBasis(
        sellTransaction.Symbol,
        sellTransaction,
        allTransactions,
    )
    if err != nil {
        return fmt.Errorf("failed to calculate cost basis: %w", err)
    }

    // 準備賣出手續費
    sellFee := 0.0
    if sellTransaction.Fee != nil {
        sellFee = *sellTransaction.Fee
    }

    // 建立已實現損益記錄
    input := &models.CreateRealizedProfitInput{
        TransactionID: sellTransaction.ID.String(),
        Symbol:        sellTransaction.Symbol,
        AssetType:     sellTransaction.AssetType,
        SellDate:      sellTransaction.Date,
        Quantity:      sellTransaction.Quantity,
        SellPrice:     sellTransaction.Price,
        SellAmount:    sellTransaction.Amount,
        SellFee:       sellFee,
        CostBasis:     costBasis,
        Currency:      string(sellTransaction.Currency),
    }

    _, err = s.realizedProfitRepo.Create(input)
    if err != nil {
        return fmt.Errorf("failed to create realized profit record: %w", err)
    }

    return nil
}
```

### 2. 測試更新

**檔案：** `backend/internal/service/transaction_service_test.go`

#### 新增 Mock 實作

```go
// MockRealizedProfitRepository 模擬的 RealizedProfitRepository
type MockRealizedProfitRepository struct {
    mock.Mock
}

// MockFIFOCalculator 模擬的 FIFOCalculator
type MockFIFOCalculator struct {
    mock.Mock
}
```

#### 更新現有測試

所有現有測試都已更新，加入新的依賴：

```go
mockRepo := new(MockTransactionRepository)
mockRealizedProfitRepo := new(MockRealizedProfitRepository)
mockFIFOCalc := new(MockFIFOCalculator)
service := NewTransactionService(mockRepo, mockRealizedProfitRepo, mockFIFOCalc)
```

#### 新增賣出交易測試

```go
func TestCreateTransaction_SellWithRealizedProfit(t *testing.T) {
    // 測試建立賣出交易並自動建立已實現損益
    // ...
}
```

### 3. Main.go 更新

**檔案：** `backend/cmd/api/main.go`

```go
// 初始化 Repository
transactionRepo := repository.NewTransactionRepository(database)
exchangeRateRepo := repository.NewExchangeRateRepository(database)
realizedProfitRepo := repository.NewRealizedProfitRepository(database)

// 初始化 FIFO Calculator（需要在 TransactionService 之前初始化）
fifoCalculator := service.NewFIFOCalculator()

// 初始化 Service
transactionService := service.NewTransactionService(transactionRepo, realizedProfitRepo, fifoCalculator)
```

## 📊 測試結果

```bash
=== RUN   TestCreateTransaction_Success
--- PASS: TestCreateTransaction_Success (0.00s)
=== RUN   TestCreateTransaction_InvalidAssetType
--- PASS: TestCreateTransaction_InvalidAssetType (0.00s)
=== RUN   TestCreateTransaction_InvalidTransactionType
--- PASS: TestCreateTransaction_InvalidTransactionType (0.00s)
=== RUN   TestCreateTransaction_NegativeQuantity
--- PASS: TestCreateTransaction_NegativeQuantity (0.00s)
=== RUN   TestCreateTransaction_SellWithRealizedProfit
--- PASS: TestCreateTransaction_SellWithRealizedProfit (0.00s)
PASS
ok    github.com/chienchuanw/asset-manager/internal/service command-line-arguments  0.341s
```

**✅ 所有測試通過！**

## 🔍 實作細節

### 自動化流程

當使用者建立賣出交易時，系統會自動：

1. **建立交易記錄**

   - 呼叫 `TransactionRepository.Create()`
   - 儲存賣出交易到資料庫

2. **計算成本基礎**

   - 取得該標的的所有交易記錄
   - 使用 FIFO Calculator 計算成本基礎

3. **建立已實現損益記錄**
   - 計算已實現損益：`(sell_amount - sell_fee) - cost_basis`
   - 計算已實現損益百分比：`(realized_pl / cost_basis) × 100`
   - 儲存到 `realized_profits` 表

### 錯誤處理

- 如果計算或記錄已實現損益失敗，會記錄警告訊息
- 不會影響交易記錄的建立（交易仍然成功）
- 未來可考慮使用資料庫事務（Transaction）來確保一致性

### 測試策略

使用 Mock 物件進行單元測試：

- `MockTransactionRepository` - 模擬交易記錄存取
- `MockRealizedProfitRepository` - 模擬已實現損益存取
- `MockFIFOCalculator` - 模擬成本基礎計算

## 🎯 使用範例

### API 請求

```bash
POST /api/transactions
Content-Type: application/json

{
  "date": "2025-10-24",
  "asset_type": "tw_stock",
  "symbol": "2330",
  "name": "台積電",
  "transaction_type": "sell",
  "quantity": 100,
  "price": 620,
  "amount": 62000,
  "fee": 28,
  "currency": "TWD"
}
```

### 自動化結果

1. **建立交易記錄**

   - `transactions` 表新增一筆賣出記錄

2. **自動建立已實現損益**
   - `realized_profits` 表新增一筆記錄
   - 包含成本基礎、已實現損益、損益百分比等資訊

## 📝 下一步：Phase 5

Phase 5 將建立 Analytics Service 和 API，提供：

- 整體投資組合摘要
- 各資產類型績效分析
- 最佳/最差表現資產排行

## 🎓 學習重點

1. **依賴注入**：透過建構函式注入依賴，提高可測試性
2. **自動化業務邏輯**：在適當的時機自動執行相關操作
3. **錯誤處理策略**：區分關鍵錯誤和非關鍵錯誤
4. **Mock 測試**：使用 Mock 物件隔離測試單元
5. **測試驅動開發**：先寫測試，確保功能正確

## ⚠️ 注意事項

### 事務一致性

目前的實作中，如果建立已實現損益失敗，交易記錄仍會保留。未來可考慮：

1. **使用資料庫事務**

   ```go
   tx, err := db.Begin()
   // 建立交易
   // 建立已實現損益
   tx.Commit() // 或 tx.Rollback()
   ```

2. **補償機制**

   - 定期掃描沒有對應已實現損益的賣出交易
   - 自動補建缺失的記錄

3. **事件驅動架構**
   - 發送事件到訊息佇列
   - 非同步處理已實現損益計算

---

**Phase 4 完成時間：** 2025-10-24  
**測試通過率：** 100% (5/5)  
**編譯狀態：** ✅ 成功
