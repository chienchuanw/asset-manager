package models

import "time"

// ExchangeRate 匯率模型
type ExchangeRate struct {
	ID           int       `json:"id" db:"id"`
	FromCurrency Currency  `json:"from_currency" db:"from_currency"`
	ToCurrency   Currency  `json:"to_currency" db:"to_currency"`
	Rate         float64   `json:"rate" db:"rate"`
	Date         time.Time `json:"date" db:"date"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ExchangeRateInput 建立匯率的輸入
type ExchangeRateInput struct {
	FromCurrency Currency  `json:"from_currency" binding:"required"`
	ToCurrency   Currency  `json:"to_currency" binding:"required"`
	Rate         float64   `json:"rate" binding:"required,gt=0"`
	Date         time.Time `json:"date" binding:"required"`
}

