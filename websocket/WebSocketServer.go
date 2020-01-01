package websocket

import (
	"fmt"
	"sync"
)

func (client *WebSocketClient) ReceivedClose(response []byte) error {
	_, err := client.Client.Write(
		client.Encode(
			response,
			OpcodeClose,
			true,
		),
	)
	return err
}

func (client *WebSocketClient) ReceivedPing(response []byte) error {
	_, err := client.Client.Write(
		client.Encode(
			response,
			OpcodePong,
			true,
		),
	)
	return err
}


func (client *WebSocketClient) StartListener(clients *[]WebSocketClient, mutex *sync.Mutex) {
	mutex.Lock()
	*clients = append(*clients, *client)
	mutex.Unlock()

	go func () {
		for {
			receivedResponse, opcode, err := client.Decode()
			if err != nil {
				err = client.Client.Close()
				*clients = client.RemoveFromClientsWithLock(*clients, mutex)
				return
			}

			switch opcode {
			case OpcodeClose:
				err = client.ReceivedClose(receivedResponse)
				*clients = client.RemoveFromClientsWithLock(*clients, mutex)
				return
			case OpcodePing:
				err = client.ReceivedPing(receivedResponse)
				if err != nil {
					err = client.Client.Close()
					*clients = client.RemoveFromClientsWithLock(*clients, mutex)
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

func Broadcast(data *[][]byte, size int, clients *[]WebSocketClient, mutex *sync.Mutex) {
	for _, client := range *clients {
		go func () {
			for i := 0; i < size; i++ {
				opcode := OpcodeBinary
				if i > 0 {
					opcode = OpcodeFin
				}
				_, err := client.Client.Write(
					client.Encode(
						(*data)[i],
						opcode,
						(i + 1) == size,
					),
				)
				if err != nil {
					// Recreate new clients slice.
					fmt.Printf("Failed to write%v\n", client.Client.RemoteAddr())
					*clients = client.RemoveFromClientsWithLock(*clients, mutex)
					return
				}
			}
		}()
	}
}