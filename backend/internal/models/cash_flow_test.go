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

