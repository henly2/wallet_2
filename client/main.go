package main

import (
	"../base/jrpc"
	"../base/common"
	"fmt"
	"encoding/json" // for json get
)

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
		}else if input == "d1" {
			dispatchData := common.ModuleDispatchData{}
			dispatchData.Name = "n1"
			dispatchData.Params = "[{\"A\":1, \"B\":2}]"
			b,err := json.Marshal(dispatchData);
			if err != nil {
				fmt.Println("Error: ", err.Error())
				continue;
			}
			params = string(b[:])
			jrpc.CallJRPCToHttpServer2("127.0.0.1:8080", "", common.MethodServerCenterDispatch, params, &req)
		}else if input == "d2" {
			params = "[{\"A\":1, \"\":2}]"
			jrpc.CallJRPCToTcpServer("127.0.0.1:8090", common.MethodServerNodeCall, params, &req)
		}
	}
}
