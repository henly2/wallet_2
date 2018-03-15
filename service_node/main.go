package main

import (
	"net/rpc"
	"../base/jrpc"
	"../base/method"
	"../base/common"
	"fmt"
	"encoding/json"
)

type ModuleInstance struct{
	name string
}

/////////////////////////////////////////////////////////////////////
func (mi *ModuleInstance)HandleRegister(req *string, res *string) error {

	fmt.Println("error")
	return nil
}
func (mi *ModuleInstance)HandleDispatch(req *string, res *string) error {

	fmt.Println("error")
	return nil
}

func (mi *ModuleInstance)HandleCall(req *string, res *string) error {

	//jrpc.CallJRPCToTcpServer("127.0.0.1:8081", "Module.Do", *req, res);
	//jrpc.CallJRPCToTcpServer("192.168.43.123:8081", "Module.Do", *req, res);
	fmt.Println("A module call in")
	*res = "haha"
	return nil
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
			module := new(method.Module)
			module.Instance = &ModuleInstance{name:"m1"}
			rpc.Register(module)

			go jrpc.StartJRPCTcpServer(":8090");

			var mdata common.ModuleRegisterData
			mdata.Name = "m1";
			mdata.Addr = "127.0.0.1:8090"
			b,err := json.Marshal(mdata);
			if err != nil {
				fmt.Println("Error: ", err.Error())
				continue;
			}
			params = string(b[:])
			fmt.Println("params: ", params)
			go jrpc.CallJRPCToTcpServer("127.0.0.1:8081", common.Method_Module_Register, params, &req)
		}else if input == "http" {
			go jrpc.CallJRPCToHttpServer2("127.0.0.1:8080", "", common.Method_Module_Dispatch, params, &req)
		}else if input == "dispatch" {
			var mdata common.ModuleDispatchData
			mdata.Name = "m1";
			b,err := json.Marshal(mdata);
			if err != nil {
				fmt.Println("Error: ", err.Error())
				continue;
			}
			params = string(b[:])
			fmt.Println("params: ", params)
			go jrpc.CallJRPCToHttpServer2("127.0.0.1:8080", "", common.Method_Module_Dispatch, params, &req)
		}else if input == "dispatch2" {
			params = "{\"name\":\"m1\", \"params\":\"pp\"}"
			fmt.Println("params: ", params)
			go jrpc.CallJRPCToHttpServer2("127.0.0.1:8080", "", common.Method_Module_Dispatch, params, &req)
		}else if input == "dispatch2_1" {

			for i:=0; i<100;i++  {
				params = "{\"name\":\"m1\", \"params\":\"pp\"}"
				fmt.Println("params: ", params)
				go jrpc.CallJRPCToHttpServer2("127.0.0.1:8080", "", common.Method_Module_Dispatch, params, &req)
			}

		}
	}
}
