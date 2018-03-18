package jrpc

import (
	"net"
	"log"
)

// Start a JRPC tcp server
// @parameter: port string, like: ":8080"
// @return: error
func CreateTcpServer(port string) (*net.TCPListener, error){
	log.Println("Start JRPC Tcp server...", port)

	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		log.Println("Error: ", err.Error())
		return nil, err;
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println("Error: ", err.Error())
		return nil, err;
	}

	return listener, nil

}
