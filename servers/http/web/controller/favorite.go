package controller

import (
	"../../../http"
	"../../../../share"
)

func RequestFavorite(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
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
	return &http.HttpBody{
		Payload: map[string]interface{}{
			"status": 200,
			"dates": map[string]interface{}{
				"20190105": []map[string]interface{}{
					{
						"src": "hoghoge",
					},
					{
						"src": "hoghoge",
					},
				},
			},
		},
	},
	&http.HttpHeader{
		Status: 200,
	}
}

func requestFavoriteByPost(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	share.AddProcedure(share.Procedure{
		Callback: func(data string) {

		},
	})
	return &http.HttpBody{
		Payload: map[string]interface{}{
			"status": 200,
		},
	},
	&http.HttpHeader{
		Status: 200,
	}
}
