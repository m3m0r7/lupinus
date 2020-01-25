package controller

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"lupinus/servers/http"
	"os"
	"os/exec"
	"strings"
)

func RequestSlack(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	filePath := os.Getenv("DEPLOY_DIRECTORY") + "/server.log"
	handle, _ := os.OpenFile(
		filePath,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)

	defer handle.Close()

	readReceivedEventIdAll, _ := ioutil.ReadFile(filePath)

	successReturn := http.HttpBody{
		Payload: http.Payload{
			"message": "OK",
		},
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
					&http.HttpHeader{
						Status: 200,
					}
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
	eventId := event["client_msg_id"]

	// Check already received event id.
	if strings.Contains(string(readReceivedEventIdAll), eventId.(string)+"\n") {
		return &http.HttpBody{
				Payload: http.Payload{
					"message": "Already received.",
				},
			}, &http.HttpHeader{
				Status: 200,
			}
	}

	defer func() {
		handle.Write(
			[]byte(eventId.(string) + "\n"),
		)
	}()

	// check for the bot mention
	if !strings.Contains(text, os.Getenv("SLACK_BOT_NAME")) {
		return &http.HttpBody{
				Payload: http.Payload{
					"message": "OK",
				},
			}, &http.HttpHeader{
				Status: 200,
			}
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
			return &http.HttpBody{

					Payload: http.Payload{
						"message": "OK",
					},
				}, &http.HttpHeader{
					Status: 200,
				}
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
						"Deploy failed :crying_cat_face: \n\n> [Reason] %v",
						deployErr,
					),
					false,
				),
			)
			return &http.HttpBody{

					Payload: http.Payload{
						"message": "OK",
					},
				}, &http.HttpHeader{
					Status: 200,
				}
		}

		_, _, err = slackApi.PostMessage(
			os.Getenv("SLACK_CHANNEL"),
			slack.MsgOptionText(
				"Deploy finished. I will restart the system.",
				false,
			),
		)

		if err != nil {
			return &http.HttpBody{
					Payload: http.Payload{
						"message": "OK",
					},
				}, &http.HttpHeader{
					Status: 200,
				}
		}

		return &successReturn, &http.HttpHeader{
			Status: 200,
		}
	}

	if strings.Contains(text, "restart") {
		_, _, err = slackApi.PostMessage(
			os.Getenv("SLACK_CHANNEL"),
			slack.MsgOptionText(
				"OK! I will restart the system.",
				false,
			),
		)

		return &successReturn, &http.HttpHeader{
			Status: 200,
		}
	}

	if strings.Contains(text, "clear") {
		handle.Close()

		// Remove server.log
		err := os.Remove(filePath)
		if err != nil {
			// Nothing to do...
		}

		_, _, err = slackApi.PostMessage(
			os.Getenv("SLACK_CHANNEL"),
			slack.MsgOptionText(
				"OK! I will clear caches.",
				false,
			),
		)

		return &successReturn, &http.HttpHeader{
			Status: 200,
		}
	}


	if strings.Contains(text, "cert") {
		output, err := exec.Command(
			"certbot renew",
		).Output()

		if err != nil {
			_, _, err = slackApi.PostMessage(
				os.Getenv("SLACK_CHANNEL"),
				slack.MsgOptionText(
					"Failed to update certificate file :crying_cat_face: \n\n" +
					"> [Reason] " + fmt.Sprintf("%v", err) + "\n" +
					"> [Log] " + string(output),
					false,
				),
			)
			return &http.HttpBody{
					Payload: http.Payload{
						"message": "OK",
					},
				}, &http.HttpHeader{
				Status: 200,
			}
		}

		updateErr := exec.Command(
			"cd "+os.Getenv("DEPLOY_DIRECTORY")+" && php update_certificate.php",
		).Run()

		if updateErr != nil {
			_, _, err = slackApi.PostMessage(
				os.Getenv("SLACK_CHANNEL"),
				slack.MsgOptionText(
					"Failed to update certificate file :crying_cat_face: \n\n" +
						"> [Reason] " + fmt.Sprintf("%v", updateErr) + "\n",
					false,
				),
			)
			return &http.HttpBody{
					Payload: http.Payload{
						"message": "OK",
					},
				}, &http.HttpHeader{
					Status: 200,
				}
		}

		_, _, err = slackApi.PostMessage(
			os.Getenv("SLACK_CHANNEL"),
			slack.MsgOptionText(
				"Certificate file was updated :key: I will restart the system.",
				false,
			),
		)

		return &successReturn, &http.HttpHeader{
			Status: 200,
		}
	}

	if strings.Contains(text, "help") {
		_, _, err = slackApi.PostMessage(
			os.Getenv("SLACK_CHANNEL"),
			slack.MsgOptionText(
				"*clear* - Clear bot caches from the app.\n" +
					"*restart* - Restart the server.\n" +
					"*deploy* - Deploy new system from GitHub repository.\n" +
					"*cert* - Update certificate file for SSL.\n" +
					"*stats* - Show the application stats.\n" +
					"*help* - Show this help.\n",
				false,
			),
		)
		return &http.HttpBody{
				Payload: http.Payload{
					"message": "OK",
				},
			}, &http.HttpHeader{
			Status: 200,
		}
	}

	_, _, err = slackApi.PostMessage(
		os.Getenv("SLACK_CHANNEL"),
		slack.MsgOptionText(
			"Sorry, entered command is not available :scream_cat: if you want to know commands then mention me with `help`.",
			false,
		),
	)

	return &http.HttpBody{
			Payload: http.Payload{
				"message": "OK",
			},
		}, &http.HttpHeader{
			Status: 200,
		}
}
