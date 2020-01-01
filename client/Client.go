package client

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"../servers/http"
)

const (
	maxHeadersLine = 128
	maxRetryToWriteCounter = 3
)

func FindHeaderByKey(headers []http.ClientHeader, key string) (*http.ClientHeader, error) {
	for _, clientHeader := range headers {
		if clientHeader.Key == key {
			return &clientHeader, nil
		}
	}
	return nil, errors.New("Not found value from header with key.")
}

func GetAllHeaders(conn net.Conn) ([]http.ClientHeader, error) {

	headers := []http.ClientHeader{}
	scanner := bufio.NewScanner(conn)

	remaining := maxHeadersLine
	for scanner.Scan() {
		if remaining == 0 {
			return nil, errors.New("Requested headers are overflow.")
		}
		line := scanner.Text()
		if line == "" {
			break
		}
		result := strings.Split(line, ":")

		// If not exists :, set key to zero value
		clientHeader := http.ClientHeader{}
		if len(result) <= 1 {
			clientHeader = http.ClientHeader {
				Value: strings.Trim(result[0], " "),
			}
		} else {
			clientHeader = http.ClientHeader {
				Key: strings.ToLower(strings.Trim(result[0], " ")),
				Value: strings.Trim(result[1], " "),
			}
		}
		headers = append(headers, clientHeader)
		remaining--
	}

	return headers, nil
}

func Write(conn net.Conn, data []byte) error {
	realSize := len(data)
	remaining := len(data)
	read := 0
	writeRetryCount := maxRetryToWriteCounter
	for remaining > 0 {
		n, err := conn.Write(data[read:realSize])
		if err != nil {
			if writeRetryCount == 0 {
				// If write is missed, close connection
				_ = conn.Close()
				return err
			}
			writeRetryCount--
		}
		read += n
		remaining -= n
	}
	return nil
}