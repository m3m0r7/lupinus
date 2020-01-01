package controller

import (
	"../../../http"
	"fmt"
	"io/ioutil"
	"os"
)

func RequestLogin(clientMeta http.HttpClientMeta) (*http.HttpBody, *http.HttpHeader) {
	body := http.HttpBody {
		Payload: map[string]interface{}{
			"code": -1,
			"message": "youkoso",
		},
	}

	// Read user file
	dir, _ := os.Getwd()
	userJsonPathHandler, _ := os.Open(dir + "/users.json")
	userData, _ := ioutil.ReadAll(userJsonPathHandler)

	fmt.Printf("%v\n", string(userData))

	return &body, &http.HttpHeader{
		Status: 404,
	}
}
