package web

import (
	"../../http"
	"./controller"
)

func Connect(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader, error) {
	responseBody := &http.HttpBody{}
	responseHeader := &http.HttpHeader{}

	if clientMeta.Method != "OPTION" {
		switch clientMeta.Path.Path {
		case "/":
			responseBody, responseHeader = controller.RequestRoot(clientMeta)
			break
		case "/api/v1/signin":
			responseBody, responseHeader = controller.RequestLogin(clientMeta)
			break
		case "/api/v1/user":
			// TODO: implement here
			break
		case "/image":
			// TODO: implement here
			break
		case "/info":
			// TODO: implement here
			break
		case "/favorite":
			// TODO: implement here
			break
		case "/capture":
			// TODO: implement here
			break
		case "/favicon.ico":
			return nil, nil, nil
		default:
			responseBody, responseHeader = controller.RequestFallback(clientMeta)
			break
		}
	}

	return responseBody, responseHeader, nil
}
