package media_processer

import (
	"bytes"
	"fmt"
	"gifka-bot/config"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-telegram/bot/models"
)

func VideoProcess(file *models.File, text string) (io.Reader, error) {
	cfg := config.New()
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", cfg.TG.Token, file.FilePath)
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	extension := filepath.Ext(file.FilePath)

	tempInput := "temp_input" + extension
	tempOutput := "temp_output" + extension
	framePNG := "frame.png"
	bgPNG := "background.png"

	// сохраняем webm-файл
	outFile, err := os.Create(tempInput)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempInput)

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return nil, err
	}
	outFile.Close()

	// извлекаем первый кадр
	cmdFrame := exec.Command("ffmpeg", "-y", "-i", tempInput, "-frames:v", "1", framePNG)
	if _, err := cmdFrame.CombinedOutput(); err != nil {
		return nil, err
	}
	defer os.Remove(framePNG)

	// создаем фон с текстом
	f, err := os.Open(framePNG)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(f)
	f.Close()
	if err != nil {
		return nil, err
	}

	if err := CreateBlackBox(img, bgPNG, text); err != nil {
		return nil, err
	}
	defer os.Remove(bgPNG)

	// накладываем фон на видео и кодируем обратно в webm
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-loop", "1", "-i", bgPNG,
		"-i", tempInput,
		"-filter_complex", "[0:v][1:v]overlay=50:50:shortest=1",
		"-c:v", "libvpx-vp9",
		"-b:v", "1M",
		"-pix_fmt", "yuv420p",
		"-an", "-shortest",
		tempOutput,
	)
	defer os.Remove(tempOutput)

	if _, err := cmd.CombinedOutput(); err != nil {
		return nil, err
	}

	// читаем готовое видео в память
	processedData, err := os.ReadFile(tempOutput)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(processedData), nil
}
