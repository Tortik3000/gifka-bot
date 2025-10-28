package media_processer

import (
	"bytes"
	"fmt"
	"gifka-bot/config"
	"image"
	"io"
	"net/http"
	"os"

	"github.com/go-telegram/bot/models"
	"golang.org/x/image/webp"
)

func WEBPProcessor(file *models.File, text string) (io.Reader, error) {
	img, err := getImgFromFile(file)

	outputPath := "background.png"
	err = CreateBlackBox(img, outputPath, text)
	if err != nil {
		return nil, err
	}

	processedData, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(processedData), nil
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
