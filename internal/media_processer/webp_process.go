package media_processer

import (
	"bytes"
	"fmt"
	"gifka-bot/config"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"

	"github.com/go-telegram/bot/models"
)

func WEBPProcessor(file *models.File, text string) (io.Reader, error) {
	img, err := getImgFromFile(file)
	if err != nil {
		return nil, err
	}

	outputPath := "background.png"
	err = CreateBlackBox(img, outputPath, text)
	if err != nil {
		return nil, err
	}

	pngFile, err := os.Open(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PNG file: %w", err)
	}
	defer pngFile.Close()
	defer os.Remove(outputPath)

	pngImg, err := png.Decode(pngFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}
	pngImg = resizeToStickerSize(pngImg)

	var processedData bytes.Buffer
	options := &webp.Options{Lossless: true, Quality: 100}
	err = webp.Encode(&processedData, pngImg, options)
	if err != nil {
		return nil, err
	}

	return &processedData, nil
}

func getImgFromFile(file *models.File) (image.Image, error) {
	cfg := config.New()
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", cfg.TG.Token, file.FilePath)
	resp, err := http.Get(fileURL)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return webp.Decode(resp.Body)
}

func resizeToStickerSize(img image.Image) image.Image {
	const maxSize = 512
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	if w > h {
		return imaging.Resize(img, maxSize, 0, imaging.Lanczos)
	}
	return imaging.Resize(img, 0, maxSize, imaging.Lanczos)
}
