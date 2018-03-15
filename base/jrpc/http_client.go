package jrpc

import (
	"log"
	"net/rpc"
	"fmt"
	"net"
)

// Call a JRPC to Http server2
// @parameter: addr string, like "127.0.0.1:8080"
// @parameter: method string
// @parameter: params string
// @parameter: res *string
// @return: error
func CallJRPCToHttpServer2(addr string, path string, method string, params string, res *string) error {
	log.Println("Call JRPC to Http server...", addr)

	realpath := path
	if  realpath == ""{
		realpath = rpc.DefaultRPCPath
	}
	client, err := rpc.DialHTTPPath("tcp", addr, realpath)
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

// FIXME: xxxx
// Call a JRPC to Http server3
// @parameter: addr string, like "127.0.0.1:8080"
// @parameter: method string
// @parameter: params string
// @parameter: res *string
// @return: error
func CallJRPCToHttpServer3(addr string, method string, params string, res *string) error {
	log.Println("Call JRPC to Http server...", addr)

	addr2, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr2)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}
	defer conn.Close()

	client := rpc.NewClient(conn)
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