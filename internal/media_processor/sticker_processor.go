package media_processor

import (
	"gifka-bot/internal/entity"
	"io"
	"path/filepath"
)

func StickerProcessor(filePath string, text string, typeGif entity.TypeGif) (processed io.Reader, err error) {
	extension := filepath.Ext(filePath)
	if extension == ".webp" {
		processed, err = WEBPProcessor(filePath, text, typeGif)
	} else if extension == ".webm" {
		processed, err = WEBMProcessor(filePath, text, typeGif)
	} else if extension == ".tgs" {

	}

	return processed, err

}
