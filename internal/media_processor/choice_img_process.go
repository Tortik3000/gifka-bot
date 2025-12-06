package media_processor

import (
	"gifka-bot/internal/entity"
	"image"
)

func choiceImgProcess(img image.Image, outputPath, text string, typeGif entity.TypeGif) error {
	switch typeGif {
	case entity.AddText:
		return createWhiteText(img, outputPath, text)
	case entity.BlackBox:
		return createBlackBox(img, outputPath, text)
	}

	return nil
}
