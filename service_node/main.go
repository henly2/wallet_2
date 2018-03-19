package main

import (
	"net/rpc"
	"../base/method"
	"../base/service"
	//"../business_center"
	"fmt"
	"strconv"
	"context"
)

const ServiceNodeName = "service"
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

func CallNodeApi(api *string, data *string, result *string){
	// dispatch your func
}

func main() {
	node := new(method.ServiceNode)

	nodeInstance, _:= service.NewServiceNodeInstance(ServiceNodeName)
	nodeInstance.RegisterData.Addr = "127.0.0.1:8090"
	nodeInstance.RegisterData.RegisterApi(new(MyFunc1))
	nodeInstance.RegisterData.RegisterApi(new(MyFunc2))

	nodeInstance.ServiceCenterAddr = "127.0.0.1:8081"
	nodeInstance.Handler = CallNodeApi

	node.Instance = nodeInstance
	rpc.Register(node)

	// start routine
	ctx, cancel := context.WithCancel(context.Background())
	go nodeInstance.Start(ctx)

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
