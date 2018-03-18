package main

import (
	"net/rpc"
	"../base/jrpc"
	"../base/method"
	"../base/common"
	"fmt"
	"time"
	"encoding/json" // for json get
	"log"
	"errors"
	"sync"
)

type ServiceNodeInfo struct{
	registerData common.ServiceCenterRegisterData

	rwmu sync.RWMutex
	client *rpc.Client
}

type ModuleBusiness struct{
	nodes map[string]*ServiceNodeInfo
}

const ServiceCenterName = "root"
type ServiceCenterInstance struct{
	name string

	mu sync.RWMutex
	apiBusinessMap map[string]string
	moduleBusinessMap map[string]*ModuleBusiness
}

func (mi *ServiceCenterInstance)Init() error {
	mi.apiBusinessMap = make(map[string]string)
	mi.moduleBusinessMap = make(map[string]*ModuleBusiness)

	return nil
}

func (mi *ServiceCenterInstance)RegisterServiceNodeInfo(registerData *common.ServiceCenterRegisterData) error{
	mi.mu.Lock()
	defer mi.mu.Unlock()

	business := mi.moduleBusinessMap[registerData.Name]
	if business == nil {
		business = new(ModuleBusiness)
		mi.moduleBusinessMap[registerData.Name] = business;
	}

	for i := 0; i < len(registerData.Apis); i++ {
		mi.apiBusinessMap[registerData.Apis[i]] = registerData.Name;
	}

	if business.nodes == nil {
		business.nodes = make(map[string]*ServiceNodeInfo)
	}

	if business.nodes[registerData.Addr] == nil {
		business.nodes[registerData.Addr] = &ServiceNodeInfo{registerData:*registerData, client:nil};
	}

	fmt.Println("nodes = ", len(business.nodes))
	return nil
}

func (mi *ServiceCenterInstance)GetServiceNodeInfoByApi(api string) *ServiceNodeInfo{
	mi.mu.RLock()
	defer mi.mu.RUnlock()

	name := mi.apiBusinessMap[api]
	if name == ""{
		return nil
	}

	business := mi.moduleBusinessMap[name]
	if business == nil || business.nodes == nil{
		return nil
	}

	var nodeInfo *ServiceNodeInfo
	nodeInfo = nil
	for _, v := range business.nodes{
		nodeInfo = v
		break
	}

	// first we return index 0
	return nodeInfo
}

func (mi *ServiceCenterInstance)RemoveServiceNodeInfo(nodeInfo *ServiceNodeInfo) error{
	mi.mu.Lock()
	defer mi.mu.Unlock()

	business := mi.moduleBusinessMap[nodeInfo.registerData.Name]
	if business == nil || business.nodes == nil{
		return nil
	}

	delete(business.nodes, nodeInfo.registerData.Addr)

	fmt.Println("nodes = ", len(business.nodes))
	return nil
}

func (mi *ServiceCenterInstance)OpenClient(nodeInfo *ServiceNodeInfo) error{
	nodeInfo.rwmu.Lock()
	defer nodeInfo.rwmu.Unlock()

	if nodeInfo.client == nil{
		client, err := rpc.Dial("tcp", nodeInfo.registerData.Addr)
		if err != nil {
			log.Println("Error Open client: ", err.Error())
			return err
		}

		nodeInfo.client = client
	}

	return nil
}

func (mi *ServiceCenterInstance)CloseClient(nodeInfo *ServiceNodeInfo) error{
	nodeInfo.rwmu.Lock()
	defer nodeInfo.rwmu.Unlock()

	if nodeInfo.client != nil{
		nodeInfo.client.Close()
		nodeInfo.client = nil
	}

	//mi.RemoveServiceNodeInfo(nodeInfo)

	return nil
}

/////////////////////////////////////////////////////////////////////
func (mi *ServiceCenterInstance)HandleRegister(req *string, res *string) error {

	RegisterData := common.ServiceCenterRegisterData{}
	err := json.Unmarshal([]byte(*req), &RegisterData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err;
	}

	err = mi.RegisterServiceNodeInfo(&RegisterData)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}

	return nil
}
func (mi *ServiceCenterInstance)HandleDispatch(req *string, res *string) error {
	//jrpc.CallJRPCToTcpServer("127.0.0.1:8081", common.Method_Module_Call, *req, res);
	//jrpc.CallJRPCToTcpServer("192.168.43.123:8081", "Module.Do", *req, res);
	dispatchData := &common.ServiceCenterDispatchData{}
	err := json.Unmarshal([]byte(*req), &dispatchData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err;
	}

	fmt.Println("A module dispatch in...", *req)
	nodeInfo := mi.GetServiceNodeInfoByApi(dispatchData.Api)
	if nodeInfo == nil {
		fmt.Println("Error: no find api")
		return errors.New("no find api")
	}

	if nodeInfo.client == nil {
		mi.OpenClient(nodeInfo)
	}

	err = func() error {
		if nodeInfo.client != nil {
			nodeInfo.rwmu.RLock()
			defer nodeInfo.rwmu.RUnlock()
			return jrpc.CallJRPCToTcpServerOnClient(nodeInfo.client, common.MethodServiceNodeCall, dispatchData, res)
		}
		return nil
	}()

	if err != nil {
		fmt.Println("Call service api failed close client, ", err.Error())

		mi.CloseClient(nodeInfo)
		return err;
	}

	fmt.Println("A module dispatch in callback")
	return err
}

func StartServiceCenter(){
	l, err := jrpc.CreateTcpServer(":8081")
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

	fmt.Println("i am graceful quit StartServiceCenter")
}

func main() {
	center := new(method.ServiceCenter)

	centerInstance := &ServiceCenterInstance{name: ServiceCenterName}
	centerInstance.Init()
	center.Instance = centerInstance
	rpc.Register(center)

	rpc.HandleHTTP();
	go jrpc.StartJRPCHttpServer(":8080")

	go StartServiceCenter()

	//go jrpc.StartJRPCHttpServer2(":8080")
	//newServer := rpc.NewServer();
	//newServer.Register(module)
	//newServer.HandleHTTP("/path", "/debug")
	//go jrpc.StartJRPCHttpServer3(newServer,":8080")

	time.Sleep(time.Second*2)
	for ; ;  {
		fmt.Println("Please input command: ")
		var input string
		fmt.Scanln(&input)

		if input == "quit" {
			fmt.Println("I do quit")
			break;
		}
	}
}

