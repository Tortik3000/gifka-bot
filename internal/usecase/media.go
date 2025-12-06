package usecase

import (
	"io"

	"go.uber.org/zap"

	"gifka-bot/internal/entity"
)

type MediaService struct {
	processor mediaProcessor
	logger    *zap.Logger
}

func NewMediaService(p mediaProcessor, logger *zap.Logger) *MediaService {
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
