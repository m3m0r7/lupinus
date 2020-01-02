package http

import (
	"strconv"
	"strings"
)

type Cookie struct {
	Name string
	Path string
	Value string
	Domain string
	MaxAge int
	HttpOnly bool
}

var cookies = []string{}

func AddCookie(cookie Cookie) {
	tmp := []string{}
	tmp = append(tmp, cookie.Name + "=" + cookie.Value)
	tmp = append(tmp, "Path=" + cookie.Path)

	if cookie.Domain != "" {
		tmp = append(tmp, "Domain="+cookie.Domain)
	}

	tmp = append(tmp, "Expires=" + strconv.Itoa(cookie.MaxAge))
	cookies = append(cookies, strings.Join(tmp, "; ") + ";")
}

func GetCookies() []string {
	return cookies
}