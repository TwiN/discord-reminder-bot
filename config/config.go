package config

import (
	"os"
	"strings"
)

var cfg *Config

type Config struct {
	DiscordToken  string
	CommandPrefix string
}

func load() {
	cfg = &Config{
		DiscordToken:  strings.TrimSpace(os.Getenv("DISCORD_BOT_TOKEN")),
		CommandPrefix: strings.TrimSpace(os.Getenv("COMMAND_PREFIX")),
	}
	if len(cfg.DiscordToken) == 0 {
		panic("environment variable 'DISCORD_BOT_TOKEN' must not be empty")
	}
	if len(cfg.CommandPrefix) == 0 {
		cfg.CommandPrefix = "!"
	}
}

func Get() *Config {
	if cfg == nil {
		load()
	}
	return cfg
}
