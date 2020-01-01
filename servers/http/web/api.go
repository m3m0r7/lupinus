package web

import (
	"../../../client"
	"../../../util"
	"../../http"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
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
			defer connection.Close()

			fmt.Printf("Connected from: %v\n", connection.RemoteAddr())

			// Get headers
			headers, _ := client.GetAllHeaders(connection)
			result, err := client.FindHeaderByKey(headers, "")

			if err != nil {
				return
			}

			status := util.SplitWithFiltered(result.Value, " ")
			if len(status) != 3 {
				return
			}

			method := strings.ToLower(status[0])
			path := status[1]
			protocol := status[2]

			// Next find body size
			result, err = client.FindHeaderByKey(headers, "content-length")

			requestBody := []byte{}
			if result != nil {
				number, _ := strconv.Atoi(result.Value)
				reads := make([]byte, number)
				connection.Read(reads)
				requestBody = append(requestBody, reads...)
			}

			urlObject, parseError := url.Parse(path)

			if parseError != nil {
				fmt.Printf("Invalid header")
				return
			}

			// TODO: Validator

			clientMeta := http.HttpClientMeta{
				Pipe: connection,
				Method: strings.ToUpper(method),
				Path: *urlObject,
				Protocol: protocol,
				Payload: requestBody,
			}

			responseBody, responseHeader, _ := Connect(clientMeta)

			// If which is invalid processing, connection to shutdown.
			if responseBody == nil || responseHeader == nil {
				writeData := "" +
					clientMeta.Protocol + " " + util.GetStatusCodeWithNameByCode(404) + "\n" +
					"Content-Length: 0\n" +
					"Connection: close\n" +
					"\n"
				clientMeta.Pipe.Write([]byte(writeData))
				return
			}

			resultJSON, _ := json.Marshal(responseBody.Payload)

			stringifiedJSON := string(resultJSON)
			statusWithName := util.GetStatusCodeWithNameByCode(
				responseHeader.Status,
			)
			// Write buffer
			writeData := "" +
				clientMeta.Protocol + " " + statusWithName + "\n" +
				"Content-Length: " + strconv.Itoa(len(stringifiedJSON)) + "\n" +
				"Content-Type: application/json\n" +
				"Connection: close" +
				"\n" +
				stringifiedJSON

			connection.Write([]byte(writeData))
		}()
	}
}
