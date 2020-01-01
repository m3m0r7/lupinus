package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"net"
	"reflect"
	"../util"
	parent "../client"
	"sync"
)

// WebSocket key ID.
const ID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// Operation codes.
const (
	OpcodeFin     = 0x00
	OpcodeMessage = 0x01
	OpcodeBinary  = 0x02
	OpcodeClose   = 0x08
	OpcodePing    = 0x09
	OpcodePong    = 0x0A
)

const (
	maxHeadersLine = 128
	maxRetryToWriteCounter = 3
)

type WebSocketClient struct {
	Pipe net.Conn
	Headers []parent.ClientHeader
	Handshake bool
}

func (client *WebSocketClient) RemoveFromClientsWithLock(clients []WebSocketClient, mutex *sync.Mutex) []WebSocketClient {
	mutex.Lock()
	defer mutex.Unlock()
	return client.RemoveFromClients(clients)
}

func (client *WebSocketClient) RemoveFromClients(clients []WebSocketClient) []WebSocketClient {
	tmpClients := []WebSocketClient{}
	for _, tmpClient := range clients {
		if !reflect.DeepEqual(tmpClient, client) {
			tmpClients = append(tmpClients, tmpClient)
		}
	}
	return tmpClients
}

func (client *WebSocketClient) findHeaderByKey(key string) (*parent.ClientHeader, error) {
	for _, clientHeader := range client.Headers {
		if clientHeader.Key == key {
			return &clientHeader, nil
		}
	}
	return nil, errors.New("Not found value from header with key.")
}

func (client *WebSocketClient) Decode() ([]byte, int, error) {
	kindByte := make([]byte, 2)
	_, err := client.Pipe.Read(kindByte)
	if err != nil {
		return nil, -1, err
	}

	isFin := kindByte[0] >> 7
	if isFin > 1 || isFin < 0 {
		return nil, -1, errors.New("Invalid fin flag.")
	}

	opcode := ((kindByte[0] << 4) & 0xff) >> 4

	switch opcode {
	case OpcodeFin:
	case OpcodeMessage:
	case OpcodeBinary:
	case OpcodeClose:
	case OpcodePing:
	case OpcodePong:
		// nothing to do
		break
	default:
		return nil, -1, errors.New("Invalid Operation Code")
	}

	maskFlag := (kindByte[1] >> 7) & 0xff
	if maskFlag > 1 || maskFlag < 0 {
		return nil, -1, errors.New("Invalid mask flag.")
	}

	receivedType := ((kindByte[1] << 1) & 0xff) >> 1
	var length int
	if receivedType == 126 {
		readMore := make([]byte, 2)
		_, err = client.Pipe.Read(readMore)
		length = int(binary.BigEndian.Uint16(readMore))
	} else if receivedType == 127 {
		readMore := make([]byte, 8)
		_, err = client.Pipe.Read(readMore)
		length = int(binary.BigEndian.Uint64(readMore))
	} else {
		length = int(receivedType)
	}

	maskData := make([]byte, 4)
	if maskFlag == 1 {
		client.Pipe.Read(maskData)
	}

	payload := make([]byte, length)
	client.Pipe.Read(payload)

	if maskFlag == 1 {
		for i, char := range payload {
			payload[i] = char ^ maskData[i % 4]
		}
	}

	return payload, int(opcode), nil
}

func (client *WebSocketClient) Write(data []byte) error {
	return parent.Write(client.Pipe, data)
}

func (client *WebSocketClient) Encode(payload []byte, opcode int, isFin bool) []byte {
	length := len(payload)
	sendType := 0
	if length > 0xffff {
		sendType = 127
	} else if length <= 0xffff && length >= 126 {
		sendType = 126
	} else {
		sendType = length
	}

	body := []byte{}
	finFlag := 0
	if isFin {
		finFlag = 128
	}

	body = append(body, util.Uint2bytes(finFlag + opcode, 1)...)
	body = append(body, util.Uint2bytes(sendType, 1)...)

	if sendType == 126 {
		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, uint16(length))
		body = append(body, size...)
	} else if sendType == 127 {
		size := make([]byte, 8)
		binary.BigEndian.PutUint64(size, uint64(length))
		body = append(body, size...)
	}

	return append(body, payload...)
}

func Upgrade(conn net.Conn) (*WebSocketClient, error) {
	wsClient := WebSocketClient{
		Pipe: conn,
	}

	headers, _ := parent.GetAllHeaders(conn)
	result, err := parent.FindHeaderByKey(
		&headers,
		"sec-websocket-key",
	)

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

	sendHeaders :=  "HTTP/1.1 101 Switching Protocols\n" +
		"Upgrade: websocket\n" +
		"Connection: upgrade\n" +
		"Sec-WebSocket-Accept: " + wsAcceptHeader + "\n" +
		"\n"

	_, err = conn.Write([]byte(sendHeaders))
	if err != nil {
		return nil, errors.New("Failed to write")
	}

	wsClient.Handshake = true
	wsClient.Headers = headers

	return &wsClient, nil
}
