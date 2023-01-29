package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()
	cache := make(map[string]string)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		fmt.Printf("Connection from %s opened.\n", conn.RemoteAddr())
		go handle(conn, cache)
	}
}

func handle(conn net.Conn, cache map[string]string) {
	for {
		message := make([]byte, 128)
		n, err := conn.Read(message)
		if err != nil {
			switch err.Error() {
			case "EOF":
				fmt.Printf("Connection from %s closed.\n", conn.RemoteAddr())
			default:
				fmt.Println("Error reading: ", err.Error())
			}
			return
		}
		if n != 0 {
			s := string(message[:])
			req := strings.Split(s, "\r\n")
			fmt.Printf("%s >> %v\n", conn.RemoteAddr(), req)
			command := strings.ToLower(req[2])
			switch command {
			case "echo":
				msg := req[4]
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)))
			case "set":
				key := req[4]
				value := req[6]
				cache[key] = value
				msg := "OK"
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)))
			default:
				conn.Write([]byte("+PONG\r\n"))
			}
			fmt.Printf("%+v\n", cache)
		}
	}
}
