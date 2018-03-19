package main

import (
	"net/rpc"
	"../base/service"
	"../base/common"
	"./function"
	//"../business_center"
	"fmt"
	"strconv"
	"context"
	"net/rpc/jsonrpc"
	"bufio"
	"strings"
)

const ServiceNodeName = "service"

// rpcRequest represents a RPC request.
// rpcRequest implements the io.ReadWriteCloser interface.
type apiRpcRequest struct {
	r    bufio.Reader
	res  *string
	done chan bool     // signals then end of the RPC request
}

// Read implements the io.ReadWriteCloser Read method.
func (r *apiRpcRequest) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

// Write implements the io.ReadWriteCloser Write method.
func (r *apiRpcRequest) Write(p []byte) (n int, err error) {
	//return r.rw.Write(p)
	*r.res += string(p)
	return len(*r.res),nil
}

// Close implements the io.ReadWriteCloser Close method.
func (r *apiRpcRequest) Close() error {
	r.done <- true
	return nil
}

// Call invokes the RPC request, waits for it to complete, and returns the results.
func (r *apiRpcRequest) Call() {
	go jsonrpc.ServeConn(r)
	<-r.done
}

func callNodeApi(req *common.ServiceCenterDispatchData, result *string){
	// dispatch your func
	rpcstring := "{\"method\":\"" + req.Api + "\"," + "\"params\":" + req.Argv + ",\"id\":"+ strconv.Itoa(req.Id) + "}"

	s := strings.NewReader(rpcstring)
	br := bufio.NewReader(s)

	done := make(chan bool)
	apiRequest := &apiRpcRequest{*br, result, done}
	apiRequest.Call()

	fmt.Println("callNodeApi: ", *result)
}

func main() {
	nodeInstance, _:= service.NewServiceNode(ServiceNodeName)
	nodeInstance.RegisterData.Addr = "127.0.0.1:8090"
	nodeInstance.RegisterData.RegisterApi(new(function.MyFunc1))
	nodeInstance.RegisterData.RegisterApi(new(function.MyFunc2))

	nodeInstance.ServiceCenterAddr = "127.0.0.1:8081"
	nodeInstance.Handler = callNodeApi

	rpc.Register(nodeInstance)

	rpc.Register(new(function.MyFunc1))
	rpc.Register(new(function.MyFunc2))

	// start routine
	ctx, cancel := context.WithCancel(context.Background())
	go service.StartNode(ctx, nodeInstance)

	// console command
	for ; ;  {
		fmt.Println("Please input command: ")
		var input string
		fmt.Scanln(&input)

		if input == "quit" {
			fmt.Println("I do quit")
			cancel()
			break;
		}
	}

	nodeInstance.Wg.Wait()
}
