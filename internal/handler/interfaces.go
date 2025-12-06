package handler

import (
	"io"

	"gifka-bot/internal/entity"
	"gifka-bot/internal/session"
)

type (
	mediaUseCase interface {
		ProcessGIF(filePath, text string, t entity.TypeGif) (io.Reader, error)
		ProcessSticker(filePath, text string, t entity.TypeGif) (io.Reader, error)
	}

	conversationUseCase interface {
		StartAddText(chatID int64, t entity.TypeGif) error
		HandleText(chatID int64, text string) (*session.Session, string, error)
		ExpectMedia(chatID int64) (*session.Session, bool, error)
		Finish(chatID int64)
	}

	sessionUseCase interface {
		GetOrCreate(chatID int64) (*session.Session, error)
		Get(chatID int64) (*session.Session, error)
	}
)
