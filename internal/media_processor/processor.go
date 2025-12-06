package media_processor

import (
	"io"

	"gifka-bot/internal/entity"
)

// Здесь вы можете просто обернуть вашу существующую логику.
// В примере предполагаются функции VideoProcess и StickerProcessor, как у вас.

type Processor interface {
	ProcessVideo(filePath, text string, t entity.TypeGif) (io.Reader, error)
	ProcessSticker(filePath, text string, t entity.TypeGif) (io.Reader, error)
}

type defaultProcessor struct{}

func New() Processor {
	return &defaultProcessor{}
}

func (p *defaultProcessor) ProcessVideo(filePath, text string, t entity.TypeGif) (io.Reader, error) {
	return VideoProcess(filePath, text, t)
}

func (p *defaultProcessor) ProcessSticker(filePath, text string, t entity.TypeGif) (io.Reader, error) {
	return StickerProcessor(filePath, text, t)
}
