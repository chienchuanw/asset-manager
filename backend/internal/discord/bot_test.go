package discord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
)

func TestChannelFilter_Allowed(t *testing.T) {
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{ChannelID: "111", Author: &discordgo.User{ID: "user-1"}}}
	allowedChannels := map[string]bool{"111": true}

	require.True(t, shouldProcessMessage(msg, "bot-1", allowedChannels))
}

func TestChannelFilter_Denied(t *testing.T) {
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{ChannelID: "999", Author: &discordgo.User{ID: "user-1"}}}
	allowedChannels := map[string]bool{"111": true}

	require.False(t, shouldProcessMessage(msg, "bot-1", allowedChannels))
}

func TestChannelFilter_BotMessage(t *testing.T) {
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{ChannelID: "111", Author: &discordgo.User{ID: "bot-1"}}}
	allowedChannels := map[string]bool{"111": true}

	require.False(t, shouldProcessMessage(msg, "bot-1", allowedChannels))
}

func TestChannelFilter_EmptyWhitelist(t *testing.T) {
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{ChannelID: "111", Author: &discordgo.User{ID: "user-1"}}}

	require.False(t, shouldProcessMessage(msg, "bot-1", map[string]bool{}))
}

func TestNewBot_ValidConfig(t *testing.T) {
	bot, err := NewBot(Config{
		Token:      "test-token",
		ChannelIDs: []string{"111", "222"},
		Lang:       "en",
	})

	require.NoError(t, err)
	require.NotNil(t, bot)
	require.NotNil(t, bot.session)
	require.Equal(t, map[string]bool{"111": true, "222": true}, bot.allowedChannels)
	require.Equal(t, "en", bot.lang)
}

func TestNewBot_EmptyToken(t *testing.T) {
	bot, err := NewBot(Config{})

	require.Error(t, err)
	require.Nil(t, bot)
}

func TestParseChannelIDs(t *testing.T) {
	allowedChannels := parseChannelIDs("111,222,333")

	require.Len(t, allowedChannels, 3)
	require.Equal(t, map[string]bool{"111": true, "222": true, "333": true}, allowedChannels)
}
