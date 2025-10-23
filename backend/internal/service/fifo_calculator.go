package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// FIFOCalculator FIFO 成本計算器介面
type FIFOCalculator interface {
	// CalculateHoldingForSymbol 計算單一標的的持倉
	CalculateHoldingForSymbol(symbol string, transactions []*models.Transaction) (*models.Holding, error)

	// CalculateAllHoldings 計算所有標的的持倉
	CalculateAllHoldings(transactions []*models.Transaction) (map[string]*models.Holding, error)
}

// fifoCalculator FIFO 計算器實作
type fifoCalculator struct{}

// NewFIFOCalculator 建立新的 FIFO 計算器
func NewFIFOCalculator() FIFOCalculator {
	return &fifoCalculator{}
}

// CalculateHoldingForSymbol 計算單一標的的持倉
func (c *fifoCalculator) CalculateHoldingForSymbol(symbol string, transactions []*models.Transaction) (*models.Holding, error) {
	// 篩選出該標的的交易記錄
	symbolTransactions := filterTransactionsBySymbol(transactions, symbol)

	if len(symbolTransactions) == 0 {
		return nil, nil
	}

	// 按日期排序（FIFO 需要按時間順序處理）
	sort.Slice(symbolTransactions, func(i, j int) bool {
		return symbolTransactions[i].Date.Before(symbolTransactions[j].Date)
	})

	// 初始化成本批次列表
	costBatches := []*models.CostBatch{}

	// 記錄標的資訊
	var assetType models.AssetType
	var name string

	// 逐筆處理交易記錄
	for _, tx := range symbolTransactions {
		assetType = tx.AssetType
		name = tx.Name

		switch tx.TransactionType {
		case models.TransactionTypeBuy:
			// 買入：新增成本批次
			batch, err := c.processBuy(tx)
			if err != nil {
				return nil, err
			}
			costBatches = append(costBatches, batch)

		case models.TransactionTypeSell:
			// 賣出：使用 FIFO 扣除成本批次
			var err error
			costBatches, err = c.processSell(tx, costBatches)
			if err != nil {
				return nil, err
			}

		case models.TransactionTypeDividend:
			// 股利：不影響持倉成本，跳過
			continue

		case models.TransactionTypeFee:
			// 單獨的手續費記錄：暫時跳過（手續費已在買賣時處理）
			continue
		}
	}

	// 如果所有批次都賣完了，返回 nil
	if len(costBatches) == 0 {
		return nil, nil
	}

	// 計算總持倉和平均成本
	holding := c.calculateHoldingFromBatches(symbol, name, assetType, costBatches)

	return holding, nil
}

// CalculateAllHoldings 計算所有標的的持倉
func (c *fifoCalculator) CalculateAllHoldings(transactions []*models.Transaction) (map[string]*models.Holding, error) {
	holdings := make(map[string]*models.Holding)

	// 取得所有唯一的標的代碼
	symbols := getUniqueSymbols(transactions)

	// 逐個計算每個標的的持倉
	for _, symbol := range symbols {
		holding, err := c.CalculateHoldingForSymbol(symbol, transactions)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate holding for %s: %w", symbol, err)
		}

		// 只保留有持倉的標的
		if holding != nil {
			holdings[symbol] = holding
		}
	}

	return holdings, nil
}

// processBuy 處理買入交易，建立新的成本批次
func (c *fifoCalculator) processBuy(tx *models.Transaction) (*models.CostBatch, error) {
	// 計算含手續費的單位成本
	totalCost := tx.Amount
	if tx.Fee != nil {
		totalCost += *tx.Fee
	}

	unitCost := totalCost / tx.Quantity

	batch := &models.CostBatch{
		Date:        tx.Date,
		Quantity:    tx.Quantity,
		UnitCost:    unitCost,
		OriginalQty: tx.Quantity,
	}

	return batch, nil
}

// processSell 處理賣出交易，使用 FIFO 扣除成本批次
func (c *fifoCalculator) processSell(tx *models.Transaction, batches []*models.CostBatch) ([]*models.CostBatch, error) {
	remainingToSell := tx.Quantity
	newBatches := []*models.CostBatch{}

	for _, batch := range batches {
		if remainingToSell <= 0 {
			// 已經賣完，保留剩餘批次
			newBatches = append(newBatches, batch)
			continue
		}

		if batch.Quantity <= remainingToSell {
			// 這個批次全部賣出
			remainingToSell -= batch.Quantity
			// 不加入 newBatches（批次已清空）
		} else {
			// 這個批次部分賣出
			batch.Quantity -= remainingToSell
			remainingToSell = 0
			newBatches = append(newBatches, batch)
		}
	}

	// 如果還有剩餘要賣的數量，表示賣超了
	if remainingToSell > 0 {
		return nil, fmt.Errorf("insufficient quantity to sell: trying to sell %.2f but only have %.2f",
			tx.Quantity, tx.Quantity-remainingToSell)
	}

	return newBatches, nil
}

// calculateHoldingFromBatches 從成本批次計算持倉資訊
func (c *fifoCalculator) calculateHoldingFromBatches(symbol, name string, assetType models.AssetType, batches []*models.CostBatch) *models.Holding {
	var totalQuantity float64
	var totalCost float64

	for _, batch := range batches {
		totalQuantity += batch.Quantity
		totalCost += batch.Quantity * batch.UnitCost
	}

	avgCost := totalCost / totalQuantity

	return &models.Holding{
		Symbol:      symbol,
		Name:        name,
		AssetType:   assetType,
		Quantity:    totalQuantity,
		AvgCost:     avgCost,
		TotalCost:   totalCost,
		LastUpdated: time.Now(),
	}
}

// ==================== 輔助函式 ====================

// filterTransactionsBySymbol 篩選出特定標的的交易記錄
func filterTransactionsBySymbol(transactions []*models.Transaction, symbol string) []*models.Transaction {
	result := []*models.Transaction{}
	for _, tx := range transactions {
		if tx.Symbol == symbol {
			result = append(result, tx)
		}
	}
	return result
}

// getUniqueSymbols 取得所有唯一的標的代碼
func getUniqueSymbols(transactions []*models.Transaction) []string {
	symbolMap := make(map[string]bool)
	for _, tx := range transactions {
		symbolMap[tx.Symbol] = true
	}

	symbols := make([]string, 0, len(symbolMap))
	for symbol := range symbolMap {
		symbols = append(symbols, symbol)
	}

	return symbols
}

