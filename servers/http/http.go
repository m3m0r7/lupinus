package http

import (
	"net"
	"net/url"
)

type ClientHeader struct {
	Key   string
	Value string
}

type Client struct {
	Pipe    net.Conn
	Headers []ClientHeader
}

type HttpClientMeta struct {
	Protocol string
	Pipe     *net.Conn
	Method   string
	Headers  []ClientHeader
	Path     url.URL
	Payload  []byte
	Cookies  []Cookie
}

type HttpBody struct {
	RawMode bool
	Payload map[string]interface{}
	AfterProcess func(meta HttpClientMeta)
}

type HttpHeader struct {
	Status      int
	ContentType string
	Headers     []ClientHeader
}

type ErrorMessage struct {
	Code    int
	Message string
}

type Payload map[string]interface{}
