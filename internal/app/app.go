package app

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"

	"gifka-bot/config"
	"gifka-bot/internal/handlers"
)

func Run(logger *zap.Logger) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New()
	service := handlers.New(logger)
	opts := []bot.Option{
		bot.WithDefaultHandler(service.CreateHandler),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, service.StartHandler),
		bot.WithMiddlewares(service.ConversationMiddleware),
	}

	b, err := bot.New(cfg.TG.Token, opts...)
	if err != nil {
		logger.Fatal("can't initialize bot", zap.Error(err))
	}

	b.Start(ctx)
}
