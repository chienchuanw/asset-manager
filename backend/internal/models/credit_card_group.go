package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreditCardGroup 信用卡群組模型
type CreditCardGroup struct {
	ID                uuid.UUID `json:"id" db:"id"`
	Name              string    `json:"name" db:"name"`
	IssuingBank       string    `json:"issuing_bank" db:"issuing_bank"`
	SharedCreditLimit float64   `json:"shared_credit_limit" db:"shared_credit_limit"`
	Note              *string   `json:"note,omitempty" db:"note"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// CreditCardGroupWithCards 信用卡群組及其包含的卡片
type CreditCardGroupWithCards struct {
	CreditCardGroup
	Cards           []*CreditCard `json:"cards"`
	TotalUsedCredit float64       `json:"total_used_credit"`
}

// CreateCreditCardGroupInput 建立信用卡群組的輸入資料
type CreateCreditCardGroupInput struct {
	Name              string    `json:"name" binding:"required,max=255"`
	IssuingBank       string    `json:"issuing_bank" binding:"required,max=255"`
	SharedCreditLimit float64   `json:"shared_credit_limit" binding:"required,gt=0"`
	CardIDs           []string  `json:"card_ids" binding:"required,min=1"`
	Note              *string   `json:"note,omitempty" binding:"omitempty,max=1000"`
}

// UpdateCreditCardGroupInput 更新信用卡群組的輸入資料
type UpdateCreditCardGroupInput struct {
	Name              *string  `json:"name,omitempty" binding:"omitempty,max=255"`
	SharedCreditLimit *float64 `json:"shared_credit_limit,omitempty" binding:"omitempty,gt=0"`
	Note              *string  `json:"note,omitempty" binding:"omitempty,max=1000"`
}

// AddCardsToGroupInput 新增卡片到群組的輸入資料
type AddCardsToGroupInput struct {
	CardIDs []string `json:"card_ids" binding:"required,min=1"`
}

// RemoveCardsFromGroupInput 從群組移除卡片的輸入資料
type RemoveCardsFromGroupInput struct {
	CardIDs []string `json:"card_ids" binding:"required,min=1"`
}

// Validate 驗證建立信用卡群組的輸入資料
func (input *CreateCreditCardGroupInput) Validate() error {
	// 驗證群組名稱不能為空
	if input.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	// 驗證發卡銀行不能為空
	if input.IssuingBank == "" {
		return fmt.Errorf("issuing_bank cannot be empty")
	}

	// 驗證共享信用額度必須大於 0
	if input.SharedCreditLimit <= 0 {
		return fmt.Errorf("shared_credit_limit must be greater than 0")
	}

	// 驗證至少要有一張卡片
	if len(input.CardIDs) == 0 {
		return fmt.Errorf("at least one card is required")
	}

	// 驗證所有卡片 ID 格式
	for _, cardID := range input.CardIDs {
		if _, err := uuid.Parse(cardID); err != nil {
			return fmt.Errorf("invalid card ID format: %s", cardID)
		}
	}

	return nil
}

// Validate 驗證更新信用卡群組的輸入資料
func (input *UpdateCreditCardGroupInput) Validate() error {
	// 如果有提供群組名稱,驗證不能為空
	if input.Name != nil && *input.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	// 如果有提供共享信用額度,驗證必須大於 0
	if input.SharedCreditLimit != nil && *input.SharedCreditLimit <= 0 {
		return fmt.Errorf("shared_credit_limit must be greater than 0")
	}

	return nil
}

// Validate 驗證新增卡片到群組的輸入資料
func (input *AddCardsToGroupInput) Validate() error {
	// 驗證至少要有一張卡片
	if len(input.CardIDs) == 0 {
		return fmt.Errorf("at least one card is required")
	}

	// 驗證所有卡片 ID 格式
	for _, cardID := range input.CardIDs {
		if _, err := uuid.Parse(cardID); err != nil {
			return fmt.Errorf("invalid card ID format: %s", cardID)
		}
	}

	return nil
}

// Validate 驗證從群組移除卡片的輸入資料
func (input *RemoveCardsFromGroupInput) Validate() error {
	// 驗證至少要有一張卡片
	if len(input.CardIDs) == 0 {
		return fmt.Errorf("at least one card is required")
	}

	// 驗證所有卡片 ID 格式
	for _, cardID := range input.CardIDs {
		if _, err := uuid.Parse(cardID); err != nil {
			return fmt.Errorf("invalid card ID format: %s", cardID)
		}
	}

	return nil
}

// AvailableCredit 計算群組的可用額度
func (g *CreditCardGroupWithCards) AvailableCredit() float64 {
	return g.SharedCreditLimit - g.TotalUsedCredit
}

// CreditUtilization 計算群組的信用額度使用率(百分比)
func (g *CreditCardGroupWithCards) CreditUtilization() float64 {
	if g.SharedCreditLimit == 0 {
		return 0
	}
	return (g.TotalUsedCredit / g.SharedCreditLimit) * 100
}

