package usecase

import (
	"io"

	"gifka-bot/internal/entity"
	"gifka-bot/internal/session"
)

type (
	sessionManager interface {
		Get(chatID int64) (*session.Session, bool, error)
		Set(chatID int64, s *session.Session) error
		Delete(chatID int64) error
		Reset(chatID int64) error
	}

	mediaProcessor interface {
		ProcessVideo(filePath, text string, t entity.TypeGif) (io.Reader, error)
		ProcessSticker(filePath, text string, t entity.TypeGif) (io.Reader, error)
	}
)
