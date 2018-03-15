package jrpc

import (
	"net/rpc"
	"net"
	"log"
)

// Start a JRPC tcp server
// @parameter: port string, like: ":8080"
// @return: error
func StartJRPCTcpServer(port string) error{
	log.Println("Start JRPC Tcp server...", port)

	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err;
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err;
	}

	log.Println("Start JRPC Tcp server successfully, listen on port: ", addr)
	for{
		conn, err := listener.Accept();
		if err != nil {
			log.Println("Error: ", err.Error())
			continue
		}

		log.Println("JRPC Tcp server Accept a client: ", conn.RemoteAddr())
		go rpc.ServeConn(conn)
	}

	return nil;
}
