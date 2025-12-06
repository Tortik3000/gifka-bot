package usecase

import "gifka-bot/internal/session"

type SessionService struct {
	Manager *session.Manager
}

func NewSessionService(m *session.Manager) *SessionService {
	return &SessionService{Manager: m}
}

func (s *SessionService) GetOrCreate(chatID int64) (*session.Session, error) {
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
