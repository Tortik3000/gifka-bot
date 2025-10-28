package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func ConversationMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
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
				return

			case stageAwaitGIFOrSticker:
				if update.Message.Animation != nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatID,
						Text:   "GIF received! Text: \"" + s.Text + "\". Processing...",
					})
					GifHandler(ctx, b, update, *s)
					resetSession(chatID)
					return
				}

				if update.Message.Sticker != nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatID,
						Text:   "Sticker received! Text: \"" + s.Text + "\". Processing...",
					})

					StickerHandler(ctx, b, update, *s)
					resetSession(chatID)
					return
				}

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "Gif or Sticker expected",
				})
				resetSession(chatID)
				return
			default:
			}
		}

		next(ctx, b, update)
	}
}
