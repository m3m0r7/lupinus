package web

import (
	"../../../util"
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
		writeData := "" +
			clientMeta.Protocol + " " + util.GetStatusCodeWithNameByCode(404) + "\n" +
			"Content-Length: 0\n" +
			"\n"
		clientMeta.Pipe.Write([]byte(writeData))
		return nil, nil, nil
	default:
		responseBody, responseHeader = controller.RequestFallback(clientMeta)
		break
	}

	return responseBody, responseHeader, nil
}
