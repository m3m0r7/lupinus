package util

func GetStatusCodeWithNameByCode(code int) string {
	// Set status code
	switch code {
	case 400:
		return "400 Bad Request"
	case 401:
		return "401 Unauthorized"
	case 404:
		return "404 Not Found"
	case 403:
		return "403 Forbidden"
	case 500:
		return "500 Internal Server Error"
	}
	return "200 OK"
}
