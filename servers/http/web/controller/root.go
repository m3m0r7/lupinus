package controller

import (
	"lupinus/servers/http"
)

func RequestRoot(clientMeta http.HttpClientMeta)  (*http.HttpBody, *http.HttpHeader) {
	return &http.HttpBody{
		Payload: map[string]interface{}{
			"message": "(=^・_・^=)",
		},
	},
	&http.HttpHeader{}
}