package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func (h *Handler) CreateHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := inline.New(b).
		Row().
		Button("Add Black Box", []byte(blackBox), h.AddText)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Choose source:",
		ReplyMarkup: kb,
	})
}

func (h *Handler) AddText(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	chatID := mes.Message.Chat.ID

	s := getSession(chatID)
	s.Stage = stageAwaitText
	s.Text = ""
	s.TypeGif = TypeGif(data)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Send the text you want to add",
	})
}
