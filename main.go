package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"lupinus/servers/http/web"
	"lupinus/servers/streaming"
	"sync"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load an env file: %v\n", err)
		return
	}

	ctx := context.Background()
	defer func(ctx context.Context) {
		childCtx, cancel := context.WithCancel(ctx)
		cancel()

		err := recover()
		if err != nil {
			startup(childCtx)
		}

	}(ctx)

	wg := sync.WaitGroup{}
	wg.Add(1)

	startup(ctx)

	wg.Wait()
}

func startup(ctx context.Context) {
	// Start listen servers
	go streaming.ListenCameraStreaming(ctx)
	go web.Listen(ctx)
}
