package camera

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net"
	"os"
	"sync"
	"../../websocket"
	"../../subscriber"
)

const (
	maxIllegalPacketCounter = 5
)

func Listen() {
	mutex := sync.Mutex{}
	clients := []websocket.WebSocketClient{}
	clientChannel := make(chan websocket.WebSocketClient)

	go func() {
		listener, _ := net.Listen(
			"tcp",
			os.Getenv("CLIENT_SERVER"),
		)
		fmt.Printf("Start client accepting server %v\n", listener.Addr())
		for {
			connection, err := listener.Accept()
			if err != nil {
				fmt.Printf("Failed to listen. retry again.")
				continue
			}

			go func() {
				// Handshake
				wsClient, err := websocket.Upgrade(connection)
				if err != nil {
					fmt.Printf("Disallowed to connect: %v\n", connection.RemoteAddr())
					// Close connection
					_ = connection.Close()
					return
				}
				clientChannel <- *wsClient
			}()
		}
	}()


	// First contact, We send an image which fulfilled black.
	buffer := bytes.NewBuffer([]byte{})
	filledImage := image.NewRGBA(
		image.Rect(
			0,
			0,
			640,
			480,
		),
	)

	_ = jpeg.Encode(
		buffer,
		filledImage,
		&jpeg.Options{
			100,
		},
	)

	go func() {
		for {
			select {
			case client := <-clientChannel:
				fmt.Printf("Client connected %v\n", client.Client.RemoteAddr())

				// Send black screen
				_ = client.Write(
					client.Encode(
						buffer.Bytes(),
						websocket.OpcodeBinary,
						true,
					),
				)

				client.StartListener(
					&clients,
					&mutex,
				)
			}
		}
	}()

	listener, _ := net.Listen(
		"tcp",
		os.Getenv("CAMERA_SERVER"),
	)
	fmt.Printf("Start camera receiving server %v\n", listener.Addr())

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to listen. retry again.")
			continue
		}

		go func() {
			fmt.Printf("[CAMERA] Connected from %v\n", connection.RemoteAddr())
			illegalPacketCounter := maxIllegalPacketCounter
			for {
				if illegalPacketCounter == 0 {
					fmt.Printf("Respond invalid frame data. retry to listen.\n")
					connection.Close()
					return
				}

				data, loops, err := subscriber.SubscribeImageStream(connection)

				if err != nil {
					illegalPacketCounter--
					continue
				}
				illegalPacketCounter = maxIllegalPacketCounter

				// Broadcast to connected all clients.
				websocket.Broadcast(
					&data,
					loops,
					&clients,
					&mutex,
				)

			}
		}()
	}
}