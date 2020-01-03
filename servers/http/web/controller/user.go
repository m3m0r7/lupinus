package controller

import (
	"lupinus/servers/http/web/behavior"
	"lupinus/servers/http"
)

func RequestUser(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	session := behavior.GetSignInInfo(clientMeta)

	if session == nil {
		// Not exists a session
		return &http.HttpBody{
			Payload: http.Payload{
				"status": 500,
				"error": "Unauthorized",
			},
		},
		&http.HttpHeader{
			Status: 401,
		}
	}

	return &http.HttpBody{
		Payload: http.Payload{
			"user": session.Data,
		},
	},
	&http.HttpHeader{
		Status: 200,
	}
}