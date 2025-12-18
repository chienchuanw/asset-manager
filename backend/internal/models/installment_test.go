package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestInstallmentStatus_Validate 測試 InstallmentStatus 驗證
func TestInstallmentStatus_Validate(t *testing.T) {
	tests := []struct {
		name   string
		status InstallmentStatus
		want   bool
	}{
		{
			name:   "valid active status",
			status: InstallmentStatusActive,
			want:   true,
		},
		{
			name:   "valid completed status",
			status: InstallmentStatusCompleted,
			want:   true,
		},
		{
			name:   "valid cancelled status",
			status: InstallmentStatusCancelled,
			want:   true,
		},
		{
			name:   "invalid status",
			status: InstallmentStatus("invalid"),
			want:   false,
		},
		{
			name:   "empty status",
			status: InstallmentStatus(""),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.Validate()
			assert.Equal(t, tt.want, got, "InstallmentStatus.Validate() should return %v for %s", tt.want, tt.status)
		})
	}
}

// TestInstallmentStatus_Constants 測試 InstallmentStatus 常數值
func TestInstallmentStatus_Constants(t *testing.T) {
	assert.Equal(t, InstallmentStatus("active"), InstallmentStatusActive, "InstallmentStatusActive should be 'active'")
	assert.Equal(t, InstallmentStatus("completed"), InstallmentStatusCompleted, "InstallmentStatusCompleted should be 'completed'")
	assert.Equal(t, InstallmentStatus("cancelled"), InstallmentStatusCancelled, "InstallmentStatusCancelled should be 'cancelled'")
}

// TestInstallment_CalculateInterest 測試利息計算
func TestInstallment_CalculateInterest(t *testing.T) {
	tests := []struct {
		name              string
		totalAmount       float64
		installmentCount  int
		interestRate      float64
		wantTotalInterest float64
		wantInstallmentAmount float64
	}{
		{
			name:              "no interest",
			totalAmount:       10000,
			installmentCount:  10,
			interestRate:      0,
			wantTotalInterest: 0,
			wantInstallmentAmount: 1000, // 10000 / 10
		},
		{
			name:              "5% interest rate",
			totalAmount:       10000,
			installmentCount:  10,
			interestRate:      5,
			wantTotalInterest: 500, // 10000 * 0.05
			wantInstallmentAmount: 1050, // (10000 + 500) / 10
		},
		{
			name:              "10% interest rate",
			totalAmount:       20000,
			installmentCount:  12,
			interestRate:      10,
			wantTotalInterest: 2000, // 20000 * 0.1
			wantInstallmentAmount: 1833.33, // (20000 + 2000) / 12
		},
		{
			name:              "fractional interest rate",
			totalAmount:       15000,
			installmentCount:  6,
			interestRate:      3.5,
			wantTotalInterest: 525, // 15000 * 0.035
			wantInstallmentAmount: 2587.5, // (15000 + 525) / 6
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &Installment{
				TotalAmount:      tt.totalAmount,
				InstallmentCount: tt.installmentCount,
				InterestRate:     tt.interestRate,
			}
			inst.CalculateInterest()
			
			assert.InDelta(t, tt.wantTotalInterest, inst.TotalInterest, 0.01, "TotalInterest should be %v", tt.wantTotalInterest)
			assert.InDelta(t, tt.wantInstallmentAmount, inst.InstallmentAmount, 0.01, "InstallmentAmount should be %v", tt.wantInstallmentAmount)
		})
	}
}

// TestInstallment_RemainingAmount 測試剩餘金額計算
func TestInstallment_RemainingAmount(t *testing.T) {
	tests := []struct {
		name              string
		totalAmount       float64
		totalInterest     float64
		installmentAmount float64
		paidCount         int
		want              float64
	}{
		{
			name:              "no payment yet",
			totalAmount:       10000,
			totalInterest:     500,
			installmentAmount: 1050,
			paidCount:         0,
			want:              10500, // 10000 + 500
		},
		{
			name:              "half paid",
			totalAmount:       10000,
			totalInterest:     500,
			installmentAmount: 1050,
			paidCount:         5,
			want:              5250, // (10000 + 500) - (1050 * 5)
		},
		{
			name:              "almost complete",
			totalAmount:       10000,
			totalInterest:     500,
			installmentAmount: 1050,
			paidCount:         9,
			want:              1050, // (10000 + 500) - (1050 * 9)
		},
		{
			name:              "fully paid",
			totalAmount:       10000,
			totalInterest:     500,
			installmentAmount: 1050,
			paidCount:         10,
			want:              0, // (10000 + 500) - (1050 * 10)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &Installment{
				TotalAmount:       tt.totalAmount,
				TotalInterest:     tt.totalInterest,
				InstallmentAmount: tt.installmentAmount,
				PaidCount:         tt.paidCount,
			}
			got := inst.RemainingAmount()
			assert.InDelta(t, tt.want, got, 0.01, "RemainingAmount() should return %v", tt.want)
		})
	}
}

// TestInstallment_RemainingCount 測試剩餘期數計算
func TestInstallment_RemainingCount(t *testing.T) {
	tests := []struct {
		name             string
		installmentCount int
		paidCount        int
		want             int
	}{
		{
			name:             "no payment yet",
			installmentCount: 12,
			paidCount:        0,
			want:             12,
		},
		{
			name:             "half paid",
			installmentCount: 12,
			paidCount:        6,
			want:             6,
		},
		{
			name:             "almost complete",
			installmentCount: 12,
			paidCount:        11,
			want:             1,
		},
		{
			name:             "fully paid",
			installmentCount: 12,
			paidCount:        12,
			want:             0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &Installment{
				InstallmentCount: tt.installmentCount,
				PaidCount:        tt.paidCount,
			}
			got := inst.RemainingCount()
			assert.Equal(t, tt.want, got, "RemainingCount() should return %v", tt.want)
		})
	}
}

// TestInstallment_NextBillingDate 測試計算下次扣款日期
func TestInstallment_NextBillingDate(t *testing.T) {
	tests := []struct {
		name       string
		startDate  time.Time
		billingDay int
		paidCount  int
		fromDate   time.Time
		want       time.Time
	}{
		{
			name:       "first payment",
			startDate:  time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			billingDay: 15,
			paidCount:  0,
			fromDate:   time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			want:       time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "second payment",
			startDate:  time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			billingDay: 15,
			paidCount:  1,
			fromDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			want:       time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "third payment",
			startDate:  time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			billingDay: 15,
			paidCount:  2,
			fromDate:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			want:       time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "billing day 31 in February",
			startDate:  time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			billingDay: 31,
			paidCount:  1,
			fromDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			want:       time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC), // February has only 28 days
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &Installment{
				StartDate:  tt.startDate,
				BillingDay: tt.billingDay,
				PaidCount:  tt.paidCount,
			}
			got := inst.NextBillingDate(tt.fromDate)
			assert.Equal(t, tt.want, got, "NextBillingDate() should return %v", tt.want)
		})
	}
}

// TestInstallment_IsActive 測試分期是否進行中
func TestInstallment_IsActive(t *testing.T) {
	tests := []struct {
		name             string
		status           InstallmentStatus
		installmentCount int
		paidCount        int
		want             bool
	}{
		{
			name:             "active status, not fully paid",
			status:           InstallmentStatusActive,
			installmentCount: 12,
			paidCount:        6,
			want:             true,
		},
		{
			name:             "active status, fully paid",
			status:           InstallmentStatusActive,
			installmentCount: 12,
			paidCount:        12,
			want:             false,
		},
		{
			name:             "completed status",
			status:           InstallmentStatusCompleted,
			installmentCount: 12,
			paidCount:        12,
			want:             false,
		},
		{
			name:             "cancelled status",
			status:           InstallmentStatusCancelled,
			installmentCount: 12,
			paidCount:        6,
			want:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &Installment{
				Status:           tt.status,
				InstallmentCount: tt.installmentCount,
				PaidCount:        tt.paidCount,
			}
			got := inst.IsActive()
			assert.Equal(t, tt.want, got, "IsActive() should return %v", tt.want)
		})
	}
}

// TestInstallment_PaymentMethodValidation 測試分期付款方式驗證
func TestInstallment_PaymentMethodValidation(t *testing.T) {
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
			inst := &Installment{
				PaymentMethod: tt.paymentMethod,
				AccountID:     tt.accountID,
			}
			got := inst.ValidatePaymentMethod()
			assert.Equal(t, tt.wantValid, got, "ValidatePaymentMethod() should return %v", tt.wantValid)
		})
	}
}
