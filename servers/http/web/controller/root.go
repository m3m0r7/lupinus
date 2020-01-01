package controller

import "../../../http"

func RequestRoot(clientMeta http.HttpClientMeta)  (*http.HttpBody, *http.HttpHeader) {
	body := http.HttpBody{
		Payload: map[string]interface{}{
			"code": -1,
			"message": ":gopher:",
		},
	}

	return &body, &http.HttpHeader{}
}