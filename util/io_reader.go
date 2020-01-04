package util

import (
	"errors"
	"net"
)

const (
	maxRemainingCounter = 5
)

func ExpectToRead(stream *net.Conn, expect int) ([]byte, error) {
	data := []byte{}
	remaining := expect
	remainingCounter := maxRemainingCounter
	for remaining > 0 {
		tmpRead := make([]byte, remaining)
		size, _ := (*stream).Read(tmpRead)
		data = append(data, tmpRead[:size]...)

		remaining -= size
		if size == 0 {
			remainingCounter--
			if remainingCounter == 0 {
				return nil, errors.New("Write failed")
			}
		}
	}
	return data, nil
}