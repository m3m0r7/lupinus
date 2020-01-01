package http

import (
	"../../client"
	"net"
	"net/url"
)

type HttpClientMeta struct {
	Pipe net.Conn
	Method string
	Headers []client.ClientHeader
	Path url.URL
}

type HttpBody struct {
	Payload map[string]interface{}
}

type HttpHeader struct {
	Status int
	Headers []client.ClientHeader
}

type ErrorMessage struct {
	Code int
	Message string
}