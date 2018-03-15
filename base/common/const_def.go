package common

const(
	MethodServerCenterRegister = "ServerCenter.Register"
	MethodServerCenterDispatch = "ServerCenter.Dispatch"
	MethodServerNodeCall       = "ServerNode.Call"
)

type ModuleRegisterData struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type ModuleDispatchData struct{
	Name string `json:"name"`
	Params string `json:"params"`
}