package usecase

import "gifka-bot/internal/session"

type SessionUseCase struct {
	Manager sessionManager
}

func NewSessionService(m sessionManager) *SessionUseCase {
	return &SessionUseCase{Manager: m}
}

func (s *SessionUseCase) GetOrCreate(chatID int64) (*session.Session, error) {
	sess, ok, err := s.Manager.Get(chatID)
	if err != nil {
		return nil, err
	}
	if !ok {
		sess = &session.Session{Stage: session.StageNone}
		if err := s.Manager.Set(chatID, sess); err != nil {
			return nil, err
		}
	}
	return sess, nil
}

func (s *SessionUseCase) Get(chatID int64) (*session.Session, error) {
	sess, ok, err := s.Manager.Get(chatID)
	if err != nil || !ok {
		return nil, err
	}

	return sess, nil
}
