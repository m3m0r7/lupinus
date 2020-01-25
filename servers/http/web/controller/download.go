package controller

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"lupinus/client"
	"lupinus/config"
	"lupinus/servers/http"
	"lupinus/servers/http/web/behavior"
	"os"
	"path/filepath"
)

func RequestDownload(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {

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

	date := clientMeta.Path.Query().Get("date")
	rootDirectoryName := date
	findPath := date
	if date == "" {
		rootDirectoryName = "all"
		findPath = "*"
	}

	bytes := new(bytes.Buffer)
	writer := zip.NewWriter(bytes)

	files, _ := filepath.Glob(
		config.GetRootDir() + "/storage/" + session.Data["id"].(string) + "/" + findPath + "/*.jpg",
	)

	for _, file := range files {
		info, _ := os.Stat(file)
		header, _ := zip.FileInfoHeader(info)
		header.Name = rootDirectoryName + "/" + filepath.Base(file)
		handle, _ := writer.CreateHeader(header)

		data, _ := ioutil.ReadFile(file)
		handle.Write(data)
	}

	writer.Close()

	return &http.HttpBody{
			RawMode: true,
			Payload: http.Payload{
				"body": string(bytes.Bytes()),
			},
		},
		&http.HttpHeader{
			Status:      200,
			ContentType: "application/zip",
		}
}
