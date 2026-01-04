package media_factory

import (
	"io"

	"gifka-bot/internal/entity"
)

type WebmProcessor struct{}

func (w *WebmProcessor) Process(filePath string, text string, t entity.TypeGif) (io.Reader, error) {
	videoProcessor := VideoProcessor{}
	return videoProcessor.Process(filePath, text, t)
}
