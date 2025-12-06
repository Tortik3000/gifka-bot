package media_processor

import (
	"io"

	"gifka-bot/internal/entity"
)

func WEBMProcessor(filePath string, text string, typeGif entity.TypeGif) (io.Reader, error) {
	processed, err := VideoProcess(filePath, text, typeGif)

	return processed, err
}
