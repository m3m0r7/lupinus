package controller

import (
	"lupinus/model"
	"lupinus/util"
	"lupinus/servers/http"
	"encoding/json"
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
		return &http.HttpBody {
			Payload: http.Payload{
				"status": 401,
				"error": "Failed to authorize",
			},
		},
		&http.HttpHeader{
			Status: 401,
		}
	}

	// Create session
	session := http.CreateSession()
	session.Write("id", username.(string))

	return &http.HttpBody {
		Payload: http.Payload{
			"status": 200,
		},
	},
	&http.HttpHeader{
		Status: 200,
	}
}
