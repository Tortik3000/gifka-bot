package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) ConversationMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update != nil && update.Message != nil {
			chatID := update.Message.Chat.ID
			s := getSession(chatID)
			switch s.Stage {
			case stageAwaitText:
				if update.Message.Text != "" {
					s.Text = update.Message.Text
					s.Stage = stageAwaitGIFOrSticker
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatID,
						Text:   "Text received. Now send a GIF or sticker.",
					})
					return
				}
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "Text expected",
				})
				resetSession(chatID)
				h.CreateHandler(ctx, b, update)
				return

			case stageAwaitGIFOrSticker:
				if update.Message.Animation != nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatID,
						Text:   "GIF received! Text: \"" + s.Text + "\". Processing...",
					})
					h.gifHandler(ctx, b, update, *s)
					resetSession(chatID)
					h.CreateHandler(ctx, b, update)
					return
				}

				if update.Message.Sticker != nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatID,
						Text:   "Sticker received! Text: \"" + s.Text + "\". Processing...",
					})

					h.stickerHandler(ctx, b, update, *s)
					resetSession(chatID)
					h.CreateHandler(ctx, b, update)
					return
				}

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "Gif or Sticker expected",
				})
				resetSession(chatID)
				h.CreateHandler(ctx, b, update)
				return
			default:
			}
		}

		next(ctx, b, update)
	}
}
