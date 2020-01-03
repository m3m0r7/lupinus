package controller

import "lupinus/servers/http"

func RequestFallback(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	return &http.HttpBody {
		Payload: http.Payload{
			"code": -1,
			"message": "No data",
		},
	},
	&http.HttpHeader{
		Status: 404,
	}
}