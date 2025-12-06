package media_processor

import (
	"errors"
	"image"

	"gifka-bot/internal/entity"
)

func choiceImgProcess(img image.Image, outputPath, text string, typeGif entity.TypeGif) error {
	switch typeGif {
	case entity.AddText:
		return errors.New("not supported")
	case entity.BlackBox:
		return createBlackBox(img, outputPath, text)
	}

	return nil
}
