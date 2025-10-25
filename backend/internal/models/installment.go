package models

import (
	"time"

	"github.com/google/uuid"
)

// InstallmentStatus 分期狀態
type InstallmentStatus string

const (
	InstallmentStatusActive    InstallmentStatus = "active"    // 進行中
	InstallmentStatusCompleted InstallmentStatus = "completed" // 已完成
	InstallmentStatusCancelled InstallmentStatus = "cancelled" // 已取消
)

// Installment 分期模型
type Installment struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	Name              string            `json:"name" db:"name"`
	TotalAmount       float64           `json:"total_amount" db:"total_amount"`
	Currency          Currency          `json:"currency" db:"currency"`
	InstallmentCount  int               `json:"installment_count" db:"installment_count"`
	InstallmentAmount float64           `json:"installment_amount" db:"installment_amount"`
	InterestRate      float64           `json:"interest_rate" db:"interest_rate"`
	TotalInterest     float64           `json:"total_interest" db:"total_interest"`
	PaidCount         int               `json:"paid_count" db:"paid_count"`
	BillingDay        int               `json:"billing_day" db:"billing_day"`
	CategoryID        uuid.UUID         `json:"category_id" db:"category_id"`
	StartDate         time.Time         `json:"start_date" db:"start_date"`
	Status            InstallmentStatus `json:"status" db:"status"`
	Note              *string           `json:"note,omitempty" db:"note"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`

	// 關聯資料（Join 時使用）
	Category *CashFlowCategory `json:"category,omitempty" db:"-"`
}

// CreateInstallmentInput 建立分期的輸入資料
type CreateInstallmentInput struct {
	Name             string    `json:"name" binding:"required,max=255"`
	TotalAmount      float64   `json:"total_amount" binding:"required,gt=0"`
	InstallmentCount int       `json:"installment_count" binding:"required,gt=0"`
	InterestRate     float64   `json:"interest_rate" binding:"gte=0"`
	BillingDay       int       `json:"billing_day" binding:"required,min=1,max=31"`
	CategoryID       uuid.UUID `json:"category_id" binding:"required"`
	StartDate        time.Time `json:"start_date" binding:"required"`
	Note             *string   `json:"note,omitempty"`
}

// UpdateInstallmentInput 更新分期的輸入資料
type UpdateInstallmentInput struct {
	Name         *string    `json:"name,omitempty" binding:"omitempty,max=255"`
	InterestRate *float64   `json:"interest_rate,omitempty" binding:"omitempty,gte=0"`
	CategoryID   *uuid.UUID `json:"category_id,omitempty"`
	Note         *string    `json:"note,omitempty"`
}

// Validate 驗證 InstallmentStatus 是否有效
func (s InstallmentStatus) Validate() bool {
	switch s {
	case InstallmentStatusActive, InstallmentStatusCompleted, InstallmentStatusCancelled:
		return true
	}
	return false
}

// CalculateInterest 計算利息和每期金額
func (i *Installment) CalculateInterest() {
	// 計算總利息：本金 × 利率
	i.TotalInterest = i.TotalAmount * (i.InterestRate / 100)

	// 計算每期金額：(本金 + 總利息) / 期數
	totalWithInterest := i.TotalAmount + i.TotalInterest
	i.InstallmentAmount = totalWithInterest / float64(i.InstallmentCount)
}

// RemainingAmount 計算剩餘金額
func (i *Installment) RemainingAmount() float64 {
	totalWithInterest := i.TotalAmount + i.TotalInterest
	paidAmount := float64(i.PaidCount) * i.InstallmentAmount
	return totalWithInterest - paidAmount
}

// RemainingCount 計算剩餘期數
func (i *Installment) RemainingCount() int {
	return i.InstallmentCount - i.PaidCount
}

// NextBillingDate 計算下次扣款日期
// fromDate: 從哪個日期開始計算（通常是今天）
func (i *Installment) NextBillingDate(fromDate time.Time) time.Time {
	// 從開始日期計算下次扣款的月份
	startYear := i.StartDate.Year()
	startMonth := i.StartDate.Month()

	// 計算目標月份（開始月份 + 已付期數）
	targetMonth := int(startMonth) + i.PaidCount
	targetYear := startYear

	// 處理跨年的情況
	for targetMonth > 12 {
		targetMonth -= 12
		targetYear++
	}

	// 建立扣款日期
	billingDate := time.Date(targetYear, time.Month(targetMonth), i.BillingDay, 0, 0, 0, 0, time.UTC)

	// 處理月份天數不足的情況（例如 2 月沒有 31 日）
	if billingDate.Month() != time.Month(targetMonth) {
		// 如果日期溢出到下個月，使用當月最後一天
		billingDate = time.Date(targetYear, time.Month(targetMonth)+1, 0, 0, 0, 0, 0, time.UTC)
	}

	return billingDate
}

// IsActive 檢查分期是否進行中
func (i *Installment) IsActive() bool {
	// 狀態必須是 active
	if i.Status != InstallmentStatusActive {
		return false
	}

	// 如果已付期數等於總期數，表示已完成
	if i.PaidCount >= i.InstallmentCount {
		return false
	}

	return true
}

