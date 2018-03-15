package jrpc

import (
	"net/rpc"
	"log"
	"fmt"
)

// Call a JRPC to Tcp server
// @parameter: addr string, like "127.0.0.1:8080"
// @parameter: method string
// @parameter: params string
// @parameter: res *string
// @return: error
func CallJRPCToTcpServer(addr string, method string, params string, res *string) error {
	log.Println("Call JRPC to Tcp server...", addr)

	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}
	defer client.Close()

	err = client.Call(method, params, res)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}

	fmt.Println("Params: ", params)
	fmt.Println("Reply: ", *res)
	return nil
}

// Call a JRPC to Tcp server
// @parameter: client rpc.Client
// @parameter: addr string, like "127.0.0.1:8080"
// @parameter: method string
// @parameter: params string
// @parameter: res *string
// @return: error
func CallJRPCToTcpServerOnClient(client *rpc.Client, method string, params string, res *string) error {
	err := client.Call(method, params, res)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}

	fmt.Println("Params: ", params)
	fmt.Println("Reply: ", *res)
	return nil
}