package media_processor

import "io"

func WEBMProcessor(filePath string, text string) (io.Reader, error) {
	processed, err := VideoProcess(filePath, text)

	return processed, err
}
