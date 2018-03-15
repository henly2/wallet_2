package main

import (
	"net/rpc"
	"../base/jrpc"
	"../base/common"
	"../base/method"
	"fmt"
	"encoding/json"
)

const ServiceNodeName = "n1"
type ServiceNodeInstance struct{
	name string
}

func (ni *ServiceNodeInstance)HandleCall(req *string, res *string) error {
	fmt.Println("I got..." , *req)

	*res = "i am node";
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
			node.Instance = &ServiceNodeInstance{name: ServiceNodeName}
			rpc.Register(node)

			go jrpc.StartJRPCTcpServer(":8090");

			var registerData common.ServiceCenterRegisterData
			registerData.Name = ServiceNodeName;
			registerData.Addr = "127.0.0.1:8090"
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
