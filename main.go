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
	conn.Write([]byte(greeting))
	conn.Read(buf)
	nickname := string(buf)
	conn.Write([]byte(nickname + " joined to chat!"))
	Mfile.Write([]byte(nickname + " joined to chat!"))

	for {
		// check if connection alive?

		conn.Read(buf)

		Mfile.Write(buf)
		if len(AllConns) > 1 {
			r := Broadcast(conn, nickname, buf)
			if r == 1 {
				break
			}
		}
	}
}

func Broadcast(conn net.Conn, nickname string, buf []byte) (r int) {
	fmt.Println("This is Broadcasting")
	for _, c := range AllConns {
		if c != conn {
			//fmt.Println("Trying to write to", c.RemoteAddr().String())
			_, err := c.Write(append([]byte(nickname+": "), buf...))
			if err != nil {
				return 1
			}
		}
	}
	return 0
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
