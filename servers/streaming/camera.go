package streaming

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"lupinus/helper"
	"lupinus/share"
	"lupinus/subscriber"
	"lupinus/util"
	"lupinus/websocket"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	maxIllegalPacketCounter = 5

	// Publics
	UpdateStaticImageInterval = 30
)

var UpdateTime = time.Now().Unix()
var NextUpdateTime = time.Now().Unix()

func ListenCameraStreaming() {
	mutex := sync.Mutex{}
	clients := []websocket.WebSocketClient{}
	clientChannel := make(chan websocket.WebSocketClient)
	lostClientChannel := make(chan websocket.WebSocketClient)

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
				wsClient, err := websocket.Upgrade(&connection)
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
				fmt.Printf("Client connected %v\n", client.Pipe.RemoteAddr())

				// Send black screen
				_ = client.Write(
					client.Encode(
						util.Byte2base64URI(buffer.Bytes()),
						websocket.OpcodeMessage,
						true,
					),
				)

				mutex.Lock()
				clients = append(clients, client)
				mutex.Unlock()

				client.StartListener(
					&clients,
					lostClientChannel,
				)
				break
			case client := <-lostClientChannel:
				mutex.Lock()
				tmpClients := []websocket.WebSocketClient{}
				for _, tmpClient := range clients {
					if reflect.DeepEqual(tmpClient, client) {
						_ = client.Pipe.Close()
						continue
					}
					tmpClients = append(tmpClients, tmpClient)
				}

				// Update clients slice
				clients = tmpClients
				mutex.Unlock()
			break
			}
		}
	}()

	for {
		hostname := strings.Split(os.Getenv("CAMERA_SERVER"), ":")
		ip, _ := net.LookupIP(hostname[0])
		port, _ := strconv.Atoi(hostname[1])
		addr := net.TCPAddr{
			IP: ip[0],
			Port: port,
		}
		listener, _ := net.ListenTCP(
			"tcp",
			&addr,
		)

		fmt.Printf("Start camera receiving server %v\n", listener.Addr())

		for {
			connection, err := listener.AcceptTCP()
			connection.SetKeepAlive(true)
			if err != nil {
				fmt.Printf("Failed to listen. retry again.\n")
				break
			}

			fmt.Printf("[CAMERA] Connected from %v\n", connection.RemoteAddr())
			illegalPacketCounter := maxIllegalPacketCounter
			for {
				if illegalPacketCounter == 0 {
					fmt.Printf("Respond invalid frame data. retry to listen.\n")
					connection.Close()
					break
				}

				frameData, data, loops, err := subscriber.SubscribeImageStream(connection)

				// FIX Golang cannot read buffered data.
				if frameData == nil && data == nil && loops == -1 && err == nil {
					// io.EOF
					break
				}

				if err != nil {
					fmt.Printf("Error has occurred: %v\n", err)
					illegalPacketCounter--
					continue
				}

				// proceed favorite procedures
				share.ProceedProcedure(
					"favorite",
					frameData,
				)

				currentTime := time.Now().Unix()
				if NextUpdateTime < currentTime {
					UpdateTime = currentTime
					NextUpdateTime = currentTime + UpdateStaticImageInterval

					// create image
					helper.CreateStaticImage(frameData, "record/image.jpg")
				}

				if err != nil {
					fmt.Printf("err = %v\n", err)
					illegalPacketCounter--
					continue
				}
				illegalPacketCounter = maxIllegalPacketCounter

				go func () {
					mutex.Lock()
					defer mutex.Unlock()
					// Broadcast to connected all clients.
					websocket.Broadcast(
						data,
						loops,
						clients,
						lostClientChannel,
					)
				}()
			}
		}
	}
}
