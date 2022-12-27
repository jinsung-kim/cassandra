package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var backends []string // Used to communicate with other backends

func main() {
	fmt.Println("Starting Cassandra")

	port := "8090"

	// --listen
	port = os.Args[2]

	backends := strings.Split(os.Args[4], ",")
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(backends)
	fmt.Println("Listening on TCP")

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
