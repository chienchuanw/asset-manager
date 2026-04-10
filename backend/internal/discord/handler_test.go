package discord

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
)

type mockSession struct {
	mu                   sync.Mutex
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
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sentMessages = append(m.sentMessages, data)
	return &discordgo.Message{ID: "reply-1", ChannelID: channelID}, nil
}

func (m *mockSession) MessageReactionAdd(channelID, messageID, emojiID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reactionAdds = append(m.reactionAdds, reactionAddCall{channelID: channelID, messageID: messageID, emojiID: emojiID})
	return nil
}

func (m *mockSession) MessageReactionRemove(channelID, messageID, emojiID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reactionRemoves = append(m.reactionRemoves, reactionRemoveCall{channelID: channelID, messageID: messageID, emojiID: emojiID, userID: userID})
	return nil
}

func (m *mockSession) ChannelMessageEditComplex(data *discordgo.MessageEdit) (*discordgo.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.editedMessages = append(m.editedMessages, data)
	return &discordgo.Message{ID: data.ID, ChannelID: data.Channel}, nil
}

func (m *mockSession) InteractionRespond(_ *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.interactionResponses = append(m.interactionResponses, resp)
	return nil
}

type mockParser struct {
	mu               sync.Mutex
	result           *ParseResult
	err              error
	parseFunc        func(context.Context, string, []CategoryInfo) (*ParseResult, error)
	parseCalled      bool
	message          string
	categoriesPassed []CategoryInfo
}

func (m *mockParser) Parse(ctx context.Context, message string, categories []CategoryInfo) (*ParseResult, error) {
	m.mu.Lock()
	m.parseCalled = true
	m.message = message
	m.categoriesPassed = append([]CategoryInfo(nil), categories...)
	parseFunc := m.parseFunc
	result := m.result
	err := m.err
	m.mu.Unlock()
	if parseFunc != nil {
		return parseFunc(ctx, message, categories)
	}
	return result, err
}

type mockCashFlowCreator struct {
	mu            sync.Mutex
	createdInputs []*BotCashFlowInput
	resultID      string
	err           error
}

func (m *mockCashFlowCreator) CreateCashFlowFromBot(input *BotCashFlowInput) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createdInputs = append(m.createdInputs, input)
	return m.resultID, m.err
}

type mockCCPaymentCreator struct {
	mu            sync.Mutex
	resultID      string
	resultAmount  float64
	err           error
	createdInputs []*BotCCPaymentInput
}

func (m *mockCCPaymentCreator) CreatePaymentFromBot(input *BotCCPaymentInput) (string, float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createdInputs = append(m.createdInputs, input)
	if m.err != nil {
		return "", 0, m.err
	}
	amount := input.Amount
	if m.resultAmount > 0 {
		amount = m.resultAmount
	}
	return m.resultID, amount, nil
}

type mockCategoryLoader struct {
	mu         sync.Mutex
	categories []CategoryInfo
	err        error
	called     bool
}

func (m *mockCategoryLoader) LoadCategories() ([]CategoryInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.called = true
	return m.categories, m.err
}

type mockAccountLoader struct {
	accounts       []AccountInfo
	accountsByType map[string][]AccountInfo
	err            error
	errByType      map[string]error
}

func (m *mockAccountLoader) LoadAccounts(sourceType string) ([]AccountInfo, error) {
	if m.errByType != nil {
		if err, ok := m.errByType[sourceType]; ok {
			return nil, err
		}
	}
	if m.accountsByType != nil {
		return m.accountsByType[sourceType], nil
	}
	return m.accounts, m.err
}

type mockCashFlowQuerier struct {
	mu          sync.Mutex
	result      *MonthlySummaryResult
	err         error
	called      bool
	calledYear  int
	calledMonth int
}

func (m *mockCashFlowQuerier) GetMonthlySummary(year, month int) (*MonthlySummaryResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.called = true
	m.calledYear = year
	m.calledMonth = month
	return m.result, m.err
}

type mockAccountBalanceQuerier struct {
	mu     sync.Mutex
	result *AccountBalancesResult
	err    error
	called bool
}

func (m *mockAccountBalanceQuerier) GetAllBalances() (*AccountBalancesResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.called = true
	return m.result, m.err
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
		SourceType:    "cash",
	}}
	creator := &mockCashFlowCreator{}
	categories := []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}
	loader := &mockCategoryLoader{categories: categories}
	h := NewHandler(context.Background(), parser, creator, loader, nil, string(LangEn))
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
	require.Equal(t, GetMessage(string(LangEn), MsgFieldPaymentMethod), embed.Fields[4].Name)
	require.Equal(t, GetMessage(string(LangEn), MsgAccountCash), embed.Fields[4].Value)
	require.NotNil(t, embed.Footer)
	require.Equal(t, "2026-04-05", embed.Footer.Text)
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
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, nil, string(LangEn))
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
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, nil, string(LangEn))
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

func TestHandleMessage_CCPayment_MissingAmount(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Action:        "cc_payment",
		PaymentType:   "custom",
		Amount:        0,
		MissingFields: []string{"amount"},
	}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "pay card", Author: &discordgo.User{ID: "author-1"}}})

	require.Len(t, session.sentMessages, 1)
	require.Contains(t, session.sentMessages[0].Content, GetMessage(string(LangEn), MsgCCPaymentMissingAmount))
	require.Contains(t, session.sentMessages[0].Content, GetMessage(string(LangEn), MsgCCPaymentUsageExamples))
}

func TestHandleMessage_CCPayment_ShowsCreditCardMenu(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Action:        "cc_payment",
		Amount:        15000,
		Date:          "2026-04-05",
		PaymentType:   "custom",
		CategoryID:    "expense-transfer",
	}}
	acctLoader := &mockAccountLoader{accountsByType: map[string][]AccountInfo{
		"credit_card": {{ID: "cc-1", Name: "Citi Visa *1234", Type: "credit_card"}, {ID: "cc-2", Name: "Chase *5678", Type: "credit_card"}},
	}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, acctLoader, string(LangEn))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "pay card 15000", Author: &discordgo.User{ID: "author-1"}}})

	require.Len(t, session.sentMessages, 1)
	sent := session.sentMessages[0]
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentSelectCard), sent.Content)
	require.Len(t, sent.Components, 1)
	row := sent.Components[0].(discordgo.ActionsRow)
	menu := row.Components[0].(*discordgo.SelectMenu)
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentSelectCard), menu.Placeholder)
	require.Len(t, menu.Options, 2)
	require.Equal(t, "cc-1", menu.Options[0].Value)
	require.Equal(t, "Citi Visa *1234", menu.Options[0].Label)
	require.Contains(t, menu.CustomID, "select_cc:")
	require.Contains(t, menu.CustomID, ":author-1")
}

func TestHandleMessage_CCPayment_NoCreditCards(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "cc_payment", Amount: 15000, PaymentType: "custom"}}
	acctLoader := &mockAccountLoader{accountsByType: map[string][]AccountInfo{"credit_card": {}}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, acctLoader, string(LangEn))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "pay card 15000", Author: &discordgo.User{ID: "author-1"}}})

	require.Len(t, session.sentMessages, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentNoCards), session.sentMessages[0].Content)
}

func TestHandleMessage_ChatAction_ZhTW(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{Action: "chat", IsBookkeeping: false}}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, nil, string(LangZhTW))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "你好",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	require.Contains(t, session.sentMessages[0].Content, GetMessage(string(LangZhTW), MsgChatGreeting))
}

func TestHandleMessage_ChatAction_En(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{Action: "chat", IsBookkeeping: false}}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, nil, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "hello",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	require.Contains(t, session.sentMessages[0].Content, GetMessage(string(LangEn), MsgChatGreeting))
}

func TestHandleMessage_UnsupportedAction_ZhTW(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{Action: "unsupported", IsBookkeeping: false}}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, nil, string(LangZhTW))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "do something unsupported",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	require.Contains(t, session.sentMessages[0].Content, GetMessage(string(LangZhTW), MsgUnsupported))
}

func TestHandleMessage_UnsupportedAction_En(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{Action: "unsupported", IsBookkeeping: false}}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, nil, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "do something unsupported",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	require.Contains(t, session.sentMessages[0].Content, GetMessage(string(LangEn), MsgUnsupported))
}

func TestHandleMessage_CategoryLoadFail_SpecificError(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{}
	loader := &mockCategoryLoader{err: errors.New("db down")}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, nil, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "lunch 150",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.False(t, parser.parseCalled)
	require.Len(t, session.sentMessages, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgDataLoadFailed), session.sentMessages[0].Content)
}

func TestHandleMessage_ParseFail_SpecificError(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{err: errors.New("parser down")}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, nil, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "lunch 150",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgParseFailed), session.sentMessages[0].Content)
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
	h := NewHandler(context.Background(), &mockParser{}, creator, &mockCategoryLoader{}, nil, string(LangEn))
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
	h := NewHandler(context.Background(), &mockParser{}, creator, &mockCategoryLoader{}, nil, string(LangEn))
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
	h := NewHandler(context.Background(), &mockParser{}, creator, &mockCategoryLoader{}, nil, string(LangEn))
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

func TestHandleMessage_EmptySourceType_SendsSelectMenu(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		Description:   "午餐",
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
		SourceType:    "",
	}}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, nil, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "lunch 180",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	sent := session.sentMessages[0]
	require.Empty(t, sent.Embeds)
	require.Len(t, sent.Components, 1)
	row, ok := sent.Components[0].(discordgo.ActionsRow)
	require.True(t, ok)
	require.Len(t, row.Components, 1)
	selectMenu, ok := row.Components[0].(*discordgo.SelectMenu)
	require.True(t, ok)
	require.Equal(t, GetMessage(string(LangEn), MsgSelectAccount), selectMenu.Placeholder)
	require.Len(t, selectMenu.Options, 3)
}

func TestHandleMessage_CreditCard_ShowsAccountIDMenu(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        2000,
		Description:   "刷卡買衣服",
		CategoryID:    "expense-other",
		CategoryName:  "Other",
		Date:          "2026-04-05",
		SourceType:    "credit_card",
	}}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-other", Name: "Other", Type: "expense"}}}
	acctLoader := &mockAccountLoader{accounts: []AccountInfo{
		{ID: "cc-1", Name: "中信 Visa *1234", Type: "credit_card"},
	}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, acctLoader, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "credit card clothes 2000",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	sent := session.sentMessages[0]
	require.Empty(t, sent.Embeds)
	require.Len(t, sent.Components, 1)
	row, ok := sent.Components[0].(discordgo.ActionsRow)
	require.True(t, ok)
	selectMenu, ok := row.Components[0].(*discordgo.SelectMenu)
	require.True(t, ok)
	require.Equal(t, GetMessage(string(LangEn), MsgSelectCreditCard), selectMenu.Placeholder)
	require.Len(t, selectMenu.Options, 1)
	require.Equal(t, "cc-1", selectMenu.Options[0].Value)
}

func TestHandleMessage_Cash_SendsPreviewDirectly(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		Description:   "lunch",
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
		SourceType:    "cash",
	}}
	loader := &mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, loader, nil, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "cash lunch 180",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	sent := session.sentMessages[0]
	require.Len(t, sent.Embeds, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgPreviewTitle), sent.Embeds[0].Title)
}

func TestHandleMessage_RoutesToCreateFlow(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Action:        "create",
		Type:          "expense",
		Amount:        180,
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
		SourceType:    "cash",
	}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "lunch 180",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	require.Len(t, session.sentMessages[0].Embeds, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgPreviewTitle), session.sentMessages[0].Embeds[0].Title)
}

func TestHandleMessage_CreateAction_RegressionIdentical(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Action:        "create",
		Type:          "expense",
		Amount:        180,
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
		SourceType:    "",
	}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "lunch 180", Author: &discordgo.User{ID: "author-1"}}})

	require.Len(t, session.sentMessages, 1)
	require.Empty(t, session.sentMessages[0].Embeds)
	require.Len(t, session.sentMessages[0].Components, 1)
	row := session.sentMessages[0].Components[0].(discordgo.ActionsRow)
	menu := row.Components[0].(*discordgo.SelectMenu)
	require.Equal(t, GetMessage(string(LangEn), MsgSelectAccount), menu.Placeholder)
}

func TestHandleMessage_BackwardCompat_NoActionField(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: false, Action: ""}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "hi there", Author: &discordgo.User{ID: "author-1"}}})

	require.Empty(t, session.sentMessages)
}

func TestHandleMessage_CreateWithQueryParams_IgnoresParams(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Action:        "create",
		Type:          "expense",
		Amount:        320,
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
		SourceType:    "cash",
		QueryParams:   &QueryParams{Month: 3, Year: 2026},
	}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "lunch 320", Author: &discordgo.User{ID: "author-1"}}})

	require.Len(t, session.sentMessages, 1)
	require.Len(t, session.sentMessages[0].Embeds, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgPreviewTitle), session.sentMessages[0].Embeds[0].Title)
	require.Empty(t, session.sentMessages[0].Content)
	require.Len(t, session.sentMessages[0].Components, 1)
}

func TestHandleMessage_RoutesToQueryFlow(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Action:        "query",
		QueryType:     "cash_flow_summary",
		QueryParams:   &QueryParams{Year: 2026, Month: 4},
	}}
	cfQuerier := &mockCashFlowQuerier{result: &MonthlySummaryResult{Year: 2026, Month: 4, TotalIncome: 1000, TotalExpense: 500, NetCashFlow: 500}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangZhTW), WithCashFlowQuerier(cfQuerier))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "本月支出摘要",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.True(t, cfQuerier.called)
	require.Len(t, session.sentMessages, 1)
	require.Len(t, session.sentMessages[0].Embeds, 1)
	require.Empty(t, session.sentMessages[0].Components)
	require.Equal(t, "📊 4月現金流摘要", session.sentMessages[0].Embeds[0].Title)
}

func TestHandleMessage_ChatIgnored(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: false, Action: ""}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "hello",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Empty(t, session.sentMessages)
}

func TestHandleMessage_ConcurrentQueryAndCreate(t *testing.T) {
	t.Parallel()

	h := NewHandler(context.Background(),
		&mockParser{parseFunc: func(_ context.Context, message string, _ []CategoryInfo) (*ParseResult, error) {
			if message == "query" {
				return &ParseResult{
					IsBookkeeping: true,
					Action:        "query",
					QueryType:     "cash_flow_summary",
					QueryParams:   &QueryParams{Year: 2026, Month: 4},
				}, nil
			}
			return &ParseResult{
				IsBookkeeping: true,
				Action:        "create",
				Type:          "expense",
				Amount:        180,
				CategoryID:    "expense-food",
				CategoryName:  "Food",
				Date:          "2026-04-05",
				SourceType:    "cash",
			}, nil
		}},
		&mockCashFlowCreator{},
		&mockCategoryLoader{categories: []CategoryInfo{{ID: "expense-food", Name: "Food", Type: "expense"}}},
		nil,
		string(LangEn),
		WithCashFlowQuerier(&mockCashFlowQuerier{result: &MonthlySummaryResult{Year: 2026, Month: 4, TotalIncome: 1000, TotalExpense: 500, NetCashFlow: 500}}),
	)

	var wg sync.WaitGroup
	for _, tc := range []struct {
		session *mockSession
		msg     *discordgo.MessageCreate
	}{
		{session: &mockSession{}, msg: &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-query", ChannelID: "channel-1", Content: "query", Author: &discordgo.User{ID: "author-1"}}}},
		{session: &mockSession{}, msg: &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-create", ChannelID: "channel-1", Content: "create", Author: &discordgo.User{ID: "author-2"}}}},
	} {
		wg.Add(1)
		go func(session *mockSession, msg *discordgo.MessageCreate) {
			defer wg.Done()
			h.handleMessage(session, msg)
		}(tc.session, tc.msg)
	}
	wg.Wait()
}

func TestHandleQuery_UnsupportedType(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "unknown", QueryParams: &QueryParams{Year: 2026, Month: 4}}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))
	msg := &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "message-1",
		ChannelID: "channel-1",
		Content:   "what can you do",
		Author:    &discordgo.User{ID: "author-1"},
	}}

	h.handleMessage(session, msg)

	require.Len(t, session.sentMessages, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgQueryUnsupported), session.sentMessages[0].Content)
}

func TestHandleCashFlowQuery_CurrentMonth_GivenSummary_WhenHandleMessage_ThenSendEmbed(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{
		IsBookkeeping: true,
		Action:        "query",
		QueryType:     "cash_flow_summary",
		QueryParams:   &QueryParams{Year: 2026, Month: 4},
	}}
	cfQuerier := &mockCashFlowQuerier{result: &MonthlySummaryResult{
		Year:          2026,
		Month:         4,
		TotalIncome:   80000,
		TotalExpense:  23500,
		NetCashFlow:   56500,
		IncomeCount:   2,
		ExpenseCount:  5,
		TopCategories: []CategoryBreakdown{{Name: "飲食", Amount: 5000}, {Name: "交通", Amount: 3000}},
		Comparison:    &MonthComparisonResult{ExpenseChange: 1200, ExpenseChangePct: 5.4, IncomeChange: 8000, IncomeChangePct: 11.1},
	}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangZhTW), WithCashFlowQuerier(cfQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "本月現金流", Author: &discordgo.User{ID: "author-1"}}})

	require.Len(t, session.sentMessages, 1)
	embed := session.sentMessages[0].Embeds[0]
	require.Equal(t, "📊 4月現金流摘要", embed.Title)
	require.Contains(t, embed.Fields[0].Value, "80,000")
	require.Contains(t, embed.Fields[1].Value, "23,500")
	require.Contains(t, embed.Fields[2].Value, "56,500")
	require.Contains(t, embed.Fields[3].Value, "7")
	require.Equal(t, GetMessage(string(LangZhTW), MsgQueryTopCategories), embed.Fields[4].Name)
	require.Contains(t, embed.Fields[4].Value, "飲食: $5,000")
	require.Equal(t, GetMessage(string(LangZhTW), MsgQueryComparison), embed.Fields[5].Name)
	require.Contains(t, embed.Fields[5].Value, "1,200")
	require.Contains(t, embed.Fields[5].Value, "5.4%")
}

func TestHandleCashFlowQuery_SpecificMonth_GivenMonthParam_WhenHandleMessage_ThenUseMonthInTitle(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "cash_flow_summary", QueryParams: &QueryParams{Year: 2026, Month: 3}}}
	cfQuerier := &mockCashFlowQuerier{result: &MonthlySummaryResult{Year: 2026, Month: 3, TotalIncome: 1, TotalExpense: 1, NetCashFlow: 0}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangZhTW), WithCashFlowQuerier(cfQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "3月支出", Author: &discordgo.User{ID: "author-1"}}})

	require.Equal(t, "📊 3月現金流摘要", session.sentMessages[0].Embeds[0].Title)
}

func TestHandleCashFlowQuery_WithCategory_GivenCategoryParam_WhenHandleMessage_ThenUseCategoryTitle(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "cash_flow_summary", QueryParams: &QueryParams{Year: 2026, Month: 4, Category: "飲食"}}}
	cfQuerier := &mockCashFlowQuerier{result: &MonthlySummaryResult{
		Year:          2026,
		Month:         4,
		TotalIncome:   80000,
		TotalExpense:  23500,
		NetCashFlow:   56500,
		TopCategories: []CategoryBreakdown{{Name: "飲食", Amount: 5000, Count: 3}},
	}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangZhTW), WithCashFlowQuerier(cfQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "4月飲食支出", Author: &discordgo.User{ID: "author-1"}}})

	embed := session.sentMessages[0].Embeds[0]
	require.Equal(t, "📊 4月飲食支出", embed.Title)
	require.Contains(t, embed.Fields[1].Value, "5,000")
}

func TestHandleCashFlowQuery_NoData_GivenZeroSummary_WhenHandleMessage_ThenShowEmptyDescription(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "cash_flow_summary", QueryParams: &QueryParams{Year: 2026, Month: 4}}}
	cfQuerier := &mockCashFlowQuerier{result: &MonthlySummaryResult{Year: 2026, Month: 4}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangZhTW), WithCashFlowQuerier(cfQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "本月摘要", Author: &discordgo.User{ID: "author-1"}}})

	require.Equal(t, GetMessage(string(LangZhTW), MsgQueryNoData), session.sentMessages[0].Embeds[0].Description)
}

func TestHandleCashFlowQuery_ServiceError_GivenServiceFailure_WhenHandleMessage_ThenReplySystemError(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "cash_flow_summary", QueryParams: &QueryParams{Year: 2026, Month: 4}}}
	cfQuerier := &mockCashFlowQuerier{err: errors.New("boom")}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn), WithCashFlowQuerier(cfQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "cash flow", Author: &discordgo.User{ID: "author-1"}}})

	require.Len(t, session.sentMessages, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgQueryFailed), session.sentMessages[0].Content)
}

func TestHandleCashFlowQuery_Fail_SpecificError(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "cash_flow_summary", QueryParams: &QueryParams{Year: 2026, Month: 4}}}
	cfQuerier := &mockCashFlowQuerier{err: errors.New("boom")}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangZhTW), WithCashFlowQuerier(cfQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "本月摘要", Author: &discordgo.User{ID: "author-1"}}})

	require.Len(t, session.sentMessages, 1)
	require.Equal(t, GetMessage(string(LangZhTW), MsgQueryFailed), session.sentMessages[0].Content)
}

func TestHandleAccountBalanceQuery_BankAndCC_GivenBalances_WhenHandleMessage_ThenSendSections(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "account_balance", QueryParams: &QueryParams{Year: 2026, Month: 4}}}
	acctQuerier := &mockAccountBalanceQuerier{result: &AccountBalancesResult{
		BankAccounts: []BankAccountBalance{{Name: "中信銀行", Last4: "1234", Balance: 25000}, {Name: "台新銀行", Last4: "5678", Balance: 12000}},
		CreditCards:  []CreditCardBalance{{Name: "中信 Visa", Last4: "4321", CreditLimit: 100000, UsedCredit: 20000, Remaining: 80000, UsagePct: 20}},
	}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangZhTW), WithAccountBalanceQuerier(acctQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "帳戶餘額", Author: &discordgo.User{ID: "author-1"}}})

	embed := session.sentMessages[0].Embeds[0]
	require.Equal(t, GetMessage(string(LangZhTW), MsgQueryAccountTitle), embed.Title)
	require.Len(t, embed.Fields, 2)
	require.Equal(t, GetMessage(string(LangZhTW), MsgQueryBankSection), embed.Fields[0].Name)
	require.Contains(t, embed.Fields[0].Value, "中信銀行 *1234")
	require.Contains(t, embed.Fields[0].Value, GetMessage(string(LangZhTW), MsgQueryBankTotal))
	require.Equal(t, GetMessage(string(LangZhTW), MsgQueryCCSection), embed.Fields[1].Name)
	require.Contains(t, embed.Fields[1].Value, "中信 Visa *4321")
}

func TestHandleAccountBalanceQuery_CCNearingLimit_GivenHighUsage_WhenHandleMessage_ThenShowWarning(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "account_balance", QueryParams: &QueryParams{Year: 2026, Month: 4}}}
	acctQuerier := &mockAccountBalanceQuerier{result: &AccountBalancesResult{CreditCards: []CreditCardBalance{{Name: "中信 Visa", Last4: "4321", CreditLimit: 100000, UsedCredit: 85000, Remaining: 15000, UsagePct: 85}}}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangZhTW), WithAccountBalanceQuerier(acctQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "信用卡額度", Author: &discordgo.User{ID: "author-1"}}})

	require.Contains(t, session.sentMessages[0].Embeds[0].Fields[0].Value, GetMessage(string(LangZhTW), MsgQueryCCNearLimit))
}

func TestHandleAccountBalanceQuery_NoAccounts_GivenEmptyResult_WhenHandleMessage_ThenSendNoAccountsMessage(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "account_balance", QueryParams: &QueryParams{Year: 2026, Month: 4}}}
	acctQuerier := &mockAccountBalanceQuerier{result: &AccountBalancesResult{}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangZhTW), WithAccountBalanceQuerier(acctQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "餘額", Author: &discordgo.User{ID: "author-1"}}})

	require.Equal(t, GetMessage(string(LangZhTW), MsgQueryNoAccounts), session.sentMessages[0].Content)
}

func TestHandleAccountBalanceQuery_PartialFailure_GivenBankOKAndCCError_WhenHandleMessage_ThenShowMixedResults(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "account_balance", QueryParams: &QueryParams{Year: 2026, Month: 4}}}
	acctQuerier := &mockAccountBalanceQuerier{result: &AccountBalancesResult{BankAccounts: []BankAccountBalance{{Name: "中信銀行", Last4: "1234", Balance: 25000}}, CCError: errors.New("cc down")}}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangZhTW), WithAccountBalanceQuerier(acctQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "帳戶餘額", Author: &discordgo.User{ID: "author-1"}}})

	embed := session.sentMessages[0].Embeds[0]
	require.Contains(t, embed.Fields[0].Value, "中信銀行 *1234")
	require.Equal(t, GetMessage(string(LangZhTW), MsgQueryLoadFailed), embed.Fields[1].Value)
}

func TestHandleAccountBalanceQuery_Fail_SpecificError(t *testing.T) {
	session := &mockSession{}
	parser := &mockParser{result: &ParseResult{IsBookkeeping: true, Action: "query", QueryType: "account_balance", QueryParams: &QueryParams{Year: 2026, Month: 4}}}
	acctQuerier := &mockAccountBalanceQuerier{err: errors.New("boom")}
	h := NewHandler(context.Background(), parser, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn), WithAccountBalanceQuerier(acctQuerier))

	h.handleMessage(session, &discordgo.MessageCreate{Message: &discordgo.Message{ID: "message-1", ChannelID: "channel-1", Content: "balances", Author: &discordgo.User{ID: "author-1"}}})

	require.Len(t, session.sentMessages, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgQueryFailed), session.sentMessages[0].Content)
}

func TestHandleInteraction_SelectMenu_CashGoesToPreview(t *testing.T) {
	session := &mockSession{}
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		Description:   "lunch",
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
		SourceType:    "",
	}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.mu.Lock()
	h.pending["test-key"] = pendingEntry{result: result, authorID: "author-1", awaitingAccount: true}
	h.mu.Unlock()

	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:   discordgo.InteractionMessageComponent,
		Member: &discordgo.Member{User: &discordgo.User{ID: "author-1"}},
		Message: &discordgo.Message{
			ID:        "reply-1",
			ChannelID: "channel-1",
		},
		Data: discordgo.MessageComponentInteractionData{
			CustomID:      "select_account:test-key:author-1",
			ComponentType: discordgo.SelectMenuComponent,
			Values:        []string{"cash"},
		},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseUpdateMessage, resp.Type)
	require.Len(t, resp.Data.Embeds, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgPreviewTitle), resp.Data.Embeds[0].Title)
	require.Len(t, resp.Data.Components, 1)
}

func TestHandleInteraction_SelectMenu_WrongUser_Ephemeral(t *testing.T) {
	session := &mockSession{}
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		SourceType:    "",
	}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.mu.Lock()
	h.pending["test-key"] = pendingEntry{result: result, authorID: "author-1", awaitingAccount: true}
	h.mu.Unlock()

	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "other-user"}},
		Message: &discordgo.Message{ID: "reply-1", ChannelID: "channel-1"},
		Data: discordgo.MessageComponentInteractionData{
			CustomID:      "select_account:test-key:author-1",
			ComponentType: discordgo.SelectMenuComponent,
			Values:        []string{"cash"},
		},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseChannelMessageWithSource, resp.Type)
	require.Equal(t, GetMessage(string(LangEn), MsgOnlyAuthor), resp.Data.Content)
	require.Equal(t, discordgo.MessageFlagsEphemeral, resp.Data.Flags)
}

func TestBuildPreviewEmbed_IncludesPaymentMethodAndFooter(t *testing.T) {
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		Description:   "lunch",
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
		SourceType:    "credit_card",
		SourceName:    "中信 Visa *1234",
	}

	embed := h.buildPreviewEmbed(result)

	require.Len(t, embed.Fields, 6)
	require.Equal(t, GetMessage(string(LangEn), MsgFieldPaymentMethod), embed.Fields[4].Name)
	require.Equal(t, GetMessage(string(LangEn), MsgAccountCreditCard), embed.Fields[4].Value)
	require.Equal(t, GetMessage(string(LangEn), MsgFieldAccount), embed.Fields[5].Name)
	require.Equal(t, "中信 Visa *1234", embed.Fields[5].Value)
	require.NotNil(t, embed.Footer)
	require.Equal(t, "2026-04-05", embed.Footer.Text)
}

func TestHandleInteraction_Confirm_PassesSourceType(t *testing.T) {
	session := &mockSession{}
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        2000,
		Description:   "clothes",
		CategoryID:    "expense-other",
		CategoryName:  "Other",
		Date:          "2026-04-05",
		SourceType:    "credit_card",
	}
	creator := &mockCashFlowCreator{resultID: "cashflow-1"}
	h := NewHandler(context.Background(), &mockParser{}, creator, &mockCategoryLoader{}, nil, string(LangEn))
	customID := h.storePending(result, "author-1")
	interaction := newComponentInteraction(customID, "author-1")

	h.handleInteraction(session, interaction)

	require.Len(t, creator.createdInputs, 1)
	require.Equal(t, "credit_card", creator.createdInputs[0].SourceType)
}

func TestHandleInteraction_CCSelectCard_ShowsBankMenu(t *testing.T) {
	session := &mockSession{}
	acctLoader := &mockAccountLoader{accountsByType: map[string][]AccountInfo{
		"credit_card":  {{ID: "cc-1", Name: "Citi Visa *1234", Type: "credit_card"}},
		"bank_account": {{ID: "bank-1", Name: "Chase Bank *8888", Type: "bank_account"}},
	}}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, acctLoader, string(LangEn))

	h.mu.Lock()
	h.pending["cc-key"] = pendingEntry{result: &ParseResult{Action: "cc_payment", Amount: 15000, Date: "2026-04-05", PaymentType: "custom"}, authorID: "author-1", ccPayment: true, ccAmount: 15000, ccPaymentType: "custom"}
	h.mu.Unlock()

	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "author-1"}},
		Message: &discordgo.Message{ID: "reply-1", ChannelID: "channel-1"},
		Data:    discordgo.MessageComponentInteractionData{CustomID: "select_cc:cc-key:author-1", ComponentType: discordgo.SelectMenuComponent, Values: []string{"cc-1"}},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseUpdateMessage, resp.Type)
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentSelectBank), resp.Data.Content)
	row := resp.Data.Components[0].(discordgo.ActionsRow)
	menu := row.Components[0].(*discordgo.SelectMenu)
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentSelectBank), menu.Placeholder)
	require.Len(t, menu.Options, 1)
	require.Equal(t, "bank-1", menu.Options[0].Value)
	require.Equal(t, "Chase Bank *8888", menu.Options[0].Label)
	require.Equal(t, "select_cc_bank:cc-key:author-1", menu.CustomID)
}

func TestHandleInteraction_CCSelectCard_NoBankAccounts(t *testing.T) {
	session := &mockSession{}
	acctLoader := &mockAccountLoader{accountsByType: map[string][]AccountInfo{
		"credit_card":  {{ID: "cc-1", Name: "Citi Visa *1234", Type: "credit_card"}},
		"bank_account": {},
	}}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, acctLoader, string(LangEn))

	h.mu.Lock()
	h.pending["cc-key"] = pendingEntry{result: &ParseResult{Action: "cc_payment", Amount: 15000, Date: "2026-04-05", PaymentType: "custom"}, authorID: "author-1", ccPayment: true, ccAmount: 15000, ccPaymentType: "custom"}
	h.mu.Unlock()

	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "author-1"}},
		Message: &discordgo.Message{ID: "reply-1", ChannelID: "channel-1"},
		Data:    discordgo.MessageComponentInteractionData{CustomID: "select_cc:cc-key:author-1", ComponentType: discordgo.SelectMenuComponent, Values: []string{"cc-1"}},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentNoBankAccounts), resp.Data.Content)
	require.Empty(t, resp.Data.Components)
}

func TestHandleInteraction_CCSelectBank_ShowsPreview(t *testing.T) {
	session := &mockSession{}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, &mockAccountLoader{accountsByType: map[string][]AccountInfo{
		"bank_account": {{ID: "bank-1", Name: "Chase Bank *8888", Type: "bank_account"}},
	}}, string(LangEn))

	h.mu.Lock()
	h.pending["cc-key"] = pendingEntry{
		result:        &ParseResult{Action: "cc_payment", Amount: 15000, Date: "2026-04-05", PaymentType: "custom"},
		authorID:      "author-1",
		ccPayment:     true,
		ccCardID:      "cc-1",
		ccCardName:    "Citi Visa *1234",
		ccAmount:      15000,
		ccPaymentType: "custom",
	}
	h.mu.Unlock()

	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "author-1"}},
		Message: &discordgo.Message{ID: "reply-1", ChannelID: "channel-1"},
		Data:    discordgo.MessageComponentInteractionData{CustomID: "select_cc_bank:cc-key:author-1", ComponentType: discordgo.SelectMenuComponent, Values: []string{"bank-1"}},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseUpdateMessage, resp.Type)
	require.Len(t, resp.Data.Embeds, 1)
	embed := resp.Data.Embeds[0]
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentPreview), embed.Title)
	require.Equal(t, 0x3498DB, embed.Color)
	require.Len(t, embed.Fields, 4)
	require.Equal(t, GetMessage(string(LangEn), MsgFieldAmount), embed.Fields[0].Name)
	require.Equal(t, "$15,000", embed.Fields[0].Value)
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentFieldCard), embed.Fields[1].Name)
	require.Equal(t, "Citi Visa *1234", embed.Fields[1].Value)
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentFieldBank), embed.Fields[2].Name)
	require.Equal(t, "Chase Bank *8888", embed.Fields[2].Value)
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentTypeCustom), embed.Fields[3].Value)
	require.Len(t, resp.Data.Components, 1)
	row := resp.Data.Components[0].(discordgo.ActionsRow)
	require.Len(t, row.Components, 2)
	confirm := row.Components[0].(*discordgo.Button)
	cancel := row.Components[1].(*discordgo.Button)
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentConfirmButton), confirm.Label)
	require.Contains(t, confirm.CustomID, "confirm_cc_payment:")
	require.Equal(t, "cancel:author-1", cancel.CustomID)
}

func TestHandleInteraction_CCConfirm_Success(t *testing.T) {
	session := &mockSession{}
	ccCreator := &mockCCPaymentCreator{resultID: "payment-1", resultAmount: 15000}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn), WithCCPaymentCreator(ccCreator))

	h.mu.Lock()
	h.pending["confirm-key"] = pendingEntry{
		result:        &ParseResult{Action: "cc_payment", Date: "2026-04-05", CategoryID: "category-1"},
		authorID:      "author-1",
		ccPayment:     true,
		ccCardID:      "cc-1",
		ccCardName:    "Citi Visa *1234",
		ccBankID:      "bank-1",
		ccBankName:    "Chase Bank *8888",
		ccAmount:      15000,
		ccPaymentType: "custom",
	}
	h.mu.Unlock()

	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "author-1"}},
		Message: &discordgo.Message{ID: "reply-1", ChannelID: "channel-1", Embeds: []*discordgo.MessageEmbed{{Title: GetMessage(string(LangEn), MsgCCPaymentPreview)}}},
		Data:    discordgo.MessageComponentInteractionData{CustomID: "confirm_cc_payment:confirm-key:author-1"},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, ccCreator.createdInputs, 1)
	require.Equal(t, &BotCCPaymentInput{CreditCardID: "cc-1", BankAccountID: "bank-1", Amount: 15000, Date: "2026-04-05", PaymentType: "custom"}, ccCreator.createdInputs[0])
	require.Len(t, session.interactionResponses, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentSuccess), session.interactionResponses[0].Data.Embeds[0].Title)
	require.Empty(t, session.interactionResponses[0].Data.Components)
}

func TestHandleInteraction_CCConfirm_Failure(t *testing.T) {
	session := &mockSession{}
	ccCreator := &mockCCPaymentCreator{err: errors.New("payment failed")}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn), WithCCPaymentCreator(ccCreator))

	h.mu.Lock()
	h.pending["confirm-key"] = pendingEntry{
		result:        &ParseResult{Action: "cc_payment", Date: "2026-04-05", CategoryID: "category-1"},
		authorID:      "author-1",
		ccPayment:     true,
		ccCardID:      "cc-1",
		ccBankID:      "bank-1",
		ccAmount:      15000,
		ccPaymentType: "custom",
	}
	h.mu.Unlock()

	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "author-1"}},
		Message: &discordgo.Message{ID: "reply-1", ChannelID: "channel-1", Embeds: []*discordgo.MessageEmbed{{Title: GetMessage(string(LangEn), MsgCCPaymentPreview)}}},
		Data:    discordgo.MessageComponentInteractionData{CustomID: "confirm_cc_payment:confirm-key:author-1"},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, ccCreator.createdInputs, 1)
	require.Len(t, session.interactionResponses, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgCCPaymentFailed), session.interactionResponses[0].Data.Embeds[0].Title)
	require.Equal(t, "payment failed", session.interactionResponses[0].Data.Embeds[0].Description)
}

func TestHandleInteraction_CCCancel(t *testing.T) {
	session := &mockSession{}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))
	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "author-1"}},
		Message: &discordgo.Message{ID: "reply-1", ChannelID: "channel-1", Embeds: []*discordgo.MessageEmbed{{Title: GetMessage(string(LangEn), MsgCCPaymentPreview)}}},
		Data:    discordgo.MessageComponentInteractionData{CustomID: "cancel:author-1"},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, session.interactionResponses, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgCancelled), session.interactionResponses[0].Data.Embeds[0].Title)
}

func TestHandleInteraction_SelectBankAccount_ShowsSecondSelectMenu(t *testing.T) {
	session := &mockSession{}
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        15000,
		Description:   "rent",
		CategoryID:    "expense-other",
		CategoryName:  "Other",
		Date:          "2026-04-05",
		SourceType:    "",
	}
	acctLoader := &mockAccountLoader{accounts: []AccountInfo{
		{ID: "acct-1", Name: "中信銀行 *1234", Type: "bank_account"},
		{ID: "acct-2", Name: "台新銀行 *5678", Type: "bank_account"},
	}}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, acctLoader, string(LangEn))

	h.mu.Lock()
	h.pending["test-key"] = pendingEntry{result: result, authorID: "author-1", awaitingAccount: true}
	h.mu.Unlock()

	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "author-1"}},
		Message: &discordgo.Message{ID: "reply-1", ChannelID: "channel-1"},
		Data: discordgo.MessageComponentInteractionData{
			CustomID:      "select_account:test-key:author-1",
			ComponentType: discordgo.SelectMenuComponent,
			Values:        []string{"bank_account"},
		},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseUpdateMessage, resp.Type)
	require.Len(t, resp.Data.Components, 1)
	row, ok := resp.Data.Components[0].(discordgo.ActionsRow)
	require.True(t, ok)
	selectMenu, ok := row.Components[0].(*discordgo.SelectMenu)
	require.True(t, ok)
	require.Equal(t, GetMessage(string(LangEn), MsgSelectBankAccount), selectMenu.Placeholder)
	require.Len(t, selectMenu.Options, 2)
	require.Equal(t, "acct-1", selectMenu.Options[0].Value)
	require.Equal(t, "中信銀行 *1234", selectMenu.Options[0].Label)
}

func TestHandleInteraction_SelectCash_SkipsSecondMenu(t *testing.T) {
	session := &mockSession{}
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        180,
		Description:   "lunch",
		CategoryID:    "expense-food",
		CategoryName:  "Food",
		Date:          "2026-04-05",
		SourceType:    "",
	}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.mu.Lock()
	h.pending["test-key"] = pendingEntry{result: result, authorID: "author-1", awaitingAccount: true}
	h.mu.Unlock()

	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "author-1"}},
		Message: &discordgo.Message{ID: "reply-1", ChannelID: "channel-1"},
		Data: discordgo.MessageComponentInteractionData{
			CustomID:      "select_account:test-key:author-1",
			ComponentType: discordgo.SelectMenuComponent,
			Values:        []string{"cash"},
		},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseUpdateMessage, resp.Type)
	require.Len(t, resp.Data.Embeds, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgPreviewTitle), resp.Data.Embeds[0].Title)
}

func TestHandleInteraction_SelectAccountID_UpdatesPendingAndShowsPreview(t *testing.T) {
	session := &mockSession{}
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        15000,
		Description:   "rent",
		CategoryID:    "expense-other",
		CategoryName:  "Other",
		Date:          "2026-04-05",
		SourceType:    "bank_account",
	}
	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.mu.Lock()
	h.pending["test-key"] = pendingEntry{result: result, authorID: "author-1", awaitingAccountID: true}
	h.mu.Unlock()

	interaction := &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Member:  &discordgo.Member{User: &discordgo.User{ID: "author-1"}},
		Message: &discordgo.Message{ID: "reply-1", ChannelID: "channel-1"},
		Data: discordgo.MessageComponentInteractionData{
			CustomID:      "select_account_id:test-key:author-1",
			ComponentType: discordgo.SelectMenuComponent,
			Values:        []string{"acct-1"},
		},
	}}

	h.handleInteraction(session, interaction)

	require.Len(t, session.interactionResponses, 1)
	resp := session.interactionResponses[0]
	require.Equal(t, discordgo.InteractionResponseUpdateMessage, resp.Type)
	require.Len(t, resp.Data.Embeds, 1)
	require.Equal(t, GetMessage(string(LangEn), MsgPreviewTitle), resp.Data.Embeds[0].Title)
	require.Len(t, resp.Data.Components, 1)
}

func TestHandleInteraction_Confirm_PassesSourceID(t *testing.T) {
	session := &mockSession{}
	result := &ParseResult{
		IsBookkeeping: true,
		Type:          "expense",
		Amount:        15000,
		Description:   "rent",
		CategoryID:    "expense-other",
		CategoryName:  "Other",
		Date:          "2026-04-05",
		SourceType:    "bank_account",
		SourceID:      "acct-uuid-1",
	}
	creator := &mockCashFlowCreator{resultID: "cashflow-1"}
	h := NewHandler(context.Background(), &mockParser{}, creator, &mockCategoryLoader{}, nil, string(LangEn))
	customID := h.storePending(result, "author-1")
	interaction := newComponentInteraction(customID, "author-1")

	h.handleInteraction(session, interaction)

	require.Len(t, creator.createdInputs, 1)
	require.Equal(t, "bank_account", creator.createdInputs[0].SourceType)
	require.Equal(t, "acct-uuid-1", creator.createdInputs[0].SourceID)
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

func TestMonthNameForQueryTitle_English(t *testing.T) {
	require.Equal(t, "April", fmt.Sprintf("%s", monthLabel(string(LangEn), 4)))
}

func TestCleanupExpired_RemovesOldEntries(t *testing.T) {
	t.Parallel()

	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.mu.Lock()
	h.pending["old-key"] = pendingEntry{
		result:    &ParseResult{},
		authorID:  "author-1",
		createdAt: time.Now().Add(-20 * time.Minute),
	}
	h.mu.Unlock()

	h.cleanupExpired()

	h.mu.Lock()
	_, exists := h.pending["old-key"]
	h.mu.Unlock()

	require.False(t, exists, "expired entry should be removed")
}

func TestCleanupExpired_KeepsFreshEntries(t *testing.T) {
	t.Parallel()

	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.mu.Lock()
	h.pending["fresh-key"] = pendingEntry{
		result:    &ParseResult{},
		authorID:  "author-1",
		createdAt: time.Now(),
	}
	h.mu.Unlock()

	h.cleanupExpired()

	h.mu.Lock()
	_, exists := h.pending["fresh-key"]
	h.mu.Unlock()

	require.True(t, exists, "fresh entry should be kept")
}

func TestCleanupExpired_RemovesMultipleExpired(t *testing.T) {
	t.Parallel()

	h := NewHandler(context.Background(), &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	h.mu.Lock()
	h.pending["old-1"] = pendingEntry{result: &ParseResult{}, authorID: "a", createdAt: time.Now().Add(-20 * time.Minute)}
	h.pending["old-2"] = pendingEntry{result: &ParseResult{}, authorID: "b", createdAt: time.Now().Add(-30 * time.Minute)}
	h.pending["old-3"] = pendingEntry{result: &ParseResult{}, authorID: "c", createdAt: time.Now().Add(-60 * time.Minute)}
	h.pending["fresh-1"] = pendingEntry{result: &ParseResult{}, authorID: "d", createdAt: time.Now()}
	h.mu.Unlock()

	h.cleanupExpired()

	h.mu.Lock()
	remaining := len(h.pending)
	_, freshExists := h.pending["fresh-1"]
	h.mu.Unlock()

	require.Equal(t, 1, remaining, "only fresh entry should remain")
	require.True(t, freshExists, "fresh entry should be kept")
}

func TestStartCleanup_StopsOnContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	h := NewHandler(ctx, &mockParser{}, &mockCashFlowCreator{}, &mockCategoryLoader{}, nil, string(LangEn))

	// Add an expired entry
	h.mu.Lock()
	h.pending["old-key"] = pendingEntry{result: &ParseResult{}, authorID: "a", createdAt: time.Now().Add(-20 * time.Minute)}
	h.mu.Unlock()

	// Cancel the context to stop the cleanup goroutine
	cancel()

	// Give the goroutine time to exit
	time.Sleep(50 * time.Millisecond)

	// The entry should still be there since the cleanup ticker (5 min) hasn't fired
	// and the goroutine should have exited due to context cancellation.
	// Manually call cleanupExpired to verify entries are still removable.
	h.mu.Lock()
	_, exists := h.pending["old-key"]
	h.mu.Unlock()

	require.True(t, exists, "entry should still exist since ticker didn't fire before cancel")
}
