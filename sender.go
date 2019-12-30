package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	for {
		connection, err := net.Dial("tcp", "localhost:31000")
		if err != nil {
			fmt.Println("Retry to connect for the testing")
			time.Sleep(5 * time.Second)
			continue
		}
		for {
			fmt.Println("Write a data")
			connection.Write(
				[]byte("camera data has been sent"),
			)
			time.Sleep(5 * time.Second)
		}
	}
}