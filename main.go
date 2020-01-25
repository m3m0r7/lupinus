package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
  "github.com/nlopes/slack"
  "lupinus/servers/http/web"
	"lupinus/servers/streaming"
  "os"
  "sync"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load an env file: %v\n", err)
		return
	}

	var slackApi = slack.New(os.Getenv("SLACK_TOKEN"))
	_, _, err := slackApi.PostMessage(
		os.Getenv("SLACK_CHANNEL"),
		slack.MsgOptionText("Server started :heart_eyes: ", false),
	)
	_ = err

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
