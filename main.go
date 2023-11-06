package main

import (
	"fmt"
	//"github.com/DimitryAl/SocketChat/client"
	"log"
	"net"
)

func readClient(conn net.Conn) {
	defer conn.Close() // Close connection

	buf := make([]byte, 32) // Bufer for client message

	for {
		readMessage, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err")
			break
		}
		fmt.Println(readMessage)
	}

}

func main() {

	fmt.Println("Start")

	listener, err := net.Listen("tcp", "localhost:8080") // Open listener socket

	if err != nil {
		log.Fatal("error occured:", err)
	}

	for {
		conn, err := listener.Accept() // Accept TCP-connection from client
		if err != nil {
			continue
		}

		go readClient(conn)
	}

}
