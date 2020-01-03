package controller

import "lupinus/servers/http"

func RequestFallback(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	body := http.HttpBody {
		Payload: http.Payload{
			"code": -1,
			"message": "No data",
		},
	}
	return &body, &http.HttpHeader{
		Status: 404,
	}
}