package http

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"../../client"
	"../../util"
	"strings"
)

func Listen() {
	listener, _ := net.Listen(
		"tcp",
		os.Getenv("CLIENT_API_SERVER"),
	)
	fmt.Printf("Start client API server %v\n", listener.Addr())
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to listen. retry again.")
			continue
		}

		go func() {
			// Get headers
			headers, _ := client.GetAllHeaders(connection)
			result, _ := client.FindHeaderByKey(&headers, "")
			status := util.SplitWithFiltered(result.Value, " ")

			if len(status) != 3 {
				fmt.Printf("Invalid header")
				return
			}


			method := strings.ToLower(status[0])
			path := status[1]
			rule := status[2]

			urlObject, parseError := url.Parse(path)

			if parseError != nil {
				fmt.Printf("Invalid header")
				return
			}

			clientMeta := HttpClientMeta{
				Pipe: connection,
				Method: method,
				Path: *urlObject,
			}

			responseBody := &HttpBody{}
			responseHeader := &HttpHeader{}

			switch urlObject.Path {
			//case "/":
				//responseBody, responseHeader = controller.RequestRoot(clientMeta)
				//break
			default:
				responseBody, responseHeader = requestFallback(clientMeta)
				break
			}

			resultJSON, _ := json.Marshal(responseBody.Payload)
			_ = resultJSON
			_ = responseHeader

			stringifiedJSON := string(resultJSON)
			statusWithName := "200 OK"

			// Set status code
			if responseHeader.Status == 400 {
				statusWithName = "400 Bad Request"
			}

			if responseHeader.Status == 404 {
				statusWithName = "404 Not Found"
			}

			if responseHeader.Status == 403 {
				statusWithName = "403 Forbidden"
			}

			if responseHeader.Status == 500 {
				statusWithName = "500 Internal Server Error"
			}

			// Write buffer
			writeData := "" +
				rule + " " + statusWithName + "\n" +
				"Content-Length: " + string(len(stringifiedJSON)) + "\n" +
				"Content-Type: application/json\n" +
				"\n" +
				stringifiedJSON

			connection.Write([]byte(writeData))
			connection.Close()
		}()
	}
}


func requestFallback(client HttpClientMeta) (*HttpBody, *HttpHeader) {
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