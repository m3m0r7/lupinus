package main

import (
	"./servers/http/web"
	"./servers/streaming"
	"fmt"
	"github.com/joho/godotenv"
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
