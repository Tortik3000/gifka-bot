package media_processor

import (
	"image"
	"strings"

	"github.com/fogleman/gg"
)

func createWhiteText(img image.Image, outputPath, text string) error {
	lines := strings.Split(text, "\n")
	lineHeight := 48.0

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width%2 != 0 {
		width++
	}
	if height%2 != 0 {
		height++
	}

	dc := gg.NewContext(width, height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.DrawImage(img, 0, 0)

	dc.SetRGB(1, 1, 1)
	err := dc.LoadFontFace("./text_style/LiberationSerif-Regular.ttf", 48)
	if err != nil {
		return err
	}

	var startY float64
	if len(lines) <= 2 {
		startY = 20
	} else {
		startY = float64(height) - lineHeight*float64(len(lines)-1) - 20
	}

	for i, line := range lines {
		if i > 1 && len(lines) > 2 {
			dc.DrawStringAnchored(line, float64(width)/2, startY+float64(i)*lineHeight, 0.5, 0.5)
		} else {
			dc.DrawStringAnchored(line, float64(width)/2, startY+float64(i)*lineHeight, 0.5, 0.5)
		}
	}

	return dc.SavePNG(outputPath)
}
