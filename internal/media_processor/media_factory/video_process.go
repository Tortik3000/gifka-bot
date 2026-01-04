package media_factory

import (
	"bytes"
	"fmt"

	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"gifka-bot/config"
	"gifka-bot/internal/entity"
	"gifka-bot/internal/media_processor/image_factory"
)

type VideoProcessor struct{}

func (v *VideoProcessor) Process(filePath string, text string, t entity.TypeGif) (io.Reader, error) {
	return VideoProcess(filePath, text, t)
}

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

	factory := image_factory.NewImageProcessorFactory()
	processor, err := factory.GetProcessor(typeGif)
	if err != nil {
		return nil, err
	}

	if err := processor.Process(img, bgPNG, text); err != nil {
		return nil, err
	}
	defer os.Remove(bgPNG)

	cmd := exec.Command(
		"ffmpeg", "-y",
		"-loop", "1", "-i", bgPNG,
		"-i", tempInput,
		"-filter_complex", "[0:v][1:v]overlay=50:50:shortest=1,scale=512:512:force_original_aspect_ratio=decrease:flags=lanczos", // Добавлен scale для ресайза
		"-c:v", "libvpx-vp9",
		"-b:v", "500K",
		"-maxrate", "500K",
		"-bufsize", "1000K",
		"-pix_fmt", "yuva420p",
		"-an",
		"-quality", "good",
		"-crf", "37",
		tempOutput,
	)
	defer os.Remove(tempOutput)

	if _, err := cmd.CombinedOutput(); err != nil {
		return nil, err
	}

	processedData, err := os.ReadFile(tempOutput)
	if err != nil {
		return nil, err
	}

	// Патчим Duration на 30ms (или другое значение)
	patchedReader, err := PatchWebMDurationReader(processedData, 30.0)
	if err != nil {
		// Если Duration не найден, возвращаем оригинальные данные
		return bytes.NewReader(processedData), nil
	}

	return patchedReader, nil
}
