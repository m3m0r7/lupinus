package controller

import "../../../http"

func RequestLogin(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	body := http.HttpBody {
		Payload: map[string]interface{}{
			"code": -1,
			"message": "No data",
		},
	}
	return &body, &http.HttpHeader{
		Status: 404,
	}
}
