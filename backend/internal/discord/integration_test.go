package discord

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_MessageToConfirmToCreate(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		Description:   "ramen lunch",
		CategoryID:    "cat-food",
		CategoryName:  "飲食",
		Date:          "2026-04-05",
		MissingFields: []string{},
	}}
	creator := &mockCashFlowCreator{resultID: "cf-001"}
	loader := &mockCategoryLoader{categories: []CategoryInfo{
		{ID: "cat-food", Name: "飲食", Type: "expense"},
	}}
	h := NewHandler(parser, creator, loader, string(LangZhTW))

	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID: "msg-1", ChannelID: "ch-1", Content: "午餐拉麵 180",
		Author: &discordgo.User{ID: "user-1"},
	}}
	h.handleMessage(session, msg)

	require.True(t, parser.parseCalled)
	require.Len(t, session.sentMessages, 1)
	sent := session.sentMessages[0]
	require.Len(t, sent.Embeds, 1)
	assert.Equal(t, GetMessage(string(LangZhTW), MsgPreviewTitle), sent.Embeds[0].Title)
	assert.Equal(t, 0xFF0000, sent.Embeds[0].Color)

	require.Len(t, sent.Components, 1)
	row, ok := sent.Components[0].(discordgo.ActionsRow)
	require.True(t, ok)
	require.Len(t, row.Components, 2)

	confirmBtn := row.Components[0].(*discordgo.Button)
	require.True(t, strings.HasPrefix(confirmBtn.CustomID, "confirm:"))

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type:   discordgo.InteractionMessageComponent,
			Data:   discordgo.MessageComponentInteractionData{CustomID: confirmBtn.CustomID},
			Member: &discordgo.Member{User: &discordgo.User{ID: "user-1"}},
			Message: &discordgo.Message{
				Embeds: sent.Embeds,
			},
		},
	}
	h.handleInteraction(session, interaction)

	require.Len(t, creator.createdInputs, 1)
	created := creator.createdInputs[0]
	assert.Equal(t, "expense", created.Type)
	assert.Equal(t, 180.0, created.Amount)
	assert.Equal(t, "cat-food", created.CategoryID)

	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	assert.Equal(t, discordgo.InteractionResponseUpdateMessage, resp.Type)
	require.Len(t, resp.Data.Embeds, 1)
	assert.Equal(t, GetMessage(string(LangZhTW), MsgConfirmSuccess), resp.Data.Embeds[0].Title)
}

func TestIntegration_MessageToCancelFlow(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        50,
		Description:   "coffee",
		CategoryID:    "cat-food",
		CategoryName:  "飲食",
		Date:          "2026-04-05",
		MissingFields: []string{},
	}}
	creator := &mockCashFlowCreator{}
	loader := &mockCategoryLoader{categories: []CategoryInfo{
		{ID: "cat-food", Name: "飲食", Type: "expense"},
	}}
	h := NewHandler(parser, creator, loader, string(LangZhTW))

	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID: "msg-2", ChannelID: "ch-1", Content: "咖啡 50",
		Author: &discordgo.User{ID: "user-2"},
	}}
	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)

	cancelInteraction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type:   discordgo.InteractionMessageComponent,
			Data:   discordgo.MessageComponentInteractionData{CustomID: "cancel:user-2"},
			Member: &discordgo.Member{User: &discordgo.User{ID: "user-2"}},
			Message: &discordgo.Message{
				Embeds: session.sentMessages[0].Embeds,
			},
		},
	}
	h.handleInteraction(session, cancelInteraction)

	assert.Empty(t, creator.createdInputs)

	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Len(t, resp.Data.Embeds, 1)
	assert.Equal(t, GetMessage(string(LangZhTW), MsgCancelled), resp.Data.Embeds[0].Title)
}

func TestStorePending_CustomIDFormat(t *testing.T) {
	h := NewHandler(&mockParser{}, &mockCashFlowCreator{}, nil, string(LangZhTW))
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        100,
		Description:   strings.Repeat("很長的描述", 50),
		CategoryID:    "cat-1",
		Date:          "2026-04-05",
	}

	customID := h.storePending(result, "user-1")
	assert.True(t, strings.HasPrefix(customID, "confirm:"))
	assert.True(t, strings.HasSuffix(customID, ":user-1"))
	assert.LessOrEqual(t, len(customID), 100)
}

func TestPopPending_RetrievesAndRemoves(t *testing.T) {
	h := NewHandler(&mockParser{}, &mockCashFlowCreator{}, nil, string(LangZhTW))
	result := &ParseResult{Type: "expense", Amount: 42}

	customID := h.storePending(result, "user-1")
	parts := strings.Split(customID, ":")
	require.Len(t, parts, 3)
	key := parts[1]

	got, ok := h.popPending(key)
	require.True(t, ok)
	assert.Equal(t, 42.0, got.Amount)

	_, ok = h.popPending(key)
	assert.False(t, ok)
}

func TestFormatAmount_Commas(t *testing.T) {
	tests := []struct {
		amount   float64
		expected string
	}{
		{0, "0"},
		{100, "100"},
		{1234, "1,234"},
		{45000, "45,000"},
		{1234567, "1,234,567"},
		{99.50, "99.50"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, formatAmount(tt.amount))
		})
	}
}
