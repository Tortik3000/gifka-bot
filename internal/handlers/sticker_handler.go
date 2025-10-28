package handlers

import (
	"context"
	"io"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	mp "gifka-bot/internal/media_processer"
)

func StickerHandler(ctx context.Context, b *bot.Bot, update *models.Update, s session) {
	chatID := update.Message.Chat.ID
	sr := update.Message.Sticker
	fileID := sr.FileID
	text := s.Text

	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: fileID})
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error receiving file."})
		return
	}

	var processed io.Reader
	var fileName string
	if !sr.IsVideo {
		fileName = "sticker.webp"
		processed, err = mp.WEBPProcessor(file, text)
	} else {
		fileName = "sticker.webm"
		processed, err = mp.VideoProcess(file, text)
	}

	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error processing sticker."})
		return
	}

	sticker := &models.InputFileUpload{
		Filename: fileName,
		Data:     processed,
	}

	_, err = b.SendSticker(ctx, &bot.SendStickerParams{
		ChatID:  chatID,
		Sticker: sticker,
	})

	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error sending sticker."})
	}
}
