// internal/handler/handler.go
package handler

import (
	"context"

	"gifka-bot/internal/usecase"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type Handler struct {
	logger     *zap.Logger
	mediaSvc   *usecase.MediaService
	convSvc    *usecase.ConversationService
	sessionSvc *usecase.SessionService
}

func New(
	logger *zap.Logger,
	mediaSvc *usecase.MediaService,
	convSvc *usecase.ConversationService,
	sessSvc *usecase.SessionService,
) *Handler {
	return &Handler{
		logger:     logger,
		mediaSvc:   mediaSvc,
		convSvc:    convSvc,
		sessionSvc: sessSvc,
	}
}

// Default — если ничего не сматчилось
func (h *Handler) Default(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Use /start to get started.",
		})
	}
}

// Callback для inline‑кнопок
func (h *Handler) Callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}
	data := []byte(update.CallbackQuery.Data)
	h.AddTextCallback(ctx, b, update.CallbackQuery.Message, data)
}
