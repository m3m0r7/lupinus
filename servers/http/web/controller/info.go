package controller

import (
	"lupinus/servers/http"
	"lupinus/servers/http/web/behavior"
)

func RequestInfo(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	session := behavior.GetSignInInfo(clientMeta)

	if session == nil {
		// Not exists a session
		return &http.HttpBody{
				Payload: http.Payload{
					"message": "Unauthorized",
				},
			},
			&http.HttpHeader{
				Status: 401,
			}
	}

	temperature := 0.0
	humidity := 0.0
	pressure := 0.0
	cpuTemperature := 0.0

	return &http.HttpBody{
			Payload: http.Payload{
				"status": 200,
				"info": map[string]interface{}{
					"temperature":     temperature,
					"humidity":        humidity,
					"pressure":        pressure,
					"cpu_temperature": cpuTemperature,
				},
				"versions": map[string]interface{}{
					"device": map[string]interface{}{
						"number": "0.0.0",
						"code":   "Lupinus",
						"extra":  "Raspibian",
					},
					"app": map[string]interface{}{
						"number": "0.0.0",
						"code":   "Lupinus",
					},
					"live_streaming": map[string]interface{}{
						"number": "0.0.0",
						"code":   "Lupinus",
					},
				},
			},
		},
		&http.HttpHeader{
			Status: 200,
		}
}
