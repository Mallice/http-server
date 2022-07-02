package http

import (
	"fmt"
	"net"
	"strings"
	"unicode/utf8"
)

type Response struct {
	protocol      string
	status        int
	statusMessage string
	headers       map[string]string
	body          string
}

func (r *Response) Header(key string, value string) *Response {
	if currentValue, ok := r.headers[key]; ok {
		r.headers[key] = fmt.Sprintf("%s, %s", currentValue, value)
	} else {
		r.headers[key] = value
	}
	return r
}

func (r *Response) Body(bytes []byte) *Response {
	r.body = string(bytes)
	return r
}

func (r *Response) Send(c net.Conn) error {
	if r.status == -1 {
		return fmt.Errorf("An HTTP Response need to have a status")
	}

	r.Header("Content-Length", fmt.Sprint(utf8.RuneCountInString(r.body)))
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %d %s\r\n", r.protocol, r.status, r.statusMessage))
	for key, value := range r.headers {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	sb.WriteString("\r\n\r\n")
	sb.WriteString(r.body)
	sb.WriteString("\r\n")
	c.Write([]byte(sb.String()))
	return nil
}
