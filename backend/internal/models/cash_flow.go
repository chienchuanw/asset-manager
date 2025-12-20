package models

import (
	"time"

	"github.com/google/uuid"
)

// CashFlowType 現金流類型
type CashFlowType string

const (
	CashFlowTypeIncome     CashFlowType = "income"      // 收入
	CashFlowTypeExpense    CashFlowType = "expense"     // 支出
	CashFlowTypeTransferIn CashFlowType = "transfer_in" // 存入帳戶
	CashFlowTypeTransferOut CashFlowType = "transfer_out" // 從帳戶轉出
)

// SourceType 現金流來源類型
type SourceType string

const (
	SourceTypeManual       SourceType = "manual"       // 手動建立（現金交易）
	SourceTypeSubscription SourceType = "subscription" // 訂閱自動產生
	SourceTypeInstallment  SourceType = "installment"  // 分期自動產生
	SourceTypeBankAccount  SourceType = "bank_account" // 銀行帳戶交易
	SourceTypeCreditCard   SourceType = "credit_card"  // 信用卡交易
	SourceTypeCash         SourceType = "cash"         // 現金
)

// CashFlowCategory 現金流分類模型
type CashFlowCategory struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	Name      string       `json:"name" db:"name"`
	Type      CashFlowType `json:"type" db:"type"`
	IsSystem  bool         `json:"is_system" db:"is_system"`
	SortOrder int          `json:"sort_order" db:"sort_order"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
}

// CashFlow 現金流記錄模型
type CashFlow struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Date        time.Time    `json:"date" db:"date"`
	Type        CashFlowType `json:"type" db:"type"`
	CategoryID  uuid.UUID    `json:"category_id" db:"category_id"`
	Amount      float64      `json:"amount" db:"amount"`
	Currency    Currency     `json:"currency" db:"currency"`
	Description string       `json:"description" db:"description"`
	Note        *string      `json:"note,omitempty" db:"note"`
	SourceType  *SourceType  `json:"source_type,omitempty" db:"source_type"`
	SourceID    *uuid.UUID   `json:"source_id,omitempty" db:"source_id"`
	TargetType  *SourceType  `json:"target_type,omitempty" db:"target_type"` // 轉帳目標類型
	TargetID    *uuid.UUID   `json:"target_id,omitempty" db:"target_id"`     // 轉帳目標ID
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`

	// 關聯資料（Join 時使用）
	Category *CashFlowCategory `json:"category,omitempty" db:"-"`
}

// CreateCashFlowInput 建立現金流記錄的輸入資料
type CreateCashFlowInput struct {
	Date        time.Time    `json:"date" binding:"required"`
	Type        CashFlowType `json:"type" binding:"required"`
	CategoryID  uuid.UUID    `json:"category_id" binding:"required"`
	Amount      float64      `json:"amount" binding:"required,gt=0"`
	Description string       `json:"description" binding:"required,max=500"`
	Note        *string      `json:"note,omitempty"`
	SourceType  *SourceType  `json:"source_type,omitempty"`
	SourceID    *uuid.UUID   `json:"source_id,omitempty"`
	TargetType  *SourceType  `json:"target_type,omitempty"` // 轉帳目標類型
	TargetID    *uuid.UUID   `json:"target_id,omitempty"`   // 轉帳目標ID
}

// UpdateCashFlowInput 更新現金流記錄的輸入資料
type UpdateCashFlowInput struct {
	Date        *time.Time  `json:"date,omitempty"`
	CategoryID  *uuid.UUID  `json:"category_id,omitempty"`
	Amount      *float64    `json:"amount,omitempty" binding:"omitempty,gt=0"`
	Description *string     `json:"description,omitempty" binding:"omitempty,max=500"`
	Note        *string     `json:"note,omitempty"`
	SourceType  *SourceType `json:"source_type,omitempty"`
	SourceID    *uuid.UUID  `json:"source_id,omitempty"`
	TargetType  *SourceType `json:"target_type,omitempty"` // 轉帳目標類型
	TargetID    *uuid.UUID  `json:"target_id,omitempty"`   // 轉帳目標ID
}

// CreateCategoryInput 建立分類的輸入資料
type CreateCategoryInput struct {
	Name string       `json:"name" binding:"required,max=20"`
	Type CashFlowType `json:"type" binding:"required"`
}

// UpdateCategoryInput 更新分類的輸入資料
type UpdateCategoryInput struct {
	Name string `json:"name" binding:"required,max=20"`
}

// CategoryOrderItem 單筆分類排序項目
type CategoryOrderItem struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	SortOrder int       `json:"sort_order" binding:"gte=0"`
}

// ReorderCategoryInput 批次更新分類排序的輸入資料
type ReorderCategoryInput struct {
	Orders []CategoryOrderItem `json:"orders" binding:"required,min=1,dive"`
}

// Validate 驗證 CashFlowType 是否有效
func (t CashFlowType) Validate() bool {
	switch t {
	case CashFlowTypeIncome, CashFlowTypeExpense, CashFlowTypeTransferIn, CashFlowTypeTransferOut:
		return true
	}
	return false
}

// Validate 驗證 SourceType 是否有效
func (s SourceType) Validate() bool {
	switch s {
	case SourceTypeManual, SourceTypeSubscription, SourceTypeInstallment, SourceTypeBankAccount, SourceTypeCreditCard, SourceTypeCash:
		return true
	}
	return false
}

