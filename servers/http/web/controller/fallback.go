package controller

import "lupinus/servers/http"

func RequestFallback(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	return &http.HttpBody{
			Payload: http.Payload{
				"status": 500,
				"error":  "No Data",
			},
		},
		&http.HttpHeader{
			Status: 404,
		}
}
