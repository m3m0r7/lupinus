package controller

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"lupinus/servers/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func RequestSlack(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	// has challenge token?
	jsonData := map[string]interface{}{}
	err := json.Unmarshal(clientMeta.Payload, &jsonData)

	if err != nil {
		return &http.HttpBody{
				Payload: http.Payload{
					"status": 401,
					"error":  "Failed to decode a json",
				},
			},
			&http.HttpHeader{
				Status: 401,
			}
	}

	if _, ok := jsonData["type"]; ok {
		if _, ok := jsonData["challenge"]; ok {
			if jsonData["type"] == "url_verification" {
				return &http.HttpBody{
						RawMode: true,
						Payload: http.Payload{
							"body": jsonData["challenge"].(string),
						},
					},
					&http.HttpHeader{}
			}
		}
	}

	_, hasText := jsonData["text"]
	_, hasChannel := jsonData["channel"]

	if !hasText || !hasChannel {
		return &http.HttpBody{}, &http.HttpHeader{}
	}

	text := jsonData["text"].(string)

	// check for the bot mention
	if !strings.Contains(text, os.Getenv("SLACK_BOT_NAME")) {
		return &http.HttpBody{}, &http.HttpHeader{}
	}

	// the deploy command
	if strings.Contains(text, "deploy") {

		var slackApi = slack.New(os.Getenv("SLACK_TOKEN"))

		_, _, err := slackApi.PostMessage(
			os.Getenv("SLACK_CHANNEL"),
			slack.MsgOptionText(
				"OK! I will deploy *"+os.Getenv("APPLICATION_NAME")+"* new system! Please wait a few minutes...",
				false,
			),
		)

		if err != nil {
			return &http.HttpBody{}, &http.HttpHeader{}
		}

		deployErr := exec.Command(
			"su",
			os.Getenv("DEPLOY_USER"),
			"-c",
			"cd " + os.Getenv("DEPLOY_DIRECTORY") + " && git pull",
		).Run()

		if deployErr != nil {
			_, _, err = slackApi.PostMessage(
				os.Getenv("SLACK_CHANNEL"),
				slack.MsgOptionText(
						fmt.Sprintf(
							"Deploy failed :crying_cat_face:\n\n> [Reason] %v",
							deployErr,
						),
					false,
				),
			)
			return &http.HttpBody{}, &http.HttpHeader{}
		}

		_, _, err = slackApi.PostMessage(
			os.Getenv("SLACK_CHANNEL"),
			slack.MsgOptionText(
				"Deploy finished. I will restart the system.",
				false,
			),
		)

		if err != nil {
			return &http.HttpBody{}, &http.HttpHeader{}
		}

		return &http.HttpBody{
			AfterProcess: func(meta http.HttpClientMeta) {
				// Restart
				exec.Command(
					"kill",
					strconv.Itoa(
						os.Getpid(),
					),
				).Run()
			},
		}, &http.HttpHeader{}
	}

	return &http.HttpBody{}, &http.HttpHeader{}
}
