package main

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"./servers/streaming"
	"./servers/http"
	"sync"
)

func main() {
	dir, _ := os.Getwd()
	if err := godotenv.Load(dir + "/.env"); err != nil {
		fmt.Printf("Failed to load an env file: %v\n", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Start listen servers
	go camera.Listen()
	go http.Listen()

	wg.Wait()
}