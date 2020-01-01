package validator

import "bytes"

func IsImageJpeg(data []byte) bool {
	jpegMagicBytes := [][]byte{
		{0xff, 0xd8, 0xdd, 0xe0},
		{0xff, 0xd8, 0xff, 0xee},
		{0xff, 0xd8, 0xff, 0xdb},
	}

	for _, v := range jpegMagicBytes {
		if bytes.Equal(data[:4], v) {
			return true
		}
	}

	return false
}
