package models

// PaymentMethod 付款方式
// 用於指定訂閱或分期的扣款來源
type PaymentMethod string

const (
	PaymentMethodCash        PaymentMethod = "cash"         // 現金
	PaymentMethodBankAccount PaymentMethod = "bank_account" // 銀行帳戶
	PaymentMethodCreditCard  PaymentMethod = "credit_card"  // 信用卡
)

// Validate 驗證 PaymentMethod 是否有效
func (p PaymentMethod) Validate() bool {
	switch p {
	case PaymentMethodCash, PaymentMethodBankAccount, PaymentMethodCreditCard:
		return true
	}
	return false
}

// RequiresAccountID 檢查此付款方式是否需要指定帳戶 ID
// 現金不需要帳戶 ID，銀行帳戶和信用卡需要
func (p PaymentMethod) RequiresAccountID() bool {
	switch p {
	case PaymentMethodBankAccount, PaymentMethodCreditCard:
		return true
	}
	return false
}

// ToSourceType 將 PaymentMethod 轉換為對應的 SourceType
// 用於在產生現金流記錄時設定正確的來源類型
func (p PaymentMethod) ToSourceType() SourceType {
	switch p {
	case PaymentMethodCash:
		return SourceTypeCash
	case PaymentMethodBankAccount:
		return SourceTypeBankAccount
	case PaymentMethodCreditCard:
		return SourceTypeCreditCard
	default:
		return SourceTypeCash // 預設為現金
	}
}

