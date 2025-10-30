package media_processor

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
)

func WEBPProcessor(filePath string, text string) (io.Reader, error) {
	img, err := getImgFromWEBPFile(filePath)
	if err != nil {
		return nil, err
	}

	outPngPath := "out.png"
	err = CreateBlackBox(img, outPngPath, text)
	if err != nil {
		return nil, err
	}

	outPngFile, err := os.Open(outPngPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PNG file: %w", err)
	}
	defer outPngFile.Close()
	defer os.Remove(outPngPath)

	pngImg, err := png.Decode(outPngFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}
	pngImg, _ = resizeToExactStickerSize(pngImg)

	var processedData bytes.Buffer
	options := &webp.Options{
		Lossless: false,
		Quality:  85,
		Exact:    true,
	}
	err = webp.Encode(&processedData, pngImg, options)
	if err != nil {
		return nil, err
	}

	return &processedData, nil
}

func getImgFromWEBPFile(filePath string) (image.Image, error) {
	cfg := config.New()
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", cfg.TG.Token, filePath)
	resp, err := http.Get(fileURL)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return webp.Decode(resp.Body)
}

func resizeToExactStickerSize(img image.Image) (image.Image, error) {
	const size = 512

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	var resized image.Image
	if w > h {
		resized = imaging.Resize(img, size, 0, imaging.Lanczos)
	} else {
		resized = imaging.Resize(img, 0, size, imaging.Lanczos)
	}

	// Создаем квадратное изображение 512x512 с прозрачным фоном
	dst := image.NewNRGBA(image.Rect(0, 0, size, size))

	// Центрируем изображение
	resizedBounds := resized.Bounds()
	offsetX := (size - resizedBounds.Dx()) / 2
	offsetY := (size - resizedBounds.Dy()) / 2

	// Копируем ресайзнутое изображение в центр
	for y := 0; y < resizedBounds.Dy(); y++ {
		for x := 0; x < resizedBounds.Dx(); x++ {
			dst.Set(x+offsetX, y+offsetY, resized.At(x, y))
		}
	}

	return dst, nil
}
