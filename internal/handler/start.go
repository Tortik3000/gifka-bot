package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"gifka-bot/internal/messages"
)

func (h *Handler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	welcomeText := messages.StartMessage

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   welcomeText,
	})

	h.Default(ctx, b, update)
}
