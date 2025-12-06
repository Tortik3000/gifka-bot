package middleware

import (
	"context"
	"gifka-bot/internal/entity"
	"gifka-bot/internal/session"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type (
	handler interface {
		GIF(ctx context.Context, b *bot.Bot, update *models.Update)
		Sticker(ctx context.Context, b *bot.Bot, update *models.Update)
		Default(ctx context.Context, b *bot.Bot, update *models.Update)
	}

	conversationUseCase interface {
		StartAddText(chatID int64, t entity.TypeGif) error
		HandleText(chatID int64, text string) (*session.Session, string, error)
		ExpectMedia(chatID int64) (*session.Session, bool, error)
		Finish(chatID int64)
	}

	sessionUseCase interface {
		Get(chatID int64) (*session.Session, bool, error)
	}
)
