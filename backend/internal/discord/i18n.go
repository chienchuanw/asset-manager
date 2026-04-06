package discord

// Lang represents a supported Discord bot language.
type Lang string

const (
	LangZhTW Lang = "zh-TW"
	LangEn   Lang = "en"
)

// MsgKey identifies a localized Discord bot message.
type MsgKey string

const (
	MsgPreviewTitle       MsgKey = "preview_title"
	MsgConfirmSuccess     MsgKey = "confirm_success"
	MsgCancelled          MsgKey = "cancelled"
	MsgExpired            MsgKey = "expired"
	MsgConfirmButton      MsgKey = "confirm_button"
	MsgCancelButton       MsgKey = "cancel_button"
	MsgTypeIncome         MsgKey = "type_income"
	MsgTypeExpense        MsgKey = "type_expense"
	MsgFieldType          MsgKey = "field_type"
	MsgFieldAmount        MsgKey = "field_amount"
	MsgFieldCategory      MsgKey = "field_category"
	MsgFieldDescription   MsgKey = "field_description"
	MsgFieldDate          MsgKey = "field_date"
	MsgMissingAmount      MsgKey = "missing_amount"
	MsgSystemError        MsgKey = "system_error"
	MsgNLPNotConfigured   MsgKey = "nlp_not_configured"
	MsgBookingFailed      MsgKey = "booking_failed"
	MsgOnlyAuthor         MsgKey = "only_author"
	MsgUsageExamples      MsgKey = "usage_examples"
	MsgSelectAccount      MsgKey = "select_account"
	MsgAccountCash        MsgKey = "account_cash"
	MsgAccountBank        MsgKey = "account_bank"
	MsgAccountCreditCard  MsgKey = "account_credit_card"
	MsgFieldPaymentMethod MsgKey = "field_payment_method"
	MsgSelectBankAccount  MsgKey = "select_bank_account"
	MsgSelectCreditCard   MsgKey = "select_credit_card"
	MsgNoAccountsFound    MsgKey = "no_accounts_found"
)

var messageCatalog = map[Lang]map[MsgKey]string{
	LangZhTW: {
		MsgPreviewTitle:       "📝 記帳預覽",
		MsgConfirmSuccess:     "✅ 記帳成功",
		MsgCancelled:          "❌ 已取消",
		MsgExpired:            "⏰ 已過期",
		MsgConfirmButton:      "✅ 確認記帳",
		MsgCancelButton:       "❌ 取消",
		MsgTypeIncome:         "收入",
		MsgTypeExpense:        "支出",
		MsgFieldType:          "類型",
		MsgFieldAmount:        "金額",
		MsgFieldCategory:      "分類",
		MsgFieldDescription:   "描述",
		MsgFieldDate:          "日期",
		MsgMissingAmount:      "❓ 無法辨識金額，請提供金額資訊",
		MsgSystemError:        "⚠️ 系統錯誤，請稍後再試",
		MsgNLPNotConfigured:   "⚠️ NLP 服務未設定",
		MsgBookingFailed:      "⚠️ 記帳失敗",
		MsgOnlyAuthor:         "只有原始訊息發送者可以操作",
		MsgUsageExamples:      "使用範例：\n• 午餐 150\n• 搭捷運 35\n• 收到薪水 45000",
		MsgSelectAccount:      "請選擇付款方式",
		MsgAccountCash:        "現金",
		MsgAccountBank:        "銀行帳戶",
		MsgAccountCreditCard:  "信用卡",
		MsgFieldPaymentMethod: "付款方式",
		MsgSelectBankAccount:  "請選擇銀行帳戶",
		MsgSelectCreditCard:   "請選擇信用卡",
		MsgNoAccountsFound:    "❓ 找不到對應的帳戶，請先在系統中新增",
	},
	LangEn: {
		MsgPreviewTitle:       "📝 Booking Preview",
		MsgConfirmSuccess:     "✅ Booked Successfully",
		MsgCancelled:          "❌ Cancelled",
		MsgExpired:            "⏰ Expired",
		MsgConfirmButton:      "✅ Confirm",
		MsgCancelButton:       "❌ Cancel",
		MsgTypeIncome:         "Income",
		MsgTypeExpense:        "Expense",
		MsgFieldType:          "Type",
		MsgFieldAmount:        "Amount",
		MsgFieldCategory:      "Category",
		MsgFieldDescription:   "Description",
		MsgFieldDate:          "Date",
		MsgMissingAmount:      "❓ Could not identify the amount. Please include an amount.",
		MsgSystemError:        "⚠️ System error, please try again later",
		MsgNLPNotConfigured:   "⚠️ NLP service not configured",
		MsgBookingFailed:      "⚠️ Booking Failed",
		MsgOnlyAuthor:         "Only the original message author can interact",
		MsgUsageExamples:      "Examples:\n• lunch 150\n• taxi 35\n• salary received 45000",
		MsgSelectAccount:      "Select payment method",
		MsgAccountCash:        "Cash",
		MsgAccountBank:        "Bank Account",
		MsgAccountCreditCard:  "Credit Card",
		MsgFieldPaymentMethod: "Payment Method",
		MsgSelectBankAccount:  "Select bank account",
		MsgSelectCreditCard:   "Select credit card",
		MsgNoAccountsFound:    "❓ No accounts found. Please add one in the system first.",
	},
}

// GetMessage returns the localized message for the given language and key.
func GetMessage(lang string, key MsgKey) string {
	messages, ok := messageCatalog[Lang(lang)]
	if !ok {
		messages = messageCatalog[LangZhTW]
	}

	message, ok := messages[key]
	if !ok {
		return string(key)
	}

	return message
}
