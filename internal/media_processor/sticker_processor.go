package media_processor

import (
	"io"
	"path/filepath"
)

func StickerProcessor(filePath string, text string) (processed io.Reader, err error) {
	extension := filepath.Ext(filePath)
	if extension == ".webp" {
		processed, err = WEBPProcessor(filePath, text)
	} else if extension == ".webm" {
		processed, err = WEBMProcessor(filePath, text)
	} else if extension == ".tgs" {

	}

	return processed, err

}
