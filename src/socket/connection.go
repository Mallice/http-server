package socket

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/mallice/http-server/src/http"
)

func readHttpHeader(reader *bufio.Reader) (string, error) {
	var sb strings.Builder
	for header, err := reader.ReadString('\n'); header != "\n" && header != "\r\n"; header, err = reader.ReadString('\n') {
		if err != nil {
			return "", err
		}
		sb.WriteString(fmt.Sprintf("%s", header))
	}
	return sb.String(), nil

}
func HandleConnection(c net.Conn) {
	br := bufio.NewReader(c)

	httpHeader, err := readHttpHeader(br)

	if err != nil {
		log.Printf("Something went wrong - %v", err)
		c.Close()
		return
	}

	request := http.Build(httpHeader)

	// Get Body when content-length is defined in Headers
	if contentLength, ok := request.Headers["Content-Length"]; ok {
		length, err := strconv.Atoi(contentLength)

		if err == nil {
			buf := make([]byte, length)
			io.ReadFull(br, buf)
			request.Body = string(buf)
			log.Printf("Body set to: '%s'", request.Body)
		}

	}

	request.Response().
		Status(200, "OK").
		Header("Content-Type", "text/html;charset=utf-8").
		Body([]byte(fmt.Sprintf("<html><body><h1>%s</h1><h3>Hello World</h3></body></html>", request.Path))).
		Send(c)
	c.Close()
}
