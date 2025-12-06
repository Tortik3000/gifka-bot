// internal/service/media.go
package usecase

import (
	"io"

	"gifka-bot/internal/entity"
	"gifka-bot/internal/media_processor"

	"go.uber.org/zap"
)

type MediaService struct {
	processor media_processor.Processor
	logger    *zap.Logger
}

func NewMediaService(p media_processor.Processor, logger *zap.Logger) *MediaService {
	return &MediaService{
		processor: p,
		logger:    logger,
	}
}

func (s *MediaService) ProcessGIF(filePath, text string, t entity.TypeGif) (io.Reader, error) {
	return s.processor.ProcessVideo(filePath, text, t)
}

func (s *MediaService) ProcessSticker(filePath, text string, t entity.TypeGif) (io.Reader, error) {
	return s.processor.ProcessSticker(filePath, text, t)
}
