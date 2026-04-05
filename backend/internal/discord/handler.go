package discord

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

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
}

// CategoryLoader loads categories from the database.
type CategoryLoader interface {
	LoadCategories() ([]CategoryInfo, error)
}

type pendingEntry struct {
	result   *ParseResult
	authorID string
}

// Handler processes Discord messages and button interactions for bookkeeping.
type Handler struct {
	parser  Parser
	creator CashFlowCreator
	catRepo CategoryLoader
	lang    string
	mu      sync.Mutex
	pending map[string]pendingEntry
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

func NewHandler(parser Parser, creator CashFlowCreator, catLoader CategoryLoader, lang string) *Handler {
	if strings.TrimSpace(lang) == "" {
		lang = string(LangZhTW)
	}

	return &Handler{
		parser:  parser,
		creator: creator,
		catRepo: catLoader,
		lang:    lang,
		pending: make(map[string]pendingEntry),
	}
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
		h.sendText(s, m.ChannelID, GetMessage(h.lang, MsgSystemError))
		return
	}

	result, err := h.parser.Parse(context.Background(), m.Content, categories)
	if err != nil {
		h.sendText(s, m.ChannelID, GetMessage(h.lang, MsgSystemError))
		return
	}
	if result == nil || !result.IsBookkeeping {
		return
	}
	if hasMissingField(result.MissingFields, "amount") {
		h.sendText(s, m.ChannelID, GetMessage(h.lang, MsgMissingAmount)+"\n"+GetMessage(h.lang, MsgUsageExamples))
		return
	}

	confirmID := h.storePending(result, m.Author.ID)
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
					CustomID: "cancel:" + m.Author.ID,
				},
			}},
		},
	}
	_, _ = s.ChannelMessageSendComplex(m.ChannelID, preview)
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
		})
		if err != nil {
			h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgBookingFailed), err.Error())
			return
		}

		h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgConfirmSuccess), "")
	case "cancel":
		h.respondWithUpdatedEmbed(s, i, GetMessage(h.lang, MsgCancelled), "")
	}
}

func (h *Handler) loadCategories() ([]CategoryInfo, error) {
	if h.catRepo == nil {
		return nil, nil
	}
	return h.catRepo.LoadCategories()
}

func (h *Handler) sendText(s discordSession, channelID, content string) {
	_, _ = s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{Content: content})
}

func (h *Handler) buildPreviewEmbed(result *ParseResult) *discordgo.MessageEmbed {
	typeLabel := GetMessage(h.lang, MsgTypeExpense)
	color := 0xFF0000
	if result.Type == "income" {
		typeLabel = GetMessage(h.lang, MsgTypeIncome)
		color = 0x00FF00
	}

	return &discordgo.MessageEmbed{
		Title: GetMessage(h.lang, MsgPreviewTitle),
		Color: color,
		Fields: []*discordgo.MessageEmbedField{
			{Name: GetMessage(h.lang, MsgFieldType), Value: typeLabel, Inline: true},
			{Name: GetMessage(h.lang, MsgFieldAmount), Value: formatAmount(result.Amount), Inline: true},
			{Name: GetMessage(h.lang, MsgFieldCategory), Value: fallbackText(result.CategoryName), Inline: true},
			{Name: GetMessage(h.lang, MsgFieldDescription), Value: fallbackText(result.Description)},
			{Name: GetMessage(h.lang, MsgFieldDate), Value: fallbackText(result.Date), Inline: true},
		},
	}
}

func (h *Handler) respondWithUpdatedEmbed(s discordSession, i *discordgo.InteractionCreate, title, description string) {
	embed := cloneFirstEmbed(i.Message)
	embed.Title = title
	embed.Description = description
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
	case len(parts) == 2 && parts[0] == "cancel":
		return parts[0], "", parts[1], true
	default:
		return "", "", "", false
	}
}

func (h *Handler) storePending(result *ParseResult, authorID string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	key := hex.EncodeToString(b)

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
