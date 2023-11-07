package main

import (
	"fmt"
	"log"
	"net"
	client "socketchat/client"
)

func readClient(conn net.Conn) {
	defer conn.Close() // Close connection

	buf := make([]byte, 32) // Bufer for client message

	res := false
	test := func() {
		for _, v := range client.KnownClients {
			if conn.RemoteAddr().String() == v.Ip {
				res = true
			}
		}
	}
	test()

	for {

		if res {
			_, err := conn.Read(buf)

			if err != nil {
				fmt.Println("err")
				break
			}

			fmt.Println(buf)
		} else {
			conn.Write([]byte("Enter your nickname"))

			_, err := conn.Read(buf)

			if err != nil {
				fmt.Println("err")
				break
			}

			fmt.Println(buf)
			client.KnownClients = append(client.KnownClients, client.TClient{Ip: conn.RemoteAddr().String(), Nickname: string(buf[:]), IsOnline: true})

		}
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
