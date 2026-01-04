package image_factory

import (
	"errors"
	"image"

	"gifka-bot/internal/entity"
)

type ImageProcessor interface {
	Process(img image.Image, outputPath, text string) error
}

type ImageProcessorFactory struct{}

func NewImageProcessorFactory() *ImageProcessorFactory {
	return &ImageProcessorFactory{}
}

func (f *ImageProcessorFactory) GetProcessor(typeGif entity.TypeGif) (ImageProcessor, error) {
	switch typeGif {
	case entity.AddText:
		return nil, errors.New("not supported")
	case entity.BlackBox:
		return &blackBoxProcessor{}, nil
	default:
		return nil, errors.New("unknown gif type")
	}
}
