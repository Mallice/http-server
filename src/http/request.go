package http

import (
	"fmt"
	"log"
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
	Protocol string
	Method   METHOD
	Path     string
	Headers  map[string]string
	Body     string
}

func Build(httpHeader string) Request {
	request := Request{Body: "", Headers: map[string]string{}}

	splitted := strings.SplitN(httpHeader, "\n\n", 2)

	head := strings.Split(splitted[0], "\n")

	method, path, protocol, err := readHttpStartLine(head[0])

	if err != nil {
		log.Panicf("%e", err)
	}

	request.Method = *method
	request.Path = *path
	request.Protocol = *protocol
	for _, header := range head[1:] {
		headerKV := strings.SplitN(header, ":", 2)

		switch len(headerKV) {
		case 2:
			headerKV[0] = strings.TrimSpace(headerKV[0])
			headerKV[1] = strings.TrimSpace(headerKV[1])
			request.Headers[headerKV[0]] = headerKV[1]
			break
		case 1:
			headerKV[0] = strings.TrimSpace(headerKV[0])
			if len(headerKV[0]) > 0 {
				request.Headers[headerKV[0]] = ""
			}
			break
		case 0:
			log.Printf("Unparsed Header: - %v", fmt.Errorf("':' nowhere on '%s'", header))
		}
	}

	if len(splitted) == 2 {
		request.Body = splitted[1]
	}

	return request
}

func readHttpStartLine(startLine string) (*METHOD, *string, *string, error) {
	splitted := strings.Split(startLine, " ")
	if len(splitted) != 3 {
		return nil, nil, nil, fmt.Errorf("Invalid Start Line")
	}

	for _, method := range Methods {
		if string(method) == splitted[0] {

			return &method, &splitted[1], &splitted[2], nil
		}
	}
	return nil, nil, nil, fmt.Errorf("Unknown HTTP Method '%s'", splitted[0])

}

func (r Request) Response() *Response {
	response := Response{status: -1, protocol: "HTTP/1.0", headers: map[string]string{}, body: ""}

	response.Header("Server", "Custom")
	response.Header("Date", time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	response.Header("Access-Control-Allow-Origin", "*")
	return &response
}
