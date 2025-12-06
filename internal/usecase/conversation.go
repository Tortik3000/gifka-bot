package usecase

import (
	"go.uber.org/zap"

	"gifka-bot/internal/entity"
	"gifka-bot/internal/session"
)

type ConversationService struct {
	sessions sessionManager
	logger   *zap.Logger
}

func NewConversationService(m sessionManager, logger *zap.Logger) *ConversationService {
	return &ConversationService{
		sessions: m,
		logger:   logger,
	}
}

func (s *ConversationService) StartAddText(chatID int64, t entity.TypeGif) error {
	sess := &session.Session{
		Stage:   session.StageAwaitText,
		TypeGif: t,
		Text:    "",
	}
	return s.sessions.Set(chatID, sess)
}

func (s *ConversationService) HandleText(chatID int64, text string) (*session.Session, string, error) {
	sess, ok, err := s.sessions.Get(chatID)
	if err != nil {
		return nil, "", err
	}
	if !ok || sess.Stage != session.StageAwaitText {
		return nil, "Unexpected text. Press the button again.", nil
	}

	if text == "" {
		_ = s.sessions.Reset(chatID)
		return nil, "Expected text, please try again.", nil
	}

	sess.Text = text
	sess.Stage = session.StageAwaitGIFOrSticker
	if err := s.sessions.Set(chatID, sess); err != nil {
		return nil, "", err
	}

	return sess, "Text received. Now send a GIF or sticker.", nil
}

func (s *ConversationService) ExpectMedia(chatID int64) (*session.Session, bool, error) {
	sess, ok, err := s.sessions.Get(chatID)
	if err != nil || !ok {
		return nil, false, err
	}
	if sess.Stage != session.StageAwaitGIFOrSticker {
		return sess, false, nil
	}
	return sess, true, nil
}

func (s *ConversationService) Finish(chatID int64) {
	_ = s.sessions.Reset(chatID)
}
