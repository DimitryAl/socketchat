package client

type TClient struct {
	//Ip       string
	//Nickname string
	//IsOnline bool
	LastMessages []string
}

var KnownClients = make(map[string][]string) // List of all clients
