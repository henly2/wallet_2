package main

import (
	"net/rpc"
	"../base/jrpc"
	"../base/common"
	"../base/method"
	//"../business_center"
	"fmt"
	"encoding/json"
)

const ServiceNodeName = "service"
var ServiceNodeApis = []string{
	"getaddress",
	"querysomething",
}
type ServiceNodeInstance struct{
	name string
	// add customer ...
}

func (ni *ServiceNodeInstance)HandleCall(req *common.ServiceCenterDispatchData, res *string) error {
	fmt.Println("I got...api=" , req.Api, "...params=", req.Params)

	*res = "no"
	for i := 0; i < len(ServiceNodeApis); i++ {
		if ServiceNodeApis[i] == req.Api {
			*res = "ok"
			break
		}
	}
	//return business.HandleMsg(req, res)
	return nil;
}

func main() {
	var params, req string
	params = "hello"

	for ; ;  {
		fmt.Println("Please input command: ")
		var input string
		fmt.Scanln(&input)

		if input == "quit" {
			fmt.Println("I do quit")
			break;
		}else if input == "register" {
			node := new(method.ServiceNode)
			nodeInstance := &ServiceNodeInstance{name: ServiceNodeName}
			node.Instance = nodeInstance
			rpc.Register(node)

			go jrpc.StartJRPCTcpServer(":8090")

			var registerData common.ServiceCenterRegisterData
			registerData.Name = ServiceNodeName
			registerData.Addr = "127.0.0.1:8090"
			registerData.Apis = ServiceNodeApis
			b,err := json.Marshal(registerData);
			if err != nil {
				fmt.Println("Error: ", err.Error())
				continue;
			}
			params = string(b[:])
			fmt.Println("params: ", params)
			go jrpc.CallJRPCToTcpServer("127.0.0.1:8081", common.MethodServiceCenterRegister, params, &req)
		}
	}
}
