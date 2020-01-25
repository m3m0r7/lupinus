package controller

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"lupinus/servers/http"
	"os"
	"os/exec"
	"strings"
)

func RequestSlack(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	successReturn := http.HttpBody{
		AfterProcess: func(meta http.HttpClientMeta) {
			// Restart
			exec.Command(
				"reboot",
			).Run()
		},
	}

	// has challenge token?
	jsonData := map[string]interface{}{}
	err := json.Unmarshal(clientMeta.Payload, &jsonData)

	var slackApi = slack.New(os.Getenv("SLACK_TOKEN"))

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

	if _, ok := jsonData["event"]; !ok {
		return &http.HttpBody{}, &http.HttpHeader{}
	}

	event := jsonData["event"].(map[string]interface{})

	_, hasText := event["text"]
	_, hasChannel := event["channel"]

	if !hasText || !hasChannel {
		return &http.HttpBody{}, &http.HttpHeader{}
	}

	text := event["text"].(string)

	// check for the bot mention
	if !strings.Contains(text, os.Getenv("SLACK_BOT_NAME")) {
		return &http.HttpBody{}, &http.HttpHeader{}
	}

	// the deploy command
	if strings.Contains(text, "deploy") {
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
			"cd "+os.Getenv("DEPLOY_DIRECTORY")+" && git pull",
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

		return &successReturn, &http.HttpHeader{}
	}

	if strings.Contains(text, "restart") {
		_, _, err = slackApi.PostMessage(
			os.Getenv("SLACK_CHANNEL"),
			slack.MsgOptionText(
				"OK!. I will restart the system.",
				false,
			),
		)

		return &successReturn, &http.HttpHeader{}
	}

	return &http.HttpBody{}, &http.HttpHeader{}
}
