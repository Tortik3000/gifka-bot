package handler

import (
	"context"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) Sticker(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Sticker == nil {
		return
	}

	chatID := update.Message.Chat.ID
	sticker := update.Message.Sticker
	fileID := sticker.FileID

	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: fileID})
	if err != nil {
		h.logger.Error("get file", zap.Error(err))
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error receiving file."})
		return
	}

	session, err := h.sessionUseCase.Get(chatID)
	if err != nil || session == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Session not found, please start over."})
		return
	}

	processed, err := h.mediaUseCase.ProcessSticker(file.FilePath, session.Text, session.TypeGif)
	if err != nil {
		h.logger.Error("process sticker", zap.Error(err))
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error processing sticker."})
		return
	}

	ext := filepath.Ext(file.FilePath)
	returnSticker := &models.InputFileUpload{
		Filename: "returnSticker" + ext,
		Data:     processed,
	}

	_, err = b.SendSticker(ctx, &bot.SendStickerParams{
		ChatID:  chatID,
		Sticker: returnSticker,
	})
	if err != nil {
		h.logger.Error("send sticker", zap.Error(err))
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error sending sticker."})
		return
	}

	h.convUseCase.Finish(chatID)
	h.Default(ctx, b, update)
}
