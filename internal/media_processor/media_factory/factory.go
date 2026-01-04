package media_factory

import (
	"io"
	"path/filepath"

	"gifka-bot/internal/entity"
)

type MediaProcessor interface {
	Process(filePath, text string, t entity.TypeGif) (io.Reader, error)
}

type ProcessorFactory struct{}

func NewProcessorFactory() *ProcessorFactory {
	return &ProcessorFactory{}
}

func (f *ProcessorFactory) GetProcessor(filePath string) (MediaProcessor, error) {
	extension := filepath.Ext(filePath)
	switch extension {
	case ".webp":
		return &WebpProcessor{}, nil
	case ".webm":
		return &WebmProcessor{}, nil
	case ".mp4", ".mov", ".avi":
		return &VideoProcessor{}, nil
	default:
		return &VideoProcessor{}, nil
	}
}
