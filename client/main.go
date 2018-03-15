package main

import (
	"../base/jrpc"
	"../base/common"
	"fmt"
	"encoding/json" // for json get
	"sync/atomic"
	"net/rpc"
	"log"
	"time"
)

var timeBegin,timeEnd time.Time

func DoTest(params string, count *int64, right *int64, times int64){
	var res string
	err := jrpc.CallJRPCToHttpServer2("127.0.0.1:8080", "", common.MethodServiceCenterDispatch, params, &res)

	atomic.AddInt64(count, 1)
	if  err == nil && res == "ok"{
		atomic.AddInt64(right, 1)
	}

	if atomic.CompareAndSwapInt64(count, times, times) {
		cost := time.Now().Sub(timeBegin)
		fmt.Println("finish...", *count, "...right...", *right, "...cost...", cost)
	}
}

func DoTestTcp(params string, count *int64, right *int64, times int64){
	var res string

	err := jrpc.CallJRPCToTcpServer("127.0.0.1:8090", common.MethodServiceNodeCall, params, &res)

	atomic.AddInt64(count, 1)
	if  err == nil && res == "ok"{
		atomic.AddInt64(right, 1)
	}

	if atomic.CompareAndSwapInt64(count, times, times) {
		cost := time.Now().Sub(timeBegin)
		fmt.Println("finish...", *count, "...right...", *right, "...cost...", cost)
	}
}

func DoTestTcp2(client *rpc.Client, params string, count *int64, right *int64, times int64){
	var res string

	err := jrpc.CallJRPCToTcpServerOnClient(client, common.MethodServiceNodeCall, params, &res)

	atomic.AddInt64(count, 1)
	if  err == nil && res == "ok"{
		atomic.AddInt64(right, 1)
	}

	if atomic.CompareAndSwapInt64(count, times, times) {
		cost := time.Now().Sub(timeBegin)
		fmt.Println("finish...", *count, "...right...", *right, "...cost...", cost)
	}
}

func main() {
	var params, res string
	params = "hello"

	const times = 100;
	var count, right int64
	count = 0
	right = 0

	dispatchData := common.ServiceCenterDispatchData{}
	dispatchData.Api = "getaddress2"
	dispatchData.Params = "{\"A\":1, \"B\":2}"
	b,err := json.Marshal(dispatchData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return;
	}
	params = string(b[:])

	for ; ;  {
		count = 0
		right = 0

		fmt.Println("Please input command: ")
		var input string
		fmt.Scanln(&input)

		timeBegin = time.Now()

		if input == "quit" {
			fmt.Println("I do quit")
			break;
		}else if input == "d1" {
			go jrpc.CallJRPCToHttpServer2("127.0.0.1:8080", "", common.MethodServiceCenterDispatch, params, &res)
		}else if input == "d2" {
			go jrpc.CallJRPCToTcpServer("127.0.0.1:8090", common.MethodServiceNodeCall, dispatchData, &res)
		}else if input == "d3" {
			for i := 0; i < times; i++ {
				go DoTestTcp(params, &count, &right, times)
			}
		} else if input == "d33" {

			client, err := rpc.Dial("tcp", "127.0.0.1:8090")
			if err != nil {
				log.Println("Error: ", err.Error())
				continue
			}

			for i := 0; i < times; i++ {
				go DoTestTcp2(client, params, &count, &right, times)
			}
		}else if input == "d4" {
			for i := 0; i < times; i++ {
				go DoTest(params, &count, &right, times)
			}
		}
	}
}
