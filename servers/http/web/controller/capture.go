package controller

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"lupinus/servers/http"
	"lupinus/servers/streaming"
	"lupinus/config"
	"lupinus/util"
	"os"
)

func RequestCapture(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
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
	}

	return &http.HttpBody{
		Payload: map[string]interface{}{
			"image": util.Byte2base64URI(capturedImage),
			"updated_at": camera.UpdateTime,
			"update_interval": camera.UpdateStaticImageInterval,
			"next_update": camera.NextUpdateTime,
		},
	},
	header
}