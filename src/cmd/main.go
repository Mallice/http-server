package main

import (
	_ "embed"
	"fmt"
	"net"
	"os"

	"github.com/mallice/http-server/src/socket"
)

//go:embed usage
var usage string

// program entrypoint
// expect only one argument (the input file)
// Print usage.txt in other cases
// panic: if input file is not a readable file
func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		println(usage)
		os.Exit(0)
	}

	port := args[0]
	filepath := args[1]

	fmt.Printf("Serving '%s' at localhost:%s\n", filepath, port)

	l, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	//defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go socket.HandleConnection(c)
	}
}
