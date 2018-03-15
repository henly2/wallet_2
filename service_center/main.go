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
	client *rpc.Client
}

type ModuleBusiness struct{
	nodes []*ServiceNodeInfo
}

const ServiceCenterName = "root"
type ServiceCenterInstance struct{
	name string

	mu sync.Mutex
	apiBusinessMap map[string]string
	moduleBusinessMap map[string]*ModuleBusiness
}

func (mi *ServiceCenterInstance)Init() error {
	mi.apiBusinessMap = make(map[string]string)
	mi.moduleBusinessMap = make(map[string]*ModuleBusiness)

	return nil
}

func (mi *ServiceCenterInstance)RegisterModuleBusiness(registerData *common.ServiceCenterRegisterData) error{
	client, err := rpc.Dial("tcp", registerData.Addr)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}

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

	business.nodes = append(business.nodes, &ServiceNodeInfo{registerData:*registerData, client:client})

	fmt.Println("nodes = ", len(business.nodes))
	return nil
}

func (mi *ServiceCenterInstance)GetServiceNodeInfoByApi(api string) *ServiceNodeInfo{
	mi.mu.Lock()
	defer mi.mu.Unlock()

	name := mi.apiBusinessMap[api]
	if name == ""{
		return nil
	}

	moduleBusiness := mi.moduleBusinessMap[name]
	if moduleBusiness == nil || len(moduleBusiness.nodes) == 0 {
		return nil
	}

	// first we return index 0
	return moduleBusiness.nodes[0]
}

/////////////////////////////////////////////////////////////////////
func (mi *ServiceCenterInstance)HandleRegister(req *string, res *string) error {

	RegisterData := common.ServiceCenterRegisterData{}
	err := json.Unmarshal([]byte(*req), &RegisterData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err;
	}

	err = mi.RegisterModuleBusiness(&RegisterData)
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

	jrpc.CallJRPCToTcpServerOnClient(nodeInfo.client, common.MethodServiceNodeCall, dispatchData, res)

	fmt.Println("A module dispatch in callback")
	return nil
}

func main() {
	center := new(method.ServiceCenter)

	centerInstance := &ServiceCenterInstance{name: ServiceCenterName}
	centerInstance.Init()
	center.Instance = centerInstance
	rpc.Register(center)

	rpc.HandleHTTP();
	go jrpc.StartJRPCHttpServer(":8080")

	go jrpc.StartJRPCTcpServer(":8081")

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

