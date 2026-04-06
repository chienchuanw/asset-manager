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
	MsgPreviewTitle               MsgKey = "preview_title"
	MsgConfirmSuccess             MsgKey = "confirm_success"
	MsgCancelled                  MsgKey = "cancelled"
	MsgExpired                    MsgKey = "expired"
	MsgConfirmButton              MsgKey = "confirm_button"
	MsgCancelButton               MsgKey = "cancel_button"
	MsgTypeIncome                 MsgKey = "type_income"
	MsgTypeExpense                MsgKey = "type_expense"
	MsgFieldType                  MsgKey = "field_type"
	MsgFieldAmount                MsgKey = "field_amount"
	MsgFieldCategory              MsgKey = "field_category"
	MsgFieldDescription           MsgKey = "field_description"
	MsgFieldDate                  MsgKey = "field_date"
	MsgMissingAmount              MsgKey = "missing_amount"
	MsgSystemError                MsgKey = "system_error"
	MsgNLPNotConfigured           MsgKey = "nlp_not_configured"
	MsgBookingFailed              MsgKey = "booking_failed"
	MsgOnlyAuthor                 MsgKey = "only_author"
	MsgUsageExamples              MsgKey = "usage_examples"
	MsgSelectAccount              MsgKey = "select_account"
	MsgAccountCash                MsgKey = "account_cash"
	MsgAccountBank                MsgKey = "account_bank"
	MsgAccountCreditCard          MsgKey = "account_credit_card"
	MsgFieldPaymentMethod         MsgKey = "field_payment_method"
	MsgFieldAccount               MsgKey = "field_account"
	MsgSelectBankAccount          MsgKey = "select_bank_account"
	MsgSelectCreditCard           MsgKey = "select_credit_card"
	MsgNoAccountsFound            MsgKey = "no_accounts_found"
	MsgQueryCashFlowTitle         MsgKey = "query_cash_flow_title"
	MsgQueryCashFlowCategoryTitle MsgKey = "query_cash_flow_category_title"
	MsgQueryAccountTitle          MsgKey = "query_account_title"
	MsgQueryNoData                MsgKey = "query_no_data"
	MsgQueryUnsupported           MsgKey = "query_unsupported"
	MsgQueryBankSection           MsgKey = "query_bank_section"
	MsgQueryCCSection             MsgKey = "query_cc_section"
	MsgQueryCCNearLimit           MsgKey = "query_cc_near_limit"
	MsgQueryNoAccounts            MsgKey = "query_no_accounts"
	MsgQueryTotalIncome           MsgKey = "query_total_income"
	MsgQueryTotalExpense          MsgKey = "query_total_expense"
	MsgQueryNetCashFlow           MsgKey = "query_net_cash_flow"
	MsgQueryFieldCount            MsgKey = "query_field_count"
	MsgQueryTopCategories         MsgKey = "query_top_categories"
	MsgQueryComparison            MsgKey = "query_comparison"
	MsgQueryCCLimit               MsgKey = "query_cc_limit"
	MsgQueryCCUsed                MsgKey = "query_cc_used"
	MsgQueryCCRemaining           MsgKey = "query_cc_remaining"
	MsgQueryBankTotal             MsgKey = "query_bank_total"
	MsgQueryCategoryNotFound      MsgKey = "query_category_not_found"
	MsgQueryLoadFailed            MsgKey = "query_load_failed"
)

var messageCatalog = map[Lang]map[MsgKey]string{
	LangZhTW: {
		MsgPreviewTitle:               "📝 記帳預覽",
		MsgConfirmSuccess:             "✅ 記帳成功",
		MsgCancelled:                  "❌ 已取消",
		MsgExpired:                    "⏰ 已過期",
		MsgConfirmButton:              "✅ 確認記帳",
		MsgCancelButton:               "❌ 取消",
		MsgTypeIncome:                 "收入",
		MsgTypeExpense:                "支出",
		MsgFieldType:                  "類型",
		MsgFieldAmount:                "金額",
		MsgFieldCategory:              "分類",
		MsgFieldDescription:           "描述",
		MsgFieldDate:                  "日期",
		MsgMissingAmount:              "❓ 無法辨識金額，請提供金額資訊",
		MsgSystemError:                "⚠️ 系統錯誤，請稍後再試",
		MsgNLPNotConfigured:           "⚠️ NLP 服務未設定",
		MsgBookingFailed:              "⚠️ 記帳失敗",
		MsgOnlyAuthor:                 "只有原始訊息發送者可以操作",
		MsgUsageExamples:              "使用範例：\n• 午餐 150\n• 搭捷運 35\n• 收到薪水 45000",
		MsgSelectAccount:              "請選擇付款方式",
		MsgAccountCash:                "現金",
		MsgAccountBank:                "銀行帳戶",
		MsgAccountCreditCard:          "信用卡",
		MsgFieldPaymentMethod:         "付款方式",
		MsgFieldAccount:               "帳戶",
		MsgSelectBankAccount:          "請選擇銀行帳戶",
		MsgSelectCreditCard:           "請選擇信用卡",
		MsgNoAccountsFound:            "❓ 找不到對應的帳戶，請先在系統中新增",
		MsgQueryCashFlowTitle:         "📊 %d月現金流摘要",
		MsgQueryCashFlowCategoryTitle: "📊 %d月%s支出",
		MsgQueryAccountTitle:          "💰 帳戶餘額",
		MsgQueryNoData:                "本月尚無記錄",
		MsgQueryUnsupported:           "❓ 目前支援的查詢：月度支出摘要、帳戶餘額",
		MsgQueryBankSection:           "🏦 銀行帳戶",
		MsgQueryCCSection:             "💳 信用卡",
		MsgQueryCCNearLimit:           "⚠️ 額度即將用盡",
		MsgQueryNoAccounts:            "❓ 尚未建立任何帳戶，請先在系統中新增",
		MsgQueryTotalIncome:           "總收入",
		MsgQueryTotalExpense:          "總支出",
		MsgQueryNetCashFlow:           "淨現金流",
		MsgQueryFieldCount:            "筆數",
		MsgQueryTopCategories:         "📋 支出分類 TOP 5",
		MsgQueryComparison:            "vs 上月",
		MsgQueryCCLimit:               "額度",
		MsgQueryCCUsed:                "已用",
		MsgQueryCCRemaining:           "剩餘",
		MsgQueryBankTotal:             "銀行帳戶合計",
		MsgQueryCategoryNotFound:      "❓ 找不到「%s」分類，請確認分類名稱",
		MsgQueryLoadFailed:            "載入失敗",
	},
	LangEn: {
		MsgPreviewTitle:               "📝 Booking Preview",
		MsgConfirmSuccess:             "✅ Booked Successfully",
		MsgCancelled:                  "❌ Cancelled",
		MsgExpired:                    "⏰ Expired",
		MsgConfirmButton:              "✅ Confirm",
		MsgCancelButton:               "❌ Cancel",
		MsgTypeIncome:                 "Income",
		MsgTypeExpense:                "Expense",
		MsgFieldType:                  "Type",
		MsgFieldAmount:                "Amount",
		MsgFieldCategory:              "Category",
		MsgFieldDescription:           "Description",
		MsgFieldDate:                  "Date",
		MsgMissingAmount:              "❓ Could not identify the amount. Please include an amount.",
		MsgSystemError:                "⚠️ System error, please try again later",
		MsgNLPNotConfigured:           "⚠️ NLP service not configured",
		MsgBookingFailed:              "⚠️ Booking Failed",
		MsgOnlyAuthor:                 "Only the original message author can interact",
		MsgUsageExamples:              "Examples:\n• lunch 150\n• taxi 35\n• salary received 45000",
		MsgSelectAccount:              "Select payment method",
		MsgAccountCash:                "Cash",
		MsgAccountBank:                "Bank Account",
		MsgAccountCreditCard:          "Credit Card",
		MsgFieldPaymentMethod:         "Payment Method",
		MsgFieldAccount:               "Account",
		MsgSelectBankAccount:          "Select bank account",
		MsgSelectCreditCard:           "Select credit card",
		MsgNoAccountsFound:            "❓ No accounts found. Please add one in the system first.",
		MsgQueryCashFlowTitle:         "📊 %s Cash Flow Summary",
		MsgQueryCashFlowCategoryTitle: "📊 %s %s Spending",
		MsgQueryAccountTitle:          "💰 Account Balances",
		MsgQueryNoData:                "No records this month",
		MsgQueryUnsupported:           "❓ Supported queries: monthly spending summary, account balances",
		MsgQueryBankSection:           "🏦 Bank Accounts",
		MsgQueryCCSection:             "💳 Credit Cards",
		MsgQueryCCNearLimit:           "⚠️ Nearing limit",
		MsgQueryNoAccounts:            "❓ No accounts found. Please add one in the system first.",
		MsgQueryTotalIncome:           "Total Income",
		MsgQueryTotalExpense:          "Total Expense",
		MsgQueryNetCashFlow:           "Net Cash Flow",
		MsgQueryFieldCount:            "Count",
		MsgQueryTopCategories:         "📋 Top 5 Expense Categories",
		MsgQueryComparison:            "vs Last Month",
		MsgQueryCCLimit:               "Limit",
		MsgQueryCCUsed:                "Used",
		MsgQueryCCRemaining:           "Remaining",
		MsgQueryBankTotal:             "Bank Total",
		MsgQueryCategoryNotFound:      "❓ Category \"%s\" not found. Please check the name.",
		MsgQueryLoadFailed:            "Load failed",
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
