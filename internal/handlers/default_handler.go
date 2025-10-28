package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func (h *Handler) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := inline.New(b).
		Row().
		Button("Create gif", []byte("create"), h.onInlineKeyboardSelect)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Main Menu",
		ReplyMarkup: kb,
	})
}

func (h *Handler) onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	switch string(data) {
	case "create":
		h.CreateHandler(ctx, b, &models.Update{
			Message: mes.Message,
		})

	default:
		h.DefaultHandler(ctx, b, &models.Update{
			Message: mes.Message,
		})
	}
}
