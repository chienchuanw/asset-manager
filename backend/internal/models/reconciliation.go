package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ReconcileBankAccountInput 校準銀行帳戶餘額的輸入
type ReconcileBankAccountInput struct {
	NewBalance float64   `json:"new_balance" binding:"gte=0"`
	Date       time.Time `json:"date" binding:"required"`
	Note       *string   `json:"note,omitempty" binding:"omitempty,max=500"`
}

// ReconcileCreditCardInput 校準信用卡 used_credit / credit_limit 的輸入
type ReconcileCreditCardInput struct {
	NewUsedCredit  *float64  `json:"new_used_credit,omitempty" binding:"omitempty,gte=0"`
	NewCreditLimit *float64  `json:"new_credit_limit,omitempty" binding:"omitempty,gt=0"`
	Date           time.Time `json:"date" binding:"required"`
	Note           *string   `json:"note,omitempty" binding:"omitempty,max=500"`
}

// ReconcileBatchItem 批次校準的單一項目
type ReconcileBatchItem struct {
	TargetType     SourceType `json:"target_type" binding:"required"`
	TargetID       uuid.UUID  `json:"target_id" binding:"required"`
	NewBalance     *float64   `json:"new_balance,omitempty"`
	NewUsedCredit  *float64   `json:"new_used_credit,omitempty"`
	NewCreditLimit *float64   `json:"new_credit_limit,omitempty"`
}

// ReconcileBatchInput 批次校準的整體輸入
type ReconcileBatchInput struct {
	Date  time.Time            `json:"date" binding:"required"`
	Note  *string              `json:"note,omitempty" binding:"omitempty,max=500"`
	Items []ReconcileBatchItem `json:"items" binding:"required,min=1,dive"`
}

// ReconcileResult 校準操作的結果
type ReconcileResult struct {
	TargetType     SourceType `json:"target_type"`
	TargetID       uuid.UUID  `json:"target_id"`
	Delta          float64    `json:"delta"`
	CashFlowID     *uuid.UUID `json:"cash_flow_id,omitempty"`
	NewBalance     *float64   `json:"new_balance,omitempty"`
	NewUsedCredit  *float64   `json:"new_used_credit,omitempty"`
	NewCreditLimit *float64   `json:"new_credit_limit,omitempty"`
}

// Validate 驗證 ReconcileBankAccountInput
func (in *ReconcileBankAccountInput) Validate() error {
	if in.NewBalance < 0 {
		return fmt.Errorf("new_balance must be >= 0")
	}
	if in.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	return nil
}

// Validate 驗證 ReconcileCreditCardInput
func (in *ReconcileCreditCardInput) Validate() error {
	if in.NewUsedCredit == nil && in.NewCreditLimit == nil {
		return fmt.Errorf("at least one of new_used_credit or new_credit_limit must be provided")
	}
	if in.NewUsedCredit != nil && *in.NewUsedCredit < 0 {
		return fmt.Errorf("new_used_credit must be >= 0")
	}
	if in.NewCreditLimit != nil && *in.NewCreditLimit <= 0 {
		return fmt.Errorf("new_credit_limit must be > 0")
	}
	if in.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	return nil
}

// Validate 驗證 ReconcileBatchInput
func (in *ReconcileBatchInput) Validate() error {
	if in.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	if len(in.Items) == 0 {
		return fmt.Errorf("items must not be empty")
	}
	for i := range in.Items {
		if err := in.Items[i].Validate(); err != nil {
			return fmt.Errorf("items[%d]: %w", i, err)
		}
	}
	return nil
}

// Validate 驗證單一 batch item
func (it *ReconcileBatchItem) Validate() error {
	if it.TargetType != SourceTypeBankAccount && it.TargetType != SourceTypeCreditCard {
		return fmt.Errorf("target_type must be bank_account or credit_card")
	}
	if it.TargetID == uuid.Nil {
		return fmt.Errorf("target_id is required")
	}
	switch it.TargetType {
	case SourceTypeBankAccount:
		if it.NewBalance == nil {
			return fmt.Errorf("bank_account item requires new_balance")
		}
		if *it.NewBalance < 0 {
			return fmt.Errorf("new_balance must be >= 0")
		}
	case SourceTypeCreditCard:
		if it.NewUsedCredit == nil && it.NewCreditLimit == nil {
			return fmt.Errorf("credit_card item requires new_used_credit or new_credit_limit")
		}
		if it.NewUsedCredit != nil && *it.NewUsedCredit < 0 {
			return fmt.Errorf("new_used_credit must be >= 0")
		}
		if it.NewCreditLimit != nil && *it.NewCreditLimit <= 0 {
			return fmt.Errorf("new_credit_limit must be > 0")
		}
	}
	return nil
}
