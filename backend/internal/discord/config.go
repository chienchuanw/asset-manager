package discord

import (
	"os"
	"strings"
)

// Config holds the Discord Bot configuration read from environment variables.
type Config struct {
	Enabled    bool
	Token      string
	ChannelIDs []string
	Lang       string
	GeminiKey  string
}

// LoadConfig reads Discord Bot configuration from environment variables
// and applies sensible defaults.
func LoadConfig() Config {
	cfg := Config{
		Enabled:   os.Getenv("DISCORD_BOT_ENABLED") == "true",
		Token:     os.Getenv("DISCORD_BOT_TOKEN"),
		GeminiKey: os.Getenv("GEMINI_API_KEY"),
		Lang:      os.Getenv("DISCORD_BOT_LANG"),
	}

	if raw := os.Getenv("DISCORD_CHANNEL_IDS"); raw != "" {
		for _, id := range strings.Split(raw, ",") {
			id = strings.TrimSpace(id)
			if id != "" {
				cfg.ChannelIDs = append(cfg.ChannelIDs, id)
			}
		}
	}

	if cfg.Lang == "" {
		cfg.Lang = "zh-TW"
	}

	return cfg
}
