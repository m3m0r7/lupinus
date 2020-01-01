package controller

import (
	"../../../http"
	"../../../../model"
	"encoding/json"
	"fmt"
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

	dest := []byte{}
	err := json.Unmarshal(dest, clientMeta.Payload)

	if err != nil {
		return nil, nil
	}

	fmt.Printf("%v", dest)

	user := model.InitUser()
	(*user).Find()

	return &body, &http.HttpHeader{
		Status: 404,
	}
}
