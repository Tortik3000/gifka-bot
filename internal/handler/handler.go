package handler

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	logger         *zap.Logger
	mediaUseCase   mediaUseCase
	convUseCase    conversationUseCase
	sessionUseCase sessionUseCase
}

func New(
	logger *zap.Logger,
	mediaUseCase mediaUseCase,
	convUseCase conversationUseCase,
	sessionUseCase sessionUseCase,
) *Handler {
	return &Handler{
		logger:         logger,
		mediaUseCase:   mediaUseCase,
		convUseCase:    convUseCase,
		sessionUseCase: sessionUseCase,
	}
}

func (h *Handler) Callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}
	data := []byte(update.CallbackQuery.Data)
	h.AddTextCallback(ctx, b, update.CallbackQuery.Message, data)
}
