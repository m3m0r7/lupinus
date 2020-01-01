package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"./servers/streaming"
	"./servers/http/web"
	"sync"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load an env file: %v\n", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Start listen servers
	go camera.Listen()
	go web.Listen()

	wg.Wait()
}
