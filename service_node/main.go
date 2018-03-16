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
)

const ServiceNodeName = "service"
type ServiceNodeInstance struct{
	registerData common.ServiceCenterRegisterData
	// add other customer ...
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

func RegisterToServiceCenter(nodeInstance *ServiceNodeInstance){
	var params, res string
	b,err := json.Marshal(nodeInstance.registerData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return;
	}
	params = string(b[:])
	fmt.Println("params: ", params)

	status := 1
	for ; ; {
		if status == 1 {
			err = jrpc.CallJRPCToTcpServer("127.0.0.1:8081", common.MethodServiceCenterRegister, params, &res)
			if err != nil {
				time.Sleep(time.Second*5)
				continue
			}else{
				status = 0
			}
		}

		if status == 0{
			time.Sleep(time.Second*10)
			err = jrpc.CallJRPCToTcpServer("127.0.0.1:8081", common.MethodServiceCenterPingpong, "ping", &res)
			if err == nil && res == "pong" {
				status = 0
			}else{
				status = 1
			}
		}
	}
}

func main() {
	func (){
		node := new(method.ServiceNode)
		nodeInstance := &ServiceNodeInstance{}
		nodeInstance.registerData.Name = ServiceNodeName
		nodeInstance.registerData.Addr = "127.0.0.1:8090"
		nodeInstance.registerData.RegisterApi(new(MyFunc1))
		nodeInstance.registerData.RegisterApi(new(MyFunc2))

		node.Instance = nodeInstance
		rpc.Register(node)

		go jrpc.StartJRPCTcpServer(":8090")

		go RegisterToServiceCenter(nodeInstance)
	}()

	for ; ;  {
		fmt.Println("Please input command: ")
		var input string
		fmt.Scanln(&input)

		if input == "quit" {
			fmt.Println("I do quit")
			break;
		}
	}
}
