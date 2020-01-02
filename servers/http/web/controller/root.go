package controller

import "../../../http"

func RequestRoot(clientMeta http.HttpClientMeta)  (*http.HttpBody, *http.HttpHeader) {
	body := http.HttpBody{
		Payload: map[string]interface{}{
			"message": "(=^・_・^=)",
		},
	}

	return &body, &http.HttpHeader{}
}