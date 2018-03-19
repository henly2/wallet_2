package main

import (
	"net/rpc"
	"../base/service"
	"fmt"
	"time"
	"context"
)

const ServiceCenterName = "root"

func main() {
	centerInstance,_ := service.NewServiceCenter(ServiceCenterName)
	centerInstance.HttpPort = ":8080"
	centerInstance.TcpPort = ":8081"

	rpc.Register(centerInstance)
	rpc.HandleHTTP();

	ctx, cancel := context.WithCancel(context.Background())
	go service.StartCenter(ctx, centerInstance)

	time.Sleep(time.Second*2)
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
}

