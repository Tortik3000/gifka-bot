package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	welcomeText := `👋 Hello!

This bot can add captions to GIFs or stickers and returns them in the correct format for adding via the official @Stickers bot.

🎯 Supported formats:
• Static stickers (WebP)
• Animated stickers (WebM) 
• GIF images (MP4)
`

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   welcomeText,
	})

	h.CreateHandler(ctx, b, update)
}
