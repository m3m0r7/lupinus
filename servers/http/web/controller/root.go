package controller

import (
	"../../../http"
	"fmt"
	"os"
)

func RequestRoot(clientMeta http.HttpClientMeta)  (*http.HttpBody, *http.HttpHeader) {
	body := http.HttpBody{
		Payload: map[string]interface{}{
			"message": "(=^・_・^=)",
		},
	}

	fmt.Printf("%v", http.FindCookie(os.Getenv("SESSION_ID"), clientMeta.Cookies))

	return &body, &http.HttpHeader{}
}