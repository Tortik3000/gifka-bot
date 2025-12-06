// internal/handler/create.go
package handler

import (
	"context"

	"gifka-bot/internal/entity"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func (h *Handler) Create(ctx context.Context, b *bot.Bot, update *models.Update) {
	var chatID int64
	switch {
	case update.Message != nil:
		chatID = update.Message.Chat.ID
	case update.CallbackQuery != nil && update.CallbackQuery.Message.InaccessibleMessage != nil:
		chatID = update.CallbackQuery.Message.InaccessibleMessage.Chat.ID
	default:
		return
	}

	kb := inline.New(b).
		Row().
		Button("Add Black Box", []byte(entity.BlackBox), h.AddTextCallback).
		Row().
		Button("Add White Text(don't work)", []byte(entity.AddText), h.AddTextCallback)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Choose source:",
		ReplyMarkup: kb,
	})
}

// callback по нажатию inline‑кнопки
func (h *Handler) AddTextCallback(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	chatID := mes.Message.Chat.ID
	t := entity.TypeGif(data)

	// запускаем сценарий
	_ = h.convUseCase.StartAddText(chatID, t)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Submit the text you want to add.",
	})
}
