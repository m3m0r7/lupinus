package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	for {
		connection, err := net.Dial("tcp", "localhost:30000")
		if err != nil {
			fmt.Println("Retry to connect for the testing")
			time.Sleep(5 * time.Second)
			continue
		}
		for {
			fmt.Println("recieve a data")
			b := make([]byte, 100)
			n, err := connection.Read(b)
			if err != nil {
				fmt.Printf("err = %+v\n", err)
				os.Exit(2)
			}

			fmt.Printf("receive data = %q\n", b[:n])
		}
	}
}
