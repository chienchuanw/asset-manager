package models

import (
	"time"

	"github.com/google/uuid"
)

// BillingCycle 計費週期
type BillingCycle string

const (
	BillingCycleMonthly   BillingCycle = "monthly"   // 月繳
	BillingCycleQuarterly BillingCycle = "quarterly" // 季繳
	BillingCycleYearly    BillingCycle = "yearly"    // 年繳
)

// SubscriptionStatus 訂閱狀態
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"    // 進行中
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled" // 已取消
)

// Subscription 訂閱模型
type Subscription struct {
	ID           uuid.UUID          `json:"id" db:"id"`
	Name         string             `json:"name" db:"name"`
	Amount       float64            `json:"amount" db:"amount"`
	Currency     Currency           `json:"currency" db:"currency"`
	BillingCycle BillingCycle       `json:"billing_cycle" db:"billing_cycle"`
	BillingDay   int                `json:"billing_day" db:"billing_day"`
	CategoryID   uuid.UUID          `json:"category_id" db:"category_id"`
	StartDate    time.Time          `json:"start_date" db:"start_date"`
	EndDate      *time.Time         `json:"end_date,omitempty" db:"end_date"`
	AutoRenew    bool               `json:"auto_renew" db:"auto_renew"`
	Status       SubscriptionStatus `json:"status" db:"status"`
	Note         *string            `json:"note,omitempty" db:"note"`
	CreatedAt    time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" db:"updated_at"`

	// 關聯資料（Join 時使用）
	Category *CashFlowCategory `json:"category,omitempty" db:"-"`
}

// CreateSubscriptionInput 建立訂閱的輸入資料
type CreateSubscriptionInput struct {
	Name         string       `json:"name" binding:"required,max=255"`
	Amount       float64      `json:"amount" binding:"required,gt=0"`
	BillingCycle BillingCycle `json:"billing_cycle" binding:"required"`
	BillingDay   int          `json:"billing_day" binding:"required,min=1,max=31"`
	CategoryID   uuid.UUID    `json:"category_id" binding:"required"`
	StartDate    time.Time    `json:"start_date" binding:"required"`
	EndDate      *time.Time   `json:"end_date,omitempty"`
	AutoRenew    bool         `json:"auto_renew"`
	Note         *string      `json:"note,omitempty"`
}

// UpdateSubscriptionInput 更新訂閱的輸入資料
type UpdateSubscriptionInput struct {
	Name         *string              `json:"name,omitempty" binding:"omitempty,max=255"`
	Amount       *float64             `json:"amount,omitempty" binding:"omitempty,gt=0"`
	BillingCycle *BillingCycle        `json:"billing_cycle,omitempty"`
	BillingDay   *int                 `json:"billing_day,omitempty" binding:"omitempty,min=1,max=31"`
	CategoryID   *uuid.UUID           `json:"category_id,omitempty"`
	EndDate      *time.Time           `json:"end_date,omitempty"`
	AutoRenew    *bool                `json:"auto_renew,omitempty"`
	Status       *SubscriptionStatus  `json:"status,omitempty"`
	Note         *string              `json:"note,omitempty"`
}

// Validate 驗證 BillingCycle 是否有效
func (b BillingCycle) Validate() bool {
	switch b {
	case BillingCycleMonthly, BillingCycleQuarterly, BillingCycleYearly:
		return true
	}
	return false
}

// Validate 驗證 SubscriptionStatus 是否有效
func (s SubscriptionStatus) Validate() bool {
	switch s {
	case SubscriptionStatusActive, SubscriptionStatusCancelled:
		return true
	}
	return false
}

// NextBillingDate 計算下次扣款日期
// fromDate: 從哪個日期開始計算（通常是今天）
func (s *Subscription) NextBillingDate(fromDate time.Time) time.Time {
	// 如果 fromDate 在開始日期之前，返回開始日期
	if fromDate.Before(s.StartDate) {
		return s.StartDate
	}

	// 根據計費週期計算下次扣款日期
	switch s.BillingCycle {
	case BillingCycleMonthly:
		return s.nextMonthlyBillingDate(fromDate)
	case BillingCycleQuarterly:
		return s.nextQuarterlyBillingDate(fromDate)
	case BillingCycleYearly:
		return s.nextYearlyBillingDate(fromDate)
	default:
		return fromDate
	}
}

// nextMonthlyBillingDate 計算下次月繳扣款日期
func (s *Subscription) nextMonthlyBillingDate(fromDate time.Time) time.Time {
	year := fromDate.Year()
	month := fromDate.Month()
	day := s.BillingDay

	// 建立當月的扣款日期
	billingDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	// 處理月份天數不足的情況（例如 2 月沒有 31 日）
	if billingDate.Month() != month {
		// 如果日期溢出到下個月，使用當月最後一天
		billingDate = time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	}

	// 如果當月扣款日已過，計算下個月
	if fromDate.After(billingDate) || fromDate.Equal(billingDate) {
		month++
		if month > 12 {
			month = 1
			year++
		}
		billingDate = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		// 再次處理月份天數不足的情況
		if billingDate.Month() != month {
			billingDate = time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
		}
	}

	return billingDate
}

// nextQuarterlyBillingDate 計算下次季繳扣款日期
func (s *Subscription) nextQuarterlyBillingDate(fromDate time.Time) time.Time {
	// 從開始日期計算每季的扣款月份
	startMonth := s.StartDate.Month()
	currentYear := fromDate.Year()
	currentMonth := fromDate.Month()

	// 計算下次扣款的月份（每 3 個月一次）
	var nextMonth time.Month
	for m := startMonth; m <= 12; m += 3 {
		if m > currentMonth || (m == currentMonth && fromDate.Day() < s.BillingDay) {
			nextMonth = m
			break
		}
	}

	// 如果今年沒有找到，使用明年的第一個季度
	if nextMonth == 0 {
		nextMonth = startMonth
		currentYear++
	}

	billingDate := time.Date(currentYear, nextMonth, s.BillingDay, 0, 0, 0, 0, time.UTC)

	// 處理月份天數不足的情況
	if billingDate.Month() != nextMonth {
		billingDate = time.Date(currentYear, nextMonth+1, 0, 0, 0, 0, 0, time.UTC)
	}

	return billingDate
}

// nextYearlyBillingDate 計算下次年繳扣款日期
func (s *Subscription) nextYearlyBillingDate(fromDate time.Time) time.Time {
	startMonth := s.StartDate.Month()
	year := fromDate.Year()

	// 建立今年的扣款日期
	billingDate := time.Date(year, startMonth, s.BillingDay, 0, 0, 0, 0, time.UTC)

	// 處理月份天數不足的情況
	if billingDate.Month() != startMonth {
		billingDate = time.Date(year, startMonth+1, 0, 0, 0, 0, 0, time.UTC)
	}

	// 如果今年的扣款日已過，使用明年
	if fromDate.After(billingDate) || fromDate.Equal(billingDate) {
		year++
		billingDate = time.Date(year, startMonth, s.BillingDay, 0, 0, 0, 0, time.UTC)
		// 再次處理月份天數不足的情況
		if billingDate.Month() != startMonth {
			billingDate = time.Date(year, startMonth+1, 0, 0, 0, 0, 0, time.UTC)
		}
	}

	return billingDate
}

// IsActive 檢查訂閱是否進行中
func (s *Subscription) IsActive() bool {
	// 狀態必須是 active
	if s.Status != SubscriptionStatusActive {
		return false
	}

	// 如果有結束日期，檢查是否已過期
	if s.EndDate != nil && time.Now().After(*s.EndDate) {
		return false
	}

	return true
}

