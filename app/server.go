package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	for {
		message := make([]byte, 128)
		n, err := conn.Read(message)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			os.Exit(1)
		}
		if n != 0 {
			fmt.Println("message: ", n, message)
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
