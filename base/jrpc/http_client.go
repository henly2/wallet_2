package jrpc

import (
	"log"
	"net/rpc"
)

// Call a JRPC to Http server
// @parameter: addr string, like "127.0.0.1:8080"
// @parameter: method string
// @parameter: params string
// @parameter: res *string
// @return: error
func CallJRPCToHttpServer(addr string, method string, params interface{}, res *string) error {
	client, err := rpc.DialHTTPPath("tcp", addr, rpc.DefaultRPCPath)
	if err != nil {
		log.Println("#CallJRPCToHttpServer Error: ", err.Error())
		return err
	}
	defer client.Close()

	return CallJRPCToHttpServerOnClient(client, method, params, res)
	if err != nil {
		log.Println("#CallJRPCToHttpServer Error: ", err.Error())
		return err
	}

	return nil
}

// Call a JRPC to Http server on a client
// @parameter: client
// @parameter: method string
// @parameter: params string
// @parameter: res *string
// @return: error
func CallJRPCToHttpServerOnClient(client *rpc.Client, method string, params interface{}, res *string) error {
	err := client.Call(method, params, res)
	if err != nil {
		log.Println("#CallJRPCToHttpServerOnClient Error: ", err.Error())
		return err
	}

	return nil
}