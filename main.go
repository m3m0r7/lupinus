package main

import (
	"fmt"
	"net"
	"os"
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
				clientChannel <- connection
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

		clients := []net.Conn{}
		go func() {
			for {
				select {
				case client := <-clientChannel:
					fmt.Printf("Client connected %v\n", client.RemoteAddr())
					mutex.Lock()
					clients = append(clients, client)
					mutex.Unlock()
				}
			}
		}()
		for {
			connection, _ := listener.Accept()
			go func() {
				fmt.Printf("[CAMERA] Connected from %v\n", connection.RemoteAddr())
				read := make([]byte, 100)
				for {
					n, err := connection.Read(read)
					if err != nil {
						fmt.Printf("err = %+v\n", err)
						os.Exit(2)
					}
					data := read[:n]
					fmt.Printf("Data received: %q\n", data)
					mutex.Lock()
					for _, client := range clients {
						fmt.Printf("%v:%s\n", client.RemoteAddr(), data)
						if _, err := client.Write(data); err != nil {
							// Recreate new clients slice.
							fmt.Printf("Failed to write%v\n", client.RemoteAddr())

							tmpClients := []net.Conn{}
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