package controller

import (
	"../../../http"
	"../../../../model"
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



	user := model.InitUser()
	(*user).Find()

	return &body, &http.HttpHeader{
		Status: 404,
	}
}
