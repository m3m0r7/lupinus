package http

import (
	"fmt"
	"net"
	"os"
)

func Listen() {
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
			// Get headers
			
		}()
	}
}