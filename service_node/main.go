package main

import (
	"net/rpc"
	"../base/jrpc"
	"../base/common"
	"../base/method"
	//"../business_center"
	"fmt"
	"encoding/json"
	"strconv"
	"time"
	"sync"
	"strings"
	"log"
	"runtime"
	"context"
	"io"
)

const ServiceNodeName = "service"
type ServiceNodeInstance struct{
	registerData common.ServiceCenterRegisterData
	// add other customer ...
	serviceCenterAddr string
	stop chan bool
	wg *sync.WaitGroup
}

type Args struct {
	A int
	B int
}
type MyFunc1 int
type MyFunc2 int
func (myFunc1 *MyFunc1)Add(args *Args, res *string)  error{
	*res = strconv.Itoa(args.A + args.B)
	return nil
}
func (myFunc2 *MyFunc2)Sub(args *Args, res *string)  error{
	*res = strconv.Itoa(args.A - args.B)
	return nil
}

func (ni *ServiceNodeInstance)HandleCall(req *common.ServiceCenterDispatchData, res *string) error {
	fmt.Println("I got...api=" , req.Api, "...params=", req.Params)

	*res = "no"
	if req.Api == "MyFunc1.Add" {
		myFunc1 := new(MyFunc1)
		args := Args{A:0, B:0}
		json.Unmarshal([]byte(req.Params), &args);
		myFunc1.Add(&args, res)

	}else if req.Api == "MyFunc2.Sub" {
		myFunc2 := new(MyFunc2)
		args := Args{A:0, B:0}
		json.Unmarshal([]byte(req.Params), &args);
		myFunc2.Sub(&args, res)
	}

	//return business.HandleMsg(req, res)
	return nil;
}

func StartServiceNode(ctx context.Context, nodeInstance *ServiceNodeInstance){
	nodeInstance.wg.Add(1)
	defer nodeInstance.wg.Done()

	s :=strings.Split(nodeInstance.registerData.Addr, ":")
	if len(s) != 2{
		fmt.Println("Error: Node addr is not ip:port format")
		return
	}

	l, err := jrpc.CreateTcpServer(":"+s[1])
	if err != nil {
		fmt.Println("Error: Create tcp server, ", err.Error())
		return
	}

	var conns []io.ReadWriteCloser
	go func(){
		for{
			conn, err := l.Accept();
			if err != nil {
				log.Println("Error: ", err.Error())
				continue
			}

			log.Println("JRPC Tcp server Accept a client: ", conn.RemoteAddr())
			conns = append(conns, conn)
			go rpc.ServeConn(conn)
		}
	}()

	<- ctx.Done()
	for i := 0; i < len(conns); i++ {
		conns[i].Close()
	}
	fmt.Println("i am graceful quit StartServiceNode")
}

func KeepAlive(nodeInstance *ServiceNodeInstance, params *string, status int) int{
	var err error
	var res string
	if status == 1 {
		err = jrpc.CallJRPCToTcpServer(nodeInstance.serviceCenterAddr, common.MethodServiceCenterRegister, *params, &res)
		if err != nil {
		}else{
			status = 0
		}
	}

	if status == 0{
		err = jrpc.CallJRPCToTcpServer(nodeInstance.serviceCenterAddr, common.MethodServiceCenterPingpong, "ping", &res)
		if err == nil && res == "pong" {
			status = 0
		}else{
			status = 1
		}
	}

	return status
}

func RegisterToServiceCenter(ctx context.Context, nodeInstance *ServiceNodeInstance){
	nodeInstance.wg.Add(1)
	defer nodeInstance.wg.Done()

	var params string
	b,err := json.Marshal(nodeInstance.registerData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return;
	}
	params = string(b[:])
	fmt.Println("params: ", params)

	timeout := make(chan bool)
	go func(){
		for ; ; {
			timeout <- true
			time.Sleep(time.Second*10)
			fmt.Println("timeout... ")
		}
	}()

	status := 1
	for ; ; {
		select{
		case <-ctx.Done():
			fmt.Println("done signal")
			status = 2
		case <-timeout:
			fmt.Println("timeout signal")
			status = KeepAlive(nodeInstance, &params, status)
		}

		if status == 2{
			break
		}
	}

	fmt.Println("i am graceful quit RegisterToServiceCenter")
}

func main() {
	t := 1
	if t==2{
		runtime.GOMAXPROCS(2)
	}

	node := new(method.ServiceNode)
	nodeInstance := &ServiceNodeInstance{}

	nodeInstance.registerData.Name = ServiceNodeName
	nodeInstance.registerData.Addr = "127.0.0.1:8090"
	nodeInstance.registerData.RegisterApi(new(MyFunc1))
	nodeInstance.registerData.RegisterApi(new(MyFunc2))

	nodeInstance.serviceCenterAddr = "127.0.0.1:8081"
	nodeInstance.stop = make(chan bool, 5)
	nodeInstance.wg = &sync.WaitGroup{}

	node.Instance = nodeInstance
	rpc.Register(node)

	ctxA, cancelA := context.WithCancel(context.Background())
	go StartServiceNode(ctxA, nodeInstance)

	ctxB, cancelB := context.WithCancel(context.Background())
	go RegisterToServiceCenter(ctxB, nodeInstance)

	for ; ;  {
		fmt.Println("Please input command: ")
		var input string
		fmt.Scanln(&input)

		if input == "quit" {
			fmt.Println("I do quit")
			cancelA()
			cancelB()
			break;
		}
	}

	nodeInstance.wg.Wait()
}
