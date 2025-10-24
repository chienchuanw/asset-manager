# Analytics Feature - Phase 3: FIFO Calculator Enhancement

## 📋 概述

Phase 3 增強了 FIFO Calculator，新增了 `CalculateCostBasis()` 方法，用於計算賣出交易的成本基礎。這是實現已實現損益自動計算的關鍵功能。

## ✅ 完成項目

### 1. Interface 擴充

**檔案：** `backend/internal/service/fifo_calculator.go`

新增方法：

```go
// CalculateCostBasis 計算賣出交易的成本基礎（使用 FIFO 規則）
CalculateCostBasis(symbol string, sellTransaction *models.Transaction, allTransactions []*models.Transaction) (float64, error)
```

### 2. 實作新增

**核心方法：**

#### `CalculateCostBasis()`

- 驗證賣出交易的有效性
- 篩選賣出交易之前的所有相關交易
- 按日期排序並建立成本批次
- 使用 FIFO 規則計算成本基礎

#### `calculateCostBasisFromBatches()`

- 從成本批次中計算賣出的總成本
- 按 FIFO 順序扣除批次
- 驗證是否有足夠的持倉可賣

#### `filterTransactionsBeforeSell()`

- 篩選出賣出交易之前的所有交易
- 只保留相同標的的交易

### 3. 測試案例

**檔案：** `backend/internal/service/fifo_calculator_test.go`

新增 5 個測試案例：

1. **TestCalculateCostBasis_SingleBatch**
   - 測試從單一批次賣出
   - 驗證成本基礎計算正確

2. **TestCalculateCostBasis_MultipleBatches**
   - 測試跨多個批次賣出
   - 驗證 FIFO 規則正確應用

3. **TestCalculateCostBasis_WithPreviousSell**
   - 測試考慮之前的賣出交易
   - 驗證批次扣除邏輯正確

4. **TestCalculateCostBasis_InsufficientQuantity**
   - 測試賣出數量超過持有
   - 驗證錯誤處理正確

5. **TestCalculateCostBasis_NotSellTransaction**
   - 測試非賣出交易的錯誤處理
   - 驗證輸入驗證正確

## 📊 測試結果

```bash
=== RUN   TestCalculateCostBasis_SingleBatch
--- PASS: TestCalculateCostBasis_SingleBatch (0.00s)
=== RUN   TestCalculateCostBasis_MultipleBatches
--- PASS: TestCalculateCostBasis_MultipleBatches (0.00s)
=== RUN   TestCalculateCostBasis_WithPreviousSell
--- PASS: TestCalculateCostBasis_WithPreviousSell (0.00s)
=== RUN   TestCalculateCostBasis_InsufficientQuantity
--- PASS: TestCalculateCostBasis_InsufficientQuantity (0.00s)
=== RUN   TestCalculateCostBasis_NotSellTransaction
--- PASS: TestCalculateCostBasis_NotSellTransaction (0.00s)
PASS
```

**所有 18 個 FIFO Calculator 測試全部通過！**

## 🔍 實作細節

### FIFO 成本基礎計算邏輯

```go
// 範例：計算賣出 120 股的成本基礎
// 
// 買入記錄：
// - 1/1: 買入 100 股 @ 500.28 (含手續費)
// - 1/5: 買入 50 股 @ 520.30 (含手續費)
//
// 賣出：1/10 賣出 120 股
//
// FIFO 計算：
// - 從第一批賣出 100 股：100 × 500.28 = 50,028
// - 從第二批賣出 20 股：20 × 520.30 = 10,406
// - 總成本基礎 = 60,434
```

### 關鍵特性

1. **時間順序處理**
   - 只考慮賣出交易之前的買入/賣出記錄
   - 按日期排序確保 FIFO 順序正確

2. **批次管理**
   - 重用現有的 `processBuy()` 和 `processSell()` 方法
   - 確保與持倉計算邏輯一致

3. **錯誤處理**
   - 驗證交易類型
   - 驗證標的代碼匹配
   - 檢查持倉是否足夠

## 🎯 使用範例

```go
calculator := NewFIFOCalculator()

// 準備交易記錄
transactions := []*models.Transaction{
    // 買入記錄...
}

// 賣出交易
sellTx := &models.Transaction{
    Date:            time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
    Symbol:          "2330",
    TransactionType: models.TransactionTypeSell,
    Quantity:        100,
    Price:           620,
    Amount:          62000,
    Fee:             ptrFloat64(28),
}

// 計算成本基礎
costBasis, err := calculator.CalculateCostBasis("2330", sellTx, transactions)
if err != nil {
    // 處理錯誤
}

// costBasis 即為賣出的成本基礎
```

## 📝 下一步：Phase 4

Phase 4 將整合 FIFO Calculator 到 Transaction Service，實現：

- 在建立賣出交易時自動計算成本基礎
- 自動建立 `RealizedProfit` 記錄
- 完整的已實現損益追蹤

## 🎓 學習重點

1. **FIFO 演算法**：先進先出的成本計算方法
2. **批次處理**：如何管理和扣除成本批次
3. **時間序列處理**：按時間順序處理交易記錄
4. **錯誤處理**：完善的輸入驗證和邊界檢查
5. **測試驅動開發**：先寫測試，再實作功能

---

**Phase 3 完成時間：** 2025-10-24  
**測試通過率：** 100% (18/18)  
**程式碼覆蓋率：** 完整覆蓋新增功能

