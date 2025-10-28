package config

import (
	"fmt"

	"github.com/spf13/viper"
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
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	cfg := &Config{}

	cfg.TG.Token = viper.GetString("tg_bot_token")

	return cfg
}
