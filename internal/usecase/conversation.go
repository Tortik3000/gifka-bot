package usecase

import (
	"gifka-bot/internal/entity"
	"gifka-bot/internal/session"

	"go.uber.org/zap"
)

type ConversationService struct {
	sessions *session.Manager
	logger   *zap.Logger
}

func NewConversationService(m *session.Manager, logger *zap.Logger) *ConversationService {
	return &ConversationService{
		sessions: m,
		logger:   logger,
	}
}

// Начать сценарий "добавить текст"
func (s *ConversationService) StartAddText(chatID int64, t entity.TypeGif) error {
	sess := &session.Session{
		Stage:   session.StageAwaitText,
		TypeGif: t,
		Text:    "",
	}
	return s.sessions.Set(chatID, sess)
}

// Обработка сообщения с текстом
func (s *ConversationService) HandleText(chatID int64, text string) (*session.Session, string, error) {
	sess, ok, err := s.sessions.Get(chatID)
	if err != nil {
		return nil, "", err
	}
	if !ok || sess.Stage != session.StageAwaitText {
		return nil, "Неожиданный текст. Нажмите кнопку ещё раз.", nil
	}

	if text == "" {
		// сбрасываем сессию
		_ = s.sessions.Reset(chatID)
		return nil, "Ожидался текст, попробуйте ещё раз.", nil
	}

	sess.Text = text
	sess.Stage = session.StageAwaitGIFOrSticker
	if err := s.sessions.Set(chatID, sess); err != nil {
		return nil, "", err
	}

	return sess, "Текст получен. Теперь пришлите GIF или стикер.", nil
}

// Проверка, что сообщение с GIF/стикером ожидается
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
