package controller

import (
	"../../../../util"
	"../../../../model"
	"../../../http"
	"encoding/json"
)

func RequestLogin(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	body := http.HttpBody {
		Payload: map[string]interface{}{
			"code": -1,
			"message": "youkoso",
		},
	}

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
	user.Find(
		username.(string),
		password.(string),
	)

	return &body, &http.HttpHeader{
		Status: 404,
	}
}
