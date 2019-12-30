package main

import (
	"fmt"
	"net"
)

func main() {
	go func() {
		listener, _ := net.Listen(
			"tcp",
			"0.0.0.0:30000",
		)
		fmt.Printf("Start client accepting server %v\n", listener.Addr())
		for {
			connection, _ := listener.Accept()
			go func() {
				fmt.Printf("[Client] Connected from %v\n", connection.RemoteAddr())
				connection.Close()
			}()
		}
	}()

	listener, _ := net.Listen(
		"tcp",
		"0.0.0.0:31000",
	)

	fmt.Printf("Start camera receiving server %v\n", listener.Addr())
	for {
		connection, _ := listener.Accept()
		go func() {
			fmt.Printf("[CAMERA] Connected from %v\n", connection.RemoteAddr())
			connection.Close()
		}()
	}
}
