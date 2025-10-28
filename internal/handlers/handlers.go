package handlers

import (
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func New(
	logger *zap.Logger,
) *Handler {
	return &Handler{
		logger: logger,
	}
}
