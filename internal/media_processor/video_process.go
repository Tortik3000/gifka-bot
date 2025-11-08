package media_processor

import (
	"bytes"
	"fmt"
	"gifka-bot/config"
	"gifka-bot/internal/entity"

	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func VideoProcess(filePath string, text string, typeGif entity.TypeGif) (io.Reader, error) {
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

	if err := choiceImgProcess(img, bgPNG, text, typeGif); err != nil {
		return nil, err
	}
	defer os.Remove(bgPNG)

	cmd := exec.Command(
		"ffmpeg", "-y",
		// Короткий "фейковый" контейнер с первым кадром
		"-i", tempInput,
		"-t", "0.03", // Duration контейнера
		"-c:v", "libvpx-vp9",
		"-pix_fmt", "yuva420p",
		"-metadata", "title=TelegramSticker",
		"fake_container.webm",
	)
	if _, err := cmd.CombinedOutput(); err != nil {
		return nil, err
	}
	defer os.Remove("fake_container.webm")

	// Накладываем реальное видео поверх короткого контейнера
	cmdOverlay := exec.Command(
		"ffmpeg", "-y",
		"-loop", "1", "-i", bgPNG,
		"-i", "fake_container.webm",
		"-filter_complex",
		"[1:v]scale=512:512:force_original_aspect_ratio=decrease[vid];"+
			"[0:v][vid]overlay=(W-w)/2:(H-h)/2:shortest=1",
		"-c:v", "libvpx-vp9",
		"-pix_fmt", "yuva420p",
		"-b:v", "300K",
		"-crf", "32",
		"-deadline", "good",
		"-an",
		tempOutput,
	)
	defer os.Remove(tempOutput)

	if _, err := cmdOverlay.CombinedOutput(); err != nil {
		return nil, err
	}

	// читаем готовое видео в память
	processedData, err := os.ReadFile(tempOutput)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(processedData), nil
}
