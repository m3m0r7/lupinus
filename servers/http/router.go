package http

import (
	"errors"
	"../../util"
)

func connect(clientMeta HttpClientMeta) (*HttpBody, *HttpHeader, error) {
	responseBody := &HttpBody{}
	responseHeader := &HttpHeader{}

	switch clientMeta.Path.Path {
	case "/":
		responseBody, responseHeader = requestRoot(clientMeta)
		break

	case "/favicon.ico":
		writeData := "" +
			clientMeta.Protocol + " " + util.GetStatusCodeWithNameByCode(404) + "\n" +
			"Content-Length: 0\n" +
			"\n"
		clientMeta.Pipe.Write([]byte(writeData))
		return nil, nil, errors.New("")
	default:
		responseBody, responseHeader = requestFallback(clientMeta)
		break
	}

	return responseBody, responseHeader, nil
}

func requestFallback(clientMeta HttpClientMeta) (*HttpBody, *HttpHeader) {
	body := HttpBody {
		Payload: map[string]interface{}{
			"code": -1,
			"message": "No data",
		},
	}
	return &body, &HttpHeader{
		Status: 404,
	}
}

func requestRoot(clientMeta HttpClientMeta)  (*HttpBody, *HttpHeader) {
	body := HttpBody{
		Payload: map[string]interface{}{
			"code": -1,
			"message": ":gopher:",
		},
	}

	return &body, &HttpHeader{}
}