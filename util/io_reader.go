package util

import (
	"errors"
	"io"
	"net"
)

func ExpectToRead(stream net.Conn, expect int) ([]byte, error) {
	data := make([]byte, expect)
	n, _ := io.ReadFull(stream, data)
	if n == 0 {
		return nil, errors.New("End of Stream")
	}
	return data, nil
}
