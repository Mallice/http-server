package http

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

type METHOD string

const (
	GET     METHOD = "GET"
	HEAD           = "HEAD"
	POST           = "POST"
	PUT            = "PUT"
	DELETE         = "DELETE"
	CONNECT        = "CONNECT"
	OPTIONS        = "OPTIONS"
	TRACE          = "TRACE"
	PATCH          = "PATCH"
)

var Methods = []METHOD{GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH}

type Request struct {
	protocol string
	Method   METHOD
	Path     string
	Headers  map[string]string
	Body     string
}

func (r *Request) parseStartLine(br *bufio.Reader) error {
	startLine, err := br.ReadString('\n')
	if err != nil {
		return err
	}

	splitted := strings.Split(startLine, " ")
	if len(splitted) != 3 {
		return fmt.Errorf("Invalid StartLine")
	}

	for _, method := range Methods {
		if string(method) == splitted[0] {
			r.Method = method
			r.Path = splitted[1]
			r.protocol = splitted[2]
			return nil
		}
	}

	return fmt.Errorf("No HTTP StartLine found")
}

func (r *Request) parseHeaders(br *bufio.Reader) error {
	var sb strings.Builder
	for header, err := br.ReadString('\n'); header != "\n" && header != "\r\n"; header, err = br.ReadString('\n') {
		if err != nil {
			return err
		}
		sb.WriteString(fmt.Sprintf("%s", header))
	}

	headers := strings.Split(sb.String(), "\n")
	for _, header := range headers {
		headerKV := strings.SplitN(header, ":", 2)

		switch len(headerKV) {
		case 2:
			headerKV[0] = strings.TrimSpace(headerKV[0])
			headerKV[1] = strings.TrimSpace(headerKV[1])
			r.Headers[headerKV[0]] = headerKV[1]
			break
		case 1:
			headerKV[0] = strings.TrimSpace(headerKV[0])
			if len(headerKV[0]) > 0 {
				r.Headers[headerKV[0]] = ""
			}
			break
		case 0:
			log.Printf("Warn: Unparsed Header - %v", fmt.Errorf("':' nowhere on '%s'", header))
		}
	}
	return nil
}

func (r *Request) parseBody(br *bufio.Reader) error {
	// Get Body when content-length is defined in Headers
	if contentLength, ok := r.Headers["Content-Length"]; ok {
		length, err := strconv.Atoi(contentLength)

		if err != nil {
			return err
		}

		buffer := make([]byte, length)
		_, err = io.ReadFull(br, buffer)

		if err != nil {
			return err
		}
		r.Body = string(buffer)
	}
	return nil
}

func ParseRequest(bf *bufio.Reader) (Request, error) {
	request := Request{Body: "", Headers: map[string]string{}}

	request.parseStartLine(bf)
	request.parseHeaders(bf)
	request.parseBody(bf)
	return request, nil
}

func (r Request) Response(status int, statusMessage string) *Response {
	response := Response{status: status, statusMessage: statusMessage, protocol: "HTTP/1.0", headers: map[string]string{}, body: ""}

	response.Header("Server", "Custom")
	response.Header("Date", time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	response.Header("Access-Control-Allow-Origin", "*")
	return &response
}
