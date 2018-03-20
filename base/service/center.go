package service

import (
	"sync"
	"net/rpc"
	"../common"
	"../jrpc"
	"../restful"
	"log"
	"encoding/json"
	"fmt"
	"errors"
	"context"
	"net/http"
	"io/ioutil"
	"strings"
)

type ServiceNodeInfo struct{
	RegisterData common.ServiceCenterRegisterData

	Rwmu sync.RWMutex
	Client *rpc.Client
}

type ServiceNodeBusiness struct{
	AddrMapServiceNodeInfo map[string]*ServiceNodeInfo
}

type ServiceCenter struct{
	// 名称
	Name string

	// 端口
	HttpPort string
	TcpPort string

	// 节点信息
	Rwmu sync.RWMutex
	ApiMapServiceName map[string]string
	ServiceNameMapBusiness map[string]*ServiceNodeBusiness

	// 等待
	Wg *sync.WaitGroup
}

func (mi *ServiceCenter)HandleRequest(w http.ResponseWriter, req *http.Request) error{
	fmt.Println("a client comein...")
	fmt.Println("path=", req.URL.Path)

	api := req.URL.Path
	api = strings.TrimLeft(api, "/")
	api = strings.TrimRight(api, "/")
	api = strings.TrimLeft(api, "wallet/v1")
	api = strings.Replace(api, "/", ".", -1)

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("#HandleRequest Error: ", err.Error())
		return err
	}

	body := string(b)
	fmt.Println("body=", body)

	// 重组rpc结构json
	dispatchData := common.ServiceCenterDispatchData{}
	err = json.Unmarshal(b, &dispatchData);
	if err != nil {
		fmt.Println("#HandleRequest Error: ", err.Error())
		return err;
	}
	dispatchData.Api = api

	dispatchAckData := common.ServiceCenterDispatchAckData{}
	mi.Dispatch(&dispatchData, &dispatchAckData)

	w.Header().Set("Content-Type", "application/json")

	b, err = json.Marshal(dispatchAckData)
	if err != nil {
		fmt.Println("#HandleRequest Error: ", err.Error())
		return err;
	}

	w.Write(b)
	//res := NewRPCRequest(req.Body).Call()
	//io.Copy(w, res)

	return nil
}

// 生成一个服务中心
func NewServiceCenter(rootName string) (*ServiceCenter, error){
	serviceCenter := &ServiceCenter{}

	serviceCenter.Wg = &sync.WaitGroup{}
	serviceCenter.Name = rootName
	serviceCenter.ApiMapServiceName = make(map[string]string)
	serviceCenter.ServiceNameMapBusiness = make(map[string]*ServiceNodeBusiness)

	return serviceCenter, nil
}

// 启动服务中心
func StartCenter(ctx context.Context, serviceCenter *ServiceCenter) {
	go jrpc.StartJRPCHttpServer(serviceCenter.HttpPort)
	go restful.StartRestfulHttpServer(ctx, serviceCenter.Wg, serviceCenter, ":8070")
	go startServiceCenter(ctx, serviceCenter)

	<-ctx.Done()
}

// RPC 方法
// 服务中心方法--注册到服务中心
func (mi *ServiceCenter) Register(req *string, res *string) error {
	RegisterData := common.ServiceCenterRegisterData{}
	err := json.Unmarshal([]byte(*req), &RegisterData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err;
	}

	err = mi.registerServiceNodeInfo(&RegisterData)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}

	return nil
}

// 服务中心方法--发送请求到服务中心进行转发
/*
func (mi *ServiceCenter) Dispatch(req *string, res * string) error {
	dispatchData := &common.ServiceCenterDispatchData{}
	err := json.Unmarshal([]byte(*req), &dispatchData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err;
	}

	fmt.Println("A module dispatch in...", *req)
	nodeInfo := mi.getServiceNodeInfoByApi(dispatchData.Api)
	if nodeInfo == nil {
		fmt.Println("Error: not find api")
		return errors.New("not find api")
	}

	if nodeInfo.Client == nil {
		mi.openClient(nodeInfo)
	}

	err = func() error {
		if nodeInfo.Client != nil {
			nodeInfo.Rwmu.RLock()
			defer nodeInfo.Rwmu.RUnlock()
			return jrpc.CallJRPCToTcpServerOnClient(nodeInfo.Client, common.MethodServiceNodeCall, dispatchData, res)
		}
		return nil
	}()

	if err != nil {
		fmt.Println("Call service api failed close client, ", err.Error())

		mi.closeClient(nodeInfo)
		return err;
	}

	fmt.Println("A module dispatch in callback")
	return err
}
*/
func (mi *ServiceCenter) Dispatch(ask *common.ServiceCenterDispatchData, ack *common.ServiceCenterDispatchAckData) error {
	api := strings.ToLower(ask.Api)

	fmt.Println("A module dispatch in api...", api)
	nodeInfo := mi.getServiceNodeInfoByApi(api)
	if nodeInfo == nil {
		fmt.Println("Error: not find api")
		return errors.New("not find api")
	}

	if nodeInfo.Client == nil {
		mi.openClient(nodeInfo)
	}

	err := func() error {
		if nodeInfo.Client != nil {
			nodeInfo.Rwmu.RLock()
			defer nodeInfo.Rwmu.RUnlock()
			return jrpc.CallJRPCToTcpServerOnClient(nodeInfo.Client, common.MethodServiceNodeCall, ask, ack)
		}
		return nil
	}()

	if err != nil {
		fmt.Println("Call service api failed close client, ", err.Error())

		mi.closeClient(nodeInfo)
		return err;
	}

	return err
}

// 服务中心方法--与服务中心心跳
func (mi *ServiceCenter) Pingpong(req *string, res * string) error {
	if *req == "ping" {
		*res = "pong"
	}else{
		*res = *req
	}
	return nil;
}

// 内部方法
func startServiceCenter(ctx context.Context, serviceCenter *ServiceCenter){
	serviceCenter.Wg.Add(1)
	defer serviceCenter.Wg.Done()

	log.Println("Start JRPC Tcp server on ", serviceCenter.TcpPort)
	l, err := jrpc.CreateTcpServer(serviceCenter.TcpPort)
	if err != nil {
		fmt.Println("Error: Create tcp server, ", err.Error())
		return
	}

	go func(){
		for{
			conn, err := l.Accept();
			if err != nil {
				log.Println("Error: ", err.Error())
				continue
			}

			log.Println("JRPC Tcp server Accept a client: ", conn.RemoteAddr())
			go rpc.ServeConn(conn)
		}
	}()

	<- ctx.Done()
}

func (mi *ServiceCenter)registerServiceNodeInfo(registerData *common.ServiceCenterRegisterData) error{
	mi.Rwmu.Lock()
	defer mi.Rwmu.Unlock()

	//name := registerData.Name
	name := strings.ToLower(registerData.Name)

	business := mi.ServiceNameMapBusiness[name]
	if business == nil {
		business = new(ServiceNodeBusiness)
		mi.ServiceNameMapBusiness[name] = business;
	}

	for i := 0; i < len(registerData.Apis); i++ {
		//api := registerData.Apis[i]
		api := strings.ToLower(registerData.Apis[i])
		mi.ApiMapServiceName[api] = name;
	}

	if business.AddrMapServiceNodeInfo == nil {
		business.AddrMapServiceNodeInfo = make(map[string]*ServiceNodeInfo)
	}

	if business.AddrMapServiceNodeInfo[registerData.Addr] == nil {
		business.AddrMapServiceNodeInfo[registerData.Addr] = &ServiceNodeInfo{RegisterData:*registerData, Client:nil};
	}

	fmt.Println("nodes = ", len(business.AddrMapServiceNodeInfo))
	return nil
}

func (mi *ServiceCenter)getServiceNodeInfoByApi(api string) *ServiceNodeInfo{
	mi.Rwmu.RLock()
	defer mi.Rwmu.RUnlock()

	name := mi.ApiMapServiceName[api]
	if name == ""{
		return nil
	}

	business := mi.ServiceNameMapBusiness[name]
	if business == nil || business.AddrMapServiceNodeInfo == nil{
		return nil
	}

	var nodeInfo *ServiceNodeInfo
	nodeInfo = nil
	for _, v := range business.AddrMapServiceNodeInfo{
		nodeInfo = v
		break
	}

	// first we return index 0
	return nodeInfo
}

func (mi *ServiceCenter)removeServiceNodeInfo(nodeInfo *ServiceNodeInfo) error{
	mi.Rwmu.Lock()
	defer mi.Rwmu.Unlock()

	business := mi.ServiceNameMapBusiness[nodeInfo.RegisterData.Name]
	if business == nil || business.AddrMapServiceNodeInfo == nil{
		return nil
	}

	delete(business.AddrMapServiceNodeInfo, nodeInfo.RegisterData.Addr)

	fmt.Println("nodes = ", len(business.AddrMapServiceNodeInfo))
	return nil
}

func (mi *ServiceCenter)openClient(nodeInfo *ServiceNodeInfo) error{
	nodeInfo.Rwmu.Lock()
	defer nodeInfo.Rwmu.Unlock()

	if nodeInfo.Client == nil{
		client, err := rpc.Dial("tcp", nodeInfo.RegisterData.Addr)
		if err != nil {
			log.Println("Error Open client: ", err.Error())
			return err
		}

		nodeInfo.Client = client
	}

	return nil
}

func (mi *ServiceCenter)closeClient(nodeInfo *ServiceNodeInfo) error{
	nodeInfo.Rwmu.Lock()
	defer nodeInfo.Rwmu.Unlock()

	if nodeInfo.Client != nil{
		nodeInfo.Client.Close()
		nodeInfo.Client = nil
	}

	//mi.RemoveServiceNodeInfo(nodeInfo)

	return nil
}