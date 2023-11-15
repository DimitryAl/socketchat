package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var AllConns = make([]net.Conn, 0)
var Mfile *os.File

func closeConnection(conn net.Conn) {
	conn.Close() // Close connection
	remove := func() {
		for i, v := range AllConns {
			if v == conn {
				AllConns = append(AllConns[:i], AllConns[i+1:]...)
			}
		}
	}
	remove()
	fmt.Println("Closing connection: ", conn.RemoteAddr().String())
}

func communication(conn net.Conn) {
	defer closeConnection(conn)

	buf := make([]byte, 32) // Bufer for client message

	// Get client's name
	greeting := "Please, enter your nickname"
	conn.Write([]byte("Server: " + greeting))
	conn.Read(buf)
	nickname := string(buf)

	conn.Write([]byte("Server: " + nickname + " joined to chat!"))
	Broadcast(conn, "Server", []byte(nickname+" joined to chat!"))
	Mfile.Write([]byte(nickname + " joined to chat!"))

	for {

		buf = make([]byte, 32)
		_, err := conn.Read(buf)
		if err != nil {
			break
		}

		Mfile.Write(buf)

		if len(AllConns) > 1 {
			Broadcast(conn, nickname, buf)
		}
	}
}

func Broadcast(conn net.Conn, nickname string, buf []byte) {
	for _, c := range AllConns {
		if c != conn {
			c.Write(append([]byte(nickname+": "), buf...))
		}
	}
}

func main() {

	fmt.Println("Start")

	listener, err := net.Listen("tcp", "localhost:8080") // Open listener socket
	if err != nil {
		log.Fatal("error occured:", err)
	}

	Mfile, err := os.OpenFile("server_data.txt", os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer Mfile.Close()

	for {
		conn, err := listener.Accept() // Accept TCP-connection from client

		if err != nil {
			continue
		}

		fmt.Println(conn.RemoteAddr())
		AllConns = append(AllConns, conn)

		go communication(conn)
	}

}
