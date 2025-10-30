package media_processor

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
)

func VideoProcess(filePath string, text string) (io.Reader, error) {
	cfg := config.New()
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", cfg.TG.Token, filePath)
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	extension := filepath.Ext(filePath)

	tempInput := "temp_input" + extension
	tempOutput := "temp_output" + extension
	framePNG := "frame.png"
	bgPNG := "background.png"

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

	cmd := exec.Command(
		"ffmpeg", "-y",
		"-loop", "1", "-i", bgPNG,
		"-i", tempInput,
		"-filter_complex", "[0:v][1:v]overlay=50:50:shortest=1,scale=512:512:force_original_aspect_ratio=decrease:flags=lanczos", // Добавлен scale для ресайза
		"-c:v", "libvpx-vp9",
		"-b:v", "500K", // Битрейт снижен до 500K
		"-maxrate", "500K",
		"-bufsize", "1000K",
		"-pix_fmt", "yuva420p", // Формат с альфа-каналом для прозрачности
		"-an",
		"-t", "3", // Ограничение длительности 3 секунды
		"-quality", "good",
		"-crf", "37", // Увеличение CRF для большего сжатия
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
