package util

import "io"

func ExpectToRead(stream io.Reader, expect int) ([]byte, error) {
	data := []byte{}
	remaining := expect
	for remaining > 0 {
		tmpRead := make([]byte, remaining)
		size, err := stream.Read(tmpRead)
		data = append(data, tmpRead...)
		if err != nil {
			return nil, err
		}

		remaining -= size
	}
	return data, nil
}