package media_processor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

func patchWebMDuration(filePath string, newDurationMs float64) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return patchWebMDurationBytes(data, newDurationMs)
}

func patchWebMDurationBytes(data []byte, newDurationMs float64) ([]byte, error) {
	// EBML ID элемента Duration: 0x4489
	durationID := []byte{0x44, 0x89}
	idx := findBytesPattern(data, durationID)

	if idx == -1 {
		return nil, fmt.Errorf("duration element not found")
	}

	result := make([]byte, len(data))
	copy(result, data)

	// После ID идёт размер элемента (VINT)
	sizePos := idx + 2
	size, sizeLen := readVINT(data[sizePos:])
	valuePos := sizePos + sizeLen

	// Записываем новое значение Duration
	if size == 4 {
		bits := math.Float32bits(float32(newDurationMs))
		binary.BigEndian.PutUint32(result[valuePos:valuePos+4], bits)
	} else if size == 8 {
		bits := math.Float64bits(newDurationMs)
		binary.BigEndian.PutUint64(result[valuePos:valuePos+8], bits)
	} else {
		return nil, fmt.Errorf("unexpected duration size: %d", size)
	}

	return result, nil
}

func PatchWebMDurationReader(data []byte, newDurationMs float64) (*bytes.Reader, error) {
	patchedData, err := patchWebMDurationBytes(data, newDurationMs)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(patchedData), nil
}

func findBytesPattern(data, pattern []byte) int {
	for i := 0; i <= len(data)-len(pattern); i++ {
		found := true
		for j := 0; j < len(pattern); j++ {
			if data[i+j] != pattern[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}

// readVINT читает EBML Variable Integer и возвращает значение и длину
func readVINT(data []byte) (int, int) {
	if len(data) == 0 {
		return 0, 0
	}

	first := data[0]
	var length int
	var mask byte

	switch {
	case first&0x80 != 0:
		length, mask = 1, 0x7F
	case first&0x40 != 0:
		length, mask = 2, 0x3F
	case first&0x20 != 0:
		length, mask = 3, 0x1F
	case first&0x10 != 0:
		length, mask = 4, 0x0F
	case first&0x08 != 0:
		length, mask = 5, 0x07
	case first&0x04 != 0:
		length, mask = 6, 0x03
	case first&0x02 != 0:
		length, mask = 7, 0x01
	case first&0x01 != 0:
		length, mask = 8, 0x00
	default:
		return 0, 0
	}

	if len(data) < length {
		return 0, 0
	}

	value := int(first & mask)
	for i := 1; i < length; i++ {
		value = (value << 8) | int(data[i])
	}

	return value, length
}
