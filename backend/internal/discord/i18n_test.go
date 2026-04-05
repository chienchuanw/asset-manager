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
