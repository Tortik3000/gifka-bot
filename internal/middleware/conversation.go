// internal/middleware/conversation.go
package middleware

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"gifka-bot/internal/session"
)

var process = []string{"process.", "process..", "process..."}

type ConversationMiddleware struct {
	sessions    sessionUseCase
	convUseCase conversationUseCase
	handler     handler
}

func NewConversation(
	m sessionUseCase,
	convUseCase conversationUseCase,
	h handler,
) *ConversationMiddleware {
	return &ConversationMiddleware{
		sessions:    m,
		convUseCase: convUseCase,
		handler:     h,
	}
}

func (m *ConversationMiddleware) Handle(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update == nil || update.Message == nil {
			next(ctx, b, update)
			return
		}

		chatID := update.Message.Chat.ID

		if update.Message.Text != "" {
			sess, ok, _ := m.sessions.Get(chatID)
			if ok && sess.Stage == session.StageAwaitText {
				_, msg, err := m.convUseCase.HandleText(chatID, update.Message.Text)
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

		sess, expect, _ := m.convUseCase.ExpectMedia(chatID)
		if expect && sess != nil {
			switch {
			case update.Message.Animation != nil:
				m.handleMediaWithProgress(
					ctx,
					b,
					update,
					chatID,
					"GIF received! Processing...",
					m.handler.GIF,
				)
				return

			case update.Message.Sticker != nil:
				m.handleMediaWithProgress(
					ctx,
					b,
					update,
					chatID,
					"Sticker received! Processing...",
					m.handler.Sticker,
				)
				return
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Expected a GIF or sticker.",
			})
			m.convUseCase.Finish(chatID)
			m.handler.Default(ctx, b, update)
			return
		}

		next(ctx, b, update)
	}
}

func (m *ConversationMiddleware) handleMediaWithProgress(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	chatID int64,
	initialText string,
	fn func(context.Context, *bot.Bot, *models.Update),
) {
	msgID, cancel := m.startProgress(ctx, b, chatID, initialText)
	defer m.stopProgress(ctx, b, chatID, msgID, cancel)

	fn(ctx, b, update)
}

func (m *ConversationMiddleware) startProgress(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	initialText string,
) (messageID int, cancel context.CancelFunc) {
	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   initialText,
	})
	if err != nil {

		return 0, func() {}
	}

	pctx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context, chatID int64, msgID int) {
		i := 0
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, _ = b.EditMessageText(ctx, &bot.EditMessageTextParams{
					ChatID:    chatID,
					MessageID: msgID,
					Text:      process[i%len(process)],
				})
				i++
			case <-ctx.Done():
				return
			}
		}
	}(pctx, chatID, msg.ID)

	return msg.ID, cancel
}

func (m *ConversationMiddleware) stopProgress(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	messageID int,
	cancel context.CancelFunc,
) {
	if cancel != nil {
		cancel()
	}
	if messageID == 0 {
		return
	}

	_, _ = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatID,
		MessageID: messageID,
	})
}
