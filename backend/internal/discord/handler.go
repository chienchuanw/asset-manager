package discord

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const processingEmoji = "⏳"

// CashFlowCreator abstracts cash flow record creation for testability.
type CashFlowCreator interface {
	CreateCashFlowFromBot(input *BotCashFlowInput) (string, error)
}

// BotCashFlowInput holds the data needed to create a cash flow from the bot.
type BotCashFlowInput struct {
	Date        string
	Type        string
	CategoryID  string
	Amount      float64
	Description string
	SourceType  string
	SourceID    string
}

// CategoryLoader loads categories from the database.
type CategoryLoader interface {
	LoadCategories() ([]CategoryInfo, error)
}

// AccountInfo represents a bank account or credit card for display in Select Menu.
type AccountInfo struct {
	ID   string
	Name string
	Type string
}

// AccountLoader loads bank accounts and credit cards from the database.
type AccountLoader interface {
	LoadAccounts(sourceType string) ([]AccountInfo, error)
}

// CategoryBreakdown holds a single category's spending/income data for query results.
type CategoryBreakdown struct {
	Name   string
	Amount float64
	Count  int
}

// MonthComparisonResult holds month-over-month comparison data for query results.
type MonthComparisonResult struct {
	PreviousMonth    int
	PreviousYear     int
	ExpenseChange    float64
	ExpenseChangePct float64
	IncomeChange     float64
	IncomeChangePct  float64
}

// MonthlySummaryResult holds the data returned by a cash flow summary query.
type MonthlySummaryResult struct {
	Year          int
	Month         int
	TotalIncome   float64
	TotalExpense  float64
	NetCashFlow   float64
	IncomeCount   int
	ExpenseCount  int
	TopCategories []CategoryBreakdown
	Comparison    *MonthComparisonResult
}

// CashFlowQuerier abstracts cash flow summary queries for the bot.
type CashFlowQuerier interface {
	GetMonthlySummary(year, month int) (*MonthlySummaryResult, error)
}

// BankAccountBalance holds a single bank account's balance data for query results.
type BankAccountBalance struct {
	Name     string
	Last4    string
	Currency string
	Balance  float64
}

// CreditCardBalance holds a single credit card's balance data for query results.
type CreditCardBalance struct {
	Name        string
	Last4       string
	CreditLimit float64
	UsedCredit  float64
	Remaining   float64
	UsagePct    float64
}

// AccountBalancesResult holds the combined bank + credit card balance data.
type AccountBalancesResult struct {
	BankAccounts []BankAccountBalance
	CreditCards  []CreditCardBalance
	BankError    error
	CCError      error
}

// AccountBalanceQuerier abstracts account balance queries for the bot.
type AccountBalanceQuerier interface {
	GetAllBalances() (*AccountBalancesResult, error)
}

type pendingEntry struct {
	result            *ParseResult
	authorID          string
	awaitingAccount   bool
	awaitingAccountID bool
	ccPayment         bool
	ccCardID          string
	ccCardName        string
	ccBankID          string
	ccBankName        string
	ccAmount          float64
	ccPaymentType     string
}

// Handler processes Discord messages and button interactions for bookkeeping.
type Handler struct {
	parser           Parser
	creator          CashFlowCreator
	ccPaymentCreator CreditCardPaymentCreator
	catRepo          CategoryLoader
	acctLoader       AccountLoader
	cfQuerier        CashFlowQuerier
	acctQuerier      AccountBalanceQuerier
	lang             string
	mu               sync.Mutex
	pending          map[string]pendingEntry
}

type discordSession interface {
	ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend) (*discordgo.Message, error)
	MessageReactionAdd(channelID, messageID, emojiID string) error
	MessageReactionRemove(channelID, messageID, emojiID, userID string) error
	ChannelMessageEditComplex(data *discordgo.MessageEdit) (*discordgo.Message, error)
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error
}

type realDiscordSession struct {
	session *discordgo.Session
}

func (r realDiscordSession) ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend) (*discordgo.Message, error) {
	return r.session.ChannelMessageSendComplex(channelID, data)
}

func (r realDiscordSession) MessageReactionAdd(channelID, messageID, emojiID string) error {
	return r.session.MessageReactionAdd(channelID, messageID, emojiID)
}

func (r realDiscordSession) MessageReactionRemove(channelID, messageID, emojiID, userID string) error {
	return r.session.MessageReactionRemove(channelID, messageID, emojiID, userID)
}

func (r realDiscordSession) ChannelMessageEditComplex(data *discordgo.MessageEdit) (*discordgo.Message, error) {
	return r.session.ChannelMessageEditComplex(data)
}

func (r realDiscordSession) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	return r.session.InteractionRespond(interaction, resp)
}

func NewHandler(parser Parser, creator CashFlowCreator, catLoader CategoryLoader, acctLoader AccountLoader, lang string, opts ...HandlerOption) *Handler {
	if strings.TrimSpace(lang) == "" {
		lang = string(LangZhTW)
	}

	h := &Handler{
		parser:     parser,
		creator:    creator,
		catRepo:    catLoader,
		acctLoader: acctLoader,
		lang:       lang,
		pending:    make(map[string]pendingEntry),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// HandlerOption configures optional Handler dependencies.
type HandlerOption func(*Handler)

func WithCashFlowQuerier(q CashFlowQuerier) HandlerOption {
	return func(h *Handler) { h.cfQuerier = q }
}

func WithAccountBalanceQuerier(q AccountBalanceQuerier) HandlerOption {
	return func(h *Handler) { h.acctQuerier = q }
}

func WithCCPaymentCreator(c CreditCardPaymentCreator) HandlerOption {
	return func(h *Handler) { h.ccPaymentCreator = c }
}

func (h *Handler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s == nil {
		return
	}
	h.handleMessage(realDiscordSession{session: s}, m)
}

func (h *Handler) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if s == nil {
		return
	}
	h.handleInteraction(realDiscordSession{session: s}, i)
}

func (h *Handler) handleMessage(s discordSession, m *discordgo.MessageCreate) {
	if s == nil || m == nil || m.Message == nil {
		return
	}

	_ = s.MessageReactionAdd(m.ChannelID, m.ID, processingEmoji)
	defer func() {
		_ = s.MessageReactionRemove(m.ChannelID, m.ID, processingEmoji, "@me")
	}()

	categories, err := h.loadCategories()
	if err != nil {
		log.Printf("discord: failed to load categories: %v", err)
		h.sendText(s, m.ChannelID, GetMessage(h.lang, MsgDataLoadFailed))
		return
	}

	result, err := h.parser.Parse(context.Background(), m.Content, categories)
	if err != nil {
		log.Printf("discord: failed to parse message %q: %v", m.Content, err)
		h.sendText(s, m.ChannelID, GetMessage(h.lang, MsgParseFailed))
		return
	}
	if result == nil || (!result.IsBookkeeping && result.Action == "") {
		return
	}
	if result.Action == "chat" {
		h.sendText(s, m.ChannelID, GetMessage(h.lang, MsgChatGreeting))
		return
	}
	if result.Action == "unsupported" {
		h.sendText(s, m.ChannelID, GetMessage(h.lang, MsgUnsupported))
		return
	}
	if result.Action == "query" {
		h.handleQuery(s, m.ChannelID, result)
		return
	}
	if result.Action == "cc_payment" {
		h.handleCCPayment(s, m.ChannelID, result, m.Author.ID)
		return
	}
	if hasMissingField(result.MissingFields, "amount") {
		h.sendText(s, m.ChannelID, GetMessage(h.lang, MsgMissingAmount)+"\n"+GetMessage(h.lang, MsgUsageExamples))
		return
	}

	if result.SourceType == "" {
		h.sendAccountSelectMenu(s, m.ChannelID, result, m.Author.ID)
		return
	}

	if result.SourceType != "cash" {
		h.sendAccountIDMenu(s, m.ChannelID, result, m.Author.ID)
		return
	}

	h.sendPreview(s, m.ChannelID, result, m.Author.ID)
}

func (h *Handler) handleInteraction(s discordSession, i *discordgo.InteractionCreate) {
	if s == nil || i == nil || i.Interaction == nil {
		return
	}

	data := i.MessageComponentData()
	action, payload, authorID, ok := parseCustomID(data.CustomID)
	if !ok {
		return
	}
	if interactionUserID(i) != authorID {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: GetMessage(h.lang, MsgOnlyAuthor),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	switch action {
	case "select_account":
		h.handleAccountSelection(s, i, payload, data.Values)
	case "select_account_id":
		h.handleAccountIDSelection(s, i, payload, data.Values)
	case "select_cc":
		h.handleCCCardSelection(s, i, payload, data.Values)
	case "select_cc_bank":
		h.handleCCBankSelection(s, i, payload, data.Values)
	case "confirm":
		result, ok := h.popPending(payload)
		if !ok {
			h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgExpired), "")
			return
		}

		_, err := h.creator.CreateCashFlowFromBot(&BotCashFlowInput{
			Date:        result.Date,
			Type:        result.Type,
			CategoryID:  result.CategoryID,
			Amount:      result.Amount,
			Description: result.Description,
			SourceType:  result.SourceType,
			SourceID:    result.SourceID,
		})
		if err != nil {
			log.Printf("discord: failed to create cash flow: %v", err)
			h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgBookingFailed), err.Error())
			return
		}

		h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgConfirmSuccess), "")
	case "confirm_cc_payment":
		entry, ok := h.popPendingEntry(payload)
		if !ok {
			h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgExpired), "")
			return
		}
		if h.ccPaymentCreator == nil {
			h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgCCPaymentFailed), GetMessage(h.lang, MsgSystemError))
			return
		}

		_, _, err := h.ccPaymentCreator.CreatePaymentFromBot(&BotCCPaymentInput{
			CreditCardID:  entry.ccCardID,
			BankAccountID: entry.ccBankID,
			Amount:        entry.ccAmount,
			Date:          entry.result.Date,
			PaymentType:   entry.ccPaymentType,
			CategoryID:    entry.result.CategoryID,
		})
		if err != nil {
			log.Printf("discord: failed to create cc payment: %v", err)
			h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgCCPaymentFailed), err.Error())
			return
		}

		h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgCCPaymentSuccess), "")
	case "cancel":
		h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgCancelled), "")
	}
}

func (h *Handler) handleCCPayment(s discordSession, channelID string, result *ParseResult, authorID string) {
	if result.PaymentType != "full" && hasMissingField(result.MissingFields, "amount") {
		h.sendText(s, channelID, GetMessage(h.lang, MsgCCPaymentMissingAmount)+"\n"+GetMessage(h.lang, MsgCCPaymentUsageExamples))
		return
	}

	if h.acctLoader == nil {
		h.sendText(s, channelID, GetMessage(h.lang, MsgCCPaymentNoCards))
		return
	}

	cards, err := h.acctLoader.LoadAccounts("credit_card")
	if err != nil || len(cards) == 0 {
		h.sendText(s, channelID, GetMessage(h.lang, MsgCCPaymentNoCards))
		return
	}

	key := randomHexKey()
	h.mu.Lock()
	h.pending[key] = pendingEntry{
		result:        result,
		authorID:      authorID,
		ccPayment:     true,
		ccAmount:      result.Amount,
		ccPaymentType: result.PaymentType,
	}
	h.mu.Unlock()

	options := make([]discordgo.SelectMenuOption, 0, len(cards))
	for _, card := range cards {
		options = append(options, discordgo.SelectMenuOption{Label: card.Name, Value: card.ID})
	}

	customID := "select_cc:" + key + ":" + authorID
	msg := &discordgo.MessageSend{
		Content: GetMessage(h.lang, MsgCCPaymentSelectCard),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				&discordgo.SelectMenu{CustomID: customID, Placeholder: GetMessage(h.lang, MsgCCPaymentSelectCard), Options: options},
			}},
		},
	}
	_, _ = s.ChannelMessageSendComplex(channelID, msg)
}

func (h *Handler) handleAccountSelection(s discordSession, i *discordgo.InteractionCreate, pendingKey string, values []string) {
	if len(values) == 0 {
		return
	}

	selectedType := values[0]

	h.mu.Lock()
	entry, ok := h.pending[pendingKey]
	if !ok {
		h.mu.Unlock()
		return
	}
	entry.result.SourceType = selectedType
	entry.awaitingAccount = false

	if selectedType == "cash" {
		delete(h.pending, pendingKey)
		h.mu.Unlock()
		h.respondWithPreview(s, i, entry.result, entry.authorID)
		return
	}

	entry.awaitingAccountID = true
	h.pending[pendingKey] = entry
	h.mu.Unlock()

	h.respondWithAccountIDMenu(s, i, pendingKey, entry.authorID, selectedType)
}

func (h *Handler) handleAccountIDSelection(s discordSession, i *discordgo.InteractionCreate, pendingKey string, values []string) {
	if len(values) == 0 {
		return
	}

	h.mu.Lock()
	entry, ok := h.pending[pendingKey]
	if !ok {
		h.mu.Unlock()
		return
	}
	entry.result.SourceID = values[0]
	entry.result.SourceName = h.lookupAccountName(entry.result.SourceType, values[0])
	entry.awaitingAccountID = false
	delete(h.pending, pendingKey)
	h.mu.Unlock()

	h.respondWithPreview(s, i, entry.result, entry.authorID)
}

func (h *Handler) handleCCCardSelection(s discordSession, i *discordgo.InteractionCreate, pendingKey string, values []string) {
	if len(values) == 0 {
		return
	}

	h.mu.Lock()
	entry, ok := h.pending[pendingKey]
	if !ok {
		h.mu.Unlock()
		return
	}
	entry.ccCardID = values[0]
	entry.ccCardName = h.lookupAccountName("credit_card", values[0])
	h.pending[pendingKey] = entry
	h.mu.Unlock()

	var accounts []AccountInfo
	if h.acctLoader != nil {
		accounts, _ = h.acctLoader.LoadAccounts("bank_account")
	}
	if len(accounts) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    GetMessage(h.lang, MsgCCPaymentNoBankAccounts),
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	options := make([]discordgo.SelectMenuOption, 0, len(accounts))
	for _, acct := range accounts {
		options = append(options, discordgo.SelectMenuOption{Label: acct.Name, Value: acct.ID})
	}
	customID := "select_cc_bank:" + pendingKey + ":" + entry.authorID
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: GetMessage(h.lang, MsgCCPaymentSelectBank),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					&discordgo.SelectMenu{CustomID: customID, Placeholder: GetMessage(h.lang, MsgCCPaymentSelectBank), Options: options},
				}},
			},
		},
	})
}

func (h *Handler) handleCCBankSelection(s discordSession, i *discordgo.InteractionCreate, pendingKey string, values []string) {
	if len(values) == 0 {
		return
	}

	h.mu.Lock()
	entry, ok := h.pending[pendingKey]
	if !ok {
		h.mu.Unlock()
		return
	}
	entry.ccBankID = values[0]
	entry.ccBankName = h.lookupAccountName("bank_account", values[0])
	delete(h.pending, pendingKey)
	h.mu.Unlock()

	confirmKey := randomHexKey()
	h.mu.Lock()
	h.pending[confirmKey] = entry
	h.mu.Unlock()

	confirmID := "confirm_cc_payment:" + confirmKey + ":" + entry.authorID
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{h.buildCCPaymentPreviewEmbed(&entry)},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					&discordgo.Button{Label: GetMessage(h.lang, MsgCCPaymentConfirmButton), Style: discordgo.SuccessButton, CustomID: confirmID},
					&discordgo.Button{Label: GetMessage(h.lang, MsgCancelButton), Style: discordgo.DangerButton, CustomID: "cancel:" + entry.authorID},
				}},
			},
		},
	})
}

func (h *Handler) respondWithPreview(s discordSession, i *discordgo.InteractionCreate, result *ParseResult, authorID string) {
	confirmID := h.storePending(result, authorID)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "",
			Embeds:  []*discordgo.MessageEmbed{h.buildPreviewEmbed(result)},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label:    GetMessage(h.lang, MsgConfirmButton),
						Style:    discordgo.SuccessButton,
						CustomID: confirmID,
					},
					&discordgo.Button{
						Label:    GetMessage(h.lang, MsgCancelButton),
						Style:    discordgo.DangerButton,
						CustomID: "cancel:" + authorID,
					},
				}},
			},
		},
	})
}

func (h *Handler) respondWithAccountIDMenu(s discordSession, i *discordgo.InteractionCreate, pendingKey, authorID, sourceType string) {
	placeholder := GetMessage(h.lang, MsgSelectBankAccount)
	if sourceType == "credit_card" {
		placeholder = GetMessage(h.lang, MsgSelectCreditCard)
	}

	var options []discordgo.SelectMenuOption
	if h.acctLoader != nil {
		accounts, err := h.acctLoader.LoadAccounts(sourceType)
		if err == nil {
			for _, acct := range accounts {
				options = append(options, discordgo.SelectMenuOption{
					Label: acct.Name,
					Value: acct.ID,
				})
			}
		}
	}

	if len(options) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    GetMessage(h.lang, MsgNoAccountsFound),
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	customID := "select_account_id:" + pendingKey + ":" + authorID
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: placeholder,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					&discordgo.SelectMenu{
						CustomID:    customID,
						Placeholder: placeholder,
						Options:     options,
					},
				}},
			},
		},
	})
}

func (h *Handler) loadCategories() ([]CategoryInfo, error) {
	if h.catRepo == nil {
		return nil, nil
	}
	return h.catRepo.LoadCategories()
}

func (h *Handler) handleQuery(s discordSession, channelID string, result *ParseResult) {
	switch result.QueryType {
	case "cash_flow_summary":
		h.handleCashFlowQuery(s, channelID, result)
	case "account_balance":
		h.handleAccountBalanceQuery(s, channelID, result)
	default:
		h.sendText(s, channelID, GetMessage(h.lang, MsgQueryUnsupported))
	}
}

func (h *Handler) handleCashFlowQuery(s discordSession, channelID string, result *ParseResult) {
	if h.cfQuerier == nil || result == nil || result.QueryParams == nil {
		h.sendText(s, channelID, GetMessage(h.lang, MsgQueryFailed))
		return
	}

	params := result.QueryParams
	summary, err := h.cfQuerier.GetMonthlySummary(params.Year, params.Month)
	if err != nil || summary == nil {
		if err != nil {
			log.Printf("discord: failed to query cash flow summary: %v", err)
		}
		h.sendText(s, channelID, GetMessage(h.lang, MsgQueryFailed))
		return
	}

	if category := strings.TrimSpace(params.Category); category != "" {
		matched := false
		filtered := *summary
		for _, breakdown := range summary.TopCategories {
			if breakdown.Name != category {
				continue
			}
			filtered.TotalExpense = breakdown.Amount
			filtered.ExpenseCount = breakdown.Count
			filtered.NetCashFlow = filtered.TotalIncome - breakdown.Amount
			filtered.TopCategories = []CategoryBreakdown{breakdown}
			summary = &filtered
			matched = true
			break
		}
		if !matched {
			h.sendText(s, channelID, fmt.Sprintf(GetMessage(h.lang, MsgQueryCategoryNotFound), category))
			return
		}
	}

	_, _ = s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{Embeds: []*discordgo.MessageEmbed{h.buildCashFlowQueryEmbed(result, summary)}})
}

func (h *Handler) handleAccountBalanceQuery(s discordSession, channelID string, _ *ParseResult) {
	if h.acctQuerier == nil {
		h.sendText(s, channelID, GetMessage(h.lang, MsgQueryFailed))
		return
	}

	balances, err := h.acctQuerier.GetAllBalances()
	if err != nil || balances == nil {
		if err != nil {
			log.Printf("discord: failed to query account balances: %v", err)
		}
		h.sendText(s, channelID, GetMessage(h.lang, MsgQueryFailed))
		return
	}

	if len(balances.BankAccounts) == 0 && len(balances.CreditCards) == 0 && balances.BankError == nil && balances.CCError == nil {
		h.sendText(s, channelID, GetMessage(h.lang, MsgQueryNoAccounts))
		return
	}

	_, _ = s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{Embeds: []*discordgo.MessageEmbed{h.buildAccountBalanceEmbed(balances)}})
}

func (h *Handler) buildCashFlowQueryEmbed(result *ParseResult, summary *MonthlySummaryResult) *discordgo.MessageEmbed {
	title := fmt.Sprintf(GetMessage(h.lang, MsgQueryCashFlowTitle), monthTitleArg(h.lang, summary.Month))
	if result != nil && result.QueryParams != nil && strings.TrimSpace(result.QueryParams.Category) != "" {
		title = fmt.Sprintf(GetMessage(h.lang, MsgQueryCashFlowCategoryTitle), monthTitleArg(h.lang, summary.Month), strings.TrimSpace(result.QueryParams.Category))
	}

	count := summary.IncomeCount + summary.ExpenseCount
	fields := []*discordgo.MessageEmbedField{
		{Name: GetMessage(h.lang, MsgQueryTotalIncome), Value: "$" + formatAmount(summary.TotalIncome), Inline: true},
		{Name: GetMessage(h.lang, MsgQueryTotalExpense), Value: "$" + formatAmount(summary.TotalExpense), Inline: true},
		{Name: GetMessage(h.lang, MsgQueryNetCashFlow), Value: "$" + formatAmount(summary.NetCashFlow), Inline: true},
		{Name: GetMessage(h.lang, MsgQueryFieldCount), Value: strconv.Itoa(count), Inline: true},
	}

	if len(summary.TopCategories) > 0 {
		lines := make([]string, 0, minInt(len(summary.TopCategories), 5))
		for i, category := range summary.TopCategories {
			if i >= 5 {
				break
			}
			lines = append(lines, fmt.Sprintf("%s: $%s", category.Name, formatAmount(category.Amount)))
		}
		fields = append(fields, &discordgo.MessageEmbedField{Name: GetMessage(h.lang, MsgQueryTopCategories), Value: strings.Join(lines, "\n")})
	}

	if summary.Comparison != nil {
		comparison := fmt.Sprintf(
			"%s: $%s (%.1f%%)\n%s: $%s (%.1f%%)",
			GetMessage(h.lang, MsgQueryTotalExpense),
			formatAmount(summary.Comparison.ExpenseChange),
			summary.Comparison.ExpenseChangePct,
			GetMessage(h.lang, MsgQueryTotalIncome),
			formatAmount(summary.Comparison.IncomeChange),
			summary.Comparison.IncomeChangePct,
		)
		fields = append(fields, &discordgo.MessageEmbedField{Name: GetMessage(h.lang, MsgQueryComparison), Value: comparison})
	}

	embed := &discordgo.MessageEmbed{
		Title:  title,
		Color:  0x3498DB,
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{Text: time.Now().Format("2006-01-02")},
	}

	if summary.TotalIncome == 0 && summary.TotalExpense == 0 && summary.NetCashFlow == 0 && count == 0 {
		embed.Description = GetMessage(h.lang, MsgQueryNoData)
	}

	return embed
}

func (h *Handler) buildAccountBalanceEmbed(result *AccountBalancesResult) *discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{}

	if result.BankError != nil || len(result.BankAccounts) > 0 {
		bankValue := GetMessage(h.lang, MsgQueryLoadFailed)
		if result.BankError == nil {
			var lines []string
			var total float64
			for _, account := range result.BankAccounts {
				lines = append(lines, fmt.Sprintf("%s *%s: $%s", account.Name, account.Last4, formatAmount(account.Balance)))
				total += account.Balance
			}
			lines = append(lines, fmt.Sprintf("%s: $%s", GetMessage(h.lang, MsgQueryBankTotal), formatAmount(total)))
			bankValue = strings.Join(lines, "\n")
		}
		fields = append(fields, &discordgo.MessageEmbedField{Name: GetMessage(h.lang, MsgQueryBankSection), Value: bankValue})
	}

	if result.CCError != nil || len(result.CreditCards) > 0 {
		ccValue := GetMessage(h.lang, MsgQueryLoadFailed)
		if result.CCError == nil {
			var lines []string
			for _, card := range result.CreditCards {
				line := fmt.Sprintf("%s *%s\n%s: $%s | %s: $%s | %s: $%s (%.0f%%)",
					card.Name,
					card.Last4,
					GetMessage(h.lang, MsgQueryCCLimit),
					formatAmount(card.CreditLimit),
					GetMessage(h.lang, MsgQueryCCUsed),
					formatAmount(card.UsedCredit),
					GetMessage(h.lang, MsgQueryCCRemaining),
					formatAmount(card.Remaining),
					card.UsagePct,
				)
				if card.UsagePct > 80 {
					line = "⚠️ " + line + " " + GetMessage(h.lang, MsgQueryCCNearLimit)
				}
				lines = append(lines, line)
			}
			ccValue = strings.Join(lines, "\n\n")
		}
		fields = append(fields, &discordgo.MessageEmbedField{Name: GetMessage(h.lang, MsgQueryCCSection), Value: ccValue})
	}

	return &discordgo.MessageEmbed{
		Title:  GetMessage(h.lang, MsgQueryAccountTitle),
		Color:  0x3498DB,
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{Text: time.Now().Format("2006-01-02")},
	}
}

func (h *Handler) sendText(s discordSession, channelID, content string) {
	_, _ = s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{Content: content})
}

func (h *Handler) sendAccountSelectMenu(s discordSession, channelID string, result *ParseResult, authorID string) {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	key := hex.EncodeToString(b)

	h.mu.Lock()
	h.pending[key] = pendingEntry{result: result, authorID: authorID, awaitingAccount: true}
	h.mu.Unlock()

	customID := "select_account:" + key + ":" + authorID
	msg := &discordgo.MessageSend{
		Content: GetMessage(h.lang, MsgSelectAccount),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				&discordgo.SelectMenu{
					CustomID:    customID,
					Placeholder: GetMessage(h.lang, MsgSelectAccount),
					Options: []discordgo.SelectMenuOption{
						{Label: GetMessage(h.lang, MsgAccountCash), Value: "cash"},
						{Label: GetMessage(h.lang, MsgAccountBank), Value: "bank_account"},
						{Label: GetMessage(h.lang, MsgAccountCreditCard), Value: "credit_card"},
					},
				},
			}},
		},
	}
	_, _ = s.ChannelMessageSendComplex(channelID, msg)
}

func (h *Handler) sendAccountIDMenu(s discordSession, channelID string, result *ParseResult, authorID string) {
	placeholder := GetMessage(h.lang, MsgSelectBankAccount)
	if result.SourceType == "credit_card" {
		placeholder = GetMessage(h.lang, MsgSelectCreditCard)
	}

	var options []discordgo.SelectMenuOption
	if h.acctLoader != nil {
		accounts, err := h.acctLoader.LoadAccounts(result.SourceType)
		if err == nil {
			for _, acct := range accounts {
				options = append(options, discordgo.SelectMenuOption{
					Label: acct.Name,
					Value: acct.ID,
				})
			}
		}
	}

	if len(options) == 0 {
		h.sendText(s, channelID, GetMessage(h.lang, MsgNoAccountsFound))
		return
	}

	b := make([]byte, 8)
	_, _ = rand.Read(b)
	key := hex.EncodeToString(b)

	h.mu.Lock()
	h.pending[key] = pendingEntry{result: result, authorID: authorID, awaitingAccountID: true}
	h.mu.Unlock()

	customID := "select_account_id:" + key + ":" + authorID
	msg := &discordgo.MessageSend{
		Content: placeholder,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				&discordgo.SelectMenu{
					CustomID:    customID,
					Placeholder: placeholder,
					Options:     options,
				},
			}},
		},
	}
	_, _ = s.ChannelMessageSendComplex(channelID, msg)
}

func (h *Handler) sendPreview(s discordSession, channelID string, result *ParseResult, authorID string) {
	confirmID := h.storePending(result, authorID)
	preview := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{h.buildPreviewEmbed(result)},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    GetMessage(h.lang, MsgConfirmButton),
					Style:    discordgo.SuccessButton,
					CustomID: confirmID,
				},
				&discordgo.Button{
					Label:    GetMessage(h.lang, MsgCancelButton),
					Style:    discordgo.DangerButton,
					CustomID: "cancel:" + authorID,
				},
			}},
		},
	}
	_, _ = s.ChannelMessageSendComplex(channelID, preview)
}

func (h *Handler) buildPreviewEmbed(result *ParseResult) *discordgo.MessageEmbed {
	typeLabel := GetMessage(h.lang, MsgTypeExpense)
	color := 0xFF0000
	if result.Type == "income" {
		typeLabel = GetMessage(h.lang, MsgTypeIncome)
		color = 0x00FF00
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: GetMessage(h.lang, MsgFieldType), Value: typeLabel, Inline: true},
		{Name: GetMessage(h.lang, MsgFieldAmount), Value: formatAmount(result.Amount), Inline: true},
		{Name: GetMessage(h.lang, MsgFieldCategory), Value: fallbackText(result.CategoryName), Inline: true},
		{Name: GetMessage(h.lang, MsgFieldDescription), Value: fallbackText(result.Description)},
		{Name: GetMessage(h.lang, MsgFieldPaymentMethod), Value: h.sourceTypeLabel(result.SourceType), Inline: true},
	}

	if result.SourceName != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name: GetMessage(h.lang, MsgFieldAccount), Value: result.SourceName, Inline: true,
		})
	}

	return &discordgo.MessageEmbed{
		Title:  GetMessage(h.lang, MsgPreviewTitle),
		Color:  color,
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{Text: fallbackText(result.Date)},
	}
}

func (h *Handler) buildCCPaymentPreviewEmbed(entry *pendingEntry) *discordgo.MessageEmbed {
	typeLabel := GetMessage(h.lang, MsgCCPaymentTypeCustom)
	switch entry.ccPaymentType {
	case "full":
		typeLabel = GetMessage(h.lang, MsgCCPaymentTypeFull)
	case "minimum":
		typeLabel = GetMessage(h.lang, MsgCCPaymentTypeMinimum)
	}

	return &discordgo.MessageEmbed{
		Title: GetMessage(h.lang, MsgCCPaymentPreview),
		Color: 0x3498DB,
		Fields: []*discordgo.MessageEmbedField{
			{Name: GetMessage(h.lang, MsgFieldAmount), Value: "$" + formatAmount(entry.ccAmount), Inline: true},
			{Name: GetMessage(h.lang, MsgCCPaymentFieldCard), Value: fallbackText(entry.ccCardName), Inline: true},
			{Name: GetMessage(h.lang, MsgCCPaymentFieldBank), Value: fallbackText(entry.ccBankName), Inline: true},
			{Name: GetMessage(h.lang, MsgCCPaymentFieldType), Value: typeLabel, Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: fallbackText(entry.result.Date)},
	}
}

func (h *Handler) lookupAccountName(sourceType, accountID string) string {
	if h.acctLoader == nil {
		return ""
	}
	accounts, err := h.acctLoader.LoadAccounts(sourceType)
	if err != nil {
		return ""
	}
	for _, acct := range accounts {
		if acct.ID == accountID {
			return acct.Name
		}
	}
	return ""
}

func (h *Handler) respondWithUpdatedEmbed(s discordSession, i *discordgo.InteractionCreate, title, description string) {
	embed := cloneFirstEmbed(i.Message)
	embed.Title = title
	embed.Description = description
	if title == GetMessage(h.lang, MsgConfirmSuccess) || title == GetMessage(h.lang, MsgCCPaymentSuccess) {
		embed.Color = 0x00CC00
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{},
		},
	})
}

func cloneFirstEmbed(message *discordgo.Message) *discordgo.MessageEmbed {
	if message != nil && len(message.Embeds) > 0 && message.Embeds[0] != nil {
		clone := *message.Embeds[0]
		return &clone
	}
	return &discordgo.MessageEmbed{}
}

func parseCustomID(customID string) (action string, payload string, authorID string, ok bool) {
	parts := strings.Split(customID, ":")
	switch {
	case len(parts) == 3 && parts[0] == "confirm":
		return parts[0], parts[1], parts[2], true
	case len(parts) == 3 && parts[0] == "confirm_cc_payment":
		return parts[0], parts[1], parts[2], true
	case len(parts) == 3 && parts[0] == "select_account":
		return parts[0], parts[1], parts[2], true
	case len(parts) == 3 && parts[0] == "select_account_id":
		return parts[0], parts[1], parts[2], true
	case len(parts) == 3 && parts[0] == "select_cc":
		return parts[0], parts[1], parts[2], true
	case len(parts) == 3 && parts[0] == "select_cc_bank":
		return parts[0], parts[1], parts[2], true
	case len(parts) == 2 && parts[0] == "cancel":
		return parts[0], "", parts[1], true
	default:
		return "", "", "", false
	}
}

func randomHexKey() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *Handler) storePending(result *ParseResult, authorID string) string {
	key := randomHexKey()

	h.mu.Lock()
	h.pending[key] = pendingEntry{result: result, authorID: authorID}
	h.mu.Unlock()

	return "confirm:" + key + ":" + authorID
}

func (h *Handler) popPending(key string) (*ParseResult, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	entry, ok := h.pending[key]
	if !ok {
		return nil, false
	}
	delete(h.pending, key)
	return entry.result, true
}

func (h *Handler) popPendingEntry(key string) (pendingEntry, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	entry, ok := h.pending[key]
	if !ok {
		return pendingEntry{}, false
	}
	delete(h.pending, key)
	return entry, true
}

func interactionUserID(i *discordgo.InteractionCreate) string {
	if i == nil || i.Interaction == nil {
		return ""
	}
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}
	if i.User != nil {
		return i.User.ID
	}
	return ""
}

func hasMissingField(fields []string, target string) bool {
	for _, field := range fields {
		if field == target {
			return true
		}
	}
	return false
}

func (h *Handler) sourceTypeLabel(sourceType string) string {
	switch sourceType {
	case "cash":
		return GetMessage(h.lang, MsgAccountCash)
	case "bank_account":
		return GetMessage(h.lang, MsgAccountBank)
	case "credit_card":
		return GetMessage(h.lang, MsgAccountCreditCard)
	default:
		return "-"
	}
}

func fallbackText(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func formatAmount(amount float64) string {
	if amount == math.Trunc(amount) {
		return formatIntegerWithCommas(int64(amount))
	}

	formatted := strconv.FormatFloat(amount, 'f', 2, 64)
	parts := strings.SplitN(formatted, ".", 2)
	if len(parts) == 1 {
		return formatIntegerWithCommas(parseInt64(parts[0]))
	}
	return fmt.Sprintf("%s.%s", formatIntegerWithCommas(parseInt64(parts[0])), parts[1])
}

func formatIntegerWithCommas(value int64) string {
	negative := value < 0
	if negative {
		value = -value
	}

	digits := strconv.FormatInt(value, 10)
	if len(digits) <= 3 {
		if negative {
			return "-" + digits
		}
		return digits
	}

	var parts []string
	for len(digits) > 3 {
		parts = append([]string{digits[len(digits)-3:]}, parts...)
		digits = digits[:len(digits)-3]
	}
	parts = append([]string{digits}, parts...)
	joined := strings.Join(parts, ",")
	if negative {
		return "-" + joined
	}
	return joined
}

func parseInt64(value string) int64 {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func monthLabel(lang string, month int) string {
	if month < 1 || month > 12 {
		return strconv.Itoa(month)
	}
	if lang == string(LangEn) {
		return time.Month(month).String()
	}
	return fmt.Sprintf("%d月", month)
}

func monthTitleArg(lang string, month int) any {
	if lang == string(LangEn) {
		return monthLabel(lang, month)
	}
	return month
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
