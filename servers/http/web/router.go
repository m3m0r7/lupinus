package web

import (
	"../../http"
	"./controller"
)

func Connect(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader, error) {
	responseBody := &http.HttpBody{}
	responseHeader := &http.HttpHeader{}

	switch clientMeta.Path.Path {
	case "/":
		responseBody, responseHeader = controller.RequestRoot(clientMeta)
		break
	case "/api/v1/signin":
		responseBody, responseHeader = controller.RequestLogin(clientMeta)
		break
	case "/favicon.ico":
		return nil, nil, nil
	default:
		responseBody, responseHeader = controller.RequestFallback(clientMeta)
		break
	}

	return responseBody, responseHeader, nil
}
