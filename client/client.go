package client

type TClient struct {
	ip       int
	nickname string
	isOnline bool
}

var curClients = make([]TClient, 10) // All clients
