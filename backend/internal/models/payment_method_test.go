package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPaymentMethod_Validate 測試 PaymentMethod 驗證
func TestPaymentMethod_Validate(t *testing.T) {
	tests := []struct {
		name   string
		method PaymentMethod
		want   bool
	}{
		{
			name:   "valid cash method",
			method: PaymentMethodCash,
			want:   true,
		},
		{
			name:   "valid bank_account method",
			method: PaymentMethodBankAccount,
			want:   true,
		},
		{
			name:   "valid credit_card method",
			method: PaymentMethodCreditCard,
			want:   true,
		},
		{
			name:   "invalid method",
			method: PaymentMethod("invalid"),
			want:   false,
		},
		{
			name:   "empty method",
			method: PaymentMethod(""),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method.Validate()
			assert.Equal(t, tt.want, got, "PaymentMethod.Validate() should return %v for %s", tt.want, tt.method)
		})
	}
}

// TestPaymentMethod_Constants 測試 PaymentMethod 常數值
func TestPaymentMethod_Constants(t *testing.T) {
	assert.Equal(t, PaymentMethod("cash"), PaymentMethodCash, "PaymentMethodCash should be 'cash'")
	assert.Equal(t, PaymentMethod("bank_account"), PaymentMethodBankAccount, "PaymentMethodBankAccount should be 'bank_account'")
	assert.Equal(t, PaymentMethod("credit_card"), PaymentMethodCreditCard, "PaymentMethodCreditCard should be 'credit_card'")
}

// TestPaymentMethod_RequiresAccountID 測試付款方式是否需要帳戶 ID
func TestPaymentMethod_RequiresAccountID(t *testing.T) {
	tests := []struct {
		name   string
		method PaymentMethod
		want   bool
	}{
		{
			name:   "cash does not require account ID",
			method: PaymentMethodCash,
			want:   false,
		},
		{
			name:   "bank_account requires account ID",
			method: PaymentMethodBankAccount,
			want:   true,
		},
		{
			name:   "credit_card requires account ID",
			method: PaymentMethodCreditCard,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method.RequiresAccountID()
			assert.Equal(t, tt.want, got, "PaymentMethod.RequiresAccountID() should return %v for %s", tt.want, tt.method)
		})
	}
}

// TestPaymentMethod_ToSourceType 測試付款方式轉換為 SourceType
func TestPaymentMethod_ToSourceType(t *testing.T) {
	tests := []struct {
		name   string
		method PaymentMethod
		want   SourceType
	}{
		{
			name:   "cash to SourceTypeCash",
			method: PaymentMethodCash,
			want:   SourceTypeCash,
		},
		{
			name:   "bank_account to SourceTypeBankAccount",
			method: PaymentMethodBankAccount,
			want:   SourceTypeBankAccount,
		},
		{
			name:   "credit_card to SourceTypeCreditCard",
			method: PaymentMethodCreditCard,
			want:   SourceTypeCreditCard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method.ToSourceType()
			assert.Equal(t, tt.want, got, "PaymentMethod.ToSourceType() should return %v for %s", tt.want, tt.method)
		})
	}
}

