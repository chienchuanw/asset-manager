package discord

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMessage_ZhTW(t *testing.T) {
	require.Equal(t, "📝 記帳預覽", GetMessage(string(LangZhTW), MsgPreviewTitle))
}

func TestGetMessage_En(t *testing.T) {
	require.Equal(t, "📝 Booking Preview", GetMessage(string(LangEn), MsgPreviewTitle))
}

func TestGetMessage_UnknownKey(t *testing.T) {
	const unknownKey MsgKey = "unknown.key"

	require.Equal(t, string(unknownKey), GetMessage(string(LangZhTW), unknownKey))
}

func TestGetMessage_UnknownLang(t *testing.T) {
	require.Equal(t, "✅ 記帳成功", GetMessage("ja", MsgConfirmSuccess))
}

func TestGetMessage_AllKeysHaveBothLanguages(t *testing.T) {
	zhMessages, ok := messageCatalog[LangZhTW]
	require.True(t, ok)

	enMessages, ok := messageCatalog[LangEn]
	require.True(t, ok)

	for key, zhMessage := range zhMessages {
		require.NotEmpty(t, zhMessage)
		require.Contains(t, enMessages, key)
		require.NotEmpty(t, enMessages[key])
	}
}

func TestGetMessage_AccountSelectionKeys(t *testing.T) {
	accountKeys := []MsgKey{
		MsgSelectAccount,
		MsgAccountCash,
		MsgAccountBank,
		MsgAccountCreditCard,
		MsgFieldPaymentMethod,
	}

	for _, key := range accountKeys {
		zhMsg := GetMessage(string(LangZhTW), key)
		require.NotEqual(t, string(key), zhMsg, "zh-TW missing key: %s", key)

		enMsg := GetMessage(string(LangEn), key)
		require.NotEqual(t, string(key), enMsg, "en missing key: %s", key)
	}
}

func TestGetMessage_AccountLabels_ZhTW(t *testing.T) {
	require.Equal(t, "現金", GetMessage(string(LangZhTW), MsgAccountCash))
	require.Equal(t, "銀行帳戶", GetMessage(string(LangZhTW), MsgAccountBank))
	require.Equal(t, "信用卡", GetMessage(string(LangZhTW), MsgAccountCreditCard))
	require.Equal(t, "請選擇付款方式", GetMessage(string(LangZhTW), MsgSelectAccount))
	require.Equal(t, "付款方式", GetMessage(string(LangZhTW), MsgFieldPaymentMethod))
}

func TestGetMessage_AccountLabels_En(t *testing.T) {
	require.Equal(t, "Cash", GetMessage(string(LangEn), MsgAccountCash))
	require.Equal(t, "Bank Account", GetMessage(string(LangEn), MsgAccountBank))
	require.Equal(t, "Credit Card", GetMessage(string(LangEn), MsgAccountCreditCard))
	require.Equal(t, "Select payment method", GetMessage(string(LangEn), MsgSelectAccount))
	require.Equal(t, "Payment Method", GetMessage(string(LangEn), MsgFieldPaymentMethod))
}

func TestGetMessage_NewQueryKeys(t *testing.T) {
	queryKeys := []MsgKey{
		MsgQueryCashFlowTitle,
		MsgQueryCashFlowCategoryTitle,
		MsgQueryAccountTitle,
		MsgQueryNoData,
		MsgQueryUnsupported,
		MsgQueryBankSection,
		MsgQueryCCSection,
		MsgQueryCCNearLimit,
		MsgQueryNoAccounts,
		MsgQueryTotalIncome,
		MsgQueryTotalExpense,
		MsgQueryNetCashFlow,
		MsgQueryFieldCount,
		MsgQueryTopCategories,
		MsgQueryComparison,
		MsgQueryCCLimit,
		MsgQueryCCUsed,
		MsgQueryCCRemaining,
		MsgQueryBankTotal,
		MsgQueryCategoryNotFound,
		MsgQueryLoadFailed,
	}

	for _, key := range queryKeys {
		zhMsg := GetMessage(string(LangZhTW), key)
		require.NotEmpty(t, zhMsg)
		require.NotEqual(t, string(key), zhMsg, "zh-TW missing key: %s", key)

		enMsg := GetMessage(string(LangEn), key)
		require.NotEmpty(t, enMsg)
		require.NotEqual(t, string(key), enMsg, "en missing key: %s", key)
	}
}
