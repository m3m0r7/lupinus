package controller

import (
	"io/ioutil"
	"lupinus/config"
	"lupinus/servers/http"
	"lupinus/servers/http/web/behavior"
	"os"
	"strconv"
	"time"
)

func RequestImage(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
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

	id := clientMeta.Path.Query().Get("id")
	integerId, err := strconv.Atoi(id)

	if err != nil {
		integerId = 0
	}

	imagePath := config.GetRootDir() +
		"/storage/" +
		session.Data["id"].(string) +
		"/" +
		time.Unix(int64(integerId), 0).Format("20060102") +
		"/" +
		id +
		".jpg"

	handle, err := os.Open(imagePath)

	if err != nil {
		// Not exists a session
		return &http.HttpBody{
			Payload: http.Payload{
				"message": "No Image",
			},
		},
		&http.HttpHeader{
			Status: 404,
		}
	}

	data, _ := ioutil.ReadAll(handle)

	// Not exists a session
	return &http.HttpBody{
		RawMode: true,
		Payload: http.Payload{
			"body": string(data),
		},
	},
	&http.HttpHeader{
		ContentType: "image/jpeg",
		Status: 200,
	}
}