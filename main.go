package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	clientChannel := make(chan net.Conn)

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
				fmt.Printf("[CLIENT] Connected from %v\n", connection.RemoteAddr())
				clientChannel <- connection
			}()
		}
	}()

	go func() {
		listener, _ := net.Listen(
			"tcp",
			"0.0.0.0:31000",
		)

		fmt.Printf("Start camera receiving server %v\n", listener.Addr())

		clients := []net.Conn{}
		cameraImageData := make(chan []byte)
		for {
			connection, _ := listener.Accept()
			go func() {
				fmt.Printf("[CAMERA] Connected from %v\n", connection.RemoteAddr())
				for {
					var read []byte
					status, _ := connection.Read(read)
					_ = status
					fmt.Printf("Data received: %s\n", read)
					cameraImageData <- read
					select {
						case imageData := <- cameraImageData:
							for _, client := range clients {
								fmt.Printf("%v%s\n", client.RemoteAddr(), imageData)
							}
						case client := <- clientChannel:
							fmt.Printf("[CLIENT] Connected from %v\n", client.RemoteAddr())
							clients = append(clients, client)
					}
				}
			}()
		}
	}()

	wg.Wait()
}
