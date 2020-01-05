package controller

import (
	"lupinus/servers/http"
)

func RequestRoot(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	return &http.HttpBody{
			Payload: http.Payload{
				"message": "(=^・_・^=)",
			},
		},
		&http.HttpHeader{}
}
