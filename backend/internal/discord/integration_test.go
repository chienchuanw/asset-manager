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
		SourceType:    "cash",
		MissingFields: []string{},
	}}
	creator := &mockCashFlowCreator{resultID: "cf-001"}
	loader := &mockCategoryLoader{categories: []CategoryInfo{
		{ID: "cat-food", Name: "飲食", Type: "expense"},
	}}
	h := NewHandler(parser, creator, loader, nil, string(LangZhTW))

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
		SourceType:    "cash",
		MissingFields: []string{},
	}}
	creator := &mockCashFlowCreator{}
	loader := &mockCategoryLoader{categories: []CategoryInfo{
		{ID: "cat-food", Name: "飲食", Type: "expense"},
	}}
	h := NewHandler(parser, creator, loader, nil, string(LangZhTW))

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
	h := NewHandler(&mockParser{}, &mockCashFlowCreator{}, nil, nil, string(LangZhTW))
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
	h := NewHandler(&mockParser{}, &mockCashFlowCreator{}, nil, nil, string(LangZhTW))
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

func TestIntegration_YesterdayDate_PassedToRecord(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		Description:   "午餐",
		CategoryID:    "cat-food",
		CategoryName:  "飲食",
		Date:          "2026-04-04",
		SourceType:    "cash",
		MissingFields: []string{},
	}}
	creator := &mockCashFlowCreator{resultID: "cf-date-1"}
	loader := &mockCategoryLoader{categories: []CategoryInfo{
		{ID: "cat-food", Name: "飲食", Type: "expense"},
	}}
	h := NewHandler(parser, creator, loader, nil, string(LangZhTW))

	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID: "msg-date-1", ChannelID: "ch-1", Content: "昨天午餐 180",
		Author: &discordgo.User{ID: "user-1"},
	}}
	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	sent := session.sentMessages[0]
	require.Len(t, sent.Embeds, 1)
	assert.Equal(t, "2026-04-04", sent.Embeds[0].Fields[4].Value)

	confirmBtn := sent.Components[0].(discordgo.ActionsRow).Components[0].(*discordgo.Button)
	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Data:    discordgo.MessageComponentInteractionData{CustomID: confirmBtn.CustomID},
		Member:  &discordgo.Member{User: &discordgo.User{ID: "user-1"}},
		Message: &discordgo.Message{Embeds: sent.Embeds},
	}}
	h.handleInteraction(session, interaction)

	require.Len(t, creator.createdInputs, 1)
	assert.Equal(t, "2026-04-04", creator.createdInputs[0].Date)
}

func TestIntegration_SpecificDate_PassedToRecord(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		Description:   "午餐",
		CategoryID:    "cat-food",
		CategoryName:  "飲食",
		Date:          "2026-04-03",
		SourceType:    "cash",
		MissingFields: []string{},
	}}
	creator := &mockCashFlowCreator{resultID: "cf-date-2"}
	loader := &mockCategoryLoader{categories: []CategoryInfo{
		{ID: "cat-food", Name: "飲食", Type: "expense"},
	}}
	h := NewHandler(parser, creator, loader, nil, string(LangZhTW))

	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID: "msg-date-2", ChannelID: "ch-1", Content: "4/3 午餐 180",
		Author: &discordgo.User{ID: "user-1"},
	}}
	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	sent := session.sentMessages[0]
	assert.Equal(t, "2026-04-03", sent.Embeds[0].Fields[4].Value)

	confirmBtn := sent.Components[0].(discordgo.ActionsRow).Components[0].(*discordgo.Button)
	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Data:    discordgo.MessageComponentInteractionData{CustomID: confirmBtn.CustomID},
		Member:  &discordgo.Member{User: &discordgo.User{ID: "user-1"}},
		Message: &discordgo.Message{Embeds: sent.Embeds},
	}}
	h.handleInteraction(session, interaction)

	require.Len(t, creator.createdInputs, 1)
	assert.Equal(t, "2026-04-03", creator.createdInputs[0].Date)
}

func TestIntegration_DefaultTodayDate(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		Description:   "午餐吃拉麵",
		CategoryID:    "cat-food",
		CategoryName:  "飲食",
		Date:          "2026-04-05",
		SourceType:    "cash",
		MissingFields: []string{},
	}}
	creator := &mockCashFlowCreator{resultID: "cf-date-3"}
	loader := &mockCategoryLoader{categories: []CategoryInfo{
		{ID: "cat-food", Name: "飲食", Type: "expense"},
	}}
	h := NewHandler(parser, creator, loader, nil, string(LangZhTW))

	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID: "msg-date-3", ChannelID: "ch-1", Content: "午餐 180",
		Author: &discordgo.User{ID: "user-1"},
	}}
	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	assert.Equal(t, "2026-04-05", session.sentMessages[0].Embeds[0].Fields[4].Value)

	confirmBtn := session.sentMessages[0].Components[0].(discordgo.ActionsRow).Components[0].(*discordgo.Button)
	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Data:    discordgo.MessageComponentInteractionData{CustomID: confirmBtn.CustomID},
		Member:  &discordgo.Member{User: &discordgo.User{ID: "user-1"}},
		Message: &discordgo.Message{Embeds: session.sentMessages[0].Embeds},
	}}
	h.handleInteraction(session, interaction)

	require.Len(t, creator.createdInputs, 1)
	assert.Equal(t, "2026-04-05", creator.createdInputs[0].Date)
}

func TestIntegration_SelectMenuToConfirmFlow(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		Description:   "午餐",
		CategoryID:    "cat-food",
		CategoryName:  "飲食",
		Date:          "2026-04-05",
		SourceType:    "",
		MissingFields: []string{},
	}}
	creator := &mockCashFlowCreator{resultID: "cf-select-1"}
	loader := &mockCategoryLoader{categories: []CategoryInfo{
		{ID: "cat-food", Name: "飲食", Type: "expense"},
	}}
	acctLoader := &mockAccountLoader{accounts: []AccountInfo{
		{ID: "cc-uuid-1", Name: "中信 Visa *1234", Type: "credit_card"},
	}}
	h := NewHandler(parser, creator, loader, acctLoader, string(LangZhTW))

	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID: "msg-select-1", ChannelID: "ch-1", Content: "午餐 180",
		Author: &discordgo.User{ID: "user-1"},
	}}
	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	sent := session.sentMessages[0]
	require.Empty(t, sent.Embeds)
	require.Len(t, sent.Components, 1)
	row := sent.Components[0].(discordgo.ActionsRow)
	selectMenu := row.Components[0].(*discordgo.SelectMenu)
	require.Equal(t, GetMessage(string(LangZhTW), MsgSelectAccount), selectMenu.Placeholder)

	selectCustomID := selectMenu.CustomID
	selectInteraction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "user-1"}},
		Message: &discordgo.Message{ID: "reply-select", ChannelID: "ch-1"},
		Data: discordgo.MessageComponentInteractionData{
			CustomID:      selectCustomID,
			ComponentType: discordgo.SelectMenuComponent,
			Values:        []string{"credit_card"},
		},
	}}
	h.handleInteraction(session, selectInteraction)

	require.Len(t, session.interactionResponses, 1)
	acctResp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseUpdateMessage, acctResp.Type)
	require.Len(t, acctResp.Data.Components, 1)
	acctRow := acctResp.Data.Components[0].(discordgo.ActionsRow)
	acctMenu := acctRow.Components[0].(*discordgo.SelectMenu)
	require.Len(t, acctMenu.Options, 1)
	require.Equal(t, "cc-uuid-1", acctMenu.Options[0].Value)

	acctInteraction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "user-1"}},
		Message: &discordgo.Message{ID: "reply-acct", ChannelID: "ch-1"},
		Data: discordgo.MessageComponentInteractionData{
			CustomID:      acctMenu.CustomID,
			ComponentType: discordgo.SelectMenuComponent,
			Values:        []string{"cc-uuid-1"},
		},
	}}
	h.handleInteraction(session, acctInteraction)

	require.Len(t, session.interactionResponses, 2)
	previewResp := session.interactionResponses[1]
	require.Len(t, previewResp.Data.Embeds, 1)
	assert.Equal(t, GetMessage(string(LangZhTW), MsgPreviewTitle), previewResp.Data.Embeds[0].Title)

	require.Len(t, previewResp.Data.Components, 1)
	confirmRow := previewResp.Data.Components[0].(discordgo.ActionsRow)
	confirmBtn := confirmRow.Components[0].(*discordgo.Button)

	confirmInteraction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Data:    discordgo.MessageComponentInteractionData{CustomID: confirmBtn.CustomID},
		Member:  &discordgo.Member{User: &discordgo.User{ID: "user-1"}},
		Message: &discordgo.Message{Embeds: previewResp.Data.Embeds},
	}}
	h.handleInteraction(session, confirmInteraction)

	require.Len(t, creator.createdInputs, 1)
	created := creator.createdInputs[0]
	assert.Equal(t, "expense", created.Type)
	assert.Equal(t, 180.0, created.Amount)
	assert.Equal(t, "cat-food", created.CategoryID)
	assert.Equal(t, "credit_card", created.SourceType)
	assert.Equal(t, "cc-uuid-1", created.SourceID)
	assert.Equal(t, "2026-04-05", created.Date)
}

func TestIntegration_KnownSourceType_StillAsksAccountID(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        2000,
		Description:   "刷卡買衣服",
		CategoryID:    "cat-other",
		CategoryName:  "其他支出",
		Date:          "2026-04-05",
		SourceType:    "credit_card",
		MissingFields: []string{},
	}}
	creator := &mockCashFlowCreator{resultID: "cf-known-1"}
	loader := &mockCategoryLoader{categories: []CategoryInfo{
		{ID: "cat-other", Name: "其他支出", Type: "expense"},
	}}
	acctLoader := &mockAccountLoader{accounts: []AccountInfo{
		{ID: "cc-uuid-1", Name: "中信 Visa *1234", Type: "credit_card"},
	}}
	h := NewHandler(parser, creator, loader, acctLoader, string(LangZhTW))

	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID: "msg-known-1", ChannelID: "ch-1", Content: "刷卡買衣服 2000",
		Author: &discordgo.User{ID: "user-1"},
	}}
	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	sent := session.sentMessages[0]
	require.Empty(t, sent.Embeds)
	require.Len(t, sent.Components, 1)
	acctMenu := sent.Components[0].(discordgo.ActionsRow).Components[0].(*discordgo.SelectMenu)
	require.Equal(t, GetMessage(string(LangZhTW), MsgSelectCreditCard), acctMenu.Placeholder)

	acctInteraction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "user-1"}},
		Message: &discordgo.Message{ID: "reply-acct", ChannelID: "ch-1"},
		Data: discordgo.MessageComponentInteractionData{
			CustomID:      acctMenu.CustomID,
			ComponentType: discordgo.SelectMenuComponent,
			Values:        []string{"cc-uuid-1"},
		},
	}}
	h.handleInteraction(session, acctInteraction)

	require.Len(t, session.interactionResponses, 1)
	previewResp := session.interactionResponses[0]
	require.Len(t, previewResp.Data.Embeds, 1)
	assert.Equal(t, GetMessage(string(LangZhTW), MsgPreviewTitle), previewResp.Data.Embeds[0].Title)

	confirmBtn := previewResp.Data.Components[0].(discordgo.ActionsRow).Components[0].(*discordgo.Button)
	confirmInteraction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Data:    discordgo.MessageComponentInteractionData{CustomID: confirmBtn.CustomID},
		Member:  &discordgo.Member{User: &discordgo.User{ID: "user-1"}},
		Message: &discordgo.Message{Embeds: previewResp.Data.Embeds},
	}}
	h.handleInteraction(session, confirmInteraction)

	require.Len(t, creator.createdInputs, 1)
	assert.Equal(t, "credit_card", creator.createdInputs[0].SourceType)
	assert.Equal(t, "cc-uuid-1", creator.createdInputs[0].SourceID)
}

func TestAdapter_InvalidDateFallsBackToNow(t *testing.T) {
	input := &BotCashFlowInput{
		Date:        "invalid-date",
		Type:        "expense",
		CategoryID:  "not-a-uuid",
		Amount:      100,
		Description: "test",
		SourceType:  "cash",
	}

	adapter := &CashFlowServiceAdapter{}
	_, err := adapter.CreateCashFlowFromBot(input)

	require.Error(t, err)
}

func TestMapSourceType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"cash", "cash"},
		{"bank_account", "bank_account"},
		{"credit_card", "credit_card"},
		{"", "cash"},
		{"unknown", "cash"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mapSourceType(tt.input)
			assert.Equal(t, tt.expected, string(result))
		})
	}
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
