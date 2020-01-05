package controller

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"lupinus/client"
	"lupinus/config"
	"lupinus/servers/http"
	"lupinus/servers/http/web/behavior"
	"lupinus/servers/streaming"
	"lupinus/util"
	"os"
)

var blackScreenImage []byte

func RequestCapture(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	session := behavior.GetSignInInfo(clientMeta)
	authKeyHeader, err := client.FindHeaderByKey(clientMeta.Headers, "x-auth-key")
	if err != nil ||
		os.Getenv("AUTH_KEY") != (*authKeyHeader).Value ||
		session == nil {
		// Not exists a session
		return &http.HttpBody{
				Payload: http.Payload{
					"status": 500,
					"error":  "Unauthorized",
				},
			},
			&http.HttpHeader{
				Status: 401,
			}
	}

	header := &http.HttpHeader{
		Status: 200,
	}

	// Check file exists
	var capturedImage []byte
	path := config.GetRootDir() + "/storage/record/image.jpg"
	if _, err := os.Stat(path); err == nil {
		handle, _ := os.Open(path)
		capturedImage, _ = ioutil.ReadAll(handle)
	} else {
		if blackScreenImage == nil {
			buffer := bytes.NewBuffer([]byte{})
			filledImage := image.NewRGBA(
				image.Rect(
					0,
					0,
					640,
					480,
				),
			)

			_ = jpeg.Encode(
				buffer,
				filledImage,
				&jpeg.Options{
					100,
				},
			)

			capturedImage = buffer.Bytes()
			blackScreenImage = capturedImage
		} else {
			capturedImage = blackScreenImage
		}
	}

	return &http.HttpBody{
			Payload: http.Payload{
				"status":          200,
				"image":           string(util.Byte2base64URI(capturedImage)),
				"updated_at":      streaming.UpdateTime,
				"update_interval": streaming.UpdateStaticImageInterval,
				"next_update":     streaming.NextUpdateTime,
			},
		},
		header
}
