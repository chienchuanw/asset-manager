package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreditCard 信用卡模型
type CreditCard struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	IssuingBank      string     `json:"issuing_bank" db:"issuing_bank"`
	CardName         string     `json:"card_name" db:"card_name"`
	CardNumberLast4  string     `json:"card_number_last4" db:"card_number_last4"`
	BillingDay       int        `json:"billing_day" db:"billing_day"`
	PaymentDueDay    int        `json:"payment_due_day" db:"payment_due_day"`
	CreditLimit      float64    `json:"credit_limit" db:"credit_limit"`
	UsedCredit       float64    `json:"used_credit" db:"used_credit"`
	GroupID          *uuid.UUID `json:"group_id,omitempty" db:"group_id"`
	Note             *string    `json:"note,omitempty" db:"note"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateCreditCardInput 建立信用卡的輸入資料
type CreateCreditCardInput struct {
	IssuingBank     string  `json:"issuing_bank" binding:"required,max=255"`
	CardName        string  `json:"card_name" binding:"required,max=255"`
	CardNumberLast4 string  `json:"card_number_last4" binding:"required,len=4"`
	BillingDay      int     `json:"billing_day" binding:"required,min=1,max=31"`
	PaymentDueDay   int     `json:"payment_due_day" binding:"required,min=1,max=31"`
	CreditLimit     float64 `json:"credit_limit" binding:"required,gt=0"`
	UsedCredit      float64 `json:"used_credit" binding:"gte=0"`
	Note            *string `json:"note,omitempty" binding:"omitempty,max=1000"`
}

// UpdateCreditCardInput 更新信用卡的輸入資料
type UpdateCreditCardInput struct {
	IssuingBank     *string  `json:"issuing_bank,omitempty" binding:"omitempty,max=255"`
	CardName        *string  `json:"card_name,omitempty" binding:"omitempty,max=255"`
	CardNumberLast4 *string  `json:"card_number_last4,omitempty" binding:"omitempty,len=4"`
	BillingDay      *int     `json:"billing_day,omitempty" binding:"omitempty,min=1,max=31"`
	PaymentDueDay   *int     `json:"payment_due_day,omitempty" binding:"omitempty,min=1,max=31"`
	CreditLimit     *float64 `json:"credit_limit,omitempty" binding:"omitempty,gt=0"`
	UsedCredit      *float64 `json:"used_credit,omitempty" binding:"omitempty,gte=0"`
	Note            *string  `json:"note,omitempty" binding:"omitempty,max=1000"`
}

// Validate 驗證建立信用卡的輸入資料
func (input *CreateCreditCardInput) Validate() error {
	// 驗證卡號後四碼格式（必須是 4 位數字）
	if len(input.CardNumberLast4) != 4 {
		return fmt.Errorf("card_number_last4 must be exactly 4 characters")
	}

	// 驗證帳單日範圍
	if input.BillingDay < 1 || input.BillingDay > 31 {
		return fmt.Errorf("billing_day must be between 1 and 31")
	}

	// 驗證繳款截止日範圍
	if input.PaymentDueDay < 1 || input.PaymentDueDay > 31 {
		return fmt.Errorf("payment_due_day must be between 1 and 31")
	}

	// 驗證信用額度必須大於 0
	if input.CreditLimit <= 0 {
		return fmt.Errorf("credit_limit must be greater than 0")
	}

	// 驗證已使用額度不能為負數
	if input.UsedCredit < 0 {
		return fmt.Errorf("used_credit cannot be negative")
	}

	// 驗證已使用額度不能超過信用額度
	if input.UsedCredit > input.CreditLimit {
		return fmt.Errorf("used_credit cannot exceed credit_limit")
	}

	return nil
}

// Validate 驗證更新信用卡的輸入資料
func (input *UpdateCreditCardInput) Validate() error {
	// 如果有提供卡號後四碼，驗證格式
	if input.CardNumberLast4 != nil && len(*input.CardNumberLast4) != 4 {
		return fmt.Errorf("card_number_last4 must be exactly 4 characters")
	}

	// 如果有提供帳單日，驗證範圍
	if input.BillingDay != nil && (*input.BillingDay < 1 || *input.BillingDay > 31) {
		return fmt.Errorf("billing_day must be between 1 and 31")
	}

	// 如果有提供繳款截止日，驗證範圍
	if input.PaymentDueDay != nil && (*input.PaymentDueDay < 1 || *input.PaymentDueDay > 31) {
		return fmt.Errorf("payment_due_day must be between 1 and 31")
	}

	// 如果有提供信用額度，驗證必須大於 0
	if input.CreditLimit != nil && *input.CreditLimit <= 0 {
		return fmt.Errorf("credit_limit must be greater than 0")
	}

	// 如果有提供已使用額度，驗證不能為負數
	if input.UsedCredit != nil && *input.UsedCredit < 0 {
		return fmt.Errorf("used_credit cannot be negative")
	}

	return nil
}

// AvailableCredit 計算可用額度
func (c *CreditCard) AvailableCredit() float64 {
	return c.CreditLimit - c.UsedCredit
}

// CreditUtilization 計算信用額度使用率（百分比）
func (c *CreditCard) CreditUtilization() float64 {
	if c.CreditLimit == 0 {
		return 0
	}
	return (c.UsedCredit / c.CreditLimit) * 100
}

