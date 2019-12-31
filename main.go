package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
	"./websocket"
)

func main() {
	var wg sync.WaitGroup

	clientChannel := make(chan websocket.WebSocketClient)

	wg.Add(1)
	go func() {
		listener, _ := net.Listen(
			"tcp",
			"0.0.0.0:30000",
		)
		fmt.Printf("Start client accepting server %v\n", listener.Addr())
		for {
			connection, _ := listener.Accept()
			go func() {
				// Handshake
				wsClient, err := websocket.UpgradeToWebSocketClient(connection)
				if err != nil {
					fmt.Printf("Disallowed to connect: %v\n", connection.RemoteAddr())
					return
				}
				clientChannel <- wsClient
			}()
		}
	}()

	go func() {
		var mutex sync.Mutex
		listener, _ := net.Listen(
			"tcp",
			"0.0.0.0:31000",
		)

		fmt.Printf("Start camera receiving server %v\n", listener.Addr())

		clients := []websocket.WebSocketClient{}
		go func() {
			for {
				select  {
				case client := <-clientChannel:
					fmt.Printf("Client connected %v\n", client.Client.RemoteAddr())
					mutex.Lock()
					clients = append(clients, client)
					mutex.Unlock()
				}
			}
		}()

		authKey := os.Getenv("AUTH_KEY")
		authKeySize := len(authKey)

		for {
			connection, _ := listener.Accept()

			go func() {
				fmt.Printf("[CAMERA] Connected from %v\n", connection.RemoteAddr())
				for {
					readAuthKey := make([]byte, authKeySize)
					receivedAuthKeySize, err := connection.Read(readAuthKey)
					if err != nil {
						fmt.Printf("err = %+v\n", err)

						// Retry to listen from the camera server.
						return
					}

					// Compare the received auth key and settled auth key.
					if string(readAuthKey[:receivedAuthKeySize]) != authKey {
						fmt.Printf("err = %+v\n", err)

						// Retry to listen from the camera server.
						return
					}

					// Receive frame size
					frameSize := make([]byte, 4)
					_, errReceivingFrameSize := connection.Read(frameSize)
					if errReceivingFrameSize != nil {
						fmt.Printf("err = %+v\n", err)

						// Retry to listen from the camera server.
						return
					}

					realFrameSize := binary.BigEndian.Uint32(frameSize)
					realFrame := make([]byte, realFrameSize)

					receivedImageDataSize, errReceivingRealFrame := connection.Read(realFrame)
					if errReceivingRealFrame != nil {
						fmt.Printf("err = %+v\n", err)

						// Retry to listen from the camera server.
						return
					}

					data := realFrame[:receivedImageDataSize]

					mutex.Lock()
					for _, client := range clients {
						if _, err := client.Client.Write(data); err != nil {
							// Recreate new clients slice.
							fmt.Printf("Failed to write%v\n", client.Client.RemoteAddr())

							tmpClients := []websocket.WebSocketClient{}
							for _, tmpClient := range clients {
								if tmpClient != client {
									tmpClients = append(tmpClients, tmpClient)
								}
							}
							clients = tmpClients
						}
					}
					mutex.Unlock()
				}
			}()
		}
	}()

	wg.Wait()
}