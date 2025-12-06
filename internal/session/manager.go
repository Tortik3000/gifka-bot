package session

import "sync"

type Manager struct {
	storage Storage
	mu      sync.RWMutex
}

func NewManager(storage Storage) *Manager {
	return &Manager{storage: storage}
}

func (m *Manager) Get(chatID int64) (*Session, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.storage.Get(chatID)
}

func (m *Manager) Set(chatID int64, s *Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.storage.Set(chatID, s)
}

func (m *Manager) Delete(chatID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.storage.Delete(chatID)
}

func (m *Manager) Reset(chatID int64) error {
	return m.Delete(chatID)
}
