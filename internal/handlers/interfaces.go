package handlers

import (
	"gifka-bot/internal/entity"
	"sync"
)

type convStage int

const (
	stageNone convStage = iota
	stageAwaitText
	stageAwaitGIFOrSticker
)

type session struct {
	Stage   convStage
	TypeGif entity.TypeGif
	Text    string
}

var (
	sessionsMu sync.Mutex
	sessions   = make(map[int64]*session)
)

func getSession(chatID int64) *session {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	s, ok := sessions[chatID]
	if !ok {
		s = &session{Stage: stageNone}
		sessions[chatID] = s
	}
	return s
}

func resetSession(chatID int64) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	delete(sessions, chatID)
}
