package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var AllConns = make([]net.Conn, 0)
var Mfile *os.File

func communication(conn net.Conn) {
	defer func() {
		conn.Close() // Close connection
		remove := func() {
			for i, v := range AllConns {
				if v == conn {
					AllConns = append(AllConns[:i], AllConns[i+1:]...)
				}
			}
		}
		remove()
	}()

	buf := make([]byte, 64) // Bufer for client message

	// Get client's name
	greeting := "Please, enter your nickname"
	conn.Write([]byte(greeting))
	conn.Read(buf)
	nickname := string(buf)
	conn.Write([]byte(nickname + " joined to chat!"))
	Mfile.Write([]byte(nickname + " joined to chat!"))

	for {
		conn.Read(buf)
		Mfile.Write(buf)
		Broadcast(nickname, buf)
	}

	//fmt.Println("Close connection:", conn.RemoteAddr())
}

func Broadcast(nickname string, buf []byte) {
	for _, conn := range AllConns {
		conn.Write(append([]byte(nickname+": "), buf...))
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
		// allConns = append(allConns, conn)

		go communication(conn)
	}

}
