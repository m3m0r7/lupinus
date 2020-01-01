package util

import (
	"bytes"
	"io"
)

func ReadTo(reader io.Reader, dest *[]byte, delim []byte) error {
	packet := make([]byte, 1)
	for {
		_, err := reader.Read(packet)
		if err != nil {
			return err
		}
		*dest = append(*dest, packet...)
		if bytes.Equal(packet, delim) {
			break
		}
	}
	return nil
}