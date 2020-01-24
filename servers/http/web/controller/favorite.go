package controller

import (
  "github.com/nlopes/slack"
  "lupinus/client"
  "lupinus/config"
  "lupinus/helper"
  "lupinus/servers/http"
  "lupinus/servers/http/web/behavior"
  "lupinus/share"
  "os"
  "path/filepath"
  "strconv"
  "strings"
  "time"
)

func RequestFavorite(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	session := behavior.GetSignInInfo(clientMeta)

	authKeyHeader, err := client.FindHeaderByKey(clientMeta.Headers, "x-auth-key")
	if err != nil ||
		os.Getenv("AUTH_KEY") != (*authKeyHeader).Value ||
		session == nil {
		// Not exists a session
		return &http.HttpBody{
				Payload: http.Payload{
					"status": 401,
					"error":  "Failed to authorize",
				},
			},
			&http.HttpHeader{
				Status: 401,
			}
	}

	body := &http.HttpBody{}
	header := &http.HttpHeader{}
	switch clientMeta.Method {
	case "GET":
		body, header = requestFavoriteByGet(clientMeta)
		break
	case "POST":
		body, header = requestFavoriteByPost(clientMeta)
		break
	default:
		body = nil
		header = nil
	}
	return body, header
}

func requestFavoriteByGet(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	session := behavior.GetSignInInfo(clientMeta)

	files, _ := filepath.Glob(
		config.GetRootDir() + "/storage/" + session.Data["id"].(string) + "/*/*.jpg",
	)

	dates := map[string]interface{}{}

	for _, file := range files {
		unixTime, _ := strconv.Atoi(
			strings.Replace(
				filepath.Base(file),
				".jpg",
				"",
				-1),
		)

		date := time.Unix(int64(unixTime), 0)
		id := date.Format("20060102")

		type dateObject = map[string]interface{}
		if _, ok := dates[id]; !ok {
			dates[id] = []dateObject{}
		}
		dates[id] = append(
			dates[id].([]dateObject),
			dateObject{
				"src": "/api/v1/image?id=" + strconv.Itoa(unixTime),
			},
		)
	}

	return &http.HttpBody{
			Payload: http.Payload{
				"status": 200,
				"dates":  dates,
			},
		},
		&http.HttpHeader{
			Status: 200,
		}
}

func requestFavoriteByPost(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	share.AddProcedure(
		"favorite",
		share.Procedure{
			Callback: func(data []byte) {
				session := behavior.GetSignInInfo(clientMeta)
				id := session.Data["id"].(string)

				path := id + "/" + time.Now().Format("20060102") + "/" + strconv.Itoa(int(time.Now().Unix())) + ".jpg"

				filePath, _ := helper.CreateStaticImage(
					data,
					path,
				)

				fileHandle, fileOpenErr := os.Open(*filePath)

				if fileOpenErr != nil {
				  return
        }

        var slackApi = slack.New(os.Getenv("SLACK_TOKEN"))

        _, _, err := slackApi.PostMessage(
         os.Getenv("SLACK_CHANNEL"),
         slack.MsgOptionText("Favorite now :heart:", false),
        )

        _, err = slackApi.UploadFile(
					slack.FileUploadParameters{
						Reader:   fileHandle,
						Filename: id,
						Channels: []string{
						  os.Getenv("SLACK_CHANNEL"),
						},
					},
				)

				_ = err

			},
		},
	)
	return &http.HttpBody{
			Payload: http.Payload{
				"status": 200,
			},
		},
		&http.HttpHeader{
			Status: 200,
		}
}
