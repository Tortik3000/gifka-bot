// internal/app/app.go
package app

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"

	"gifka-bot/config"
	"gifka-bot/internal/handler"
	"gifka-bot/internal/handler/middleware"
	"gifka-bot/internal/media_processor"
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
	// 1. инфраструктурные зависимости
	mp := media_processor.New() // ваш процессор GIF/стикеров
	sessionStorage := session.NewInMemoryStorage()
	sessionManager := session.NewManager(sessionStorage)

	// 2. сервисы (бизнес-логика)
	mediaUseCase := usecase.NewMediaService(mp, logger)
	convUseCase := usecase.NewConversationService(sessionManager, logger)
	sessionUseCase := usecase.NewSessionService(sessionManager)

	// 3. HTTP/Telegram handlers
	h := handler.New(logger, mediaUseCase, convUseCase, sessionUseCase)

	//
	conversation := middleware.NewConversation(sessionManager, convUseCase, h)

	// 4. сборка опций бота
	opts := []bot.Option{
		bot.WithDefaultHandler(h.Create), // раньше CreateHandler
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, h.Start),
		bot.WithMiddlewares(conversation.Handle), // middleware как метод Handler
		//bot.WithHTTPClient(10*time.Second, bot.HttpClient()),
	}

	b, err := bot.New(cfg.TG.Token, opts...)
	if err != nil {
		return err
	}

	logger.Info("telegram bot started")
	b.Start(ctx)
	return nil
}
