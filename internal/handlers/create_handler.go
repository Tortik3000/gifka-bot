package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func (h *Handler) CreateHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var chatID int64

	switch {
	case update.Message != nil:
		chatID = update.Message.Chat.ID
	case update.CallbackQuery != nil && update.CallbackQuery.Message.InaccessibleMessage != nil:
		chatID = update.CallbackQuery.Message.InaccessibleMessage.Chat.ID
	default:
	}

	kb := inline.New(b).
		Row().
		Button("Add Black Box", []byte(blackBox), h.addText)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Choose source:",
		ReplyMarkup: kb,
	})
	if err != nil {
		return
	}
}

func (h *Handler) addText(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
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
