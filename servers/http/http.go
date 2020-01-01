package http

import (
	"net"
	"net/url"
)

type ClientHeader struct {
	Key string
	Value string
}

type Client struct {
	Pipe net.Conn
	Headers []ClientHeader
}

type HttpClientMeta struct {
	Protocol string
	Pipe net.Conn
	Method string
	Headers []ClientHeader
	Path url.URL
	Payload []byte
}

type HttpBody struct {
	Payload map[string]interface{}
}

type HttpHeader struct {
	Status int
	Headers []ClientHeader
}

type ErrorMessage struct {
	Code int
	Message string
}
