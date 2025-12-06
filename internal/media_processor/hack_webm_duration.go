package media_processor

import (
	"os"

	"github.com/ebml-go/webm"
)

func HackWebMDuration(inPath, outPath string) error {
	f, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// m будет заполнен структурой WebM (Segment, Info, Tracks, и т.п.)
	var m webm.WebM

	// Parse: читает контейнер, заполняет m и возвращает Reader
	// сигнатура из документации: func Parse(r io.ReadSeeker, m *WebM) (wr *Reader, err error)
	r, err := webm.Parse(f, &m)
	if err != nil {
		return err
	}
	_ = r // Reader можно использовать, если нужно проходить по фреймам

	// 1. Укорачиваем контейнерную длительность
	// смотри в GoDoc, как именно хранится Duration (обычно float64/float32 в Info)

	m.Segment.Duration = 0.03 // 30 ms – то, что ffmpeg покажет как Duration в Input

	// 2. Вешаем длинный DURATION в метаданные видео‑трека
	// ниже – общая идея, реальные поля зависят от структуры WebM

	// 3. Сериализация обратно
	// В этом месте нужно использовать реальный writer из ebml-go/webm,
	// если он есть (например, webm.NewWriter / (*Writer).Write),
	// или отдельный EBML‑writer из этого же репозитория.
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Псевдо‑строка: замените на реальные вызовы из пакета
	// err = webm.Write(out, &m)
	// if err != nil {
	//     return err
	// }

	return nil
}
