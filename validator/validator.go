package validator

import (
	"bytes"
)

func IsImageJpeg(data []byte) bool {
	jpegMagicBytes := [][]byte{
		{0xff, 0xd8},
	}

	for _, v := range jpegMagicBytes {
		if bytes.Equal(data[:2], v) {
			return true
		}
	}

	return false
}
