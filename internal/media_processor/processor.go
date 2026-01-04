package media_processor

import (
	"io"

	"gifka-bot/internal/entity"
	"gifka-bot/internal/media_processor/media_factory"
)

type Processor interface {
	ProcessVideo(filePath, text string, t entity.TypeGif) (io.Reader, error)
	ProcessSticker(filePath, text string, t entity.TypeGif) (io.Reader, error)
}

type defaultProcessor struct {
	factory *media_factory.ProcessorFactory
}

func New() Processor {
	return &defaultProcessor{
		factory: media_factory.NewProcessorFactory(),
	}
}

func (p *defaultProcessor) ProcessVideo(filePath, text string, t entity.TypeGif) (io.Reader, error) {
	processor, err := p.factory.GetProcessor(filePath)
	if err != nil {
		return nil, err
	}
	return processor.Process(filePath, text, t)
}

func (p *defaultProcessor) ProcessSticker(filePath, text string, t entity.TypeGif) (io.Reader, error) {
	processor, err := p.factory.GetProcessor(filePath)
	if err != nil {
		return nil, err
	}
	return processor.Process(filePath, text, t)
}
