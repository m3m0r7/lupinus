package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net"
	"os"
	"sync"
	"./websocket"
	"./subscriber"
)

const (
	maxIllegalPacketCounter = 5
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load an env file: %v\n", err)
		return
	}

	var wg sync.WaitGroup

	clientChannel := make(chan websocket.WebSocketClient)

	wg.Add(1)
	go func() {
		listener, _ := net.Listen(
			"tcp",
			os.Getenv("CLIENT_SERVER"),
		)
		fmt.Printf("Start client accepting server %v\n", listener.Addr())
		for {
			connection, _ := listener.Accept()
			go func() {
				// Handshake
				wsClient, err := websocket.Upgrade(connection)
				if err != nil {
					fmt.Printf("Disallowed to connect: %v\n", connection.RemoteAddr())

					// Close connection
					connection.Close()
					return
				}
				clientChannel <- *wsClient
			}()
		}
	}()

	go func() {
		mutex := sync.Mutex{}
		listener, _ := net.Listen(
			"tcp",
			os.Getenv("CAMERA_SERVER"),
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

					go func () {
						for {
							result, opcode, err := client.Decode()
							if err != nil {
								err = client.Client.Close()
								mutex.Lock()
								clients = client.RemoveFromClients(clients)
								mutex.Unlock()
								return
							}

							switch opcode {
							case websocket.OpcodeClose:
								_, err := client.Client.Write(
									client.Encode(
										result,
										websocket.OpcodeClose,
										true,
									),
								)
								err = client.Client.Close()
								_ = err

								mutex.Lock()
								clients = client.RemoveFromClients(clients)
								mutex.Unlock()
								return
							case websocket.OpcodePing:
								_, err := client.Client.Write(
									client.Encode(
										result,
										websocket.OpcodePong,
										true,
									),
								)
								if err != nil {
									err = client.Client.Close()
									mutex.Lock()
									clients = client.RemoveFromClients(clients)
									mutex.Unlock()
									return
								}
								break
							default:
								// Nothing to do
								break
							}
						}
					}()
				}
			}
		}()

		for {
			connection, _ := listener.Accept()

			go func() {
				fmt.Printf("[CAMERA] Connected from %v\n", connection.RemoteAddr())
				illegalPacketCounter := maxIllegalPacketCounter
				for {
					if illegalPacketCounter == 0 {
						fmt.Printf("connected from illegal connection.")
						connection.Close()
						return
					}

					data, loops, err := subscriber.SubscribeImageStream(connection)
					if err != nil {
						illegalPacketCounter--
						continue
					}

					illegalPacketCounter = maxIllegalPacketCounter

					for _, client := range clients {
						go func () {
							for i := 0; i < loops; i++ {
								opcode := websocket.OpcodeBinary
								if i > 0 {
									opcode = websocket.OpcodeFin
								}
								_, err := client.Client.Write(
									client.Encode(
										data[i],
										opcode,
										i == loops,
									),
								)
								if err != nil {
									// Recreate new clients slice.
									fmt.Printf("Failed to write%v\n", client.Client.RemoteAddr())

									mutex.Lock()
									clients = client.RemoveFromClients(clients)
									mutex.Unlock()
								}
							}
						}()
					}
				}
			}()
		}
	}()

	wg.Wait()
}