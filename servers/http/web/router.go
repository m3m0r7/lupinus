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
			responseBody, responseHeader = controller.RequestSignin(clientMeta)
			break
		case "/api/v1/user":
			responseBody, responseHeader = controller.RequestUser(clientMeta)
			break
		case "/api/v1/info":
			responseBody, responseHeader = controller.RequestInfo(clientMeta)
			break
		case "/api/v1/favorite":
			// TODO: implement here
			break
		case "/image":
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
