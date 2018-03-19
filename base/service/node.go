package service

import (
	"sync"
	"../common"
	"../jrpc"
	"fmt"
	"errors"
	"encoding/json"
	"context"
	"strings"
	"io"
	"log"
	"net/rpc"
	"time"
)

type CallNodeApi func(api *string, data *string, result *string)

type ServiceNodeInstance struct{
	// 注册的信息
	RegisterData common.ServiceCenterRegisterData
	// 回掉
	Handler CallNodeApi
	// 服务中心
	ServiceCenterAddr string
	// 等待
	Wg *sync.WaitGroup
}

func (ni *ServiceNodeInstance)HandleCall(req *common.ServiceCenterDispatchData, res *string) error {
	ack := common.ServiceCenterDispatchAckData{}
	ack.Api = req.Api
	ack.Err = common.ServiceDispatchErrOk
	if ni.Handler != nil {
		ni.Handler(&req.Api, &req.Params, &ack.Result)
	}else{
		fmt.Println("Error api call (no handler)--api=" , req.Api, ",params=", req.Params)

		ack.Err = common.ServiceDispatchErrNotFindHanlder
		ack.Errmsg = "Not find handler"
	}

	// to json string
	b,err := json.Marshal(ack);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		*res = "System error"
		return errors.New(*res);
	}

	*res = string(b);
	return nil
}

func (ni *ServiceNodeInstance)Start(ctx context.Context) {
	go startServiceNode(ctx, ni)

	go RegisterToServiceCenter(ctx, ni)

	<-ctx.Done()
}

func NewServiceNodeInstance(serviceName string) (*ServiceNodeInstance, error){
	serviceNodeInstance := &ServiceNodeInstance{}

	serviceNodeInstance.Wg = &sync.WaitGroup{}
	serviceNodeInstance.RegisterData.Name = serviceName

	return serviceNodeInstance, nil
}

///////////////////////////////////////////////////////////////////////
func startServiceNode(ctx context.Context, nodeInstance *ServiceNodeInstance){
	nodeInstance.Wg.Add(1)
	defer nodeInstance.Wg.Done()

	s :=strings.Split(nodeInstance.RegisterData.Addr, ":")
	if len(s) != 2{
		fmt.Println("Error: Node addr is not ip:port format")
		return
	}

	l, err := jrpc.CreateTcpServer(":"+s[1])
	if err != nil {
		fmt.Println("Error: Create tcp server, ", err.Error())
		return
	}

	var conns []io.ReadWriteCloser
	go func(){
		for{
			conn, err := l.Accept();
			if err != nil {
				log.Println("Error: ", err.Error())
				continue
			}

			log.Println("JRPC Tcp server Accept a client: ", conn.RemoteAddr())
			conns = append(conns, conn)
			go rpc.ServeConn(conn)
		}
	}()

	<- ctx.Done()
	for i := 0; i < len(conns); i++ {
		conns[i].Close()
	}
	fmt.Println("i am graceful quit StartServiceNode")
}

func KeepAlive(nodeInstance *ServiceNodeInstance, params *string, status int) int{
	var err error
	var res string
	if status == 1 {
		err = jrpc.CallJRPCToTcpServer(nodeInstance.ServiceCenterAddr, common.MethodServiceCenterRegister, *params, &res)
		if err != nil {
		}else{
			status = 0
		}
	}

	if status == 0{
		err = jrpc.CallJRPCToTcpServer(nodeInstance.ServiceCenterAddr, common.MethodServiceCenterPingpong, "ping", &res)
		if err == nil && res == "pong" {
			status = 0
		}else{
			status = 1
		}
	}

	return status
}

func RegisterToServiceCenter(ctx context.Context, nodeInstance *ServiceNodeInstance){
	nodeInstance.Wg.Add(1)
	defer nodeInstance.Wg.Done()

	var params string
	b,err := json.Marshal(nodeInstance.RegisterData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return;
	}
	params = string(b[:])
	fmt.Println("params: ", params)

	timeout := make(chan bool)
	go func(){
		for ; ; {
			timeout <- true
			time.Sleep(time.Second*10)
			fmt.Println("timeout... ")
		}
	}()

	status := 1
	for ; ; {
		select{
		case <-ctx.Done():
			fmt.Println("done signal")
			status = 2
		case <-timeout:
			fmt.Println("timeout signal")
			status = KeepAlive(nodeInstance, &params, status)
		}

		if status == 2{
			break
		}
	}

	fmt.Println("i am graceful quit RegisterToServiceCenter")
}
