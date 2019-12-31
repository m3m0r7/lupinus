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
		counter := 1
		for {
			fmt.Printf("Write a data %d\n", counter)

			connection.Write(
				[]byte(fmt.Sprintf("camera data has been sent %d", counter)),
			)
			counter++
			time.Sleep(5 * time.Second)
		}
	}
}