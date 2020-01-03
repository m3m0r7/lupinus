package controller

import (
	"../../../../model"
	"../../../../util"
	"../../../http"
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

	username := util.GetFromMap("username", jsonData)
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
			Payload: map[string]interface{}{
				"code": 100,
				"message": "Failed to authorize",
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
		Payload: map[string]interface{}{
			"message": "Sign-in was succeeded",
		},
	},
	&http.HttpHeader{
		Status: 404,
	}
}
