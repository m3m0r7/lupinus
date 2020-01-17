package controller

import "lupinus/servers/http"

func RequestSignout(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	http.DestroySession(clientMeta)
	return &http.HttpBody{
			Payload: http.Payload{
				"status": 200,
			},
		},
		&http.HttpHeader{
			Status: 200,
		}
}
