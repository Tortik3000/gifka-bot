// internal/session/session.go
package session

import "gifka-bot/internal/entity"

type Stage int

const (
	StageNone Stage = iota
	StageAwaitText
	StageAwaitGIFOrSticker
)

type Session struct {
	Stage   Stage
	TypeGif entity.TypeGif
	Text    string
}
