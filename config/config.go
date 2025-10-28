package config

import (
	"os"
)

type (
	Config struct {
		TG
	}
	TG struct {
		Token string `env:"TELEGRAM_BOT_TOKEN"`
	}
)

func New() *Config {
	cfg := &Config{}

	cfg.TG.Token = os.Getenv("TELEGRAM_BOT_TOKEN")

	return cfg
}
