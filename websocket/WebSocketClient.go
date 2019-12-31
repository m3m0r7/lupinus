package websocket

import (
	"net"
)

const ID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

type WebSocketClient struct {
	Client net.Conn
	Handshake bool
}

func UpgradeToWebSocketClient(client net.Conn) (WebSocketClient, error) {
	wsClient := WebSocketClient{
		Client: client,
	}

	// TODO: Read headers

	client.Write(
		[]byte("HTTP/1.1 101 Switching Protocols\n"),
	)
	client.Write(
		[]byte("Upgrade: websocket\n"),
	)
	client.Write(
		[]byte("Connection: upgrade\n"),
	)
	client.Write(
		[]byte("Sec-WebSocket-Accept: \n"),
	)
	client.Write(
		[]byte("\n"),
	)

	return wsClient, nil
}
