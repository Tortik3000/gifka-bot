package media_factory

import (
	"io"

	"gifka-bot/internal/entity"
)

type WebmProcessor struct{}

func (w *WebmProcessor) Process(filePath string, text string, t entity.TypeGif) (io.Reader, error) {
	return WEBMProcessor(filePath, text, t)
}

func WEBMProcessor(filePath string, text string, typeGif entity.TypeGif) (io.Reader, error) {
	processed, err := VideoProcess(filePath, text, typeGif)

	return processed, err
}
