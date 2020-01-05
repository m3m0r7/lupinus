package controller

import (
	"lupinus/client"
	"lupinus/servers/http"
	"lupinus/servers/http/web/behavior"
	"os"
)

func RequestUser(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	session := behavior.GetSignInInfo(clientMeta)
	authKeyHeader, err := client.FindHeaderByKey(clientMeta.Headers, "x-auth-key")
	if err != nil ||
		os.Getenv("AUTH_KEY") != (*authKeyHeader).Value ||
		session == nil {
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
			"status": 200,
			"user": session.Data,
		},
	},
	&http.HttpHeader{
		Status: 200,
	}
}