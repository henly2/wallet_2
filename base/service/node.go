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

// 服务节点回调接口
type CallNodeApi func(req *common.ServiceCenterDispatchData, result *string)

// 服务节点信息
type ServiceNode struct{
	// 注册的信息
	RegisterData common.ServiceCenterRegisterData
	// 回掉
	Handler CallNodeApi
	// 服务中心
	ServiceCenterAddr string
	// 等待
	Wg *sync.WaitGroup
}

// 生成一个服务节点
func NewServiceNode(serviceName string) (*ServiceNode, error){
	serviceNode := &ServiceNode{}

	serviceNode.Wg = &sync.WaitGroup{}
	serviceNode.RegisterData.Name = serviceName

	return serviceNode, nil
}

// 启动服务节点
func StartNode(ctx context.Context, serviceNode *ServiceNode) {
	go startServiceNode(ctx, serviceNode)

	go registerToServiceCenter(ctx, serviceNode)

	<-ctx.Done()
}

// RPC 方法
// 服务节点RPC--调用节点方法ServiceNodeInstance.Call
func (ni *ServiceNode) Call(req *common.ServiceCenterDispatchData, res * string) error {
	ack := common.ServiceCenterDispatchAckData{}
	ack.Api = req.Api
	ack.Err = common.ServiceDispatchErrOk
	if ni.Handler != nil {
		ni.Handler(req, &ack.Result)
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

// 内部方法
///////////////////////////////////////////////////////////////////////
func startServiceNode(ctx context.Context, serviceNode *ServiceNode){
	serviceNode.Wg.Add(1)
	defer serviceNode.Wg.Done()

	s :=strings.Split(serviceNode.RegisterData.Addr, ":")
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

func keepAlive(serviceNode *ServiceNode, params *string, status int) int{
	var err error
	var res string
	if status == 1 {
		err = jrpc.CallJRPCToTcpServer(serviceNode.ServiceCenterAddr, common.MethodServiceCenterRegister, *params, &res)
		if err != nil {
		}else{
			status = 0
		}
	}

	if status == 0{
		err = jrpc.CallJRPCToTcpServer(serviceNode.ServiceCenterAddr, common.MethodServiceCenterPingpong, "ping", &res)
		if err == nil && res == "pong" {
			status = 0
		}else{
			status = 1
		}
	}

	return status
}

func registerToServiceCenter(ctx context.Context, serviceNode *ServiceNode){
	serviceNode.Wg.Add(1)
	defer serviceNode.Wg.Done()

	var params string
	b,err := json.Marshal(serviceNode.RegisterData);
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
			status = keepAlive(serviceNode, &params, status)
		}

		if status == 2{
			break
		}
	}

	fmt.Println("i am graceful quit RegisterToServiceCenter")
}
