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
		var i int
		for {
			i++
			fmt.Println("Write a data")
			connection.Write(
				[]byte(fmt.Sprintf("camera data has been sent %d", i)),
			)
			time.Sleep(5 * time.Second)
		}
	}
}
