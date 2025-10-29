package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	mp "gifka-bot/internal/media_processer"
)

func (h *Handler) gifHandler(ctx context.Context, b *bot.Bot, update *models.Update, s session) {
	chatID := update.Message.Chat.ID
	fileID := update.Message.Animation.FileID
	text := s.Text

	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: fileID})
	if err != nil {
		h.logger.Error(err.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error receiving file."})
		return
	}

	processed, err := mp.VideoProcess(file, text)
	if err != nil {
		h.logger.Error(err.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error processing video."})
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
		h.logger.Error(err.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error sending animation."})
		return
	}
}
