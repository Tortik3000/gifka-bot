package media_processer

import (
	"image"

	"github.com/fogleman/gg"
)

func CreateBlackBox(img image.Image, outputPath, text string) error {
	width := img.Bounds().Dx() + 100
	height := img.Bounds().Dy() + 100
	if width%2 != 0 {
		width++
	}
	if height%2 != 0 {
		height++
	}

	dc := gg.NewContext(width, height)

	// Чёрный фон
	dc.SetRGB(0, 0, 0)
	dc.Clear()

	// Белая рамка
	dc.SetRGB(1, 1, 1)
	dc.SetLineWidth(4)
	dc.DrawRectangle(40, 40, float64(width-80), float64(img.Bounds().Dy())+10)
	dc.Stroke()

	// Само изображение
	dc.DrawImage(img, 50, 50)

	// Белый текст
	dc.SetRGB(1, 1, 1)
	if err := dc.LoadFontFace("./text_style/LiberationSerif-Regular.ttf", 28); err == nil {
		dc.DrawStringAnchored(text, float64(width)/2, float64(height)-30, 0.5, 0.5)
	}

	return dc.SavePNG(outputPath)
}
