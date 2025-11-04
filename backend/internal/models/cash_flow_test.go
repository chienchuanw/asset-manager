package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCashFlowType_Validate 測試 CashFlowType 驗證
func TestCashFlowType_Validate(t *testing.T) {
	tests := []struct {
		name     string
		flowType CashFlowType
		want     bool
	}{
		{
			name:     "valid income type",
			flowType: CashFlowTypeIncome,
			want:     true,
		},
		{
			name:     "valid expense type",
			flowType: CashFlowTypeExpense,
			want:     true,
		},
		{
			name:     "invalid type",
			flowType: CashFlowType("invalid"),
			want:     false,
		},
		{
			name:     "empty type",
			flowType: CashFlowType(""),
			want:     false,
		},
		{
			name:     "random string",
			flowType: CashFlowType("random"),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.flowType.Validate()
			assert.Equal(t, tt.want, got, "CashFlowType.Validate() should return %v for %s", tt.want, tt.flowType)
		})
	}
}

// TestCashFlowType_Constants 測試 CashFlowType 常數值
func TestCashFlowType_Constants(t *testing.T) {
	assert.Equal(t, CashFlowType("income"), CashFlowTypeIncome, "CashFlowTypeIncome should be 'income'")
	assert.Equal(t, CashFlowType("expense"), CashFlowTypeExpense, "CashFlowTypeExpense should be 'expense'")
}

// TestSourceType_Validate 測試 SourceType 驗證
func TestSourceType_Validate(t *testing.T) {
	tests := []struct {
		name       string
		sourceType SourceType
		want       bool
	}{
		{
			name:       "valid manual type",
			sourceType: SourceTypeManual,
			want:       true,
		},
		{
			name:       "valid subscription type",
			sourceType: SourceTypeSubscription,
			want:       true,
		},
		{
			name:       "valid installment type",
			sourceType: SourceTypeInstallment,
			want:       true,
		},
		{
			name:       "valid bank account type",
			sourceType: SourceTypeBankAccount,
			want:       true,
		},
		{
			name:       "valid credit card type",
			sourceType: SourceTypeCreditCard,
			want:       true,
		},
		{
			name:       "invalid type",
			sourceType: SourceType("invalid"),
			want:       false,
		},
		{
			name:       "empty type",
			sourceType: SourceType(""),
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sourceType.Validate()
			assert.Equal(t, tt.want, got, "SourceType.Validate() should return %v for %s", tt.want, tt.sourceType)
		})
	}
}

// TestSourceType_Constants 測試 SourceType 常數值
func TestSourceType_Constants(t *testing.T) {
	assert.Equal(t, SourceType("manual"), SourceTypeManual, "SourceTypeManual should be 'manual'")
	assert.Equal(t, SourceType("subscription"), SourceTypeSubscription, "SourceTypeSubscription should be 'subscription'")
	assert.Equal(t, SourceType("installment"), SourceTypeInstallment, "SourceTypeInstallment should be 'installment'")
	assert.Equal(t, SourceType("bank_account"), SourceTypeBankAccount, "SourceTypeBankAccount should be 'bank_account'")
	assert.Equal(t, SourceType("credit_card"), SourceTypeCreditCard, "SourceTypeCreditCard should be 'credit_card'")
}

