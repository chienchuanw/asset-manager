package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// FIFOCalculatorResult FIFO 計算結果（包含持倉和警告）
type FIFOCalculatorResult struct {
	Holdings map[string]*models.Holding // 成功計算的持倉
	Warnings []*models.Warning           // 計算過程中的警告
}

// FIFOCalculator FIFO 成本計算器介面
type FIFOCalculator interface {
	// CalculateHoldingForSymbol 計算單一標的的持倉
	CalculateHoldingForSymbol(symbol string, transactions []*models.Transaction) (*models.Holding, error)

	// CalculateAllHoldings 計算所有標的的持倉（返回結果包含警告）
	CalculateAllHoldings(transactions []*models.Transaction) (*FIFOCalculatorResult, error)

	// CalculateCostBasis 計算賣出交易的成本基礎（使用 FIFO 規則）
	CalculateCostBasis(symbol string, sellTransaction *models.Transaction, allTransactions []*models.Transaction) (float64, error)
}

// fifoCalculator FIFO 計算器實作
type fifoCalculator struct {
	exchangeRateService ExchangeRateService
}

// NewFIFOCalculator 建立新的 FIFO 計算器
func NewFIFOCalculator(exchangeRateService ExchangeRateService) FIFOCalculator {
	return &fifoCalculator{
		exchangeRateService: exchangeRateService,
	}
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

// CalculateAllHoldings 計算所有標的的持倉（返回結果包含警告）
func (c *fifoCalculator) CalculateAllHoldings(transactions []*models.Transaction) (*FIFOCalculatorResult, error) {
	holdings := make(map[string]*models.Holding)
	warnings := []*models.Warning{}

	// 取得所有唯一的標的代碼
	symbols := getUniqueSymbols(transactions)

	// 逐個計算每個標的的持倉
	for _, symbol := range symbols {
		holding, err := c.CalculateHoldingForSymbol(symbol, transactions)
		if err != nil {
			// 檢查是否為數量不足錯誤
			if isInsufficientQuantityError(err) {
				// 記錄警告並跳過此標的
				warning := createInsufficientQuantityWarning(symbol, err)
				warnings = append(warnings, warning)
				continue
			}
			// 其他錯誤直接返回
			return nil, fmt.Errorf("failed to calculate holding for %s: %w", symbol, err)
		}

		// 只保留有持倉的標的
		if holding != nil {
			holdings[symbol] = holding
		}
	}

	return &FIFOCalculatorResult{
		Holdings: holdings,
		Warnings: warnings,
	}, nil
}

// processBuy 處理買入交易，建立新的成本批次
func (c *fifoCalculator) processBuy(tx *models.Transaction) (*models.CostBatch, error) {
	// 計算含手續費的總成本（原幣別）
	totalCostOriginal := tx.Amount
	if tx.Fee != nil {
		totalCostOriginal += *tx.Fee
	}

	// 計算原幣別的單位成本
	unitCostOriginal := totalCostOriginal / tx.Quantity

	// 計算 TWD 的總成本和單位成本
	var totalCostTWD float64
	var exchangeRate float64

	if tx.Currency == models.CurrencyUSD {
		// 使用交易當天的匯率轉換
		var err error
		totalCostTWD, err = c.exchangeRateService.ConvertToTWD(totalCostOriginal, tx.Currency, tx.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to convert cost to TWD for %s on %s: %w", tx.Symbol, tx.Date.Format("2006-01-02"), err)
		}
		// 計算匯率（TWD / 原幣別）
		exchangeRate = totalCostTWD / totalCostOriginal
	} else {
		// TWD 交易，直接使用原值
		totalCostTWD = totalCostOriginal
		exchangeRate = 1.0
	}

	unitCostTWD := totalCostTWD / tx.Quantity

	batch := &models.CostBatch{
		Date:             tx.Date,
		Quantity:         tx.Quantity,
		UnitCost:         unitCostTWD,
		UnitCostOriginal: unitCostOriginal,
		OriginalQty:      tx.Quantity,
		Currency:         tx.Currency,
		ExchangeRate:     exchangeRate,
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
	var totalCostTWD float64
	var totalCostOriginal float64

	for _, batch := range batches {
		totalQuantity += batch.Quantity
		totalCostTWD += batch.Quantity * batch.UnitCost
		totalCostOriginal += batch.Quantity * batch.UnitCostOriginal
	}

	avgCostTWD := totalCostTWD / totalQuantity
	avgCostOriginal := totalCostOriginal / totalQuantity

	return &models.Holding{
		Symbol:          symbol,
		Name:            name,
		AssetType:       assetType,
		Quantity:        totalQuantity,
		AvgCost:         avgCostTWD,
		AvgCostOriginal: avgCostOriginal,
		TotalCost:       totalCostTWD,
		LastUpdated:     time.Now(),
	}
}

// CalculateCostBasis 計算賣出交易的成本基礎（使用 FIFO 規則）
func (c *fifoCalculator) CalculateCostBasis(symbol string, sellTransaction *models.Transaction, allTransactions []*models.Transaction) (float64, error) {
	// 驗證賣出交易
	if sellTransaction.TransactionType != models.TransactionTypeSell {
		return 0, fmt.Errorf("transaction is not a sell transaction")
	}

	if sellTransaction.Symbol != symbol {
		return 0, fmt.Errorf("transaction symbol %s does not match requested symbol %s", sellTransaction.Symbol, symbol)
	}

	// 篩選出該標的在賣出交易之前的所有交易
	symbolTransactions := filterTransactionsBeforeSell(allTransactions, symbol, sellTransaction.Date)

	// 按日期排序（FIFO 需要按時間順序處理）
	sort.Slice(symbolTransactions, func(i, j int) bool {
		return symbolTransactions[i].Date.Before(symbolTransactions[j].Date)
	})

	// 建立成本批次
	costBatches := []*models.CostBatch{}

	for _, tx := range symbolTransactions {
		switch tx.TransactionType {
		case models.TransactionTypeBuy:
			batch, err := c.processBuy(tx)
			if err != nil {
				return 0, err
			}
			costBatches = append(costBatches, batch)

		case models.TransactionTypeSell:
			var err error
			costBatches, err = c.processSell(tx, costBatches)
			if err != nil {
				return 0, err
			}
		}
	}

	// 使用 FIFO 規則計算賣出的成本基礎
	costBasis, err := c.calculateCostBasisFromBatches(sellTransaction.Quantity, costBatches)
	if err != nil {
		return 0, err
	}

	return costBasis, nil
}

// calculateCostBasisFromBatches 從成本批次計算賣出的成本基礎
func (c *fifoCalculator) calculateCostBasisFromBatches(sellQuantity float64, batches []*models.CostBatch) (float64, error) {
	remainingToSell := sellQuantity
	totalCostBasis := 0.0

	for _, batch := range batches {
		if remainingToSell <= 0 {
			break
		}

		if batch.Quantity <= remainingToSell {
			// 這個批次全部賣出
			totalCostBasis += batch.Quantity * batch.UnitCost
			remainingToSell -= batch.Quantity
		} else {
			// 這個批次部分賣出
			totalCostBasis += remainingToSell * batch.UnitCost
			remainingToSell = 0
		}
	}

	// 如果還有剩餘要賣的數量，表示賣超了
	if remainingToSell > 0 {
		return 0, fmt.Errorf("insufficient quantity to sell: trying to sell %.2f but only have %.2f available",
			sellQuantity, sellQuantity-remainingToSell)
	}

	return totalCostBasis, nil
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

// filterTransactionsBeforeSell 篩選出賣出交易之前的所有交易
func filterTransactionsBeforeSell(transactions []*models.Transaction, symbol string, sellDate time.Time) []*models.Transaction {
	result := []*models.Transaction{}
	for _, tx := range transactions {
		// 只保留相同標的且在賣出日期之前的交易
		if tx.Symbol == symbol && tx.Date.Before(sellDate) {
			result = append(result, tx)
		}
	}
	return result
}

// isInsufficientQuantityError 檢查錯誤是否為數量不足錯誤
func isInsufficientQuantityError(err error) bool {
	if err == nil {
		return false
	}
	// 檢查錯誤訊息是否包含 "insufficient quantity"
	return fmt.Sprintf("%v", err) != "" &&
		(fmt.Sprintf("%v", err)[:len("insufficient quantity")] == "insufficient quantity" ||
		 fmt.Sprintf("%v", err)[:len("failed to calculate holding")] == "failed to calculate holding")
}

// createInsufficientQuantityWarning 建立數量不足警告
func createInsufficientQuantityWarning(symbol string, err error) *models.Warning {
	// 解析錯誤訊息以取得詳細資訊
	// 錯誤格式: "insufficient quantity to sell: trying to sell X but only have Y"
	errMsg := err.Error()

	var required, available float64
	// 嘗試從錯誤訊息中解析數字
	fmt.Sscanf(errMsg, "insufficient quantity to sell: trying to sell %f but only have %f", &required, &available)

	missing := required - available
	if missing < 0 {
		missing = 0
	}

	return &models.Warning{
		Code:    models.WarningCodeInsufficientQuantity,
		Symbol:  symbol,
		Message: fmt.Sprintf("標的 %s 的交易記錄不完整：嘗試賣出 %.2f 股，但只有 %.2f 股可用", symbol, required, available),
		Details: map[string]interface{}{
			"required":  required,
			"available": available,
			"missing":   missing,
		},
	}
}
