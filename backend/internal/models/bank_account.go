package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BankAccount 銀行帳戶模型
type BankAccount struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	BankName           string    `json:"bank_name" db:"bank_name"`
	AccountType        string    `json:"account_type" db:"account_type"`
	AccountNumberLast4 string    `json:"account_number_last4" db:"account_number_last4"`
	Currency           Currency  `json:"currency" db:"currency"`
	Balance            float64   `json:"balance" db:"balance"`
	Note               *string   `json:"note,omitempty" db:"note"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// CreateBankAccountInput 建立銀行帳戶的輸入資料
type CreateBankAccountInput struct {
	BankName           string   `json:"bank_name" binding:"required,max=255"`
	AccountType        string   `json:"account_type" binding:"required,max=50"`
	AccountNumberLast4 string   `json:"account_number_last4" binding:"required,len=4"`
	Currency           Currency `json:"currency" binding:"required,oneof=TWD USD"`
	Balance            float64  `json:"balance" binding:"gte=0"`
	Note               *string  `json:"note,omitempty" binding:"omitempty,max=1000"`
}

// UpdateBankAccountInput 更新銀行帳戶的輸入資料
type UpdateBankAccountInput struct {
	BankName           *string   `json:"bank_name,omitempty" binding:"omitempty,max=255"`
	AccountType        *string   `json:"account_type,omitempty" binding:"omitempty,max=50"`
	AccountNumberLast4 *string   `json:"account_number_last4,omitempty" binding:"omitempty,len=4"`
	Currency           *Currency `json:"currency,omitempty" binding:"omitempty,oneof=TWD USD"`
	Balance            *float64  `json:"balance,omitempty" binding:"omitempty,gte=0"`
	Note               *string   `json:"note,omitempty" binding:"omitempty,max=1000"`
}

// Validate 驗證建立銀行帳戶的輸入資料
func (input *CreateBankAccountInput) Validate() error {
	// 驗證帳號後四碼格式（必須是 4 位數字）
	if len(input.AccountNumberLast4) != 4 {
		return fmt.Errorf("account_number_last4 must be exactly 4 characters")
	}

	// 驗證餘額不能為負數
	if input.Balance < 0 {
		return fmt.Errorf("balance cannot be negative")
	}

	return nil
}

// Validate 驗證更新銀行帳戶的輸入資料
func (input *UpdateBankAccountInput) Validate() error {
	// 如果有提供帳號後四碼，驗證格式
	if input.AccountNumberLast4 != nil && len(*input.AccountNumberLast4) != 4 {
		return fmt.Errorf("account_number_last4 must be exactly 4 characters")
	}

	// 如果有提供餘額，驗證不能為負數
	if input.Balance != nil && *input.Balance < 0 {
		return fmt.Errorf("balance cannot be negative")
	}

	return nil
}

