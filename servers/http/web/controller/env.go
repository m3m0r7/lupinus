package controller

import (
	"encoding/json"
	"lupinus/servers/http"
	"lupinus/share"
)

func RequestEnv(clientMeta http.HttpClientMeta)  (*http.HttpBody, *http.HttpHeader) {
	if clientMeta.Method != "PUT" {
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

	data := map[string]interface{}{}
	json.Unmarshal(clientMeta.Payload, &data)

	share.SetCameraEnv(
		data["temperature"].(float64),
		data["humidity"].(float64),
		data["cpu_temperature"].(float64),
		data["pressure"].(float64),
	)

	return &http.HttpBody {
		Payload: http.Payload{
			"status": 200,
		},
	},
	&http.HttpHeader{
		Status: 200,
	}
}
