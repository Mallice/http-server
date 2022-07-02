package socket

import (
	"bufio"
	"fmt"
	"net"

	"github.com/mallice/http-server/src/http"
)

func HandleConnection(c net.Conn) {
	br := bufio.NewReader(c)

	request, err := http.ParseRequest(br)

	if err != nil {
		request.Response(500, "Fatal").Header("Content-Type", "text/html;charset=utf-8").Body([]byte(fmt.Sprintf("<html><body><h1>Something went wrong: %s</h1></body></html>", err.Error())))

	}
	request.Response(200, "OK").
		Header("Content-Type", "text/html;charset=utf-8").
		Body([]byte(fmt.Sprintf("<html><body><h1>%s</h1><h3>Hello World</h3></body></html>", request.Path))).
		Send(c)
	c.Close()
}
