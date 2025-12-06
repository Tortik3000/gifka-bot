// internal/middleware/conversation.go
package middleware

import (
	"context"

	"gifka-bot/internal/session"
	"gifka-bot/internal/usecase"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ConversationMiddleware struct {
	sessions *session.Manager
	convSvc  *usecase.ConversationService

	h interface {
		GIF(ctx context.Context, b *bot.Bot, update *models.Update)
		Sticker(ctx context.Context, b *bot.Bot, update *models.Update)
		Create(ctx context.Context, b *bot.Bot, update *models.Update)
	}
}

func NewConversation(
	m *session.Manager,
	convSvc *usecase.ConversationService,
	h interface {
		GIF(ctx context.Context, b *bot.Bot, update *models.Update)
		Sticker(ctx context.Context, b *bot.Bot, update *models.Update)
		Create(ctx context.Context, b *bot.Bot, update *models.Update)
	},
) *ConversationMiddleware {
	return &ConversationMiddleware{
		sessions: m,
		convSvc:  convSvc,
		h:        h,
	}
}

func (m *ConversationMiddleware) Handle(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		// если нет Message — пропускаем
		if update == nil || update.Message == nil {
			next(ctx, b, update)
			return
		}

		chatID := update.Message.Chat.ID

		// 1. если ожидаем текст
		if update.Message.Text != "" {
			sess, ok, _ := m.sessions.Get(chatID)
			if ok && sess.Stage == session.StageAwaitText {
				_, msg, err := m.convSvc.HandleText(chatID, update.Message.Text)
				if err == nil && msg != "" {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatID,
						Text:   msg,
					})
					return
				}
			}
			next(ctx, b, update)
			return
		}

		// 2. если ожидаем GIF/стикер
		sess, expect, _ := m.convSvc.ExpectMedia(chatID)
		if expect && sess != nil {
			if update.Message.Animation != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "GIF получен! Обработка...",
				})
				m.h.GIF(ctx, b, update)
				return
			}
			if update.Message.Sticker != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "Стикер получен! Обработка...",
				})
				m.h.Sticker(ctx, b, update)
				return
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Ожидался GIF или стикер.",
			})
			m.convSvc.Finish(chatID)
			m.h.Create(ctx, b, update)
			return
		}

		// по умолчанию
		next(ctx, b, update)
	}
}
