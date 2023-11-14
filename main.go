package main

import (
	"fmt"
	"log"
	"net"
	"os"
	client "socketchat/client"
)

var AllConns = make([]net.Conn, 0)
var mfile *os.File

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
	greeting := "Please, enter your nickname (1 word)"
	conn.Write([]byte(greeting))
	mfile.Write([]byte(greeting))
	conn.Read(buf)

	client_name := string(buf)
	// Сделать это отдельной функцией
	res := false
	test := func() {
		if _, ok := client.KnownClients[client_name]; ok == true {
			res = true
		}
	}
	test()

	// If nickname is known then load few last messages
	if res {
		//message_cnt := 0
		first_message := 0
		if len(client.KnownClients[client_name]) < 10 {
			first_message = 0
		} else {
			first_message = len(client.KnownClients[client_name]) - 1 - 10
		}
		for i := first_message; i < first_message+10; i++ {
			//message_cnt++
			if i < len(client.KnownClients[client_name]) {
				conn.Write([]byte(client.KnownClients[client_name][i]))
				mfile.Write()
			}
		}
	} else { // else add him to list
		answer := "Your nickname is " + string(buf)
		conn.Write([]byte(answer))
		client.KnownClients[client_name] = append(client.KnownClients[client_name], greeting, client_name, answer)
	}
	for {

	}
	fmt.Println("Close connection:", conn.RemoteAddr())
}

func main() {

	fmt.Println("Start")

	listener, err := net.Listen("tcp", "localhost:8080") // Open listener socket
	if err != nil {
		log.Fatal("error occured:", err)
	}

	mfile, err := os.OpenFile("server_data.txt", os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer mfile.Close()

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
