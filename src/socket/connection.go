package socket

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
)

func HandleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	bytes, err := ioutil.ReadAll(bufio.NewReader(c))

	if err != nil {
		fmt.Println(err)
		c.Close()
		return
	}

	fmt.Printf("%s", bytes)

	c.Close()
}
