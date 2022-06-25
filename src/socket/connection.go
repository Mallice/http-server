package socket

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func HandleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}

		result := "5\n"
		c.Write([]byte(string(result)))
	}
	c.Close()
}
