package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestBillingCycle_Validate 測試 BillingCycle 驗證
func TestBillingCycle_Validate(t *testing.T) {
	tests := []struct {
		name  string
		cycle BillingCycle
		want  bool
	}{
		{
			name:  "valid monthly cycle",
			cycle: BillingCycleMonthly,
			want:  true,
		},
		{
			name:  "valid quarterly cycle",
			cycle: BillingCycleQuarterly,
			want:  true,
		},
		{
			name:  "valid yearly cycle",
			cycle: BillingCycleYearly,
			want:  true,
		},
		{
			name:  "invalid cycle",
			cycle: BillingCycle("invalid"),
			want:  false,
		},
		{
			name:  "empty cycle",
			cycle: BillingCycle(""),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cycle.Validate()
			assert.Equal(t, tt.want, got, "BillingCycle.Validate() should return %v for %s", tt.want, tt.cycle)
		})
	}
}

// TestBillingCycle_Constants 測試 BillingCycle 常數值
func TestBillingCycle_Constants(t *testing.T) {
	assert.Equal(t, BillingCycle("monthly"), BillingCycleMonthly, "BillingCycleMonthly should be 'monthly'")
	assert.Equal(t, BillingCycle("quarterly"), BillingCycleQuarterly, "BillingCycleQuarterly should be 'quarterly'")
	assert.Equal(t, BillingCycle("yearly"), BillingCycleYearly, "BillingCycleYearly should be 'yearly'")
}

// TestSubscriptionStatus_Validate 測試 SubscriptionStatus 驗證
func TestSubscriptionStatus_Validate(t *testing.T) {
	tests := []struct {
		name   string
		status SubscriptionStatus
		want   bool
	}{
		{
			name:   "valid active status",
			status: SubscriptionStatusActive,
			want:   true,
		},
		{
			name:   "valid cancelled status",
			status: SubscriptionStatusCancelled,
			want:   true,
		},
		{
			name:   "invalid status",
			status: SubscriptionStatus("invalid"),
			want:   false,
		},
		{
			name:   "empty status",
			status: SubscriptionStatus(""),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.Validate()
			assert.Equal(t, tt.want, got, "SubscriptionStatus.Validate() should return %v for %s", tt.want, tt.status)
		})
	}
}

// TestSubscriptionStatus_Constants 測試 SubscriptionStatus 常數值
func TestSubscriptionStatus_Constants(t *testing.T) {
	assert.Equal(t, SubscriptionStatus("active"), SubscriptionStatusActive, "SubscriptionStatusActive should be 'active'")
	assert.Equal(t, SubscriptionStatus("cancelled"), SubscriptionStatusCancelled, "SubscriptionStatusCancelled should be 'cancelled'")
}

// TestSubscription_NextBillingDate 測試計算下次扣款日期
func TestSubscription_NextBillingDate(t *testing.T) {
	tests := []struct {
		name        string
		startDate   time.Time
		billingDay  int
		billingCycle BillingCycle
		fromDate    time.Time
		want        time.Time
	}{
		{
			name:         "monthly - same month, before billing day",
			startDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			billingDay:   15,
			billingCycle: BillingCycleMonthly,
			fromDate:     time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "monthly - same month, after billing day",
			startDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			billingDay:   15,
			billingCycle: BillingCycleMonthly,
			fromDate:     time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "monthly - next month",
			startDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			billingDay:   15,
			billingCycle: BillingCycleMonthly,
			fromDate:     time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "quarterly - first quarter",
			startDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			billingDay:   15,
			billingCycle: BillingCycleQuarterly,
			fromDate:     time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "quarterly - after first billing",
			startDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			billingDay:   15,
			billingCycle: BillingCycleQuarterly,
			fromDate:     time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "yearly - first year",
			startDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			billingDay:   15,
			billingCycle: BillingCycleYearly,
			fromDate:     time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "yearly - after first billing",
			startDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			billingDay:   15,
			billingCycle: BillingCycleYearly,
			fromDate:     time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "monthly - billing day 31 in February",
			startDate:    time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			billingDay:   31,
			billingCycle: BillingCycleMonthly,
			fromDate:     time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC), // February has only 28 days
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{
				StartDate:    tt.startDate,
				BillingDay:   tt.billingDay,
				BillingCycle: tt.billingCycle,
			}
			got := sub.NextBillingDate(tt.fromDate)
			assert.Equal(t, tt.want, got, "NextBillingDate() should return %v", tt.want)
		})
	}
}

// TestSubscription_PaymentMethodValidation 測試訂閱付款方式驗證
func TestSubscription_PaymentMethodValidation(t *testing.T) {
	bankAccountID := uuid.New()
	creditCardID := uuid.New()

	tests := []struct {
		name          string
		paymentMethod PaymentMethod
		accountID     *uuid.UUID
		wantValid     bool
	}{
		{
			name:          "cash without account ID is valid",
			paymentMethod: PaymentMethodCash,
			accountID:     nil,
			wantValid:     true,
		},
		{
			name:          "cash with account ID is also valid (ignored)",
			paymentMethod: PaymentMethodCash,
			accountID:     &bankAccountID,
			wantValid:     true,
		},
		{
			name:          "bank_account with account ID is valid",
			paymentMethod: PaymentMethodBankAccount,
			accountID:     &bankAccountID,
			wantValid:     true,
		},
		{
			name:          "bank_account without account ID is invalid",
			paymentMethod: PaymentMethodBankAccount,
			accountID:     nil,
			wantValid:     false,
		},
		{
			name:          "credit_card with account ID is valid",
			paymentMethod: PaymentMethodCreditCard,
			accountID:     &creditCardID,
			wantValid:     true,
		},
		{
			name:          "credit_card without account ID is invalid",
			paymentMethod: PaymentMethodCreditCard,
			accountID:     nil,
			wantValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{
				PaymentMethod: tt.paymentMethod,
				AccountID:     tt.accountID,
			}
			got := sub.ValidatePaymentMethod()
			assert.Equal(t, tt.wantValid, got, "ValidatePaymentMethod() should return %v", tt.wantValid)
		})
	}
}

// TestSubscription_IsActive 測試訂閱是否進行中
func TestSubscription_IsActive(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		status    SubscriptionStatus
		endDate   *time.Time
		want      bool
	}{
		{
			name:    "active status, no end date",
			status:  SubscriptionStatusActive,
			endDate: nil,
			want:    true,
		},
		{
			name:    "active status, end date in future",
			status:  SubscriptionStatusActive,
			endDate: &[]time.Time{now.AddDate(0, 1, 0)}[0],
			want:    true,
		},
		{
			name:    "active status, end date in past",
			status:  SubscriptionStatusActive,
			endDate: &[]time.Time{now.AddDate(0, -1, 0)}[0],
			want:    false,
		},
		{
			name:    "cancelled status",
			status:  SubscriptionStatusCancelled,
			endDate: nil,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{
				Status:  tt.status,
				EndDate: tt.endDate,
			}
			got := sub.IsActive()
			assert.Equal(t, tt.want, got, "IsActive() should return %v", tt.want)
		})
	}
}

