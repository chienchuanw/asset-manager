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
	MsgCCPaymentPreview           MsgKey = "cc_payment_preview"
	MsgCCPaymentSuccess           MsgKey = "cc_payment_success"
	MsgCCPaymentFailed            MsgKey = "cc_payment_failed"
	MsgCCPaymentConfirmButton     MsgKey = "cc_payment_confirm_button"
	MsgCCPaymentTypeFull          MsgKey = "cc_payment_type_full"
	MsgCCPaymentTypeMinimum       MsgKey = "cc_payment_type_minimum"
	MsgCCPaymentTypeCustom        MsgKey = "cc_payment_type_custom"
	MsgCCPaymentFieldCard         MsgKey = "cc_payment_field_card"
	MsgCCPaymentFieldBank         MsgKey = "cc_payment_field_bank"
	MsgCCPaymentFieldType         MsgKey = "cc_payment_field_type"
	MsgCCPaymentMissingAmount     MsgKey = "cc_payment_missing_amount"
	MsgCCPaymentUsageExamples     MsgKey = "cc_payment_usage_examples"
	MsgCCPaymentNoCards           MsgKey = "cc_payment_no_cards"
	MsgCCPaymentNoBankAccounts    MsgKey = "cc_payment_no_bank_accounts"
	MsgCCPaymentSelectCard        MsgKey = "cc_payment_select_card"
	MsgCCPaymentSelectBank        MsgKey = "cc_payment_select_bank"
	MsgUnsupported                MsgKey = "unsupported"
	MsgChatGreeting               MsgKey = "chat_greeting"
	MsgDataLoadFailed             MsgKey = "data_load_failed"
	MsgParseFailed                MsgKey = "parse_failed"
	MsgQueryFailed                MsgKey = "query_failed"
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
		MsgCCPaymentPreview:           "💳 繳款預覽",
		MsgCCPaymentSuccess:           "✅ 繳款成功",
		MsgCCPaymentFailed:            "⚠️ 繳款失敗",
		MsgCCPaymentConfirmButton:     "✅ 確認繳款",
		MsgCCPaymentTypeFull:          "全額繳款",
		MsgCCPaymentTypeMinimum:       "最低應繳",
		MsgCCPaymentTypeCustom:        "自訂金額",
		MsgCCPaymentFieldCard:         "信用卡",
		MsgCCPaymentFieldBank:         "扣款帳戶",
		MsgCCPaymentFieldType:         "繳款類型",
		MsgCCPaymentMissingAmount:     "❓ 請提供繳款金額，或指定「全額」繳款",
		MsgCCPaymentUsageExamples:     "使用範例：\n• 繳中信卡 15000\n• 繳玉山卡全額\n• 繳中信卡最低 3000",
		MsgCCPaymentNoCards:           "❓ 尚未建立信用卡，請先在系統中新增",
		MsgCCPaymentNoBankAccounts:    "❓ 尚未建立銀行帳戶，請先在系統中新增",
		MsgCCPaymentSelectCard:        "請選擇要繳款的信用卡",
		MsgCCPaymentSelectBank:        "請選擇扣款銀行帳戶",
		MsgUnsupported:                "❓ 目前不支援這項操作\n\n可用功能：\n• 記帳（收入/支出）：午餐 150\n• 信用卡繳款：繳中信卡 15000\n• 查詢月度摘要：這個月花了多少？\n• 查詢帳戶餘額：我的餘額多少？",
		MsgChatGreeting:               "👋 嗨！我是記帳小幫手，你可以告訴我消費紀錄，我會幫你記帳喔！\n\n試試看：\n• 午餐 150\n• 繳中信卡全額\n• 這個月花了多少？",
		MsgDataLoadFailed:             "⚠️ 無法載入資料，請稍後再試",
		MsgParseFailed:                "⚠️ 訊息解析失敗，請稍後再試或直接使用網頁版記帳",
		MsgQueryFailed:                "⚠️ 查詢失敗，請稍後再試",
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
		MsgCCPaymentPreview:           "💳 Payment Preview",
		MsgCCPaymentSuccess:           "✅ Payment Successful",
		MsgCCPaymentFailed:            "⚠️ Payment Failed",
		MsgCCPaymentConfirmButton:     "✅ Confirm Payment",
		MsgCCPaymentTypeFull:          "Full Payment",
		MsgCCPaymentTypeMinimum:       "Minimum Payment",
		MsgCCPaymentTypeCustom:        "Custom Amount",
		MsgCCPaymentFieldCard:         "Credit Card",
		MsgCCPaymentFieldBank:         "Bank Account",
		MsgCCPaymentFieldType:         "Payment Type",
		MsgCCPaymentMissingAmount:     "❓ Please provide a payment amount, or specify \"full\" payment",
		MsgCCPaymentUsageExamples:     "Examples:\n• pay credit card 15000\n• pay credit card in full\n• pay credit card minimum 3000",
		MsgCCPaymentNoCards:           "❓ No credit cards found. Please add one in the system first.",
		MsgCCPaymentNoBankAccounts:    "❓ No bank accounts found. Please add one in the system first.",
		MsgCCPaymentSelectCard:        "Select credit card to pay",
		MsgCCPaymentSelectBank:        "Select bank account for payment",
		MsgUnsupported:                "❓ This operation is not supported yet\n\nAvailable features:\n• Record income/expense: lunch 150\n• Credit card payment: pay credit card 15000\n• Monthly summary: how much did I spend this month?\n• Account balances: what's my balance?",
		MsgChatGreeting:               "👋 Hi! I'm your bookkeeping assistant. Tell me about your expenses and I'll help you track them!\n\nTry:\n• lunch 150\n• pay credit card in full\n• how much did I spend this month?",
		MsgDataLoadFailed:             "⚠️ Unable to load data, please try again later",
		MsgParseFailed:                "⚠️ Failed to parse message. Please try again later or use the web app.",
		MsgQueryFailed:                "⚠️ Query failed, please try again later",
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
