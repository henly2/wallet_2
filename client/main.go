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
