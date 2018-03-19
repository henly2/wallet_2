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
	"strconv"
)

var timeBegin,timeEnd time.Time

func DoTest(params interface{}, str *string, count *int64, right *int64, times int64){
	ackData := common.ServiceCenterDispatchAckData{}
	err := jrpc.CallJRPCToHttpServer("127.0.0.1:8080", common.MethodServiceCenterDispatch, params, &ackData)

	atomic.AddInt64(count, 1)
	if  err == nil && ackData.Err==100{
		atomic.AddInt64(right, 1)
	}

	if atomic.CompareAndSwapInt64(count, times, times) {
		cost := time.Now().Sub(timeBegin)
		fmt.Println("finish...", *count, "...right...", *right, "...cost...", cost)
	}
}

func DoTest2(client *rpc.Client, params interface{}, str *string, count *int64, right *int64, times int64){
	ackData := common.ServiceCenterDispatchAckData{}
	err := jrpc.CallJRPCToHttpServerOnClient(client, common.MethodServiceCenterDispatch, params, &ackData)

	atomic.AddInt64(count, 1)
	if  err == nil && ackData.Err==100{
		atomic.AddInt64(right, 1)
	}

	if atomic.CompareAndSwapInt64(count, times, times) {
		cost := time.Now().Sub(timeBegin)
		fmt.Println("finish...", *count, "...right...", *right, "...cost...", cost)
	}
}

func DoTestTcp(params interface{}, str *string, count *int64, right *int64, times int64){
	ackData := common.ServiceCenterDispatchAckData{}

	err := jrpc.CallJRPCToTcpServer("127.0.0.1:8090", common.MethodServiceNodeCall, params, &ackData)

	atomic.AddInt64(count, 1)
	if  err == nil && ackData.Err==100{
		atomic.AddInt64(right, 1)
	}

	if atomic.CompareAndSwapInt64(count, times, times) {
		cost := time.Now().Sub(timeBegin)
		fmt.Println("finish...", *count, "...right...", *right, "...cost...", cost)
	}
}

func DoTestTcp2(client *rpc.Client, params interface{}, str *string, count *int64, right *int64, times int64){
	ackData := common.ServiceCenterDispatchAckData{}
	err := jrpc.CallJRPCToTcpServerOnClient(client, common.MethodServiceNodeCall, params, &ackData)

	atomic.AddInt64(count, 1)
	if  err == nil && ackData.Err==100{
		atomic.AddInt64(right, 1)
	}

	if atomic.CompareAndSwapInt64(count, times, times) {
		cost := time.Now().Sub(timeBegin)
		fmt.Println("finish...", *count, "...right...", *right, "...cost...", cost)
	}
}

// curl -d '{"method":"ServiceCenter.Dispatch", "params":["{\"api\":\"MyFunc2.Sub\",\"params\":\"[{\\\"A\\\":\\\"hello, \\\", \\\"B\\\":\\\"world\\\"}]\",\"id\":1}"], "id": 1 }' http://localhost:8080

// curl -d '{
// "method":"ServiceCenter.Dispatch",
// "params":["{\"api\":\"MyFunc2.Sub\",\"params\":\"[{\\\"A\\\":\\\"hello, \\\", \\\"B\\\":\\\"world\\\"}]\", \"id\":1}"],
// "id": 1
// }'
// http://localhost:8080

// curl -d '{"method":"ServiceCenter.Test", "params":[{"api":"MyFunc2.Sub","argv":"[{\"A\":\"hello, \", \"B\":\"world\"}]", "id":1}], "id": 1}' http://localhost:8080
// curl -d '{
// "method":"ServiceCenter.Dispatch",
// "params":[{"api":"MyFunc2.Sub","argv":"[{\"A\":\"hello, \", \"B\":\"world\"}], "id":1}],
// "id": 1
// }'
// http://localhost:8080
func main() {
	const times = 100;
	var count, right int64
	count = 0
	right = 0

	var testdata string
	for i := 0; i < 1000; i++ {
		testdata += strconv.Itoa(i)
	}
	testdata = "hello, world"

	dispatchData := common.ServiceCenterDispatchData{}
	dispatchData.Api = "MyFunc2.Sub"
	dispatchData.Argv = "[{\"A\":\"hello, \", \"B\":\"world\"}]"
	dispatchData.Id = 1;
	b,err := json.Marshal(dispatchData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return;
	}

	fmt.Println("argv:", string(b[:]))

	ackData := common.ServiceCenterDispatchAckData{}

	for ; ;  {
		fmt.Println("Please input command: ")
		var input string
		fmt.Scanln(&input)

		fmt.Println("Execute input command: ")
		count = 0
		right = 0
		timeBegin = time.Now();

		if input == "quit" {
			fmt.Println("I do quit")
			break;
		}else if input == "d1" {
			jrpc.CallJRPCToHttpServer("127.0.0.1:8080", common.MethodServiceCenterDispatch, dispatchData, &ackData)
			fmt.Println("ack==", ackData.Ack)
		}else if input == "d2" {
			jrpc.CallJRPCToTcpServer("127.0.0.1:8090", common.MethodServiceNodeCall, dispatchData, &ackData)
			fmt.Println("ack==", ackData.Ack)
		}else if input == "d3" {
			for i := 0; i < times; i++ {
				dispatchData.Id = i
				go DoTestTcp(dispatchData, &testdata, &count, &right, times)
			}
		} else if input == "d33" {

			client, err := rpc.Dial("tcp", "127.0.0.1:8090")
			if err != nil {
				log.Println("Error: ", err.Error())
				continue
			}

			for i := 0; i < times; i++ {
				dispatchData.Id = i
				go DoTestTcp2(client, dispatchData, &testdata, &count, &right, times)
			}
		}else if input == "d4" {
			for i := 0; i < times; i++ {
				dispatchData.Id = i
				go DoTest(dispatchData, &testdata, &count, &right, times)
			}
		} else if input == "d44" {

			addr := "127.0.0.1:8080"
			log.Println("Call JRPC to Http server...", addr)

			realpath := ""
			if  realpath == ""{
				realpath = rpc.DefaultRPCPath
			}
			client, err := rpc.DialHTTPPath("tcp", addr, realpath)
			if err != nil {
				log.Println("Error: ", err.Error())
				continue
			}
			for i := 0; i < times; i++ {
				dispatchData.Id = i
				go DoTest2(client, dispatchData, &testdata, &count, &right, times)
			}
		}
	}
}
