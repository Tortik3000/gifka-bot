package app

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"

	"gifka-bot/config"
	"gifka-bot/internal/handler"
	"gifka-bot/internal/media_processor"
	"gifka-bot/internal/middleware"
	"gifka-bot/internal/session"
	"gifka-bot/internal/usecase"
)

func Run(logger *zap.Logger) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New()

	if err := runBot(ctx, cfg, logger); err != nil {
		logger.Fatal("bot stopped with error", zap.Error(err))
	}
}

func runBot(ctx context.Context, cfg *config.Config, logger *zap.Logger) error {
	mp := media_processor.New()
	sessionStorage := session.NewInMemoryStorage()
	sessionManager := session.NewManager(sessionStorage)

	mediaUseCase := usecase.NewMediaService(mp, logger)
	convUseCase := usecase.NewConversationService(sessionManager, logger)
	sessionUseCase := usecase.NewSessionService(sessionManager)

	h := handler.New(logger, mediaUseCase, convUseCase, sessionUseCase)

	conversation := middleware.NewConversation(sessionManager, convUseCase, h)

	opts := []bot.Option{
		bot.WithDefaultHandler(h.Default),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, h.Start),
		bot.WithMiddlewares(conversation.Handle),
	}

	b, err := bot.New(cfg.TG.Token, opts...)
	if err != nil {
		return err
	}

	logger.Info("telegram bot started")
	b.Start(ctx)
	return nil
}
