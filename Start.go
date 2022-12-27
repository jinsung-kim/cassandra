package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var backends []string // Used to communicate with other backends

func removeSelf(selfPort string, backend []string) []string {
	res := []string{}

	for _, port := range backend {
		fmt.Println(port[1:])
		if selfPort != port[1:] {
			res = append(res, port)
		}
	}

	return res
}

func main() {
	fmt.Println("Starting Cassandra")

	port := "8090"

	// --listen
	port = os.Args[2]

	backends := removeSelf(port, strings.Split(os.Args[4], ","))
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(backends)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		// In case of restart after crash, make sure to sync the logs
		fmt.Println(conn)
	}
}
