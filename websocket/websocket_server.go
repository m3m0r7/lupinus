package websocket

import (
	"fmt"
)

func (client *WebSocketClient) ReceivedClose(response []byte) error {
	err := client.Write(
		client.Encode(
			response,
			OpcodeClose,
			true,
		),
	)
	return err
}

func (client *WebSocketClient) ReceivedPing(response []byte) error {
	err := client.Write(
		client.Encode(
			response,
			OpcodePong,
			true,
		),
	)
	return err
}


func (client *WebSocketClient) StartListener(clients *[]WebSocketClient, lostClientChannel chan WebSocketClient) {
	go func () {
		for {
			receivedResponse, opcode, err := client.Decode()
			if err != nil {
				err = client.Pipe.Close()
				lostClientChannel <- *client
				return
			}

			switch opcode {
			case OpcodeClose:
				err = client.ReceivedClose(receivedResponse)
				lostClientChannel <- *client
				return
			case OpcodePing:
				err = client.ReceivedPing(receivedResponse)
				if err != nil {
					err = client.Pipe.Close()
					lostClientChannel <- *client
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

func Broadcast(data [][]byte, size int, clients *[]WebSocketClient, lostClientChannel chan WebSocketClient) {
	for _, client := range *clients {
		go func () {
			for i := 0; i < size; i++ {
				opcode := OpcodeMessage
				if i > 0 {
					opcode = OpcodeFin
				}
				err := client.Write(
					client.Encode(
						data[i],
						opcode,
						(i + 1) == size,
					),
				)
				if err != nil {
					// Recreate new clients slice.
					fmt.Printf("Failed to write %v, %v\n", client.Pipe.RemoteAddr(), err)
					lostClientChannel <- client
					break
				}
			}
		}()
	}
}