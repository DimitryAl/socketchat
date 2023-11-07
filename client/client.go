package client

type TClient struct {
	Ip       string
	Nickname string
	IsOnline bool
}

var KnownClients = make([]TClient, 0) // All clients
