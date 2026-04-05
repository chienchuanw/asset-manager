package discord

import (
	"errors"
	"log"
	"runtime/debug"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// MessageHandler handles Discord messages and interactions.
type MessageHandler interface {
	HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate)
	HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Bot manages the Discord gateway connection and event routing.
type Bot struct {
	session         *discordgo.Session
	allowedChannels map[string]bool
	botUserID       string
	lang            string
	handler         MessageHandler
	stopCh          chan struct{}
}

// NewBot creates a Discord bot connection manager.
func NewBot(cfg Config) (*Bot, error) {
	if strings.TrimSpace(cfg.Token) == "" {
		return nil, errors.New("discord token is required")
	}

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	return &Bot{
		session:         session,
		allowedChannels: parseChannelIDs(strings.Join(cfg.ChannelIDs, ",")),
		lang:            cfg.Lang,
		stopCh:          make(chan struct{}),
	}, nil
}

// Start opens the Discord gateway connection and registers handlers.
func (b *Bot) Start() error {
	b.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		b.onMessage(s, m)
	})
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		b.onInteraction(s, i)
	})

	if err := b.session.Open(); err != nil {
		return err
	}

	if b.session.State != nil && b.session.State.User != nil {
		b.botUserID = b.session.State.User.ID
	}

	log.Print("Discord bot connected")
	return nil
}

// Stop closes the Discord gateway connection.
func (b *Bot) Stop() error {
	select {
	case <-b.stopCh:
	default:
		close(b.stopCh)
	}

	if err := b.session.Close(); err != nil {
		return err
	}

	log.Print("Discord bot disconnected")
	return nil
}

// SetHandler sets the runtime message handler.
func (b *Bot) SetHandler(h MessageHandler) {
	b.handler = h
}

func (b *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("discord message handler panic: %v\n%s", r, debug.Stack())
		}
	}()

	if !shouldProcessMessage(m, b.botUserID, b.allowedChannels) {
		return
	}

	if b.handler != nil {
		b.handler.HandleMessage(s, m)
	}
}

func (b *Bot) onInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("discord interaction handler panic: %v", r)
		}
	}()

	if b.handler != nil {
		b.handler.HandleInteraction(s, i)
	}
}

func shouldProcessMessage(msg *discordgo.MessageCreate, botUserID string, allowedChannels map[string]bool) bool {
	if msg == nil || msg.Message == nil || msg.Author == nil {
		return false
	}

	if msg.Author.ID == botUserID {
		return false
	}

	if msg.Author.Bot {
		return false
	}

	if len(allowedChannels) == 0 {
		return false
	}

	return allowedChannels[msg.ChannelID]
}

func parseChannelIDs(raw string) map[string]bool {
	allowedChannels := make(map[string]bool)
	for _, channelID := range strings.Split(raw, ",") {
		channelID = strings.TrimSpace(channelID)
		if channelID != "" {
			allowedChannels[channelID] = true
		}
	}
	return allowedChannels
}
