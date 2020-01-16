package controller

import (
	"encoding/json"
	"lupinus/client"
	"lupinus/model"
	"lupinus/servers/http"
	"lupinus/util"
	"os"
)

func RequestSignin(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	// Do not allowed received method.
	if clientMeta.Method != "POST" {
		return nil, nil
	}

	jsonData := map[string]interface{}{}
	err := json.Unmarshal(clientMeta.Payload, &jsonData)

	if err != nil {
		return nil, nil
	}

	authKeyHeader, err := client.FindHeaderByKey(clientMeta.Headers, "x-auth-key")
	if err != nil ||
		os.Getenv("AUTH_KEY") != (*authKeyHeader).Value {
		// Not exists a session
		return &http.HttpBody{
				Payload: http.Payload{
					"status": 401,
					"error":  "Unauthorized",
				},
			},
			&http.HttpHeader{
				Status: 401,
			}
	}

	username := util.GetFromMap("id", jsonData)
	password := util.GetFromMap("password", jsonData)

	if username == nil || password == nil {
		return nil, nil
	}

	user := *model.InitUser()
	userData := user.Find(
		username.(string),
		password.(string),
	)

	// Not found
	if userData == nil {
		return &http.HttpBody{
				Payload: http.Payload{
					"status": 401,
					"error":  "Failed to authorize",
				},
			},
			&http.HttpHeader{
				Status: 401,
			}
	}

	// Create session
	session := http.CreateSession()
	session.Write("id", username.(string))

	return &http.HttpBody{
			Payload: http.Payload{
				"status": 200,
			},
		},
		&http.HttpHeader{
			Status: 200,
		}
}
