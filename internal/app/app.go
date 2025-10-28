package app

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"

	"gifka-bot/config"
	"gifka-bot/internal/handlers"
)

func Run(logger *zap.Logger) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.New()
	service := handlers.New(logger)
	opts := []bot.Option{
		bot.WithDefaultHandler(service.DefaultHandler),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, service.StartHandler),
		bot.WithMiddlewares(service.ConversationMiddleware),
	}

	b, err := bot.New(cfg.TG.Token, opts...)
	if err != nil {
		logger.Fatal("can't initialize bot", zap.Error(err))
	}

	b.Start(ctx)
}
