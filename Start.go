package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var NODE Node = Node{}

func removeSelf(selfPort string, backend []string) []string {
	res := []string{}

	for _, port := range backend {
		if selfPort != port[1:] {
			res = append(res, port)
		}
	}

	return res
}

func convertLogToString() string {
	res := ""

	for _, command := range NODE.CommitLog {
		res += (command + "+")
	}

	if len(res) > 0 {
		res = res[:len(res)-1] + ";"
	}
	return res
}

func handleConnection(conn net.Conn) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		data := strings.Split(message, "\n")
		args := strings.Split(data[0], "-")

		if args[0] == "GET" {
			if val, ok := NODE.Store[args[1]]; ok {
				conn.Write([]byte(val + "\n"))
			} else {
				conn.Write([]byte("DNE; \n"))
			}
		} else if args[0] == "DELETE" {
			if _, ok := NODE.Store[args[1]]; ok {
				delete(NODE.Store, args[1])
				NODE.CommitLog = append(NODE.CommitLog, "DELETE-"+string(args[1]))
				syncBackends(args)
				conn.Write([]byte("SUCCESS; \n"))
			} else {
				conn.Write([]byte("DNE; \n"))
			}
		} else if args[0] == "INSERT" {
			NODE.Store[args[1]] = args[2]
			NODE.CommitLog = append(NODE.CommitLog, "INSERT-"+string(args[1])+"-"+string(args[2]))
			syncBackends(args)
			conn.Write([]byte("SUCCESS; \n"))
		} else if args[0] == "REQUEST" {
			// Return list of commands up to this point
			conn.Write([]byte(convertLogToString()))
		} else {
			conn.Write([]byte("ERROR; \n"))
		}
	}
}

func syncBackends(args []string) {
	for i := 0; i < len(NODE.ReplicationAddresses); i++ {
		conn, err := net.Dial("tcp", NODE.ReplicationAddresses[i])

		if err != nil {
			continue
		}

		if args[0] == "INSERT" {
			conn.Write([]byte("INSERT-" + string(args[1]) + "-" + string(args[2]) + "\n"))
		} else if args[0] == "DELETE" {
			conn.Write([]byte("DELETE-" + string(args[1]) + "\n"))
		}
	}
}

func executeCommand(command string) {
	args := strings.Split(string(command), "-")

	if args[0] == "INSERT" {
		NODE.Store[args[1]] = args[2]
		NODE.CommitLog = append(NODE.CommitLog, "INSERT-"+string(args[1])+"-"+string(args[2]))
	} else if args[1] == "DELETE" {
		if _, ok := NODE.Store[args[1]]; ok {
			delete(NODE.Store, args[1])
			NODE.CommitLog = append(NODE.CommitLog, "DELETE-"+string(args[1]))
		}
	}
}

func handleSync(ls string) {
	// Reset
	NODE.Store = make(map[string]string)
	NODE.CommitLog = make([]string, 0)

	// Split by + for each log -> Grab the original and execute logs in that order to sync
	l := strings.Split(ls, "+")

	for i := 0; i < len(l); i++ {
		executeCommand(l[i])
	}

	fmt.Println("Sync complete")
}

func requestSync() {
	for i := 0; i < len(NODE.ReplicationAddresses); i++ {
		conn, err := net.Dial("tcp", NODE.ReplicationAddresses[i])

		if err != nil {
			continue
		}

		conn.Write([]byte("REQUEST-\n"))

		message, _ := bufio.NewReader(conn).ReadString('\n')
		data := strings.Split(message, "\n")
		args := strings.Split(data[0], ";")

		n, _ := strconv.Atoi(args[0])
		if n > len(NODE.CommitLog) {
			handleSync(args[1])
		}
	}
}

func main() {
	fmt.Println("Starting Cassandra")

	port := "8090"

	// --listen
	port = os.Args[2]

	backends := removeSelf(port, strings.Split(os.Args[4], ","))
	listener, err := net.Listen("tcp", ":"+port)

	// Initialize node
	NODE.PartitionKey = 0 // TODO: Implement hash
	NODE.Store = make(map[string]string)
	NODE.ReplicationAddresses = backends
	NODE.CommitLog = make([]string, 0)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer listener.Close()
	requestSync()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		// In case of restart after crash, make sure to sync the logs
		handleConnection(conn)
	}
}
