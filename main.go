package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var AllConns = make([]net.Conn, 0)
var ConnsNicks = make(map[net.Conn]string)

func clearMessage(text []byte) []byte {
	res := make([]byte, 32)
	j := 0
	for i := 0; i < len(text); i++ {
		if text[i] != byte('\x00') {
			res[j] = text[i]
			j++
		}
	}
	return res[:j]
}

func writeFile(text []byte) {
	Mfile, err := os.OpenFile("server_data.txt", os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer Mfile.Close()

	Mfile.Write(append(text, byte('\n')))
}

func readFile() []string {

	var res []string

	Mfile, err := os.OpenFile("server_data.txt", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer Mfile.Close()

	scanner := bufio.NewScanner(Mfile)
	scanner.Split(bufio.ScanLines)
	var fileLines []string
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}

	if len(fileLines) > 10 {
		start := len(fileLines) - 11
		for i := start; i < start+10; i++ {
			res = append(res, fileLines[i]+"\n")
		}
	} else {
		for _, line := range fileLines {
			res = append(res, line+"\n")
		}
	}

	return res
}

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
	Broadcast(conn, "Server", []byte(ConnsNicks[conn]+" leaves chat!"))
	writeFile([]byte("Server: " + ConnsNicks[conn] + " leaves chat!"))
	fmt.Println("Closing connection: ", conn.RemoteAddr().String())
}

func communication(conn net.Conn) {
	defer closeConnection(conn)

	buf := make([]byte, 32) // Bufer for client message

	// Get client's name
	greeting := "Please, enter your nickname"
	conn.Write([]byte("Server: " + greeting))
	conn.Read(buf)
	nickname := string(clearMessage(buf))
	ConnsNicks[conn] = nickname

	// send last 10 messages
	lastMessages := readFile()
	for _, message := range lastMessages {
		conn.Write([]byte(message))
	}

	conn.Write([]byte("Server: " + nickname + " joined to chat!"))
	writeFile([]byte("Server: " + nickname + " joined to chat!"))
	Broadcast(conn, "Server", []byte(nickname+" joined to chat!"))

	for {

		buf = make([]byte, 32)
		_, err := conn.Read(buf)
		if err != nil {
			break
		}

		if len(buf) != 0 {
			buf = clearMessage(buf)
			temp := []byte(nickname + ": " + string(buf))
			writeFile(temp)
		}
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

	if _, err := os.Stat("server_data.txt"); err != nil {
		os.Create("server_data.txt")
	}

	listener, err := net.Listen("tcp", "localhost:8080") // Open listener socket
	if err != nil {
		log.Fatal("error occured:", err)
	}

	for {
		conn, err := listener.Accept() // Accept TCP-connection from client

		if err != nil {
			continue
		}

		fmt.Print("Opening connection: ")
		fmt.Println(conn.RemoteAddr())
		AllConns = append(AllConns, conn)

		go communication(conn)
	}

}
