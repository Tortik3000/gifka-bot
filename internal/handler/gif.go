// internal/handler/gif.go
package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handler) GIF(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Animation == nil {
		return
	}

	chatID := update.Message.Chat.ID
	fileID := update.Message.Animation.FileID

	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: fileID})
	if err != nil {
		h.logger.Error("get file", zap.Error(err))
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error receiving file."})
		return
	}

	sess, _, err := h.sessionSvc.Manager.Get(chatID)
	if err != nil || sess == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Session not found, please start over."})
		return
	}

	processed, err := h.mediaSvc.ProcessGIF(file.FilePath, sess.Text, sess.TypeGif)
	if err != nil {
		h.logger.Error("process gif", zap.Error(err))
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "GIF processing error."})
		return
	}

	video := &models.InputFileUpload{
		Filename: "animation.mp4",
		Data:     processed,
	}

	_, err = b.SendAnimation(ctx, &bot.SendAnimationParams{
		ChatID:    chatID,
		Animation: video,
	})
	if err != nil {
		h.logger.Error("send animation", zap.Error(err))
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error sending GIF."})
		return
	}

	h.convSvc.Finish(chatID)
	h.Create(ctx, b, update)
}
