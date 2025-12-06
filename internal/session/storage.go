package session

type Storage interface {
	Get(chatID int64) (*Session, bool, error)
	Set(chatID int64, s *Session) error
	Delete(chatID int64) error
}

type inMemoryStorage struct {
	data map[int64]*Session
}

func NewInMemoryStorage() Storage {
	return &inMemoryStorage{
		data: make(map[int64]*Session),
	}
}

func (s *inMemoryStorage) Get(chatID int64) (*Session, bool, error) {
	sess, ok := s.data[chatID]
	return sess, ok, nil
}

func (s *inMemoryStorage) Set(chatID int64, sess *Session) error {
	s.data[chatID] = sess
	return nil
}

func (s *inMemoryStorage) Delete(chatID int64) error {
	delete(s.data, chatID)
	return nil
}
