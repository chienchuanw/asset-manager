package models

// CSVTransactionRow 代表 CSV 檔案中的一行交易記錄
type CSVTransactionRow struct {
	Date            string `csv:"date"`             // 交易日期 (YYYY-MM-DD)
	AssetType       string `csv:"asset_type"`       // 資產類別
	Symbol          string `csv:"symbol"`           // 交易標的代碼
	Name            string `csv:"name"`             // 交易標的名稱
	TransactionType string `csv:"transaction_type"` // 交易類型
	Quantity        string `csv:"quantity"`         // 數量
	Price           string `csv:"price"`            // 單價
	Fee             string `csv:"fee"`              // 手續費（選填）
	Tax             string `csv:"tax"`              // 交易稅（選填）
	Currency        string `csv:"currency"`         // 幣別
	Note            string `csv:"note"`             // 備註（選填）
}

// CSVImportResult CSV 匯入結果
type CSVImportResult struct {
	Success      bool                      `json:"success"`
	Transactions []*CreateTransactionInput `json:"transactions,omitempty"`
	Errors       []CSVValidationError      `json:"errors,omitempty"`
}

// CSVValidationError CSV 驗證錯誤
type CSVValidationError struct {
	Row     int    `json:"row"`     // 錯誤發生的行號（從 1 開始，不含 header）
	Field   string `json:"field"`   // 錯誤欄位
	Message string `json:"message"` // 錯誤訊息
}

