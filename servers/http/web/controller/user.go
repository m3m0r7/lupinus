package controller

import (
	"../../../http"
	"../behavior"
)

func RequestUser(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	session := behavior.GetSignInInfo(clientMeta)

	if session == nil {
		// Not exists a session
		return &http.HttpBody{
			Payload: map[string]interface{}{
				"message": "Unauthorized",
			},
		},
		&http.HttpHeader{
			Status: 401,
		}
	}

	return &http.HttpBody{
		Payload: map[string]interface{}{
			"user": session.Data,
		},
	},
	&http.HttpHeader{
		Status: 200,
	}
}