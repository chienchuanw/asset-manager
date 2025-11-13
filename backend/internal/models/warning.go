package models

// WarningCode 警告代碼
type WarningCode string

const (
	// WarningCodeInsufficientQuantity 數量不足警告
	WarningCodeInsufficientQuantity WarningCode = "INSUFFICIENT_QUANTITY"
)

// Warning API 警告訊息
type Warning struct {
	Code    WarningCode            `json:"code"`              // 警告代碼
	Symbol  string                 `json:"symbol"`            // 相關標的代碼
	Message string                 `json:"message"`           // 警告訊息
	Details map[string]interface{} `json:"details,omitempty"` // 詳細資訊
}

// InsufficientQuantityDetails 數量不足的詳細資訊
type InsufficientQuantityDetails struct {
	Required  float64 `json:"required"`  // 需要的數量
	Available float64 `json:"available"` // 可用的數量
	Missing   float64 `json:"missing"`   // 缺少的數量
}

