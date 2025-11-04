package media_processor

import (
	"gifka-bot/internal/entity"
	"io"
)

func WEBMProcessor(filePath string, text string, typeGif entity.TypeGif) (io.Reader, error) {
	processed, err := VideoProcess(filePath, text, typeGif)

	return processed, err
}
