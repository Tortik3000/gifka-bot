package handlers

import "sync"

type convStage int

const (
	stageNone convStage = iota
	stageAwaitText
	stageAwaitGIFOrSticker
)

type TypeGif string

const (
	addText  TypeGif = "addText"
	blackBox TypeGif = "blackBox"
)

type session struct {
	Stage   convStage
	TypeGif TypeGif
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
