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

type ModuleNode struct{
	name string
	addr string
	client *rpc.Client
}

type ModuleBusiness struct{
	nodes []ModuleNode
}

const ServerCenterName = "root"
type ServerCenterInstance struct{
	name string

	mu sync.Mutex
	moduleBusinessMap map[string]*ModuleBusiness
}

func (mi *ServerCenterInstance)Init() error {
	mi.moduleBusinessMap = make(map[string]*ModuleBusiness)

	return nil
}

func (mi *ServerCenterInstance)CreateAndGetModuleBusinessByName(name string) *ModuleBusiness{
	var business *ModuleBusiness
	mi.mu.Lock()
	defer mi.mu.Unlock()

	business = mi.moduleBusinessMap[name]
	if business == nil {
		business = new(ModuleBusiness)
		mi.moduleBusinessMap[name] = business;
	}

	return business
}

func (mi *ServerCenterInstance)GetModuleBusinessByName(name string) *ModuleBusiness{
	mi.mu.Lock()
	defer mi.mu.Unlock()

	return mi.moduleBusinessMap[name]
}

/////////////////////////////////////////////////////////////////////
func (mi *ServerCenterInstance)HandleRegister(req *string, res *string) error {

	RegisterData := &common.ModuleRegisterData{}
	err := json.Unmarshal([]byte(*req), &RegisterData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err;
	}

	fmt.Println("A module register in...", mi.name, "--", RegisterData.Name)
	client, err := rpc.Dial("tcp", RegisterData.Addr)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}

	business := mi.CreateAndGetModuleBusinessByName(RegisterData.Name)
	business.nodes = append(business.nodes, ModuleNode{name:RegisterData.Name, addr:RegisterData.Addr, client:client})

	fmt.Println("nodes = ", len(business.nodes))
	return nil
}
func (mi *ServerCenterInstance)HandleDispatch(req *string, res *string) error {

	fmt.Println("A module dispatch in callback2...", mi.name)
	//jrpc.CallJRPCToTcpServer("127.0.0.1:8081", common.Method_Module_Call, *req, res);
	//jrpc.CallJRPCToTcpServer("192.168.43.123:8081", "Module.Do", *req, res);
	dispatchData := &common.ModuleDispatchData{}
	err := json.Unmarshal([]byte(*req), &dispatchData);
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err;
	}

	fmt.Println("A module dispatch in...", dispatchData.Name)
	business := mi.GetModuleBusinessByName(dispatchData.Name)
	if business==nil || len(business.nodes) == 0{
		fmt.Println("Error: no module")
		return errors.New("no module")
	}

	fmt.Println("A module dispatch in callback1")

	node := business.nodes[0]
	jrpc.CallJRPCToTcpServerOnClient(node.client, common.MethodServerNodeCall, dispatchData.Params, res)

	fmt.Println("A module dispatch in callback")
	return nil
}

func main() {
	center := new(method.ServerCenter)

	centerInstance := &ServerCenterInstance{name:ServerCenterName}
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

