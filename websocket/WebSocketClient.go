package websocket

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net"
	"strings"
)

const ID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

type ClientHeader struct {
	Key string
	Value string
}

type WebSocketClient struct {
	Client net.Conn
	Headers []ClientHeader
	Handshake bool
}

func (client *WebSocketClient) findHeaderByKey(key string) (*ClientHeader, error) {
	for _, clientHeader := range client.Headers {
		if clientHeader.Key == key {
			return &clientHeader, nil
		}
	}
	return nil, errors.New("Not found value from header with key.")
}

func Upgrade(client net.Conn) (*WebSocketClient, error) {
	headers := []ClientHeader{}
	scanner := bufio.NewScanner(client)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		result := strings.Split(line, ":")

		// If not exists :, set key to zero value
		clientHeader := ClientHeader{}
		if len(result) <= 1 {
			clientHeader = ClientHeader {
				Value: strings.Trim(result[0], " "),
			}
		} else {
			clientHeader = ClientHeader {
				Key: strings.ToLower(strings.Trim(result[0], " ")),
				Value: strings.Trim(result[1], " "),
			}
		}
		headers = append(headers, clientHeader)
	}

	wsClient := WebSocketClient{
		Client: client,
		Handshake: true,
		Headers: headers,
	}

	client.Write(
		[]byte("HTTP/1.1 101 Switching Protocols\n"),
	)
	client.Write(
		[]byte("Upgrade: websocket\n"),
	)
	client.Write(
		[]byte("Connection: upgrade\n"),
	)

	result, err := wsClient.findHeaderByKey("sec-websocket-key")

	if err != nil {
		return nil, errors.New("Connected client is invalid.")
	}

	cryptedToSha1 := sha1.New()
	cryptedToSha1.Write(
		[]byte(result.Value + ID),
	)

	wsAcceptHeader := base64.StdEncoding.EncodeToString(
		cryptedToSha1.Sum(
			nil,
		),
	)

	client.Write(
		[]byte("Sec-WebSocket-Accept: " + wsAcceptHeader + "\n"),
	)
	client.Write(
		[]byte("\n"),
	)

	return &wsClient, nil
}
