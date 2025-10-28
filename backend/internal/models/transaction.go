package models

import (
	"time"

	"github.com/google/uuid"
)

// AssetType 資產類型
type AssetType string

const (
	AssetTypeCash     AssetType = "cash"
	AssetTypeTWStock  AssetType = "tw-stock"
	AssetTypeUSStock  AssetType = "us-stock"
	AssetTypeCrypto   AssetType = "crypto"
)

// Currency 幣別
type Currency string

const (
	CurrencyTWD Currency = "TWD" // 新台幣
	CurrencyUSD Currency = "USD" // 美金
)

// TransactionType 交易類型
type TransactionType string

const (
	TransactionTypeBuy      TransactionType = "buy"
	TransactionTypeSell     TransactionType = "sell"
	TransactionTypeDividend TransactionType = "dividend"
	TransactionTypeFee      TransactionType = "fee"
)

// Transaction 交易記錄模型
type Transaction struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	Date            time.Time       `json:"date" db:"date"`
	AssetType       AssetType       `json:"asset_type" db:"asset_type"`
	Symbol          string          `json:"symbol" db:"symbol"`
	Name            string          `json:"name" db:"name"`
	TransactionType TransactionType `json:"type" db:"transaction_type"`
	Quantity        float64         `json:"quantity" db:"quantity"`
	Price           float64         `json:"price" db:"price"`
	Amount          float64         `json:"amount" db:"amount"`
	Fee             *float64        `json:"fee,omitempty" db:"fee"`
	Tax             *float64        `json:"tax,omitempty" db:"tax"`
	Currency        Currency        `json:"currency" db:"currency"`
	ExchangeRateID  *int            `json:"exchange_rate_id,omitempty" db:"exchange_rate_id"`
	Note            *string         `json:"note,omitempty" db:"note"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateTransactionInput 建立交易的輸入資料
type CreateTransactionInput struct {
	Date            time.Time       `json:"date" binding:"required"`
	AssetType       AssetType       `json:"asset_type" binding:"required"`
	Symbol          string          `json:"symbol" binding:"required"`
	Name            string          `json:"name" binding:"required"`
	TransactionType TransactionType `json:"type" binding:"required"`
	Quantity        float64         `json:"quantity" binding:"required,gte=0"`
	Price           float64         `json:"price" binding:"required,gte=0"`
	Amount          float64         `json:"amount" binding:"required"`
	Fee             *float64        `json:"fee,omitempty" binding:"omitempty,gte=0"`
	Tax             *float64        `json:"tax,omitempty" binding:"omitempty,gte=0"`
	Currency        Currency        `json:"currency" binding:"required"`
	Note            *string         `json:"note,omitempty"`
}

// UpdateTransactionInput 更新交易的輸入資料
type UpdateTransactionInput struct {
	Date            *time.Time       `json:"date,omitempty"`
	AssetType       *AssetType       `json:"asset_type,omitempty"`
	Symbol          *string          `json:"symbol,omitempty"`
	Name            *string          `json:"name,omitempty"`
	TransactionType *TransactionType `json:"type,omitempty"`
	Quantity        *float64         `json:"quantity,omitempty" binding:"omitempty,gte=0"`
	Price           *float64         `json:"price,omitempty" binding:"omitempty,gte=0"`
	Amount          *float64         `json:"amount,omitempty"`
	Fee             *float64         `json:"fee,omitempty" binding:"omitempty,gte=0"`
	Tax             *float64         `json:"tax,omitempty" binding:"omitempty,gte=0"`
	Currency        *Currency        `json:"currency,omitempty"`
	Note            *string          `json:"note,omitempty"`
}

// Validate 驗證 AssetType 是否有效
func (a AssetType) Validate() bool {
	switch a {
	case AssetTypeCash, AssetTypeTWStock, AssetTypeUSStock, AssetTypeCrypto:
		return true
	}
	return false
}

// Validate 驗證 TransactionType 是否有效
func (t TransactionType) Validate() bool {
	switch t {
	case TransactionTypeBuy, TransactionTypeSell, TransactionTypeDividend, TransactionTypeFee:
		return true
	}
	return false
}

// Validate 驗證 Currency 是否有效
func (c Currency) Validate() bool {
	switch c {
	case CurrencyTWD, CurrencyUSD:
		return true
	}
	return false
}

