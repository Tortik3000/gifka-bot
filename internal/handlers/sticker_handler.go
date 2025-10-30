package handlers

import (
	"context"
	"io"
	"path/filepath"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	mp "gifka-bot/internal/media_processor"
)

func (h *Handler) stickerHandler(ctx context.Context, b *bot.Bot, update *models.Update, s session) {
	chatID := update.Message.Chat.ID
	sticker := update.Message.Sticker
	fileID := sticker.FileID
	text := s.Text

	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: fileID})
	if err != nil {
		h.logger.Error(err.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error receiving file."})
		return
	}

	var processed io.Reader
	processed, err = mp.StickerProcessor(file.FilePath, text)

	if err != nil {
		h.logger.Error(err.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error processing returnSticker."})
		return
	}

	extension := filepath.Ext(file.FilePath)
	returnSticker := &models.InputFileUpload{
		Filename: "returnSticker" + extension,
		Data:     processed,
	}

	_, err = b.SendSticker(ctx, &bot.SendStickerParams{
		ChatID:  chatID,
		Sticker: returnSticker,
	})
	if err != nil {
		h.logger.Error(err.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error sending returnSticker."})
	}
}
