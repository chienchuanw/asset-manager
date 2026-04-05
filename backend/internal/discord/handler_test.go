package discord

import (
	"context"
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
)

type mockSession struct {
	sentMessages         []*discordgo.MessageSend
	editedMessages       []*discordgo.MessageEdit
	interactionResponses []*discordgo.InteractionResponse
	reactionAdds         []reactionAddCall
	reactionRemoves      []reactionRemoveCall
}

type reactionAddCall struct {
	channelID string
	messageID string
	emojiID   string
}

type reactionRemoveCall struct {
	channelID string
	messageID string
	emojiID   string
	userID    string
}

func (m *mockSession) ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend) (*discordgo.Message, error) {
	m.sentMessages = append(m.sentMessages, data)
	return &discordgo.Message{ID: "reply-1", ChannelID: channelID}, nil
}

func (m *mockSession) MessageReactionAdd(channelID, messageID, emojiID string) error {
	m.reactionAdds = append(m.reactionAdds, reactionAddCall{channelID: channelID, messageID: messageID, emojiID: emojiID})
	return nil
}

func (m *mockSession) MessageReactionRemove(channelID, messageID, emojiID, userID string) error {
	m.reactionRemoves = append(m.reactionRemoves, reactionRemoveCall{channelID: channelID, messageID: messageID, emojiID: emojiID, userID: userID})
	return nil
}

func (m *mockSession) ChannelMessageEditComplex(data *discordgo.MessageEdit) (*discordgo.Message, error) {
	m.editedMessages = append(m.editedMessages, data)
	return &discordgo.Message{ID: data.ID, ChannelID: data.Channel}, nil
}

func (m *mockSession) InteractionRespond(_ *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	m.interactionResponses = append(m.interactionResponses, resp)
	return nil
}

type mockParser struct {
	result           *ParseResult
	err              error
	parseCalled      bool
	message          string
	categoriesPassed []CategoryInfo
}

func (m *mockParser) Parse(_ context.Context, message string, categories []CategoryInfo) (*ParseResult, error) {
	m.parseCalled = true
	m.message = message
	m.categoriesPassed = append([]CategoryInfo(nil), categories...)
	return m.result, m.err
}

type mockCashFlowCreator struct {
	createdInputs []*BotCashFlowInput
	resultID      string
	err           error
}

func (m *mockCashFlowCreator) CreateCashFlowFromBot(input *BotCashFlowInput) (string, error) {
	m.createdInputs = append(m.createdInputs, input)
	return m.resultID, m.err
}

type mockCategoryLoader struct {
	categories []CategoryInfo
	err        error
	called     bool
}

func (m *mockCategoryLoader) LoadCategories() ([]CategoryInfo, error) {
	m.called = true
	return m.categories, m.err
}

func TestHandler_ImplementsMessageHandler(t *testing.T) {
	var _ MessageHandler = (*Handler)(nil)
}

func TestHandleMessage_SuccessfulParse_SendsPreview(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        1234,
		Description:   "lunch with team",
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
	}}
	creator := &mockCashFlowCreator{}
	categories := []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}
	loader := &mockCategoryLoader{categories: categories}
	h := NewHandler(parser, creator, loader, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "lunch 1234",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.True(t, loader.called)
	require.True(t, parser.parseCalled)
	require.Equal(t, "lunch 1234", parser.message)
	require.Equal(t, categories, parser.categoriesPassed)
	require.Len(t, session.reactionAdds, 1)
	require.Equal(t, "⏳", session.reactionAdds[0].emojiID)
	require.Len(t, session.reactionRemoves, 1)
	require.Equal(t, "⏳", session.reactionRemoves[0].emojiID)
	require.Equal(t, "@me", session.reactionRemoves[0].userID)
	require.Len(t, session.sentMessages, 1)

	sent := session.sentMessages[0]
	require.Len(t, sent.Embeds, 1)
	embed := sent.Embeds[0]
	require.Equal(t, GetMessage(string(LangEn), MsgPreviewTitle), embed.Title)
	require.Equal(t, 0xFF0000, embed.Color)
	require.Len(t, embed.Fields, 5)
	require.Equal(t, GetMessage(string(LangEn), MsgFieldType), embed.Fields[0].Name)
	require.Equal(t, GetMessage(string(LangEn), MsgTypeExpense), embed.Fields[0].Value)
	require.Equal(t, GetMessage(string(LangEn), MsgFieldAmount), embed.Fields[1].Name)
	require.Contains(t, embed.Fields[1].Value, "1,234")
	require.Equal(t, "Food", embed.Fields[2].Value)
	require.Equal(t, "lunch with team", embed.Fields[3].Value)
	require.Equal(t, "2026-04-05", embed.Fields[4].Value)
	require.Len(t, sent.Components, 1)

	row, ok := sent.Components[0].(discordgo.ActionsRow)
	require.True(t, ok)
	require.Len(t, row.Components, 2)
	confirm, ok := row.Components[0].(*discordgo.Button)
	require.True(t, ok)
	require.Equal(t, GetMessage(string(LangEn), MsgConfirmButton), confirm.Label)
	require.Contains(t, confirm.CustomID, "confirm:")
	require.Contains(t, confirm.CustomID, ":author-1")
	cancel, ok := row.Components[1].(*discordgo.Button)
	require.True(t, ok)
	require.Equal(t, GetMessage(string(LangEn), MsgCancelButton), cancel.Label)
	require.Equal(t, "cancel:author-1", cancel.CustomID)
}

func TestHandleMessage_NonBookkeeping_Silent(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: false}}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}}
	h := NewHandler(parser, &mockCashFlowCreator{}, loader, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "hello there",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 0)
	require.Len(t, session.reactionAdds, 1)
	require.Len(t, session.reactionRemoves, 1)
}

func TestHandleMessage_MissingAmount_SendsHint(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		MissingFields: []string{"amount"},
	}}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}}
	h := NewHandler(parser, &mockCashFlowCreator{}, loader, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "lunch",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	require.Contains(t, session.sentMessages[0].Content, GetMessage(string(LangEn), MsgMissingAmount))
	require.Contains(t, session.sentMessages[0].Content, GetMessage(string(LangEn), MsgUsageExamples))
}

func TestHandleMessage_ParseError_SendsSystemError(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{err: errors.New("parser down")}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}}
	h := NewHandler(parser, &mockCashFlowCreator{}, loader, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "lunch 150",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgSystemError), session.sentMessages[0].Content)
}

func TestHandleInteraction_Confirm_CreatesRecord(t *testing.T) {
	session := &mockSession{}
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        150,
		Description:   "",
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
	}
	creator := &mockCashFlowCreator{resultID: "cashflow-1"}
	h := NewHandler(&mockParser{}, creator, &mockCategoryLoader{}, string(LangEn))
	customID := h.storePending(result, "author-1")
	interaction := newComponentInteraction(customID, "author-1")

	h.handleInteraction(session, interaction)

	require.Len(t, creator.createdInputs, 1)
	require.Equal(t, &BotCashFlowInput{
		Date:        "2026-04-05",
		Type:        "expense",
		CategoryID:  "expense-food",
		Amount:      150,
		Description: "",
	}, creator.createdInputs[0])
	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseUpdateMessage, resp.Type)
	require.Len(t, resp.Data.Embeds, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgConfirmSuccess), resp.Data.Embeds[0].Title)
	require.Empty(t, resp.Data.Components)
}

func TestHandleInteraction_Cancel_UpdatesEmbed(t *testing.T) {
	session := &mockSession{}
	creator := &mockCashFlowCreator{}
	h := NewHandler(&mockParser{}, creator, &mockCategoryLoader{}, string(LangEn))
	interaction := newComponentInteraction("cancel:author-1", "author-1")

	h.handleInteraction(session, interaction)

	require.Empty(t, creator.createdInputs)
	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseUpdateMessage, resp.Type)
	require.Len(t, resp.Data.Embeds, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgCancelled), resp.Data.Embeds[0].Title)
	require.Empty(t, resp.Data.Components)
}

func TestHandleInteraction_WrongUser_Ephemeral(t *testing.T) {
	session := &mockSession{}
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        150,
		Description:   "lunch",
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
	}
	creator := &mockCashFlowCreator{}
	h := NewHandler(&mockParser{}, creator, &mockCategoryLoader{}, string(LangEn))
	customID := h.storePending(result, "author-1")
	interaction := newComponentInteraction(customID, "other-user")

	h.handleInteraction(session, interaction)

	require.Empty(t, creator.createdInputs)
	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseChannelMessageWithSource, resp.Type)
	require.Equal(t, GetMessage(string(LangEn), MsgOnlyAuthor), resp.Data.Content)
	require.Equal(t, discordgo.MessageFlagsEphemeral, resp.Data.Flags)
}

func newComponentInteraction(customID string, userID string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:   discordgo.InteractionMessageComponent,
		Member: &discordgo.Member{User: &discordgo.User{ID: userID}},
		Message: &discordgo.Message{
			ID:        "reply-1",
			ChannelID: "channel-1",
			Embeds:    []*discordgo.MessageEmbed{{Title: GetMessage(string(LangEn), MsgPreviewTitle)}},
		},
		Data: discordgo.MessageComponentInteractionData{CustomID: customID},
	}}
}
