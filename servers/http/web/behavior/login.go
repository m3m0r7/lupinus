package behavior

import (
	"../../../http"
)

func GetSignInInfo(clientMeta http.HttpClientMeta) *http.Session {
	session := http.LoadSessionFromCookie(clientMeta.Cookies)

	if session == nil {
		// Not exists a session
		return nil
	}
	if _, isExist := session.Data["id"]; !isExist {
		// Not exists a session
		return nil
	}

	return session
}